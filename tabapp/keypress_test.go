package tabapp

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

// TestKeyPressRoundTrips checks that the message KeyPress builds reports the
// name it was built from, which is what the dispatcher matches bindings on.
func TestKeyPressRoundTrips(t *testing.T) {
	names := []string{
		"a", "s", "?", "1", "+", "-", "/", "é",
		"A", "Z",
		"space", "enter", "tab", "esc", "backspace",
		"up", "down", "left", "right", "home", "end", "pgup", "pgdown", "insert", "delete",
		"f1", "f5", "f12",
		"shift+tab", "shift+enter", "ctrl+c", "ctrl+q", "ctrl+h", "ctrl+right", "ctrl+left",
		"alt+x", "ctrl+alt+x", "ctrl+shift+a", "super+space", "ctrl++", "meta+m", "hyper+h",
	}
	for _, name := range names {
		if got := KeyPress(name).String(); got != name {
			t.Errorf("KeyPress(%q).String() = %q", name, got)
		}
	}
}

func TestKeyPressDetails(t *testing.T) {
	// The space bar may be written " "; it is the key named space.
	if got := KeyPress(" "); got.String() != "space" || got.Code != tea.KeySpace || got.Text != " " {
		t.Errorf(`KeyPress(" ") = %+v, want the space key with its text`, got)
	}

	// A capital is its key with shift held, as a terminal reports it.
	capital := KeyPress("A")
	if capital.Code != 'a' || capital.ShiftedCode != 'A' || capital.Mod != tea.ModShift || capital.Text != "A" {
		t.Errorf(`KeyPress("A") = %+v, want code a, shifted A, shift held`, capital)
	}

	// A chord carries no text, so it never reads as typed input.
	if chord := KeyPress("ctrl+c"); chord.Text != "" || chord.Code != 'c' || chord.Mod != tea.ModCtrl {
		t.Errorf(`KeyPress("ctrl+c") = %+v, want code c with ctrl and no text`, chord)
	}

	// Several characters that name no key are text typed at once.
	if text := KeyPress("héllo"); text.Code != tea.KeyExtended || text.Text != "héllo" || text.String() != "héllo" {
		t.Errorf(`KeyPress("héllo") = %+v, want extended text`, text)
	}

	// An unknown modifier adds nothing.
	if odd := KeyPress("turbo+x"); odd.Mod != 0 || odd.Code != 'x' {
		t.Errorf(`KeyPress("turbo+x") = %+v, want plain x`, odd)
	}
}

// TestViewDeclaresTheFrame checks View against Render: the same frame, shown
// on the alternate screen.
func TestViewDeclaresTheFrame(t *testing.T) {
	app := rootMixApp()
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	v := app.View()
	if v.Content != app.Render() {
		t.Error("View should carry the frame Render returns")
	}
	if !v.AltScreen {
		t.Error("View should ask for the alternate screen")
	}
}
