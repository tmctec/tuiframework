package main

import (
	"fmt"
	"image/color"
	"slices"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// demoThemes is every built-in theme plus the demo's Cyberpunk theme.
// Model.Theme is the source of truth; the index is derived from it by pointer.
var demoThemes = append(slices.Clone(styles.AllThemes), CyberpunkTheme)

// fakeTabNames label the tab preview on the Themes page.
var fakeTabNames = []string{"Home", "Products", "Orders", "Reports", "Settings"}

// themeIndex returns the position of the active theme in demoThemes.
func themeIndex(m *tabapp.Model) int {
	if i := slices.Index(demoThemes, m.Theme); i >= 0 {
		return i
	}
	return 0
}

// cyberpunkIndex is the position of the Cyberpunk theme in demoThemes.
func cyberpunkIndex() int {
	return len(demoThemes) - 1
}

// setTheme emits the library message that activates demoThemes[i], wrapping i.
func setTheme(i int) tea.Cmd {
	n := len(demoThemes)
	i = ((i % n) + n) % n
	return emit(tabapp.SetThemeMsg{Theme: demoThemes[i]})
}

const themesHelp = `Themes

Left/Right  previous / next theme
a           move focus in the tab preview
n           cycle the 8 notification levels
?           this help

The whole app re-renders from Model.Theme,
so tabs, footer and overlays follow along.
Settings (s) offers the same list.`

// buildThemesPage shows the active palette: base colors with light and dark
// variants, background and foreground levels, hero colors, and the tab styles.
func buildThemesPage() *page.Page {
	p := page.New("Themes").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			palette := theme.Palette
			st := getState(m)
			themed(fb, theme)

			title := fixedRow(fb, centered(theme.AccentStyle).Bold(true), 1,
				fmt.Sprintf("Theme %d/%d: %s", themeIndex(m)+1, len(demoThemes), theme.Name()))

			fgBg := fixedRow(fb, centered(theme.MutedStyle), 1,
				fmt.Sprintf("Fg0 %s   Bg0 %s", styles.Hex(palette.Fg0), styles.Hex(palette.Bg0)))

			base := fixedRow(fb, centered(theme.BaseStyle), 4, renderBaseColors(palette))

			bgSwatch := func(bg, fg color.Color, label string) string {
				return lipgloss.NewStyle().
					Background(bg).
					Foreground(fg).
					Padding(0, 1).
					MarginRight(1).
					Border(lipgloss.NormalBorder()).
					BorderForeground(palette.Fg3).
					Render(fmt.Sprintf("%s %s", label, styles.Hex(bg)))
			}
			bgLevels := fixedRow(fb, centered(theme.BaseStyle), 3, lipgloss.JoinHorizontal(lipgloss.Center,
				bgSwatch(palette.Bg0, palette.Fg1, "Bg0"),
				bgSwatch(palette.Bg1, palette.Fg1, "Bg1"),
				bgSwatch(palette.Bg2, palette.Fg1, "Bg2"),
				bgSwatch(palette.Bg3, palette.Fg0, "Bg3"),
			))

			fgSwatch := func(fg color.Color, label string) string {
				return lipgloss.NewStyle().Foreground(fg).Padding(0, 1).MarginRight(1).
					Render(fmt.Sprintf("%s %s", label, styles.Hex(fg)))
			}
			fgLevels := fixedRow(fb, centered(theme.BaseStyle), 1, lipgloss.JoinHorizontal(lipgloss.Center,
				fgSwatch(palette.Fg0, "Fg0"),
				fgSwatch(palette.Fg1, "Fg1"),
				fgSwatch(palette.Fg2, "Fg2"),
				fgSwatch(palette.Fg3, "Fg3"),
			))

			heroSwatch := func(c color.Color, label string) string {
				return lipgloss.NewStyle().Background(c).Foreground(palette.Bg0).Padding(0, 2).MarginRight(1).Bold(true).
					Render(fmt.Sprintf("%s %s", label, styles.Hex(c)))
			}
			hero := fixedRow(fb, centered(theme.BaseStyle), 1, lipgloss.JoinHorizontal(lipgloss.Center,
				theme.MutedStyle.Render("Hero (tabs, borders): "),
				heroSwatch(palette.Hero, "Hero"),
				heroSwatch(palette.HeroDk, "HeroDk"),
			))

			legend := fixedRow(fb, centered(theme.BaseStyle), 1, lipgloss.JoinHorizontal(lipgloss.Top,
				theme.TabSelectedMuted.Render("Selected"),
				theme.TabActive.Render("Focused"),
				theme.TabInactive.Render("Unselected"),
			))

			preview := fixedRow(fb, centered(theme.BaseStyle), 2,
				renderFakeTabs(theme, st.FakeTabIndex)+"\n"+theme.MutedStyle.Render("a moves the focus"))

			fb.AddRows([]*flexbox.Row{
				title, fgBg, spacer(fb),
				base, spacer(fb),
				bgLevels, fgLevels, spacer(fb),
				hero, legend, spacer(fb),
				preview, filler(fb),
			})
		}).
		BindKey("left", "", StepThemeMsg{Delta: -1}).
		BindKey("right", "", StepThemeMsg{Delta: 1}).
		BindKey("a", "Tab focus", CycleFakeTabMsg{}).
		BindKey("n", "Notify", CycleNotificationMsg{}).
		SetFooter(map[string]string{
			"←/→": "Theme",
			"a":   "Tab focus",
			"n":   "Notify",
		})
	return withHelp(p, themesHelp)
}

// renderBaseColors draws the seven base hues in their normal, light and dark variants.
func renderBaseColors(palette styles.Palette) string {
	swatch := func(c color.Color) string {
		return lipgloss.NewStyle().Background(c).Foreground(palette.Bg0).Width(9).Align(lipgloss.Center).Render(styles.Hex(c)) + " "
	}
	label := func(c color.Color, name string) string {
		return lipgloss.NewStyle().Foreground(c).Bold(true).Width(9).Align(lipgloss.Center).Render(name) + " "
	}
	rowLabel := lipgloss.NewStyle().Foreground(palette.Fg2).Width(7).Align(lipgloss.Right).PaddingRight(1)

	type hue struct {
		name           string
		normal, lt, dk color.Color
	}
	hues := []hue{
		{"Red", palette.Red, palette.RedLt, palette.RedDk},
		{"Orange", palette.Orange, palette.OrangeLt, palette.OrangeDk},
		{"Yellow", palette.Yellow, palette.YellowLt, palette.YellowDk},
		{"Green", palette.Green, palette.GreenLt, palette.GreenDk},
		{"Blue", palette.Blue, palette.BlueLt, palette.BlueDk},
		{"Purple", palette.Purple, palette.PurpleLt, palette.PurpleDk},
		{"Aqua", palette.Aqua, palette.AquaLt, palette.AquaDk},
	}

	header, normal, light, dark := rowLabel.Render(""), rowLabel.Render("normal"), rowLabel.Render("light"), rowLabel.Render("dark")
	for _, h := range hues {
		header += label(h.normal, h.name)
		normal += swatch(h.normal)
		light += swatch(h.lt)
		dark += swatch(h.dk)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, normal, light, dark)
}

// renderFakeTabs draws a nav bar preview: the first tab is "selected", the
// focused one moves with a.
func renderFakeTabs(theme *styles.Theme, focused int) string {
	tabs := make([]string, 0, len(fakeTabNames))
	for i, name := range fakeTabNames {
		style := theme.TabInactive
		switch i {
		case focused:
			style = theme.TabActive
		case 0:
			style = theme.TabSelectedMuted
		}
		tabs = append(tabs, style.Render(name))
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}
