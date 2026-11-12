// Package tabapp runs a multi-page terminal application: it owns the layout
// (header, nav bar, content, footer), the navigation tree, the key
// dispatch and the overlay stack, and renders the page that is active.
//
// # Building an app
//
// Pages come from package page; sections group them into a tree of any
// depth. TabApp implements tea.Model and Run starts the Bubble Tea program:
//
//	app := tabapp.New()
//	app.MustAddPage(page.New("Home").SetContent(renderHome))
//	app.MustAddSection(tabapp.NewSection("Catalog").AddPages(phones, laptops))
//	if err := app.Run(); err != nil {
//	    log.Fatal(err)
//	}
//
// Every setter returns the TabApp, so configuration chains:
//
//	app.SetNavBarMode(tabapp.NavBarPath).SetNavBarAutoHide(false).UseTheme(styles.Nord)
//
// Render returns the frame as text, exactly the terminal's size. View wraps
// it in the tea.View that Bubble Tea draws.
//
// # State and messages
//
// The framework follows the Elm architecture. Settings that views read live
// in Model, which Model returns and every page content function receives.
// Keys are bound to messages, not callbacks: BindKey binds a key on every
// page, page.BindKey on one page, and Update applies the message in
// updateModel. The messages in messages.go cover the nav bar, the header
// and footer, theme, navigation and overlays; an application defines its own messages and
// handles them in the function given to SetCustomMessageHandler, keeping
// its state in Model.Custom.
//
// A key is named as a key press reports it: "s", "?", "space", "enter",
// "shift+tab", "ctrl+c". Only presses dispatch. A paste, a mouse event and
// a key release reach the application's handler as the Bubble Tea messages
// they are. KeyPress builds the message for a key name, which drives an
// application from a test:
//
//	app.Update(tabapp.KeyPress("s"))
//
// # The terminal
//
// What the application wants from the terminal is state in the Model, and
// View declares it on every frame: SetWindowTitle, SetTerminalColors (the
// terminal's own background and foreground follow the theme),
// SetReportFocus, SetMouseMode and SetProgress, each with a message. What
// the terminal reports is in the Model as well: Focused, DarkBackground,
// ColorProfile and KeyDisambiguation, which tells whether chords such as
// shift+enter can be bound.
//
// UseAdaptiveTheme switches between a light and a dark theme as the
// terminal's background says. SetProgramOptions passes options to the
// Bubble Tea program, which runs an application over SSH or without a
// terminal, and SetInitCmds adds commands that run at the start.
//
// # Navigation
//
// The tab-navigation keys (Tab and Shift+Tab by default) start navigation
// and move the highlight along the bar; the arrows, Enter and Backspace then
// move through the tree in the shape of the bar's layout (NavBarMode), and
// Esc or a timeout ends navigation. Digits 1-9 select a page at the current
// level. NavigateToPath jumps anywhere in the tree programmatically.
//
// # Header and footer
//
// SetHeader gives the app a header band drawn by a function that receives a
// flexbox and the Model; the footer lists the keys of the active page. Each
// band has a height in lines (SetHeaderHeight, SetFooterHeight; 0 hides it)
// and a Frame (SetHeaderFrame, SetFooterFrame): bare, inside a single-line,
// double-line or theme-shaped border, or filled as a solid bar in the
// theme's fill colors with the text reversed.
//
// # Overlays
//
// ShowOverlay, ShowModal and ShowNotification stack overlays from package
// overlays by z-index; keys go to the topmost overlay first and a blocking
// overlay stops them there. Notifications dismiss themselves after their
// timeout through the returned command.
package tabapp
