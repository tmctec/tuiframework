# Contributing

## Setup

```bash
make prereqs   # checks Go and golangci-lint
make check     # fmt, vet, test
make lint      # golangci-lint with .golangci.yml
make demo      # run the combined demo
make security  # secrets, vulnerabilities, gosec and banned names
```

CI runs `gofmt`, `go vet`, `golangci-lint`, `go test -race`, `make build` and
`make security` on every push and pull request; run `make lint` and
`make test-race` before opening one.

`make security` reads the names that must not appear in the repository as a
regular expression from `BANNED_NAMES`, or from
`~/.config/tuiframework/banned-names`, one name per line. Without either it
skips that check. CI reads `BANNED_NAMES` from a repository secret.

## Conventions

**Fluent builders.** Struct fields are private. Setters are `SetX(v) *T` and
return the receiver; getters are `X()`. A setter that can fail returns an
error and has a `MustX` twin that panics. The exceptions are the types that
are plain data and therefore have exported fields: `tabapp.Model`, which
views read; the message types; `styles.Palette`, `styles.Theme` and
`styles.NotificationStyle`; and `overlays.Cell`.

**Elm pattern.** State lives in `tabapp.Model` (library settings) and
`Model.Custom` (application state, which the library never touches). Keys bind to messages, never to
callbacks. `updateModel` in `tabapp/update.go` is the only place library
state changes in response to a message; the setters are the only place it
changes in response to configuration. Views read the Model and build a
flexbox; they do not mutate anything. New settings get a Model field, a
setter that writes it, a getter that reads it, and a message that toggles or
sets it. See [docs/architecture.md](docs/architecture.md).

**One source of truth.** Do not mirror a Model field into another struct.
Every change to `navPath` goes through `selectAtLevel` or ends with
`syncActivePage`, so `Model.ActivePage` is always right.

**Files by concern.** `tabapp` is split by responsibility (see the package
layout in docs/architecture.md). Put new code in the file whose concern it
belongs to; add a file rather than growing one past its concern.

**Documentation describes the present.** Code comments, docs and READMEs
state what the code does, in active voice, as if it had always been that
way. Change history belongs in commit messages. Every exported identifier
has a doc comment; package overviews live in `doc.go`.

**Tests.** Plain `testing`, table-driven where natural, behavior-level:
drive the app with `tabapp.KeyPress("s")` and `tea.WindowSizeMsg`, and
assert on `Render()`, `View()` or the Model. Lip Gloss styles always render
in full color, so a test sees the escape sequences an application emits.
To test what the app declares to the terminal, run the real program without
one: give `SetProgramOptions` a pipe with `tea.WithInput`, a buffer with
`tea.WithOutput`, `tea.WithWindowSize` and `tea.WithColorProfile`, type keys
into the pipe and read the stream; `runHeadless` in `tabapp/terminal_test.go`
does this. Keep the race detector happy.

**Performance.** `Render` runs on every key press and timer tick. It draws
every band and overlay as a layer on one `overlays.Canvas`, which keeps its
cells between frames; avoid per-frame allocations that scale with the screen
area, and extra passes over rendered text. A Lip Gloss style with a width
wraps and pads its text in several passes, so a band without a border goes
to the canvas as plain lines, and an overlay with fixed content renders once
per size and per change. `make bench` reports the cost of a frame with and
without overlays. On an Apple M2 Pro at 160×48 a frame takes about 130 µs
and 41 KB with the bar on top, 150 µs and 45 KB with the side bar, and
230 µs and 50 KB with three overlays. Keep a change within a quarter of
those.

**Text width.** Measure rendered text with `lipgloss.Width`, or let the
canvas place it. Both count grapheme clusters the same way; a rune-by-rune
measure breaks on combining marks and joined emoji. The locale does not
change a width. `RUNEWIDTH_EASTASIAN=1` makes ambiguous characters such as
box-drawing lines two cells wide in Lip Gloss, Bubble Tea and the canvas
alike. Lip Gloss still sizes a border as one cell there, so a test that
asserts the exact size of a bordered box calls `needsNarrowBorders` and
skips in that mode. Run that mode with `-count=1`; the test cache does not
see the variable.

## Commit messages

A short imperative subject, a body that says what changed and why. The body
is where change history lives; the documentation is not.

## Security issues

Report a vulnerability privately, as [SECURITY.md](SECURITY.md) describes,
not in a public issue or pull request.
