package main

// Demo messages. Pages bind keys to these and the handler in handler.go
// applies them. Library behavior is driven by tabapp messages the handler
// emits in response. An application defines its own messages the same way.

// ToggleSettingsMsg opens Settings, or returns to the page it was opened from.
type ToggleSettingsMsg struct{}

// LeaveSettingsMsg returns from Settings to the previous page (Esc on Settings).
type LeaveSettingsMsg struct{}

// SettingsMoveMsg moves the Settings cursor by Delta rows, wrapping.
type SettingsMoveMsg struct{ Delta int }

// SettingsAdjustMsg changes the selected Settings row by Delta steps.
type SettingsAdjustMsg struct{ Delta int }

// SettingsResetMsg returns every setting to its startup value.
type SettingsResetMsg struct{}

// StepThemeMsg selects the theme Delta positions away in the demo theme list.
type StepThemeMsg struct{ Delta int }

// ToggleCyberpunkMsg switches between the Cyberpunk theme and the first theme.
type ToggleCyberpunkMsg struct{}

// CycleFakeTabMsg moves the focus in the Themes page tab preview.
type CycleFakeTabMsg struct{}

// MoveCursorMsg moves the Catalog product cursor.
type MoveCursorMsg struct{ DX, DY int }

// ToggleDemoModalMsg opens or closes one of the Overlays page modals (A, B or C).
type ToggleDemoModalMsg struct{ ID string }

// StackDemoModalsMsg opens three stacked modals at once, or closes them.
type StackDemoModalsMsg struct{}

// FrontDemoModalMsg brings one of the Overlays page modals (A, B or C) to
// the front and says so in a notification.
type FrontDemoModalMsg struct{ ID string }

// ClearOverlaysMsg hides every overlay and says so in a notification.
type ClearOverlaysMsg struct{}

// CycleNotificationMsg shows a sample notification of the next level, from
// info through error and around again.
type CycleNotificationMsg struct{}

// RunProgressMsg runs the terminal's progress bar from 0 to 100 and clears
// it. A press while it runs changes nothing.
type RunProgressMsg struct{}

// progressTickMsg advances the running progress bar one step.
type progressTickMsg struct{}

// CopyThemeNameMsg copies the name of the active theme to the clipboard and
// says so in a notification.
type CopyThemeNameMsg struct{}
