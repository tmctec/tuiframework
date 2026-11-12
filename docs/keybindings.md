# Key bindings

Keys dispatch messages. A binding is a key name, a label for the footer, and
the message to deliver when the key is pressed. There are no callbacks: the
message reaches `TabApp.Update`, which applies library messages itself and
hands application messages to the custom message handler.

## Key names

A key is named as Bubble Tea reports the press:

| Key | Name |
|---|---|
| A character | `"a"`, `"A"`, `"?"`, `"1"` |
| The space bar | `"space"` |
| Named keys | `"enter"`, `"tab"`, `"esc"`, `"backspace"`, `"up"`, `"down"`, `"left"`, `"right"`, `"home"`, `"end"`, `"pgup"`, `"pgdown"`, `"delete"`, `"f1"` to `"f12"` |
| With modifiers | `"ctrl+s"`, `"shift+tab"`, `"alt+x"`, `"ctrl+shift+a"` |

A binding written `" "` is the space bar: every binding call stores it as
`"space"`, and the footer shows `space`.

Only presses dispatch. A paste arrives as one `tea.PasteMsg` at the custom
message handler and never as keys, so pasted text cannot trigger a binding.
Mouse events and key releases reach the handler the same way.

Some chords exist only on terminals that tell them apart from older keys:
`"shift+enter"` from `"enter"`, `"ctrl+h"` from `"backspace"`, `"ctrl+m"`
from `"enter"`. `Model.KeyDisambiguation` is true once the terminal reports
that it does. Bind the chord and a fallback key to the same message, and
show the chord in the interface only when `KeyDisambiguation` is true:

```go
p.BindKey("shift+enter", "", NewLineMsg{}).
	BindKey("ctrl+j", "New line", NewLineMsg{})
```

`tabapp.KeyPress(name)` builds the key press message for a name. It drives
an application from a test or from code:

```go
app.Update(tabapp.KeyPress("s"))
app.Update(tabapp.KeyPress("shift+tab"))
```

## Universal bindings

A universal binding works on every page.

```go
app.MustBindKey("s", "Settings", ToggleSettingsMsg{})

if err := app.BindKey("ctrl+c", "Quit", QuitMsg{}); err != nil {
	// "key \"ctrl+c\" is hardcoded and cannot be overridden"
}

app.UnbindKey("q")          // let pages use q; chains like the setters
app.IsUniversalKey("s")     // true
app.IsHardcoded("ctrl+c")   // true
app.UniversalBindings()     // a copy of the bindings by key
```

`New` binds `q` to `QuitMsg`, and `tab` and `shift+tab` to the
tab-navigation keys. `Ctrl+C` and `Ctrl+Q` are hardcoded: they always quit
and no binding can take them.

The labels appear in the footer. An empty label hides the key from the
footer while the binding still works.

## Page bindings

A page binding works while that page is active.

```go
p := page.New("Orders").
	BindKey("up", "", MoveCursorMsg{DY: -1}).
	BindKey("down", "", MoveCursorMsg{DY: 1}).
	BindKey("r", "Refresh", RefreshMsg{})

p.UnbindKey("r")
p.HasKey("up")   // true
p.Bindings()     // map[string]page.KeyBinding, read-only
```

`AddPage`, `AddSection` and `SetRoot` validate every page's keys against the
universal and hardcoded bindings and return an error on a collision;
`MustAddPage`, `MustAddSection` and `MustSetRoot` panic instead. Bind
universal keys before adding pages, so the check sees them.

Bubble Tea delivers `esc` on every exit from tab navigation, so a page
should expect `esc` to arrive without a preceding key of its own.

## Dispatch order

1. Hardcoded `Ctrl+C` and `Ctrl+Q`.
2. Overlays, topmost first: exit key, registered keys, then a blocking
   overlay swallows the key.
3. The tab-navigation keys (`Tab`, `Shift+Tab`; change them with
   `SetTabNavKeyForward` and `SetTabNavKeyReverse`).
4. While tab navigation is active: the arrows, `Enter`, `Backspace`; `Esc`
   ends navigation and falls through.
5. Universal bindings.
6. `Ctrl+←` / `Ctrl+→` in the Expanded layout: switch sections.
7. The active page's bindings.
8. Digits `1`-`9`: the page at that position on the current level.

## Footer

The footer lists the keys of the active page: universal bindings first,
then page bindings, both with their labels, then the entries of
`page.SetFooter`, which override labels or add keys the framework does not
know about (`"←→↑↓": "Browse"`). Keys in the order `↑/↓`, `←/→`, `Enter`,
`Space`, `Tab`, `Esc`, `q` come first; the rest follow sorted. With key glow
on, the key just pressed renders in the glow color for
`SetKeyGlowDuration` (300ms by default).

```go
app.SetFooterSeparator(" │ ").
	SetKeyGlowEnabled(true).
	SetKeyGlowDuration(200 * time.Millisecond)
```

## Overlay keys

An overlay has an exit key and registered keys of its own:

```go
m := overlays.NewModal().
	SetTitle("Delete?").
	SetExitKey("esc").
	RegisterKey("y", ConfirmDeleteMsg{}).
	RegisterKey("n", tabapp.DismissOverlayMsg{})
```

`ShowOverlay` rejects an overlay without an exit key (`overlays.ErrNoExitKey`)
and one whose keys collide with the hardcoded quit keys
(`overlays.ErrKeyConflict`).
