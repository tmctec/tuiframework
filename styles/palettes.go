package styles

import "charm.land/lipgloss/v2"

// Built-in palettes for themes.
// Users can use these directly with NewTheme() or clone and modify them.

// =============================================================================
// DRACULA
// Source: https://spec.draculatheme.com/
// =============================================================================

// DraculaPalette is the official Dracula theme.
var DraculaPalette = Palette{
	Name: "Dracula",

	// Background depth
	Bg0: lipgloss.Color("#282A36"), // Background
	Bg1: lipgloss.Color("#44475A"), // Selection
	Bg2: lipgloss.Color("#44475A"), // Current line
	Bg3: lipgloss.Color("#6272A4"), // Comment

	// Foreground depth
	Fg0: lipgloss.Color("#F8F8F2"), // Foreground
	Fg1: lipgloss.Color("#F8F8F2"), // Foreground
	Fg2: lipgloss.Color("#BFBFBF"), // Muted
	Fg3: lipgloss.Color("#6272A4"), // Comment

	// Base colors
	Gray: lipgloss.Color("#6272A4"), GrayLt: lipgloss.Color("#7E8AAF"), GrayDk: lipgloss.Color("#4A5568"),
	Red: lipgloss.Color("#FF5555"), RedLt: lipgloss.Color("#FF7777"), RedDk: lipgloss.Color("#CC4444"),
	Orange: lipgloss.Color("#FFB86C"), OrangeLt: lipgloss.Color("#FFCF9C"), OrangeDk: lipgloss.Color("#CC9356"),
	Green: lipgloss.Color("#50FA7B"), GreenLt: lipgloss.Color("#70FF95"), GreenDk: lipgloss.Color("#40C862"),
	Yellow: lipgloss.Color("#F1FA8C"), YellowLt: lipgloss.Color("#F5FCA8"), YellowDk: lipgloss.Color("#C1C870"),
	Blue: lipgloss.Color("#6272F4"), BlueLt: lipgloss.Color("#8290F6"), BlueDk: lipgloss.Color("#4E5BC3"),
	Purple: lipgloss.Color("#BD93F9"), PurpleLt: lipgloss.Color("#CDA8FA"), PurpleDk: lipgloss.Color("#8B6BC4"),
	Pink: lipgloss.Color("#FF79C6"), PinkLt: lipgloss.Color("#FF9BD4"), PinkDk: lipgloss.Color("#CC619E"),
	Aqua: lipgloss.Color("#8BE9FD"), AquaLt: lipgloss.Color("#A5EFFE"), AquaDk: lipgloss.Color("#6FBBCA"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#BD93F9"), HeroDk: lipgloss.Color("#8B6BC4"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#8BE9FD"), // Aqua
	Success: lipgloss.Color("#50FA7B"), // Green
	Warning: lipgloss.Color("#FFB86C"), // Orange
	Error:   lipgloss.Color("#FF5555"), // Red
	Info:    lipgloss.Color("#8BE9FD"), // Aqua

	Border: lipgloss.RoundedBorder(),
}

// DraculaLightPalette is the Alucard (light) variant of Dracula.
// Source: https://spec.draculatheme.com/
var DraculaLightPalette = Palette{
	Name: "DraculaLight",

	// Background depth
	Bg0: lipgloss.Color("#FFFBEB"), // Background
	Bg1: lipgloss.Color("#CFCFDE"), // Selection
	Bg2: lipgloss.Color("#E8E8F0"), // Elevated
	Bg3: lipgloss.Color("#6C664B"), // Comment

	// Foreground depth
	Fg0: lipgloss.Color("#1F1F1F"), // Foreground
	Fg1: lipgloss.Color("#1F1F1F"), // Foreground
	Fg2: lipgloss.Color("#4A4A4A"), // Muted
	Fg3: lipgloss.Color("#6C664B"), // Comment

	// Base colors
	Gray: lipgloss.Color("#6C664B"), GrayLt: lipgloss.Color("#8A8470"), GrayDk: lipgloss.Color("#4E4A38"),
	Red: lipgloss.Color("#CB3A2A"), RedLt: lipgloss.Color("#E05545"), RedDk: lipgloss.Color("#A22E22"),
	Orange: lipgloss.Color("#A34D14"), OrangeLt: lipgloss.Color("#C46018"), OrangeDk: lipgloss.Color("#823E10"),
	Green: lipgloss.Color("#14710A"), GreenLt: lipgloss.Color("#1D8F10"), GreenDk: lipgloss.Color("#0F5F06"),
	Yellow: lipgloss.Color("#D4A017"), YellowLt: lipgloss.Color("#E8B832"), YellowDk: lipgloss.Color("#A98012"),
	Blue: lipgloss.Color("#0460A8"), BlueLt: lipgloss.Color("#0578CC"), BlueDk: lipgloss.Color("#034A84"),
	Purple: lipgloss.Color("#644AC9"), PurpleLt: lipgloss.Color("#8B6FD9"), PurpleDk: lipgloss.Color("#503BA1"),
	Pink: lipgloss.Color("#A3144D"), PinkLt: lipgloss.Color("#C42A66"), PinkDk: lipgloss.Color("#82103E"),
	Aqua: lipgloss.Color("#036A96"), AquaLt: lipgloss.Color("#0488BE"), AquaDk: lipgloss.Color("#025578"),

	// Hero colors (tabs/borders cascade from these - light theme uses Lt for muted)
	Hero: lipgloss.Color("#644AC9"), HeroDk: lipgloss.Color("#8B6FD9"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#036A96"), // Aqua
	Success: lipgloss.Color("#14710A"), // Green
	Warning: lipgloss.Color("#A34D14"), // Orange
	Error:   lipgloss.Color("#CB3A2A"), // Red
	Info:    lipgloss.Color("#036A96"), // Aqua

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// CATPPUCCIN
// Source: https://catppuccin.com/palette
// =============================================================================

// CatppuccinMochaPalette is the darkest Catppuccin flavor.
var CatppuccinMochaPalette = Palette{
	Name: "CatMocha",

	// Background depth
	Bg0: lipgloss.Color("#1e1e2e"), // Base
	Bg1: lipgloss.Color("#181825"), // Mantle
	Bg2: lipgloss.Color("#313244"), // Surface0
	Bg3: lipgloss.Color("#45475a"), // Surface1

	// Foreground depth
	Fg0: lipgloss.Color("#cdd6f4"), // Text
	Fg1: lipgloss.Color("#bac2de"), // Subtext1
	Fg2: lipgloss.Color("#a6adc8"), // Subtext0
	Fg3: lipgloss.Color("#6c7086"), // Overlay0

	// Base colors
	Gray: lipgloss.Color("#6c7086"), GrayLt: lipgloss.Color("#9399b2"), GrayDk: lipgloss.Color("#585b70"),
	Red: lipgloss.Color("#f38ba8"), RedLt: lipgloss.Color("#f5a6ba"), RedDk: lipgloss.Color("#c27088"),
	Orange: lipgloss.Color("#fab387"), OrangeLt: lipgloss.Color("#fcc5a5"), OrangeDk: lipgloss.Color("#c88f6c"),
	Green: lipgloss.Color("#a6e3a1"), GreenLt: lipgloss.Color("#b8eab4"), GreenDk: lipgloss.Color("#85b681"),
	Yellow: lipgloss.Color("#f9e2af"), YellowLt: lipgloss.Color("#faeac4"), YellowDk: lipgloss.Color("#c7b58c"),
	Blue: lipgloss.Color("#89b4fa"), BlueLt: lipgloss.Color("#a4c5fb"), BlueDk: lipgloss.Color("#6e90c8"),
	Purple: lipgloss.Color("#cba6f7"), PurpleLt: lipgloss.Color("#d8bbf9"), PurpleDk: lipgloss.Color("#9B7AC4"),
	Pink: lipgloss.Color("#f5c2e7"), PinkLt: lipgloss.Color("#f8d4ed"), PinkDk: lipgloss.Color("#c49bb9"),
	Aqua: lipgloss.Color("#94e2d5"), AquaLt: lipgloss.Color("#aae9df"), AquaDk: lipgloss.Color("#76b5aa"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#cba6f7"), HeroDk: lipgloss.Color("#9B7AC4"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#313244"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#89b4fa"), // Blue
	Success: lipgloss.Color("#a6e3a1"), // Green
	Warning: lipgloss.Color("#fab387"), // Orange
	Error:   lipgloss.Color("#f38ba8"), // Red
	Info:    lipgloss.Color("#89dceb"), // Sky

	Border: lipgloss.RoundedBorder(),
}

// CatppuccinMacchiatoPalette is the medium-dark Catppuccin flavor.
var CatppuccinMacchiatoPalette = Palette{
	Name: "CatMacchiato",

	// Background depth
	Bg0: lipgloss.Color("#24273a"), // Base
	Bg1: lipgloss.Color("#1e2030"), // Mantle
	Bg2: lipgloss.Color("#363a4f"), // Surface0
	Bg3: lipgloss.Color("#494d64"), // Surface1

	// Foreground depth
	Fg0: lipgloss.Color("#cad3f5"), // Text
	Fg1: lipgloss.Color("#b8c0e0"), // Subtext1
	Fg2: lipgloss.Color("#a5adcb"), // Subtext0
	Fg3: lipgloss.Color("#6e738d"), // Overlay0

	// Base colors
	Gray: lipgloss.Color("#6e738d"), GrayLt: lipgloss.Color("#939ab7"), GrayDk: lipgloss.Color("#5b5f77"),
	Red: lipgloss.Color("#ed8796"), RedLt: lipgloss.Color("#f1a1ad"), RedDk: lipgloss.Color("#be6c78"),
	Orange: lipgloss.Color("#f5a97f"), OrangeLt: lipgloss.Color("#f7bb9c"), OrangeDk: lipgloss.Color("#c48766"),
	Green: lipgloss.Color("#a6da95"), GreenLt: lipgloss.Color("#b9e3ac"), GreenDk: lipgloss.Color("#85ae77"),
	Yellow: lipgloss.Color("#eed49f"), YellowLt: lipgloss.Color("#f2deb5"), YellowDk: lipgloss.Color("#beaa7f"),
	Blue: lipgloss.Color("#8aadf4"), BlueLt: lipgloss.Color("#a5c0f6"), BlueDk: lipgloss.Color("#6e8ac3"),
	Purple: lipgloss.Color("#c6a0f6"), PurpleLt: lipgloss.Color("#d4b5f8"), PurpleDk: lipgloss.Color("#9676C3"),
	Pink: lipgloss.Color("#f5bde6"), PinkLt: lipgloss.Color("#f8d0ed"), PinkDk: lipgloss.Color("#c497b8"),
	Aqua: lipgloss.Color("#8bd5ca"), AquaLt: lipgloss.Color("#a5dfd6"), AquaDk: lipgloss.Color("#6faaa2"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#c6a0f6"), HeroDk: lipgloss.Color("#9676C3"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#363a4f"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#8aadf4"), // Blue
	Success: lipgloss.Color("#a6da95"), // Green
	Warning: lipgloss.Color("#f5a97f"), // Orange
	Error:   lipgloss.Color("#ed8796"), // Red
	Info:    lipgloss.Color("#91d7e3"), // Sky

	Border: lipgloss.RoundedBorder(),
}

// CatppuccinFrappePalette is the muted dark Catppuccin flavor.
var CatppuccinFrappePalette = Palette{
	Name: "CatFrappe",

	// Background depth
	Bg0: lipgloss.Color("#303446"), // Base
	Bg1: lipgloss.Color("#292c3c"), // Mantle
	Bg2: lipgloss.Color("#414559"), // Surface0
	Bg3: lipgloss.Color("#51576d"), // Surface1

	// Foreground depth
	Fg0: lipgloss.Color("#c6d0f5"), // Text
	Fg1: lipgloss.Color("#b5bfe2"), // Subtext1
	Fg2: lipgloss.Color("#a5adce"), // Subtext0
	Fg3: lipgloss.Color("#737994"), // Overlay0

	// Base colors
	Gray: lipgloss.Color("#737994"), GrayLt: lipgloss.Color("#949cbb"), GrayDk: lipgloss.Color("#5c6179"),
	Red: lipgloss.Color("#e78284"), RedLt: lipgloss.Color("#ec9c9d"), RedDk: lipgloss.Color("#b9686a"),
	Orange: lipgloss.Color("#ef9f76"), OrangeLt: lipgloss.Color("#f3b594"), OrangeDk: lipgloss.Color("#bf7f5e"),
	Green: lipgloss.Color("#a6d189"), GreenLt: lipgloss.Color("#b9dca1"), GreenDk: lipgloss.Color("#85a76e"),
	Yellow: lipgloss.Color("#e5c890"), YellowLt: lipgloss.Color("#ebd5a8"), YellowDk: lipgloss.Color("#b7a073"),
	Blue: lipgloss.Color("#8caaee"), BlueLt: lipgloss.Color("#a6bcf1"), BlueDk: lipgloss.Color("#7088be"),
	Purple: lipgloss.Color("#ca9ee6"), PurpleLt: lipgloss.Color("#d7b4eb"), PurpleDk: lipgloss.Color("#9A74B6"),
	Pink: lipgloss.Color("#f4b8e4"), PinkLt: lipgloss.Color("#f7ccea"), PinkDk: lipgloss.Color("#c393b6"),
	Aqua: lipgloss.Color("#81c8be"), AquaLt: lipgloss.Color("#9bd5cd"), AquaDk: lipgloss.Color("#67a098"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#ca9ee6"), HeroDk: lipgloss.Color("#9A74B6"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#8caaee"), // Blue
	Success: lipgloss.Color("#a6d189"), // Green
	Warning: lipgloss.Color("#ef9f76"), // Orange
	Error:   lipgloss.Color("#e78284"), // Red
	Info:    lipgloss.Color("#99d1db"), // Sky

	Border: lipgloss.RoundedBorder(),
}

// CatppuccinLattePalette is the light Catppuccin flavor.
var CatppuccinLattePalette = Palette{
	Name: "CatLatte",

	// Background depth
	Bg0: lipgloss.Color("#eff1f5"), // Base
	Bg1: lipgloss.Color("#e6e9ef"), // Mantle
	Bg2: lipgloss.Color("#ccd0da"), // Surface0
	Bg3: lipgloss.Color("#bcc0cc"), // Surface1

	// Foreground depth
	Fg0: lipgloss.Color("#4c4f69"), // Text
	Fg1: lipgloss.Color("#5c5f77"), // Subtext1
	Fg2: lipgloss.Color("#6c6f85"), // Subtext0
	Fg3: lipgloss.Color("#9ca0b0"), // Overlay0

	// Base colors
	Gray: lipgloss.Color("#9ca0b0"), GrayLt: lipgloss.Color("#bcc0cc"), GrayDk: lipgloss.Color("#7c8091"),
	Red: lipgloss.Color("#d20f39"), RedLt: lipgloss.Color("#e53e5f"), RedDk: lipgloss.Color("#a80c2e"),
	Orange: lipgloss.Color("#fe640b"), OrangeLt: lipgloss.Color("#fe8541"), OrangeDk: lipgloss.Color("#cb5009"),
	Green: lipgloss.Color("#40a02b"), GreenLt: lipgloss.Color("#56c23b"), GreenDk: lipgloss.Color("#338022"),
	Yellow: lipgloss.Color("#df8e1d"), YellowLt: lipgloss.Color("#e9a63e"), YellowDk: lipgloss.Color("#b27217"),
	Blue: lipgloss.Color("#1e66f5"), BlueLt: lipgloss.Color("#4a84f7"), BlueDk: lipgloss.Color("#1852c4"),
	Purple: lipgloss.Color("#8839ef"), PurpleLt: lipgloss.Color("#A563F7"), PurpleDk: lipgloss.Color("#6d2ebf"),
	Pink: lipgloss.Color("#ea76cb"), PinkLt: lipgloss.Color("#ee92d6"), PinkDk: lipgloss.Color("#bb5ea2"),
	Aqua: lipgloss.Color("#179299"), AquaLt: lipgloss.Color("#1db5be"), AquaDk: lipgloss.Color("#12757a"),

	// Hero colors (tabs/borders cascade from these - light theme uses Lt for muted)
	Hero: lipgloss.Color("#8839ef"), HeroDk: lipgloss.Color("#A563F7"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#1e66f5"), // Blue
	Success: lipgloss.Color("#40a02b"), // Green
	Warning: lipgloss.Color("#fe640b"), // Orange
	Error:   lipgloss.Color("#d20f39"), // Red
	Info:    lipgloss.Color("#04a5e5"), // Sky

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// TOKYO NIGHT
// Source: https://github.com/enkia/tokyo-night-vscode-theme
// =============================================================================

// TokyoNightPalette is the darkest Tokyo Night variant.
var TokyoNightPalette = Palette{
	Name: "TokyoNight",

	// Background depth
	Bg0: lipgloss.Color("#1a1b26"), // bg
	Bg1: lipgloss.Color("#16161e"), // bg_dark
	Bg2: lipgloss.Color("#292e42"), // bg_highlight
	Bg3: lipgloss.Color("#3b4261"), // fg_gutter

	// Foreground depth
	Fg0: lipgloss.Color("#c0caf5"), // fg
	Fg1: lipgloss.Color("#a9b1d6"), // fg_dark
	Fg2: lipgloss.Color("#737aa2"), // dark3
	Fg3: lipgloss.Color("#565f89"), // comment

	// Base colors
	Gray: lipgloss.Color("#565f89"), GrayLt: lipgloss.Color("#737aa2"), GrayDk: lipgloss.Color("#414868"),
	Red: lipgloss.Color("#f7768e"), RedLt: lipgloss.Color("#f9919f"), RedDk: lipgloss.Color("#c55e72"),
	Orange: lipgloss.Color("#ff9e64"), OrangeLt: lipgloss.Color("#ffb386"), OrangeDk: lipgloss.Color("#cc7e50"),
	Green: lipgloss.Color("#9ece6a"), GreenLt: lipgloss.Color("#b3d987"), GreenDk: lipgloss.Color("#7ea555"),
	Yellow: lipgloss.Color("#e0af68"), YellowLt: lipgloss.Color("#e8c388"), YellowDk: lipgloss.Color("#b38c53"),
	Blue: lipgloss.Color("#7aa2f7"), BlueLt: lipgloss.Color("#97b6f9"), BlueDk: lipgloss.Color("#5578C4"),
	Purple: lipgloss.Color("#bb9af7"), PurpleLt: lipgloss.Color("#cbb0f9"), PurpleDk: lipgloss.Color("#957bc6"),
	Pink: lipgloss.Color("#ff007c"), PinkLt: lipgloss.Color("#ff4d9e"), PinkDk: lipgloss.Color("#cc0063"),
	Aqua: lipgloss.Color("#7dcfff"), AquaLt: lipgloss.Color("#9adbff"), AquaDk: lipgloss.Color("#64a6cc"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#7aa2f7"), HeroDk: lipgloss.Color("#5578C4"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#292e42"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#7dcfff"), // Aqua
	Success: lipgloss.Color("#9ece6a"), // Green
	Warning: lipgloss.Color("#ff9e64"), // Orange
	Error:   lipgloss.Color("#f7768e"), // Red
	Info:    lipgloss.Color("#2ac3de"), // Cyan

	Border: lipgloss.RoundedBorder(),
}

// TokyoStormPalette is the storm variant of Tokyo Night.
var TokyoStormPalette = Palette{
	Name: "TokyoStorm",

	// Background depth
	Bg0: lipgloss.Color("#24283b"), // bg
	Bg1: lipgloss.Color("#1f2335"), // bg_dark
	Bg2: lipgloss.Color("#292e42"), // bg_highlight
	Bg3: lipgloss.Color("#3b4261"), // fg_gutter

	// Foreground depth
	Fg0: lipgloss.Color("#c0caf5"), // fg
	Fg1: lipgloss.Color("#a9b1d6"), // fg_dark
	Fg2: lipgloss.Color("#737aa2"), // dark3
	Fg3: lipgloss.Color("#565f89"), // comment

	// Base colors
	Gray: lipgloss.Color("#565f89"), GrayLt: lipgloss.Color("#737aa2"), GrayDk: lipgloss.Color("#414868"),
	Red: lipgloss.Color("#f7768e"), RedLt: lipgloss.Color("#f9919f"), RedDk: lipgloss.Color("#c55e72"),
	Orange: lipgloss.Color("#ff9e64"), OrangeLt: lipgloss.Color("#ffb386"), OrangeDk: lipgloss.Color("#cc7e50"),
	Green: lipgloss.Color("#9ece6a"), GreenLt: lipgloss.Color("#b3d987"), GreenDk: lipgloss.Color("#7ea555"),
	Yellow: lipgloss.Color("#e0af68"), YellowLt: lipgloss.Color("#e8c388"), YellowDk: lipgloss.Color("#b38c53"),
	Blue: lipgloss.Color("#7aa2f7"), BlueLt: lipgloss.Color("#97b6f9"), BlueDk: lipgloss.Color("#5578C4"),
	Purple: lipgloss.Color("#bb9af7"), PurpleLt: lipgloss.Color("#cbb0f9"), PurpleDk: lipgloss.Color("#957bc6"),
	Pink: lipgloss.Color("#ff007c"), PinkLt: lipgloss.Color("#ff4d9e"), PinkDk: lipgloss.Color("#cc0063"),
	Aqua: lipgloss.Color("#7dcfff"), AquaLt: lipgloss.Color("#9adbff"), AquaDk: lipgloss.Color("#64a6cc"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#7aa2f7"), HeroDk: lipgloss.Color("#5578C4"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#292e42"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#7dcfff"), // Aqua
	Success: lipgloss.Color("#9ece6a"), // Green
	Warning: lipgloss.Color("#ff9e64"), // Orange
	Error:   lipgloss.Color("#f7768e"), // Red
	Info:    lipgloss.Color("#2ac3de"), // Cyan

	Border: lipgloss.RoundedBorder(),
}

// TokyoMoonPalette is the moon variant of Tokyo Night.
var TokyoMoonPalette = Palette{
	Name: "TokyoMoon",

	// Background depth
	Bg0: lipgloss.Color("#222436"), // bg
	Bg1: lipgloss.Color("#1e2030"), // bg_dark
	Bg2: lipgloss.Color("#2f334d"), // bg_highlight
	Bg3: lipgloss.Color("#3b4261"), // fg_gutter

	// Foreground depth
	Fg0: lipgloss.Color("#c8d3f5"), // fg
	Fg1: lipgloss.Color("#828bb8"), // fg_dark
	Fg2: lipgloss.Color("#828bb8"), // muted
	Fg3: lipgloss.Color("#636da6"), // comment

	// Base colors
	Gray: lipgloss.Color("#636da6"), GrayLt: lipgloss.Color("#828bb8"), GrayDk: lipgloss.Color("#4e5685"),
	Red: lipgloss.Color("#ff757f"), RedLt: lipgloss.Color("#ff9199"), RedDk: lipgloss.Color("#cc5e66"),
	Orange: lipgloss.Color("#ff966c"), OrangeLt: lipgloss.Color("#ffae8c"), OrangeDk: lipgloss.Color("#cc7856"),
	Green: lipgloss.Color("#c3e88d"), GreenLt: lipgloss.Color("#d1eda6"), GreenDk: lipgloss.Color("#9cba71"),
	Yellow: lipgloss.Color("#ffc777"), YellowLt: lipgloss.Color("#ffd599"), YellowDk: lipgloss.Color("#cc9f5f"),
	Blue: lipgloss.Color("#82aaff"), BlueLt: lipgloss.Color("#9ebbff"), BlueDk: lipgloss.Color("#5C7DCC"),
	Purple: lipgloss.Color("#c099ff"), PurpleLt: lipgloss.Color("#d0b0ff"), PurpleDk: lipgloss.Color("#9a7acc"),
	Pink: lipgloss.Color("#ff007c"), PinkLt: lipgloss.Color("#ff4d9e"), PinkDk: lipgloss.Color("#cc0063"),
	Aqua: lipgloss.Color("#86e1fc"), AquaLt: lipgloss.Color("#a3e9fd"), AquaDk: lipgloss.Color("#6bb4ca"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#82aaff"), HeroDk: lipgloss.Color("#5C7DCC"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#2f334d"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#86e1fc"), // Aqua
	Success: lipgloss.Color("#c3e88d"), // Green
	Warning: lipgloss.Color("#ff966c"), // Orange
	Error:   lipgloss.Color("#ff757f"), // Red
	Info:    lipgloss.Color("#65bcff"), // Cyan

	Border: lipgloss.RoundedBorder(),
}

// TokyoDayPalette is the light variant of Tokyo Night.
var TokyoDayPalette = Palette{
	Name: "TokyoDay",

	// Background depth
	Bg0: lipgloss.Color("#e1e2e7"), // bg
	Bg1: lipgloss.Color("#d0d5e3"), // bg_dark
	Bg2: lipgloss.Color("#c4c8da"), // bg_highlight
	Bg3: lipgloss.Color("#b7c1e3"), // visual

	// Foreground depth
	Fg0: lipgloss.Color("#3760bf"), // fg
	Fg1: lipgloss.Color("#3760bf"), // fg
	Fg2: lipgloss.Color("#6172b0"), // muted
	Fg3: lipgloss.Color("#848cb5"), // comment

	// Base colors
	Gray: lipgloss.Color("#848cb5"), GrayLt: lipgloss.Color("#9ba2c4"), GrayDk: lipgloss.Color("#6a7091"),
	Red: lipgloss.Color("#f52a65"), RedLt: lipgloss.Color("#f75584"), RedDk: lipgloss.Color("#c42251"),
	Orange: lipgloss.Color("#b15c00"), OrangeLt: lipgloss.Color("#d17800"), OrangeDk: lipgloss.Color("#8d4900"),
	Green: lipgloss.Color("#587539"), GreenLt: lipgloss.Color("#6d9147"), GreenDk: lipgloss.Color("#465e2e"),
	Yellow: lipgloss.Color("#8c6c3e"), YellowLt: lipgloss.Color("#a98450"), YellowDk: lipgloss.Color("#705632"),
	Blue: lipgloss.Color("#2e7de9"), BlueLt: lipgloss.Color("#5A9BF0"), BlueDk: lipgloss.Color("#2564ba"),
	Purple: lipgloss.Color("#9854f1"), PurpleLt: lipgloss.Color("#af79f4"), PurpleDk: lipgloss.Color("#7a43c1"),
	Pink: lipgloss.Color("#d23d8a"), PinkLt: lipgloss.Color("#e065a5"), PinkDk: lipgloss.Color("#a8316e"),
	Aqua: lipgloss.Color("#007197"), AquaLt: lipgloss.Color("#008db8"), AquaDk: lipgloss.Color("#005a79"),

	// Hero colors (tabs/borders cascade from these - light theme uses Lt for muted)
	Hero: lipgloss.Color("#2e7de9"), HeroDk: lipgloss.Color("#5A9BF0"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#007197"), // Aqua
	Success: lipgloss.Color("#587539"), // Green
	Warning: lipgloss.Color("#b15c00"), // Orange
	Error:   lipgloss.Color("#f52a65"), // Red
	Info:    lipgloss.Color("#007197"), // Aqua

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// ROSE PINE
// Source: https://rosepinetheme.com/palette
// =============================================================================

// RosePinePalette is the main dark Rose Pine variant.
var RosePinePalette = Palette{
	Name: "RosePine",

	// Background depth
	Bg0: lipgloss.Color("#191724"), // Base
	Bg1: lipgloss.Color("#1f1d2e"), // Surface
	Bg2: lipgloss.Color("#26233a"), // Overlay
	Bg3: lipgloss.Color("#6e6a86"), // Muted

	// Foreground depth
	Fg0: lipgloss.Color("#e0def4"), // Text
	Fg1: lipgloss.Color("#e0def4"), // Text
	Fg2: lipgloss.Color("#908caa"), // Subtle
	Fg3: lipgloss.Color("#6e6a86"), // Muted

	// Base colors
	Gray: lipgloss.Color("#6e6a86"), GrayLt: lipgloss.Color("#908caa"), GrayDk: lipgloss.Color("#56536b"),
	Red: lipgloss.Color("#eb6f92"), RedLt: lipgloss.Color("#f08ca8"), RedDk: lipgloss.Color("#bc5975"),
	Orange: lipgloss.Color("#ebbcba"), OrangeLt: lipgloss.Color("#f0ccc9"), OrangeDk: lipgloss.Color("#bc9695"),
	Green: lipgloss.Color("#9ccfd8"), GreenLt: lipgloss.Color("#b3dbe2"), GreenDk: lipgloss.Color("#7da6ad"),
	Yellow: lipgloss.Color("#f6c177"), YellowLt: lipgloss.Color("#f8cf95"), YellowDk: lipgloss.Color("#c59a5f"),
	Blue: lipgloss.Color("#31748f"), BlueLt: lipgloss.Color("#4590ad"), BlueDk: lipgloss.Color("#275d72"),
	Purple: lipgloss.Color("#c4a7e7"), PurpleLt: lipgloss.Color("#d3bbec"), PurpleDk: lipgloss.Color("#947DB7"),
	Pink: lipgloss.Color("#ebbcba"), PinkLt: lipgloss.Color("#f0ccc9"), PinkDk: lipgloss.Color("#bc9695"),
	Aqua: lipgloss.Color("#9ccfd8"), AquaLt: lipgloss.Color("#b3dbe2"), AquaDk: lipgloss.Color("#7da6ad"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#c4a7e7"), HeroDk: lipgloss.Color("#947DB7"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#26233a"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#31748f"), // Pine (Blue)
	Success: lipgloss.Color("#9ccfd8"), // Foam (Aqua)
	Warning: lipgloss.Color("#f6c177"), // Gold (Yellow)
	Error:   lipgloss.Color("#eb6f92"), // Love (Red)
	Info:    lipgloss.Color("#9ccfd8"), // Foam (Aqua)

	Border: lipgloss.RoundedBorder(),
}

// RosePineMoonPalette is the moon variant of Rose Pine.
var RosePineMoonPalette = Palette{
	Name: "RosePineMoon",

	// Background depth
	Bg0: lipgloss.Color("#232136"), // Base
	Bg1: lipgloss.Color("#2a273f"), // Surface
	Bg2: lipgloss.Color("#393552"), // Overlay
	Bg3: lipgloss.Color("#6e6a86"), // Muted

	// Foreground depth
	Fg0: lipgloss.Color("#e0def4"), // Text
	Fg1: lipgloss.Color("#e0def4"), // Text
	Fg2: lipgloss.Color("#908caa"), // Subtle
	Fg3: lipgloss.Color("#6e6a86"), // Muted

	// Base colors
	Gray: lipgloss.Color("#6e6a86"), GrayLt: lipgloss.Color("#908caa"), GrayDk: lipgloss.Color("#56536b"),
	Red: lipgloss.Color("#eb6f92"), RedLt: lipgloss.Color("#f08ca8"), RedDk: lipgloss.Color("#bc5975"),
	Orange: lipgloss.Color("#ea9a97"), OrangeLt: lipgloss.Color("#efb2af"), OrangeDk: lipgloss.Color("#bb7b79"),
	Green: lipgloss.Color("#9ccfd8"), GreenLt: lipgloss.Color("#b3dbe2"), GreenDk: lipgloss.Color("#7da6ad"),
	Yellow: lipgloss.Color("#f6c177"), YellowLt: lipgloss.Color("#f8cf95"), YellowDk: lipgloss.Color("#c59a5f"),
	Blue: lipgloss.Color("#3e8fb0"), BlueLt: lipgloss.Color("#5aa8c7"), BlueDk: lipgloss.Color("#32728d"),
	Purple: lipgloss.Color("#c4a7e7"), PurpleLt: lipgloss.Color("#d3bbec"), PurpleDk: lipgloss.Color("#947DB7"),
	Pink: lipgloss.Color("#ea9a97"), PinkLt: lipgloss.Color("#efb2af"), PinkDk: lipgloss.Color("#bb7b79"),
	Aqua: lipgloss.Color("#9ccfd8"), AquaLt: lipgloss.Color("#b3dbe2"), AquaDk: lipgloss.Color("#7da6ad"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#c4a7e7"), HeroDk: lipgloss.Color("#947DB7"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#3e8fb0"), // Pine (Blue)
	Success: lipgloss.Color("#9ccfd8"), // Foam (Aqua)
	Warning: lipgloss.Color("#f6c177"), // Gold (Yellow)
	Error:   lipgloss.Color("#eb6f92"), // Love (Red)
	Info:    lipgloss.Color("#9ccfd8"), // Foam (Aqua)

	Border: lipgloss.RoundedBorder(),
}

// RosePineDawnPalette is the light variant of Rose Pine.
var RosePineDawnPalette = Palette{
	Name: "RosePineDawn",

	// Background depth
	Bg0: lipgloss.Color("#faf4ed"), // Base
	Bg1: lipgloss.Color("#fffaf3"), // Surface
	Bg2: lipgloss.Color("#f2e9e1"), // Overlay
	Bg3: lipgloss.Color("#9893a5"), // Muted

	// Foreground depth
	Fg0: lipgloss.Color("#575279"), // Text
	Fg1: lipgloss.Color("#575279"), // Text
	Fg2: lipgloss.Color("#797593"), // Subtle
	Fg3: lipgloss.Color("#9893a5"), // Muted

	// Base colors
	Gray: lipgloss.Color("#9893a5"), GrayLt: lipgloss.Color("#b0acba"), GrayDk: lipgloss.Color("#7a7684"),
	Red: lipgloss.Color("#b4637a"), RedLt: lipgloss.Color("#c77e92"), RedDk: lipgloss.Color("#904f62"),
	Orange: lipgloss.Color("#d7827e"), OrangeLt: lipgloss.Color("#e09c98"), OrangeDk: lipgloss.Color("#ac6865"),
	Green: lipgloss.Color("#56949f"), GreenLt: lipgloss.Color("#71aab4"), GreenDk: lipgloss.Color("#45767f"),
	Yellow: lipgloss.Color("#ea9d34"), YellowLt: lipgloss.Color("#eeb25c"), YellowDk: lipgloss.Color("#bb7e2a"),
	Blue: lipgloss.Color("#286983"), BlueLt: lipgloss.Color("#3a849f"), BlueDk: lipgloss.Color("#205469"),
	Purple: lipgloss.Color("#907aa9"), PurpleLt: lipgloss.Color("#B09DC4"), PurpleDk: lipgloss.Color("#736287"),
	Pink: lipgloss.Color("#d7827e"), PinkLt: lipgloss.Color("#e09c98"), PinkDk: lipgloss.Color("#ac6865"),
	Aqua: lipgloss.Color("#56949f"), AquaLt: lipgloss.Color("#71aab4"), AquaDk: lipgloss.Color("#45767f"),

	// Hero colors (tabs/borders cascade from these - light theme uses Lt for muted)
	Hero: lipgloss.Color("#907aa9"), HeroDk: lipgloss.Color("#B09DC4"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#f2e9e1"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#286983"), // Pine (Blue)
	Success: lipgloss.Color("#56949f"), // Foam (Aqua)
	Warning: lipgloss.Color("#ea9d34"), // Gold (Yellow)
	Error:   lipgloss.Color("#b4637a"), // Love (Red)
	Info:    lipgloss.Color("#56949f"), // Foam (Aqua)

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// CLASSIC THEMES
// =============================================================================

// RetroPalette is a classic green-only terminal theme (no other colors).
// Inspired by early CRT monitors.
var RetroPalette = Palette{
	Name: "Retro",

	// Background depth (green-tinted blacks)
	Bg0: lipgloss.Color("#000000"),
	Bg1: lipgloss.Color("#001100"),
	Bg2: lipgloss.Color("#002200"),
	Bg3: lipgloss.Color("#003300"),

	// Foreground depth (green variants)
	Fg0: lipgloss.Color("#00ff00"),
	Fg1: lipgloss.Color("#00dd00"),
	Fg2: lipgloss.Color("#00aa00"),
	Fg3: lipgloss.Color("#006600"),

	// Base colors (all map to 3 green tiers: Lt=Hero, Normal=HeroDk, Dk=Unselected)
	Gray: lipgloss.Color("#00bb00"), GrayLt: lipgloss.Color("#00ff00"), GrayDk: lipgloss.Color("#008800"),
	Red: lipgloss.Color("#00bb00"), RedLt: lipgloss.Color("#00ff00"), RedDk: lipgloss.Color("#008800"),
	Orange: lipgloss.Color("#00bb00"), OrangeLt: lipgloss.Color("#00ff00"), OrangeDk: lipgloss.Color("#008800"),
	Green: lipgloss.Color("#00bb00"), GreenLt: lipgloss.Color("#00ff00"), GreenDk: lipgloss.Color("#008800"),
	Yellow: lipgloss.Color("#00bb00"), YellowLt: lipgloss.Color("#00ff00"), YellowDk: lipgloss.Color("#008800"),
	Blue: lipgloss.Color("#00bb00"), BlueLt: lipgloss.Color("#00ff00"), BlueDk: lipgloss.Color("#008800"),
	Purple: lipgloss.Color("#00bb00"), PurpleLt: lipgloss.Color("#00ff00"), PurpleDk: lipgloss.Color("#008800"),
	Pink: lipgloss.Color("#00bb00"), PinkLt: lipgloss.Color("#00ff00"), PinkDk: lipgloss.Color("#008800"),
	Aqua: lipgloss.Color("#00bb00"), AquaLt: lipgloss.Color("#00ff00"), AquaDk: lipgloss.Color("#008800"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#00ff00"), HeroDk: lipgloss.Color("#00bb00"),

	// Semantic colors (all green for retro look)
	Accent:  lipgloss.Color("#00bb00"),
	Success: lipgloss.Color("#00ff00"),
	Warning: lipgloss.Color("#00bb00"),
	Error:   lipgloss.Color("#00ff00"),
	Info:    lipgloss.Color("#00bb00"),

	// Notification background (filled for retro CRT look)
	NotifBg: lipgloss.Color("#001100"), // Bg1 - green-tinted fill

	Border:        lipgloss.DoubleBorder(),
	IconInfo:      "[i]",
	IconNote:      "[*]",
	IconTip:       "[?]",
	IconImportant: "[!]",
	IconSuccess:   "[OK]",
	IconWarning:   "[!!]",
	IconCaution:   "[!!!]",
	IconError:     "[ERR]",
}

// RetroMonoPalette is a classic amber/orange terminal theme.
// Inspired by early amber phosphor CRT monitors.
var RetroMonoPalette = Palette{
	Name: "RetroMono",

	// Background depth (amber-tinted blacks)
	Bg0: lipgloss.Color("#000000"),
	Bg1: lipgloss.Color("#1a0a00"),
	Bg2: lipgloss.Color("#2a1500"),
	Bg3: lipgloss.Color("#3a2000"),

	// Foreground depth (amber/orange variants)
	Fg0: lipgloss.Color("#ffaa00"),
	Fg1: lipgloss.Color("#dd9500"),
	Fg2: lipgloss.Color("#aa7000"),
	Fg3: lipgloss.Color("#775000"),

	// Base colors (all map to 3 amber tiers: Lt=Hero, Normal=HeroDk, Dk=Unselected)
	Gray: lipgloss.Color("#cc8800"), GrayLt: lipgloss.Color("#ffaa00"), GrayDk: lipgloss.Color("#996600"),
	Red: lipgloss.Color("#cc8800"), RedLt: lipgloss.Color("#ffaa00"), RedDk: lipgloss.Color("#996600"),
	Orange: lipgloss.Color("#cc8800"), OrangeLt: lipgloss.Color("#ffaa00"), OrangeDk: lipgloss.Color("#996600"),
	Green: lipgloss.Color("#cc8800"), GreenLt: lipgloss.Color("#ffaa00"), GreenDk: lipgloss.Color("#996600"),
	Yellow: lipgloss.Color("#cc8800"), YellowLt: lipgloss.Color("#ffaa00"), YellowDk: lipgloss.Color("#996600"),
	Blue: lipgloss.Color("#cc8800"), BlueLt: lipgloss.Color("#ffaa00"), BlueDk: lipgloss.Color("#996600"),
	Purple: lipgloss.Color("#cc8800"), PurpleLt: lipgloss.Color("#ffaa00"), PurpleDk: lipgloss.Color("#996600"),
	Pink: lipgloss.Color("#cc8800"), PinkLt: lipgloss.Color("#ffaa00"), PinkDk: lipgloss.Color("#996600"),
	Aqua: lipgloss.Color("#cc8800"), AquaLt: lipgloss.Color("#ffaa00"), AquaDk: lipgloss.Color("#996600"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#ffaa00"), HeroDk: lipgloss.Color("#cc8800"),

	// Semantic colors (all amber for retro look)
	Accent:  lipgloss.Color("#cc8800"),
	Success: lipgloss.Color("#ffaa00"),
	Warning: lipgloss.Color("#cc8800"),
	Error:   lipgloss.Color("#ffaa00"),
	Info:    lipgloss.Color("#cc8800"),

	// Notification background (filled for retro CRT look)
	NotifBg: lipgloss.Color("#1a0a00"), // Bg1 - amber-tinted fill

	Border:        lipgloss.DoubleBorder(),
	IconInfo:      "[i]",
	IconNote:      "[*]",
	IconTip:       "[?]",
	IconImportant: "[!]",
	IconSuccess:   "[OK]",
	IconWarning:   "[!!]",
	IconCaution:   "[!!!]",
	IconError:     "[ERR]",
}

// MinimalPalette is a clean, understated theme with hidden borders.
// Uses white/gray as the hero for a monochrome aesthetic.
var MinimalPalette = Palette{
	Name: "Minimal",

	// Background depth
	Bg0: lipgloss.Color("#1a1a1a"),
	Bg1: lipgloss.Color("#2a2a2a"),
	Bg2: lipgloss.Color("#3a3a3a"),
	Bg3: lipgloss.Color("#4a4a4a"),

	// Foreground depth
	Fg0: lipgloss.Color("#ffffff"),
	Fg1: lipgloss.Color("#cccccc"),
	Fg2: lipgloss.Color("#888888"),
	Fg3: lipgloss.Color("#666666"),

	// Base colors (grayscale with subtle tints)
	Gray: lipgloss.Color("#888888"), GrayLt: lipgloss.Color("#aaaaaa"), GrayDk: lipgloss.Color("#666666"),
	Red: lipgloss.Color("#ff8888"), RedLt: lipgloss.Color("#ffaaaa"), RedDk: lipgloss.Color("#cc6666"),
	Orange: lipgloss.Color("#ffaa88"), OrangeLt: lipgloss.Color("#ffccaa"), OrangeDk: lipgloss.Color("#cc8866"),
	Green: lipgloss.Color("#88ff88"), GreenLt: lipgloss.Color("#aaffaa"), GreenDk: lipgloss.Color("#66cc66"),
	Yellow: lipgloss.Color("#ffcc88"), YellowLt: lipgloss.Color("#ffddaa"), YellowDk: lipgloss.Color("#cca366"),
	Blue: lipgloss.Color("#88aaff"), BlueLt: lipgloss.Color("#aaccff"), BlueDk: lipgloss.Color("#6688cc"),
	Purple: lipgloss.Color("#ffffff"), PurpleLt: lipgloss.Color("#ffffff"), PurpleDk: lipgloss.Color("#aaaaaa"),
	Pink: lipgloss.Color("#cccccc"), PinkLt: lipgloss.Color("#dddddd"), PinkDk: lipgloss.Color("#aaaaaa"),
	Aqua: lipgloss.Color("#88dddd"), AquaLt: lipgloss.Color("#aaeeff"), AquaDk: lipgloss.Color("#66aaaa"),

	// Hero colors (tabs/borders cascade from these - white/gray for minimal)
	Hero: lipgloss.Color("#ffffff"), HeroDk: lipgloss.Color("#aaaaaa"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#3a3a3a"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#aaaaaa"), // Gray
	Success: lipgloss.Color("#88ff88"), // Green
	Warning: lipgloss.Color("#ffaa88"), // Orange
	Error:   lipgloss.Color("#ff8888"), // Red
	Info:    lipgloss.Color("#88aaff"), // Blue

	Border:        lipgloss.HiddenBorder(),
	IconInfo:      "ℹ",
	IconNote:      "•",
	IconTip:       "→",
	IconImportant: "!",
	IconSuccess:   "✓",
	IconWarning:   "⚠",
	IconCaution:   "⚠",
	IconError:     "✗",
}

// =============================================================================
// OTHER POPULAR THEMES
// =============================================================================

// SolarizedPalette is based on Ethan Schoonover's Solarized Dark palette.
// Source: https://ethanschoonover.com/solarized/
var SolarizedPalette = Palette{
	Name: "Solarized",

	// Background depth
	Bg0: lipgloss.Color("#002b36"), // base03
	Bg1: lipgloss.Color("#073642"), // base02
	Bg2: lipgloss.Color("#586e75"), // base01
	Bg3: lipgloss.Color("#657b83"), // base00

	// Foreground depth
	Fg0: lipgloss.Color("#fdf6e3"), // base3
	Fg1: lipgloss.Color("#eee8d5"), // base2
	Fg2: lipgloss.Color("#93a1a1"), // base1
	Fg3: lipgloss.Color("#839496"), // base0

	// Base colors
	Gray: lipgloss.Color("#839496"), GrayLt: lipgloss.Color("#93a1a1"), GrayDk: lipgloss.Color("#657b83"),
	Red: lipgloss.Color("#dc322f"), RedLt: lipgloss.Color("#e35d5b"), RedDk: lipgloss.Color("#b02826"),
	Orange: lipgloss.Color("#cb4b16"), OrangeLt: lipgloss.Color("#d96d42"), OrangeDk: lipgloss.Color("#a23c12"),
	Green: lipgloss.Color("#859900"), GreenLt: lipgloss.Color("#9db300"), GreenDk: lipgloss.Color("#6a7a00"),
	Yellow: lipgloss.Color("#b58900"), YellowLt: lipgloss.Color("#cda100"), YellowDk: lipgloss.Color("#916e00"),
	Blue: lipgloss.Color("#268bd2"), BlueLt: lipgloss.Color("#4aa3e0"), BlueDk: lipgloss.Color("#1a6399"),
	Purple: lipgloss.Color("#6c71c4"), PurpleLt: lipgloss.Color("#8b8fd3"), PurpleDk: lipgloss.Color("#565a9d"),
	Pink: lipgloss.Color("#d33682"), PinkLt: lipgloss.Color("#dc5a9a"), PinkDk: lipgloss.Color("#a92b68"),
	Aqua: lipgloss.Color("#2aa198"), AquaLt: lipgloss.Color("#45b9af"), AquaDk: lipgloss.Color("#22817a"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#268bd2"), HeroDk: lipgloss.Color("#1a6399"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#6c71c4"), // Violet
	Success: lipgloss.Color("#859900"), // Green
	Warning: lipgloss.Color("#cb4b16"), // Orange
	Error:   lipgloss.Color("#dc322f"), // Red
	Info:    lipgloss.Color("#268bd2"), // Blue

	Border: lipgloss.RoundedBorder(),
}

// NordPalette is inspired by Arctic/Nordic landscapes.
// Source: https://www.nordtheme.com/docs/colors-and-palettes
var NordPalette = Palette{
	Name: "Nord",

	// Background depth (Polar Night)
	Bg0: lipgloss.Color("#2e3440"), // nord0
	Bg1: lipgloss.Color("#3b4252"), // nord1
	Bg2: lipgloss.Color("#434c5e"), // nord2
	Bg3: lipgloss.Color("#4c566a"), // nord3

	// Foreground depth (Snow Storm)
	Fg0: lipgloss.Color("#eceff4"), // nord6
	Fg1: lipgloss.Color("#e5e9f0"), // nord5
	Fg2: lipgloss.Color("#d8dee9"), // nord4
	Fg3: lipgloss.Color("#8892a6"), // muted

	// Base colors
	Gray: lipgloss.Color("#4c566a"), GrayLt: lipgloss.Color("#616e88"), GrayDk: lipgloss.Color("#3b4252"),
	Red: lipgloss.Color("#bf616a"), RedLt: lipgloss.Color("#d08770"), RedDk: lipgloss.Color("#994d55"),
	Orange: lipgloss.Color("#d08770"), OrangeLt: lipgloss.Color("#dba18d"), OrangeDk: lipgloss.Color("#a66c5a"),
	Green: lipgloss.Color("#a3be8c"), GreenLt: lipgloss.Color("#b5cca0"), GreenDk: lipgloss.Color("#829870"),
	Yellow: lipgloss.Color("#ebcb8b"), YellowLt: lipgloss.Color("#f0d8a5"), YellowDk: lipgloss.Color("#bca26f"),
	Blue: lipgloss.Color("#81a1c1"), BlueLt: lipgloss.Color("#88c0d0"), BlueDk: lipgloss.Color("#5e7a94"),
	Purple: lipgloss.Color("#b48ead"), PurpleLt: lipgloss.Color("#c5a5bf"), PurpleDk: lipgloss.Color("#90718a"),
	Pink: lipgloss.Color("#b48ead"), PinkLt: lipgloss.Color("#c5a5bf"), PinkDk: lipgloss.Color("#90718a"),
	Aqua: lipgloss.Color("#8fbcbb"), AquaLt: lipgloss.Color("#a8cece"), AquaDk: lipgloss.Color("#729695"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#81a1c1"), HeroDk: lipgloss.Color("#5e7a94"),

	// Tab colors (override defaults)
	TabInactiveBg: lipgloss.Color("#434c5e"), // Bg2 instead of Bg1

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#5e81ac"), // nord10
	Success: lipgloss.Color("#a3be8c"), // Green
	Warning: lipgloss.Color("#d08770"), // Orange
	Error:   lipgloss.Color("#bf616a"), // Red
	Info:    lipgloss.Color("#81a1c1"), // Blue

	Border: lipgloss.RoundedBorder(),
}

// MonokaiPalette is based on Wimer Hazenberg's iconic Monokai palette.
// Source: https://monokai.pro/vscode
var MonokaiPalette = Palette{
	Name: "Monokai",

	// Background depth
	Bg0: lipgloss.Color("#272822"), // bg
	Bg1: lipgloss.Color("#3e3d32"), // dark gray
	Bg2: lipgloss.Color("#49483e"), // lighter gray
	Bg3: lipgloss.Color("#75715e"), // comment gray

	// Foreground depth
	Fg0: lipgloss.Color("#f8f8f2"), // white
	Fg1: lipgloss.Color("#f8f8f0"), // white variant
	Fg2: lipgloss.Color("#cfcfc2"), // muted
	Fg3: lipgloss.Color("#75715e"), // comment gray

	// Base colors
	Gray: lipgloss.Color("#75715e"), GrayLt: lipgloss.Color("#908b7a"), GrayDk: lipgloss.Color("#5e5a4b"),
	Red: lipgloss.Color("#f92672"), RedLt: lipgloss.Color("#fa518e"), RedDk: lipgloss.Color("#c41d5a"),
	Orange: lipgloss.Color("#fd971f"), OrangeLt: lipgloss.Color("#fdad4c"), OrangeDk: lipgloss.Color("#ca7919"),
	Green: lipgloss.Color("#a6e22e"), GreenLt: lipgloss.Color("#b8e858"), GreenDk: lipgloss.Color("#85b525"),
	Yellow: lipgloss.Color("#e6db74"), YellowLt: lipgloss.Color("#ece390"), YellowDk: lipgloss.Color("#b8af5d"),
	Blue: lipgloss.Color("#66d9ef"), BlueLt: lipgloss.Color("#88e2f3"), BlueDk: lipgloss.Color("#52aebf"),
	Purple: lipgloss.Color("#ae81ff"), PurpleLt: lipgloss.Color("#c19dff"), PurpleDk: lipgloss.Color("#8b67cc"),
	Pink: lipgloss.Color("#f92672"), PinkLt: lipgloss.Color("#fa518e"), PinkDk: lipgloss.Color("#c41d5a"),
	Aqua: lipgloss.Color("#66d9ef"), AquaLt: lipgloss.Color("#88e2f3"), AquaDk: lipgloss.Color("#52aebf"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#f92672"), HeroDk: lipgloss.Color("#c41d5a"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#ae81ff"), // Purple
	Success: lipgloss.Color("#a6e22e"), // Green
	Warning: lipgloss.Color("#fd971f"), // Orange
	Error:   lipgloss.Color("#f92672"), // Red/Pink
	Info:    lipgloss.Color("#66d9ef"), // Aqua

	Border: lipgloss.RoundedBorder(),
}

// GruvboxPalette is based on Pavel Pertsev's Gruvbox palette.
// Source: https://github.com/morhetz/gruvbox
var GruvboxPalette = Palette{
	Name: "Gruvbox",

	// Background depth
	Bg0: lipgloss.Color("#282828"), // bg
	Bg1: lipgloss.Color("#3c3836"), // bg1
	Bg2: lipgloss.Color("#504945"), // bg2
	Bg3: lipgloss.Color("#665c54"), // bg3

	// Foreground depth
	Fg0: lipgloss.Color("#fbf1c7"), // fg0
	Fg1: lipgloss.Color("#ebdbb2"), // fg
	Fg2: lipgloss.Color("#d5c4a1"), // fg2
	Fg3: lipgloss.Color("#a89984"), // gray

	// Base colors
	Gray: lipgloss.Color("#928374"), GrayLt: lipgloss.Color("#a89984"), GrayDk: lipgloss.Color("#7c6f64"),
	Red: lipgloss.Color("#fb4934"), RedLt: lipgloss.Color("#fc6d5c"), RedDk: lipgloss.Color("#cc241d"),
	Orange: lipgloss.Color("#fe8019"), OrangeLt: lipgloss.Color("#fe9a47"), OrangeDk: lipgloss.Color("#d65d0e"),
	Green: lipgloss.Color("#b8bb26"), GreenLt: lipgloss.Color("#c9cb51"), GreenDk: lipgloss.Color("#98971a"),
	Yellow: lipgloss.Color("#fabd2f"), YellowLt: lipgloss.Color("#fbcd59"), YellowDk: lipgloss.Color("#d79921"),
	Blue: lipgloss.Color("#83a598"), BlueLt: lipgloss.Color("#9db9ad"), BlueDk: lipgloss.Color("#458588"),
	Purple: lipgloss.Color("#d3869b"), PurpleLt: lipgloss.Color("#dea1b0"), PurpleDk: lipgloss.Color("#b16286"),
	Pink: lipgloss.Color("#d3869b"), PinkLt: lipgloss.Color("#dea1b0"), PinkDk: lipgloss.Color("#b16286"),
	Aqua: lipgloss.Color("#8ec07c"), AquaLt: lipgloss.Color("#a8cf99"), AquaDk: lipgloss.Color("#689d6a"),

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#fe8019"), HeroDk: lipgloss.Color("#d65d0e"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#d3869b"), // Purple
	Success: lipgloss.Color("#b8bb26"), // Green
	Warning: lipgloss.Color("#fe8019"), // Orange
	Error:   lipgloss.Color("#fb4934"), // Red
	Info:    lipgloss.Color("#83a598"), // Blue

	Border: lipgloss.RoundedBorder(),
}

// GruvboxLightPalette is the light variant of Pavel Pertsev's Gruvbox palette.
// Source: https://github.com/morhetz/gruvbox
var GruvboxLightPalette = Palette{
	Name: "GruvboxLight",

	// Background depth (light backgrounds)
	Bg0: lipgloss.Color("#fbf1c7"), // light0
	Bg1: lipgloss.Color("#ebdbb2"), // light1
	Bg2: lipgloss.Color("#d5c4a1"), // light2
	Bg3: lipgloss.Color("#bdae93"), // light3

	// Foreground depth (dark text on light bg)
	Fg0: lipgloss.Color("#282828"), // dark0 - darkest text
	Fg1: lipgloss.Color("#3c3836"), // dark1
	Fg2: lipgloss.Color("#504945"), // dark2
	Fg3: lipgloss.Color("#665c54"), // dark3 - muted

	// Base colors
	Gray: lipgloss.Color("#7c6f64"), GrayLt: lipgloss.Color("#928374"), GrayDk: lipgloss.Color("#665c54"),
	Red: lipgloss.Color("#cc241d"), RedLt: lipgloss.Color("#e04942"), RedDk: lipgloss.Color("#9d2517"),
	Orange: lipgloss.Color("#d65d0e"), OrangeLt: lipgloss.Color("#e8854a"), OrangeDk: lipgloss.Color("#af4709"),
	Green: lipgloss.Color("#98971a"), GreenLt: lipgloss.Color("#b0af30"), GreenDk: lipgloss.Color("#79740e"),
	Yellow: lipgloss.Color("#d79921"), YellowLt: lipgloss.Color("#e3ae41"), YellowDk: lipgloss.Color("#b57614"),
	Blue: lipgloss.Color("#458588"), BlueLt: lipgloss.Color("#5a9fa2"), BlueDk: lipgloss.Color("#076678"),
	Purple: lipgloss.Color("#b16286"), PurpleLt: lipgloss.Color("#c27d9e"), PurpleDk: lipgloss.Color("#8f3f71"),
	Pink: lipgloss.Color("#b16286"), PinkLt: lipgloss.Color("#c27d9e"), PinkDk: lipgloss.Color("#8f3f71"),
	Aqua: lipgloss.Color("#689d6a"), AquaLt: lipgloss.Color("#82b284"), AquaDk: lipgloss.Color("#427b58"),

	// Hero colors (tabs/borders cascade from these - light theme uses Lt for muted)
	Hero: lipgloss.Color("#d65d0e"), HeroDk: lipgloss.Color("#e8854a"),

	// Semantic colors (map to base colors for direct palette access)
	Accent:  lipgloss.Color("#b16286"), // Purple
	Success: lipgloss.Color("#98971a"), // Green
	Warning: lipgloss.Color("#d65d0e"), // Orange
	Error:   lipgloss.Color("#cc241d"), // Red
	Info:    lipgloss.Color("#458588"), // Blue

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// CHARM
// Colors: the charmtone palette, github.com/charmbracelet/x/exp/charmtone
// (MIT License, Copyright (c) 2023 Charmbracelet, Inc.). The values are
// copied here; the comment on each line gives the color's charmtone name.
// =============================================================================

// CharmPalette is a dark theme in Charm's own colors: a Pepper background
// with the Charple purple as its hero. Text on the purple is light.
var CharmPalette = Palette{
	Name: "Charm",

	// Background depth
	Bg0: lipgloss.Color("#201F26"), // Pepper
	Bg1: lipgloss.Color("#2D2C36"), // BBQ
	Bg2: lipgloss.Color("#3A3943"), // Char
	Bg3: lipgloss.Color("#4D4C57"), // Iron

	// Foreground depth
	Fg0: lipgloss.Color("#F7F6FB"), // Salt
	Fg1: lipgloss.Color("#D6D3DC"), // Steep
	Fg2: lipgloss.Color("#BFBCC8"), // Smoke
	Fg3: lipgloss.Color("#858392"), // Squid

	// Base colors
	Gray: lipgloss.Color("#605F6B"), GrayLt: lipgloss.Color("#858392"), GrayDk: lipgloss.Color("#4D4C57"), // Oyster, Squid, Iron
	Red: lipgloss.Color("#EB4268"), RedLt: lipgloss.Color("#FF577D"), RedDk: lipgloss.Color("#AB2454"), // Sriracha, Coral, Pom
	Orange: lipgloss.Color("#FF985A"), OrangeLt: lipgloss.Color("#FFB587"), OrangeDk: lipgloss.Color("#D36C64"), // Tang, Yam, Paprika
	Green: lipgloss.Color("#12C78F"), GreenLt: lipgloss.Color("#00FFB2"), GreenDk: lipgloss.Color("#00A475"), // Guac, Julep, Pickle
	Yellow: lipgloss.Color("#F5EF34"), YellowLt: lipgloss.Color("#E8FE96"), YellowDk: lipgloss.Color("#BF976F"), // Mustard, Zest, Cumin
	Blue: lipgloss.Color("#00A4FF"), BlueLt: lipgloss.Color("#4FBEFE"), BlueDk: lipgloss.Color("#007AB8"), // Malibu, Sardine, Damson
	Purple: lipgloss.Color("#6B50FF"), PurpleLt: lipgloss.Color("#8B75FF"), PurpleDk: lipgloss.Color("#5B40EC"), // Charple, Hazy, Darple
	Pink: lipgloss.Color("#FF60FF"), PinkLt: lipgloss.Color("#FF84FF"), PinkDk: lipgloss.Color("#C337E0"), // Dolly, Blush, Urchin
	Aqua: lipgloss.Color("#0ADCD9"), AquaLt: lipgloss.Color("#5CDFEA"), AquaDk: lipgloss.Color("#10B1AE"), // Turtle, Lichen, Zinc

	// Hero colors (tabs/borders cascade from these)
	Hero: lipgloss.Color("#6B50FF"), HeroDk: lipgloss.Color("#5B40EC"), // Charple, Darple

	// Semantic colors
	Secondary: lipgloss.Color("#FF60FF"), // Dolly
	Accent:    lipgloss.Color("#68FFD6"), // Bok
	Success:   lipgloss.Color("#12C78F"), // Guac
	Warning:   lipgloss.Color("#E8FE96"), // Zest
	Error:     lipgloss.Color("#EB4268"), // Sriracha
	Info:      lipgloss.Color("#00A4FF"), // Malibu

	// Charple is a deep purple, so the text on it is Salt, not the background.
	TabFocusedFg:         lipgloss.Color("#F7F6FB"),
	TabSelectedFg:        lipgloss.Color("#F7F6FB"),
	HeaderFillFg:         lipgloss.Color("#F7F6FB"),
	FooterFillFg:         lipgloss.Color("#F7F6FB"),
	BreadcrumbSelectedFg: lipgloss.Color("#F7F6FB"),

	Border: lipgloss.RoundedBorder(),
}

// =============================================================================
// THEME INSTANCES AND UTILITIES
// =============================================================================

// Built-in themes (pre-computed from palettes)
var (
	// Dracula
	Dracula      = NewTheme(DraculaPalette)
	DraculaLight = NewTheme(DraculaLightPalette)

	// Catppuccin
	CatMocha     = NewTheme(CatppuccinMochaPalette)
	CatMacchiato = NewTheme(CatppuccinMacchiatoPalette)
	CatFrappe    = NewTheme(CatppuccinFrappePalette)
	CatLatte     = NewTheme(CatppuccinLattePalette)

	// Tokyo Night
	TokyoNight = NewTheme(TokyoNightPalette)
	TokyoStorm = NewTheme(TokyoStormPalette)
	TokyoMoon  = NewTheme(TokyoMoonPalette)
	TokyoDay   = NewTheme(TokyoDayPalette)

	// Rose Pine
	RosePine     = NewTheme(RosePinePalette)
	RosePineMoon = NewTheme(RosePineMoonPalette)
	RosePineDawn = NewTheme(RosePineDawnPalette)

	// Classic
	Retro     = NewTheme(RetroPalette)
	RetroMono = NewTheme(RetroMonoPalette)
	Minimal   = NewTheme(MinimalPalette)

	// Other popular
	Solarized    = NewTheme(SolarizedPalette)
	Nord         = NewTheme(NordPalette)
	Monokai      = NewTheme(MonokaiPalette)
	Gruvbox      = NewTheme(GruvboxPalette)
	GruvboxLight = NewTheme(GruvboxLightPalette)
	Charm        = NewTheme(CharmPalette)
)

// AllThemes is a slice of all built-in themes for iteration.
var AllThemes = []*Theme{
	// Dracula
	Dracula, DraculaLight,
	// Catppuccin
	CatMocha, CatMacchiato, CatFrappe, CatLatte,
	// Tokyo Night
	TokyoNight, TokyoStorm, TokyoMoon, TokyoDay,
	// Rose Pine
	RosePine, RosePineMoon, RosePineDawn,
	// Classic
	Retro, RetroMono, Minimal,
	// Other popular
	Solarized, Nord, Monokai, Gruvbox, GruvboxLight,
	// Charm
	Charm,
}

// counterparts pairs each built-in light theme with the dark theme of its
// family.
var counterparts = [][2]*Theme{
	{DraculaLight, Dracula},
	{CatLatte, CatMocha},
	{TokyoDay, TokyoNight},
	{RosePineDawn, RosePine},
	{GruvboxLight, Gruvbox},
}

// Counterpart returns the built-in theme of the same family for the other
// kind of background: the light theme for a dark one and the dark theme for
// a light one. It returns nil for a theme without one. A pair suits
// TabApp.UseAdaptiveTheme, which follows the terminal's background.
func Counterpart(t *Theme) *Theme {
	for _, pair := range counterparts {
		switch t {
		case pair[0]:
			return pair[1]
		case pair[1]:
			return pair[0]
		}
	}
	return nil
}

// ThemeByName returns a theme by name, or nil if not found.
func ThemeByName(name string) *Theme {
	for _, t := range AllThemes {
		if t.Palette.Name == name {
			return t
		}
	}
	return nil
}

// AllPalettes is a slice of all built-in palettes for iteration.
var AllPalettes = []Palette{
	// Dracula
	DraculaPalette, DraculaLightPalette,
	// Catppuccin
	CatppuccinMochaPalette, CatppuccinMacchiatoPalette, CatppuccinFrappePalette, CatppuccinLattePalette,
	// Tokyo Night
	TokyoNightPalette, TokyoStormPalette, TokyoMoonPalette, TokyoDayPalette,
	// Rose Pine
	RosePinePalette, RosePineMoonPalette, RosePineDawnPalette,
	// Classic
	RetroPalette, RetroMonoPalette, MinimalPalette,
	// Other popular
	SolarizedPalette, NordPalette, MonokaiPalette, GruvboxPalette, GruvboxLightPalette,
	// Charm
	CharmPalette,
}

// PaletteByName returns a palette by name, or nil if not found.
func PaletteByName(name string) *Palette {
	for i := range AllPalettes {
		if AllPalettes[i].Name == name {
			return &AllPalettes[i]
		}
	}
	return nil
}
