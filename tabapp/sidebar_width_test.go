package tabapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/styles"
)

// TestModalBorderColorSurvivesThemeChange opens a modal with its own border
// color, changes the theme, and checks the override is still in place while
// the rest of the frame took the new theme.
func TestModalBorderColorSurvivesThemeChange(t *testing.T) {
	app := rootMixApp()
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app.Update(ToggleModalMsg{Title: "A", Content: "body", SizePct: 50, BorderColor: lipgloss.Color("#ff0000")})

	app.UseTheme(styles.Gruvbox)

	style := app.TopOverlay().Style()
	if got := style.GetBorderTopForeground(); got != lipgloss.Color("#ff0000") {
		t.Errorf("border = %v, want the #ff0000 override", got)
	}
	if got := style.GetBackground(); got != styles.Gruvbox.ModalBackground {
		t.Errorf("background = %v, want Gruvbox's %v", got, styles.Gruvbox.ModalBackground)
	}
}

// TestSideBarWidthHasAMinimum checks a side bar too narrow for its border,
// the tab padding and four columns of text is widened, so the bar and the
// content share the frame width without overflowing it.
func TestSideBarWidthHasAMinimum(t *testing.T) {
	needsNarrowBorders(t)
	app := rootMixApp()
	app.SetNavBarAutoHide(false).SetNavBarPosition(NavBarLeft).SetNavBarWidth(4).SetShowNavBarBorder(true)
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	if app.NavBarWidth() != minNavBarWidth {
		t.Fatalf("NavBarWidth = %d, want %d", app.NavBarWidth(), minNavBarWidth)
	}
	checkFrame(t, "bordered", app.Render(), 80, 24)
	// The bar shows the level-0 items, each on one line, inside its border.
	plain := plainView(app)
	for _, want := range []string{"│", "▸ A", "▸ B", "S…"} {
		if !strings.Contains(plain, want) {
			t.Errorf("%q missing from the side bar:\n%s", want, plain)
		}
	}

	// Without the border the names get its two columns.
	app.SetShowNavBarBorder(false)
	checkFrame(t, "bare", app.Render(), 80, 24)
	plain = plainView(app)
	if strings.Contains(plain, "│") || !strings.Contains(plain, "Set…") {
		t.Errorf("a bare side bar should have no border and wider names:\n%s", plain)
	}

	app.Update(SetNavBarWidthMsg{Width: 3})
	if app.NavBarWidth() != minNavBarWidth {
		t.Errorf("SetNavBarWidthMsg below the minimum left %d, want %d", app.NavBarWidth(), minNavBarWidth)
	}
}

// TestNavBarBorderCoversBothPositions checks the one border setting: off,
// neither the top bar nor the side bar has a border; on, both do, in the nav
// bar border's shape, from the theme or from SetNavBarBorderStyle.
func TestNavBarBorderCoversBothPositions(t *testing.T) {
	needsNarrowBorders(t)
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		app := rootMixApp()
		app.SetNavBarPosition(pos).SetNavBarAutoHide(false).SetNavBarHeight(2)
		app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

		bar := func() string {
			if pos == NavBarLeft {
				return ansiSequence.ReplaceAllString(app.renderSideNavBar(10), "")
			}
			return ansiSequence.ReplaceAllString(app.renderNavSection(), "")
		}
		corner := func() string { return string([]rune(bar())[0]) }
		size := func() (int, int) { return lipgloss.Width(bar()), lipgloss.Height(bar()) }
		wantW, wantH := 80, 2
		if pos == NavBarLeft {
			wantW, wantH = app.NavBarWidth(), 10
		}

		// A bare bar is only its text: the canvas cuts it to the bar's
		// rectangle and fills the rest. A top bar keeps its height.
		wrongSize := func(w, h int) bool {
			if pos == NavBarTop {
				return w > wantW || h != wantH
			}
			return w > wantW || h > wantH
		}
		if w, h := size(); strings.ContainsAny(bar(), "╭╮│─") || wrongSize(w, h) {
			t.Errorf("pos=%v border off: %dx%d, want a bare bar of %dx%d:\n%s", pos, w, h, wantW, wantH, bar())
		}

		app.Update(ToggleNavBarBorderMsg{})
		if pos == NavBarTop {
			wantH += 2 // the border takes a row above and below a top bar
		}
		if w, h := size(); corner() != styles.Dracula.NavBarBorder.TopLeft || w != wantW || h != wantH {
			t.Errorf("pos=%v border on: corner %q size %dx%d, want %q and %dx%d", pos, corner(), w, h, styles.Dracula.NavBarBorder.TopLeft, wantW, wantH)
		}
		app.UseTheme(styles.Retro)
		if got := corner(); got != styles.Retro.NavBarBorder.TopLeft {
			t.Errorf("pos=%v: corner = %q, want Retro's %q", pos, got, styles.Retro.NavBarBorder.TopLeft)
		}
		app.SetNavBarBorderStyle(lipgloss.ThickBorder())
		if got := corner(); got != lipgloss.ThickBorder().TopLeft {
			t.Errorf("pos=%v: corner = %q, want the thick border's %q", pos, got, lipgloss.ThickBorder().TopLeft)
		}
		checkFrame(t, "bordered", app.Render(), 80, 24)
	}
}
