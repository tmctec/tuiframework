package styles

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestNewThemeResolvesSemanticColors(t *testing.T) {
	th := NewTheme(DraculaPalette)

	if !isSet(th.Palette.Primary) {
		t.Fatal("Primary should be resolved")
	}
	if th.Palette.Primary != DraculaPalette.Hero {
		t.Errorf("Primary = %v, want Hero %v", th.Palette.Primary, DraculaPalette.Hero)
	}
	if th.Palette.Secondary != DraculaPalette.Blue {
		t.Errorf("Secondary = %v, want Blue %v", th.Palette.Secondary, DraculaPalette.Blue)
	}
	if th.Palette.Success != DraculaPalette.Green || th.Palette.Error != DraculaPalette.Red {
		t.Error("Success/Error should resolve to Green/Red")
	}
}

func TestNewThemeKeepsExplicitSemanticColors(t *testing.T) {
	p := DraculaPalette
	p.Primary = lipgloss.Color("#123456")

	th := NewTheme(p)

	if th.Palette.Primary != lipgloss.Color("#123456") {
		t.Errorf("Primary = %v, want explicit #123456", th.Palette.Primary)
	}
	if DraculaPalette.Primary != nil {
		t.Error("the package-level palette must stay untouched")
	}
}

func TestAllThemesHaveResolvedPrimary(t *testing.T) {
	for _, th := range AllThemes {
		if !isSet(th.Palette.Primary) {
			t.Errorf("%s: Primary is empty", th.Name())
		}
	}
}
