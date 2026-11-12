package tabapp

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/styles"
)

// This file holds what a TabApp declares to the terminal and what it learns
// from it. The Model carries both; View reads the Model and sets the fields
// of the tea.View, so every one of these is a setting, never a command.

// SetWindowTitle sets the title of the terminal window. An empty title
// leaves the terminal's own. Returns self for chaining.
func (a *TabApp) SetWindowTitle(title string) *TabApp {
	a.model.WindowTitle = title
	return a
}

// SetTerminalColors makes the terminal's own background and foreground
// follow the theme: the base background and the primary text color. The
// space a terminal keeps around the cell grid then matches the frame.
// Bubble Tea restores the terminal's colors when the program exits. Off by
// default. Returns self for chaining.
func (a *TabApp) SetTerminalColors(on bool) *TabApp {
	a.model.TerminalColors = on
	return a
}

// SetReportFocus asks the terminal to report when its window gains and
// loses focus. While on, Model.Focused follows the reports, and the
// application's message handler receives tea.FocusMsg and tea.BlurMsg.
// Returns self for chaining.
func (a *TabApp) SetReportFocus(on bool) *TabApp {
	a.model.ReportFocus = on
	if !on {
		a.model.Focused = true
	}
	return a
}

// SetMouseMode sets which mouse events the terminal reports. The library
// binds no mouse action; it hands tea.MouseClickMsg, tea.MouseReleaseMsg,
// tea.MouseWheelMsg and tea.MouseMotionMsg to the application's message
// handler. The default is tea.MouseModeNone, which leaves text selection to
// the terminal. Returns self for chaining.
func (a *TabApp) SetMouseMode(mode tea.MouseMode) *TabApp {
	a.model.MouseMode = mode
	return a
}

// SetProgress shows the terminal's own progress bar, where the terminal has
// one: in its tab, title bar or taskbar entry. The value is a percentage
// and is kept between 0 and 100. tea.ProgressBarNone clears the bar.
// Returns self for chaining.
func (a *TabApp) SetProgress(state tea.ProgressBarState, value int) *TabApp {
	a.model.ProgressState = state
	a.model.ProgressValue = min(max(value, 0), 100)
	return a
}

// SetProgramOptions sets the options Run gives to Bubble Tea, replacing any
// set before. Use them to run over SSH or without a terminal:
//
//	app.SetProgramOptions(tea.WithInput(in), tea.WithOutput(out), tea.WithWindowSize(80, 24))
//
// Returns self for chaining.
func (a *TabApp) SetProgramOptions(opts ...tea.ProgramOption) *TabApp {
	a.programOptions = opts
	return a
}

// SetInitCmds sets commands that run once when the program starts, beside
// the library's own, replacing any set before. Use them for a request such
// as tea.RequestTerminalVersion; the answer reaches the application's
// message handler. Returns self for chaining.
func (a *TabApp) SetInitCmds(cmds ...tea.Cmd) *TabApp {
	a.initCmds = cmds
	return a
}

// UseAdaptiveTheme makes the theme follow the terminal's background: light
// on a light background, dark on a dark one. The app asks the terminal for
// its background when it starts and switches on every answer. Until the
// terminal answers, and on terminals that never do, the dark theme is
// active. Choosing a theme with UseTheme, SetThemeMsg or CycleThemeMsg ends
// the following. Returns self for chaining.
func (a *TabApp) UseAdaptiveTheme(light, dark *styles.Theme) *TabApp {
	a.adaptiveLight, a.adaptiveDark = light, dark
	a.model.FollowsBackground = true
	a.followBackground()
	return a
}

// followBackground applies the member of the adaptive pair that suits the
// terminal's background. Without a pair it does nothing.
func (a *TabApp) followBackground() {
	if !a.model.FollowsBackground {
		return
	}
	theme := a.adaptiveLight
	if a.model.DarkBackground {
		theme = a.adaptiveDark
	}
	a.applyTheme(theme)
}

// forward hands a message the library has recorded to the application's
// message handler, so the application can react to it as well.
func (a *TabApp) forward(msg tea.Msg) tea.Cmd {
	if a.customMsgHandler == nil {
		return nil
	}
	_, cmd := a.customMsgHandler(msg)
	return cmd
}

// declare sets the fields of a view from the Model: the terminal features
// the application wants, stated on every frame.
func (a *TabApp) declare(v *tea.View) {
	m := a.model
	v.AltScreen = true
	v.WindowTitle = m.WindowTitle
	v.ReportFocus = m.ReportFocus
	v.MouseMode = m.MouseMode
	if m.TerminalColors {
		v.BackgroundColor = m.Theme.Palette.Bg0
		v.ForegroundColor = m.Theme.Palette.Fg1
	}
	if m.ProgressState != tea.ProgressBarNone {
		v.ProgressBar = tea.NewProgressBar(m.ProgressState, m.ProgressValue)
	}
}
