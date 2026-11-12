package styles

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

// builtIns lists every built-in theme in registry order with its exported
// theme and palette variables.
var builtIns = []struct {
	name    string
	theme   *Theme
	palette Palette
}{
	{"Dracula", Dracula, DraculaPalette},
	{"DraculaLight", DraculaLight, DraculaLightPalette},
	{"CatMocha", CatMocha, CatppuccinMochaPalette},
	{"CatMacchiato", CatMacchiato, CatppuccinMacchiatoPalette},
	{"CatFrappe", CatFrappe, CatppuccinFrappePalette},
	{"CatLatte", CatLatte, CatppuccinLattePalette},
	{"TokyoNight", TokyoNight, TokyoNightPalette},
	{"TokyoStorm", TokyoStorm, TokyoStormPalette},
	{"TokyoMoon", TokyoMoon, TokyoMoonPalette},
	{"TokyoDay", TokyoDay, TokyoDayPalette},
	{"RosePine", RosePine, RosePinePalette},
	{"RosePineMoon", RosePineMoon, RosePineMoonPalette},
	{"RosePineDawn", RosePineDawn, RosePineDawnPalette},
	{"Retro", Retro, RetroPalette},
	{"RetroMono", RetroMono, RetroMonoPalette},
	{"Minimal", Minimal, MinimalPalette},
	{"Solarized", Solarized, SolarizedPalette},
	{"Nord", Nord, NordPalette},
	{"Monokai", Monokai, MonokaiPalette},
	{"Gruvbox", Gruvbox, GruvboxPalette},
	{"GruvboxLight", GruvboxLight, GruvboxLightPalette},
	{"Charm", Charm, CharmPalette},
}

// namedColor pairs a field name with its color for table checks.
type namedColor struct {
	field string
	color color.Color
}

// requiredPaletteColors lists the tier 1 colors every palette must set.
func requiredPaletteColors(p Palette) []namedColor {
	return []namedColor{
		{"Bg0", p.Bg0}, {"Bg1", p.Bg1}, {"Bg2", p.Bg2}, {"Bg3", p.Bg3},
		{"Fg0", p.Fg0}, {"Fg1", p.Fg1}, {"Fg2", p.Fg2}, {"Fg3", p.Fg3},
		{"Gray", p.Gray}, {"GrayLt", p.GrayLt}, {"GrayDk", p.GrayDk},
		{"Red", p.Red}, {"RedLt", p.RedLt}, {"RedDk", p.RedDk},
		{"Orange", p.Orange}, {"OrangeLt", p.OrangeLt}, {"OrangeDk", p.OrangeDk},
		{"Green", p.Green}, {"GreenLt", p.GreenLt}, {"GreenDk", p.GreenDk},
		{"Yellow", p.Yellow}, {"YellowLt", p.YellowLt}, {"YellowDk", p.YellowDk},
		{"Blue", p.Blue}, {"BlueLt", p.BlueLt}, {"BlueDk", p.BlueDk},
		{"Purple", p.Purple}, {"PurpleLt", p.PurpleLt}, {"PurpleDk", p.PurpleDk},
		{"Pink", p.Pink}, {"PinkLt", p.PinkLt}, {"PinkDk", p.PinkDk},
		{"Aqua", p.Aqua}, {"AquaLt", p.AquaLt}, {"AquaDk", p.AquaDk},
	}
}

// themeColors lists every plain color field a rendered theme uses.
func themeColors(th *Theme) []namedColor {
	notifs := []struct {
		level string
		style NotificationStyle
	}{
		{"NotifInfo", th.NotifInfo}, {"NotifNote", th.NotifNote},
		{"NotifTip", th.NotifTip}, {"NotifImportant", th.NotifImportant},
		{"NotifSuccess", th.NotifSuccess}, {"NotifWarning", th.NotifWarning},
		{"NotifCaution", th.NotifCaution}, {"NotifError", th.NotifError},
	}
	colors := make([]namedColor, 0, 6+3*len(notifs))
	colors = append(colors,
		namedColor{"NavBarBorderColor", th.NavBarBorderColor},
		namedColor{"HeaderBorderColor", th.HeaderBorderColor},
		namedColor{"ContentBorderColor", th.ContentBorderColor},
		namedColor{"FooterBorderColor", th.FooterBorderColor},
		namedColor{"ModalBackground", th.ModalBackground},
		namedColor{"ModalBorderColor", th.ModalBorderColor},
	)
	for _, n := range notifs {
		colors = append(colors,
			namedColor{n.level + " foreground", n.style.Foreground},
			namedColor{n.level + " background", n.style.Background},
			namedColor{n.level + " border", n.style.Border},
		)
	}
	return colors
}

// themeStyles lists every lipgloss style a theme exposes.
func themeStyles(th *Theme) []struct {
	name  string
	style lipgloss.Style
} {
	return []struct {
		name  string
		style lipgloss.Style
	}{
		{"TabInactive", th.TabInactive}, {"TabActive", th.TabActive},
		{"TabSelectedMuted", th.TabSelectedMuted},
		{"HeaderFill", th.HeaderFill},
		{"Footer", th.Footer}, {"FooterKey", th.FooterKey},
		{"FooterValue", th.FooterValue}, {"FooterGlow", th.FooterGlow},
		{"FooterFill", th.FooterFill},
		{"Breadcrumb", th.Breadcrumb}, {"BreadcrumbSelected", th.BreadcrumbSelected},
		{"ModalTitle", th.ModalTitle},
		{"TextPrimary", th.TextPrimary}, {"TextSecondary", th.TextSecondary},
		{"TextMuted", th.TextMuted}, {"TextHeading", th.TextHeading},
		{"BaseStyle", th.BaseStyle}, {"ContentStyle", th.ContentStyle},
		{"MutedStyle", th.MutedStyle}, {"DisabledStyle", th.DisabledStyle},
		{"AccentStyle", th.AccentStyle},
		{"CardStyle", th.CardStyle}, {"PanelStyle", th.PanelStyle},
		{"HeaderStyle", th.HeaderStyle}, {"LabelStyle", th.LabelStyle},
		{"StatusSuccess", th.StatusSuccess}, {"StatusWarning", th.StatusWarning},
		{"StatusError", th.StatusError}, {"StatusInfo", th.StatusInfo},
	}
}

func TestBuiltInRegistriesMatch(t *testing.T) {
	if len(AllThemes) != len(builtIns) {
		t.Fatalf("len(AllThemes) = %d, want %d", len(AllThemes), len(builtIns))
	}
	if len(AllPalettes) != len(builtIns) {
		t.Fatalf("len(AllPalettes) = %d, want %d", len(AllPalettes), len(builtIns))
	}
	for i, b := range builtIns {
		if AllThemes[i] != b.theme {
			t.Errorf("AllThemes[%d] = %s, want %s", i, AllThemes[i].Name(), b.name)
		}
		if AllPalettes[i].Name != b.name {
			t.Errorf("AllPalettes[%d].Name = %q, want %q", i, AllPalettes[i].Name, b.name)
		}
		if AllPalettes[i] != b.palette {
			t.Errorf("AllPalettes[%d] differs from the exported %s palette", i, b.name)
		}
		if AllThemes[i].Palette.Name != AllPalettes[i].Name {
			t.Errorf("AllThemes[%d] is %q but AllPalettes[%d] is %q",
				i, AllThemes[i].Palette.Name, i, AllPalettes[i].Name)
		}
		if b.palette.Name != b.name {
			t.Errorf("%s palette has Name %q", b.name, b.palette.Name)
		}
	}
}

func TestBuiltInNamesAreUnique(t *testing.T) {
	seen := make(map[string]int, len(AllPalettes))
	for i := range AllPalettes {
		name := AllPalettes[i].Name
		if first, dup := seen[name]; dup {
			t.Errorf("AllPalettes[%d] and AllPalettes[%d] share the name %q", first, i, name)
		}
		seen[name] = i
	}
}

func TestBuiltInPalettesSetRequiredFields(t *testing.T) {
	for i := range AllPalettes {
		p := AllPalettes[i]
		t.Run(p.Name, func(t *testing.T) {
			if p.Name == "" {
				t.Error("Name is empty")
			}
			if p.Border == (lipgloss.Border{}) {
				t.Error("Border is unset")
			}
			for _, c := range requiredPaletteColors(p) {
				if !isSet(c.color) {
					t.Errorf("%s is empty", c.field)
				}
			}
		})
	}
}

func TestBuiltInThemesRenderNoEmptyColor(t *testing.T) {
	for _, th := range AllThemes {
		t.Run(th.Name(), func(t *testing.T) {
			for _, c := range themeColors(th) {
				if !isSet(c.color) {
					t.Errorf("%s is empty", c.field)
				}
			}
			for _, s := range themeStyles(th) {
				if !isSet(s.style.GetForeground()) {
					t.Errorf("%s has no foreground color", s.name)
				}
			}
			for _, n := range []struct {
				level string
				icon  string
			}{
				{"NotifInfo", th.NotifInfo.Icon}, {"NotifNote", th.NotifNote.Icon},
				{"NotifTip", th.NotifTip.Icon}, {"NotifImportant", th.NotifImportant.Icon},
				{"NotifSuccess", th.NotifSuccess.Icon}, {"NotifWarning", th.NotifWarning.Icon},
				{"NotifCaution", th.NotifCaution.Icon}, {"NotifError", th.NotifError.Icon},
			} {
				if n.icon == "" {
					t.Errorf("%s icon is empty", n.level)
				}
			}
		})
	}
}

func TestThemeByNameFindsEveryBuiltIn(t *testing.T) {
	for _, b := range builtIns {
		if got := ThemeByName(b.name); got != b.theme {
			t.Errorf("ThemeByName(%q) = %v, want the exported %s theme", b.name, got, b.name)
		}
	}
}

func TestThemeByNameReturnsNilForUnknown(t *testing.T) {
	for _, name := range []string{"", "dracula", "NoSuchTheme"} {
		if got := ThemeByName(name); got != nil {
			t.Errorf("ThemeByName(%q) = %s, want nil", name, got.Name())
		}
	}
}

func TestPaletteByNameReturnsPointerIntoAllPalettes(t *testing.T) {
	for i, b := range builtIns {
		got := PaletteByName(b.name)
		if got == nil {
			t.Errorf("PaletteByName(%q) = nil", b.name)
			continue
		}
		if got != &AllPalettes[i] {
			t.Errorf("PaletteByName(%q) points at %p, want &AllPalettes[%d] %p", b.name, got, i, &AllPalettes[i])
		}
		if *got != b.palette {
			t.Errorf("PaletteByName(%q) differs from the exported %s palette", b.name, b.name)
		}
	}
}

func TestPaletteByNameReturnsNilForUnknown(t *testing.T) {
	for _, name := range []string{"", "dracula", "NoSuchPalette"} {
		if got := PaletteByName(name); got != nil {
			t.Errorf("PaletteByName(%q) = %s, want nil", name, got.Name)
		}
	}
}
