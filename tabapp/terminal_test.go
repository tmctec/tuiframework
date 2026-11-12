package tabapp

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// runHeadless runs the real program without a terminal: input from a pipe,
// output into a buffer, a window of width by height. It types each key
// after the first frame is out and returns everything the program wrote.
// Every script ends by quitting, or the test fails on the timeout.
func runHeadless(t *testing.T, app *TabApp, width, height int, keys ...string) string {
	t.Helper()
	in, typist := io.Pipe()
	var out bytes.Buffer
	app.SetProgramOptions(
		tea.WithInput(in), tea.WithOutput(&out),
		tea.WithWindowSize(width, height),
		tea.WithColorProfile(colorprofile.TrueColor),
		tea.WithEnvironment([]string{"TERM=xterm-256color"}),
	)

	go func() {
		for _, k := range keys {
			time.Sleep(60 * time.Millisecond)
			_, _ = typist.Write([]byte(k))
		}
	}()
	done := make(chan error, 1)
	go func() { done <- app.Run() }()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run returned %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("the program did not quit")
	}
	_ = typist.Close()
	return out.String()
}

// Terminal sequences the program writes, as Bubble Tea emits them.
const (
	oscSetForeground   = "\x1b]10;"
	oscSetBackground   = "\x1b]11;"
	oscResetForeground = "\x1b]110\a"
	oscResetBackground = "\x1b]111\a"
	oscTitle           = "\x1b]2;"
	oscProgress        = "\x1b]9;4;"
	oscProgressClear   = "\x1b]9;4;0\a"
)

func terminalApp() (*TabApp, *recorder) {
	rec := &recorder{}
	app := sizedApp(New().MustAddPage(page.New("One")).SetNavBarAutoHide(false))
	app.SetCustomMessageHandler(rec.handle)
	return app, rec
}

// tc-1 to tc-4: the terminal's own colors follow the theme while the setting
// is on, change with the theme, and go back to nil when it is off.
func TestTerminalColorsFollowTheTheme(t *testing.T) {
	app, _ := terminalApp()
	app.UseTheme(styles.Nord)

	if v := app.View(); v.BackgroundColor != nil || v.ForegroundColor != nil {
		t.Errorf("off by default: colors = %v %v, want nil", v.BackgroundColor, v.ForegroundColor)
	}

	app.SetTerminalColors(true)
	v := app.View()
	if v.BackgroundColor != styles.Nord.Palette.Bg0 || v.ForegroundColor != styles.Nord.Palette.Fg1 {
		t.Errorf("colors = %v %v, want Nord's Bg0 and Fg1", v.BackgroundColor, v.ForegroundColor)
	}
	if !app.Model().TerminalColors {
		t.Error("the Model should record the setting")
	}

	app.Update(SetThemeMsg{Theme: styles.CatLatte})
	v = app.View()
	if v.BackgroundColor != styles.CatLatte.Palette.Bg0 || v.ForegroundColor != styles.CatLatte.Palette.Fg1 {
		t.Errorf("after a theme change: colors = %v %v, want CatLatte's", v.BackgroundColor, v.ForegroundColor)
	}

	app.Update(ToggleTerminalColorsMsg{})
	if v := app.View(); v.BackgroundColor != nil || v.ForegroundColor != nil || app.Model().TerminalColors {
		t.Error("toggled off: the view should carry no colors")
	}
	app.Update(ToggleTerminalColorsMsg{})
	if !app.Model().TerminalColors {
		t.Error("toggled again: the setting should be on")
	}
}

// tc-5 and pg-4: in a real run the program sets the terminal's colors to the
// theme's and shows the progress bar, and when it quits it resets the colors
// and clears the bar.
func TestProgramSetsAndRestoresTheTerminal(t *testing.T) {
	app := New().MustAddPage(page.New("One")).UseTheme(styles.Nord).
		SetTerminalColors(true).SetWindowTitle("axinite").SetProgress(tea.ProgressBarDefault, 40)

	out := runHeadless(t, app, 60, 10, "q")

	bg, fg := styles.Hex(styles.Nord.Palette.Bg0), styles.Hex(styles.Nord.Palette.Fg1)
	for name, seq := range map[string]string{
		"set background":   oscSetBackground + bg + "\a",
		"set foreground":   oscSetForeground + fg + "\a",
		"window title":     oscTitle + "axinite\a",
		"progress at 40":   oscProgress + "1;40\a",
		"alternate screen": "\x1b[?1049h",
	} {
		if !strings.Contains(out, seq) {
			t.Errorf("the stream lacks the %s sequence %q", name, seq)
		}
	}

	// What it set, it puts back, after the last frame.
	frameEnd := strings.LastIndex(out, "Quit")
	for name, seq := range map[string]string{
		"background reset": oscResetBackground,
		"foreground reset": oscResetForeground,
		"progress cleared": oscProgressClear,
		"title cleared":    oscTitle + "\a",
		"main screen":      "\x1b[?1049l",
	} {
		if at := strings.LastIndex(out, seq); at < frameEnd {
			t.Errorf("the stream should end with the %s sequence %q", name, seq)
		}
	}
}

// A theme change during a real run sends the new background to the terminal.
func TestProgramFollowsAThemeChange(t *testing.T) {
	app := New().MustAddPage(page.New("One")).UseTheme(styles.Nord).SetTerminalColors(true).
		MustBindKey("t", "Theme", SetThemeMsg{Theme: styles.Gruvbox})

	out := runHeadless(t, app, 60, 10, "t", "q")

	nord := oscSetBackground + styles.Hex(styles.Nord.Palette.Bg0) + "\a"
	gruvbox := oscSetBackground + styles.Hex(styles.Gruvbox.Palette.Bg0) + "\a"
	first, second := strings.Index(out, nord), strings.Index(out, gruvbox)
	if first < 0 || second < first {
		t.Errorf("the stream should set Nord's background, then Gruvbox's (at %d and %d)", first, second)
	}
}

// wt-1 to wt-3: the window title comes from the Model.
func TestWindowTitle(t *testing.T) {
	app, _ := terminalApp()
	if got := app.View().WindowTitle; got != "" {
		t.Errorf("no title set: WindowTitle = %q, want none", got)
	}

	app.SetWindowTitle("Inventory")
	if got := app.View().WindowTitle; got != "Inventory" || app.Model().WindowTitle != "Inventory" {
		t.Errorf("WindowTitle = %q, want Inventory in the view and the Model", got)
	}

	app.Update(SetWindowTitleMsg{Title: "Inventory: Phones"})
	if got := app.View().WindowTitle; got != "Inventory: Phones" {
		t.Errorf("after the message: WindowTitle = %q", got)
	}
}

// pg-1 to pg-3: the terminal's progress bar comes from the Model, with the
// value kept between 0 and 100.
func TestProgressBar(t *testing.T) {
	app, _ := terminalApp()
	if app.View().ProgressBar != nil {
		t.Fatal("no progress set: the view should carry no bar")
	}

	app.Update(SetProgressMsg{State: tea.ProgressBarDefault, Value: 40})
	bar := app.View().ProgressBar
	if bar == nil || bar.State != tea.ProgressBarDefault || bar.Value != 40 {
		t.Fatalf("bar = %+v, want the default state at 40", bar)
	}

	for _, tc := range []struct{ value, want int }{{140, 100}, {100, 100}, {0, 0}, {-5, 0}, {73, 73}} {
		app.SetProgress(tea.ProgressBarError, tc.value)
		if got := app.View().ProgressBar; got.Value != tc.want || got.State != tea.ProgressBarError {
			t.Errorf("value %d: bar = %+v, want the error state at %d", tc.value, got, tc.want)
		}
	}
	if app.Model().ProgressValue != 73 || app.Model().ProgressState != tea.ProgressBarError {
		t.Errorf("Model = %v at %d, want the last state set", app.Model().ProgressState, app.Model().ProgressValue)
	}

	app.Update(SetProgressMsg{State: tea.ProgressBarNone})
	if app.View().ProgressBar != nil {
		t.Error("ProgressBarNone should clear the bar")
	}
}

// fo-1 to fo-3: with focus reporting on, the Model follows the terminal's
// focus reports and the application's handler receives them.
func TestFocusReporting(t *testing.T) {
	app, rec := terminalApp()
	if app.View().ReportFocus || !app.Model().Focused {
		t.Fatal("off by default, and a window is focused until told otherwise")
	}

	// Off: a stray blur changes nothing.
	app.Update(tea.BlurMsg{})
	if !app.Model().Focused {
		t.Error("without focus reporting the Model stays focused")
	}

	app.SetReportFocus(true)
	if !app.View().ReportFocus || !app.Model().ReportFocus {
		t.Fatal("the view should ask for focus reports")
	}
	app.Update(tea.BlurMsg{})
	if app.Model().Focused {
		t.Error("a blur should clear Model.Focused")
	}
	app.Update(tea.FocusMsg{})
	if !app.Model().Focused {
		t.Error("a focus should set Model.Focused")
	}
	if len(rec.other) != 3 {
		t.Errorf("the handler received %d reports, want all 3", len(rec.other))
	}

	// Turning it off while blurred restores the default.
	app.Update(tea.BlurMsg{})
	app.Update(ToggleReportFocusMsg{})
	if app.Model().ReportFocus || !app.Model().Focused || app.View().ReportFocus {
		t.Error("toggled off: no reports asked for, and the Model focused again")
	}
	app.Update(ToggleReportFocusMsg{})
	if !app.Model().ReportFocus {
		t.Error("toggled again: focus reporting should be on")
	}
}

// The handler sees a report after the Model has taken it.
func TestReportsReachTheHandlerAfterTheModel(t *testing.T) {
	app := sizedApp(New().MustAddPage(page.New("One"))).SetReportFocus(true)
	var seen []bool
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		if _, ok := msg.(tea.BlurMsg); ok {
			seen = append(seen, app.Model().Focused)
		}
		return false, func() tea.Msg { return nil }
	})

	_, cmd := app.Update(tea.BlurMsg{})
	if len(seen) != 1 || seen[0] {
		t.Errorf("the handler saw Focused = %v, want false already", seen)
	}
	if cmd == nil {
		t.Error("the handler's command should be returned")
	}

	// Without a handler a report is simply recorded.
	bare := sizedApp(New().MustAddPage(page.New("One")))
	if _, cmd := bare.Update(tea.ColorProfileMsg{Profile: colorprofile.ANSI256}); cmd != nil {
		t.Error("no handler: no command")
	}
	if bare.Model().ColorProfile != colorprofile.ANSI256 {
		t.Errorf("ColorProfile = %v, want ANSI256", bare.Model().ColorProfile)
	}
}

// mo-1 and mo-2: the application picks the mouse mode and receives the
// mouse messages; the library binds none.
func TestMousePassThrough(t *testing.T) {
	app, rec := terminalApp()
	if got := app.View().MouseMode; got != tea.MouseModeNone {
		t.Fatalf("default mouse mode = %v, want none", got)
	}

	app.SetMouseMode(tea.MouseModeCellMotion)
	if got := app.View().MouseMode; got != tea.MouseModeCellMotion || app.Model().MouseMode != got {
		t.Errorf("mouse mode = %v, want cell motion in the view and the Model", got)
	}

	app.Update(tea.MouseClickMsg(tea.Mouse{X: 3, Y: 1, Button: tea.MouseLeft}))
	if len(rec.other) != 1 {
		t.Fatalf("the handler received %d messages, want the click", len(rec.other))
	}
	if click, ok := rec.other[0].(tea.MouseClickMsg); !ok || click.Mouse().X != 3 {
		t.Errorf("the handler received %#v, want the click", rec.other[0])
	}

	app.Update(SetMouseModeMsg{Mode: tea.MouseModeAllMotion})
	if got := app.View().MouseMode; got != tea.MouseModeAllMotion {
		t.Errorf("after the message: mouse mode = %v, want all motion", got)
	}
}

// po-1: with program options the app runs without a terminal, at the size
// given, and the commands from SetInitCmds run at the start.
func TestRunWithProgramOptions(t *testing.T) {
	type startedMsg struct{}
	started := false
	app := New().MustAddPage(page.New("One"))
	app.SetInitCmds(func() tea.Msg { return startedMsg{} })
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		if _, ok := msg.(startedMsg); ok {
			started = true
		}
		return false, nil
	})

	out := runHeadless(t, app, 80, 24, "q")

	if app.width != 80 || app.height != 24 {
		t.Errorf("the app is %dx%d, want the 80x24 given", app.width, app.height)
	}
	frame := ansiSequence.ReplaceAllString(app.Render(), "")
	if lipgloss.Width(frame) != 80 || lipgloss.Height(frame) != 24 {
		t.Errorf("the frame is %dx%d, want 80x24", lipgloss.Width(frame), lipgloss.Height(frame))
	}
	if !started {
		t.Error("the command from SetInitCmds should have run")
	}
	if !strings.Contains(out, oscSetBackground+"?\a") {
		t.Error("the program should ask the terminal for its background color at the start")
	}
	if strings.Contains(out, oscSetBackground+"#") {
		t.Error("terminal colors are off: the program should not set the background")
	}
}

// po-2: without options and without a terminal, Run returns the error from
// Bubble Tea instead of hanging or panicking.
func TestRunWithoutATerminalReturnsAnError(t *testing.T) {
	if tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0); err == nil {
		_ = tty.Close()
		t.Skip("a terminal is attached: Run would take it over")
	}

	app := New().MustAddPage(page.New("One"))
	done := make(chan error, 1)
	go func() { done <- app.Run() }()
	select {
	case err := <-done:
		if err == nil {
			t.Error("Run without a terminal should return an error")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("Run without a terminal should return, not hang")
	}
}

// ad-1 to ad-6: an adaptive pair follows the terminal's background.
func TestAdaptiveThemeFollowsTheBackground(t *testing.T) {
	light, dark := styles.CatLatte, styles.CatMocha
	black, white := lipgloss.Color("#101010"), lipgloss.Color("#fafafa")
	report := func(app *TabApp, c interface {
		RGBA() (r, g, b, a uint32)
	}) {
		app.Update(tea.BackgroundColorMsg{Color: c})
	}

	// ad-4: until the terminal answers, the dark member is active.
	app, rec := terminalApp()
	app.UseAdaptiveTheme(light, dark)
	if app.Model().Theme != dark || !app.Model().FollowsBackground || app.Model().BackgroundKnown {
		t.Fatalf("before any report: theme %s, want the dark member", app.Model().Theme.Name())
	}

	// ad-2, ad-1 and ad-6: each report switches to the matching member.
	report(app, white)
	if app.Model().Theme != light || app.Model().DarkBackground || !app.Model().BackgroundKnown {
		t.Errorf("light background: theme %s, want the light member", app.Model().Theme.Name())
	}
	report(app, black)
	if app.Model().Theme != dark || !app.Model().DarkBackground {
		t.Errorf("dark background: theme %s, want the dark member", app.Model().Theme.Name())
	}
	report(app, white)
	if app.Model().Theme != light {
		t.Errorf("the background flipped again: theme %s, want the light member", app.Model().Theme.Name())
	}
	if len(rec.other) != 3 {
		t.Errorf("the handler received %d reports, want 3", len(rec.other))
	}
	if app.tabStyle.GetBackground() != light.TabInactive.GetBackground() {
		t.Error("the styles should come from the theme now active")
	}

	// ad-5: a theme the application picks stays, whatever is reported later.
	app.Update(SetThemeMsg{Theme: styles.Nord})
	report(app, black)
	if app.Model().Theme != styles.Nord || app.Model().FollowsBackground {
		t.Errorf("after an explicit choice: theme %s, want Nord kept", app.Model().Theme.Name())
	}
	if !app.Model().DarkBackground {
		t.Error("the background is still recorded")
	}

	// Following resumes on request, with the background already known.
	app.Update(UseAdaptiveThemeMsg{Light: light, Dark: dark})
	if app.Model().Theme != dark || !app.Model().FollowsBackground {
		t.Errorf("following resumed on a dark background: theme %s, want the dark member", app.Model().Theme.Name())
	}
	app.Update(UseAdaptiveThemeMsg{Light: light})
	if app.Model().Theme != dark {
		t.Error("a message without both themes changes nothing")
	}

	// A cycle is an explicit choice too.
	app.Update(CycleThemeMsg{Direction: 1, Themes: styles.AllThemes})
	if app.Model().FollowsBackground {
		t.Error("cycling themes ends the following")
	}

	// ad-3: without a pair a report changes no theme and is still recorded.
	plain, _ := terminalApp()
	plain.UseTheme(styles.Gruvbox)
	report(plain, white)
	if plain.Model().Theme != styles.Gruvbox || plain.Model().DarkBackground || !plain.Model().BackgroundKnown {
		t.Errorf("no pair: theme %s dark %v, want Gruvbox and a light background recorded",
			plain.Model().Theme.Name(), plain.Model().DarkBackground)
	}
}
