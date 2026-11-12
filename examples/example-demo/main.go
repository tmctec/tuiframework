// example-demo combines the framework's demos into one application.
//
// Layout: a header band (app name, theme and nav bar status), the
// nav bar below or above it, the page content, and the footer with
// the page's keys.
//
// Root tabs:
//
//	Showcase   Welcome · Themes · Custom Theme · Overlays · Terminal
//	Catalog    Electronics (Phones, PCs, Deals) · Fashion · Home (Furniture, Kitchen, Garden)
//
// Settings is a centered popup (s): theme, nav bar position, order and
// mode, autohide and its timeout, tab-nav timeout, header and footer height
// and border (none, single line, double line or fill), nav bar and content
// borders, key glow, terminal colors, and following the terminal's
// light or dark background.
//
// Global keys:
//
//	s                open or close the Settings popup
//	?                help for the current page
//	Tab / Shift+Tab  tab navigation: arrows move while it is active, Esc leaves it
//	1-9              jump to a page at the current level
//	q, Ctrl+C        quit
//
// Every page reads its state from the tabapp.Model. The message handler in
// handler.go is the only place that changes demo state, and it changes
// library state by emitting tabapp messages.
package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/tabapp"
)

const (
	headerHeight     = 1  // app name and status on one line
	footerHeight     = 2  // room for the keys to wrap once
	compactNavHeight = 2  // tab row plus the key hint line
	sideBarWidth     = 22 // width of the nav bar in the Left position

	windowTitle = "tuiframework demo" // title of the terminal window
)

// buildApp assembles the demo.
func buildApp() *tabapp.TabApp {
	app := tabapp.New()
	app.SetHeader(renderHeader).
		SetHeaderHeight(headerHeight).
		SetHeaderFrame(tabapp.FrameNone).
		SetNavBarHeight(compactNavHeight).
		SetNavBarBelowHeader(true).
		SetShowNavBarBorder(false).
		SetFooterHeight(footerHeight).
		SetFooterFrame(tabapp.FrameNone).
		SetShowContentBorder(false).
		SetNavBarMode(tabapp.NavBarCompact).
		SetNavBarAutoHide(true).
		SetNavBarWidth(sideBarWidth).
		SetKeyGlowEnabled(true).
		SetWindowTitle(windowTitle).
		SetTerminalColors(true).
		SetReportFocus(true).
		SetInitCmds(tea.RequestTerminalVersion).
		UseTheme(demoThemes[0])

	// Universal keys go first: AddPage and AddSection reject page keys that
	// collide with them.
	app.MustBindKey("s", "Settings", ToggleSettingsMsg{})

	app.MustAddSection(buildShowcase()).
		MustAddSection(buildCatalog())

	setState(app.Model(), defaultState())
	app.SetCustomMessageHandler(newHandler(app))
	return app
}

// buildShowcase is the tour: one page per framework feature.
func buildShowcase() *tabapp.Section {
	return tabapp.NewSection("Showcase").
		AddPage(buildWelcomePage()).
		AddPage(buildThemesPage()).
		AddPage(buildCustomThemePage()).
		AddPage(buildOverlaysPage()).
		AddPage(buildTerminalPage())
}

func main() {
	if err := buildApp().Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
