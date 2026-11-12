package page

import (
	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/internal/keyname"
)

// ContentFunc builds a page's content: it adds rows and cells to a flexbox
// sized to the content area, on every render. The model is the app's
// *tabapp.Model, passed as interface{} because tabapp imports this package;
// assert it with model.(*tabapp.Model). navBarHeight and footerHeight are the
// heights of the nav bar and the footer in lines of content.
type ContentFunc func(fb *flexbox.FlexBox, model interface{}, navBarHeight, footerHeight int)

// KeyBinding is a page-scoped key binding: the message it dispatches and the
// label the footer shows for it.
type KeyBinding struct {
	key   string
	label string
	msg   interface{} // Message to dispatch (interface{} to avoid circular dependency)
}

// Key returns the key string.
func (kb KeyBinding) Key() string { return kb.key }

// Label returns the label for footer display.
func (kb KeyBinding) Label() string { return kb.label }

// Msg returns the message to dispatch.
func (kb KeyBinding) Msg() interface{} { return kb.msg }

// Page is one page of a TabApp: a name for its tab, a function that builds
// its content, the keys it binds, and footer labels for them.
type Page struct {
	name    string            // Page name for the tab
	content ContentFunc       // Content builder function receiving model
	footer  map[string]string // Manual footer overrides (optional)

	// Key bindings with labels (Elm architecture)
	bindings map[string]KeyBinding

	// Optional style overrides (nil means use app defaults)
	borderColor color.Color
	borderStyle *lipgloss.Border
}

// New creates a new Page with the given name.
func New(name string) *Page {
	return &Page{
		name:     name,
		footer:   make(map[string]string),
		bindings: make(map[string]KeyBinding),
	}
}

// SetContent sets the function that builds the page's content. It runs on
// every render with a fresh flexbox, the Model, and the heights of the
// nav bar and the footer:
//
//	page.SetContent(func(fb *flexbox.FlexBox, model interface{}, navBarHeight, footerHeight int) {
//	    m := model.(*tabapp.Model)
//	    // add rows and cells to fb from m
//	})
func (p *Page) SetContent(fn ContentFunc) *Page {
	p.content = fn
	return p
}

// Content returns the content builder function.
func (p *Page) Content() ContentFunc {
	return p.content
}

// SetFooter sets footer entries by key: a label for a key the footer lists,
// or an extra entry such as "↑/↓": "Select" for keys bound without a label.
// They replace the labels given to BindKey.
func (p *Page) SetFooter(footer map[string]string) *Page {
	p.footer = footer
	return p
}

// AddFooterItem adds or replaces one footer entry; see SetFooter.
func (p *Page) AddFooterItem(key, value string) *Page {
	if p.footer == nil {
		p.footer = make(map[string]string)
	}
	p.footer[key] = value
	return p
}

// SetBorderColor gives this page its own content border color, used in
// place of the app's while the page is active and the content border is on.
func (p *Page) SetBorderColor(c color.Color) *Page {
	p.borderColor = c
	return p
}

// SetBorderStyle gives this page its own content border shape, used in
// place of the app's while the page is active and the content border is on.
func (p *Page) SetBorderStyle(border lipgloss.Border) *Page {
	p.borderStyle = &border
	return p
}

// Name returns the page name for the tab.
func (p *Page) Name() string {
	return p.name
}

// Footer returns the footer key-value pairs.
func (p *Page) Footer() map[string]string {
	return p.footer
}

// BorderColor returns the override border color, or nil to use app default.
func (p *Page) BorderColor() color.Color {
	return p.borderColor
}

// BorderStyle returns the override border style, or nil to use app default.
func (p *Page) BorderStyle() *lipgloss.Border {
	return p.borderStyle
}

// BindKey binds a key to a message while this page is active. Pressing the
// key dispatches the message: the library applies its own messages and
// hands any other to the app's custom message handler. The label is shown
// in the footer; an empty label keeps the key out of it. The key is named
// as a key press reports it: "s", "space", "shift+tab". A key written " "
// is the space bar.
//
//	page.BindKey("p", "Toggle position", tabapp.ToggleNavBarPositionMsg{})
func (p *Page) BindKey(key, label string, msg interface{}) *Page {
	if p.bindings == nil {
		p.bindings = make(map[string]KeyBinding)
	}
	key = keyname.Normalize(key)
	p.bindings[key] = KeyBinding{key: key, label: label, msg: msg}
	return p
}

// UnbindKey removes a key binding.
// Returns self for chaining.
func (p *Page) UnbindKey(key string) *Page {
	delete(p.bindings, keyname.Normalize(key))
	return p
}

// Bindings returns the page's key bindings.
// Returns the internal map directly for efficiency; do not modify.
func (p *Page) Bindings() map[string]KeyBinding {
	return p.bindings
}

// HasKey checks if a key binding exists for this key.
func (p *Page) HasKey(key string) bool {
	_, ok := p.bindings[keyname.Normalize(key)]
	return ok
}
