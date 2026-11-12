package tabapp

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/colorprofile"
	"github.com/tmctec/tuiframework/styles"
)

// Model represents the centralized application state following the Elm architecture.
// All application state lives here, and pages render based on this model.
type Model struct {
	// Nav bar state
	NavBarPosition NavBarPosition
	NavBarAutoHide bool
	NavBarTimeout  time.Duration
	NavBarMode     NavBarMode // Compact, Expanded, Path or Columns

	// Tab navigation state
	TabNavTimeout time.Duration // Timeout before exiting tab nav mode
	TabNavActive  bool          // Whether tab navigation mode is currently active

	// Theme
	Theme *styles.Theme // Current active theme

	// Navigation state
	ActivePage int // Index of the active page among the pages at its level

	// Chrome
	NavBarBelowHeader bool  // A top nav bar sits under the header
	HeaderHeight      int   // Lines of header content; 0 hides the header
	NavBarHeight      int   // Lines of nav bar content
	FooterHeight      int   // Lines of footer content; 0 hides the footer
	HeaderFrame       Frame // How the header band is set off from the content
	FooterFrame       Frame // How the footer is set off from the content
	ShowNavBarBorder  bool
	ShowContentBorder bool
	KeyGlowEnabled    bool // Pressed keys briefly highlight in the footer

	// What the application declares to the terminal. View states these on
	// every frame; see the setters of the same names.
	WindowTitle    string               // Title of the terminal window; "" leaves the terminal's own
	TerminalColors bool                 // The terminal's own background and foreground follow the theme
	ReportFocus    bool                 // The terminal reports focus changes
	MouseMode      tea.MouseMode        // Which mouse events the terminal reports
	ProgressState  tea.ProgressBarState // The terminal's progress bar; ProgressBarNone hides it
	ProgressValue  int                  // Progress in percent, 0 to 100

	// What the terminal reported.
	Focused           bool                 // The terminal window has focus; true unless ReportFocus says otherwise
	DarkBackground    bool                 // The terminal's background is dark; true until it reports otherwise
	BackgroundKnown   bool                 // The terminal reported its background color
	FollowsBackground bool                 // The theme follows DarkBackground; see UseAdaptiveTheme
	ColorProfile      colorprofile.Profile // The colors the terminal can show; Unknown until reported
	KeyDisambiguation bool                 // The terminal tells chords such as shift+enter and ctrl+h apart

	// Overlays
	OverlayCount int              // Number of overlays currently shown
	OpenModals   map[string]int64 // Overlay ids of the modals ToggleModalMsg has open, by title

	// Custom holds the application's own state. The library never reads or
	// writes it; use the GetCustom accessors or type assertions to read it.
	Custom map[string]interface{}
}

// NewModel creates a new Model with default values.
func NewModel() *Model {
	return &Model{
		NavBarPosition:    NavBarTop,
		NavBarAutoHide:    true,
		NavBarTimeout:     2 * time.Second,
		NavBarMode:        NavBarCompact,
		TabNavTimeout:     2 * time.Second,
		Theme:             styles.Dracula,
		ActivePage:        0,
		NavBarBelowHeader: true,
		HeaderHeight:      1,
		NavBarHeight:      3,
		FooterHeight:      3,
		KeyGlowEnabled:    true,
		Focused:           true,
		DarkBackground:    true,
		OpenModals:        make(map[string]int64),
		Custom:            make(map[string]interface{}),
	}
}

// SetCustom sets a custom state value by key.
func (m *Model) SetCustom(key string, value interface{}) *Model {
	m.Custom[key] = value
	return m
}

// GetCustom retrieves a custom state value by key.
// Returns nil if not found.
func (m *Model) GetCustom(key string) interface{} {
	return m.Custom[key]
}

// GetCustomString retrieves a custom string value.
// Returns empty string if not found or not a string.
func (m *Model) GetCustomString(key string) string {
	if v, ok := m.Custom[key].(string); ok {
		return v
	}
	return ""
}

// GetCustomInt retrieves a custom int value.
// Returns 0 if not found or not an int.
func (m *Model) GetCustomInt(key string) int {
	if v, ok := m.Custom[key].(int); ok {
		return v
	}
	return 0
}

// GetCustomBool retrieves a custom bool value.
// Returns false if not found or not a bool.
func (m *Model) GetCustomBool(key string) bool {
	if v, ok := m.Custom[key].(bool); ok {
		return v
	}
	return false
}
