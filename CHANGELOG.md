# Changelog

All notable changes to this project are recorded here, newest first. The
format follows [Keep a Changelog](https://keepachangelog.com).

## [0.1.0] - 2026-11-11

The first tagged release.

### Changed
- The module is `github.com/tmctec/tuiframework`.
- The library runs on `charm.land/bubbletea/v2`, `charm.land/lipgloss/v2` and `github.com/tmctec/stickers/v2` v2.0.0, and needs Go 1.26.
- `TabApp.View` returns a `tea.View`; `TabApp.Render` returns the frame as text.
- Keys arrive as `tea.KeyPressMsg` and are named as Bubble Tea v2 names them on every binding call; the space bar is `space`, and a binding written `" "` means the same key.
- Colors in `styles.Palette`, `styles.Theme`, the overlay setters, `overlays.Canvas` and `ToggleModalMsg.BorderColor` are `image/color.Color` values; an unset color is nil.
- An overlay's configured size is its whole footprint, border included, as Lip Gloss v2 counts it.
- `overlays.Cell` carries `Attrs`, `Underline`, `UnderlineColor` and `Link`; `HasBackground` is a method.
- A frame costs less: at 160x48 the top bar layout takes about 131 µs and 41 KB per frame (v1: 161 µs, 86 KB) and a frame with three overlays about 223 µs and 50 KB (v1: 436 µs, 180 KB).
- `make demo` opens the demo in a new iTerm2 or Terminal window on macOS when no terminal is attached; `NO_WINDOW=1` keeps the error.

### Added
- Terminal features declared from the Model, each with a setter and a message: `SetWindowTitle`, `SetTerminalColors`, `SetReportFocus` with `Model.Focused`, `SetMouseMode`, `SetProgress`. `Model.ColorProfile` and `Model.KeyDisambiguation` record what the terminal reports.
- `SetProgramOptions` runs an app over SSH or without a terminal; `SetInitCmds` adds startup commands.
- `UseAdaptiveTheme(light, dark)` follows the terminal's light or dark background; `Model.DarkBackground`, `Theme.IsDark()`, `styles.Counterpart`.
- `styles.PaletteFrom` and `MustPaletteFrom` build a palette from a background, a text color and an accent; `styles.Contrast` and `styles.Hex`.
- Gradient borders: `Palette.BorderBlend`, `TabApp.SetBorderBlend`, `Overlay.SetBorderBlend`.
- `styles.Charm`, a dark theme in Charm's charmtone colors; 22 built-in themes.
- `tabapp.KeyPress(name)` builds the key press message for a key name.
- Demo: a Terminal page, Settings rows for terminal colors and following the terminal's background, a scrolling Settings list on short terminals, a gradient border on the Cyberpunk theme.
- `make smoke`, `make test-eastasian` and `make release-verify`.

### Fixed
- Italic, underlined, faint, reversed and struck-through text, and hyperlinks, keep their styling in the frame; the compositor kept only colors and bold.

### Removed
- The `termenv` dependency and the v1 Charm modules.

### Security
- `make security` scans the tree and every commit for leaked secrets, and runs govulncheck, gosec and a banned-name check; CI runs it on every push and pull request.
- CI runs with a read-only token and actions pinned to commit SHAs; dependabot tracks Go modules and actions.
- The ANSI encoder writes a single byte only for runes in the ASCII range.
- SECURITY.md gives the private route for reporting a vulnerability.
