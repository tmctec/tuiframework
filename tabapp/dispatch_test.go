package tabapp

import (
	"slices"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// runeKey presses a printable key and returns the command Update produced.
func runeKey(app *TabApp, r rune) tea.Cmd {
	_, cmd := app.Update(KeyPress(string(r)))
	return cmd
}

// markerMsg is an application message tabapp does not know.
type markerMsg struct{ n int }

// assertFrame checks the rendered frame is exactly w columns by h lines.
func assertFrame(t *testing.T, app *TabApp, w, h int) {
	t.Helper()
	lines := strings.Split(app.Render(), "\n")
	if len(lines) != h {
		t.Fatalf("%d lines, want %d", len(lines), h)
	}
	for i, l := range lines {
		if lw := lipgloss.Width(l); lw != w {
			t.Fatalf("line %d is %d wide, want %d", i, lw, w)
		}
	}
}

func TestNewDefaultModel(t *testing.T) {
	app := New()
	m := app.Model()

	checks := []struct {
		name string
		ok   bool
	}{
		{"NavBarPosition Top", m.NavBarPosition == NavBarTop && app.NavBarPosition() == NavBarTop},
		{"NavBarAutoHide on", m.NavBarAutoHide && app.NavBarAutoHide()},
		{"NavBarTimeout 2s", m.NavBarTimeout == 2*time.Second && app.NavBarAutoHideTimeout() == 2*time.Second},
		{"TabNavTimeout 2s", m.TabNavTimeout == 2*time.Second && app.TabNavKeyTimeout() == 2*time.Second},
		{"NavBarMode Compact", m.NavBarMode == NavBarCompact && app.NavBarMode() == NavBarCompact},
		{"Theme Dracula", m.Theme == styles.Dracula},
		{"NavBarBelowHeader on", m.NavBarBelowHeader && app.NavBarBelowHeader()},
		{"KeyGlowEnabled on", m.KeyGlowEnabled && app.KeyGlowEnabled()},
		{"TabNavActive off", !m.TabNavActive && !app.TabNavActive()},
		{"ActivePage 0", m.ActivePage == 0},
		{"OverlayCount 0", m.OverlayCount == 0},
		{"no borders", !app.ShowHeaderBorder() && !m.ShowNavBarBorder && !app.ShowFooterBorder() && !m.ShowContentBorder},
		{"no frames", app.HeaderFrame() == FrameNone && app.FooterFrame() == FrameNone},
		{"band heights", app.HeaderHeight() == 1 && app.NavBarHeight() == 3 && app.FooterHeight() == 3},
		{"NavBarWidth 20", app.NavBarWidth() == 20},
		{"KeyGlowDuration 300ms", app.KeyGlowDuration() == 300*time.Millisecond},
		{"KeyGlowColor from theme", app.KeyGlowColor() == styles.Dracula.KeyGlowColor()},
		{"no content", !app.HasContent() && !app.HasSections() && app.Root() == nil},
		{"NavPath [0]", slices.Equal(app.NavPath(), []int{0})},
		{"Breadcrumbs empty", len(app.Breadcrumbs()) == 0},
	}
	for _, c := range checks {
		if !c.ok {
			t.Errorf("default %s does not hold", c.name)
		}
	}
}

func TestHardcodedKeysQuit(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	for _, k := range []string{"ctrl+c", "ctrl+q"} {
		_, cmd := app.Update(KeyPress(k))
		if cmd == nil {
			t.Fatalf("%v should return a command", k)
		}
		if _, ok := cmd().(tea.QuitMsg); !ok {
			t.Errorf("%v produced %T, want tea.QuitMsg", k, cmd())
		}
	}
	_, cmd := app.Update(QuitMsg{})
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Errorf("QuitMsg produced %T, want tea.QuitMsg", cmd())
	}
	if _, cmd := app.Update(NoOpMsg{}); cmd != nil {
		t.Error("NoOpMsg should produce no command")
	}
}

func TestKeyDispatchOrder(t *testing.T) {
	app := New()
	app.MustBindKey("n", "Nav border", ToggleNavBarBorderMsg{})
	app.MustAddPage(page.New("One").BindKey("z", "Header border", ToggleHeaderBorderMsg{}))
	app.MustAddPage(page.New("Two"))

	if cmd := runeKey(app, 'n'); cmd == nil {
		t.Error("a universal key returns the glow command")
	}
	if !app.Model().ShowNavBarBorder {
		t.Error("the universal key's message should be applied")
	}
	if app.keyGlowKey != "n" {
		t.Errorf("glowing key = %q, want n", app.keyGlowKey)
	}

	if cmd := runeKey(app, 'z'); cmd == nil {
		t.Error("a page key returns the glow command")
	}
	if app.Model().HeaderFrame != FrameTheme {
		t.Error("the page key's message should be applied")
	}
	if app.keyGlowKey != "z" {
		t.Errorf("glowing key = %q, want z", app.keyGlowKey)
	}

	before := app.NavPath()
	if cmd := runeKey(app, 'w'); cmd != nil {
		t.Error("an unbound key returns no command")
	}
	if app.keyGlowKey != "z" || !slices.Equal(app.NavPath(), before) {
		t.Error("an unbound key must not change the state")
	}

	runeKey(app, '2')
	if app.Active() != 1 {
		t.Errorf("Active = %d after '2', want 1", app.Active())
	}
	runeKey(app, '9')
	runeKey(app, '0')
	if app.Active() != 1 {
		t.Errorf("Active = %d after out-of-range digits, want 1", app.Active())
	}
}

func TestUpdateReturnsTheSameModel(t *testing.T) {
	app := New()
	m, _ := app.Update(NoOpMsg{})
	if m.(*TabApp) != app {
		t.Error("Update should return the receiver")
	}
}

func TestEscLeavesTabNavigation(t *testing.T) {
	app := rootMixApp() // autohide on
	key(app, "tab")
	if !app.TabNavActive() || app.navBarHidden {
		t.Fatalf("after tab: active %v hidden %v, want true and false", app.TabNavActive(), app.navBarHidden)
	}

	key(app, "esc")

	if app.TabNavActive() {
		t.Error("esc should end tab navigation")
	}
	if !app.navBarHidden {
		t.Error("leaving navigation with autohide on hides the bar at once")
	}
	if app.navFocusCrumb != -1 {
		t.Errorf("crumb focus = %d, want -1", app.navFocusCrumb)
	}

	app.SetNavBarAutoHide(false)
	key(app, "tab")
	key(app, "esc")
	if app.navBarHidden {
		t.Error("with autohide off the bar stays visible after esc")
	}
}

func TestAutohideTimerIgnoresStaleID(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	if cmd := app.Init(); cmd == nil {
		t.Fatal("Init should arm the autohide timer")
	}

	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID - 1})
	if app.navBarHidden {
		t.Error("a stale autohide timer must not hide the bar")
	}
	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID})
	if !app.navBarHidden {
		t.Error("the current autohide timer hides the bar")
	}

	app.SetNavBarAutoHide(false)
	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID})
	if app.navBarHidden {
		t.Error("with autohide off the timer message is ignored")
	}
	if app.startAutohideTimer() != nil {
		t.Error("no autohide timer is armed while autohide is off")
	}
}

func TestTabNavTimeoutIgnoresStaleID(t *testing.T) {
	app := stripApp()
	key(app, "tab")
	key(app, "left")
	key(app, "left") // onto crumb A
	if !app.TabNavActive() || app.navFocusCrumb != 0 {
		t.Fatalf("precondition: active %v focus %d", app.TabNavActive(), app.navFocusCrumb)
	}

	app.Update(tabNavTimeoutMsg{id: app.tabNavTimerID - 1})
	if !app.TabNavActive() {
		t.Error("a stale navigation timeout must not end navigation")
	}

	app.Update(tabNavTimeoutMsg{id: app.tabNavTimerID})
	if app.TabNavActive() || app.Model().TabNavActive {
		t.Error("the current navigation timeout ends navigation")
	}
	if app.navFocusCrumb != -1 {
		t.Errorf("crumb focus = %d after timeout, want -1", app.navFocusCrumb)
	}
}

func TestKeyGlowTimeoutClearsOnlyTheCurrentGlow(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.MustBindKey("x", "Mark", NoOpMsg{})

	if cmd := runeKey(app, 'x'); cmd == nil {
		t.Fatal("a universal key returns the glow command while glow is on")
	}
	if app.keyGlowKey != "x" {
		t.Fatalf("glowing key = %q, want x", app.keyGlowKey)
	}

	app.Update(keyGlowTimeoutMsg{id: app.keyGlowTimerID - 1})
	if app.keyGlowKey != "x" {
		t.Error("a stale glow timer must not clear the glow")
	}
	app.Update(keyGlowTimeoutMsg{id: app.keyGlowTimerID})
	if app.keyGlowKey != "" {
		t.Errorf("glowing key = %q after timeout, want none", app.keyGlowKey)
	}

	app.SetKeyGlowEnabled(false)
	if cmd := runeKey(app, 'x'); cmd != nil {
		t.Error("with glow off a NoOp universal key returns no command")
	}
	if app.keyGlowKey != "" {
		t.Error("with glow off no key glows")
	}
}

func TestWindowSizeSetsTheFrame(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Only"))
	app.SetNavBarAutoHide(false)

	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	assertFrame(t, app, 80, 24)

	app.SetWidth(60).SetHeight(20)
	assertFrame(t, app, 60, 20)
}

func TestCtrlArrowsSwitchSectionsInExpandedMode(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarAutoHide(false)
	key(app, "tab")       // a2
	key(app, "backspace") // strip: up to the root level
	app.Update(tabNavTimeoutMsg{id: app.tabNavTimerID})
	if app.TabNavActive() || app.tabNavLevel != 0 {
		t.Fatalf("precondition: active %v level %d", app.TabNavActive(), app.tabNavLevel)
	}

	key(app, "ctrl+right")
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("compact mode: ctrl+right changed the path to %v", got)
	}

	app.SetNavBarMode(NavBarExpanded)
	key(app, "tab") // level 1: a1
	key(app, "up")  // level 0
	app.Update(tabNavTimeoutMsg{id: app.tabNavTimerID})
	if app.tabNavLevel != 0 {
		t.Fatalf("level = %d, want 0", app.tabNavLevel)
	}

	steps := []struct {
		k    string
		want []int
		page string
	}{
		{"ctrl+right", []int{1, 0}, "b1"},
		{"ctrl+l", []int{2, 0}, "Settings"},
		{"ctrl+h", []int{1, 0}, "b1"},
		{"ctrl+left", []int{0, 0}, "a1"},
	}
	for _, s := range steps {
		key(app, s.k)
		if got := app.NavPath(); !slices.Equal(got, s.want) {
			t.Errorf("after %v: path %v, want %v", s.k, got, s.want)
		}
		if got := app.currentPage().Name(); got != s.page {
			t.Errorf("after %v: page %q, want %q", s.k, got, s.page)
		}
	}

	flat := New()
	flat.MustAddPage(page.New("p1")).MustAddPage(page.New("p2"))
	flat.SetNavBarMode(NavBarExpanded)
	key(flat, "ctrl+right")
	if flat.Active() != 0 {
		t.Error("without sections ctrl+right does nothing")
	}
}

func TestTabNavKeysAreConfigurable(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("p1")).MustAddPage(page.New("p2")).MustAddPage(page.New("p3"))
	app.SetTabNavKeyForward("n").SetTabNavKeyReverse("p")

	runeKey(app, 'n')
	if app.Active() != 1 || !app.TabNavActive() {
		t.Fatalf("after n: active %d nav %v, want 1 and true", app.Active(), app.TabNavActive())
	}
	runeKey(app, 'p')
	runeKey(app, 'p')
	if app.Active() != 2 {
		t.Errorf("after p p: active %d, want 2 (wrapped)", app.Active())
	}
	key(app, "tab")
	if app.Active() != 2 {
		t.Error("tab is an ordinary universal key once the navigation key changes")
	}
}

func TestCustomHandlerResultIsHonoured(t *testing.T) {
	app := New()
	var seen []tea.Msg
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		seen = append(seen, msg)
		if m, ok := msg.(markerMsg); ok && m.n > 0 {
			return true, func() tea.Msg { return markerMsg{n: m.n + 1} }
		}
		return false, nil
	})

	_, cmd := app.Update(markerMsg{n: 1})
	if cmd == nil {
		t.Fatal("a handled message returns the handler's command")
	}
	if got, ok := cmd().(markerMsg); !ok || got.n != 2 {
		t.Errorf("command produced %#v, want markerMsg{n: 2}", cmd())
	}

	if _, cmd := app.Update(markerMsg{n: 0}); cmd != nil {
		t.Error("an unhandled message yields no command")
	}
	if len(seen) != 2 {
		t.Errorf("handler called %d times, want 2", len(seen))
	}

	if _, cmd := New().Update(markerMsg{n: 1}); cmd != nil {
		t.Error("without a handler an unknown message yields no command")
	}
}

// Settings messages update the TabApp field and its Model mirror together.

func TestSettingMessages(t *testing.T) {
	app := New()

	app.Update(ToggleNavBarPositionMsg{LeftWidth: 30})
	if app.NavBarPosition() != NavBarLeft || app.Model().NavBarPosition != NavBarLeft || app.NavBarWidth() != 30 {
		t.Errorf("toggle to left: position %v width %d, want Left and 30", app.NavBarPosition(), app.NavBarWidth())
	}
	app.Update(ToggleNavBarPositionMsg{LeftWidth: 44})
	if app.NavBarPosition() != NavBarTop || app.NavBarWidth() != 30 {
		t.Errorf("toggle to top: position %v width %d, want Top and 30", app.NavBarPosition(), app.NavBarWidth())
	}
	app.Update(ToggleNavBarPositionMsg{})
	if app.NavBarPosition() != NavBarLeft || app.NavBarWidth() != 30 {
		t.Error("toggling without LeftWidth keeps the width")
	}

	app.Update(SetNavBarWidthMsg{Width: 25})
	app.Update(SetNavBarWidthMsg{Width: 0})
	app.Update(SetNavBarWidthMsg{Width: -3})
	if app.NavBarWidth() != 25 {
		t.Errorf("width = %d, want 25 (non-positive widths ignored)", app.NavBarWidth())
	}
	app.Update(SetNavBarWidthMsg{Width: 4})
	if app.NavBarWidth() != minNavBarWidth {
		t.Errorf("width = %d, want the %d minimum for a narrow request", app.NavBarWidth(), minNavBarWidth)
	}
	app.Update(ToggleNavBarPositionMsg{}) // to the top
	app.Update(ToggleNavBarPositionMsg{LeftWidth: 7})
	if app.NavBarPosition() != NavBarLeft || app.NavBarWidth() != minNavBarWidth {
		t.Errorf("toggle with a narrow LeftWidth: position %v width %d, want Left and %d", app.NavBarPosition(), app.NavBarWidth(), minNavBarWidth)
	}

	app.Update(SetNavBarTimeoutMsg{Duration: 5 * time.Second})
	if app.NavBarAutoHideTimeout() != 5*time.Second || app.Model().NavBarTimeout != 5*time.Second {
		t.Errorf("nav bar timeout = %v, want 5s", app.NavBarAutoHideTimeout())
	}
	app.Update(AdjustNavBarTimeoutMsg{Delta: -time.Hour})
	if app.Model().NavBarTimeout != time.Second {
		t.Errorf("nav bar timeout = %v, want the 1s floor", app.Model().NavBarTimeout)
	}
	app.Update(AdjustNavBarTimeoutMsg{Delta: 500 * time.Millisecond})
	if app.Model().NavBarTimeout != 1500*time.Millisecond {
		t.Errorf("nav bar timeout = %v, want 1.5s", app.Model().NavBarTimeout)
	}

	app.Update(SetTabNavTimeoutMsg{Duration: 3 * time.Second})
	if app.TabNavKeyTimeout() != 3*time.Second || app.Model().TabNavTimeout != 3*time.Second {
		t.Errorf("tab nav timeout = %v, want 3s", app.TabNavKeyTimeout())
	}
	app.Update(AdjustTabNavTimeoutMsg{Delta: time.Hour})
	if app.Model().TabNavTimeout != 10*time.Second {
		t.Errorf("tab nav timeout = %v, want the 10s ceiling", app.Model().TabNavTimeout)
	}
	app.Update(AdjustTabNavTimeoutMsg{Delta: -time.Second})
	if app.Model().TabNavTimeout != 9*time.Second {
		t.Errorf("tab nav timeout = %v, want 9s", app.Model().TabNavTimeout)
	}

	app.Update(SetNavBarModeMsg{Mode: NavBarColumns})
	if app.NavBarMode() != NavBarColumns || app.Model().NavBarMode != NavBarColumns {
		t.Errorf("mode = %v, want Columns", app.NavBarMode())
	}

	toggles := []struct {
		name string
		msg  tea.Msg
		get  func() bool
	}{
		{"header border", ToggleHeaderBorderMsg{}, app.ShowHeaderBorder},
		{"nav border", ToggleNavBarBorderMsg{}, app.ShowNavBarBorder},
		{"footer border", ToggleFooterBorderMsg{}, app.ShowFooterBorder},
		{"content border", ToggleContentBorderMsg{}, app.ShowContentBorder},
	}
	for _, tg := range toggles {
		app.Update(tg.msg)
		if !tg.get() {
			t.Errorf("%s should be on after one toggle", tg.name)
		}
		app.Update(tg.msg)
		if tg.get() {
			t.Errorf("%s should be off after two toggles", tg.name)
		}
	}

	app.Update(ToggleNavBarBelowHeaderMsg{})
	if app.NavBarBelowHeader() || app.Model().NavBarBelowHeader {
		t.Error("nav below header should be off after one toggle")
	}
	app.Update(ToggleNavBarBelowHeaderMsg{})
	if !app.NavBarBelowHeader() {
		t.Error("nav below header should be on after two toggles")
	}

	app.keyGlowKey = "k"
	app.Update(ToggleKeyGlowMsg{})
	if app.KeyGlowEnabled() || app.keyGlowKey != "" {
		t.Error("disabling key glow clears the glowing key")
	}
	app.Update(ToggleKeyGlowMsg{})
	if !app.KeyGlowEnabled() || !app.Model().KeyGlowEnabled {
		t.Error("key glow should be on after two toggles")
	}
}

func TestToggleAutohideMessages(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.Init()
	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID})
	if !app.navBarHidden {
		t.Fatal("precondition: the bar should be hidden")
	}

	_, cmd := app.Update(ToggleNavBarAutoHideMsg{})
	if app.NavBarAutoHide() || app.Model().NavBarAutoHide {
		t.Error("toggle should switch autohide off")
	}
	if app.navBarHidden {
		t.Error("switching autohide off shows the bar")
	}
	if cmd != nil {
		t.Error("switching autohide off arms no timer")
	}

	_, cmd = app.Update(ToggleNavBarAutoHideMsg{})
	if !app.NavBarAutoHide() || cmd == nil {
		t.Error("switching autohide on arms the timer")
	}
}

func TestThemeMessages(t *testing.T) {
	defer overlays.NotificationConfig.ApplyTheme(styles.Dracula)
	app := New()

	_, cmd := app.Update(SetThemeMsg{Theme: styles.Nord})
	if cmd != nil {
		t.Error("SetThemeMsg yields no command")
	}
	if app.Model().Theme != styles.Nord || app.TabStyle().GetForeground() != styles.Nord.TabInactive.GetForeground() {
		t.Error("SetThemeMsg should apply the theme")
	}

	themes := []*styles.Theme{styles.Dracula, styles.Nord, styles.Gruvbox}
	steps := []struct {
		dir  int
		want *styles.Theme
		idx  int
	}{
		{1, styles.Gruvbox, 2},
		{1, styles.Dracula, 0},
		{-1, styles.Gruvbox, 2},
		{-1, styles.Nord, 1},
	}
	for _, s := range steps {
		app.Update(CycleThemeMsg{Direction: s.dir, Themes: themes})
		if app.Model().Theme != s.want || s.want != themes[s.idx] {
			t.Errorf("direction %d: theme %s, want %s at index %d", s.dir, app.Model().Theme.Name(), s.want.Name(), s.idx)
		}
	}

	app.UseTheme(styles.Monokai) // not in the list
	app.Update(CycleThemeMsg{Direction: -1, Themes: themes})
	if app.Model().Theme != styles.Gruvbox {
		t.Errorf("unknown theme backward: %s, want the last (Gruvbox)", app.Model().Theme.Name())
	}
	app.UseTheme(styles.Monokai)
	app.Update(CycleThemeMsg{Direction: 1, Themes: themes})
	if app.Model().Theme != styles.Dracula {
		t.Errorf("unknown theme forward: %s, want the first (Dracula)", app.Model().Theme.Name())
	}

	app.Update(CycleThemeMsg{Direction: 1, Themes: nil})
	if app.Model().Theme != styles.Dracula {
		t.Error("an empty theme list is a no-op")
	}

	cp := *styles.Nord
	app.UseTheme(&cp)
	app.Update(CycleThemeMsg{Direction: 1, Themes: themes})
	if app.Model().Theme != styles.Gruvbox {
		t.Errorf("a copy matches by name: got %s, want Gruvbox", app.Model().Theme.Name())
	}
	if themeIndex(themes, nil) != -1 {
		t.Error("themeIndex of nil should be -1")
	}
}

func TestNavigationMessages(t *testing.T) {
	app := rootMixApp()

	app.Update(NavigateToPageMsg{Index: 1})
	if app.currentPage().Name() != "a2" || app.Model().ActivePage != 1 {
		t.Fatalf("page %q active %d, want a2 and 1", app.currentPage().Name(), app.Model().ActivePage)
	}
	app.Update(NavigateToPageMsg{Index: 7})
	app.Update(NavigateToPageMsg{Index: -1})
	if app.currentPage().Name() != "a2" {
		t.Error("out-of-range page indices leave the position alone")
	}
	app.Update(NavigateNextPageMsg{})
	if app.currentPage().Name() != "a1" {
		t.Errorf("next from the last page: %q, want a1", app.currentPage().Name())
	}
	app.Update(NavigatePrevPageMsg{})
	if app.currentPage().Name() != "a2" {
		t.Errorf("prev from the first page: %q, want a2", app.currentPage().Name())
	}

	// A fresh app focuses level 0, a section level, so the section messages act.
	app.Update(NavigateNextSectionMsg{})
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) {
		t.Errorf("next section: path %v, want [1 0]", got)
	}
	app.Update(NavigatePrevSectionMsg{})
	app.Update(NavigatePrevSectionMsg{})
	if got := app.NavPath(); !slices.Equal(got, []int{2, 0}) || app.currentPage().Name() != "Settings" {
		t.Errorf("prev section twice: path %v page %q, want [2 0] Settings", got, app.currentPage().Name())
	}

	// At the page level the section messages do nothing.
	if err := app.NavigateToPath([]int{0, 1}); err != nil {
		t.Fatal(err)
	}
	if app.TabNavLevel() != 1 {
		t.Fatalf("level = %d, want 1", app.TabNavLevel())
	}
	app.Update(NavigateNextSectionMsg{})
	app.Update(NavigatePrevSectionMsg{})
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("section messages at the page level moved to %v", got)
	}

	app.Update(EnterSectionMsg{Level: 0, Index: 1})
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || len(app.Breadcrumbs()) != 1 {
		t.Fatalf("enter section: path %v crumbs %d, want [1 0] and 1", got, len(app.Breadcrumbs()))
	}
	app.Update(EnterSectionMsg{Level: 3, Index: 0})
	app.Update(PopToBreadcrumbMsg{BreadcrumbIndex: 4})
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || len(app.Breadcrumbs()) != 1 {
		t.Errorf("invalid enter/pop: path %v crumbs %d, want unchanged", got, len(app.Breadcrumbs()))
	}
	app.Update(PopToBreadcrumbMsg{BreadcrumbIndex: 0})
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) || len(app.Breadcrumbs()) != 0 {
		t.Errorf("pop: path %v crumbs %d, want [0 1] and 0", got, len(app.Breadcrumbs()))
	}
	if app.currentPage().Name() != "a2" || app.Model().ActivePage != 1 {
		t.Errorf("after pop: page %q active %d, want a2 and 1", app.currentPage().Name(), app.Model().ActivePage)
	}
}

// deliver runs a command and feeds what it produces back into Update,
// opening batches so every timer inside them fires.
func deliver(app *TabApp, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	msg := cmd()
	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			deliver(app, c)
		}
		return
	}
	app.Update(msg)
}

func TestTimerCommandsDeliverTheirMessages(t *testing.T) {
	app := rootMixApp()
	app.SetKeyGlowDuration(time.Millisecond).
		SetTabNavKeyTimeout(time.Millisecond).
		SetNavBarAutoHideTimeout(time.Millisecond)
	app.MustBindKey("x", "Mark", NoOpMsg{})

	deliver(app, runeKey(app, 'x'))
	if app.keyGlowKey != "" {
		t.Errorf("glowing key = %q after the glow timer fired, want none", app.keyGlowKey)
	}

	_, cmd := app.Update(KeyPress("tab"))
	if !app.TabNavActive() || app.navBarHidden {
		t.Fatalf("precondition: active %v hidden %v", app.TabNavActive(), app.navBarHidden)
	}
	deliver(app, cmd)
	if app.TabNavActive() {
		t.Error("the navigation timer should end tab navigation")
	}
	if !app.navBarHidden {
		t.Error("the autohide timer should hide the bar")
	}
}

func TestEmptyAppIgnoresNavigation(t *testing.T) {
	app := New()
	app.Update(KeyPress("tab"))
	app.Update(KeyPress("shift+tab"))
	runeKey(app, '1')
	for _, msg := range []tea.Msg{
		NavigateToPageMsg{Index: 0},
		NavigateNextPageMsg{},
		NavigatePrevPageMsg{},
		NavigateNextSectionMsg{},
		NavigatePrevSectionMsg{},
		EnterSectionMsg{},
		PopToBreadcrumbMsg{},
		NavigateToPathMsg{Path: []int{0}},
	} {
		app.Update(msg)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0}) {
		t.Errorf("NavPath = %v, want [0]", got)
	}
	if app.currentPage() != nil || app.Model().ActivePage != 0 {
		t.Error("an empty app has no page")
	}
}
