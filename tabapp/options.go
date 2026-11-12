package tabapp

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/internal/keyname"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/styles"
)

// HeaderFunc renders the header band. Like a page content function it
// receives a fresh flexbox sized to the header and the Model, as
// interface{}, on every render; assert the model to *Model.
type HeaderFunc func(fb *flexbox.FlexBox, model interface{})

// HeaderHeight returns the header height in lines of content.
func (a *TabApp) HeaderHeight() int {
	return a.model.HeaderHeight
}

// NavBarHeight returns the nav bar height in lines of content.
func (a *TabApp) NavBarHeight() int {
	return a.model.NavBarHeight
}

// NavBarBelowHeader reports whether a top nav bar sits under the header.
func (a *TabApp) NavBarBelowHeader() bool {
	return a.model.NavBarBelowHeader
}

// FooterHeight returns the footer height in lines of content.
func (a *TabApp) FooterHeight() int {
	return a.model.FooterHeight
}

// NavBarAutoHide returns whether nav bar autohide is enabled.
func (a *TabApp) NavBarAutoHide() bool {
	return a.model.NavBarAutoHide
}

// NavBarAutoHideTimeout returns the nav bar autohide timeout duration.
func (a *TabApp) NavBarAutoHideTimeout() time.Duration {
	return a.model.NavBarTimeout
}

// NavBarPosition returns the current nav bar position.
func (a *TabApp) NavBarPosition() NavBarPosition {
	return a.model.NavBarPosition
}

// NavBarWidth returns the width of the side nav bar.
func (a *TabApp) NavBarWidth() int {
	return a.navBarWidth
}

// SetHeader sets the function that renders the header band above the content.
// A nil function means no header. Returns self for chaining.
func (a *TabApp) SetHeader(fn HeaderFunc) *TabApp {
	a.headerFunc = fn
	return a
}

// SetHeaderHeight sets the header height in lines of content; a border adds
// a row above and below. A height of 0 hides the header, and negative heights
// count as 0. Returns self for chaining.
func (a *TabApp) SetHeaderHeight(h int) *TabApp {
	a.model.HeaderHeight = max(h, 0)
	return a
}

// SetNavBarHeight sets the nav bar height in lines of content.
// Negative heights count as 0. Returns self for chaining.
func (a *TabApp) SetNavBarHeight(h int) *TabApp {
	a.model.NavBarHeight = max(h, 0)
	return a
}

// SetNavBarBelowHeader puts a top nav bar under the header (true) or
// above it (false). A side nav bar is unaffected.
// Returns self for chaining.
func (a *TabApp) SetNavBarBelowHeader(below bool) *TabApp {
	a.model.NavBarBelowHeader = below
	return a
}

// SetFooterHeight sets the footer height in lines of content; a border adds
// a row above and below. A height of 0 hides the footer, and negative heights
// count as 0. Returns self for chaining.
func (a *TabApp) SetFooterHeight(h int) *TabApp {
	a.model.FooterHeight = max(h, 0)
	return a
}

// SetHeaderFrame sets how the header band is set off from the content: bare,
// inside a border, or as a filled bar. Returns self for chaining.
//
// A header function draws a filled header on Theme.HeaderFill; cells it
// leaves without a background take the fill color.
func (a *TabApp) SetHeaderFrame(f Frame) *TabApp {
	a.model.HeaderFrame = f
	return a
}

// HeaderFrame returns how the header band is framed.
func (a *TabApp) HeaderFrame() Frame {
	return a.model.HeaderFrame
}

// SetFooterFrame sets how the footer is set off from the content: bare,
// inside a border, or as a filled bar with its keys and labels reversed.
// Returns self for chaining.
func (a *TabApp) SetFooterFrame(f Frame) *TabApp {
	a.model.FooterFrame = f
	return a
}

// FooterFrame returns how the footer is framed.
func (a *TabApp) FooterFrame() Frame {
	return a.model.FooterFrame
}

// SetShowHeaderBorder draws the header band inside the theme's border (true)
// or bare (false): shorthand for SetHeaderFrame with FrameTheme or FrameNone.
// Returns self for chaining.
func (a *TabApp) SetShowHeaderBorder(show bool) *TabApp {
	return a.SetHeaderFrame(themeFrame(show))
}

// ShowHeaderBorder reports whether the header frame is a border.
func (a *TabApp) ShowHeaderBorder() bool {
	return a.model.HeaderFrame.bordered()
}

// SetShowNavBarBorder sets whether the nav bar is drawn inside its border,
// across the top or down the left side. The border takes a row above and
// below the bar and a column on each side. Returns self for chaining.
func (a *TabApp) SetShowNavBarBorder(show bool) *TabApp {
	a.model.ShowNavBarBorder = show
	return a
}

// ShowNavBarBorder returns whether the nav bar border is shown.
func (a *TabApp) ShowNavBarBorder() bool {
	return a.model.ShowNavBarBorder
}

// SetShowFooterBorder draws the footer inside the theme's border (true) or
// bare (false): shorthand for SetFooterFrame with FrameTheme or FrameNone.
// Returns self for chaining.
func (a *TabApp) SetShowFooterBorder(show bool) *TabApp {
	return a.SetFooterFrame(themeFrame(show))
}

// ShowFooterBorder reports whether the footer frame is a border.
func (a *TabApp) ShowFooterBorder() bool {
	return a.model.FooterFrame.bordered()
}

// SetShowContentBorder sets whether to show the content border and returns self for chaining.
func (a *TabApp) SetShowContentBorder(show bool) *TabApp {
	a.model.ShowContentBorder = show
	return a
}

// ShowContentBorder returns whether the content border is shown.
func (a *TabApp) ShowContentBorder() bool {
	return a.model.ShowContentBorder
}

// SetNavBarBorderStyle sets the tab border style and returns self for chaining.
func (a *TabApp) SetNavBarBorderStyle(style lipgloss.Border) *TabApp {
	a.navBarBorderStyle = style
	return a
}

// SetNavBarBorderColor sets the tab border color and returns self for chaining.
func (a *TabApp) SetNavBarBorderColor(c color.Color) *TabApp {
	a.navBarBorderColor = c
	return a
}

// SetContentBorderStyle sets the content border style and returns self for chaining.
func (a *TabApp) SetContentBorderStyle(style lipgloss.Border) *TabApp {
	a.contentBorderStyle = style
	return a
}

// SetContentBorderColor sets the content border color and returns self for chaining.
func (a *TabApp) SetContentBorderColor(c color.Color) *TabApp {
	a.contentBorderColor = c
	return a
}

// SetBorderBlend draws the content and nav bar borders as a gradient through
// two or more colors, in place of their single colors. Fewer than two colors
// remove the gradient. A theme's Palette.BorderBlend sets the same thing.
// Returns self for chaining.
func (a *TabApp) SetBorderBlend(colors ...color.Color) *TabApp {
	if len(colors) < 2 {
		colors = nil
	}
	a.borderBlend = colors
	return a
}

// BorderBlend returns the colors of the border gradient, or nil without one.
func (a *TabApp) BorderBlend() []color.Color {
	return a.borderBlend
}

// SetContentBackground sets the content area background color and returns self for chaining.
// This fills the entire content area with the specified color, useful for light themes.
func (a *TabApp) SetContentBackground(c color.Color) *TabApp {
	a.contentBgColor = c
	return a
}

// ContentBackground returns the content area background color.
func (a *TabApp) ContentBackground() color.Color {
	return a.contentBgColor
}

// SetFooterBorderStyle sets the border shape FrameTheme draws around the
// footer and returns self for chaining.
func (a *TabApp) SetFooterBorderStyle(style lipgloss.Border) *TabApp {
	a.footerBorderStyle = style
	return a
}

// SetFooterBorderColor sets the footer border color and returns self for chaining.
func (a *TabApp) SetFooterBorderColor(c color.Color) *TabApp {
	a.footerBorderColor = c
	return a
}

// SetFooterStyle sets the footer style and returns self for chaining.
func (a *TabApp) SetFooterStyle(style lipgloss.Style) *TabApp {
	a.footerStyle = style
	return a
}

// SetWidth sets the width and returns self for chaining.
func (a *TabApp) SetWidth(w int) *TabApp {
	a.width = w
	return a
}

// SetHeight sets the height and returns self for chaining.
func (a *TabApp) SetHeight(h int) *TabApp {
	a.height = h
	return a
}

// SetTabStyle sets the inactive tab style and returns self for chaining.
func (a *TabApp) SetTabStyle(style lipgloss.Style) *TabApp {
	a.tabStyle = style
	return a
}

// SetActiveTabStyle sets the active tab style and returns self for chaining.
func (a *TabApp) SetActiveTabStyle(style lipgloss.Style) *TabApp {
	a.activeTabStyle = style
	return a
}

// SetFooterKeyColor sets the color for footer keys when using Page objects.
// Returns self for chaining.
func (a *TabApp) SetFooterKeyColor(c color.Color) *TabApp {
	a.footerKeyColor = c
	return a
}

// SetFooterValueColor sets the color for footer values when using Page objects.
// Returns self for chaining.
func (a *TabApp) SetFooterValueColor(c color.Color) *TabApp {
	a.footerValueColor = c
	return a
}

// TabStyle returns the style for inactive/stationary tabs.
func (a *TabApp) TabStyle() lipgloss.Style {
	return a.tabStyle
}

// ActiveTabStyle returns the style for the currently active/navigating tab.
func (a *TabApp) ActiveTabStyle() lipgloss.Style {
	return a.activeTabStyle
}

// SelectedInactiveStyle returns the style for selected but not navigating tabs.
func (a *TabApp) SelectedInactiveStyle() lipgloss.Style {
	return a.selectedInactiveStyle
}

// FooterKeyColor returns the color used for footer key labels.
func (a *TabApp) FooterKeyColor() color.Color {
	return a.footerKeyColor
}

// FooterValueColor returns the color used for footer value text.
func (a *TabApp) FooterValueColor() color.Color {
	return a.footerValueColor
}

// SetFooterSeparator sets the separator between footer items when using Page objects.
// Returns self for chaining.
func (a *TabApp) SetFooterSeparator(sep string) *TabApp {
	a.footerSeparator = sep
	return a
}

// SetKeyGlowEnabled enables or disables the key glow effect in the footer.
// When enabled, pressed keys briefly highlight in the footer.
// Returns self for chaining.
func (a *TabApp) SetKeyGlowEnabled(enabled bool) *TabApp {
	a.model.KeyGlowEnabled = enabled
	if !enabled {
		a.keyGlowKey = ""
	}
	return a
}

// KeyGlowEnabled returns whether the key glow effect is enabled.
func (a *TabApp) KeyGlowEnabled() bool {
	return a.model.KeyGlowEnabled
}

// SetKeyGlowDuration sets how long the key glow effect lasts.
// Default is 300ms. Returns self for chaining.
func (a *TabApp) SetKeyGlowDuration(d time.Duration) *TabApp {
	a.keyGlowDuration = d
	return a
}

// KeyGlowDuration returns the key glow effect duration.
func (a *TabApp) KeyGlowDuration() time.Duration {
	return a.keyGlowDuration
}

// SetKeyGlowColor sets the color for the key glow effect.
// Default is green (#50FA7B). Returns self for chaining.
func (a *TabApp) SetKeyGlowColor(c color.Color) *TabApp {
	a.keyGlowColor = c
	return a
}

// KeyGlowColor returns the key glow effect color.
func (a *TabApp) KeyGlowColor() color.Color {
	return a.keyGlowColor
}

// UseTheme applies a theme to every part of the interface: tabs, borders,
// footer, band fills, the modals on screen and the defaults for new
// notifications. Pages read it from Model.Theme on the next render. Use a
// built-in theme such as styles.Dracula, or styles.NewTheme with a palette
// of your own. Returns self for chaining.
func (a *TabApp) UseTheme(theme *styles.Theme) *TabApp {
	// A theme chosen by name stays: the terminal's background no longer
	// switches it.
	a.adaptiveLight, a.adaptiveDark = nil, nil
	a.model.FollowsBackground = false
	return a.applyTheme(theme)
}

// applyTheme makes theme the active theme and takes every style from it.
func (a *TabApp) applyTheme(theme *styles.Theme) *TabApp {
	// Sync to model for View access
	a.model.Theme = theme

	// Modals shown through TabApp whose frame comes from a theme follow it
	for _, e := range a.overlayStack {
		if e.modal != nil && e.modal.FollowsTheme() {
			e.modal.ApplyTheme(theme)
		}
	}

	// Tab styles
	a.tabStyle = theme.TabInactive
	a.activeTabStyle = theme.TabActive
	a.selectedInactiveStyle = theme.TabSelectedMuted
	a.navBarBorderStyle = theme.NavBarBorder
	a.navBarBorderColor = theme.NavBarBorderColor

	// Header
	a.headerBorderColor = theme.HeaderBorderColor

	// Content
	a.contentBorderStyle = theme.ContentBorder
	a.contentBorderColor = theme.ContentBorderColor
	a.borderBlend = theme.BorderBlend
	a.contentBgColor = theme.Palette.Bg0 // Fill content area with theme background

	// Footer
	a.footerBorderStyle = theme.FooterBorder
	a.footerBorderColor = theme.FooterBorderColor
	a.footerStyle = theme.Footer
	a.footerKeyColor = theme.Palette.FooterKeyColor
	a.footerValueColor = theme.Palette.FooterValueColor
	a.keyGlowColor = theme.KeyGlowColor()

	// New notifications take the theme's colors, icons and border.
	overlays.NotificationConfig.ApplyTheme(theme)

	return a
}

// SetNavBarAutoHide enables or disables autohide: the nav bar hides
// after the autohide timeout without navigation, and the tab-navigation
// keys bring it back. Returns self for chaining.
func (a *TabApp) SetNavBarAutoHide(enabled bool) *TabApp {
	a.model.NavBarAutoHide = enabled
	if !enabled {
		a.navBarHidden = false // without autohide the bar is always visible
	}
	return a
}

// SetNavBarAutoHideTimeout sets the duration before the nav bar auto-hides.
// Default is 2 seconds. Returns self for chaining.
func (a *TabApp) SetNavBarAutoHideTimeout(d time.Duration) *TabApp {
	a.model.NavBarTimeout = d
	return a
}

// SetNavBarPosition puts the nav bar across the top or down the left
// side. Returns self for chaining.
func (a *TabApp) SetNavBarPosition(pos NavBarPosition) *TabApp {
	a.model.NavBarPosition = pos
	a.navFocusCrumb = -1
	// For side-expanded mode, start at current page level
	if pos == NavBarLeft && a.model.NavBarMode == NavBarExpanded {
		a.tabNavLevel = a.currentPageLevel()
	}
	return a
}

// minNavBarWidth is the narrowest side bar: room for the border, the tab
// padding and four columns of text.
const minNavBarWidth = 10

// SetNavBarWidth sets the width of the side nav bar, used in the Left
// position. Widths below 10 are raised to 10 so the bar fits its border, the
// tab padding and four columns of text. Returns self for chaining.
func (a *TabApp) SetNavBarWidth(w int) *TabApp {
	a.navBarWidth = max(w, minNavBarWidth)
	return a
}

// NavBarMode returns the current nav bar mode.
func (a *TabApp) NavBarMode() NavBarMode {
	return a.model.NavBarMode
}

// SetNavBarMode sets the layout of the nav bar: NavBarCompact shows
// one strip of crumbs and the focused level, NavBarExpanded a row per level,
// NavBarPath the active route with a drop-down, and NavBarColumns a column
// per level. The navigation focus returns to the active page's level.
// Returns self for chaining.
func (a *TabApp) SetNavBarMode(mode NavBarMode) *TabApp {
	a.model.NavBarMode = mode
	a.navFocusCrumb = -1
	if mode == NavBarExpanded {
		a.navBarHidden = false
	}
	a.tabNavLevel = a.currentPageLevel()
	return a
}

// SetTabNavKeyForward sets the key to move forward in tab navigation (default: "tab").
// Returns self for chaining.
func (a *TabApp) SetTabNavKeyForward(key string) *TabApp {
	a.tabNavKeyFwd = keyname.Normalize(key)
	return a
}

// SetTabNavKeyReverse sets the key to move backward in tab navigation (default: "shift+tab").
// Returns self for chaining.
func (a *TabApp) SetTabNavKeyReverse(key string) *TabApp {
	a.tabNavKeyRev = keyname.Normalize(key)
	return a
}

// SetTabNavKeyTimeout sets the timeout before exiting tab navigation mode (default: 2s).
// When tab navigation is active, arrow keys navigate tabs. After this timeout
// with no key press, arrow keys return to page content handling.
// Returns self for chaining.
func (a *TabApp) SetTabNavKeyTimeout(d time.Duration) *TabApp {
	a.model.TabNavTimeout = d
	return a
}

// TabNavKeyTimeout returns the current tab navigation timeout duration.
func (a *TabApp) TabNavKeyTimeout() time.Duration {
	return a.model.TabNavTimeout
}

// TabNavActive returns whether tab navigation mode is currently active.
func (a *TabApp) TabNavActive() bool {
	return a.model.TabNavActive
}

// SetCustomMessageHandler sets a handler for application-specific messages.
// The handler is called for messages that tabapp doesn't handle internally.
// Return (true, cmd) to indicate the message was handled, or (false, nil) to pass through.
// This enables applications to implement custom message types (like tickers, events, etc.)
// while using tabapp as the main orchestrator.
// Returns self for chaining.
func (a *TabApp) SetCustomMessageHandler(handler func(msg tea.Msg) (bool, tea.Cmd)) *TabApp {
	a.customMsgHandler = handler
	return a
}
