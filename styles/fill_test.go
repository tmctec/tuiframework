package styles

import (
	"testing"

	"charm.land/lipgloss/v2"
)

func TestNewThemeDerivesBandFills(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	// A fill reverses the text: the base background on the hero color.
	wantStyleColors(t, "HeaderFill", th.HeaderFill, p.Bg0, p.Purple)
	wantStyleColors(t, "FooterFill", th.FooterFill, p.Bg0, p.Purple)

	p.Hero = lipgloss.Color("#b00000")
	th = NewTheme(p)
	wantStyleColors(t, "HeaderFill with a hero", th.HeaderFill, p.Bg0, p.Hero)
	wantStyleColors(t, "FooterFill with a hero", th.FooterFill, p.Bg0, p.Hero)
}

func TestNewThemeUsesExplicitBandFills(t *testing.T) {
	p := fullPalette()
	p.HeaderFillBg, p.HeaderFillFg = lipgloss.Color("#f00000"), lipgloss.Color("#f00001")
	p.FooterFillBg, p.FooterFillFg = lipgloss.Color("#f00002"), lipgloss.Color("#f00003")
	th := NewTheme(p)

	wantStyleColors(t, "HeaderFill", th.HeaderFill, p.HeaderFillFg, p.HeaderFillBg)
	wantStyleColors(t, "FooterFill", th.FooterFill, p.FooterFillFg, p.FooterFillBg)
}

// TestBandFillsAreReadable holds every built-in theme to the WCAG minimum
// for interface text, 3:1, on its filled header and footer. The theme's
// ordinary text color on the same fill falls short, which is why a fill
// reverses the text.
func TestBandFillsAreReadable(t *testing.T) {
	for _, th := range AllThemes {
		checkBandFills(t, th)
	}
}

// checkBandFills holds a theme to the contrast rules for its filled bands.
func checkBandFills(t *testing.T, th *Theme) {
	t.Helper()
	const minimum = 3.0
	fills := []struct {
		name  string
		style lipgloss.Style
	}{
		{"HeaderFill", th.HeaderFill},
		{"FooterFill", th.FooterFill},
	}
	for _, f := range fills {
		reversed := Contrast(f.style.GetForeground(), f.style.GetBackground())
		if reversed < minimum {
			t.Errorf("%s %s: contrast %.2f, want at least %.1f", th.Name(), f.name, reversed, minimum)
		}
		if plain := Contrast(th.Palette.Fg1, f.style.GetBackground()); plain > reversed {
			t.Errorf("%s %s: body text reads better (%.2f) than the fill's text (%.2f)", th.Name(), f.name, plain, reversed)
		}
	}
}
