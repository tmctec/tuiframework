package styles

import (
	"image/color"
	"testing"

	"charm.land/lipgloss/v2"
)

// basePalette returns a palette that sets only tier 1 colors and a border.
// Every color is distinct, so a style that cascades from the wrong source
// produces a visible mismatch.
func basePalette() Palette {
	return Palette{
		Name: "Base",

		Bg0: lipgloss.Color("#000000"), Bg1: lipgloss.Color("#000001"), Bg2: lipgloss.Color("#000002"), Bg3: lipgloss.Color("#000003"),
		Fg0: lipgloss.Color("#100000"), Fg1: lipgloss.Color("#100001"), Fg2: lipgloss.Color("#100002"), Fg3: lipgloss.Color("#100003"),

		Gray: lipgloss.Color("#200000"), GrayLt: lipgloss.Color("#200001"), GrayDk: lipgloss.Color("#200002"),
		Red: lipgloss.Color("#300000"), RedLt: lipgloss.Color("#300001"), RedDk: lipgloss.Color("#300002"),
		Orange: lipgloss.Color("#400000"), OrangeLt: lipgloss.Color("#400001"), OrangeDk: lipgloss.Color("#400002"),
		Green: lipgloss.Color("#500000"), GreenLt: lipgloss.Color("#500001"), GreenDk: lipgloss.Color("#500002"),
		Yellow: lipgloss.Color("#600000"), YellowLt: lipgloss.Color("#600001"), YellowDk: lipgloss.Color("#600002"),
		Blue: lipgloss.Color("#700000"), BlueLt: lipgloss.Color("#700001"), BlueDk: lipgloss.Color("#700002"),
		Purple: lipgloss.Color("#800000"), PurpleLt: lipgloss.Color("#800001"), PurpleDk: lipgloss.Color("#800002"),
		Pink: lipgloss.Color("#900000"), PinkLt: lipgloss.Color("#900001"), PinkDk: lipgloss.Color("#900002"),
		Aqua: lipgloss.Color("#a00000"), AquaLt: lipgloss.Color("#a00001"), AquaDk: lipgloss.Color("#a00002"),

		Border: lipgloss.RoundedBorder(),
	}
}

// fullPalette returns basePalette with every optional field set to a value
// that differs from the default it replaces.
func fullPalette() Palette {
	p := basePalette()
	p.Name = "Full"

	p.Hero, p.HeroDk = lipgloss.Color("#b00000"), lipgloss.Color("#b00001")

	p.Primary, p.PrimaryDark = lipgloss.Color("#c00000"), lipgloss.Color("#c00001")
	p.Secondary, p.Accent = lipgloss.Color("#c00002"), lipgloss.Color("#c00003")
	p.Success, p.Warning = lipgloss.Color("#c00004"), lipgloss.Color("#c00005")
	p.Error, p.Info = lipgloss.Color("#c00006"), lipgloss.Color("#c00007")

	p.TabFocusedBg, p.TabFocusedFg = lipgloss.Color("#d00000"), lipgloss.Color("#d00001")
	p.TabSelectedBg, p.TabSelectedFg = lipgloss.Color("#d00002"), lipgloss.Color("#d00003")
	p.TabInactiveBg, p.TabInactiveFg = lipgloss.Color("#d00004"), lipgloss.Color("#d00005")
	p.NavBarBorderColor = lipgloss.Color("#d00006")

	p.HeaderBorderColor = lipgloss.Color("#e00000")
	p.ContentBorderColor = lipgloss.Color("#e00001")
	p.FooterBorderColor = lipgloss.Color("#e00002")

	p.FooterKeyColor = lipgloss.Color("#e00003")
	p.FooterValueColor = lipgloss.Color("#e00004")
	p.KeyPressHighlight = lipgloss.Color("#e00005")

	p.HeaderFillBg, p.HeaderFillFg = lipgloss.Color("#e0000a"), lipgloss.Color("#e0000b")
	p.FooterFillBg, p.FooterFillFg = lipgloss.Color("#e0000c"), lipgloss.Color("#e0000d")

	p.BreadcrumbBg, p.BreadcrumbFg = lipgloss.Color("#e00006"), lipgloss.Color("#e00007")
	p.BreadcrumbSelectedBg, p.BreadcrumbSelectedFg = lipgloss.Color("#e00008"), lipgloss.Color("#e00009")

	p.ModalBg = lipgloss.Color("#f00000")
	p.ModalBorderColor = lipgloss.Color("#f00001")
	p.ModalTitleColor = lipgloss.Color("#f00002")
	p.ModalBorder = lipgloss.ThickBorder()

	p.NotifBg = lipgloss.Color("#f00003")

	p.IconInfo, p.IconNote = "[i]", "[*]"
	p.IconTip, p.IconImportant = "[?]", "[!]"
	p.IconSuccess, p.IconWarning = "[OK]", "[!!]"
	p.IconCaution, p.IconError = "[!!!]", "[ERR]"

	return p
}

// wantColor fails the test when got differs from want.
func wantColor(t *testing.T, field string, got, want color.Color) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", field, got, want)
	}
}

// wantBorder fails the test when got differs from want.
func wantBorder(t *testing.T, field string, got, want lipgloss.Border) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %+v, want %+v", field, got, want)
	}
}

// wantStyleColors checks a style's foreground and background.
func wantStyleColors(t *testing.T, name string, s lipgloss.Style, fg, bg color.Color) {
	t.Helper()
	wantColor(t, name+" foreground", s.GetForeground(), fg)
	wantColor(t, name+" background", s.GetBackground(), bg)
}

// wantBold checks a style's bold attribute.
func wantBold(t *testing.T, name string, s lipgloss.Style, bold bool) {
	t.Helper()
	if got := s.GetBold(); got != bold {
		t.Errorf("%s bold = %v, want %v", name, got, bold)
	}
}

// wantPadding checks a style's padding as top, right, bottom, left.
func wantPadding(t *testing.T, name string, s lipgloss.Style, top, right, bottom, left int) {
	t.Helper()
	gotTop, gotRight, gotBottom, gotLeft := s.GetPadding()
	if gotTop != top || gotRight != right || gotBottom != bottom || gotLeft != left {
		t.Errorf("%s padding = %d %d %d %d, want %d %d %d %d",
			name, gotTop, gotRight, gotBottom, gotLeft, top, right, bottom, left)
	}
}

// wantIcon fails the test when a notification icon differs from want.
func wantIcon(t *testing.T, level string, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("%s icon = %q, want %q", level, got, want)
	}
}
