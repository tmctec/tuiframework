package overlays

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/stickers/v2/flexbox"
)

func TestNewModal(t *testing.T) {
	m := NewModal()

	if !m.IsBlocking() {
		t.Error("Modal should be blocking by default")
	}
	if m.Size() != 80 {
		t.Errorf("Size() = %d, want 80", m.Size())
	}
	if m.Position() != PositionCenter {
		t.Errorf("Position() = %s, want CC", m.Position())
	}
}

func TestModalSetTitle(t *testing.T) {
	m := NewModal().SetTitle("My Title")

	if m.Title() != "My Title" {
		t.Errorf("Title() = %q, want %q", m.Title(), "My Title")
	}
}

func TestModalSetContent(t *testing.T) {
	contentCalled := false
	m := NewModal().SetContent(func(fb *flexbox.FlexBox) {
		contentCalled = true
		// Add some content to flexbox
		row := fb.NewRow()
		cell := flexbox.NewCell(1, 1)
		cell.SetContent("Test")
		row.AddCells(cell)
		fb.AddRows([]*flexbox.Row{row})
	})

	if !contentCalled {
		t.Error("Content function should be called immediately")
	}
	if m.Flexbox() == nil {
		t.Error("Flexbox() should not be nil after SetContent")
	}
}

func TestModalChaining(t *testing.T) {
	m := NewModal()
	result := m.
		SetTitle("Title").
		SetSize(60).
		SetPosition(PositionTopCenter).
		SetExitKey("esc").
		SetBlocking(true)

	if result != m {
		t.Error("Chaining broken - setter did not return self")
	}
}

// testModalKeyMsg is a message for testing modal key registration
type testModalKeyMsg struct{ key string }

func TestModalKeyRegistration(t *testing.T) {
	// Register messages to dispatch when keys are pressed (Elm pattern)
	m := NewModal().
		RegisterKey("y", testModalKeyMsg{key: "y"}).
		RegisterKey("n", testModalKeyMsg{key: "n"})

	if !m.HasKey("y") {
		t.Error("HasKey('y') = false, want true")
	}
	if !m.HasKey("n") {
		t.Error("HasKey('n') = false, want true")
	}

	// Verify the stored message is correct
	if msg, ok := m.Keys()["y"].(testModalKeyMsg); !ok || msg.key != "y" {
		t.Error("Keys()['y'] should contain testModalKeyMsg{key: 'y'}")
	}
}

func TestModalRegisterMultipleKeys(t *testing.T) {
	// Register multiple messages at once (Elm pattern)
	m := NewModal().RegisterKeys(map[string]tea.Msg{
		"y": testModalKeyMsg{key: "y"},
		"n": testModalKeyMsg{key: "n"},
	})

	if !m.HasKey("y") || !m.HasKey("n") {
		t.Error("RegisterKeys did not register all keys")
	}
}

func TestModalInheritsOverlay(t *testing.T) {
	m := NewModal().
		SetExitKey("esc").
		SetPosition(PositionBottomCenter).
		SetSize(70)

	if m.ExitKey() != "esc" {
		t.Errorf("ExitKey() = %q, want %q", m.ExitKey(), "esc")
	}
	if m.Position() != PositionBottomCenter {
		t.Errorf("Position() = %s, want CB", m.Position())
	}
	if m.Size() != 70 {
		t.Errorf("Size() = %d, want 70", m.Size())
	}
}

func TestModalValidate(t *testing.T) {
	m := NewModal()

	err := m.Validate()
	if err != ErrNoExitKey {
		t.Errorf("Validate() = %v, want ErrNoExitKey", err)
	}

	m.SetExitKey("esc")
	err = m.Validate()
	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestModalRender(t *testing.T) {
	m := NewModal().
		SetTitle("Test Modal").
		SetExitKey("esc")

	rendered := m.Render(100, 50)
	if rendered == "" {
		t.Error("Render() returned empty string")
	}
}

func TestModalRenderWithContent(t *testing.T) {
	m := NewModal().
		SetTitle("Test").
		SetContent(func(fb *flexbox.FlexBox) {
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetContent("Content here")
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		}).
		SetExitKey("esc")

	rendered := m.Render(100, 50)
	if rendered == "" {
		t.Error("Render() returned empty string")
	}
}

func TestModalSetBlocking(t *testing.T) {
	m := NewModal().SetBlocking(false)

	if m.IsBlocking() {
		t.Error("SetBlocking(false) should make modal non-blocking")
	}
}
