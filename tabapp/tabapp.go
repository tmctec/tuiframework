package tabapp

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/styles"
)

// NavBarPosition defines where the nav bar is rendered.
type NavBarPosition int

const (
	// NavBarTop renders tabs horizontally at the top (default).
	NavBarTop NavBarPosition = iota
	// NavBarLeft renders tabs vertically on the left side.
	NavBarLeft
)

// NavBarMode controls how the nav bar lays out the navigation tree.
// Autohide applies in every mode; see SetNavBarAutoHide.
type NavBarMode int

const (
	// NavBarCompact shows one strip: crumbs for the levels above the focus,
	// then the items at the focused level (default).
	NavBarCompact NavBarMode = iota
	// NavBarExpanded shows a row per level, so every level is reachable at once.
	NavBarExpanded
	// NavBarPath shows the active route on one line. While navigating, the
	// focused segment drops a list of its siblings beneath it.
	NavBarPath
	// NavBarColumns shows one column per level, Finder style: each level's
	// items listed vertically, with the selection cascading to the right.
	NavBarColumns
)

// String returns the mode's name.
func (m NavBarMode) String() string {
	switch m {
	case NavBarExpanded:
		return "Expanded"
	case NavBarPath:
		return "Path"
	case NavBarColumns:
		return "Columns"
	default:
		return "Compact"
	}
}

// next returns the mode after m, wrapping from Columns back to Compact.
func (m NavBarMode) next() NavBarMode {
	return (m + 1) % (NavBarColumns + 1)
}

// TabApp coordinates switching between pages using the Page framework.
//
// Settings that views read live in the Model and nowhere else; the fields
// below hold the content tree, the rendering caches derived from the theme,
// and transient interaction state.
type TabApp struct {
	model *Model // Elm architecture: the single source of truth for settings

	// Content tree. navPath holds one index per level, counting a level's
	// sections first and then its pages.
	root    *Section
	navPath []int

	// breadcrumbs is the trail EnterSection pushes and PopToBreadcrumb pops.
	breadcrumbs []Breadcrumb

	width  int
	height int
	canvas *overlays.Canvas // the frame View draws; its cells are reused between frames

	// Header: an optional band the app renders itself, above the content.
	headerFunc        HeaderFunc
	headerBorderColor color.Color

	// Navigation (tab) bar styles, set by UseTheme or the Set*Style setters.
	tabStyle              lipgloss.Style
	activeTabStyle        lipgloss.Style
	selectedInactiveStyle lipgloss.Style // selected but not being navigated
	navBarBorderStyle     lipgloss.Border
	navBarBorderColor     color.Color
	navBarWidth           int // side bar width (default 20)
	navBarHidden          bool
	navBarAutoHideID      int64 // id of the pending autohide timer

	// Content area.
	contentBorderStyle lipgloss.Border
	contentBorderColor color.Color
	contentBgColor     color.Color
	borderBlend        []color.Color // gradient for the content and nav bar borders; empty for none

	// Footer.
	footerBorderStyle lipgloss.Border
	footerBorderColor color.Color
	footerStyle       lipgloss.Style
	footerKeyColor    color.Color
	footerValueColor  color.Color
	footerSeparator   string

	// Key glow: the pressed key briefly highlights in the footer.
	keyGlowKey      string // key currently glowing, "" for none
	keyGlowTimerID  int64
	keyGlowDuration time.Duration // default 300ms
	keyGlowColor    color.Color

	// Key registries.
	universalBindings map[string]KeyBinding
	hardcodedKeys     map[string]internalKeyHandler // ctrl+c and ctrl+q; cannot be overridden

	// Overlay stack, sorted by z-index with the topmost last.
	overlayStack []overlayEntry
	overlaySeq   int64 // last overlay id handed out

	// Tab navigation: Tab enters it, the arrows move through the bar, and a
	// timeout or Esc leaves it.
	tabNavLevel   int    // level the highlight is on
	navFocusCrumb int    // compact strip: crumb level in focus, -1 = the level's items
	tabNavTimerID int64  // id of the pending timeout
	tabNavKeyFwd  string // default "tab"
	tabNavKeyRev  string // default "shift+tab"

	// customMsgHandler receives the messages tabapp does not handle itself.
	customMsgHandler func(msg tea.Msg) (bool, tea.Cmd)

	// The themes UseAdaptiveTheme switches between, by terminal background.
	adaptiveLight, adaptiveDark *styles.Theme

	programOptions []tea.ProgramOption // given to Bubble Tea by Run
	initCmds       []tea.Cmd           // the application's own startup commands
}

// New creates a TabApp with the Dracula theme, autohide on, and the default
// universal keys (q quits; Tab and Shift+Tab start tab navigation). Ctrl+C
// and Ctrl+Q always quit. The terminal size arrives through Bubble Tea.
func New() *TabApp {
	app := &TabApp{
		model:             NewModel(),
		canvas:            overlays.NewCanvas(),
		navPath:           []int{0},
		breadcrumbs:       make([]Breadcrumb, 0),
		navFocusCrumb:     -1,
		footerSeparator:   " • ",
		keyGlowDuration:   300 * time.Millisecond,
		universalBindings: make(map[string]KeyBinding),
		hardcodedKeys:     make(map[string]internalKeyHandler),
		navBarWidth:       20,
		tabNavKeyFwd:      "tab",
		tabNavKeyRev:      "shift+tab",
	}

	// The theme sets every color and style, including the footer and glow colors.
	app.UseTheme(styles.Dracula)

	app.hardcodedKeys["ctrl+c"] = func() tea.Cmd { return tea.Quit }
	app.hardcodedKeys["ctrl+q"] = func() tea.Cmd { return tea.Quit }

	// Universal keys dispatch messages; the tab-navigation keys are matched
	// by name in Update before the bindings are consulted.
	app.universalBindings["q"] = KeyBinding{key: "q", label: "Quit", msg: QuitMsg{}}
	app.universalBindings["tab"] = KeyBinding{key: "tab", msg: NoOpMsg{}}
	app.universalBindings["shift+tab"] = KeyBinding{key: "shift+tab", msg: NoOpMsg{}}

	return app
}

// Model returns the centralized application model (Elm architecture).
// Use this in View functions to access application state.
func (a *TabApp) Model() *Model {
	return a.model
}

// Run starts the application and blocks until it exits. It runs a Bubble
// Tea program with the options from SetProgramOptions.
func (a *TabApp) Run() error {
	p := tea.NewProgram(a, a.programOptions...)
	_, err := p.Run()
	return err
}

// Init initializes the application (Bubble Tea interface). It arms the
// autohide timer, asks the terminal for its background color, and runs the
// commands from SetInitCmds.
func (a *TabApp) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0, 2+len(a.initCmds))
	cmds = append(cmds, a.startAutohideTimer(), tea.RequestBackgroundColor)
	return tea.Batch(append(cmds, a.initCmds...)...)
}

// startAutohideTimer shows the nav bar and arms the autohide timer.
// Returns nil when autohide is disabled.
func (a *TabApp) startAutohideTimer() tea.Cmd {
	if !a.model.NavBarAutoHide {
		return nil
	}
	a.navBarHidden = false
	a.navBarAutoHideID++
	id := a.navBarAutoHideID
	return tea.Tick(a.model.NavBarTimeout, func(_ time.Time) tea.Msg {
		return navBarAutoHideMsg{id: id}
	})
}

// Update handles messages and updates the application state (Bubble Tea interface).
//
// A key press is resolved in this order: the hardcoded quit keys, the
// overlay stack (top first), the tab-navigation keys, the arrows while tab
// navigation is active, the universal bindings, the active page's bindings,
// and finally the digits 1-9, which select a page at the current level.
// Only presses dispatch; a key release, which a terminal reports on request,
// is not a key press. Every other message is an Elm message: built-in ones
// are applied by updateModel, and unknown ones reach the custom message
// handler from there. A paste, a mouse event and a key release are among
// the unknown ones.
func (a *TabApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		return a, nil

	case tea.KeyPressMsg:
		return a, a.handleKey(msg.String())

	case overlayTimeoutMsg:
		a.HideOverlayByID(msg.overlayID)
		return a, nil

	case navBarAutoHideMsg:
		// A stale timer (one re-armed since) has an older id and is ignored.
		if a.model.NavBarAutoHide && msg.id == a.navBarAutoHideID {
			a.navBarHidden = true
		}
		return a, nil

	case tabNavTimeoutMsg:
		if a.model.TabNavActive && msg.id == a.tabNavTimerID {
			a.model.TabNavActive = false
			a.navFocusCrumb = -1
		}
		return a, nil

	case keyGlowTimeoutMsg:
		if msg.id == a.keyGlowTimerID {
			a.keyGlowKey = ""
		}
		return a, nil
	}

	return a.updateModel(msg)
}

// handleKey resolves one key press; see Update for the order.
func (a *TabApp) handleKey(key string) tea.Cmd {
	if handler, ok := a.hardcodedKeys[key]; ok {
		return handler()
	}

	if a.HasOverlay() {
		if handled, cmd := a.handleOverlayKeys(key); handled {
			return cmd
		}
	}

	// The tab-navigation keys enter navigation and move the highlight. They
	// show the bar and re-arm the autohide timer.
	if key == a.tabNavKeyFwd || key == a.tabNavKeyRev {
		a.enterTabNav()
		delta := 1
		if key == a.tabNavKeyRev {
			delta = -1
		}
		if a.navShape() == navShapeStrip {
			a.navStripMove(delta)
		} else {
			a.moveAtLevel(delta)
		}
		return a.navTimers()
	}

	if a.model.TabNavActive {
		// Esc leaves navigation and is then released to the handlers below.
		if key == "esc" {
			a.leaveTabNav()
		}
		if handled, done := a.handleTabNavKey(key); handled {
			if done {
				a.leaveTabNav()
				return nil
			}
			return a.navTimers()
		}
	}

	if binding, ok := a.universalBindings[key]; ok {
		_, cmd := a.updateModel(binding.msg)
		return tea.Batch(a.startKeyGlow(key), cmd)
	}

	// Ctrl+arrows switch sections at the focused level in the expanded bar.
	if a.HasSections() && a.model.NavBarMode == NavBarExpanded {
		switch key {
		case "ctrl+right", "ctrl+l":
			_, cmd := a.updateModel(NavigateNextSectionMsg{})
			return cmd
		case "ctrl+left", "ctrl+h":
			_, cmd := a.updateModel(NavigatePrevSectionMsg{})
			return cmd
		}
	}

	if p := a.currentPage(); p != nil {
		if binding, ok := p.Bindings()[key]; ok {
			_, cmd := a.updateModel(binding.Msg())
			return tea.Batch(a.startKeyGlow(key), cmd)
		}
	}

	if len(key) == 1 && key >= "1" && key <= "9" {
		a.selectPageAtCurrentLevel(int(key[0] - '1'))
	}
	return nil
}
