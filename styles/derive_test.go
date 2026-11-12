package styles

import (
	"errors"
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

// The seeds the derived-palette tests use: one dark, one light.
var (
	darkSeed  = [3]color.Color{lipgloss.Color("#101418"), lipgloss.Color("#d7dde4"), lipgloss.Color("#7aa2f7")}
	lightSeed = [3]color.Color{lipgloss.Color("#fafafa"), lipgloss.Color("#24292f"), lipgloss.Color("#0969da")}
)

func TestContrast(t *testing.T) {
	black, white := lipgloss.Color("#000000"), lipgloss.Color("#ffffff")
	if got := Contrast(black, white); got < 20.9 || got > 21.1 {
		t.Errorf("black on white = %.2f, want 21", got)
	}
	if got := Contrast(white, black); got < 20.9 || got > 21.1 {
		t.Errorf("white on black = %.2f, want 21 whichever comes first", got)
	}
	if got := Contrast(white, white); got != 1 {
		t.Errorf("a color on itself = %.2f, want 1", got)
	}
	if got := Contrast(nil, white); got < 20.9 {
		t.Errorf("an unset color counts as black: %.2f, want 21", got)
	}
}

// ps-1 and ps-2: three colors give a whole palette, dark or light with its
// background, that meets the contrast the built-in themes meet.
func TestPaletteFromDerivesAReadableTheme(t *testing.T) {
	for name, seed := range map[string][3]color.Color{"dark": darkSeed, "light": lightSeed} {
		t.Run(name, func(t *testing.T) {
			p, err := PaletteFrom(seed[0], seed[1], seed[2])
			if err != nil {
				t.Fatal(err)
			}
			if p.Name != "Custom" || p.Border != lipgloss.RoundedBorder() {
				t.Errorf("name %q border %v, want Custom with a rounded border", p.Name, p.Border)
			}
			if p.Bg0 != seed[0] || p.Fg1 != seed[1] || p.Hero != seed[2] || p.Purple != seed[2] {
				t.Error("the palette should keep the background, the text and the accent as given")
			}
			for _, c := range requiredPaletteColors(p) {
				if !isSet(c.color) {
					t.Errorf("%s is not set", c.field)
				}
			}

			th := NewTheme(p)
			if got, want := th.IsDark(), name == "dark"; got != want {
				t.Errorf("IsDark = %v, want %v", got, want)
			}
			for _, c := range themeColors(th) {
				if !isSet(c.color) {
					t.Errorf("%s is not set", c.field)
				}
			}
			checkBandFills(t, th)

			// Text keeps its order of emphasis and stays readable.
			for _, text := range []struct {
				field string
				c     color.Color
				min   float64
			}{
				{"Fg0", p.Fg0, 4.5}, {"Fg1", p.Fg1, 4.5}, {"Fg2", p.Fg2, 4.5}, {"Fg3", p.Fg3, 3},
				{"Red", p.Red, 3}, {"Orange", p.Orange, 3}, {"Green", p.Green, 3}, {"Yellow", p.Yellow, 3},
				{"Blue", p.Blue, 3}, {"Pink", p.Pink, 3}, {"Aqua", p.Aqua, 3},
			} {
				if got := Contrast(text.c, p.Bg0); got < text.min {
					t.Errorf("%s on the background is %.2f:1, want at least %.1f:1", text.field, got, text.min)
				}
			}
			if Contrast(p.Fg0, p.Bg0) < Contrast(p.Fg1, p.Bg0) || Contrast(p.Fg2, p.Bg0) < Contrast(p.Fg3, p.Bg0) {
				t.Error("the text levels should run from the strongest, Fg0, to the weakest, Fg3")
			}
			// The surfaces step away from the background, a little at a time.
			if c1, c3 := Contrast(p.Bg1, p.Bg0), Contrast(p.Bg3, p.Bg0); c1 <= 1 || c3 <= c1 || c3 > 3 {
				t.Errorf("Bg1 and Bg3 stand %.2f and %.2f from the background, want small rising steps", c1, c3)
			}
		})
	}
}

// ps-3: colors that would not be readable are refused.
func TestPaletteFromRefusesLowContrast(t *testing.T) {
	bg := lipgloss.Color("#101418")
	tests := []struct {
		name           string
		bg, fg, accent color.Color
	}{
		{"text too close to the background", bg, lipgloss.Color("#40464d"), lipgloss.Color("#7aa2f7")},
		{"accent too close to the background", bg, lipgloss.Color("#d7dde4"), lipgloss.Color("#2a3038")},
		{"no background", nil, lipgloss.Color("#d7dde4"), lipgloss.Color("#7aa2f7")},
		{"no text", bg, lipgloss.NoColor{}, lipgloss.Color("#7aa2f7")},
		{"no accent", bg, lipgloss.Color("#d7dde4"), nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := PaletteFrom(tt.bg, tt.fg, tt.accent)
			if !errors.Is(err, ErrLowContrast) {
				t.Fatalf("err = %v, want ErrLowContrast", err)
			}
			if p != (Palette{}) {
				t.Error("a refused palette should be empty")
			}
			defer func() {
				if recover() == nil {
					t.Error("MustPaletteFrom should panic on the same colors")
				}
			}()
			MustPaletteFrom(tt.bg, tt.fg, tt.accent)
		})
	}

	if p := MustPaletteFrom(darkSeed[0], darkSeed[1], darkSeed[2]); p.Hero != darkSeed[2] {
		t.Error("MustPaletteFrom should return the palette for readable colors")
	}
}

// The theme takes a border gradient only from a palette that sets both ends.
func TestThemeBorderBlend(t *testing.T) {
	from, to := lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")

	p := basePalette()
	if th := NewTheme(p); th.BorderBlend != nil {
		t.Errorf("no gradient set: BorderBlend = %v, want nil", th.BorderBlend)
	}
	p.BorderBlend = [2]color.Color{from, nil}
	if th := NewTheme(p); th.BorderBlend != nil {
		t.Errorf("one end set: BorderBlend = %v, want nil", th.BorderBlend)
	}
	p.BorderBlend = [2]color.Color{from, to}
	th := NewTheme(p)
	if len(th.BorderBlend) != 2 || th.BorderBlend[0] != from || th.BorderBlend[1] != to {
		t.Errorf("BorderBlend = %v, want the two ends", th.BorderBlend)
	}
	if th.Palette.BorderBlend != p.BorderBlend {
		t.Error("the resolved palette should keep the gradient")
	}
	for _, built := range AllThemes {
		if built.BorderBlend != nil {
			t.Errorf("%s has a border gradient; the built-in themes have none", built.Name())
		}
	}
}
