# Architecture

tuiframework is four packages. `tabapp` is the application; the others are
the pieces it assembles.

```
tabapp ──► page      pages: name, content function, key bindings, footer labels
   │  └──► overlays  Overlay, Notification, Modal, the Canvas that composites them
   └─────► styles    Palette → Theme; the 22 built-in themes
```

`page` and `overlays` do not import `tabapp`. A page content function
therefore receives the model as `interface{}` and asserts it to
`*tabapp.Model`.

## Elm pattern

The framework follows the Elm architecture: one model, messages that
describe what can happen, one update function that applies them, and views
that read the model.

**Model.** `tabapp.Model` holds every setting a view can read: nav bar
position, mode, autohide and timeouts, the theme, the active page index, the
band heights and frames, the border and key glow switches, the overlay
count and the open modals, what the application declares to the terminal
and what the terminal reported. `Custom` is a map for application state; the
library never reads or writes it. `TabApp.Model()` returns the Model and
every page content function receives it. The Model is the single source of truth: `TabApp`
keeps no second copy of a setting. The setters (`SetNavBarMode`,
`SetShowFooterBorder`, ...) write the Model; the getters read it.

**Messages.** `tabapp/messages.go` defines the library messages: nav bar
toggles and adjustments, band heights and frames, navigation jumps, theme
changes, terminal features, notifications and modals. An application defines its own message types the same way, as plain
structs.

**Update.** `TabApp.Update` resolves a key press to a message (see the
dispatch order below) and hands every message to `updateModel`, the one
place where library state changes. A message `updateModel` does not know
goes to the function given to `SetCustomMessageHandler`, which is the
application's update function. It changes application state in
`Model.Custom` and changes library state by returning commands that emit
library messages.

**View.** `TabApp.Render` draws the frame from the Model and the content
tree, and `TabApp.View` returns it to Bubble Tea with the terminal features
the Model asks for. Page content functions run on every render with a fresh
flexbox, so they stay pure: read the Model, add rows and cells, return.

Timers follow the same loop. Autohide, tab-navigation timeout, key glow and
notification dismissal each arm a `tea.Tick` that carries an id; the
resulting message is applied only when its id is still the current one, so
a re-armed timer silently cancels the stale one.

## Fluent builders

Every configurable type (`TabApp`, `Page`, `Section`, `Overlay`,
`Notification`, `Modal`, `NotificationSettings`) keeps its fields private
and exposes setters that return the receiver, so configuration chains:

```go
app := tabapp.New().
	SetNavBarMode(tabapp.NavBarColumns).
	SetNavBarAutoHide(false).
	UseTheme(styles.Gruvbox)
```

Getters have the field's name without `Set`. A setter that can fail
(`AddPage`, `AddSection`, `SetRoot`, `BindKey`, `ShowOverlay`) returns an
error and has a `Must` twin that panics, for static setup. The types that
are plain data by design have exported fields: `Model`, which views read;
the messages, which are values; and in the other packages `styles.Palette`,
which applications write as a literal, `styles.Theme` and
`styles.NotificationStyle`, which views read, and `overlays.Cell`.

## Key dispatch

`TabApp.Update` resolves a key press in this order and stops at the first
match. A key release, a paste and a mouse event are not key presses: they
go to the custom message handler as the Bubble Tea messages they are.

1. **Hardcoded keys.** `Ctrl+C` and `Ctrl+Q` quit. Nothing overrides them.
2. **Overlays**, topmost first. The overlay's exit key hides it; one of its
   registered keys dispatches that message; a blocking overlay swallows any
   other key.
3. **Tab-navigation keys** (`Tab`, `Shift+Tab` by default) enter navigation
   and move the highlight.
4. **Navigation keys while navigation is active**: the arrows, `Enter` and
   `Backspace`, in the shape of the bar's layout. `Esc` leaves navigation
   and then falls through.
5. **Universal bindings** from `BindKey`.
6. **Ctrl+arrows** switch sections in the Expanded layout.
7. **The active page's bindings** from `page.BindKey`.
8. **Digits 1-9** select a page at the current level.

`AddPage`, `AddSection` and `SetRoot` reject a page whose keys collide with a
universal or hardcoded binding, so conflicts surface at startup. See
[keybindings.md](keybindings.md).

## Navigation tree

`Section` is recursive: it holds child sections and pages, in any mix. The
root is an anonymous section. `navPath` holds one index per level of the
tree's depth, counting a level's sections first and then its pages; once an
index selects a page, the deeper indices are ignored. `NavigateToPath`
jumps to any position and `NavPath` reports the current one. Digit keys,
`SetActive` and the `NavigateToPageMsg` family work on the pages of the
level the active page is on, so a page beside subsections is as reachable
as a leaf page. `Active` and `Model.ActivePage` are the active page's index
among those pages.

## Rendering

`Render` draws the frame on an `overlays.Canvas`, a grid of cells the size of
the terminal. The bands are layers, each drawn in its rectangle: header
(optional), nav bar (unless autohide has hidden it), content, footer.
In the top position the bands stack; in the left position the bar sits
beside the content. The overlays follow in z order. A layer replaces what is
under it and fills the cells its text leaves without a background: the
theme's `Bg0` under the bands, the overlay's own background under an
overlay. A transparent terminal never shows through.

The canvas parses each layer's text into cells that carry their colors,
attributes, underline and hyperlink, and writes the frame once, with escape
sequences only where the style changes. It keeps its cells between frames,
so a frame allocates little more than its own text. Text is measured per
grapheme cluster with the measure Lip Gloss lays text out with, so wide
characters, combining marks and box-drawing lines stay aligned.

`Render` returns the frame as text. `View` wraps it in a `tea.View` and
declares the terminal features the Model asks for; see
[Terminal features](#terminal-features).

`tabapp/bench_test.go` measures `Render`; `make bench` runs it.

### Terminal features

What the application wants from the terminal is state in the Model, and
`View` declares it on every frame through the fields of the `tea.View`.
Nothing is switched on or off by a command.

| Setter | Model field | Message | Effect |
|---|---|---|---|
| `SetWindowTitle(title)` | `WindowTitle` | `SetWindowTitleMsg` | The title of the terminal window |
| `SetTerminalColors(on)` | `TerminalColors` | `ToggleTerminalColorsMsg` | The terminal's own background and foreground take the theme's `Bg0` and `Fg1` |
| `SetReportFocus(on)` | `ReportFocus`, `Focused` | `ToggleReportFocusMsg` | The terminal reports focus changes; `Focused` follows them |
| `SetMouseMode(mode)` | `MouseMode` | `SetMouseModeMsg` | The terminal reports mouse events; the handler receives them |
| `SetProgress(state, value)` | `ProgressState`, `ProgressValue` | `SetProgressMsg` | The terminal's own progress bar |

Bubble Tea puts back what it changed when the program exits.

The Model also records what the terminal reports: `DarkBackground` and
`BackgroundKnown`, `ColorProfile`, and `KeyDisambiguation`. Each report
reaches the application's message handler after the Model has it.

`SetProgramOptions` passes options to the Bubble Tea program, which runs an
application over SSH or without a terminal. `SetInitCmds` adds commands that
run at the start; the answer to a request such as
`tea.RequestTerminalVersion` reaches the message handler.

### Header and footer

The header is a band the application draws: `SetHeader` takes a function
that receives a fresh flexbox and the Model on every render, the way a page
content function does. The footer lists the keys of the active page. Each
band has a height in lines of content and a frame.

| Frame | Band |
|---|---|
| `FrameNone` | Bare (default) |
| `FrameSingle` | Inside a single-line border |
| `FrameDouble` | Inside a double-line border |
| `FrameTheme` | Inside a border of the theme's shape, the palette's `Border` |
| `FrameFill` | A solid bar in the theme's fill colors, with reversed text |

```go
app.SetHeader(renderHeader).
	SetHeaderHeight(1).
	SetHeaderFrame(tabapp.FrameFill).
	SetFooterHeight(2).
	SetFooterFrame(tabapp.FrameSingle)
```

A border adds a row above and below the content lines and a column on each
side. A fill adds no rows and insets the content one column on each side,
so text starts in the same column under every frame but `FrameNone`. A
height of 0 hides the band. `SetShowHeaderBorder` and `SetShowFooterBorder`
are shorthand for `FrameTheme` and `FrameNone`.

`SetHeaderHeightMsg`, `SetFooterHeightMsg`, `SetHeaderFrameMsg` and
`SetFooterFrameMsg` change the bands from a key binding. Views read the
current values from `Model.HeaderHeight`, `Model.FooterHeight`,
`Model.HeaderFrame` and `Model.FooterFrame`.

A filled footer draws its keys and labels in `Theme.FooterFill` and shows
the glowing key as a chip of the content background. A header function
draws a filled header on `Theme.HeaderFill`; cells it leaves without a
background take the fill color.

```go
func renderHeader(fb *flexbox.FlexBox, model interface{}) {
	m := model.(*tabapp.Model)
	base := m.Theme.BaseStyle
	if m.HeaderFrame == tabapp.FrameFill {
		base = m.Theme.HeaderFill
	}
	fb.UseStyle(base)

	row := fb.NewRow().StylePassing(true) // cells inherit the base colors
	title := flexbox.NewCell(1, 1)
	title.SetContent("Orders")
	row.AddCells(title)
	fb.AddRows([]*flexbox.Row{row})
}
```

## Package layout

```
tabapp/
  doc.go         package overview
  tabapp.go      TabApp, New, Update, key dispatch
  model.go       Model
  messages.go    library messages and the internal timer messages
  section.go     Section and Breadcrumb
  content.go     AddPage, AddSection, SetRoot and key validation
  navigation.go  tree queries and position changes
  tabnav.go      tab-navigation mode: entering, leaving, the per-layout keys
  keys.go        universal bindings, footer text, key glow
  keypress.go    KeyPress: the key press message for a key name
  terminal.go    terminal features, adaptive themes, program options
  options.go     setters and getters, UseTheme
  frame.go       Frame: how the header and footer are set off
  view.go        View and Render: the bands and overlays as canvas layers
  navbar.go      the four bar layouts and the side bar
  overlays.go    the overlay stack
  update.go      updateModel
page/            Page and KeyBinding
overlays/        Overlay, Notification, Modal, Position, Canvas and Cell
styles/          Palette, Theme, NewTheme, PaletteFrom, the built-in palettes
internal/keyname the rule every binding call shares for naming a key
examples/        example-demo, one program that shows every feature
```
