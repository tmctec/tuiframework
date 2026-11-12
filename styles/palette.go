package styles

import (
	"image/color"

	"charm.land/lipgloss/v2"
)

// Palette defines the base colors for a theme.
// All UI styles are computed from these values via NewTheme().
// Users can create custom themes by defining a Palette and calling NewTheme().
//
// # Color Hierarchy (4-tier cascade)
//
//  1. BASE COLORS: Gray, Red, Orange, Green, Yellow, Blue, Purple, Pink, Aqua
//     Each has three variants: normal, Lt (light), Dk (dark)
//
//  2. HERO COLORS: The signature accent colors for this theme
//     Hero (bright) and HeroDk (muted) define the theme's personality
//
//  3. SEMANTIC COLORS: Primary, Secondary, Accent, Success, Warning, Error, Info
//     Default to Hero colors or base colors if not set
//
//  4. UI ELEMENT COLORS: TabFocusedBg, ContentBorderColor, etc.
//     Default to semantic colors. Tabs get colors from Hero via Primary.
//
// This allows themes to define just base colors + hero, and let everything derive.
type Palette struct {
	// Name identifies the palette (e.g., "Dracula", "Gruvbox", "Nord")
	Name string

	// === BACKGROUND LEVELS (layered depth) ===
	Bg0 color.Color // Base background (darkest)
	Bg1 color.Color // Elevated surfaces (tabs, cards)
	Bg2 color.Color // Selected/hover state
	Bg3 color.Color // Borders, subtle elements

	// === FOREGROUND LEVELS (text hierarchy) ===
	Fg0 color.Color // Brightest (headings, emphasis)
	Fg1 color.Color // Primary text
	Fg2 color.Color // Secondary/muted text
	Fg3 color.Color // Disabled/hint text

	// === BASE COLORS (the foundation - each with light/dark variants) ===
	Gray   color.Color // Neutral gray
	GrayLt color.Color // Light gray
	GrayDk color.Color // Dark gray

	Red   color.Color // Base red
	RedLt color.Color // Light red
	RedDk color.Color // Dark red

	Orange   color.Color // Base orange
	OrangeLt color.Color // Light orange
	OrangeDk color.Color // Dark orange

	Green   color.Color // Base green
	GreenLt color.Color // Light green
	GreenDk color.Color // Dark green

	Yellow   color.Color // Base yellow
	YellowLt color.Color // Light yellow
	YellowDk color.Color // Dark yellow

	Blue   color.Color // Base blue
	BlueLt color.Color // Light blue
	BlueDk color.Color // Dark blue

	Purple   color.Color // Base purple/magenta
	PurpleLt color.Color // Light purple
	PurpleDk color.Color // Dark purple

	Pink   color.Color // Base pink/magenta
	PinkLt color.Color // Light pink
	PinkDk color.Color // Dark pink

	Aqua   color.Color // Base aqua/cyan/teal
	AquaLt color.Color // Light aqua
	AquaDk color.Color // Dark aqua

	// === HERO COLORS (the signature accent colors for this theme) ===
	// Hero colors define the theme's personality and cascade to tabs/borders.
	// NOTE: Tab styles (focused, selected) derive their colors from Hero/HeroDk.
	//       Header border also uses Hero by default.
	Hero   color.Color // Primary hero color (e.g., Purple for Dracula). Default: Purple
	HeroDk color.Color // Darker hero variant for selected states. Default: PurpleDk

	// === SEMANTIC COLORS (optional - default to hero/base colors) ===
	Primary     color.Color // Main accent. Default: Hero
	PrimaryDark color.Color // Darker variant. Default: HeroDk
	Secondary   color.Color // Secondary accent. Default: Blue
	Accent      color.Color // Tertiary accent. Default: Aqua
	Success     color.Color // Success, tips. Default: Green
	Warning     color.Color // Warnings. Default: Yellow
	Error       color.Color // Errors, caution. Default: Red
	Info        color.Color // Info, notes. Default: Aqua

	// === TAB COLORS (optional - cascade from Hero colors) ===
	// These control nav bar appearance. If not set, they cascade from Hero:
	//   TabFocusedBg  <- Primary <- Hero
	//   TabSelectedBg <- PrimaryDark <- HeroDk
	TabFocusedBg      color.Color // Tab bg when focused. Default: Primary (from Hero)
	TabFocusedFg      color.Color // Tab text when focused. Default: Bg0
	TabSelectedBg     color.Color // Tab bg when selected but not focused. Default: PrimaryDark (from HeroDk)
	TabSelectedFg     color.Color // Tab text when selected but not focused. Default: Bg0
	TabInactiveBg     color.Color // Tab bg for unselected tabs. Default: Bg1
	TabInactiveFg     color.Color // Tab text for unselected tabs. Default: Fg3
	NavBarBorderColor color.Color // Border around nav bar. Default: Primary (from Hero)

	// === BORDER COLORS (optional - cascade from Hero or Bg3) ===
	// Header border uses Hero color for accent; content/footer use subtle Bg3.
	HeaderBorderColor  color.Color // Border below header/tabs. Default: Primary (from Hero)
	ContentBorderColor color.Color // Border around content area. Default: Bg3 (subtle)
	FooterBorderColor  color.Color // Border above footer. Default: Bg3 (subtle)

	// === FOOTER (optional) ===
	FooterKeyColor    color.Color // Footer key labels. Default: Primary
	FooterValueColor  color.Color // Footer value text. Default: Fg2
	KeyPressHighlight color.Color // Key press highlight in footer. Default: GreenLt

	// === BAND FILLS (optional) ===
	// A header or footer drawn as a filled bar reverses its text: the fill
	// takes the Hero color and the text the base background.
	HeaderFillBg color.Color // Filled header background. Default: Primary (from Hero)
	HeaderFillFg color.Color // Text on a filled header and on labels. Default: Bg0
	FooterFillBg color.Color // Filled footer background. Default: Primary (from Hero)
	FooterFillFg color.Color // Text on a filled footer. Default: Bg0

	// === BREADCRUMBS (optional) ===
	BreadcrumbBg         color.Color // Inactive breadcrumb background. Default: Accent
	BreadcrumbFg         color.Color // Inactive breadcrumb text. Default: Bg0
	BreadcrumbSelectedBg color.Color // Active breadcrumb background. Default: Primary
	BreadcrumbSelectedFg color.Color // Active breadcrumb text. Default: Bg0

	// === MODALS (optional) ===
	ModalBg          color.Color // Modal background. Default: Bg0
	ModalBorderColor color.Color // Modal border color. Default: HeaderBorderColor (Primary)
	ModalTitleColor  color.Color // Modal title text. Default: Primary

	// === NOTIFICATIONS (optional) ===
	// Controls notification visibility. Default: Bg0 (transparent/blends with content).
	// Set to Bg1 or Bg2 for a filled/contrasting appearance.
	NotifBg color.Color // Notification background. Default: Bg0

	// === BORDER STYLES ===
	Border      lipgloss.Border // Default border style (Rounded, Double, Normal, Hidden)
	ModalBorder lipgloss.Border // Modal border style. Default: DoubleBorder()

	// === BORDER GRADIENT (optional) ===
	// Two colors draw the content, nav bar and modal borders as a gradient
	// from the first to the second, in place of the single border colors.
	// Default: no gradient.
	BorderBlend [2]color.Color

	// === NOTIFICATION ICONS (per-theme customizable) ===
	IconInfo      string // e.g., "💬 Info:" or "[i]"
	IconNote      string // e.g., "📌 Note:" or "[*]"
	IconTip       string // e.g., "💡 Tip:" or "[?]"
	IconImportant string // e.g., "❗ Important:" or "[!]"
	IconSuccess   string // e.g., "✅ Success:" or "[OK]"
	IconWarning   string // e.g., "🔶 Warning:" or "[!!]"
	IconCaution   string // e.g., "⛔ Caution:" or "[!!!]"
	IconError     string // e.g., "❌ Error:" or "[ERR]"
}
