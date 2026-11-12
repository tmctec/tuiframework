# Overlays

Overlays float above the page: notifications for transient messages, modals
for dialogs, and the base `Overlay` for anything else. They stack by
z-index, receive keys before the page, and are composited over the frame
with their own background so the content below never bleeds through.

```go
import "github.com/tmctec/tuiframework/overlays"

cmd := app.ShowNotification(overlays.Success("Saved"))   // auto-dismiss after 3s
id, err := app.ShowModal(overlays.NewModal().SetTitle("Confirm").SetExitKey("esc"))
id, err = app.ShowOverlay(overlays.New().SetContent("...").SetExitKey("esc"))
```

## Overlay

`overlays.New()` creates the base overlay: 80% of the screen, centered,
blocking, with a rounded border and padding, and no exit key. An exit key is
required; `Validate` reports `ErrNoExitKey` and `ShowOverlay` refuses the
overlay without one.

| Method | Effect |
|---|---|
| `SetContent(string)` | Text content |
| `SetFlexbox(*flexbox.FlexBox)` | A flexbox resized to the content area on each render |
| `SetContentFunc(func(*flexbox.FlexBox))` | Builds a fresh flexbox on every render, like a page |
| `SetTitle`, `SetTitleStyle` | A title line above the content; it takes two rows |
| `SetSize(pct)` | Width and height as one percentage of the screen |
| `SetWidthPct`, `SetHeightPct` | Separate percentages |
| `SetWidth`, `SetHeight` | Fixed cells and lines; fixed beats percentage |
| `SetPosition(Position)` | One of nine anchors: `PositionTopLeft` ... `PositionBottomRight`; `PositionCenter` is the default |
| `SetBlocking(bool)` | Whether unknown keys stop here instead of reaching the page |
| `SetExitKey(key)` | The key that hides the overlay |
| `RegisterKey(key, msg)`, `RegisterKeys(map)` | Keys that dispatch messages while the overlay is shown |
| `SetZIndex(int)` | Stacking order; higher renders on top |
| `SetStyle`, `SetBorder`, `SetBorderForeground`, `SetBorderBackground`, `SetBackground`, `SetForeground`, `SetPadding` | Frame styling |

`CalculateDimensions(w, h)` returns the footprint in cells for a screen
size, never larger than the screen. `Render(w, h)` draws the overlay; the
frame, padding and title come out of the footprint, the rest is content.

## Notification

A notification is an overlay with a level, a timeout and a message. The
level picks its colors and icon; the message is prefixed with the icon.

```go
overlays.Info("Syncing")                   // also Note, Tip, Important,
overlays.Warning("Disk at 90%").            // Success, Warning, Caution, Error
	SetTimeout(0).                          // stays until its exit key
	SetPosition(overlays.PositionBottomRight).
	SetZIndex(20)
```

`app.ShowNotification(n)` shows it and returns a command that dismisses it
after `n.Timeout()`; return that command from your update function. A
timeout of 0 returns nil and the notification stays until its exit key.

### Defaults: NotificationConfig

`overlays.NotificationConfig` is the `*NotificationSettings` every new
notification starts from: top right, 50 cells by 3 lines, non-blocking,
`esc` dismisses, three-second timeout, and the colors, icons and border of
the Dracula theme. `TabApp.UseTheme` calls `NotificationConfig.ApplyTheme`
so new notifications match the active theme.

```go
overlays.NotificationConfig.
	SetPosition(overlays.PositionBottomCenter).
	SetWidthPct(60).SetWidth(0).           // width 0 makes the percentage apply
	SetTimeout(5 * time.Second).
	SetLevelStyle(overlays.LevelTip, styles.NotificationStyle{
		Foreground: lipgloss.Color("#50FA7B"), Background: lipgloss.Color("#282A36"),
		Border: lipgloss.Color("#50FA7B"), Icon: "[tip]",
	})
```

`LevelStyle(level)` returns a level's `styles.NotificationStyle`. The
settings are process-wide; a theme change restyles notifications created
afterwards, not those already on screen.

## Modal

A modal is a blocking overlay styled from the theme's `Modal*` fields:
background, border, border color and title style. `NewModal()` applies
`styles.Dracula`; `ApplyTheme(t)` applies another theme, and a modal shown
through `TabApp.ShowModal` is re-themed when the app's theme changes. A
modal given its own `SetStyle` stops following the theme (`FollowsTheme`
reports which).

```go
m := overlays.NewModal().
	SetTitle("Delete order 42?").
	SetSize(50).
	SetContent(func(fb *flexbox.FlexBox) {
		row := fb.NewRow()
		cell := flexbox.NewCell(1, 1)
		cell.SetContent("This cannot be undone.\n\n[y] delete   [n] keep")
		row.AddCells(cell)
		fb.AddRows([]*flexbox.Row{row})
	}).
	SetExitKey("esc").
	RegisterKey("y", DeleteOrderMsg{ID: 42}).
	RegisterKey("n", tabapp.DismissOverlayMsg{})

id, err := app.ShowModal(m)
```

## The stack

`TabApp` keeps overlays sorted by z-index, lowest first; among equal
z-indices the later one is on top.

| Method | Effect |
|---|---|
| `ShowOverlay(o) (int64, error)` | Validates and shows; returns the overlay's id |
| `MustShowOverlay(o) int64` | Panics on error |
| `ShowNotification(n) tea.Cmd` | Shows and returns the auto-dismiss command |
| `ShowModal(m) (int64, error)` | Shows and keeps the modal in step with the theme |
| `HideOverlay()` | Hides the topmost |
| `HideOverlayByID(id) bool` | Hides one; false when the id is not on the stack |
| `HideAllOverlays()` | Clears the stack |
| `BringToFront(id)` | Moves one to the top with z-index max+1 |
| `TopOverlay() *Overlay` | The topmost, or nil |
| `HasOverlay() bool`, `OverlayCount() int` | Stack state; `Model.OverlayCount` mirrors the count for views |

### Keys with overlays

A key press checks the overlays from the top down. The exit key hides that
overlay; a registered key dispatches its message; then a blocking overlay
stops the key and a non-blocking one lets it continue to the overlay below
and finally to the page. `ShowOverlay` rejects an overlay whose keys collide
with `Ctrl+C` or `Ctrl+Q`; the error wraps `overlays.ErrKeyConflict`.

### Messages

Key bindings drive the stack through messages:

| Message | Effect |
|---|---|
| `ShowNotificationMsg{Level, Message, ZIndex}` | Show a notification of that level; `Level` is an `overlays.NotificationLevel` and its zero value is `LevelInfo` |
| `ShowModalMsg{Title, Content, ZIndex, SizePct}` | Open a blocking themed modal with text content, closed by `esc`. A `SizePct` of 0 keeps the default size |
| `ToggleModalMsg{Title, Content, ZIndex, SizePct, BorderColor}` | Open a themed modal that lets keys through to the page, or close it if a modal with that title is open. `BorderColor` is a `color.Color`; nil keeps the theme's border. `Model.OpenModals` lists the open ones by title |
| `BringOverlayToFrontMsg{Name}` | Bring the `ToggleModalMsg` modal with that title to the front |
| `ClearAllOverlaysMsg{}` | Hide every overlay |
| `DismissOverlayMsg{}` | Hide the topmost overlay |

```go
app.MustBindKey("ctrl+s", "Save", tabapp.ShowNotificationMsg{
	Level: overlays.LevelSuccess, Message: "Saved",
})
```

Each message does one thing. To do two in a fixed order, such as clearing
the stack and then showing a notification, emit both from the custom message
handler with `tea.Sequence`; the demo's `handler.go` does.

## Compositing

`View` draws the frame on an `overlays.Canvas`, a grid of cells. Every band
and every overlay is a layer drawn in its own rectangle. A layer replaces
what is under it and gives the cells its text leaves without a background a
fill color, so nothing shows through. The frame becomes text once, however
many layers it took.

The canvas is exported for applications that draw their own layers:

```go
c := overlays.NewCanvas().Reset(width, height, theme.Palette.Bg0)
c.Draw(pageLines, 0, 0, width, height, theme.Palette.Bg0)

lines, w := overlays.GetLines(popup.Render(width, height))
x, y := overlays.CalculatePosition(popup.Position(), w, len(lines), width, height)
c.Draw(lines, x, y, w, len(lines), theme.ModalBackground)

frame := c.String()
```

| Function | Effect |
|---|---|
| `NewCanvas()` | An empty canvas. It keeps its cells between frames, so a frame allocates only its text |
| `Canvas.Reset(w, h, color)` | Starts a frame of that size, every cell a space on the color |
| `Canvas.Draw(lines, x, y, w, h, color)` | Draws a layer in the rectangle. Text outside it is dropped, short lines are padded, and cells without a background take the color |
| `Canvas.String()` | The frame, one line per row, every line exactly as wide as the canvas |
| `ParseCells(line, width)` | Reads a line into `Cell`s, one per screen column. A cell carries its colors, its attributes (`Attrs`: bold, faint, italic, blink, reverse, conceal, strikethrough), its underline shape and color, and its hyperlink |
| `FillBackground(lines, width, color)` | Pads lines to the width and gives every cell without a background that color |
| `GetLines(block)` | Splits rendered text into lines and reports the width of the widest |
| `CalculatePosition(pos, w, h, screenW, screenH)` | An anchor's offset |

Text is measured per grapheme cluster, the way Lip Gloss measures it. A wide
character takes two cells, and combining marks and joined emoji stay with
their character. The locale does not change a width. Setting
`RUNEWIDTH_EASTASIAN=1` makes ambiguous characters such as box-drawing lines
two cells wide, here as in Lip Gloss and Bubble Tea.

Everything Lip Gloss styles survives compositing: italic and the other
attributes, curly, dotted and colored underlines, gradient borders, and
hyperlinks made with `lipgloss.NewStyle().Hyperlink(url)`. A layer drawn across a hyperlink
ends the link at its edge, and the link resumes on the far side.
