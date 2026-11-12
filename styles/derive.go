package styles

import (
	"errors"
	"fmt"
	"image/color"
	"math"

	"charm.land/lipgloss/v2"
)

// ErrLowContrast is returned by PaletteFrom for colors that would not be
// readable together.
var ErrLowContrast = errors.New("styles: colors are too close to read")

// The least contrast PaletteFrom accepts, as WCAG ratios: body text on the
// background, and the accent beside the background, where it fills tabs and
// bars that carry text in the background color.
const (
	minTextContrast   = 4.5
	minAccentContrast = 3.0
)

// Contrast returns the WCAG contrast ratio of two colors, from 1 for the
// same color to 21 for black on white. Body text wants 4.5 or more; large
// text and the edges of controls want 3.
func Contrast(a, b color.Color) float64 {
	la, lb := luminance(a), luminance(b)
	if la < lb {
		la, lb = lb, la
	}
	return (la + 0.05) / (lb + 0.05)
}

// luminance returns the WCAG relative luminance of a color, 0 for black to
// 1 for white. A color that is not set counts as black.
func luminance(c color.Color) float64 {
	if !isSet(c) {
		return 0
	}
	r, g, b, _ := c.RGBA()
	channel := func(x uint32) float64 {
		f := float64(x>>8) / 255
		if f <= 0.03928 {
			return f / 12.92
		}
		return math.Pow((f+0.055)/1.055, 2.4)
	}
	return 0.2126*channel(r) + 0.7152*channel(g) + 0.0722*channel(b)
}

// isDark reports whether c is darker than mid lightness, the test Bubble
// Tea applies to a terminal's background.
func isDark(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return max(r, g, b)+min(r, g, b) < 0xffff
}

// PaletteFrom derives a whole palette from three colors: the background,
// the body text and one accent. The background and text levels are steps
// between the two, the accent becomes the hero color that tabs, borders and
// fills take, and the base colors are a set chosen for a dark or a light
// background, whichever bg is. The palette is named "Custom" and has a
// rounded border; set Name, Border or any other field on the result before
// passing it to NewTheme.
//
// It returns ErrLowContrast when the text has less than 4.5:1 contrast on
// the background, or the accent less than 3:1.
func PaletteFrom(bg, fg, accent color.Color) (Palette, error) {
	if !isSet(bg) || !isSet(fg) || !isSet(accent) {
		return Palette{}, fmt.Errorf("%w: the background, text and accent must all be set", ErrLowContrast)
	}
	if c := Contrast(fg, bg); c < minTextContrast {
		return Palette{}, fmt.Errorf("%w: text on background is %.1f:1, want at least %.1f:1", ErrLowContrast, c, minTextContrast)
	}
	if c := Contrast(accent, bg); c < minAccentContrast {
		return Palette{}, fmt.Errorf("%w: accent on background is %.1f:1, want at least %.1f:1", ErrLowContrast, c, minAccentContrast)
	}

	dark := isDark(bg)
	pick := lipgloss.LightDark(dark) // pick(light, dark)
	hex := lipgloss.Color

	// toward moves a color away from the background: lighter on a dark
	// theme, darker on a light one. away does the opposite.
	toward, away := lipgloss.Lighten, lipgloss.Darken
	if !dark {
		toward, away = away, toward
	}

	// Thirteen steps from the background to the text: index 0 is bg, 12 is fg.
	ramp := lipgloss.Blend1D(13, bg, fg)

	p := Palette{
		Name:   "Custom",
		Border: lipgloss.RoundedBorder(),

		Bg0: bg, Bg1: ramp[1], Bg2: ramp[2], Bg3: ramp[4],
		Fg0: toward(fg, 0.08), Fg1: fg, Fg2: ramp[10], Fg3: ramp[8],

		Gray:   ramp[6],
		Red:    pick(hex("#D20F39"), hex("#F7768E")),
		Orange: pick(hex("#C2570C"), hex("#FF9E64")),
		Green:  pick(hex("#2F7D32"), hex("#9ECE6A")),
		Yellow: pick(hex("#9A6700"), hex("#E0AF68")),
		Blue:   pick(hex("#1E66F5"), hex("#7AA2F7")),
		Purple: accent,
		Pink:   pick(hex("#C2338F"), hex("#F5A9D0")),
		Aqua:   pick(hex("#0B7C93"), hex("#7DCFFF")),

		Hero: accent, HeroDk: away(accent, 0.25),
	}

	// Each base color gets a lighter and a darker shade.
	for _, shade := range []struct{ base, lt, dk *color.Color }{
		{&p.Gray, &p.GrayLt, &p.GrayDk}, {&p.Red, &p.RedLt, &p.RedDk},
		{&p.Orange, &p.OrangeLt, &p.OrangeDk}, {&p.Green, &p.GreenLt, &p.GreenDk},
		{&p.Yellow, &p.YellowLt, &p.YellowDk}, {&p.Blue, &p.BlueLt, &p.BlueDk},
		{&p.Purple, &p.PurpleLt, &p.PurpleDk}, {&p.Pink, &p.PinkLt, &p.PinkDk},
		{&p.Aqua, &p.AquaLt, &p.AquaDk},
	} {
		*shade.lt, *shade.dk = lipgloss.Lighten(*shade.base, 0.15), lipgloss.Darken(*shade.base, 0.2)
	}
	return p, nil
}

// MustPaletteFrom is PaletteFrom for colors known to be readable; it panics
// on an error.
func MustPaletteFrom(bg, fg, accent color.Color) Palette {
	p, err := PaletteFrom(bg, fg, accent)
	if err != nil {
		panic(err)
	}
	return p
}
