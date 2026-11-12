package tabapp

import (
	"testing"
	"time"

	"github.com/tmctec/tuiframework/styles"
)

func TestNewModelDefaults(t *testing.T) {
	m := NewModel()
	checks := []struct {
		name string
		ok   bool
	}{
		{"NavBarPosition", m.NavBarPosition == NavBarTop},
		{"NavBarAutoHide", m.NavBarAutoHide},
		{"NavBarTimeout", m.NavBarTimeout == 2*time.Second},
		{"NavBarMode", m.NavBarMode == NavBarCompact},
		{"TabNavTimeout", m.TabNavTimeout == 2*time.Second},
		{"TabNavActive", !m.TabNavActive},
		{"Theme", m.Theme == styles.Dracula},
		{"ActivePage", m.ActivePage == 0},
		{"NavBarBelowHeader", m.NavBarBelowHeader},
		{"band heights", m.HeaderHeight == 1 && m.NavBarHeight == 3 && m.FooterHeight == 3},
		{"frames none", m.HeaderFrame == FrameNone && m.FooterFrame == FrameNone},
		{"borders off", !m.ShowNavBarBorder && !m.ShowContentBorder},
		{"KeyGlowEnabled", m.KeyGlowEnabled},
		{"OverlayCount", m.OverlayCount == 0},
		{"Custom empty", m.Custom != nil && len(m.Custom) == 0},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("default %s does not hold", c.name)
		}
	}
}

func TestModelCustomAccessors(t *testing.T) {
	m := NewModel()
	if m.SetCustom("s", "text") != m {
		t.Error("SetCustom should return the model for chaining")
	}
	m.SetCustom("i", 7).SetCustom("b", true)

	if got := m.GetCustom("s"); got != "text" {
		t.Errorf("GetCustom(s) = %v, want text", got)
	}
	if got := m.GetCustom("missing"); got != nil {
		t.Errorf("GetCustom(missing) = %v, want nil", got)
	}

	if got := m.GetCustomString("s"); got != "text" {
		t.Errorf("GetCustomString(s) = %q", got)
	}
	if got := m.GetCustomString("i"); got != "" {
		t.Errorf("GetCustomString of an int = %q, want empty", got)
	}
	if got := m.GetCustomString("missing"); got != "" {
		t.Errorf("GetCustomString(missing) = %q, want empty", got)
	}

	if got := m.GetCustomInt("i"); got != 7 {
		t.Errorf("GetCustomInt(i) = %d", got)
	}
	if got := m.GetCustomInt("s"); got != 0 {
		t.Errorf("GetCustomInt of a string = %d, want 0", got)
	}
	if got := m.GetCustomInt("missing"); got != 0 {
		t.Errorf("GetCustomInt(missing) = %d, want 0", got)
	}

	if !m.GetCustomBool("b") {
		t.Error("GetCustomBool(b) should be true")
	}
	if m.GetCustomBool("i") {
		t.Error("GetCustomBool of an int should be false")
	}
	if m.GetCustomBool("missing") {
		t.Error("GetCustomBool(missing) should be false")
	}
}
