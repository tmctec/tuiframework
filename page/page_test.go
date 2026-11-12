package page

import (
	"testing"

	"github.com/tmctec/stickers/v2/flexbox"
)

// TestMsg is a test message type
type TestMsg struct {
	Value string
}

func TestNewPage(t *testing.T) {
	p := New("TestPage")
	if p.Name() != "TestPage" {
		t.Errorf("expected name 'TestPage', got %q", p.Name())
	}
	if p.Bindings() == nil {
		t.Error("bindings map should be initialized")
	}
	if len(p.Bindings()) != 0 {
		t.Error("bindings map should be empty initially")
	}
}

func TestBindKey(t *testing.T) {
	msg := TestMsg{Value: "test"}

	p := New("Test").BindKey("a", "Action A", msg)

	if !p.HasKey("a") {
		t.Error("key 'a' should be registered")
	}
	if p.HasKey("b") {
		t.Error("key 'b' should not be registered")
	}

	// Verify the binding is stored correctly
	binding := p.Bindings()["a"]
	if binding.Label() != "Action A" {
		t.Errorf("expected label 'Action A', got %q", binding.Label())
	}
	storedMsg := binding.Msg().(TestMsg)
	if storedMsg.Value != "test" {
		t.Errorf("expected message value 'test', got %q", storedMsg.Value)
	}
}

func TestBindKeyMultiple(t *testing.T) {
	p := New("Test").
		BindKey("a", "Action A", TestMsg{Value: "a"}).
		BindKey("b", "Action B", TestMsg{Value: "b"}).
		BindKey("c", "Action C", TestMsg{Value: "c"})

	if !p.HasKey("a") || !p.HasKey("b") || !p.HasKey("c") {
		t.Error("all keys should be registered")
	}
	if len(p.Bindings()) != 3 {
		t.Errorf("expected 3 bindings, got %d", len(p.Bindings()))
	}
}

func TestUnbindKey(t *testing.T) {
	p := New("Test").
		BindKey("a", "Action A", TestMsg{Value: "a"}).
		BindKey("b", "Action B", TestMsg{Value: "b"})

	if !p.HasKey("a") {
		t.Error("key 'a' should be registered")
	}

	p.UnbindKey("a")

	if p.HasKey("a") {
		t.Error("key 'a' should be unbound")
	}
	if !p.HasKey("b") {
		t.Error("key 'b' should still be registered")
	}
}

func TestBindKeyChaining(t *testing.T) {
	p := New("Test").
		BindKey("a", "Action A", TestMsg{Value: "a"}).
		BindKey("b", "Action B", TestMsg{Value: "b"}).
		UnbindKey("a").
		BindKey("c", "Action C", TestMsg{Value: "c"})

	if p.HasKey("a") {
		t.Error("key 'a' should be unbound")
	}
	if !p.HasKey("b") || !p.HasKey("c") {
		t.Error("keys 'b' and 'c' should be registered")
	}
}

func TestBindKeyOverwrite(t *testing.T) {
	p := New("Test").
		BindKey("a", "First", TestMsg{Value: "first"}).
		BindKey("a", "Second", TestMsg{Value: "second"})

	binding := p.Bindings()["a"]
	if binding.Label() != "Second" {
		t.Errorf("expected label 'Second' (overwrite), got %q", binding.Label())
	}
	storedMsg := binding.Msg().(TestMsg)
	if storedMsg.Value != "second" {
		t.Errorf("expected message value 'second' (overwrite), got %q", storedMsg.Value)
	}
}

func TestBindKeyEmptyLabel(t *testing.T) {
	p := New("Test").BindKey("a", "", TestMsg{Value: "test"})

	binding := p.Bindings()["a"]
	if binding.Label() != "" {
		t.Errorf("expected empty label, got %q", binding.Label())
	}
}

func TestSetContent(t *testing.T) {
	p := New("Test").SetContent(func(_ *flexbox.FlexBox, _ interface{}, _, _ int) {
		// Content function set
	})

	// Note: Content() returns the function, it doesn't execute it
	if p.Content() == nil {
		t.Error("content function should be set")
	}
}

func TestSetFooter(t *testing.T) {
	footer := map[string]string{
		"Enter": "Select",
		"q":     "Quit",
	}

	p := New("Test").SetFooter(footer)

	if len(p.Footer()) != 2 {
		t.Errorf("expected 2 footer items, got %d", len(p.Footer()))
	}
	if p.Footer()["Enter"] != "Select" {
		t.Error("footer 'Enter' should be 'Select'")
	}
}

func TestAddFooterItem(t *testing.T) {
	p := New("Test").
		AddFooterItem("a", "first").
		AddFooterItem("b", "second")

	if len(p.Footer()) != 2 {
		t.Errorf("expected 2 footer items, got %d", len(p.Footer()))
	}
}
