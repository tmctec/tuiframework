package styles_test

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/styles"
)

// A palette needs its base colors and a border; the hero and semantic
// colors cascade from them and the Theme's Palette carries the result.
func ExampleNewTheme() {
	c := lipgloss.Color
	p := styles.Palette{
		Name:   "Mine",
		Border: lipgloss.RoundedBorder(),
		Bg0:    c("#1a1b26"), Bg1: c("#24283b"), Bg2: c("#414868"), Bg3: c("#565f89"),
		Fg0: c("#c0caf5"), Fg1: c("#a9b1d6"), Fg2: c("#9aa5ce"), Fg3: c("#565f89"),
		Red: c("#f7768e"), Green: c("#9ece6a"), Yellow: c("#e0af68"), Blue: c("#7aa2f7"),
		Purple: c("#bb9af7"), Aqua: c("#7dcfff"),
	}
	t := styles.NewTheme(p)

	fmt.Println(t.Name())
	fmt.Println("hero:", styles.Hex(t.Palette.Hero), "primary:", styles.Hex(t.Palette.Primary))
	fmt.Println("success:", styles.Hex(t.Palette.Success), "info:", styles.Hex(t.Palette.Info))
	fmt.Println("palette var untouched:", p.Primary == nil)

	t2 := styles.NewTheme(styles.Palette{Name: "Blue hero", Border: lipgloss.NormalBorder(), Purple: c("#ff00ff"), Hero: c("#0000ff")})
	fmt.Println("hero wins:", styles.Hex(t2.Palette.Primary))
	// Output:
	// Mine
	// hero: #bb9af7 primary: #bb9af7
	// success: #9ece6a info: #7dcfff
	// palette var untouched: true
	// hero wins: #0000ff
}

// The built-in themes are listed in AllThemes and found by name.
func ExampleThemeByName() {
	fmt.Println(len(styles.AllThemes), len(styles.AllPalettes))
	fmt.Println(styles.ThemeByName("Nord").Name(), styles.ThemeByName("Nope") == nil)
	fmt.Println(styles.AllThemes[0].Name(), styles.AllPalettes[0].Name)
	// Output:
	// 22 22
	// Nord true
	// Dracula Dracula
}
