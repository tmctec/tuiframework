package main

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/tabapp"
)

const terminalHelp = `Terminal

p   run the terminal's progress bar from 0 to 100
c   copy the theme name to the clipboard

The first table is what the terminal reported about
itself. "not reported" means it did not answer.

The second is what the app declares to the terminal.
Terminal colors and following the background are
rows in the Settings popup (s).`

const (
	progressStep     = 5                      // percent a tick adds to the progress bar
	progressInterval = 120 * time.Millisecond // time between ticks
	notReported      = "not reported"
	terminalLabel    = 20 // width of the label column
)

// reportRow is one line of a table on the Terminal page.
type reportRow struct{ label, value string }

// terminalReported lists what the terminal told the app about itself. Every
// value is a word, so the table reads the same without color.
func terminalReported(m *tabapp.Model, st demoState) []reportRow {
	background := notReported + ", dark assumed"
	if m.BackgroundKnown {
		background = "light"
		if m.DarkBackground {
			background = "dark"
		}
	}
	profile := notReported
	if m.ColorProfile != colorprofile.Unknown {
		profile = m.ColorProfile.String()
	}
	chords := "no: bind a fallback beside shift+enter or ctrl+h"
	if m.KeyDisambiguation {
		chords = "yes: shift+enter and ctrl+h are keys of their own"
	}
	focus := "focused"
	if !m.Focused {
		focus = "not focused"
	}
	name := notReported
	if st.TerminalName != "" {
		name = st.TerminalName
	}
	return []reportRow{
		{"Background", background},
		{"Color profile", profile},
		{"Key chords", chords},
		{"Window", focus},
		{"Name and version", name},
	}
}

// terminalDeclared lists what the app declares to the terminal.
func terminalDeclared(m *tabapp.Model) []reportRow {
	progress := "idle"
	if m.ProgressState != tea.ProgressBarNone {
		progress = fmt.Sprintf("%d%%", m.ProgressValue)
	}
	title := "none"
	if m.WindowTitle != "" {
		title = m.WindowTitle
	}
	return []reportRow{
		{"Window title", title},
		{"Terminal colors", strings.ToLower(onOff(m.TerminalColors))},
		{"Follow terminal", strings.ToLower(onOff(m.FollowsBackground))},
		{"Progress bar", progress},
	}
}

// renderReport draws rows as two aligned columns: muted labels, plain values.
func renderReport(label, value lipgloss.Style, rows []reportRow) string {
	label = label.Width(terminalLabel)
	lines := make([]string, 0, len(rows))
	for _, r := range rows {
		lines = append(lines, label.Render(r.label)+value.Render(r.value))
	}
	return strings.Join(lines, "\n")
}

// buildTerminalPage shows what Bubble Tea learned from the terminal and
// what the app declares to it, with keys for the progress bar and the
// clipboard.
func buildTerminalPage() *page.Page {
	p := page.New("Terminal").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			themed(fb, theme)

			panel := card(theme).Align(lipgloss.Left)
			label := lipgloss.NewStyle().Foreground(theme.Palette.Fg2).Background(theme.Palette.Bg1)
			value := lipgloss.NewStyle().Foreground(theme.Palette.Fg0).Background(theme.Palette.Bg1)
			reported := terminalReported(m, getState(m))
			declared := terminalDeclared(m)

			fb.AddRows([]*flexbox.Row{
				fixedRow(fb, centered(theme.TextHeading), 1, "Terminal"),
				fixedRow(fb, centered(theme.MutedStyle), 1, "What this terminal reported, and what the app declares to it"),
				spacer(fb),
				fixedRow(fb, theme.TextHeading.PaddingLeft(2), 1, "Reported by the terminal"),
				fixedRow(fb, panel, len(reported)+2, renderReport(label, value, reported)),
				spacer(fb),
				fixedRow(fb, theme.TextHeading.PaddingLeft(2), 1, "Declared by the app"),
				fixedRow(fb, panel, len(declared)+2, renderReport(label, value, declared)),
				spacer(fb),
				fixedRow(fb, centered(theme.MutedStyle), 1, "p runs the progress bar · c copies the theme name · s opens Settings"),
				filler(fb),
			})
		}).
		BindKey("p", "Progress", RunProgressMsg{}).
		BindKey("c", "Copy theme", CopyThemeNameMsg{})
	return withHelp(p, terminalHelp)
}

// progressTick schedules the next step of the progress bar.
func progressTick() tea.Cmd {
	return tea.Tick(progressInterval, func(time.Time) tea.Msg { return progressTickMsg{} })
}

// stepProgress returns the message for the bar's next state: progressStep
// further, or cleared once it is full. done reports the clearing.
func stepProgress(m *tabapp.Model) (msg tabapp.SetProgressMsg, done bool) {
	next := m.ProgressValue + progressStep
	if m.ProgressState == tea.ProgressBarNone || next > 100 {
		return tabapp.SetProgressMsg{State: tea.ProgressBarNone}, true
	}
	return tabapp.SetProgressMsg{State: tea.ProgressBarDefault, Value: next}, false
}
