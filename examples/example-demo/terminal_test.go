package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/charmbracelet/x/ansi"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// terminalPagePath leads to the Terminal page: Showcase, its fifth page.
var terminalPagePath = []int{0, 4, 0, 0}

// onTerminalPage returns the demo sized and showing the Terminal page.
func onTerminalPage(t *testing.T) *tabapp.TabApp {
	t.Helper()
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if err := app.NavigateToPath(terminalPagePath); err != nil {
		t.Fatal(err)
	}
	return app
}

func plainFrame(app *tabapp.TabApp) string {
	return ansiSequence.ReplaceAllString(app.Render(), "")
}

// wantLine fails unless one line of the frame holds every part, in order.
func wantLine(t *testing.T, frame string, parts ...string) {
	t.Helper()
	for _, line := range strings.Split(frame, "\n") {
		at := 0
		for _, part := range parts {
			i := strings.Index(line[at:], part)
			if i < 0 {
				at = -1
				break
			}
			at += i + len(part)
		}
		if at >= 0 {
			return
		}
	}
	t.Errorf("no line holds %q in order:\n%s", parts, frame)
}

// The Terminal page states what the terminal reported in words, and says so
// when it reported nothing.
func TestTerminalPageReportsTheTerminal(t *testing.T) {
	app := onTerminalPage(t)

	frame := plainFrame(app)
	wantLine(t, frame, "Background", "not reported, dark assumed")
	wantLine(t, frame, "Color profile", "not reported")
	wantLine(t, frame, "Key chords", "no: bind a fallback")
	wantLine(t, frame, "Window", "focused")
	wantLine(t, frame, "Name and version", "not reported")
	wantLine(t, frame, "Window title", windowTitle)
	wantLine(t, frame, "Terminal colors", "on")
	wantLine(t, frame, "Follow terminal", "off")
	wantLine(t, frame, "Progress bar", "idle")

	app.Update(tea.BackgroundColorMsg{Color: lipgloss.Color("#fafafa")})
	app.Update(tea.ColorProfileMsg{Profile: colorprofile.TrueColor})
	app.Update(tea.KeyboardEnhancementsMsg{Flags: ansi.KittyDisambiguateEscapeCodes})
	app.Update(tea.BlurMsg{})
	app.Update(tea.TerminalVersionMsg{Name: "ghostty 1.2.0"})

	frame = plainFrame(app)
	wantLine(t, frame, "Background", "light")
	wantLine(t, frame, "Color profile", "TrueColor")
	wantLine(t, frame, "Key chords", "yes: shift+enter and ctrl+h")
	wantLine(t, frame, "Window", "not focused")
	wantLine(t, frame, "Name and version", "ghostty 1.2.0")

	app.Update(tea.BackgroundColorMsg{Color: lipgloss.Color("#101010")})
	wantLine(t, plainFrame(app), "Background", "dark")
}

// The demo declares its window title, terminal colors and focus reports,
// and asks for the terminal's name at the start.
func TestDemoDeclaresTerminalFeatures(t *testing.T) {
	app := buildApp()
	v := app.View()
	if v.WindowTitle != windowTitle || !v.ReportFocus || v.BackgroundColor != demoThemes[0].Palette.Bg0 {
		t.Errorf("view title %q focus %v background %v, want the demo's declarations",
			v.WindowTitle, v.ReportFocus, v.BackgroundColor)
	}
	if app.Init() == nil {
		t.Error("Init should return the startup requests")
	}
}

// p runs the progress bar from 0 to 100 and clears it; a second press while
// it runs changes nothing.
func TestProgressRun(t *testing.T) {
	app := onTerminalPage(t)
	m := app.Model()

	press(app, tabapp.KeyPress("p"))
	if m.ProgressState != tea.ProgressBarDefault || m.ProgressValue != 0 {
		t.Fatalf("after p: progress %v at %d, want running at 0", m.ProgressState, m.ProgressValue)
	}
	wantLine(t, plainFrame(app), "Progress bar", "0%")

	// The ticks arrive on a timer; step them by hand.
	_, cmd := app.Update(progressTickMsg{})
	drain(app, cmd)
	if m.ProgressValue != progressStep {
		t.Errorf("after one tick: %d, want %d", m.ProgressValue, progressStep)
	}
	press(app, tabapp.KeyPress("p")) // running: ignored
	if m.ProgressValue != progressStep {
		t.Errorf("a press while running moved the bar to %d", m.ProgressValue)
	}

	for i := 0; i < 100/progressStep+1 && m.ProgressState != tea.ProgressBarNone; i++ {
		_, cmd = app.Update(progressTickMsg{})
		drain(app, cmd)
	}
	if m.ProgressState != tea.ProgressBarNone {
		t.Errorf("the run should end cleared, progress %v at %d", m.ProgressState, m.ProgressValue)
	}
	wantLine(t, plainFrame(app), "Progress bar", "idle")

	// A stray tick after the end changes nothing.
	_, cmd = app.Update(progressTickMsg{})
	drain(app, cmd)
	if m.ProgressState != tea.ProgressBarNone {
		t.Error("a tick with no run should leave the bar cleared")
	}
}

func TestStepProgress(t *testing.T) {
	m := buildApp().Model()
	if msg, done := stepProgress(m); !done || msg.State != tea.ProgressBarNone {
		t.Errorf("no run: %+v done %v, want cleared and done", msg, done)
	}
	m.ProgressState, m.ProgressValue = tea.ProgressBarDefault, 95
	if msg, done := stepProgress(m); done || msg.Value != 100 {
		t.Errorf("at 95: %+v done %v, want 100 and not done", msg, done)
	}
	m.ProgressValue = 100
	if msg, done := stepProgress(m); !done || msg.State != tea.ProgressBarNone {
		t.Errorf("at 100: %+v done %v, want cleared and done", msg, done)
	}
}

// c copies the theme name and says so.
func TestCopyThemeName(t *testing.T) {
	app := onTerminalPage(t)
	_, cmd := app.Update(tabapp.KeyPress("c"))
	if cmd == nil {
		t.Fatal("c should return the clipboard command")
	}
	drain(app, cmd)
	top := app.TopOverlay()
	if top == nil || !strings.Contains(top.Content(), "Copied "+demoThemes[0].Name()) {
		t.Errorf("want a notification naming the copied theme, top = %v", top)
	}
}

// The two terminal rows of the Settings popup drive the library state the
// snapshot predicts.
func TestTerminalSettingsRows(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	m := app.Model()
	apply := func(row settingsRow) {
		t.Helper()
		before := snapshotOf(m)
		want, msgs := adjustRow(before, row, 1)
		drain(app, emit(msgs...))
		if got := snapshotOf(m); got != want {
			t.Errorf("row %s: app reached %+v\nthe reducer predicted %+v", rowLabel(row), got, want)
		}
	}

	apply(rowTerminalColors)
	if m.TerminalColors || app.View().BackgroundColor != nil {
		t.Error("terminal colors should be off")
	}
	apply(rowTerminalColors)
	if !m.TerminalColors {
		t.Error("terminal colors should be on again")
	}

	// Following: Dracula has a light sibling; on a dark background it stays.
	apply(rowFollowBackground)
	if !m.FollowsBackground || m.Theme != styles.Dracula {
		t.Fatalf("following on a dark background: theme %s follows %v", m.Theme.Name(), m.FollowsBackground)
	}
	app.Update(tea.BackgroundColorMsg{Color: lipgloss.Color("#fafafa")})
	if m.Theme != styles.DraculaLight {
		t.Errorf("a light background should switch to DraculaLight, theme %s", m.Theme.Name())
	}

	// Off: the active theme stays and no longer follows.
	apply(rowFollowBackground)
	if m.FollowsBackground || m.Theme != styles.DraculaLight {
		t.Errorf("following off: theme %s follows %v, want DraculaLight kept", m.Theme.Name(), m.FollowsBackground)
	}

	// On a light background, turning it on from a dark theme picks the light one.
	drain(app, emit(tabapp.SetThemeMsg{Theme: styles.Gruvbox}))
	apply(rowFollowBackground)
	if m.Theme != styles.GruvboxLight {
		t.Errorf("following Gruvbox on a light background: theme %s, want GruvboxLight", m.Theme.Name())
	}

	// Picking a theme by hand ends the following.
	apply(rowTheme)
	if m.FollowsBackground {
		t.Error("the Theme row should end the following")
	}

	// A theme without a sibling follows with the first theme's family.
	light, dark := adaptivePair(styles.Nord)
	if light != styles.DraculaLight || dark != styles.Dracula {
		t.Errorf("pair for Nord = %s, %s, want the Dracula family", light.Name(), dark.Name())
	}
	if light, dark = adaptivePair(styles.CatLatte); light != styles.CatLatte || dark != styles.CatMocha {
		t.Errorf("pair for CatLatte = %s, %s, want CatLatte and CatMocha", light.Name(), dark.Name())
	}
}

func TestRowWindow(t *testing.T) {
	n := int(rowCount)
	tests := []struct{ cursor, height, first, count int }{
		{0, 10, 0, 10},
		{5, 10, 0, 10},
		{9, 10, 4, 10},
		{n - 1, 10, n - 10, 10},
		{3, n, 0, n},
		{3, n + 5, 0, n},
		{7, 0, 7, 1},
		{7, -3, 7, 1},
	}
	for _, tt := range tests {
		first, count := rowWindow(tt.cursor, tt.height)
		if first != tt.first || count != tt.count {
			t.Errorf("rowWindow(%d, %d) = %d, %d, want %d, %d", tt.cursor, tt.height, first, count, tt.first, tt.count)
		}
		if tt.cursor < first || tt.cursor >= first+count {
			t.Errorf("rowWindow(%d, %d) hides the cursor", tt.cursor, tt.height)
		}
	}
}

// On an 80x24 terminal the popup shows part of the list; the cursor row is
// in view for every row, and the hint tells where it is.
func TestSettingsPopupScrollsOnAShortTerminal(t *testing.T) {
	needsNarrowBorders(t)
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	press(app, tabapp.KeyPress("s"))

	for row := settingsRow(0); row < rowCount; row++ {
		frame := plainFrame(app)
		if !strings.Contains(frame, rowLabel(row)) {
			t.Errorf("row %d (%s) is the cursor and is not in view:\n%s", row, rowLabel(row), frame)
		}
		if !strings.Contains(frame, settingsHint(int(row))) {
			t.Errorf("row %d: the hint should read %q", row, settingsHint(int(row)))
		}
		press(app, tabapp.KeyPress("down"))
	}
}

// runHeadless runs the real demo without a terminal and types the keys,
// each given as the bytes a terminal sends. It returns what the program
// wrote.
func runHeadless(t *testing.T, app *tabapp.TabApp, width, height int, keys ...string) string {
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
			time.Sleep(80 * time.Millisecond)
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
	case <-time.After(15 * time.Second):
		t.Fatal("the demo did not quit")
	}
	_ = typist.Close()
	return out.String()
}

// The real demo, run without a terminal: Welcome, the Terminal page, then
// Settings, where the theme changes and both terminal rows are switched.
func TestDemoRunsHeadless(t *testing.T) {
	const up, right, enter = "\x1b[A", "\x1b[C", "\r"
	app := buildApp()

	out := runHeadless(t, app, 100, 30,
		"5",       // the Terminal page
		"s",       // Settings, cursor on Theme
		right,     // next theme: DraculaLight
		up, enter, // Follow terminal: on
		up, enter, // Terminal colors: off
		"s", // close Settings
		"q",
	)

	plain := ansiSequence.ReplaceAllString(out, "")
	for _, text := range []string{"Welcome", "Reported by the terminal", "Settings", "Follow terminal"} {
		if !strings.Contains(plain, text) {
			t.Errorf("the stream never showed %q", text)
		}
	}

	setBg := func(th *styles.Theme) string { return "\x1b]11;" + styles.Hex(th.Palette.Bg0) + "\a" }
	first, second := strings.Index(out, setBg(styles.Dracula)), strings.Index(out, setBg(styles.DraculaLight))
	if first < 0 || second < first {
		t.Errorf("the terminal background should be set to Dracula's, then DraculaLight's (at %d and %d)", first, second)
	}
	if !strings.Contains(out, "\x1b]2;"+windowTitle+"\a") {
		t.Error("the stream should set the window title")
	}

	m := app.Model()
	if m.TerminalColors {
		t.Error("Terminal colors should have been switched off")
	}
	if !m.FollowsBackground || m.Theme != styles.Dracula {
		t.Errorf("following on with no report: theme %s follows %v, want the dark sibling Dracula", m.Theme.Name(), m.FollowsBackground)
	}
	if got := app.NavPath(); got[0] != terminalPagePath[0] || got[1] != terminalPagePath[1] {
		t.Errorf("the demo ended on %v, want the Terminal page", got)
	}
	if !strings.HasSuffix(strings.TrimRight(out, "\r\n"), "\a") || !strings.Contains(out, "\x1b[?1049l") {
		t.Error("the stream should end by restoring the terminal")
	}
}
