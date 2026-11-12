package tabapp

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

func TestOptionGettersMirrorSetters(t *testing.T) {
	app := New()
	red := lipgloss.Color("#ff0000")
	blue := lipgloss.Color("#0000ff")
	green := lipgloss.Color("#00ff00")
	bold := lipgloss.NewStyle().Foreground(red).Bold(true)
	italic := lipgloss.NewStyle().Foreground(blue).Italic(true)
	handler := func(_ tea.Msg) (bool, tea.Cmd) { return false, nil }

	app.SetNavBarPosition(NavBarLeft).
		SetNavBarMode(NavBarPath).
		SetNavBarAutoHide(false).
		SetNavBarAutoHideTimeout(5 * time.Second).
		SetTabNavKeyTimeout(4 * time.Second).
		SetNavBarWidth(22).
		SetHeader(testHeader).
		SetHeaderHeight(2).
		SetNavBarHeight(4).
		SetFooterHeight(2).
		SetNavBarBelowHeader(false).
		SetShowHeaderBorder(true).
		SetShowNavBarBorder(true).
		SetShowFooterBorder(true).
		SetShowContentBorder(true).
		SetKeyGlowEnabled(false).
		SetKeyGlowDuration(time.Second).
		SetKeyGlowColor(green).
		SetFooterKeyColor(red).
		SetFooterValueColor(blue).
		SetContentBackground(green).
		SetTabStyle(bold).
		SetActiveTabStyle(italic).
		SetNavBarBorderStyle(lipgloss.DoubleBorder()).
		SetNavBarBorderColor(red).
		SetContentBorderStyle(lipgloss.ThickBorder()).
		SetContentBorderColor(blue).
		SetFooterBorderStyle(lipgloss.NormalBorder()).
		SetFooterBorderColor(green).
		SetFooterStyle(italic).
		SetCustomMessageHandler(handler)

	m := app.Model()
	checks := []struct {
		name string
		ok   bool
	}{
		{"NavBarPosition", app.NavBarPosition() == NavBarLeft && m.NavBarPosition == NavBarLeft},
		{"NavBarMode", app.NavBarMode() == NavBarPath && m.NavBarMode == NavBarPath},
		{"NavBarAutoHide", !app.NavBarAutoHide() && !m.NavBarAutoHide},
		{"NavBarAutoHideTimeout", app.NavBarAutoHideTimeout() == 5*time.Second && m.NavBarTimeout == 5*time.Second},
		{"TabNavKeyTimeout", app.TabNavKeyTimeout() == 4*time.Second && m.TabNavTimeout == 4*time.Second},
		{"NavBarWidth", app.NavBarWidth() == 22},
		{"Header", app.headerFunc != nil},
		{"HeaderHeight", app.HeaderHeight() == 2},
		{"NavBarHeight", app.NavBarHeight() == 4},
		{"FooterHeight", app.FooterHeight() == 2},
		{"NavBarBelowHeader", !app.NavBarBelowHeader() && !m.NavBarBelowHeader},
		{"ShowHeaderBorder", app.ShowHeaderBorder() && m.HeaderFrame == FrameTheme},
		{"ShowNavBarBorder", app.ShowNavBarBorder() && m.ShowNavBarBorder},
		{"ShowFooterBorder", app.ShowFooterBorder() && m.FooterFrame == FrameTheme},
		{"ShowContentBorder", app.ShowContentBorder() && m.ShowContentBorder},
		{"KeyGlowEnabled", !app.KeyGlowEnabled() && !m.KeyGlowEnabled},
		{"KeyGlowDuration", app.KeyGlowDuration() == time.Second},
		{"KeyGlowColor", app.KeyGlowColor() == green},
		{"FooterKeyColor", app.FooterKeyColor() == red},
		{"FooterValueColor", app.FooterValueColor() == blue},
		{"ContentBackground", app.ContentBackground() == green},
		{"TabStyle", app.TabStyle().GetForeground() == red && app.TabStyle().GetBold()},
		{"ActiveTabStyle", app.ActiveTabStyle().GetForeground() == blue && app.ActiveTabStyle().GetItalic()},
		{"NavBarBorderStyle", app.navBarBorderStyle == lipgloss.DoubleBorder()},
		{"NavBarBorderColor", app.navBarBorderColor == red},
		{"ContentBorderStyle", app.contentBorderStyle == lipgloss.ThickBorder()},
		{"ContentBorderColor", app.contentBorderColor == blue},
		{"FooterBorderStyle", app.footerBorderStyle == lipgloss.NormalBorder()},
		{"FooterBorderColor", app.footerBorderColor == green},
		{"FooterStyle", app.footerStyle.GetForeground() == blue},
		{"CustomMessageHandler", app.customMsgHandler != nil},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("%s does not reflect its setter", c.name)
		}
	}
}

func TestSetFooterSeparator(t *testing.T) {
	app := New()
	p := page.New("P").BindKey("a", "Alpha", NoOpMsg{}).BindKey("b", "Beta", NoOpMsg{})
	app.MustAddPage(p)
	app.SetFooterSeparator(" | ")

	got := ansiSequence.ReplaceAllString(app.renderPageFooter(p), "")
	if strings.Count(got, " | ") != 2 {
		t.Errorf("footer %q should join three items with the custom separator", got)
	}
}

func TestSetNavBarModeResetsCrumbFocusAndLevel(t *testing.T) {
	app := stripApp()
	key(app, "tab")
	key(app, "left")
	key(app, "left") // onto crumb A
	if app.navFocusCrumb != 0 {
		t.Fatalf("precondition: focus = %d, want crumb 0", app.navFocusCrumb)
	}

	app.SetNavBarMode(NavBarPath)

	if app.navFocusCrumb != -1 {
		t.Errorf("crumb focus = %d after a mode change, want -1", app.navFocusCrumb)
	}
	if app.TabNavLevel() != app.currentPageLevel() || app.TabNavLevel() != 1 {
		t.Errorf("level = %d, want the page's level 1", app.TabNavLevel())
	}

	hidden := rootMixApp()
	hidden.Init()
	hidden.Update(navBarAutoHideMsg{id: hidden.navBarAutoHideID})
	if !hidden.navBarHidden {
		t.Fatal("precondition: the bar should be hidden")
	}
	hidden.SetNavBarMode(NavBarExpanded)
	if hidden.navBarHidden {
		t.Error("the expanded bar shows at once")
	}
}

func TestSetNavBarPositionLeftInExpandedModeSetsLevel(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarMode(NavBarExpanded).SetNavBarAutoHide(false)
	key(app, "tab")
	key(app, "up")
	if app.TabNavLevel() != 0 {
		t.Fatalf("precondition: level = %d, want 0", app.TabNavLevel())
	}

	app.SetNavBarPosition(NavBarLeft)

	if app.TabNavLevel() != 1 || app.navFocusCrumb != -1 {
		t.Errorf("level %d focus %d, want the page level 1 and no crumb focus", app.TabNavLevel(), app.navFocusCrumb)
	}

	compact := stripApp()
	key(compact, "tab")
	key(compact, "backspace")
	if compact.TabNavLevel() != 0 {
		t.Fatalf("precondition: level = %d, want 0", compact.TabNavLevel())
	}
	compact.SetNavBarPosition(NavBarLeft)
	if compact.TabNavLevel() != 0 {
		t.Errorf("in compact mode the position change keeps level 0, got %d", compact.TabNavLevel())
	}
}

func TestUseThemeAppliesEverywhere(t *testing.T) {
	defer overlays.NotificationConfig.ApplyTheme(styles.Dracula)
	app := New()
	th := styles.Gruvbox

	if app.UseTheme(th) != app {
		t.Error("UseTheme should return the app for chaining")
	}

	checks := []struct {
		name string
		ok   bool
	}{
		{"Model.Theme", app.Model().Theme == th},
		{"TabStyle", app.TabStyle().GetForeground() == th.TabInactive.GetForeground()},
		{"ActiveTabStyle", app.ActiveTabStyle().GetForeground() == th.TabActive.GetForeground()},
		{"SelectedInactiveStyle", app.SelectedInactiveStyle().GetForeground() == th.TabSelectedMuted.GetForeground()},
		{"tab border", app.navBarBorderStyle == th.NavBarBorder && app.navBarBorderColor == th.NavBarBorderColor},
		{"header border", app.headerBorderColor == th.HeaderBorderColor},
		{"content border", app.contentBorderStyle == th.ContentBorder && app.contentBorderColor == th.ContentBorderColor},
		{"ContentBackground", app.ContentBackground() == th.Palette.Bg0},
		{"footer border", app.footerBorderStyle == th.FooterBorder && app.footerBorderColor == th.FooterBorderColor},
		{"footer style", app.footerStyle.GetBackground() == th.Footer.GetBackground()},
		{"FooterKeyColor", app.FooterKeyColor() == th.Palette.Primary},
		{"FooterValueColor", app.FooterValueColor() == th.Palette.Fg2},
		{"KeyGlowColor", app.KeyGlowColor() == th.KeyGlowColor()},
		{"notification info", overlays.NotificationConfig.LevelStyle(overlays.LevelInfo) == th.NotifInfo},
		{"notification error", overlays.NotificationConfig.LevelStyle(overlays.LevelError) == th.NotifError},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("UseTheme did not apply %s", c.name)
		}
	}
}

func TestNavBarModeString(t *testing.T) {
	cases := map[NavBarMode]string{
		NavBarCompact:  "Compact",
		NavBarExpanded: "Expanded",
		NavBarPath:     "Path",
		NavBarColumns:  "Columns",
		NavBarMode(99): "Compact",
	}
	for mode, want := range cases {
		if got := mode.String(); got != want {
			t.Errorf("NavBarMode(%d).String() = %q, want %q", int(mode), got, want)
		}
	}
	if NavBarColumns.next() != NavBarCompact {
		t.Error("next after Columns wraps to Compact")
	}
}

// TestUseThemeTakesFooterColorsFromThePalette checks that a palette's footer
// colors reach the footer: the key, label and glow colors are the palette's
// own, not the defaults they would cascade to.
func TestUseThemeTakesFooterColorsFromThePalette(t *testing.T) {
	p := styles.DraculaPalette
	p.FooterKeyColor = lipgloss.Color("#ff0001")
	p.FooterValueColor = lipgloss.Color("#ff0002")
	p.KeyPressHighlight = lipgloss.Color("#ff0003")
	app := New().UseTheme(styles.NewTheme(p))

	if got := app.FooterKeyColor(); got != p.FooterKeyColor {
		t.Errorf("FooterKeyColor = %v, want %v", got, p.FooterKeyColor)
	}
	if got := app.FooterValueColor(); got != p.FooterValueColor {
		t.Errorf("FooterValueColor = %v, want %v", got, p.FooterValueColor)
	}
	if got := app.KeyGlowColor(); got != p.KeyPressHighlight {
		t.Errorf("KeyGlowColor = %v, want %v", got, p.KeyPressHighlight)
	}

	// A palette that sets none of them gets the cascade's defaults.
	app.UseTheme(styles.Dracula)
	d := styles.Dracula.Palette
	if app.FooterKeyColor() != d.Primary || app.FooterValueColor() != d.Fg2 || app.KeyGlowColor() != d.GreenLt {
		t.Errorf("defaults = %v %v %v, want Primary, Fg2 and GreenLt", app.FooterKeyColor(), app.FooterValueColor(), app.KeyGlowColor())
	}
}
