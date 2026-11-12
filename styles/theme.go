package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// NotificationStyle defines colors and icon for a notification level.
type NotificationStyle struct {
	Foreground color.Color
	Background color.Color
	Border     color.Color
	Icon       string
}

// Theme contains all computed styles for the UI.
// Create via NewTheme(palette) or use built-in themes (Dracula, Gruvbox, etc.).
type Theme struct {
	// Palette is the source palette with every optional field resolved
	Palette Palette

	// === TAB BAR ===
	TabInactive       lipgloss.Style  // Stationary tabs
	TabActive         lipgloss.Style  // Currently navigating
	TabSelectedMuted  lipgloss.Style  // Selected but not navigating
	NavBarBorder      lipgloss.Border // Border style for nav bar
	NavBarBorderColor color.Color     // Border color for nav bar

	// === HEADER ===
	HeaderBorderColor color.Color    // Border color below header/tabs
	HeaderFill        lipgloss.Style // Filled header band: reversed text on the fill color

	// === CONTENT ===
	ContentBorder      lipgloss.Border // Border style for content area
	ContentBorderColor color.Color     // Border color for content area
	BorderBlend        []color.Color   // Gradient for the content, nav bar and modal borders; empty for none

	// === FOOTER ===
	Footer            lipgloss.Style  // Base footer style
	FooterBorder      lipgloss.Border // Border style for footer
	FooterBorderColor color.Color     // Border color for footer
	FooterKey         lipgloss.Style  // Key labels style
	FooterValue       lipgloss.Style  // Value text style
	FooterGlow        lipgloss.Style  // Pressed key highlight
	FooterFill        lipgloss.Style  // Filled footer: reversed text on the fill color

	// === BREADCRUMBS ===
	Breadcrumb         lipgloss.Style // Inactive breadcrumb
	BreadcrumbSelected lipgloss.Style // Active breadcrumb

	// === NOTIFICATIONS (8 levels) ===
	NotifInfo      NotificationStyle
	NotifNote      NotificationStyle
	NotifTip       NotificationStyle
	NotifImportant NotificationStyle
	NotifSuccess   NotificationStyle
	NotifWarning   NotificationStyle
	NotifCaution   NotificationStyle
	NotifError     NotificationStyle
	NotifBorder    lipgloss.Border // Shared border style for notifications

	// === MODALS ===
	ModalBackground  color.Color     // Modal background color
	ModalBorder      lipgloss.Border // Modal border style
	ModalBorderColor color.Color     // Modal border color
	ModalTitle       lipgloss.Style  // Modal title style

	// === TEXT ===
	TextPrimary   lipgloss.Style // Primary text (Fg1)
	TextSecondary lipgloss.Style // Secondary text (Fg2)
	TextMuted     lipgloss.Style // Muted text (Fg3)
	TextHeading   lipgloss.Style // Heading text (Fg0, bold)

	// === CELL/FLEXBOX STYLES (for StylePassing inheritance) ===
	BaseStyle     lipgloss.Style // Fg0/Bg0 - ideal for FlexBox-level inheritance
	ContentStyle  lipgloss.Style // Fg1/Bg0 - primary body text
	MutedStyle    lipgloss.Style // Fg2/Bg0 - secondary content
	DisabledStyle lipgloss.Style // Fg3/Bg0 - disabled/tertiary
	AccentStyle   lipgloss.Style // Hero/Bg0 - highlighted content

	// === COMPONENT PRESETS ===
	CardStyle   lipgloss.Style // Elevated container with subtle border
	PanelStyle  lipgloss.Style // Bordered panel with primary accent
	HeaderStyle lipgloss.Style // Section header with underline
	LabelStyle  lipgloss.Style // Pill-shaped badge/label

	// === STATUS STYLES ===
	StatusSuccess lipgloss.Style // Success indicator (green)
	StatusWarning lipgloss.Style // Warning indicator (yellow)
	StatusError   lipgloss.Style // Error indicator (red)
	StatusInfo    lipgloss.Style // Info indicator (blue/aqua)
}

// borderOr returns the border if set (non-empty Top), otherwise the fallback.
// This enables optional border style fields with documented defaults.
func borderOr(border, fallback lipgloss.Border) lipgloss.Border {
	if border.Top == "" {
		return fallback
	}
	return border
}

// NewTheme creates a complete theme from a palette.
// All styles are computed automatically from the palette colors.
//
// Color resolution order (4-tier cascade):
//  1. Base colors (Gray, Red, Orange, Green, Yellow, Blue, Purple, Pink, Aqua + Lt/Dk variants)
//  2. Hero colors (Hero, HeroDk) - the theme's signature accent colors
//  3. Semantic colors (Primary, Secondary, Accent, Success, etc.) - default to hero or base
//  4. UI element colors (TabFocusedBg, ContentBorderColor, etc.) - default to semantic
//
// The Theme's Palette is the resolved palette: every optional field holds the
// color, border or icon the cascade gave it, so a reader of theme.Palette
// never sees an empty one.
func NewTheme(palette Palette) *Theme {
	p := resolve(palette)

	return &Theme{
		Palette: p,

		// Inactive: muted appearance for unselected tabs
		TabInactive: lipgloss.NewStyle().
			Background(p.TabInactiveBg).
			Foreground(p.TabInactiveFg).
			Padding(0, 2).
			MarginRight(1),
		// Focused: brightest (Primary accent by default)
		TabActive: lipgloss.NewStyle().
			Background(p.TabFocusedBg).
			Foreground(p.TabFocusedFg).
			Padding(0, 2).
			MarginRight(1).
			Bold(true),
		// Selected: darker version of focused (PrimaryDark by default)
		TabSelectedMuted: lipgloss.NewStyle().
			Background(p.TabSelectedBg).
			Foreground(p.TabSelectedFg).
			Padding(0, 2).
			MarginRight(1),
		NavBarBorder:      p.Border,
		NavBarBorderColor: p.NavBarBorderColor,

		HeaderBorderColor: p.HeaderBorderColor,
		HeaderFill:        lipgloss.NewStyle().Foreground(p.HeaderFillFg).Background(p.HeaderFillBg),

		ContentBorder:      p.Border,
		ContentBorderColor: p.ContentBorderColor,
		BorderBlend:        borderBlend(p),

		Footer:            lipgloss.NewStyle().Foreground(p.FooterValueColor),
		FooterBorder:      p.Border,
		FooterBorderColor: p.FooterBorderColor,
		FooterKey:         lipgloss.NewStyle().Foreground(p.FooterKeyColor).Bold(true),
		FooterValue:       lipgloss.NewStyle().Foreground(p.FooterValueColor),
		FooterGlow:        lipgloss.NewStyle().Foreground(p.KeyPressHighlight).Bold(true),
		FooterFill:        lipgloss.NewStyle().Foreground(p.FooterFillFg).Background(p.FooterFillBg),

		Breadcrumb: lipgloss.NewStyle().
			Background(p.BreadcrumbBg).
			Foreground(p.BreadcrumbFg).
			Padding(0, 2).
			MarginRight(1),
		BreadcrumbSelected: lipgloss.NewStyle().
			Background(p.BreadcrumbSelectedBg).
			Foreground(p.BreadcrumbSelectedFg).
			Padding(0, 2).
			MarginRight(1),

		// Notifications share NotifBg: Bg0 blends with the content, Bg1 or
		// Bg2 in the palette gives a filled look.
		NotifInfo:      NotificationStyle{p.Fg1, p.NotifBg, p.Secondary, p.IconInfo},
		NotifNote:      NotificationStyle{p.Info, p.NotifBg, p.Info, p.IconNote},
		NotifTip:       NotificationStyle{p.Success, p.NotifBg, p.Success, p.IconTip},
		NotifImportant: NotificationStyle{p.Accent, p.NotifBg, p.Accent, p.IconImportant},
		NotifSuccess:   NotificationStyle{p.Success, p.NotifBg, p.Success, p.IconSuccess},
		NotifWarning:   NotificationStyle{p.Warning, p.NotifBg, p.Warning, p.IconWarning},
		NotifCaution:   NotificationStyle{p.Warning, p.NotifBg, p.Warning, p.IconCaution},
		NotifError:     NotificationStyle{p.Error, p.NotifBg, p.Error, p.IconError},
		NotifBorder:    p.Border,

		ModalBackground:  p.ModalBg,
		ModalBorder:      p.ModalBorder,
		ModalBorderColor: p.ModalBorderColor,
		ModalTitle:       lipgloss.NewStyle().Bold(true).Foreground(p.ModalTitleColor),

		TextPrimary:   lipgloss.NewStyle().Foreground(p.Fg1),
		TextSecondary: lipgloss.NewStyle().Foreground(p.Fg2),
		TextMuted:     lipgloss.NewStyle().Foreground(p.Fg3),
		TextHeading:   lipgloss.NewStyle().Foreground(p.Fg0).Bold(true),

		// Cell/FlexBox styles for inheritance
		BaseStyle:     lipgloss.NewStyle().Foreground(p.Fg0).Background(p.Bg0),
		ContentStyle:  lipgloss.NewStyle().Foreground(p.Fg1).Background(p.Bg0),
		MutedStyle:    lipgloss.NewStyle().Foreground(p.Fg2).Background(p.Bg0),
		DisabledStyle: lipgloss.NewStyle().Foreground(p.Fg3).Background(p.Bg0),
		AccentStyle:   lipgloss.NewStyle().Foreground(p.Hero).Background(p.Bg0),

		// Component presets
		CardStyle: lipgloss.NewStyle().
			Foreground(p.Fg1).
			Background(p.Bg1).
			Border(p.Border).
			BorderForeground(p.Bg3).
			Padding(1, 2),
		PanelStyle: lipgloss.NewStyle().
			Foreground(p.Fg1).
			Background(p.Bg0).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(p.Primary).
			Padding(1, 2),
		HeaderStyle: lipgloss.NewStyle().
			Foreground(p.Fg0).
			Background(p.Bg0).
			Bold(true).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(p.Primary).
			PaddingBottom(1).
			MarginBottom(1),
		LabelStyle: lipgloss.NewStyle().
			Foreground(p.HeaderFillFg).
			Background(p.Primary).
			Padding(0, 1),

		// Status styles
		StatusSuccess: lipgloss.NewStyle().Foreground(p.Success).Background(p.Bg0),
		StatusWarning: lipgloss.NewStyle().Foreground(p.Warning).Background(p.Bg0),
		StatusError:   lipgloss.NewStyle().Foreground(p.Error).Background(p.Bg0),
		StatusInfo:    lipgloss.NewStyle().Foreground(p.Info).Background(p.Bg0),
	}
}

// resolve returns p with every optional field filled in. Each tier defaults
// to the one before it, so this is the only place the cascade is written.
func resolve(p Palette) Palette {
	// Tier 2: the hero colors, the theme's signature accent.
	p.Hero = colorOr(p.Hero, p.Purple)
	p.HeroDk = colorOr(p.HeroDk, p.PurpleDk)

	// Tier 3: semantic colors.
	p.Primary = colorOr(p.Primary, p.Hero)
	p.PrimaryDark = colorOr(p.PrimaryDark, p.HeroDk)
	p.Secondary = colorOr(p.Secondary, p.Blue)
	p.Accent = colorOr(p.Accent, p.Aqua)
	p.Success = colorOr(p.Success, p.Green)
	p.Warning = colorOr(p.Warning, p.Yellow)
	p.Error = colorOr(p.Error, p.Red)
	p.Info = colorOr(p.Info, p.Aqua)

	// Tier 4: UI elements. A focused tab is the base background on the hero
	// color; the breadcrumbs and the band fills use the same reversed pairing.
	p.TabFocusedBg = colorOr(p.TabFocusedBg, p.Primary)
	p.TabFocusedFg = colorOr(p.TabFocusedFg, p.Bg0)
	p.TabSelectedBg = colorOr(p.TabSelectedBg, p.PrimaryDark)
	p.TabSelectedFg = colorOr(p.TabSelectedFg, p.Bg0)
	p.TabInactiveBg = colorOr(p.TabInactiveBg, p.Bg1)
	p.TabInactiveFg = colorOr(p.TabInactiveFg, p.Fg3)
	p.NavBarBorderColor = colorOr(p.NavBarBorderColor, p.Primary)

	// The header border takes the hero color; content and footer stay subtle.
	p.HeaderBorderColor = colorOr(p.HeaderBorderColor, p.Primary)
	p.ContentBorderColor = colorOr(p.ContentBorderColor, p.Bg3)
	p.FooterBorderColor = colorOr(p.FooterBorderColor, p.Bg3)

	p.FooterKeyColor = colorOr(p.FooterKeyColor, p.Primary)
	p.FooterValueColor = colorOr(p.FooterValueColor, p.Fg2)
	p.KeyPressHighlight = colorOr(p.KeyPressHighlight, p.GreenLt)

	p.HeaderFillBg = colorOr(p.HeaderFillBg, p.Primary)
	p.HeaderFillFg = colorOr(p.HeaderFillFg, p.Bg0)
	p.FooterFillBg = colorOr(p.FooterFillBg, p.Primary)
	p.FooterFillFg = colorOr(p.FooterFillFg, p.Bg0)

	p.BreadcrumbBg = colorOr(p.BreadcrumbBg, p.Accent)
	p.BreadcrumbFg = colorOr(p.BreadcrumbFg, p.Bg0)
	p.BreadcrumbSelectedBg = colorOr(p.BreadcrumbSelectedBg, p.Primary)
	p.BreadcrumbSelectedFg = colorOr(p.BreadcrumbSelectedFg, p.Bg0)

	p.ModalBg = colorOr(p.ModalBg, p.Bg0)
	p.ModalBorderColor = colorOr(p.ModalBorderColor, p.HeaderBorderColor)
	p.ModalTitleColor = colorOr(p.ModalTitleColor, p.Primary)
	p.ModalBorder = borderOr(p.ModalBorder, lipgloss.DoubleBorder())

	p.NotifBg = colorOr(p.NotifBg, p.Bg0)

	p.IconInfo = iconOr(p.IconInfo, "💬 Info:")
	p.IconNote = iconOr(p.IconNote, "📌 Note:")
	p.IconTip = iconOr(p.IconTip, "💡 Tip:")
	p.IconImportant = iconOr(p.IconImportant, "❗ Important:")
	p.IconSuccess = iconOr(p.IconSuccess, "✅ Success:")
	p.IconWarning = iconOr(p.IconWarning, "🔶 Warning:")
	p.IconCaution = iconOr(p.IconCaution, "⛔ Caution:")
	p.IconError = iconOr(p.IconError, "❌ Error:")

	return p
}

// borderBlend returns the palette's border gradient as the colors a style
// takes, or nil when the palette does not set both ends.
func borderBlend(p Palette) []color.Color {
	if !isSet(p.BorderBlend[0]) || !isSet(p.BorderBlend[1]) {
		return nil
	}
	return []color.Color{p.BorderBlend[0], p.BorderBlend[1]}
}

// iconOr returns the icon if set, otherwise the fallback.
func iconOr(icon, fallback string) string {
	if icon == "" {
		return fallback
	}
	return icon
}

// Name returns the theme name from the underlying palette.
func (t *Theme) Name() string {
	return t.Palette.Name
}

// IsDark reports whether the theme is a dark one: whether its base
// background is darker than mid lightness. It is the test Bubble Tea applies
// to a terminal's background, so a theme and a terminal are classified the
// same way.
func (t *Theme) IsDark() bool {
	return isDark(t.Palette.Bg0)
}

// KeyGlowColor returns the color a pressed key takes in the footer: the
// palette's KeyPressHighlight, which defaults to GreenLt.
func (t *Theme) KeyGlowColor() color.Color {
	return t.Palette.KeyPressHighlight
}
