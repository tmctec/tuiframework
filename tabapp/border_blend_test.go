package tabapp

import (
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// borderColorsIn returns the distinct colors of the box-drawing characters
// in a frame.
func borderColorsIn(frame string, width int) map[string]bool {
	colors := map[string]bool{}
	for _, line := range strings.Split(frame, "\n") {
		for _, cell := range overlays.ParseCells(line, width) {
			if strings.ContainsRune("─│╭╮╰╯═║╔╗╚╝┌┐└┘", cell.Char) {
				colors[cell.FgColor] = true
			}
		}
	}
	return colors
}

// blendTheme is Nord with a border gradient from red to blue.
func blendTheme() *styles.Theme {
	p := styles.NordPalette
	p.Name = "NordBlend"
	p.BorderBlend = [2]color.Color{lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")}
	return styles.NewTheme(p)
}

// borderApp is a sized one-page app with only the content border drawn.
func borderApp(p *page.Page) *TabApp {
	app := New().MustAddPage(p).SetNavBarAutoHide(false).
		SetShowContentBorder(true).SetShowNavBarBorder(false).
		SetHeaderFrame(FrameNone).SetFooterFrame(FrameNone)
	return sizedApp(app)
}

// ps-4: a theme with a border gradient draws the content border in more
// than one color; a theme without draws it in exactly one.
func TestContentBorderBlend(t *testing.T) {
	needsNarrowBorders(t)
	app := borderApp(page.New("One"))

	app.UseTheme(styles.Nord)
	if got := borderColorsIn(app.Render(), app.width); len(got) != 1 {
		t.Errorf("without a gradient the border has %d colors, want 1", len(got))
	}

	app.UseTheme(blendTheme())
	if got := borderColorsIn(app.Render(), app.width); len(got) < 3 {
		t.Errorf("with a gradient the border has %d colors, want several", len(got))
	}
	if len(app.BorderBlend()) != 2 {
		t.Errorf("BorderBlend = %v, want the theme's two colors", app.BorderBlend())
	}

	// A theme without a gradient takes it away again.
	app.UseTheme(styles.Nord)
	if got := borderColorsIn(app.Render(), app.width); len(got) != 1 || app.BorderBlend() != nil {
		t.Errorf("back on a plain theme the border has %d colors, want 1", len(got))
	}
}

// The app can set the gradient itself, and fewer than two colors remove it.
func TestSetBorderBlend(t *testing.T) {
	needsNarrowBorders(t)
	app := borderApp(page.New("One"))

	app.SetBorderBlend(lipgloss.Color("#00ff00"), lipgloss.Color("#ff00ff"), lipgloss.Color("#ffff00"))
	if got := borderColorsIn(app.Render(), app.width); len(got) < 3 {
		t.Errorf("the border has %d colors, want a gradient", len(got))
	}
	app.SetBorderBlend(lipgloss.Color("#00ff00"))
	if app.BorderBlend() != nil {
		t.Error("one color is no gradient")
	}
	if got := borderColorsIn(app.Render(), app.width); len(got) != 1 {
		t.Errorf("the border has %d colors, want the single border color", len(got))
	}
}

// A page's own border color beats the gradient.
func TestPageBorderColorBeatsTheBlend(t *testing.T) {
	needsNarrowBorders(t)
	app := borderApp(page.New("Own").SetBorderColor(lipgloss.Color("#123456")))
	app.UseTheme(blendTheme())

	got := borderColorsIn(app.Render(), app.width)
	if len(got) != 1 || !got["38;2;18;52;86"] {
		t.Errorf("border colors = %v, want only the page's #123456", got)
	}
}

// The nav bar border takes the gradient in both positions.
func TestNavBarBorderBlend(t *testing.T) {
	needsNarrowBorders(t)
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		app := borderApp(page.New("One"))
		app.SetShowContentBorder(false).SetShowNavBarBorder(true).SetNavBarPosition(pos)

		app.UseTheme(styles.Nord)
		if got := borderColorsIn(app.Render(), app.width); len(got) != 1 {
			t.Errorf("position %v: %d border colors without a gradient, want 1", pos, len(got))
		}
		app.UseTheme(blendTheme())
		if got := borderColorsIn(app.Render(), app.width); len(got) < 3 {
			t.Errorf("position %v: %d border colors with a gradient, want several", pos, len(got))
		}
	}
}

// A modal shown through the app follows the theme's gradient, and a border
// color of its own replaces it.
func TestModalBorderBlend(t *testing.T) {
	app := borderApp(page.New("One"))
	app.UseTheme(blendTheme())

	app.Update(ToggleModalMsg{Title: "Blend", Content: "body", SizePct: 50})
	if got := app.TopOverlay().Style().GetBorderForegroundBlend(); len(got) != 2 {
		t.Errorf("the modal's border blend = %v, want the theme's", got)
	}

	app.Update(ToggleModalMsg{Title: "Own", Content: "body", SizePct: 40, ZIndex: 5, BorderColor: lipgloss.Color("#ff0000")})
	top := app.TopOverlay().Style()
	if len(top.GetBorderForegroundBlend()) != 0 || top.GetBorderTopForeground() != lipgloss.Color("#ff0000") {
		t.Errorf("a modal with its own border color should drop the gradient: blend %v color %v",
			top.GetBorderForegroundBlend(), top.GetBorderTopForeground())
	}

	// A theme without a gradient restyles the first modal plainly.
	app.UseTheme(styles.Nord)
	for _, e := range app.overlayStack {
		if len(e.overlay.Style().GetBorderForegroundBlend()) != 0 {
			t.Errorf("modal %q keeps a gradient after a theme without one", e.overlay.Title())
		}
	}
}
