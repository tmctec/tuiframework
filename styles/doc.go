// Package styles is the theme system: a Palette of colors becomes a Theme of
// ready-made Lip Gloss styles through NewTheme, and TabApp applies a Theme to
// every part of the interface with UseTheme.
//
// Colors cascade in four tiers, so a palette needs only its base colors and
// a border; everything else derives from them:
//
//  1. base colors: Bg0-Bg3, Fg0-Fg3, and Gray..Aqua with Lt and Dk variants
//  2. hero colors: Hero and HeroDk, the theme's signature accent
//  3. semantic colors: Primary, Secondary, Accent, Success, Warning, Error, Info
//  4. UI element colors: tab, border, footer, breadcrumb, modal and
//     notification colors
//
// Twenty-two themes are built in, from Dracula to Charm; AllThemes lists
// them and ThemeByName looks one up:
//
//	app.UseTheme(styles.Nord)
//
// A color is an image/color.Color, as lipgloss.Color returns it, and a color
// that is not set is nil. A palette of your own names its base colors:
//
//	c := lipgloss.Color
//	var MyPalette = styles.Palette{
//	    Name: "Mine", Border: lipgloss.RoundedBorder(),
//	    Bg0: c("#1a1b26"), Bg1: c("#24283b"), Bg2: c("#414868"), Bg3: c("#565f89"),
//	    Fg0: c("#c0caf5"), Fg1: c("#a9b1d6"), Fg2: c("#9aa5ce"), Fg3: c("#565f89"),
//	    Red: c("#f7768e"), Green: c("#9ece6a"), Blue: c("#7aa2f7"), // ... and the rest
//	    Hero: c("#7aa2f7"),
//	}
//	app.UseTheme(styles.NewTheme(MyPalette))
//
// PaletteFrom derives the whole palette from three colors instead: the
// background, the body text and one accent. It refuses colors that would
// not be readable together.
//
//	p, err := styles.PaletteFrom(c("#101418"), c("#d7dde4"), c("#7aa2f7"))
//
// Palette.BorderBlend draws the content, nav bar and modal borders as a
// gradient between two colors.
//
// Theme.Palette is the resolved palette, so a view can read
// theme.Palette.Primary or theme.Palette.FooterKeyColor without knowing
// whether the palette set it. Hex returns a color as "#rrggbb" and Contrast
// the WCAG contrast ratio of two colors.
//
// Theme.IsDark tells a dark theme from a light one, and Counterpart returns
// the built-in theme of the same family for the other kind of background.
// TabApp.UseAdaptiveTheme switches between such a pair as the terminal's
// background says.
package styles
