package styles

import (
	"reflect"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestThemeNameComesFromPalette(t *testing.T) {
	if got := NewTheme(basePalette()).Name(); got != "Base" {
		t.Errorf("Name() = %q, want %q", got, "Base")
	}
	for _, th := range AllThemes {
		if th.Name() != th.Palette.Name {
			t.Errorf("Name() = %q, want Palette.Name %q", th.Name(), th.Palette.Name)
		}
	}
}

func TestKeyGlowColorFallsBackToGreenLt(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	wantColor(t, "KeyGlowColor", th.KeyGlowColor(), p.GreenLt)
	wantColor(t, "FooterGlow foreground", th.FooterGlow.GetForeground(), th.KeyGlowColor())
}

func TestKeyGlowColorUsesKeyPressHighlight(t *testing.T) {
	p := basePalette()
	p.KeyPressHighlight = lipgloss.Color("#e00005")
	th := NewTheme(p)

	wantColor(t, "KeyGlowColor", th.KeyGlowColor(), p.KeyPressHighlight)
	wantColor(t, "FooterGlow foreground", th.FooterGlow.GetForeground(), th.KeyGlowColor())
}

// TestThemePaletteIsResolved checks that the Theme's Palette copy fills
// every optional field from the cascade and keeps the base colors as given.
func TestThemePaletteIsResolved(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	want := p
	want.Hero, want.HeroDk = p.Purple, p.PurpleDk
	want.Primary, want.PrimaryDark = p.Purple, p.PurpleDk
	want.Secondary, want.Accent = p.Blue, p.Aqua
	want.Success, want.Warning = p.Green, p.Yellow
	want.Error, want.Info = p.Red, p.Aqua

	want.TabFocusedBg, want.TabFocusedFg = p.Purple, p.Bg0
	want.TabSelectedBg, want.TabSelectedFg = p.PurpleDk, p.Bg0
	want.TabInactiveBg, want.TabInactiveFg = p.Bg1, p.Fg3
	want.NavBarBorderColor = p.Purple
	want.HeaderBorderColor = p.Purple
	want.ContentBorderColor, want.FooterBorderColor = p.Bg3, p.Bg3
	want.FooterKeyColor, want.FooterValueColor = p.Purple, p.Fg2
	want.KeyPressHighlight = p.GreenLt
	want.HeaderFillBg, want.HeaderFillFg = p.Purple, p.Bg0
	want.FooterFillBg, want.FooterFillFg = p.Purple, p.Bg0
	want.BreadcrumbBg, want.BreadcrumbFg = p.Aqua, p.Bg0
	want.BreadcrumbSelectedBg, want.BreadcrumbSelectedFg = p.Purple, p.Bg0
	want.ModalBg, want.ModalBorderColor, want.ModalTitleColor = p.Bg0, p.Purple, p.Purple
	want.ModalBorder = lipgloss.DoubleBorder()
	want.NotifBg = p.Bg0
	want.IconInfo, want.IconNote = "💬 Info:", "📌 Note:"
	want.IconTip, want.IconImportant = "💡 Tip:", "❗ Important:"
	want.IconSuccess, want.IconWarning = "✅ Success:", "🔶 Warning:"
	want.IconCaution, want.IconError = "⛔ Caution:", "❌ Error:"

	if th.Palette != want {
		t.Errorf("Theme.Palette = %+v\nwant %+v", th.Palette, want)
	}
}

// TestBuiltInThemePalettesHaveNoEmptyField walks every field of every
// built-in theme's resolved palette: none is left at its zero value, but
// for the border gradient, which has no default.
func TestBuiltInThemePalettesHaveNoEmptyField(t *testing.T) {
	for _, th := range AllThemes {
		v := reflect.ValueOf(th.Palette)
		for i := 0; i < v.NumField(); i++ {
			if v.Type().Field(i).Name == "BorderBlend" {
				continue
			}
			if v.Field(i).IsZero() {
				t.Errorf("%s: Palette.%s is empty after NewTheme", th.Name(), v.Type().Field(i).Name)
			}
		}
	}
}

func TestThemePaletteCopyKeepsExplicitFields(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	if th.Palette != p {
		t.Errorf("Theme.Palette = %+v\nwant the palette as given %+v", th.Palette, p)
	}
}

func TestNewThemeLeavesPackagePalettesUntouched(t *testing.T) {
	before := make([]Palette, len(AllPalettes))
	copy(before, AllPalettes)

	for i := range AllPalettes {
		NewTheme(AllPalettes[i])
	}

	for i := range AllPalettes {
		if AllPalettes[i] != before[i] {
			t.Errorf("AllPalettes[%d] (%s) changed after NewTheme", i, before[i].Name)
		}
	}
	if DraculaPalette != before[0] {
		t.Error("DraculaPalette changed after NewTheme")
	}
}

// TestThemeIsDark classifies every built-in theme. The light ones are listed
// here; every other theme is dark.
func TestThemeIsDark(t *testing.T) {
	light := map[*Theme]bool{DraculaLight: true, CatLatte: true, TokyoDay: true, RosePineDawn: true, GruvboxLight: true}
	for _, th := range AllThemes {
		if got, want := th.IsDark(), !light[th]; got != want {
			t.Errorf("%s.IsDark() = %v, want %v", th.Name(), got, want)
		}
	}
}

// TestCounterpart checks the light and dark pairs both ways.
func TestCounterpart(t *testing.T) {
	pairs := map[*Theme]*Theme{
		DraculaLight: Dracula, CatLatte: CatMocha, TokyoDay: TokyoNight,
		RosePineDawn: RosePine, GruvboxLight: Gruvbox,
	}
	for light, dark := range pairs {
		if Counterpart(light) != dark || Counterpart(dark) != light {
			t.Errorf("%s and %s should be each other's counterpart", light.Name(), dark.Name())
		}
		if light.IsDark() || !dark.IsDark() {
			t.Errorf("%s should be light and %s dark", light.Name(), dark.Name())
		}
	}
	for _, th := range []*Theme{Nord, Monokai, CatFrappe, nil} {
		if got := Counterpart(th); got != nil {
			t.Errorf("Counterpart = %s, want none", got.Name())
		}
	}
}
