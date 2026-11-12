package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/tabapp"
)

const welcomeTree = `Showcase    Welcome · Themes · Custom Theme · Overlays · Terminal
Catalog     Electronics (Phones · PCs · Deals) · Fashion · Home (Furniture · Kitchen · Garden)
s           Settings popup: theme, nav bar, timeouts, band heights, borders, key glow`

const welcomeKeys = `Tab, Shift+Tab   tab navigation: ◂▸ move, Enter opens, Backspace back
1-9              jump to a page at the current level
s                Settings popup (Esc or s closes it)
?                help for this page
q                quit`

const welcomeHelp = `Welcome

Tab / Shift+Tab   enter tab navigation
◂▸                move along the nav bar
Enter / Backspace open the highlight / go back a level
Esc               leave tab navigation
1-9               jump to a page at this level
s                 Settings popup (Esc or s closes)
?                 help for the current page
q, Ctrl+C         quit

Each page's help lists its own keys.`

// buildWelcomePage is the orientation screen and the start location.
func buildWelcomePage() *page.Page {
	p := page.New("Welcome").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			themed(fb, theme)

			title := fixedRow(fb, centered(theme.TextHeading), 1, "Welcome")
			sub := fixedRow(fb, centered(theme.MutedStyle), 1,
				fmt.Sprintf("%d themes to try, a settings popup on s, help on ? for every page", len(demoThemes)))
			tree := fixedRow(fb, card(theme).Align(lipgloss.Left), 5, welcomeTree)
			keys := fixedRow(fb, centered(theme.MutedStyle), 6, block(welcomeKeys))

			fb.AddRows([]*flexbox.Row{title, sub, spacer(fb), tree, spacer(fb), keys, filler(fb)})
		}).
		SetFooter(map[string]string{
			"Tab": "Navigate",
		})
	return withHelp(p, welcomeHelp)
}
