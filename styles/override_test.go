package styles

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

// These tests build a theme from a palette that sets every optional field and
// check that each explicit value replaces the default it would cascade from.

func TestNewThemeUsesExplicitSemanticColors(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	tests := []struct {
		field string
		got   color.Color
		want  color.Color
	}{
		{"Hero", th.Palette.Hero, p.Hero},
		{"HeroDk", th.Palette.HeroDk, p.HeroDk},
		{"Primary", th.Palette.Primary, p.Primary},
		{"PrimaryDark", th.Palette.PrimaryDark, p.PrimaryDark},
		{"Secondary", th.Palette.Secondary, p.Secondary},
		{"Accent", th.Palette.Accent, p.Accent},
		{"Success", th.Palette.Success, p.Success},
		{"Warning", th.Palette.Warning, p.Warning},
		{"Error", th.Palette.Error, p.Error},
		{"Info", th.Palette.Info, p.Info},
	}
	for _, tc := range tests {
		wantColor(t, tc.field, tc.got, tc.want)
	}

	// Styles that read semantic colors directly follow the explicit values.
	wantColor(t, "AccentStyle foreground", th.AccentStyle.GetForeground(), p.Hero)
	wantColor(t, "PanelStyle border foreground", th.PanelStyle.GetBorderTopForeground(), p.Primary)
	wantColor(t, "HeaderStyle border foreground", th.HeaderStyle.GetBorderBottomForeground(), p.Primary)
	wantColor(t, "LabelStyle background", th.LabelStyle.GetBackground(), p.Primary)
	wantColor(t, "StatusSuccess foreground", th.StatusSuccess.GetForeground(), p.Success)
	wantColor(t, "StatusWarning foreground", th.StatusWarning.GetForeground(), p.Warning)
	wantColor(t, "StatusError foreground", th.StatusError.GetForeground(), p.Error)
	wantColor(t, "StatusInfo foreground", th.StatusInfo.GetForeground(), p.Info)
}

func TestNewThemeUsesExplicitTabColors(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	wantStyleColors(t, "TabActive", th.TabActive, p.TabFocusedFg, p.TabFocusedBg)
	wantStyleColors(t, "TabSelectedMuted", th.TabSelectedMuted, p.TabSelectedFg, p.TabSelectedBg)
	wantStyleColors(t, "TabInactive", th.TabInactive, p.TabInactiveFg, p.TabInactiveBg)
	wantColor(t, "NavBarBorderColor", th.NavBarBorderColor, p.NavBarBorderColor)
}

func TestNewThemeUsesExplicitBorderColors(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	wantColor(t, "HeaderBorderColor", th.HeaderBorderColor, p.HeaderBorderColor)
	wantColor(t, "ContentBorderColor", th.ContentBorderColor, p.ContentBorderColor)
	wantColor(t, "FooterBorderColor", th.FooterBorderColor, p.FooterBorderColor)
}

func TestNewThemeUsesExplicitFooterColors(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	wantColor(t, "Footer foreground", th.Footer.GetForeground(), p.FooterValueColor)
	wantColor(t, "FooterKey foreground", th.FooterKey.GetForeground(), p.FooterKeyColor)
	wantColor(t, "FooterValue foreground", th.FooterValue.GetForeground(), p.FooterValueColor)
	wantColor(t, "FooterGlow foreground", th.FooterGlow.GetForeground(), p.KeyPressHighlight)
}

func TestNewThemeUsesExplicitBreadcrumbColors(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	wantStyleColors(t, "Breadcrumb", th.Breadcrumb, p.BreadcrumbFg, p.BreadcrumbBg)
	wantStyleColors(t, "BreadcrumbSelected", th.BreadcrumbSelected, p.BreadcrumbSelectedFg, p.BreadcrumbSelectedBg)
}

func TestNewThemeUsesExplicitModalSettings(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	wantColor(t, "ModalBackground", th.ModalBackground, p.ModalBg)
	wantBorder(t, "ModalBorder", th.ModalBorder, lipgloss.ThickBorder())
	wantColor(t, "ModalBorderColor", th.ModalBorderColor, p.ModalBorderColor)
	wantColor(t, "ModalTitle foreground", th.ModalTitle.GetForeground(), p.ModalTitleColor)
}

func TestNewThemeUsesExplicitNotificationSettings(t *testing.T) {
	p := fullPalette()
	th := NewTheme(p)

	tests := []struct {
		level  string
		got    NotificationStyle
		fg     color.Color
		border color.Color
		icon   string
	}{
		{"NotifInfo", th.NotifInfo, p.Fg1, p.Secondary, p.IconInfo},
		{"NotifNote", th.NotifNote, p.Info, p.Info, p.IconNote},
		{"NotifTip", th.NotifTip, p.Success, p.Success, p.IconTip},
		{"NotifImportant", th.NotifImportant, p.Accent, p.Accent, p.IconImportant},
		{"NotifSuccess", th.NotifSuccess, p.Success, p.Success, p.IconSuccess},
		{"NotifWarning", th.NotifWarning, p.Warning, p.Warning, p.IconWarning},
		{"NotifCaution", th.NotifCaution, p.Warning, p.Warning, p.IconCaution},
		{"NotifError", th.NotifError, p.Error, p.Error, p.IconError},
	}
	for _, tc := range tests {
		t.Run(tc.level, func(t *testing.T) {
			wantColor(t, "Foreground", tc.got.Foreground, tc.fg)
			wantColor(t, "Background", tc.got.Background, p.NotifBg)
			wantColor(t, "Border", tc.got.Border, tc.border)
			wantIcon(t, tc.level, tc.got.Icon, tc.icon)
		})
	}
}

// TestNewThemeTierFourFollowsTierThree sets tier 3 colors without any tier 4
// color and checks that UI elements take the tier 3 value, not the Hero or
// base color it would otherwise fall through to.
func TestNewThemeTierFourFollowsTierThree(t *testing.T) {
	p := basePalette()
	p.Primary, p.PrimaryDark = lipgloss.Color("#c00000"), lipgloss.Color("#c00001")
	p.Accent = lipgloss.Color("#c00003")
	p.HeaderBorderColor = lipgloss.Color("#e00000")
	th := NewTheme(p)

	tests := []struct {
		field string
		got   color.Color
		want  color.Color
	}{
		{"TabActive background", th.TabActive.GetBackground(), p.Primary},
		{"TabSelectedMuted background", th.TabSelectedMuted.GetBackground(), p.PrimaryDark},
		{"NavBarBorderColor", th.NavBarBorderColor, p.Primary},
		{"FooterKey foreground", th.FooterKey.GetForeground(), p.Primary},
		{"Breadcrumb background", th.Breadcrumb.GetBackground(), p.Accent},
		{"BreadcrumbSelected background", th.BreadcrumbSelected.GetBackground(), p.Primary},
		{"ModalBorderColor", th.ModalBorderColor, p.HeaderBorderColor},
		{"ModalTitle foreground", th.ModalTitle.GetForeground(), p.Primary},
		{"NotifImportant border", th.NotifImportant.Border, p.Accent},
		{"LabelStyle background", th.LabelStyle.GetBackground(), p.Primary},
		{"AccentStyle foreground", th.AccentStyle.GetForeground(), p.Purple},
	}
	for _, tc := range tests {
		wantColor(t, tc.field, tc.got, tc.want)
	}
}
