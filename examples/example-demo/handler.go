package main

import (
	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/tabapp"
)

// newHandler returns the custom message handler: the single place that
// changes demo state. Library state changes by emitting tabapp messages,
// which Update applies on the next pass.
func newHandler(app *tabapp.TabApp) func(tea.Msg) (bool, tea.Cmd) {
	return func(msg tea.Msg) (bool, tea.Cmd) {
		m := app.Model()
		st := getState(m)

		switch msg := msg.(type) {
		case ToggleSettingsMsg:
			if closeSettings(app, &st) {
				setState(m, st)
				return true, nil
			}
			if id, err := app.ShowModal(buildSettingsModal(app)); err == nil {
				st.SettingsID = id
			}
			setState(m, st)
			return true, nil

		case LeaveSettingsMsg:
			closeSettings(app, &st)
			setState(m, st)
			return true, nil

		case SettingsMoveMsg:
			st.SettingsCursor = moveCursor(st.SettingsCursor, msg.Delta, int(rowCount))
			setState(m, st)
			return true, nil

		case SettingsAdjustMsg:
			before := snapshotOf(m)
			after, msgs := adjustRow(before, settingsRow(st.SettingsCursor), msg.Delta)
			applyChrome(app, before, after)
			return true, emit(msgs...)

		case SettingsResetMsg:
			before := snapshotOf(m)
			applyChrome(app, before, defaultSettings)
			return true, emit(resetMsgs(before, defaultSettings)...)

		case StepThemeMsg:
			return true, setTheme(themeIndex(m) + msg.Delta)

		case ToggleCyberpunkMsg:
			if themeIndex(m) == cyberpunkIndex() {
				return true, setTheme(0)
			}
			return true, setTheme(cyberpunkIndex())

		case CycleFakeTabMsg:
			st.FakeTabIndex = (st.FakeTabIndex + 1) % len(fakeTabNames)
			setState(m, st)
			return true, nil

		case MoveCursorMsg:
			st.CursorX = clampGrid(st.CursorX + msg.DX)
			st.CursorY = clampGrid(st.CursorY + msg.DY)
			setState(m, st)
			return true, nil

		case ToggleDemoModalMsg:
			spec, ok := demoModals[msg.ID]
			if !ok {
				return true, nil
			}
			return true, emit(spec.toggleMsg(m.Theme.Palette))

		case StackDemoModalsMsg:
			return true, emit(stackMsgs(m.Theme.Palette)...)

		case FrontDemoModalMsg:
			spec, ok := demoModals[msg.ID]
			if !ok {
				return true, nil
			}
			if _, open := m.OpenModals[spec.Title]; !open {
				return true, emit(tabapp.ShowNotificationMsg{Level: overlays.LevelNote, Message: spec.Title + " not open"})
			}
			return true, emit(
				tabapp.BringOverlayToFrontMsg{Name: spec.Title},
				tabapp.ShowNotificationMsg{Level: overlays.LevelTip, Message: spec.Title + " brought to front"},
			)

		// The notification follows the clear, so the clear cannot take it.
		case ClearOverlaysMsg:
			return true, emitInOrder(
				tabapp.ClearAllOverlaysMsg{},
				tabapp.ShowNotificationMsg{Message: "All overlays cleared"},
			)

		case CycleNotificationMsg:
			level := overlays.NotificationLevel(st.NotifyLevel)
			st.NotifyLevel = (st.NotifyLevel + 1) % len(sampleNotifications)
			setState(m, st)
			return true, emit(tabapp.ShowNotificationMsg{Level: level, Message: sampleNotifications[level]})

		// The answer to the tea.RequestTerminalVersion that buildApp asks at
		// the start. The library does not keep it, so the demo does.
		case tea.TerminalVersionMsg:
			st.TerminalName = msg.Name
			setState(m, st)
			return true, nil

		// The progress bar is library state: the demo only decides the next
		// value and when to stop.
		case RunProgressMsg:
			if m.ProgressState != tea.ProgressBarNone {
				return true, nil // already running
			}
			return true, tea.Batch(emit(tabapp.SetProgressMsg{State: tea.ProgressBarDefault}), progressTick())

		case progressTickMsg:
			next, done := stepProgress(m)
			if done {
				return true, emit(next)
			}
			return true, tea.Batch(emit(next), progressTick())

		// The clipboard is the terminal's (OSC 52), so this works over SSH.
		case CopyThemeNameMsg:
			name := m.Theme.Name()
			return true, tea.Batch(
				tea.SetClipboard(name),
				emit(tabapp.ShowNotificationMsg{Level: overlays.LevelSuccess, Message: "Copied " + name}),
			)
		}
		return false, nil
	}
}

// closeSettings hides the Settings popup if it is open and reports whether it
// was. Esc closes the popup inside the library, so a stale id is expected.
func closeSettings(app *tabapp.TabApp, st *demoState) bool {
	if st.SettingsID == 0 {
		return false
	}
	wasOpen := app.HideOverlayByID(st.SettingsID)
	st.SettingsID = 0
	return wasOpen
}

// applyChrome handles the one setting with no library message: the
// nav bar height that the nav bar mode needs.
func applyChrome(app *tabapp.TabApp, before, after settingsSnapshot) {
	if before.Mode != after.Mode {
		app.SetNavBarHeight(navHeightFor(after.Mode, app.Root().Depth()))
	}
}

// navHeightFor returns the nav bar height a nav bar mode needs:
// Expanded draws one row per level plus the hint line, Compact draws two lines.
func navHeightFor(mode tabapp.NavBarMode, depth int) int {
	if mode == tabapp.NavBarExpanded {
		return depth + 1
	}
	return compactNavHeight // Path and Columns size the band themselves
}
