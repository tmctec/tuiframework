package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/tabapp"
)

// stateKey is the Model.Custom entry that holds the demo's own state.
// Library state (theme, nav bar, borders, overlay count) lives in the Model's
// own fields.
const stateKey = "demo"

// demoState is the state the demo owns.
type demoState struct {
	SettingsCursor   int    // selected Settings row
	SettingsID       int64  // overlay id of the open Settings popup; 0 when closed
	CursorX, CursorY int    // Catalog product grid cursor
	FakeTabIndex     int    // focused tab in the Themes page preview
	NotifyLevel      int    // level the next sample notification takes
	TerminalName     string // name and version the terminal reported; "" until it answers
}

func defaultState() demoState {
	return demoState{CursorX: 1, CursorY: 1}
}

func getState(m *tabapp.Model) demoState {
	if s, ok := m.GetCustom(stateKey).(demoState); ok {
		return s
	}
	return defaultState()
}

func setState(m *tabapp.Model, s demoState) {
	m.SetCustom(stateKey, s)
}

// emit returns a command that delivers every message to Update, in no
// particular order.
func emit(msgs ...tea.Msg) tea.Cmd {
	if len(msgs) == 0 {
		return nil
	}
	return tea.Batch(commands(msgs)...)
}

// emitInOrder returns a command that delivers the messages to Update one
// after the other, for messages whose effects depend on their order.
func emitInOrder(msgs ...tea.Msg) tea.Cmd {
	return tea.Sequence(commands(msgs)...)
}

// commands wraps each message in a command that returns it.
func commands(msgs []tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, 0, len(msgs))
	for _, msg := range msgs {
		cmds = append(cmds, func() tea.Msg { return msg })
	}
	return cmds
}
