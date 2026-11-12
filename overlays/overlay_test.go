package overlays

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestNewOverlay(t *testing.T) {
	o := New()

	if o.Size() != 80 {
		t.Errorf("Size() = %d, want 80", o.Size())
	}
	if o.Position() != PositionCenter {
		t.Errorf("Position() = %s, want CC", o.Position())
	}
	if !o.IsBlocking() {
		t.Error("IsBlocking() = false, want true")
	}
	if o.IsVisible() {
		t.Error("IsVisible() = true, want false")
	}
}

func TestOverlaySetters(t *testing.T) {
	o := New().
		SetContent("test content").
		SetSize(50).
		SetPosition(PositionTopRight).
		SetBlocking(false).
		SetExitKey("esc")

	if o.Content() != "test content" {
		t.Errorf("Content() = %q, want %q", o.Content(), "test content")
	}
	if o.Size() != 50 {
		t.Errorf("Size() = %d, want 50", o.Size())
	}
	if o.Position() != PositionTopRight {
		t.Errorf("Position() = %s, want RT", o.Position())
	}
	if o.IsBlocking() {
		t.Error("IsBlocking() = true, want false")
	}
	if o.ExitKey() != "esc" {
		t.Errorf("ExitKey() = %q, want %q", o.ExitKey(), "esc")
	}
}

func TestOverlayChaining(t *testing.T) {
	// Verify all setters return self for chaining
	o := New()
	result := o.SetContent("x").SetSize(50).SetPosition(PositionBottomLeft).SetBlocking(true).SetExitKey("q")
	if result != o {
		t.Error("Chaining broken - setter did not return self")
	}
}

func TestOverlaySizeClamp(t *testing.T) {
	o := New()

	o.SetSize(0)
	if o.Size() != 1 {
		t.Errorf("Size() = %d after SetSize(0), want 1", o.Size())
	}

	o.SetSize(150)
	if o.Size() != 100 {
		t.Errorf("Size() = %d after SetSize(150), want 100", o.Size())
	}
}

func TestOverlayVisibility(t *testing.T) {
	o := New()

	if o.IsVisible() {
		t.Error("New overlay should not be visible")
	}

	o.Show()
	if !o.IsVisible() {
		t.Error("Show() should make overlay visible")
	}

	o.Hide()
	if o.IsVisible() {
		t.Error("Hide() should make overlay not visible")
	}
}

// testOverlayKeyMsg is a message for testing key registration
type testOverlayKeyMsg struct{ key string }

func TestOverlayKeyRegistration(t *testing.T) {
	// Register a message to dispatch when key is pressed (Elm pattern)
	o := New().RegisterKey("a", testOverlayKeyMsg{key: "a"})

	if !o.HasKey("a") {
		t.Error("HasKey('a') = false, want true")
	}
	if o.HasKey("b") {
		t.Error("HasKey('b') = true, want false")
	}

	// Verify the stored message is correct
	if msg, ok := o.Keys()["a"].(testOverlayKeyMsg); !ok || msg.key != "a" {
		t.Error("Keys()['a'] should contain testOverlayKeyMsg{key: 'a'}")
	}
}

func TestOverlayRegisterMultipleKeys(t *testing.T) {
	// Register multiple messages at once (Elm pattern)
	o := New().RegisterKeys(map[string]tea.Msg{
		"a": testOverlayKeyMsg{key: "a"},
		"b": testOverlayKeyMsg{key: "b"},
	})

	if !o.HasKey("a") || !o.HasKey("b") {
		t.Error("RegisterKeys did not register all keys")
	}
}

func TestOverlayValidate(t *testing.T) {
	o := New()

	err := o.Validate()
	if err != ErrNoExitKey {
		t.Errorf("Validate() = %v, want ErrNoExitKey", err)
	}

	o.SetExitKey("esc")
	err = o.Validate()
	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestOverlaySetStyle(t *testing.T) {
	style := lipgloss.NewStyle().Background(lipgloss.Color("#ff0000"))
	o := New().SetStyle(style)

	// Style should be set (we can check by comparing render output behavior)
	if o.Style().GetBackground() != lipgloss.Color("#ff0000") {
		t.Error("SetStyle did not apply style")
	}
}

func TestOverlayRender(t *testing.T) {
	o := New().
		SetContent("Hello World").
		SetSize(50).
		SetExitKey("esc")

	// Render at 100x50
	rendered := o.Render(100, 50)

	// Should contain the content
	if rendered == "" {
		t.Error("Render() returned empty string")
	}
}
