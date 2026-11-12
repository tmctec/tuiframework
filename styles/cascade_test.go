package styles

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

// These tests build a theme from a palette that sets only tier 1 colors and
// check that every tier 2, 3 and 4 value derives from its documented source.

func TestNewThemeCascadesHeroAndSemanticColors(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		field string
		got   color.Color
		want  color.Color
	}{
		{"Hero", th.Palette.Hero, p.Purple},
		{"HeroDk", th.Palette.HeroDk, p.PurpleDk},
		{"Primary", th.Palette.Primary, p.Purple},
		{"PrimaryDark", th.Palette.PrimaryDark, p.PurpleDk},
		{"Secondary", th.Palette.Secondary, p.Blue},
		{"Accent", th.Palette.Accent, p.Aqua},
		{"Success", th.Palette.Success, p.Green},
		{"Warning", th.Palette.Warning, p.Yellow},
		{"Error", th.Palette.Error, p.Red},
		{"Info", th.Palette.Info, p.Aqua},
	}
	for _, tc := range tests {
		wantColor(t, tc.field, tc.got, tc.want)
	}
}

func TestNewThemePrimaryFollowsExplicitHero(t *testing.T) {
	p := basePalette()
	p.Hero, p.HeroDk = lipgloss.Color("#b00000"), lipgloss.Color("#b00001")
	th := NewTheme(p)

	wantColor(t, "Primary", th.Palette.Primary, p.Hero)
	wantColor(t, "PrimaryDark", th.Palette.PrimaryDark, p.HeroDk)
	wantColor(t, "TabActive background", th.TabActive.GetBackground(), p.Hero)
	wantColor(t, "TabSelectedMuted background", th.TabSelectedMuted.GetBackground(), p.HeroDk)
	wantColor(t, "AccentStyle foreground", th.AccentStyle.GetForeground(), p.Hero)
	wantColor(t, "HeaderBorderColor", th.HeaderBorderColor, p.Hero)
}

func TestNewThemeDerivesTabStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		fg    color.Color
		bg    color.Color
		bold  bool
	}{
		{"TabActive", th.TabActive, p.Bg0, p.Purple, true},
		{"TabSelectedMuted", th.TabSelectedMuted, p.Bg0, p.PurpleDk, false},
		{"TabInactive", th.TabInactive, p.Fg3, p.Bg1, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			wantStyleColors(t, tc.name, tc.style, tc.fg, tc.bg)
			wantBold(t, tc.name, tc.style, tc.bold)
			wantPadding(t, tc.name, tc.style, 0, 2, 0, 2)
			if got := tc.style.GetMarginRight(); got != 1 {
				t.Errorf("%s margin right = %d, want 1", tc.name, got)
			}
		})
	}
}

func TestNewThemeDerivesBorders(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	wantBorder(t, "NavBarBorder", th.NavBarBorder, p.Border)
	wantBorder(t, "ContentBorder", th.ContentBorder, p.Border)
	wantBorder(t, "FooterBorder", th.FooterBorder, p.Border)
	wantBorder(t, "NotifBorder", th.NotifBorder, p.Border)

	wantColor(t, "NavBarBorderColor", th.NavBarBorderColor, p.Purple)
	wantColor(t, "HeaderBorderColor", th.HeaderBorderColor, p.Purple)
	wantColor(t, "ContentBorderColor", th.ContentBorderColor, p.Bg3)
	wantColor(t, "FooterBorderColor", th.FooterBorderColor, p.Bg3)
}

func TestNewThemeDerivesFooterStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		fg    color.Color
		bold  bool
	}{
		{"Footer", th.Footer, p.Fg2, false},
		{"FooterKey", th.FooterKey, p.Purple, true},
		{"FooterValue", th.FooterValue, p.Fg2, false},
		{"FooterGlow", th.FooterGlow, p.GreenLt, true},
	}
	for _, tc := range tests {
		wantColor(t, tc.name+" foreground", tc.style.GetForeground(), tc.fg)
		wantBold(t, tc.name, tc.style, tc.bold)
	}
}

func TestNewThemeDerivesBreadcrumbStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		bg    color.Color
	}{
		{"Breadcrumb", th.Breadcrumb, p.Aqua},
		{"BreadcrumbSelected", th.BreadcrumbSelected, p.Purple},
	}
	for _, tc := range tests {
		wantStyleColors(t, tc.name, tc.style, p.Bg0, tc.bg)
		wantPadding(t, tc.name, tc.style, 0, 2, 0, 2)
		if got := tc.style.GetMarginRight(); got != 1 {
			t.Errorf("%s margin right = %d, want 1", tc.name, got)
		}
	}
}

func TestNewThemeDerivesModalStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	wantColor(t, "ModalBackground", th.ModalBackground, p.Bg0)
	wantBorder(t, "ModalBorder", th.ModalBorder, lipgloss.DoubleBorder())
	wantColor(t, "ModalBorderColor", th.ModalBorderColor, th.HeaderBorderColor)
	wantColor(t, "ModalTitle foreground", th.ModalTitle.GetForeground(), p.Purple)
	wantBold(t, "ModalTitle", th.ModalTitle, true)
}

func TestNewThemeDerivesNotificationStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		level  string
		got    NotificationStyle
		fg     color.Color
		border color.Color
		icon   string
	}{
		{"NotifInfo", th.NotifInfo, p.Fg1, p.Blue, "💬 Info:"},
		{"NotifNote", th.NotifNote, p.Aqua, p.Aqua, "📌 Note:"},
		{"NotifTip", th.NotifTip, p.Green, p.Green, "💡 Tip:"},
		{"NotifImportant", th.NotifImportant, p.Aqua, p.Aqua, "❗ Important:"},
		{"NotifSuccess", th.NotifSuccess, p.Green, p.Green, "✅ Success:"},
		{"NotifWarning", th.NotifWarning, p.Yellow, p.Yellow, "🔶 Warning:"},
		{"NotifCaution", th.NotifCaution, p.Yellow, p.Yellow, "⛔ Caution:"},
		{"NotifError", th.NotifError, p.Red, p.Red, "❌ Error:"},
	}
	for _, tc := range tests {
		t.Run(tc.level, func(t *testing.T) {
			wantColor(t, "Foreground", tc.got.Foreground, tc.fg)
			wantColor(t, "Background", tc.got.Background, p.Bg0)
			wantColor(t, "Border", tc.got.Border, tc.border)
			wantIcon(t, tc.level, tc.got.Icon, tc.icon)
		})
	}
}

func TestNewThemeDerivesTextStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		fg    color.Color
		bold  bool
	}{
		{"TextPrimary", th.TextPrimary, p.Fg1, false},
		{"TextSecondary", th.TextSecondary, p.Fg2, false},
		{"TextMuted", th.TextMuted, p.Fg3, false},
		{"TextHeading", th.TextHeading, p.Fg0, true},
	}
	for _, tc := range tests {
		wantColor(t, tc.name+" foreground", tc.style.GetForeground(), tc.fg)
		wantBold(t, tc.name, tc.style, tc.bold)
	}
}

func TestNewThemeDerivesBaseStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		fg    color.Color
	}{
		{"BaseStyle", th.BaseStyle, p.Fg0},
		{"ContentStyle", th.ContentStyle, p.Fg1},
		{"MutedStyle", th.MutedStyle, p.Fg2},
		{"DisabledStyle", th.DisabledStyle, p.Fg3},
		{"AccentStyle", th.AccentStyle, p.Purple},
	}
	for _, tc := range tests {
		wantStyleColors(t, tc.name, tc.style, tc.fg, p.Bg0)
	}
}

func TestNewThemeDerivesComponentStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	t.Run("CardStyle", func(t *testing.T) {
		wantStyleColors(t, "CardStyle", th.CardStyle, p.Fg1, p.Bg1)
		wantBorder(t, "CardStyle border", th.CardStyle.GetBorderStyle(), p.Border)
		wantColor(t, "CardStyle border foreground", th.CardStyle.GetBorderTopForeground(), p.Bg3)
		wantPadding(t, "CardStyle", th.CardStyle, 1, 2, 1, 2)
	})

	t.Run("PanelStyle", func(t *testing.T) {
		wantStyleColors(t, "PanelStyle", th.PanelStyle, p.Fg1, p.Bg0)
		wantBorder(t, "PanelStyle border", th.PanelStyle.GetBorderStyle(), lipgloss.RoundedBorder())
		wantColor(t, "PanelStyle border foreground", th.PanelStyle.GetBorderTopForeground(), p.Purple)
		wantPadding(t, "PanelStyle", th.PanelStyle, 1, 2, 1, 2)
	})

	t.Run("HeaderStyle", func(t *testing.T) {
		wantStyleColors(t, "HeaderStyle", th.HeaderStyle, p.Fg0, p.Bg0)
		wantBold(t, "HeaderStyle", th.HeaderStyle, true)
		wantBorder(t, "HeaderStyle border", th.HeaderStyle.GetBorderStyle(), lipgloss.NormalBorder())
		wantColor(t, "HeaderStyle border foreground", th.HeaderStyle.GetBorderBottomForeground(), p.Purple)
		_, top, right, bottom, left := th.HeaderStyle.GetBorder()
		if top || right || left || !bottom {
			t.Errorf("HeaderStyle border sides top/right/bottom/left = %v/%v/%v/%v, want only bottom",
				top, right, bottom, left)
		}
		if got := th.HeaderStyle.GetPaddingBottom(); got != 1 {
			t.Errorf("HeaderStyle padding bottom = %d, want 1", got)
		}
		if got := th.HeaderStyle.GetMarginBottom(); got != 1 {
			t.Errorf("HeaderStyle margin bottom = %d, want 1", got)
		}
	})

	t.Run("LabelStyle", func(t *testing.T) {
		wantStyleColors(t, "LabelStyle", th.LabelStyle, p.Bg0, p.Purple)
		wantPadding(t, "LabelStyle", th.LabelStyle, 0, 1, 0, 1)
	})
}

func TestNewThemeDerivesStatusStyles(t *testing.T) {
	p := basePalette()
	th := NewTheme(p)

	tests := []struct {
		name  string
		style lipgloss.Style
		fg    color.Color
	}{
		{"StatusSuccess", th.StatusSuccess, p.Green},
		{"StatusWarning", th.StatusWarning, p.Yellow},
		{"StatusError", th.StatusError, p.Red},
		{"StatusInfo", th.StatusInfo, p.Aqua},
	}
	for _, tc := range tests {
		wantStyleColors(t, tc.name, tc.style, tc.fg, p.Bg0)
	}
}
