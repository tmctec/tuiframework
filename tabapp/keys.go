package tabapp

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/internal/keyname"
	"github.com/tmctec/tuiframework/page"
)

// KeyHandler is the message a key binding dispatches.
//
// Deprecated: use tea.Msg directly.
type KeyHandler = tea.Msg

// internalKeyHandler runs for the hardcoded keys (ctrl+c, ctrl+q).
type internalKeyHandler func() tea.Cmd

// KeyBinding is a universal key binding: the message it dispatches and the
// label the footer shows for it.
type KeyBinding struct {
	key   string
	label string
	msg   tea.Msg
}

// Key returns the key string.
func (kb KeyBinding) Key() string { return kb.key }

// Label returns the label for footer display.
func (kb KeyBinding) Label() string { return kb.label }

// Msg returns the message to dispatch when this key is pressed.
func (kb KeyBinding) Msg() tea.Msg { return kb.msg }

// BindKey registers a universal key binding that works on all pages (Elm pattern).
// The label is displayed in the footer; use empty string to hide from footer.
// The msg is dispatched when the key is pressed: the library applies its own
// messages and hands any other to the custom message handler.
// The key is named as a key press reports it: "s", "space", "shift+tab",
// "ctrl+h". A key written " " is the space bar.
// Returns error if key is already registered as hardcoded.
func (a *TabApp) BindKey(key, label string, msg tea.Msg) error {
	key = keyname.Normalize(key)
	if _, isHardcoded := a.hardcodedKeys[key]; isHardcoded {
		return fmt.Errorf("key %q is hardcoded and cannot be overridden", key)
	}
	a.universalBindings[key] = KeyBinding{key: key, label: label, msg: msg}
	return nil
}

// MustBindKey registers a universal key binding, panicking on error (Elm pattern).
// Returns self for chaining.
func (a *TabApp) MustBindKey(key, label string, msg tea.Msg) *TabApp {
	if err := a.BindKey(key, label, msg); err != nil {
		panic(err)
	}
	return a
}

// UnbindKey removes a universal key binding, which frees the key for pages
// to bind. Returns self for chaining.
func (a *TabApp) UnbindKey(key string) *TabApp {
	delete(a.universalBindings, keyname.Normalize(key))
	return a
}

// IsUniversalKey checks if a key is registered as a universal binding.
func (a *TabApp) IsUniversalKey(key string) bool {
	_, ok := a.universalBindings[keyname.Normalize(key)]
	return ok
}

// UniversalBindings returns a copy of the universal key bindings by key.
func (a *TabApp) UniversalBindings() map[string]KeyBinding {
	return maps.Clone(a.universalBindings)
}

// startKeyGlow highlights key in the footer for the glow duration.
// Returns nil when the glow effect is disabled.
func (a *TabApp) startKeyGlow(key string) tea.Cmd {
	if !a.model.KeyGlowEnabled {
		return nil
	}
	a.keyGlowKey = key
	a.keyGlowTimerID++
	id := a.keyGlowTimerID
	return tea.Tick(a.keyGlowDuration, func(_ time.Time) tea.Msg {
		return keyGlowTimeoutMsg{id: id}
	})
}

// IsHardcoded checks if a key is hardcoded (cannot be overridden).
func (a *TabApp) IsHardcoded(key string) bool {
	_, ok := a.hardcodedKeys[key]
	return ok
}

// footerKeyOrder lists the keys the footer shows first, in this order; the
// rest follow sorted.
var footerKeyOrder = []string{"↑/↓", "←/→", "Enter", "Space", "Tab", "Esc", "q"}

// renderPageFooter auto-builds footer from key bindings and merges with manual overrides.
// Priority: manual SetFooter() overrides > page bindings > universal bindings
func (a *TabApp) renderPageFooter(p *page.Page) string {
	// Build combined footer map: start with universal, overlay page, then manual overrides
	combined := make(map[string]string)

	// 1. Add universal bindings (with non-empty labels)
	for key, binding := range a.universalBindings {
		if binding.label != "" {
			combined[key] = binding.label
		}
	}

	// 2. Add page bindings (with non-empty labels) - overrides universal
	for key, binding := range p.Bindings() {
		if binding.Label() != "" {
			combined[key] = binding.Label()
		}
	}

	// 3. Apply manual footer overrides from SetFooter()
	for key, value := range p.Footer() {
		combined[key] = value
	}

	if len(combined) == 0 {
		return "Tab: Navigate • Ctrl+Q: Quit" // Default
	}

	// Build footer with styled keys and values
	var items []string
	keyStyle, keyGlowStyle, valueStyle := a.footerStyles()

	// Helper to get appropriate key style (glow or normal)
	getKeyStyle := func(key string) lipgloss.Style {
		if a.keyGlowKey == key {
			return keyGlowStyle
		}
		return keyStyle
	}

	// First add ordered keys if they exist
	for _, key := range footerKeyOrder {
		if value, exists := combined[key]; exists {
			items = append(items, fmt.Sprintf("%s %s",
				getKeyStyle(key).Render(key),
				valueStyle.Render(value)))
		}
	}

	// Then add any remaining keys (sorted for consistency)
	var remainingKeys []string
	for key := range combined {
		if !slices.Contains(footerKeyOrder, key) {
			remainingKeys = append(remainingKeys, key)
		}
	}
	slices.Sort(remainingKeys)
	for _, key := range remainingKeys {
		items = append(items, fmt.Sprintf("%s %s",
			getKeyStyle(key).Render(key),
			valueStyle.Render(combined[key])))
	}

	// Style the separator so it doesn't lose color after ANSI resets
	return strings.Join(items, valueStyle.Render(a.footerSeparator))
}

// footerStyles returns the styles of the footer's keys, of the key that is
// glowing, and of its labels and separators. On a filled footer the keys and
// labels take the fill's reversed colors, and the glowing key shows as a
// chip of the content background.
func (a *TabApp) footerStyles() (key, glow, value lipgloss.Style) {
	glow = lipgloss.NewStyle().Foreground(a.keyGlowColor).Bold(true)
	if a.model.FooterFrame == FrameFill {
		fill := a.model.Theme.FooterFill
		return fill.Bold(true), glow.Background(a.contentBgColor), fill
	}
	return lipgloss.NewStyle().Foreground(a.footerKeyColor).Bold(true),
		glow,
		lipgloss.NewStyle().Foreground(a.footerValueColor)
}

// renderFooterContent returns the footer content without the outer footerStyle.
func (a *TabApp) renderFooterContent() string {
	currentPage := a.currentPage()
	if currentPage != nil {
		return a.renderPageFooter(currentPage)
	}
	return "Tab: Navigate • Ctrl+Q: Quit"
}
