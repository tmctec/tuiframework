package main

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// CyberpunkPalette is a custom palette built on the 4-tier cascade:
//
//	Tier 1  base colors (Bg/Fg levels, Red..Aqua with Lt/Dk variants)
//	Tier 2  Hero, HeroDk              default: Purple, PurpleDk
//	Tier 3  Primary, Secondary, ...   default: Hero and base colors
//	Tier 4  UI element colors         default: the semantic colors
//
// Only Tier 1 and Border are needed; this palette sets a few Tier 3 and
// Tier 4 values to show how overrides look.
var CyberpunkPalette = styles.Palette{
	Name: "Cyberpunk",

	Bg0: lipgloss.Color("#0a0a12"), Bg1: lipgloss.Color("#14141f"),
	Bg2: lipgloss.Color("#1e1e2a"), Bg3: lipgloss.Color("#2a2a3a"),
	Fg0: lipgloss.Color("#ffffff"), Fg1: lipgloss.Color("#e0e0e8"),
	Fg2: lipgloss.Color("#a0a0b0"), Fg3: lipgloss.Color("#606070"),

	Gray: lipgloss.Color("#606070"), GrayLt: lipgloss.Color("#808090"), GrayDk: lipgloss.Color("#404050"),
	Red: lipgloss.Color("#ff0044"), RedLt: lipgloss.Color("#ff4477"), RedDk: lipgloss.Color("#cc0033"),
	Orange: lipgloss.Color("#ff8800"), OrangeLt: lipgloss.Color("#ffaa44"), OrangeDk: lipgloss.Color("#cc6600"),
	Green: lipgloss.Color("#00ff00"), GreenLt: lipgloss.Color("#44ff44"), GreenDk: lipgloss.Color("#00cc00"),
	Yellow: lipgloss.Color("#ffee00"), YellowLt: lipgloss.Color("#ffff66"), YellowDk: lipgloss.Color("#ccbb00"),
	Blue: lipgloss.Color("#00aaff"), BlueLt: lipgloss.Color("#44ccff"), BlueDk: lipgloss.Color("#0088cc"),
	Purple: lipgloss.Color("#ff00ff"), PurpleLt: lipgloss.Color("#ff44ff"), PurpleDk: lipgloss.Color("#aa00aa"),
	Pink: lipgloss.Color("#ff2d95"), PinkLt: lipgloss.Color("#ff6ab5"), PinkDk: lipgloss.Color("#cc1f77"),
	Aqua: lipgloss.Color("#00ffff"), AquaLt: lipgloss.Color("#44ffff"), AquaDk: lipgloss.Color("#00cccc"),

	// Tier 3 override: warnings in orange instead of the Yellow default.
	Warning: lipgloss.Color("#ff8800"),

	// Tier 4 overrides: cyan content border, green key glow.
	ContentBorderColor: lipgloss.Color("#00ffff"),
	KeyPressHighlight:  lipgloss.Color("#00ff00"),

	// A border gradient from magenta to cyan. It replaces the single border
	// colors on the content, nav bar and modal borders.
	BorderBlend: [2]color.Color{lipgloss.Color("#ff00ff"), lipgloss.Color("#00ffff")},

	Border: lipgloss.DoubleBorder(),

	IconInfo: "[i]", IconNote: "[*]", IconTip: "[>]", IconImportant: "[!]",
	IconSuccess: "[OK]", IconWarning: "[!!]", IconCaution: "[XX]", IconError: "[ERR]",
}

// CyberpunkTheme is built once so demoThemes can hold it by pointer.
var CyberpunkTheme = styles.NewTheme(CyberpunkPalette)

const customThemeHelp = `Custom Theme

c   switch between Cyberpunk and the first theme

A theme is a styles.Palette passed to
styles.NewTheme. Set the base colors and a
Border; every other color cascades from them.
The swatches show Cyberpunk's resolved colors
whatever theme is active.`

const customThemeCode = `var MyPalette = styles.Palette{
    Name: "MyTheme", Border: lipgloss.DoubleBorder(),
    Bg0: lipgloss.Color("#0a0a12"), Bg1: lipgloss.Color("#14141f"), Bg2: lipgloss.Color("#1e1e2a"), Bg3: lipgloss.Color("#2a2a3a"),
    Fg0: lipgloss.Color("#ffffff"), Fg1: lipgloss.Color("#e0e0e8"), Fg2: lipgloss.Color("#a0a0b0"), Fg3: lipgloss.Color("#606070"),
    Red: lipgloss.Color("#ff0044"), Green: lipgloss.Color("#00ff00"), Blue: lipgloss.Color("#00aaff"), Purple: lipgloss.Color("#ff00ff"),
    Hero: lipgloss.Color("#ff00ff"),   // optional: tabs and borders follow it
}
app.UseTheme(styles.NewTheme(MyPalette))   // or tabapp.SetThemeMsg`

// buildCustomThemePage explains how a palette becomes a theme and shows the
// Cyberpunk palette's resolved colors.
func buildCustomThemePage() *page.Page {
	p := page.New("Custom Theme").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			themed(fb, theme)

			status := "c applies it"
			if themeIndex(m) == cyberpunkIndex() {
				status = "active, c returns to " + demoThemes[0].Name()
			}
			title := fixedRow(fb, centered(theme.TextHeading), 1, "Custom theme: Cyberpunk")
			state := fixedRow(fb, centered(theme.MutedStyle), 1, status)

			swatch := func(c color.Color, label string) string {
				return lipgloss.NewStyle().Background(c).Foreground(CyberpunkPalette.Bg0).Bold(true).Padding(0, 1).
					MarginRight(1).Render(label)
			}
			cp := CyberpunkTheme.Palette
			swatches := fixedRow(fb, centered(theme.BaseStyle), 2, lipgloss.JoinVertical(lipgloss.Center,
				lipgloss.JoinHorizontal(lipgloss.Center,
					swatch(cp.Primary, "Primary"), swatch(cp.PrimaryDark, "PrimaryDark"),
					swatch(cp.Secondary, "Secondary"), swatch(cp.Accent, "Accent")),
				lipgloss.JoinHorizontal(lipgloss.Center,
					swatch(cp.Success, "Success"), swatch(cp.Warning, "Warning"),
					swatch(cp.Error, "Error"), swatch(cp.Info, "Info")),
			))

			code := fixedRow(fb, card(theme).Align(lipgloss.Left), 10, customThemeCode)

			note := fixedRow(fb, centered(theme.MutedStyle), 1,
				fmt.Sprintf("Cyberpunk is theme %d of %d in Settings and on the Themes page", cyberpunkIndex()+1, len(demoThemes)))

			fb.AddRows([]*flexbox.Row{title, state, spacer(fb), swatches, spacer(fb), code, note, filler(fb)})
		}).
		BindKey("c", "Cyberpunk", ToggleCyberpunkMsg{}).
		SetFooter(map[string]string{
			"c": "Toggle Cyberpunk",
		})
	return withHelp(p, customThemeHelp)
}
