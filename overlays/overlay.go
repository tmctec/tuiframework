package overlays

import (
	"image/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/internal/keyname"
)

// KeyHandler is the message a key binding dispatches.
//
// Deprecated: use tea.Msg directly.
type KeyHandler = tea.Msg

// Overlay is the base type for all overlay variants.
// It provides common functionality for positioning, sizing, and key handling.
type Overlay struct {
	// Content
	content     string                    // Simple text content
	flexbox     *flexbox.FlexBox          // Optional: complex layouts
	contentFunc func(fb *flexbox.FlexBox) // Optional: rebuilds the flexbox on every render

	// Positioning - use width/height percentages or fixed values
	size        int      // Percentage of container for both dimensions (default: 80)
	widthPct    int      // Width percentage (0 = use size or fixed)
	heightPct   int      // Height percentage (0 = use size or fixed)
	widthFixed  int      // Fixed width in characters (0 = use percentage)
	heightFixed int      // Fixed height in lines (0 = use percentage)
	position    Position // Justification anchor (default: CC)

	// Behavior
	blocking bool // Block page keybindings (default: true)

	// Keybindings - stores messages to dispatch when key is pressed (Elm pattern)
	keys    map[string]tea.Msg
	exitKey string // Required: key to dismiss overlay

	// Styling
	style lipgloss.Style

	// Title line drawn above the content (empty = none)
	title      string
	titleStyle lipgloss.Style

	// State
	visible bool

	// Z-axis stacking
	zIndex int // Higher = renders on top (CSS-style, default 0)

	// cache holds the last render of an overlay with fixed content. Every
	// setter that changes what Render returns clears it.
	cache renderCache
}

// renderCache is an overlay's rendered text for one container size.
type renderCache struct {
	valid         bool
	width, height int // the container the text was rendered for
	text          string
}

// New creates a new base Overlay with sensible defaults.
func New() *Overlay {
	return &Overlay{
		size:     80,
		position: PositionCenter,
		blocking: true,
		keys:     make(map[string]tea.Msg),
		style: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2),
		titleStyle: lipgloss.NewStyle().Bold(true),
	}
}

// SetContent sets simple text content for the overlay.
// For complex layouts, use SetFlexbox instead.
func (o *Overlay) SetContent(content string) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.content = content
	return o
}

// SetFlexbox sets a flexbox for complex content layouts.
// Takes precedence over simple text content.
func (o *Overlay) SetFlexbox(fb *flexbox.FlexBox) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.flexbox = fb
	return o
}

// SetContentFunc sets content via a builder function, like a Page: the
// function runs on every render with a fresh flexbox sized to the content
// area, so content that reads application state stays current. It also runs
// once here, with an unsized flexbox, so Flexbox is never nil afterwards.
func (o *Overlay) SetContentFunc(fn func(fb *flexbox.FlexBox)) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.contentFunc = fn
	fb := flexbox.New(0, 0)
	fn(fb)
	o.flexbox = fb
	return o
}

// clampPercent keeps a size percentage inside 1 to 100.
func clampPercent(percent int) int {
	return min(max(percent, 1), 100)
}

// SetSize sets the overlay size as a percentage of the container (1-100).
// This sets both width and height to the same percentage.
// For different width/height, use SetWidthPct and SetHeightPct.
func (o *Overlay) SetSize(percent int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.size = clampPercent(percent)
	o.widthPct = 0 // Clear individual settings
	o.heightPct = 0
	return o
}

// SetWidthPct sets the overlay width as a percentage (1-100).
// Use with SetHeightPct for rectangular overlays.
func (o *Overlay) SetWidthPct(percent int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.widthPct = clampPercent(percent)
	o.widthFixed = 0 // Clear fixed value
	return o
}

// SetHeightPct sets the overlay height as a percentage (1-100).
// Use with SetWidthPct for rectangular overlays.
func (o *Overlay) SetHeightPct(percent int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.heightPct = clampPercent(percent)
	o.heightFixed = 0 // Clear fixed value
	return o
}

// SetWidth sets a fixed overlay width in characters.
// Takes precedence over percentage-based sizing.
func (o *Overlay) SetWidth(chars int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.widthFixed = max(chars, 1)
	o.widthPct = 0 // Clear percentage
	return o
}

// SetHeight sets a fixed overlay height in lines.
// Takes precedence over percentage-based sizing.
func (o *Overlay) SetHeight(lines int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.heightFixed = max(lines, 1)
	o.heightPct = 0 // Clear percentage
	return o
}

// SetPosition sets where the overlay is anchored on screen.
func (o *Overlay) SetPosition(pos Position) *Overlay {
	o.position = pos
	return o
}

// SetBlocking sets whether the overlay blocks page keybindings.
func (o *Overlay) SetBlocking(blocking bool) *Overlay {
	o.blocking = blocking
	return o
}

// SetStyle sets the lipgloss style for the overlay container.
func (o *Overlay) SetStyle(style lipgloss.Style) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = style
	return o
}

// SetBorder sets the border style.
func (o *Overlay) SetBorder(border lipgloss.Border) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.Border(border)
	return o
}

// SetBorderForeground sets the border foreground color. One color replaces
// a gradient set with SetBorderBlend.
func (o *Overlay) SetBorderForeground(c color.Color) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.UnsetBorderForegroundBlend().BorderForeground(c)
	return o
}

// SetBorderBlend draws the border as a gradient through two or more colors.
// Fewer than two colors remove the gradient. Returns self for chaining.
func (o *Overlay) SetBorderBlend(colors ...color.Color) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	if len(colors) < 2 {
		o.style = o.style.UnsetBorderForegroundBlend()
		return o
	}
	o.style = o.style.BorderForegroundBlend(colors...)
	return o
}

// SetBorderBackground sets the border background color.
func (o *Overlay) SetBorderBackground(c color.Color) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.BorderBackground(c)
	return o
}

// SetPadding sets the padding (vertical, horizontal).
func (o *Overlay) SetPadding(vertical, horizontal int) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.Padding(vertical, horizontal)
	return o
}

// SetBackground sets the background color.
func (o *Overlay) SetBackground(c color.Color) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.Background(c)
	return o
}

// SetTitle sets the title line drawn above the content.
// An empty title draws nothing and leaves the full height to the content.
func (o *Overlay) SetTitle(title string) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.title = title
	return o
}

// Title returns the title line.
func (o *Overlay) Title() string {
	return o.title
}

// SetTitleStyle sets the style of the title line.
func (o *Overlay) SetTitleStyle(style lipgloss.Style) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.titleStyle = style
	return o
}

// SetForeground sets the foreground (text) color.
func (o *Overlay) SetForeground(c color.Color) *Overlay {
	o.cache = renderCache{} // the next Render draws this change
	o.style = o.style.Foreground(c)
	return o
}

// SetExitKey sets the key that dismisses this overlay.
// This is REQUIRED for all overlays.
func (o *Overlay) SetExitKey(key string) *Overlay {
	o.exitKey = keyname.Normalize(key)
	return o
}

// RegisterKey registers a message to dispatch when key is pressed (Elm pattern).
// The key is named as a key press reports it; a key written " " is the
// space bar. Returns self for chaining.
func (o *Overlay) RegisterKey(key string, msg tea.Msg) *Overlay {
	if o.keys == nil {
		o.keys = make(map[string]tea.Msg)
	}
	o.keys[keyname.Normalize(key)] = msg
	return o
}

// RegisterKeys registers multiple messages to dispatch when keys are pressed (Elm pattern).
// Returns self for chaining.
func (o *Overlay) RegisterKeys(bindings map[string]tea.Msg) *Overlay {
	if o.keys == nil {
		o.keys = make(map[string]tea.Msg)
	}
	for key, msg := range bindings {
		o.keys[keyname.Normalize(key)] = msg
	}
	return o
}

// Show marks the overlay as visible.
func (o *Overlay) Show() *Overlay {
	o.visible = true
	return o
}

// Hide marks the overlay as not visible.
func (o *Overlay) Hide() *Overlay {
	o.visible = false
	return o
}

// IsVisible returns whether the overlay is currently visible.
func (o *Overlay) IsVisible() bool {
	return o.visible
}

// SetZIndex sets the z-axis stacking order (CSS-style).
// Higher values render on top. Default is 0.
func (o *Overlay) SetZIndex(z int) *Overlay {
	o.zIndex = z
	return o
}

// ZIndex returns the z-axis stacking order.
func (o *Overlay) ZIndex() int {
	return o.zIndex
}

// IsBlocking returns whether the overlay blocks page keybindings.
func (o *Overlay) IsBlocking() bool {
	return o.blocking
}

// ExitKey returns the configured exit key.
func (o *Overlay) ExitKey() string {
	return o.exitKey
}

// Keys returns the overlay's registered key->message mappings (Elm pattern).
func (o *Overlay) Keys() map[string]tea.Msg {
	return o.keys
}

// HasKey checks if a key is registered on this overlay.
func (o *Overlay) HasKey(key string) bool {
	_, ok := o.keys[key]
	return ok
}

// Size returns the configured size percentage.
func (o *Overlay) Size() int {
	return o.size
}

// WidthPct returns the configured width percentage (0 if using Size() or fixed).
func (o *Overlay) WidthPct() int {
	return o.widthPct
}

// HeightPct returns the configured height percentage (0 if using Size() or fixed).
func (o *Overlay) HeightPct() int {
	return o.heightPct
}

// WidthFixed returns the configured fixed width (0 if using percentage).
func (o *Overlay) WidthFixed() int {
	return o.widthFixed
}

// HeightFixed returns the configured fixed height (0 if using percentage).
func (o *Overlay) HeightFixed() int {
	return o.heightFixed
}

// EffectiveWidth returns the width percentage in effect: WidthPct, or Size.
//
// Deprecated: use CalculateDimensions, which also honours fixed sizes.
func (o *Overlay) EffectiveWidth() int {
	if o.widthPct > 0 {
		return o.widthPct
	}
	return o.size
}

// EffectiveHeight returns the height percentage in effect: HeightPct, or Size.
//
// Deprecated: use CalculateDimensions, which also honours fixed sizes.
func (o *Overlay) EffectiveHeight() int {
	if o.heightPct > 0 {
		return o.heightPct
	}
	return o.size
}

// CalculateDimensions returns the actual overlay dimensions in characters/lines.
// Fixed values take precedence over percentages; neither exceeds the container.
func (o *Overlay) CalculateDimensions(containerWidth, containerHeight int) (width, height int) {
	// Width: fixed > widthPct > size
	if o.widthFixed > 0 {
		width = o.widthFixed
	} else if o.widthPct > 0 {
		width = containerWidth * o.widthPct / 100
	} else {
		width = containerWidth * o.size / 100
	}

	// Height: fixed > heightPct > size
	if o.heightFixed > 0 {
		height = o.heightFixed
	} else if o.heightPct > 0 {
		height = containerHeight * o.heightPct / 100
	} else {
		height = containerHeight * o.size / 100
	}

	return min(width, containerWidth), min(height, containerHeight)
}

// Position returns the configured position.
func (o *Overlay) Position() Position {
	return o.position
}

// Content returns the simple text content.
func (o *Overlay) Content() string {
	return o.content
}

// Flexbox returns the flexbox content (may be nil).
func (o *Overlay) Flexbox() *flexbox.FlexBox {
	return o.flexbox
}

// Style returns the overlay's style.
func (o *Overlay) Style() lipgloss.Style {
	return o.style
}

// Render renders the overlay content at the given container dimensions.
// An overlay whose content is text renders once per container size and per
// change; one with a flexbox or a content function renders every time,
// because what they show can change between frames.
func (o *Overlay) Render(containerWidth, containerHeight int) string {
	fixed := o.contentFunc == nil && o.flexbox == nil
	if fixed && o.cache.valid && o.cache.width == containerWidth && o.cache.height == containerHeight {
		return o.cache.text
	}

	// The computed dimensions are the overlay's footprint, frame included.
	overlayWidth, overlayHeight := o.CalculateDimensions(containerWidth, containerHeight)

	// The style's Width and Height cover the border and the padding, so they
	// take the footprint as it is; the content area is what the frame leaves,
	// and a title takes 2 more rows.
	contentWidth := overlayWidth - o.style.GetHorizontalFrameSize()
	contentHeight := overlayHeight - o.style.GetVerticalFrameSize()
	if o.title != "" {
		contentHeight -= 2
	}

	var content string
	switch {
	case o.contentFunc != nil:
		fb := flexbox.New(contentWidth, contentHeight)
		o.contentFunc(fb)
		o.flexbox = fb
		content = fb.Render()
	case o.flexbox != nil:
		o.flexbox.SetWidth(contentWidth)
		o.flexbox.SetHeight(contentHeight)
		content = o.flexbox.Render()
	default:
		content = o.content
	}
	if o.title != "" {
		content = o.titleStyle.Render(o.title) + "\n\n" + content
	}

	text := o.style.
		Width(overlayWidth).
		Height(overlayHeight).
		Render(content)
	if fixed {
		o.cache = renderCache{valid: true, width: containerWidth, height: containerHeight, text: text}
	}
	return text
}

// Validate checks that the overlay is properly configured.
// Returns an error if validation fails.
func (o *Overlay) Validate() error {
	if o.exitKey == "" {
		return ErrNoExitKey
	}
	return nil
}
