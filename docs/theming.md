# Theming System

The TUI Framework uses a 4-tier color cascade system. You only need to define what you want to customize - everything else cascades from sensible defaults.

## Quick Start

```go
// Minimal custom theme - just set base colors + hero
var MyPalette = styles.Palette{
    Name: "MyTheme",

    // Required: Background/Foreground
    Bg0: lipgloss.Color("#1a1b26"), Bg1: lipgloss.Color("#24283b"),
    Bg2: lipgloss.Color("#414868"), Bg3: lipgloss.Color("#565f89"),
    Fg0: lipgloss.Color("#c0caf5"), Fg1: lipgloss.Color("#a9b1d6"),
    Fg2: lipgloss.Color("#9aa5ce"), Fg3: lipgloss.Color("#565f89"),

    // Required: Base colors
    Red: lipgloss.Color("#f7768e"), Green: lipgloss.Color("#9ece6a"),
    Blue: lipgloss.Color("#7aa2f7"), // ... etc

    // Set your hero - everything else cascades from it
    Hero: lipgloss.Color("#7aa2f7"),

    Border: lipgloss.RoundedBorder(),
}

var MyTheme = styles.NewTheme(MyPalette)
```

A color is an `image/color.Color`: write one with `lipgloss.Color("#1a1b26")`
or `lipgloss.Color("212")`. A color left out of a palette is nil, and the
cascade fills it in. `styles.Hex(c)` returns a color as `"#rrggbb"`, and
`styles.Contrast(a, b)` the WCAG contrast ratio of two colors.

## The 4-Tier Cascade

```
TIER 1: Base Colors (foundation)
    │
    ├── Bg0, Bg1, Bg2, Bg3        (backgrounds)
    ├── Fg0, Fg1, Fg2, Fg3        (foregrounds)
    └── Red, Orange, Yellow, Green, Blue, Purple, Pink, Aqua
        (each with Lt/Dk variants)
    │
    ▼
TIER 2: Hero Colors (theme personality)
    │
    ├── Hero    ← defaults to Purple
    └── HeroDk  ← defaults to PurpleDk
    │
    ▼
TIER 3: Semantic Colors (meaning)
    │
    ├── Primary     ← defaults to Hero
    ├── PrimaryDark ← defaults to HeroDk
    ├── Secondary   ← defaults to Blue
    ├── Accent      ← defaults to Aqua
    ├── Success     ← defaults to Green
    ├── Warning     ← defaults to Yellow
    ├── Error       ← defaults to Red
    └── Info        ← defaults to Aqua
    │
    ▼
TIER 4: UI Element Colors (specific components)
    │
    ├── Tab colors     ← cascade from Primary/Hero
    ├── Border colors  ← cascade from Primary or Bg3
    ├── Notification   ← cascade from Bg0
    ├── Footer         ← keys Primary, values Fg2, key glow GreenLt
    ├── Band fills     ← Bg0 on Primary, for a filled header or footer
    └── Modal          ← cascade from Bg0 and Primary
```

## All Palette Fields

### Tier 1: Base Colors (Required)

| Field | Type | Description |
|-------|------|-------------|
| `Name` | string | Theme identifier (e.g., "Dracula", "Nord") |

#### Backgrounds (layered depth)
| Field | Description |
|-------|-------------|
| `Bg0` | Base background (darkest) |
| `Bg1` | Elevated surfaces (cards, inactive tabs) |
| `Bg2` | Selected/hover states |
| `Bg3` | Borders, subtle elements |

#### Foregrounds (text hierarchy)
| Field | Description |
|-------|-------------|
| `Fg0` | Brightest (headings, emphasis) |
| `Fg1` | Primary text |
| `Fg2` | Secondary/muted text |
| `Fg3` | Disabled/hint text |

#### Base Colors (each with Lt/Dk variants)
| Color | Fields | Typical Use |
|-------|--------|-------------|
| Gray | `Gray`, `GrayLt`, `GrayDk` | Neutral elements |
| Red | `Red`, `RedLt`, `RedDk` | Errors, destructive |
| Orange | `Orange`, `OrangeLt`, `OrangeDk` | Warnings |
| Yellow | `Yellow`, `YellowLt`, `YellowDk` | Caution, highlights |
| Green | `Green`, `GreenLt`, `GreenDk` | Success, positive |
| Blue | `Blue`, `BlueLt`, `BlueDk` | Info, links |
| Purple | `Purple`, `PurpleLt`, `PurpleDk` | Default hero |
| Pink | `Pink`, `PinkLt`, `PinkDk` | Accent |
| Aqua | `Aqua`, `AquaLt`, `AquaDk` | Secondary accent |

---

### Tier 2: Hero Colors (Optional)

The "hero" defines your theme's signature accent color.

| Field | Default | Description |
|-------|---------|-------------|
| `Hero` | `Purple` | Primary accent color (tabs, borders, highlights) |
| `HeroDk` | `PurpleDk` | Darker variant (selected states) |

**Example:** For a blue-themed app:
```go
Hero:   lipgloss.Color("#7aa2f7"),  // Bright blue
HeroDk: lipgloss.Color("#5a7abf"),  // Muted blue
```

---

### Tier 3: Semantic Colors (Optional)

Override these to customize meaning-based colors.

| Field | Default | Description |
|-------|---------|-------------|
| `Primary` | `Hero` | Main accent (buttons, active elements) |
| `PrimaryDark` | `HeroDk` | Darker primary (selected states) |
| `Secondary` | `Blue` | Secondary accent |
| `Accent` | `Aqua` | Tertiary accent |
| `Success` | `Green` | Success states, tips |
| `Warning` | `Yellow` | Warning states |
| `Error` | `Red` | Error states, destructive |
| `Info` | `Aqua` | Informational |

---

### Tier 4: UI Element Colors (Optional)

Fine-tune specific UI components.

#### Tabs and Nav Bar
| Field | Default | Description |
|-------|---------|-------------|
| `TabFocusedBg` | `Primary` | Background when tab is focused |
| `TabFocusedFg` | `Bg0` | Text color when tab is focused |
| `TabSelectedBg` | `PrimaryDark` | Background when selected but not focused |
| `TabSelectedFg` | `Bg0` | Text color when selected but not focused |
| `TabInactiveBg` | `Bg1` | Background for unselected tabs |
| `TabInactiveFg` | `Fg3` | Text color for unselected tabs |
| `NavBarBorderColor` | `Primary` | Nav bar border color |

#### Borders
| Field | Default | Description |
|-------|---------|-------------|
| `HeaderBorderColor` | `Primary` | Header band border |
| `ContentBorderColor` | `Bg3` | Border around content area |
| `FooterBorderColor` | `Bg3` | Border above footer |

#### Footer
| Field | Default | Description |
|-------|---------|-------------|
| `FooterKeyColor` | `Primary` | Footer key label color |
| `FooterValueColor` | `Fg2` | Footer value text color |
| `KeyPressHighlight` | `GreenLt` | Key press glow effect |

#### Band Fills
| Field | Default | Description |
|-------|---------|-------------|
| `HeaderFillBg` | `Primary` | Background of a filled header |
| `HeaderFillFg` | `Bg0` | Text on a filled header and on labels (`LabelStyle`) |
| `FooterFillBg` | `Primary` | Background of a filled footer |
| `FooterFillFg` | `Bg0` | Keys and labels on a filled footer |

These feed `Theme.HeaderFill` and `Theme.FooterFill`, the colors of a header
or footer whose frame is `tabapp.FrameFill`. A fill reverses the text: body
text on the hero color is hard to read in most themes, while the base
background on the hero color reaches a contrast of at least 3:1 in every
built-in theme. A palette that sets a fill background sets the text color
that suits it.

#### Breadcrumbs
| Field | Default | Description |
|-------|---------|-------------|
| `BreadcrumbBg` | `Accent` | Inactive breadcrumb background |
| `BreadcrumbFg` | `Bg0` | Inactive breadcrumb text |
| `BreadcrumbSelectedBg` | `Primary` | Active breadcrumb background |
| `BreadcrumbSelectedFg` | `Bg0` | Active breadcrumb text |

These feed `Theme.Breadcrumb` and `Theme.BreadcrumbSelected`, styles for an
application's own trails. The nav bar draws its crumbs with
`TabSelectedMuted` and the focused crumb with `TabActive`, so they match the
tabs beside them.

#### Modals
| Field | Default | Description |
|-------|---------|-------------|
| `ModalBg` | `Bg0` | Modal background |
| `ModalBorderColor` | `HeaderBorderColor` | Modal border color |
| `ModalTitleColor` | `Primary` | Modal title text color |
| `ModalBorder` | `DoubleBorder()` | Modal border style |

#### Notifications
| Field | Default | Description |
|-------|---------|-------------|
| `NotifBg` | `Bg0` | Notification background. Use `Bg1`/`Bg2` for filled look |

#### Border Styles
| Field | Default | Description |
|-------|---------|-------------|
| `Border` | (required) | Default border style: `RoundedBorder()`, `DoubleBorder()`, `NormalBorder()`, `HiddenBorder()` |
| `ModalBorder` | `DoubleBorder()` | Modal-specific border style |

#### Notification Icons
| Field | Default | Description |
|-------|---------|-------------|
| `IconInfo` | `"💬 Info:"` | Info notification prefix |
| `IconNote` | `"📌 Note:"` | Note notification prefix |
| `IconTip` | `"💡 Tip:"` | Tip notification prefix |
| `IconImportant` | `"❗ Important:"` | Important notification prefix |
| `IconSuccess` | `"✅ Success:"` | Success notification prefix |
| `IconWarning` | `"🔶 Warning:"` | Warning notification prefix |
| `IconCaution` | `"⛔ Caution:"` | Caution notification prefix |
| `IconError` | `"❌ Error:"` | Error notification prefix |

**ASCII Icons:** For terminal compatibility:
```go
IconInfo: "[i]", IconNote: "[*]", IconTip: "[?]", IconImportant: "[!]",
IconSuccess: "[OK]", IconWarning: "[!!]", IconCaution: "[!!!]", IconError: "[ERR]",
```

---

## Examples

### Example 1: Minimal Theme (Just Hero)

```go
var BluePalette = styles.Palette{
    Name: "Blue",

    // Backgrounds
    Bg0: lipgloss.Color("#0d1117"), Bg1: lipgloss.Color("#161b22"),
    Bg2: lipgloss.Color("#21262d"), Bg3: lipgloss.Color("#30363d"),

    // Foregrounds
    Fg0: lipgloss.Color("#f0f6fc"), Fg1: lipgloss.Color("#c9d1d9"),
    Fg2: lipgloss.Color("#8b949e"), Fg3: lipgloss.Color("#6e7681"),

    // Base colors
    Red: lipgloss.Color("#f85149"), RedLt: lipgloss.Color("#ff7b72"), RedDk: lipgloss.Color("#da3633"),
    // ... other base colors ...
    Blue: lipgloss.Color("#58a6ff"), BlueLt: lipgloss.Color("#79c0ff"), BlueDk: lipgloss.Color("#388bfd"),

    // Just set Hero - everything cascades
    Hero: lipgloss.Color("#58a6ff"), HeroDk: lipgloss.Color("#388bfd"),

    Border: lipgloss.RoundedBorder(),
}
```

### Example 2: Override Specific Elements

```go
var CustomPalette = styles.Palette{
    // ... base colors ...

    Hero: lipgloss.Color("#ff79c6"),

    // Override: I want different tab colors
    TabFocusedBg:  lipgloss.Color("#44475a"),  // Custom, not Hero
    TabSelectedBg: lipgloss.Color("#282a36"),  // Custom, not HeroDk

    // Override: Filled notifications
    NotifBg: lipgloss.Color("#282a36"),  // Bg1 for contrast

    // Override: Subtle header border
    HeaderBorderColor: lipgloss.Color("#6272a4"),  // Not Hero
}
```

### Example 3: Monochrome Theme

```go
var MonoPalette = styles.Palette{
    Name: "Mono",

    // All grays
    Bg0: lipgloss.Color("#1a1a1a"), Bg1: lipgloss.Color("#2a2a2a"),
    Bg2: lipgloss.Color("#3a3a3a"), Bg3: lipgloss.Color("#4a4a4a"),
    Fg0: lipgloss.Color("#ffffff"), Fg1: lipgloss.Color("#e0e0e0"),
    Fg2: lipgloss.Color("#a0a0a0"), Fg3: lipgloss.Color("#606060"),

    // All base colors = same gray (for true monochrome)
    Gray: lipgloss.Color("#808080"), GrayLt: lipgloss.Color("#a0a0a0"), GrayDk: lipgloss.Color("#606060"),
    Red: lipgloss.Color("#808080"), RedLt: lipgloss.Color("#a0a0a0"), RedDk: lipgloss.Color("#606060"),
    // ... same for all colors ...

    // Hero = white for contrast
    Hero: lipgloss.Color("#ffffff"), HeroDk: lipgloss.Color("#c0c0c0"),

    Border: lipgloss.NormalBorder(),
}
```

---

## A Palette from Three Colors

`PaletteFrom` derives the whole palette from a background, a body text color
and one accent:

```go
c := lipgloss.Color
p, err := styles.PaletteFrom(c("#101418"), c("#d7dde4"), c("#7aa2f7"))
if err != nil {
    // errors.Is(err, styles.ErrLowContrast)
}
p.Name = "Midnight"
app.UseTheme(styles.NewTheme(p))
```

The background and text levels are steps between the first two colors. The
accent becomes the hero color, which tabs, borders and fills take. The base
colors, red through aqua, are a set chosen for a dark or a light background,
whichever the background is. The palette is named `Custom` and has a rounded
border; set any field on the result before `NewTheme`.

`PaletteFrom` returns `styles.ErrLowContrast` when the text has less than
4.5:1 contrast on the background, or the accent less than 3:1.
`MustPaletteFrom` panics instead, for colors known to be readable.

## Gradient Borders

Two colors in `Palette.BorderBlend` draw the content border, the nav bar
border and modal borders as a gradient from the first to the second:

```go
p.BorderBlend = [2]color.Color{lipgloss.Color("#ff00ff"), lipgloss.Color("#00ffff")}
```

Without both colors the borders keep their single colors. On an app,
`SetBorderBlend(colors...)` sets a gradient through any number of colors for
the content and nav bar borders, and with fewer than two removes it. A page's
own `SetBorderColor` and a modal's own border color replace the gradient.

## Built-in Themes

| Theme | Hero Color | Style |
|-------|------------|-------|
| `Dracula` | Purple | Dark, vibrant |
| `DraculaLight` | Purple | Light variant |
| `CatMocha` | Mauve | Catppuccin dark |
| `CatMacchiato` | Mauve | Catppuccin medium |
| `CatFrappe` | Mauve | Catppuccin light-dark |
| `CatLatte` | Mauve | Catppuccin light |
| `TokyoNight` | Blue | Dark, cool |
| `TokyoStorm` | Blue | Darker variant |
| `TokyoMoon` | Blue | Purple tint |
| `TokyoDay` | Blue | Light variant |
| `RosePine` | Rose | Dark, warm |
| `RosePineMoon` | Rose | Medium variant |
| `RosePineDawn` | Rose | Light variant |
| `Retro` | Green | CRT green monochrome |
| `RetroMono` | Amber | CRT amber monochrome |
| `Minimal` | White | Clean, hidden borders |
| `Solarized` | Blue | Ethan Schoonover's classic |
| `Nord` | Blue | Arctic/Nordic |
| `Monokai` | Yellow | Wimer Hazenberg's classic |
| `Gruvbox` | Orange | Retro groove |
| `GruvboxLight` | Orange | Light variant |
| `Charm` | Purple | Charm's charmtone colors, dark |

---

## The Theme

`NewTheme` returns a `*Theme` of ready-made Lip Gloss styles. `TabApp.UseTheme`
applies the tab, border, footer, modal and notification parts; the rest are
for page content.

| Group | Fields |
|---|---|
| Tabs | `TabInactive`, `TabActive`, `TabSelectedMuted`, `NavBarBorder`, `NavBarBorderColor` |
| Chrome | `HeaderBorderColor`, `HeaderFill`, `ContentBorder`, `ContentBorderColor`, `BorderBlend`, `Footer`, `FooterBorder`, `FooterBorderColor`, `FooterKey`, `FooterValue`, `FooterGlow`, `FooterFill` |
| Breadcrumbs | `Breadcrumb`, `BreadcrumbSelected` |
| Notifications | `NotifInfo` ... `NotifError` (`styles.NotificationStyle`: foreground, background, border, icon), `NotifBorder` |
| Modals | `ModalBackground`, `ModalBorder`, `ModalBorderColor`, `ModalTitle` |
| Text | `TextPrimary` (Fg1), `TextSecondary` (Fg2), `TextMuted` (Fg3), `TextHeading` (Fg0, bold) |
| Cell styles | `BaseStyle` (Fg0/Bg0), `ContentStyle` (Fg1/Bg0), `MutedStyle` (Fg2/Bg0), `DisabledStyle` (Fg3/Bg0), `AccentStyle` (Hero/Bg0) |
| Components | `CardStyle`, `PanelStyle`, `HeaderStyle`, `LabelStyle` |
| Status | `StatusSuccess`, `StatusWarning`, `StatusError`, `StatusInfo` |

`Theme.Palette` is the resolved palette: every optional color, border and
icon holds the value the cascade gave it, so `theme.Palette.Primary` and
`theme.Palette.FooterKeyColor` are never empty. `Theme.Name()` returns the
palette's name, `Theme.KeyGlowColor()` the footer glow color, and
`Theme.IsDark()` whether the theme has a dark background.

Pass `BaseStyle` to a flexbox with `fb.UseStyle(theme.BaseStyle)` and every
cell inherits the theme's foreground and background.

## Registry

`AllThemes` and `AllPalettes` list the built-in themes and palettes in the
same order. `ThemeByName("Nord")` and `PaletteByName("Nord")` look one up
and return nil for an unknown name. A theme's name is the name of its
variable: `styles.RosePineMoon` is `"RosePineMoon"`.

## Applying Themes

```go
// Use a built-in theme
app.UseTheme(styles.Dracula)

// Use a custom theme
app.UseTheme(styles.NewTheme(MyPalette))

// Cycle themes with keybinding
app.MustBindKey("right", "Next Theme", tabapp.CycleThemeMsg{
    Direction: 1,
    Themes: []*styles.Theme{styles.Dracula, styles.Nord, styles.Gruvbox},
})

// Or set one from a message
app.MustBindKey("n", "Nord", tabapp.SetThemeMsg{Theme: styles.Nord})
```

`UseTheme` restyles the nav bar, borders and footer at once, re-themes the
modals shown through `ShowModal`, and updates `overlays.NotificationConfig`
so new notifications match. Pages re-render from `Model.Theme` on the next
frame.

## Following the Terminal Background

```go
app.UseAdaptiveTheme(styles.CatLatte, styles.CatMocha) // light, dark
```

The app asks the terminal for its background color when it starts and
switches to the light theme on a light background and the dark theme on a
dark one. Until the terminal answers, and on terminals that never do, the
dark theme is active. `Model.DarkBackground` holds the answer and
`Model.BackgroundKnown` tells whether one arrived, with or without an
adaptive pair.

Choosing a theme with `UseTheme`, `SetThemeMsg` or `CycleThemeMsg` ends the
following. `UseAdaptiveTheme` or `UseAdaptiveThemeMsg{Light, Dark}` starts it
again. `Model.FollowsBackground` tells which state the app is in.

`Theme.IsDark()` classifies any theme by its base background.
`styles.Counterpart(theme)` returns the built-in theme of the same family for
the other kind of background, or nil:

| Light | Dark |
|-------|------|
| `DraculaLight` | `Dracula` |
| `CatLatte` | `CatMocha` |
| `TokyoDay` | `TokyoNight` |
| `RosePineDawn` | `RosePine` |
| `GruvboxLight` | `Gruvbox` |
