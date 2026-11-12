package main

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/tabapp"
)

func TestMoveCursorWraps(t *testing.T) {
	n := int(rowCount)
	if got := moveCursor(0, -1, n); got != n-1 {
		t.Errorf("moveCursor(0,-1) = %d, want %d", got, n-1)
	}
	if got := moveCursor(n-1, 1, n); got != 0 {
		t.Errorf("moveCursor(n-1,+1) = %d, want 0", got)
	}
}

func TestAdjustThemeWrapsAndEmitsSetTheme(t *testing.T) {
	s := defaultSettings
	s.ThemeIndex = len(demoThemes) - 1

	after, msgs := adjustRow(s, rowTheme, 1)

	if after.ThemeIndex != 0 {
		t.Errorf("ThemeIndex = %d, want 0", after.ThemeIndex)
	}
	if len(msgs) != 1 {
		t.Fatalf("got %d msgs, want 1", len(msgs))
	}
	set, ok := msgs[0].(tabapp.SetThemeMsg)
	if !ok || set.Theme != demoThemes[0] {
		t.Errorf("msg = %#v, want SetThemeMsg for demoThemes[0]", msgs[0])
	}
}

func TestAdjustPositionAndModeEmitLibraryMsgs(t *testing.T) {
	after, msgs := adjustRow(defaultSettings, rowPosition, 1)
	if after.Position != tabapp.NavBarLeft {
		t.Errorf("Position = %v, want Left", after.Position)
	}
	if _, ok := msgs[0].(tabapp.ToggleNavBarPositionMsg); !ok {
		t.Errorf("msg = %#v, want ToggleNavBarPositionMsg", msgs[0])
	}

	after, msgs = adjustRow(defaultSettings, rowMode, 1)
	if after.Mode != tabapp.NavBarExpanded {
		t.Errorf("Mode = %v, want Expanded", after.Mode)
	}
	if set, ok := msgs[0].(tabapp.SetNavBarModeMsg); !ok || set.Mode != tabapp.NavBarExpanded {
		t.Errorf("msg = %#v, want SetNavBarModeMsg{Expanded}", msgs[0])
	}
	if after, _ = adjustRow(defaultSettings, rowMode, -1); after.Mode != tabapp.NavBarColumns {
		t.Errorf("Mode = %v, want Columns (the picker wraps)", after.Mode)
	}
}

func TestAdjustTimeoutsUseLibraryDeltas(t *testing.T) {
	after, msgs := adjustRow(defaultSettings, rowAutoHideTimeout, 1)
	if after.AutoHideTimeout != 3*time.Second {
		t.Errorf("AutoHideTimeout = %v, want 3s", after.AutoHideTimeout)
	}
	if adj, ok := msgs[0].(tabapp.AdjustNavBarTimeoutMsg); !ok || adj.Delta != time.Second {
		t.Errorf("msg = %#v, want AdjustNavBarTimeoutMsg{+1s}", msgs[0])
	}

	after, msgs = adjustRow(defaultSettings, rowTabNavTimeout, -1)
	if after.TabNavTimeout != 1500*time.Millisecond {
		t.Errorf("TabNavTimeout = %v, want 1.5s", after.TabNavTimeout)
	}
	if adj, ok := msgs[0].(tabapp.AdjustTabNavTimeoutMsg); !ok || adj.Delta != -500*time.Millisecond {
		t.Errorf("msg = %#v, want AdjustTabNavTimeoutMsg{-500ms}", msgs[0])
	}
}

func TestAdjustTimeoutsClampSnapshot(t *testing.T) {
	s := defaultSettings
	s.AutoHideTimeout = autoHideMax
	after, _ := adjustRow(s, rowAutoHideTimeout, 1)
	if after.AutoHideTimeout != autoHideMax {
		t.Errorf("AutoHideTimeout = %v, want clamp at %v", after.AutoHideTimeout, autoHideMax)
	}
	s.TabNavTimeout = tabNavMin
	after, _ = adjustRow(s, rowTabNavTimeout, -1)
	if after.TabNavTimeout != tabNavMin {
		t.Errorf("TabNavTimeout = %v, want clamp at %v", after.TabNavTimeout, tabNavMin)
	}
}

func TestChromeRowsEmitToggles(t *testing.T) {
	cases := []struct {
		row  settingsRow
		want tea.Msg
	}{
		{rowNavOrder, tabapp.ToggleNavBarBelowHeaderMsg{}},
		{rowAutoHide, tabapp.ToggleNavBarAutoHideMsg{}},
		{rowNavBorder, tabapp.ToggleNavBarBorderMsg{}},
		{rowContentBorder, tabapp.ToggleContentBorderMsg{}},
		{rowKeyGlow, tabapp.ToggleKeyGlowMsg{}},
	}
	for _, c := range cases {
		_, msgs := adjustRow(defaultSettings, c.row, 1)
		if len(msgs) != 1 || msgs[0] != c.want {
			t.Errorf("row %d: msgs = %#v, want [%T]", c.row, msgs, c.want)
		}
	}
}

func TestBandHeightRowsSetAndClamp(t *testing.T) {
	cases := []struct {
		row   settingsRow
		start int
		delta int
		want  int
		msg   tea.Msg
	}{
		{rowHeaderHeight, 1, 1, 2, tabapp.SetHeaderHeightMsg{Height: 2}},
		{rowHeaderHeight, 1, -1, 0, tabapp.SetHeaderHeightMsg{Height: 0}},
		{rowHeaderHeight, 0, -1, 0, tabapp.SetHeaderHeightMsg{Height: 0}},
		{rowFooterHeight, 2, 1, 3, tabapp.SetFooterHeightMsg{Height: 3}},
		{rowFooterHeight, bandHeightMax, 1, bandHeightMax, tabapp.SetFooterHeightMsg{Height: bandHeightMax}},
	}
	for _, c := range cases {
		s := defaultSettings
		s.HeaderHeight, s.FooterHeight = c.start, c.start
		after, msgs := adjustRow(s, c.row, c.delta)
		got := after.HeaderHeight
		if c.row == rowFooterHeight {
			got = after.FooterHeight
		}
		if got != c.want {
			t.Errorf("row %d from %d by %+d: height = %d, want %d", c.row, c.start, c.delta, got, c.want)
		}
		if len(msgs) != 1 || msgs[0] != c.msg {
			t.Errorf("row %d from %d by %+d: msgs = %#v, want [%#v]", c.row, c.start, c.delta, msgs, c.msg)
		}
	}
}

func TestBandFrameRowsCycleThePicker(t *testing.T) {
	// Right walks None, Single line, Double line, Fill and wraps to None.
	s := defaultSettings
	want := []tabapp.Frame{tabapp.FrameSingle, tabapp.FrameDouble, tabapp.FrameFill, tabapp.FrameNone}
	for _, frame := range want {
		var msgs []tea.Msg
		s, msgs = adjustRow(s, rowHeaderFrame, 1)
		if s.HeaderFrame != frame {
			t.Fatalf("HeaderFrame = %v, want %v", s.HeaderFrame, frame)
		}
		if len(msgs) != 1 || msgs[0] != (tabapp.SetHeaderFrameMsg{Frame: frame}) {
			t.Errorf("msgs = %#v, want SetHeaderFrameMsg{%v}", msgs, frame)
		}
	}

	// Left from None wraps to Fill.
	after, msgs := adjustRow(defaultSettings, rowFooterFrame, -1)
	if after.FooterFrame != tabapp.FrameFill {
		t.Errorf("FooterFrame = %v, want Fill", after.FooterFrame)
	}
	if len(msgs) != 1 || msgs[0] != (tabapp.SetFooterFrameMsg{Frame: tabapp.FrameFill}) {
		t.Errorf("msgs = %#v, want SetFooterFrameMsg{Fill}", msgs)
	}

	// A frame the picker does not offer counts as the first entry.
	s = defaultSettings
	s.FooterFrame = tabapp.FrameTheme
	if after, _ = adjustRow(s, rowFooterFrame, 1); after.FooterFrame != tabapp.FrameSingle {
		t.Errorf("FooterFrame after Theme = %v, want Single", after.FooterFrame)
	}
}

func TestResetDiff(t *testing.T) {
	if msgs := resetMsgs(defaultSettings, defaultSettings); len(msgs) != 0 {
		t.Errorf("reset on defaults emitted %d msgs", len(msgs))
	}

	changed := settingsSnapshot{
		ThemeIndex:        3,
		Position:          tabapp.NavBarLeft,
		NavBarBelowHeader: false,
		Mode:              tabapp.NavBarExpanded,
		AutoHide:          false,
		AutoHideTimeout:   9 * time.Second,
		TabNavTimeout:     4 * time.Second,
		HeaderHeight:      3,
		HeaderFrame:       tabapp.FrameFill,
		NavBorder:         true,
		ContentBorder:     true,
		FooterHeight:      4,
		FooterFrame:       tabapp.FrameDouble,
		KeyGlow:           false,
		TerminalColors:    false,
		FollowBackground:  true,
		DarkBackground:    true,
	}
	msgs := resetMsgs(changed, defaultSettings)
	if len(msgs) != int(rowCount) {
		t.Fatalf("got %d msgs, want one per row (%d)", len(msgs), rowCount)
	}
	if set, ok := msgs[0].(tabapp.SetThemeMsg); !ok || set.Theme != demoThemes[0] {
		t.Errorf("first msg = %#v, want SetThemeMsg for demoThemes[0]", msgs[0])
	}
	if set, ok := msgs[5].(tabapp.SetNavBarTimeoutMsg); !ok || set.Duration != 2*time.Second {
		t.Errorf("timeout msg = %#v, want SetNavBarTimeoutMsg{2s}", msgs[5])
	}
	for _, want := range []tea.Msg{
		tabapp.SetHeaderHeightMsg{Height: headerHeight},
		tabapp.SetHeaderFrameMsg{Frame: tabapp.FrameNone},
		tabapp.SetFooterHeightMsg{Height: footerHeight},
		tabapp.SetFooterFrameMsg{Frame: tabapp.FrameNone},
	} {
		if !slices.Contains(msgs, want) {
			t.Errorf("reset should emit %#v", want)
		}
	}
}

func TestNavOrderRowIgnoredForSidebar(t *testing.T) {
	s := defaultSettings
	s.Position = tabapp.NavBarLeft
	after, msgs := adjustRow(s, rowNavOrder, 1)
	if after != s || len(msgs) != 0 {
		t.Errorf("sidebar nav order should be a no-op, got %+v %v", after, msgs)
	}
	if got := rowValue(s, rowNavOrder); got != "Sidebar" {
		t.Errorf("value = %q, want Sidebar", got)
	}
}

func TestRowValueFormatting(t *testing.T) {
	s := defaultSettings
	s.ThemeIndex = 5
	if got := rowValue(s, rowTheme); got != demoThemes[5].Name()+" (6/"+itoa(len(demoThemes))+")" {
		t.Errorf("theme value = %q", got)
	}
	if got := rowValue(s, rowAutoHideTimeout); got != "2.0 s" {
		t.Errorf("timeout value = %q, want 2.0 s", got)
	}
	if got := rowValue(s, rowKeyGlow); got != "On" {
		t.Errorf("key glow value = %q, want On", got)
	}

	heights := map[int]string{0: "Hidden", 1: "1 line", 3: "3 lines"}
	for h, want := range heights {
		s.HeaderHeight, s.FooterHeight = h, h
		if got := rowValue(s, rowHeaderHeight); got != want {
			t.Errorf("header height %d value = %q, want %q", h, got, want)
		}
		if got := rowValue(s, rowFooterHeight); got != want {
			t.Errorf("footer height %d value = %q, want %q", h, got, want)
		}
	}
	frames := map[tabapp.Frame]string{
		tabapp.FrameNone:   "None",
		tabapp.FrameSingle: "Single line",
		tabapp.FrameDouble: "Double line",
		tabapp.FrameFill:   "Fill",
	}
	for f, want := range frames {
		s.HeaderFrame, s.FooterFrame = f, f
		if got := rowValue(s, rowHeaderFrame); got != want {
			t.Errorf("header frame %v value = %q, want %q", f, got, want)
		}
		if got := rowValue(s, rowFooterFrame); got != want {
			t.Errorf("footer frame %v value = %q, want %q", f, got, want)
		}
	}
}

// TestRowsHaveLabelsAndFittingDescriptions checks every row is named and
// that its description fits the popup's description lines without clipping.
func TestRowsHaveLabelsAndFittingDescriptions(t *testing.T) {
	wrap := lipgloss.NewStyle().Width(settingsTextWidth)
	for row := settingsRow(0); row < rowCount; row++ {
		if rowLabel(row) == "" || rowValue(defaultSettings, row) == "" {
			t.Errorf("row %d: label %q value %q, want both set", row, rowLabel(row), rowValue(defaultSettings, row))
		}
		if w := lipgloss.Width(rowLabel(row)); w > settingsName {
			t.Errorf("row %d label %q is %d wide, want at most %d", row, rowLabel(row), w, settingsName)
		}
		desc := rowDescription(row)
		if lines := lipgloss.Height(wrap.Render(desc)); desc == "" || lines > settingsDescLines {
			t.Errorf("row %d description takes %d lines, want 1 to %d: %q", row, lines, settingsDescLines, desc)
		}
	}
}

func itoa(n int) string {
	return string(rune('0'+n/10)) + string(rune('0'+n%10))
}

func TestDemoThemesUniqueWithCyberpunkLast(t *testing.T) {
	seen := map[string]bool{}
	for _, th := range demoThemes {
		if seen[th.Name()] {
			t.Errorf("duplicate theme %s", th.Name())
		}
		seen[th.Name()] = true
	}
	if demoThemes[cyberpunkIndex()] != CyberpunkTheme {
		t.Error("Cyberpunk should be the last demo theme")
	}
}

func TestBuildApp(t *testing.T) {
	app := buildApp() // panics on any page key colliding with s, q or tab

	if app.Root().Depth() != 4 {
		t.Errorf("depth = %d, want 4", app.Root().Depth())
	}
	if n := len(app.Root().Sections()); n != 2 {
		t.Errorf("root sections = %d, want 2 (Showcase, Catalog)", n)
	}
	if n := len(app.Root().Pages()); n != 0 {
		t.Errorf("root pages = %d, want 0: Settings is a popup, not a tab", n)
	}
	if themeIndex(app.Model()) != 0 {
		t.Errorf("start theme index = %d, want 0", themeIndex(app.Model()))
	}
}

// drain runs cmd and feeds every message it produces back into app.Update,
// the way Bubble Tea's runtime would: the commands of a batch each in turn,
// and those of a sequence in order. A command that does not return at once
// is a timer (autohide, key glow, notification dismissal) and is left
// pending, as it would be for the moment the test looks at.
func drain(app *tabapp.TabApp, cmd tea.Cmd) {
	if cmd == nil {
		return
	}
	result := make(chan tea.Msg, 1)
	go func() { result <- cmd() }()
	var msg tea.Msg
	select {
	case msg = <-result:
	case <-time.After(50 * time.Millisecond):
		return
	}

	if batch, ok := msg.(tea.BatchMsg); ok {
		for _, c := range batch {
			drain(app, c)
		}
		return
	}
	// tea.Sequence delivers its commands as an unexported slice type.
	if v := reflect.ValueOf(msg); v.IsValid() && v.Kind() == reflect.Slice && v.Type().Elem() == reflect.TypeOf(tea.Cmd(nil)) {
		for i := 0; i < v.Len(); i++ {
			drain(app, v.Index(i).Interface().(tea.Cmd))
		}
		return
	}
	_, next := app.Update(msg)
	drain(app, next)
}

func press(app *tabapp.TabApp, key tea.KeyMsg) {
	_, cmd := app.Update(key)
	drain(app, cmd)
}

func TestSettingsPopupRoundTrip(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if err := app.NavigateToPath([]int{1, 0, 0, 2}); err != nil { // deep in the Catalog
		t.Fatal(err)
	}
	before := app.NavPath()

	// s opens the popup over the current page.
	press(app, tabapp.KeyPress("s"))
	if app.Model().OverlayCount != 1 || getState(app.Model()).SettingsID == 0 {
		t.Fatalf("after s: overlays=%d id=%d, want the popup open", app.Model().OverlayCount, getState(app.Model()).SettingsID)
	}
	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	if !strings.Contains(plain, "Settings") || !strings.Contains(plain, "Nav bar position") {
		t.Error("the popup should render its title and rows")
	}

	// Right on the Theme row changes the theme; the popup frame follows it.
	press(app, tabapp.KeyPress("right"))
	if themeIndex(app.Model()) != 1 {
		t.Errorf("theme index after Right = %d, want 1", themeIndex(app.Model()))
	}
	if got := app.TopOverlay().Style().GetBorderTopForeground(); got != demoThemes[1].ModalBorderColor {
		t.Errorf("popup border = %v, want %v from the new theme", got, demoThemes[1].ModalBorderColor)
	}

	// Esc closes the popup inside the library; the page underneath is unchanged.
	press(app, tabapp.KeyPress("esc"))
	if app.Model().OverlayCount != 0 {
		t.Fatalf("after Esc: overlays=%d, want 0", app.Model().OverlayCount)
	}
	if got := app.NavPath(); !slices.Equal(got, before) {
		t.Errorf("NavPath = %v, want %v unchanged", got, before)
	}

	// s reopens; s while open closes through the popup's own key.
	press(app, tabapp.KeyPress("s"))
	if app.Model().OverlayCount != 1 {
		t.Fatalf("after second s: overlays=%d, want 1", app.Model().OverlayCount)
	}
	press(app, tabapp.KeyPress("s"))
	if app.Model().OverlayCount != 0 || getState(app.Model()).SettingsID != 0 {
		t.Errorf("after third s: overlays=%d id=%d, want closed", app.Model().OverlayCount, getState(app.Model()).SettingsID)
	}
}

func TestSettingsResetRestoresDefaults(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	_, cmd := app.Update(ToggleSettingsMsg{})
	drain(app, cmd)

	// Walk every row and step it once.
	for row := settingsRow(0); row < rowCount; row++ {
		_, cmd = app.Update(SettingsAdjustMsg{Delta: 1})
		drain(app, cmd)
		_, cmd = app.Update(SettingsMoveMsg{Delta: 1})
		drain(app, cmd)
	}
	if snapshotOf(app.Model()) == defaultSettings {
		t.Fatal("precondition: settings should differ from defaults")
	}

	_, cmd = app.Update(SettingsResetMsg{})
	drain(app, cmd)

	if got := snapshotOf(app.Model()); got != defaultSettings {
		t.Errorf("after reset: %+v\nwant %+v", got, defaultSettings)
	}
	if app.NavBarHeight() != compactNavHeight {
		t.Errorf("nav bar height = %d, want %d", app.NavBarHeight(), compactNavHeight)
	}
}
