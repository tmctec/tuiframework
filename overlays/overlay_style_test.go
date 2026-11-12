package overlays

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func TestOverlayStyleSettersReturnSelf(t *testing.T) {
	o := New()
	results := []*Overlay{
		o.SetBorder(lipgloss.DoubleBorder()),
		o.SetBorderForeground(lipgloss.Color("#ff0000")),
		o.SetBorderBackground(lipgloss.Color("#00ff00")),
		o.SetPadding(2, 3),
		o.SetBackground(lipgloss.Color("#0000ff")),
		o.SetForeground(lipgloss.Color("#ffffff")),
		o.SetTitle("t"),
		o.SetTitleStyle(lipgloss.NewStyle()),
		o.SetZIndex(1),
		o.SetWidthPct(10),
		o.SetHeightPct(10),
		o.SetWidth(10),
		o.SetHeight(10),
		o.SetFlexbox(nil),
		o.Show(),
		o.Hide(),
	}
	for i, r := range results {
		if r != o {
			t.Errorf("setter %d did not return the receiver", i)
		}
	}
}

func TestOverlayStyleSettersAccumulate(t *testing.T) {
	border := lipgloss.Color("#ff0000")
	borderBg := lipgloss.Color("#00ff00")
	bg := lipgloss.Color("#0000ff")
	fg := lipgloss.Color("#ffffff")

	o := New().
		SetBorder(lipgloss.DoubleBorder()).
		SetBorderForeground(border).
		SetBorderBackground(borderBg).
		SetPadding(2, 3).
		SetBackground(bg).
		SetForeground(fg)

	s := o.Style()
	if s.GetBorderStyle() != lipgloss.DoubleBorder() {
		t.Error("SetBorder did not apply the double border")
	}
	if s.GetBorderTopForeground() != border || s.GetBorderBottomForeground() != border ||
		s.GetBorderLeftForeground() != border || s.GetBorderRightForeground() != border {
		t.Error("SetBorderForeground should color all four sides")
	}
	if s.GetBorderTopBackground() != borderBg || s.GetBorderLeftBackground() != borderBg {
		t.Error("SetBorderBackground should color the border background")
	}
	if s.GetPaddingTop() != 2 || s.GetPaddingBottom() != 2 || s.GetPaddingLeft() != 3 || s.GetPaddingRight() != 3 {
		t.Errorf("padding = %d %d %d %d, want 2 3 2 3",
			s.GetPaddingTop(), s.GetPaddingRight(), s.GetPaddingBottom(), s.GetPaddingLeft())
	}
	if s.GetBackground() != bg {
		t.Errorf("background = %v, want %v", s.GetBackground(), bg)
	}
	if s.GetForeground() != fg {
		t.Errorf("foreground = %v, want %v", s.GetForeground(), fg)
	}
}

func TestOverlayNewStyleDefaults(t *testing.T) {
	s := New().Style()
	if s.GetBorderStyle() != lipgloss.RoundedBorder() {
		t.Error("a new overlay should have a rounded border")
	}
	if s.GetPaddingTop() != 1 || s.GetPaddingLeft() != 2 {
		t.Errorf("padding = %d,%d, want 1,2", s.GetPaddingTop(), s.GetPaddingLeft())
	}
	if s.GetHorizontalFrameSize() != 6 || s.GetVerticalFrameSize() != 4 {
		t.Errorf("frame = %dx%d, want 6x4", s.GetHorizontalFrameSize(), s.GetVerticalFrameSize())
	}
}

func TestOverlayZIndex(t *testing.T) {
	o := New()
	if o.ZIndex() != 0 {
		t.Errorf("default ZIndex() = %d, want 0", o.ZIndex())
	}
	o.SetZIndex(5)
	if o.ZIndex() != 5 {
		t.Errorf("ZIndex() = %d, want 5", o.ZIndex())
	}
	o.SetZIndex(-2)
	if o.ZIndex() != -2 {
		t.Errorf("ZIndex() = %d, want -2", o.ZIndex())
	}
}

func TestOverlayRegisterKeyOnZeroValue(t *testing.T) {
	var o Overlay
	o.RegisterKey("a", testOverlayKeyMsg{key: "a"})
	if !o.HasKey("a") {
		t.Error("RegisterKey should create the key map on a zero-value overlay")
	}

	var o2 Overlay
	o2.RegisterKeys(map[string]tea.Msg{"b": testOverlayKeyMsg{key: "b"}})
	if !o2.HasKey("b") {
		t.Error("RegisterKeys should create the key map on a zero-value overlay")
	}

	var empty Overlay
	if empty.HasKey("x") || len(empty.Keys()) != 0 {
		t.Error("a zero-value overlay has no keys")
	}
}

func TestOverlayRegisterKeyReplacesBinding(t *testing.T) {
	o := New().
		RegisterKey("a", testOverlayKeyMsg{key: "first"}).
		RegisterKey("a", testOverlayKeyMsg{key: "second"})

	if len(o.Keys()) != 1 {
		t.Errorf("Keys() has %d entries, want 1", len(o.Keys()))
	}
	if msg, ok := o.Keys()["a"].(testOverlayKeyMsg); !ok || msg.key != "second" {
		t.Errorf("Keys()[a] = %v, want the second binding", o.Keys()["a"])
	}
}

func TestOverlayRegisterKeysAddsToExisting(t *testing.T) {
	o := New().
		RegisterKey("a", testOverlayKeyMsg{key: "a"}).
		RegisterKeys(map[string]tea.Msg{
			"b": testOverlayKeyMsg{key: "b"},
			"c": testOverlayKeyMsg{key: "c"},
		})

	for _, k := range []string{"a", "b", "c"} {
		if !o.HasKey(k) {
			t.Errorf("HasKey(%q) = false, want true", k)
		}
	}
	if len(o.Keys()) != 3 {
		t.Errorf("Keys() has %d entries, want 3", len(o.Keys()))
	}
}

func TestOverlayExitKeyIsNotARegisteredKey(t *testing.T) {
	o := New().SetExitKey("esc")
	if o.HasKey("esc") {
		t.Error("the exit key is separate from the registered key bindings")
	}
	if o.ExitKey() != "esc" {
		t.Errorf("ExitKey() = %q, want esc", o.ExitKey())
	}
}

// TestOverlayBorderBlend covers the gradient setter: two or more colors set
// it, fewer remove it, and a single border color replaces it.
func TestOverlayBorderBlend(t *testing.T) {
	red, blue := lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")
	o := New().SetBorder(lipgloss.RoundedBorder()).SetBorderBlend(red, blue)
	if got := o.Style().GetBorderForegroundBlend(); len(got) != 2 {
		t.Fatalf("blend = %v, want two colors", got)
	}

	o.SetBorderForeground(red)
	if got := o.Style().GetBorderForegroundBlend(); len(got) != 0 {
		t.Errorf("a single border color should replace the gradient, blend = %v", got)
	}
	if o.Style().GetBorderTopForeground() != red {
		t.Error("the border should take the single color")
	}

	o.SetBorderBlend(red, blue).SetBorderBlend(red)
	if got := o.Style().GetBorderForegroundBlend(); len(got) != 0 {
		t.Errorf("fewer than two colors remove the gradient, blend = %v", got)
	}
}
