# tuiframework

[![Go Reference](https://pkg.go.dev/badge/github.com/tmctec/tuiframework.svg)](https://pkg.go.dev/github.com/tmctec/tuiframework)
[![CI](https://github.com/tmctec/tuiframework/actions/workflows/ci.yml/badge.svg)](https://github.com/tmctec/tuiframework/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A framework for multi-page terminal applications on
[Bubble Tea](https://charm.land/bubbletea) v2 and Lip Gloss v2. You write pages;
it provides the navigation tree and its nav bar, the header and footer, key
dispatch, themes, and stacked overlays.

- **Pages in a tree.** Sections nest to any depth and may hold pages beside
  subsections. The nav bar has four layouts (Compact, Expanded, Path,
  Columns) and sits on top or on the left, with optional autohide.
- **Keys dispatch messages.** Bind a key to a message on the app or on a
  page; the framework delivers it to your handler. Conflicts with global keys
  are reported when a page is added.
- **Themes.** Twenty-two built-in themes, or a palette of your own that
  cascades into every style, or one derived from three colors. Switch at
  runtime, or let the theme follow the terminal's light or dark background.
- **Overlays.** Notifications with eight levels and auto-dismiss, modal
  dialogs that follow the theme, and a z-indexed stack of either.
- **Terminal features.** The window title, the terminal's own colors, its
  progress bar, focus and mouse reports are settings in the Model. The same
  app runs in a terminal, over SSH, or headless in a test.
- **One source of truth.** State lives in a Model your views read, and keys
  dispatch messages instead of callbacks. See
  [docs/architecture.md](docs/architecture.md).

## Install

```bash
go get github.com/tmctec/tuiframework@latest
```

Go 1.26 or later. Layouts use [stickers/flexbox](https://github.com/tmctec/stickers).

## Quick start

```go
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// GreetMsg is what the g key dispatches; the handler below applies it.
type GreetMsg struct{}

func main() {
	app := tabapp.New()
	app.UseTheme(styles.Nord).SetNavBarAutoHide(false)

	home := page.New("Home").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model) // views read the Model, nothing else
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetStyle(m.Theme.TextHeading)
			cell.SetContent(fmt.Sprintf("Hello, %s", m.GetCustomString("name")))
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		}).
		BindKey("g", "Greet", GreetMsg{})

	about := page.New("About").
		SetContent(func(fb *flexbox.FlexBox, _ interface{}, _, _ int) {
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetContent("Tab moves between pages. q quits.")
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		})

	app.MustAddPage(home).MustAddPage(about)
	app.Model().SetCustom("name", "world")

	// Application messages arrive here; library messages are handled before.
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		if _, ok := msg.(GreetMsg); ok {
			app.Model().SetCustom("name", "friend")
			return true, nil
		}
		return false, nil
	})

	if err := app.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
```

Sections group pages into a tree:

```go
catalog := tabapp.NewSection("Catalog").
	AddSection(tabapp.NewSection("Phones").AddPages(browse, deals)).
	AddPage(overview) // a page beside a subsection
app.MustAddSection(catalog)
```

## Try the demo

```bash
make demo
```

The demo opens on a Welcome page. Press `s` for the Settings popup, where
the theme, the nav bar's position, layout, autohide and border, the header
and footer height and border, key glow and the terminal's colors change
live. The Terminal page shows what your terminal reported about itself. The
Catalog section is a mixed-depth navigation tree; `?` shows the keys for the
current page.
[examples/README.md](examples/README.md) describes it.

## Packages

| Package | Purpose | Guide |
|---|---|---|
| `tabapp` | The application: layout, navigation tree and nav bar, key dispatch, overlay stack, Bubble Tea `Model` | [docs/architecture.md](docs/architecture.md), [docs/keybindings.md](docs/keybindings.md) |
| `page` | A page: name, content function, key bindings, footer labels | [docs/keybindings.md](docs/keybindings.md) |
| `overlays` | Overlay, Notification and Modal, plus the Canvas that composites a frame | [docs/overlays.md](docs/overlays.md) |
| `styles` | Palette, Theme, the 22 built-in themes | [docs/theming.md](docs/theming.md) |

API reference: [pkg.go.dev/github.com/tmctec/tuiframework](https://pkg.go.dev/github.com/tmctec/tuiframework).

## Navigation

| Key | Action |
|---|---|
| `Tab`, `Shift+Tab` | Start tab navigation and move the highlight along the bar |
| `←` `→` | Compact bar: move along the strip, crumbs included. Expanded bar: move within the level |
| `↑` `↓` | Expanded bar: change level. Path, Columns and the side bar: move within the level |
| `Enter` | Open the highlighted crumb, section or page |
| `Backspace` | Go back a level |
| `Esc` | Leave tab navigation |
| `1`-`9` | Jump to a page at the current level |
| `q`, `Ctrl+C` | Quit (`q` can be unbound; `Ctrl+C` and `Ctrl+Q` cannot) |

Keys are named as Bubble Tea reports them: `"s"`, `"space"`, `"enter"`,
`"shift+tab"`, `"ctrl+c"`. See [docs/keybindings.md](docs/keybindings.md).

`SetNavBarMode` picks the layout: `NavBarCompact` is one strip with crumbs,
`NavBarExpanded` is a row per level, `NavBarPath` is the route with a
drop-down of siblings while navigating, and `NavBarColumns` is a column per
level in the style of a file browser. Path and Columns size the bar
themselves; `SetNavBarHeight` applies to the other two. `SetNavBarPosition`
moves the bar to the left as a tree. With autohide on, the bar appears while
navigating and hides after the timeout.

## Header and footer

```go
app.SetHeader(renderHeader).                 // a function that draws the band
	SetHeaderHeight(1).
	SetHeaderFrame(tabapp.FrameFill).        // a solid bar, text reversed
	SetFooterHeight(2).
	SetFooterFrame(tabapp.FrameSingle)       // a single-line border
```

The header is a band your function draws from the Model; the footer lists
the active page's keys. Each has a height in lines and a frame:
`FrameNone`, `FrameSingle`, `FrameDouble`, `FrameTheme` for the theme's
border shape, or `FrameFill`. A height of 0 hides the band. See
[docs/architecture.md](docs/architecture.md#header-and-footer).

## Theming

```go
app.UseTheme(styles.CatMocha)                 // a built-in theme
app.UseTheme(styles.NewTheme(myPalette))      // your palette, cascaded into a theme
app.MustBindKey("t", "Theme", tabapp.CycleThemeMsg{Direction: 1, Themes: styles.AllThemes})
```

A palette needs its base colors and a border; hero, semantic and UI colors
derive from them unless set. Views read `m.Theme` for ready-made styles such
as `TextHeading`, `CardStyle` and `StatusError`.

```go
app.UseAdaptiveTheme(styles.CatLatte, styles.CatMocha) // light or dark, as the terminal is

c := lipgloss.Color
p, err := styles.PaletteFrom(c("#101418"), c("#d7dde4"), c("#7aa2f7")) // background, text, accent
```

`UseAdaptiveTheme` follows the terminal's background. `PaletteFrom` builds a
whole palette from three colors and refuses a combination that would not be
readable. See [docs/theming.md](docs/theming.md).

## The terminal

```go
app.SetWindowTitle("Inventory").
	SetTerminalColors(true).                     // the terminal's own colors follow the theme
	SetProgress(tea.ProgressBarDefault, 40).     // the terminal's progress bar at 40%
	SetReportFocus(true)                         // Model.Focused follows the window
```

Each is a setting in the Model with a message of its own, such as
`SetProgressMsg`, so a key binding can change it. The Model also holds what
the terminal reported: `DarkBackground`, `ColorProfile`, `Focused` and
`KeyDisambiguation`.

```go
app.SetProgramOptions(tea.WithInput(in), tea.WithOutput(out), tea.WithWindowSize(80, 24))
```

`SetProgramOptions` passes options to the Bubble Tea program, which is how an
app runs over SSH or without a terminal. See
[docs/architecture.md](docs/architecture.md#terminal-features).

## Overlays

```go
cmd := app.ShowNotification(overlays.Success("Saved")) // dismisses itself after 3s
id, err := app.ShowModal(overlays.NewModal().SetTitle("Confirm").SetExitKey("esc"))
```

Notifications and modals stack by z-index; keys go to the topmost overlay
first, and a blocking overlay stops them there. Messages such as
`ShowNotificationMsg` and `ToggleModalMsg` do the same from key bindings. See [docs/overlays.md](docs/overlays.md).

## Development

```bash
make help        # every target
make check       # fmt, vet, test
make lint        # golangci-lint with .golangci.yml
make test-race   # tests under the race detector
make bench       # frame benchmarks
make build       # the demo into ./bin
make smoke       # run the demo in a pseudo-terminal and check it exits clean
make security    # secrets in the tree and every commit, vulnerabilities, gosec, banned names
make release-verify VERSION=0.1.0  # fetch, build and check a released version from outside the repo
```

CI runs format, vet, lint, race tests, the build and `make security` on every
push and pull request. Conventions for contributors are in [CONTRIBUTING.md](CONTRIBUTING.md).

## Security

Report a vulnerability privately, as [SECURITY.md](SECURITY.md) describes.

## License

[MIT](LICENSE)
