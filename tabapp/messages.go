package tabapp

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/styles"
)

// Msg is the interface for all messages in the Elm architecture.
// Messages describe what CAN happen, not how.
type Msg interface{}

// App-level messages

// QuitMsg requests the application to quit.
type QuitMsg struct{}

// NoOpMsg is a message that does nothing (for keys handled specially elsewhere).
type NoOpMsg struct{}

// ToggleNavBarPositionMsg toggles the nav bar between top and left positions.
// LeftWidth is optional - if > 0, sets the nav bar width when switching to left.
type ToggleNavBarPositionMsg struct {
	LeftWidth int // Optional: width to set when switching to left position
}

// SetNavBarWidthMsg sets the nav bar width (for left position).
type SetNavBarWidthMsg struct {
	Width int
}

// ToggleNavBarAutoHideMsg toggles the nav bar auto-hide feature.
type ToggleNavBarAutoHideMsg struct{}

// SetNavBarTimeoutMsg sets the auto-hide timeout duration.
type SetNavBarTimeoutMsg struct {
	Duration time.Duration
}

// AdjustNavBarTimeoutMsg adjusts the timeout by a delta (positive or negative).
type AdjustNavBarTimeoutMsg struct {
	Delta time.Duration
}

// SetTabNavTimeoutMsg sets the tab navigation mode timeout duration.
type SetTabNavTimeoutMsg struct {
	Duration time.Duration
}

// AdjustTabNavTimeoutMsg adjusts the tab nav timeout by a delta (positive or negative).
type AdjustTabNavTimeoutMsg struct {
	Delta time.Duration
}

// NavigateToPageMsg navigates to a specific page by index.
type NavigateToPageMsg struct {
	Index int
}

// NavigateNextPageMsg navigates to the next page.
type NavigateNextPageMsg struct{}

// NavigatePrevPageMsg navigates to the previous page.
type NavigatePrevPageMsg struct{}

// ShowNotificationMsg shows a notification of a level, with its message and
// the z-index it stacks at. The zero Level is overlays.LevelInfo.
type ShowNotificationMsg struct {
	Level   overlays.NotificationLevel
	Message string
	ZIndex  int
}

// ShowModalMsg shows a blocking modal in the active theme with a title and
// text content, closed by esc. ZIndex is where it stacks. SizePct is its
// size as a percentage of the screen; 0 keeps the overlay default.
type ShowModalMsg struct {
	Title   string
	Content string
	ZIndex  int
	SizePct int
}

// DismissOverlayMsg dismisses the current overlay.
type DismissOverlayMsg struct{}

// Section Navigation Messages

// EnterSectionMsg navigates to a section at a specific level and index.
// Pushes current location onto breadcrumb stack.
type EnterSectionMsg struct {
	Level int // Which level in the hierarchy
	Index int // Index of section at that level
}

// PopToBreadcrumbMsg pops back to a breadcrumb, discarding deeper positions.
type PopToBreadcrumbMsg struct {
	BreadcrumbIndex int
}

// NavigateNextSectionMsg navigates to the next section at the current level.
type NavigateNextSectionMsg struct{}

// NavigatePrevSectionMsg navigates to the previous section at the current level.
type NavigatePrevSectionMsg struct{}

// NavigateToPathMsg jumps to an absolute navigation path: one index per
// hierarchy level, counting sections first and then pages at each level.
// Indices below a selected page are ignored (forced to 0).
type NavigateToPathMsg struct {
	Path []int
}

// SetWindowTitleMsg sets the title of the terminal window. An empty title
// leaves the terminal's own.
type SetWindowTitleMsg struct {
	Title string
}

// ToggleTerminalColorsMsg turns on or off the terminal's own background and
// foreground following the theme.
type ToggleTerminalColorsMsg struct{}

// ToggleReportFocusMsg turns the terminal's focus reports on or off.
type ToggleReportFocusMsg struct{}

// SetMouseModeMsg sets which mouse events the terminal reports.
type SetMouseModeMsg struct {
	Mode tea.MouseMode
}

// SetProgressMsg shows the terminal's progress bar at Value percent, or
// clears it with tea.ProgressBarNone.
type SetProgressMsg struct {
	State tea.ProgressBarState
	Value int
}

// UseAdaptiveThemeMsg makes the theme follow the terminal's background:
// Light on a light background, Dark on a dark one.
type UseAdaptiveThemeMsg struct {
	Light *styles.Theme
	Dark  *styles.Theme
}

// SetThemeMsg switches the app theme.
// Use styles.Dracula, styles.Gruvbox, etc. or a custom theme from styles.NewTheme().
type SetThemeMsg struct {
	Theme *styles.Theme
}

// CycleThemeMsg cycles to next/previous theme from a provided list.
// Direction: 1 for next, -1 for previous.
type CycleThemeMsg struct {
	Direction int             // 1 for next, -1 for previous
	Themes    []*styles.Theme // List of available themes
}

// Nav bar mode and chrome messages

// SetNavBarModeMsg sets the nav bar mode: Compact, Expanded, Path or Columns.
type SetNavBarModeMsg struct {
	Mode NavBarMode
}

// ToggleNavBarModeMsg steps to the next nav bar mode: Compact, Expanded,
// Path, Columns and back to Compact.
type ToggleNavBarModeMsg struct{}

// SetHeaderHeightMsg sets the header height in lines of content; 0 hides the
// header.
type SetHeaderHeightMsg struct {
	Height int
}

// SetFooterHeightMsg sets the footer height in lines of content; 0 hides the
// footer.
type SetFooterHeightMsg struct {
	Height int
}

// SetHeaderFrameMsg sets how the header band is framed.
type SetHeaderFrameMsg struct {
	Frame Frame
}

// SetFooterFrameMsg sets how the footer is framed.
type SetFooterFrameMsg struct {
	Frame Frame
}

// ToggleHeaderBorderMsg switches the header band between the theme's border
// and no frame; a framed header of any kind becomes bare.
type ToggleHeaderBorderMsg struct{}

// ToggleNavBarBorderMsg toggles the nav bar's border.
type ToggleNavBarBorderMsg struct{}

// ToggleNavBarBelowHeaderMsg moves a top nav bar under or above the header.
type ToggleNavBarBelowHeaderMsg struct{}

// ToggleFooterBorderMsg switches the footer between the theme's border and no
// frame; a framed footer of any kind becomes bare.
type ToggleFooterBorderMsg struct{}

// ToggleContentBorderMsg toggles the content border visibility.
type ToggleContentBorderMsg struct{}

// ToggleKeyGlowMsg toggles the footer key glow effect.
type ToggleKeyGlowMsg struct{}

// Overlay stack messages

// ToggleModalMsg opens a themed modal that lets keys through to the page,
// or closes it when a modal with that title is open. Model.OpenModals lists
// the open ones by title.
type ToggleModalMsg struct {
	Title       string
	Content     string // Text content for the modal body
	ZIndex      int
	SizePct     int
	BorderColor color.Color // The border color; nil keeps the theme's
}

// BringOverlayToFrontMsg brings the ToggleModalMsg modal titled Name to the
// front. A title that is not open changes nothing.
type BringOverlayToFrontMsg struct {
	Name string
}

// ClearAllOverlaysMsg hides every overlay.
type ClearAllOverlaysMsg struct{}

// Internal messages (not exported for app use)

// tabNavTimeoutMsg is sent when the tab navigation mode timeout expires.
// This is an internal message - apps should not send it directly.
type tabNavTimeoutMsg struct {
	id int64 // Matches tabNavTimerID to prevent stale timeouts
}

// keyGlowTimeoutMsg is sent when the key glow effect should be cleared.
type keyGlowTimeoutMsg struct {
	id int64 // Matches glowTimerID to prevent stale clears
}

// overlayTimeoutMsg signals that an overlay's timeout has expired.
type overlayTimeoutMsg struct {
	overlayID int64 // Unique ID to match the overlay that started the timer
}

// navBarAutoHideMsg signals that the nav bar autohide timeout has expired.
type navBarAutoHideMsg struct {
	id int64 // Unique ID to match the timer that started the autohide
}
