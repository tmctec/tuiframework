# TUI Framework Example

`example-demo` is one program that shows every feature of the framework. It
is a `package main` under this directory and builds into `./bin` at the
repository root.

| Root tab | Contents |
|---|---|
| Showcase | Welcome, Themes gallery, Custom Theme (Cyberpunk palette with a gradient border), Overlays (z-index stacking), Terminal (what the terminal reported, progress bar, clipboard) |
| Catalog | Electronics (Phones, PCs, and a Deals page beside them), Fashion, Home (Furniture, with Kitchen and Garden beside it) |

A header band shows the app name with the current theme and nav bar status.
`s` opens the Settings popup over any page: theme, nav bar position (top or
sidebar), order (below or above the header), mode (Compact, Expanded, Path
or Columns), autohide and its timeout, tab-nav timeout, header and footer
height (0 to 5 lines; 0 hides the band) and border (none, single line,
double line or fill), nav bar and content borders, key glow, terminal colors
(the terminal's own background and text color take the theme's), and
following the terminal (the theme switches between the light and the dark
member of its family as the terminal's background says). Esc or `s` closes
it. `?` shows the keys for the current page. On a short terminal the popup
scrolls its rows, and its last line shows which row the cursor is on.

The Terminal page lists what the terminal reported about itself: a light or
dark background, its color profile, whether it tells chords such as
`shift+enter` apart, whether its window has focus, and its name and version.
A value reads "not reported" when the terminal did not answer. Below that it
lists what the demo declares to the terminal. `p` runs the terminal's own
progress bar from 0 to 100, where the terminal has one, and `c` copies the
theme name to the clipboard.

Settings is a pure reducer over a snapshot of `tabapp.Model` that emits
`tabapp` messages; `settings_test.go` covers it. The messages that exist
only for the demo, such as the one that cycles sample notifications, are
defined in `messages.go` and handled in `handler.go`: the pattern an
application follows for its own messages.

```bash
make demo
go run ./examples/example-demo   # without make
```

## Building

```bash
make build        # the demo into ./bin
go build -o bin/example-demo ./examples/example-demo
make smoke        # run it in a pseudo-terminal, quit with q, check it exits clean
```

The demo needs a terminal. When `make demo` runs without one on macOS, as it
does from an editor or an assistant's shell, it opens the demo in a new
iTerm2 or Terminal window. Set `NO_WINDOW=1` to get an error instead.
`make smoke` is the check that runs without a terminal.
