package overlays

import (
	"image/color"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/styles"
)

// Modal is an overlay for user interaction - input, confirmations, selections.
// The title line lives on the embedded Overlay, so it renders wherever the
// Overlay renders.
type Modal struct {
	*Overlay

	followsTheme bool        // frame set by ApplyTheme; TabApp re-applies on theme change
	borderColor  color.Color // SetBorderColor override, kept across ApplyTheme; nil without one
}

// NewModal creates a new modal overlay.
// Modals are blocking by default and carry the Dracula theme until
// ApplyTheme replaces it.
func NewModal() *Modal {
	m := &Modal{
		Overlay: New(),
	}
	m.Overlay.SetBlocking(true)
	return m.ApplyTheme(styles.Dracula)
}

// ApplyTheme styles the frame, background and title from the theme's Modal*
// fields and gives the body the theme's primary text color.
func (m *Modal) ApplyTheme(t *styles.Theme) *Modal {
	m.Overlay.SetStyle(lipgloss.NewStyle().
		Background(t.ModalBackground).
		Foreground(t.Palette.Fg1).
		Border(t.ModalBorder).
		BorderForeground(t.ModalBorderColor).
		BorderBackground(t.ModalBackground).
		Padding(1, 2))
	m.SetTitleStyle(t.ModalTitle)
	if len(t.BorderBlend) >= 2 {
		m.SetBorderBlend(t.BorderBlend...)
	}
	if m.borderColor != nil {
		m.SetBorderForeground(m.borderColor)
	}
	m.followsTheme = true
	return m
}

// SetBorderColor sets the border color. The color survives ApplyTheme, so a
// modal can follow the theme and still keep its own border color.
func (m *Modal) SetBorderColor(c color.Color) *Modal {
	m.borderColor = c
	m.SetBorderForeground(c)
	return m
}

// FollowsTheme reports whether the frame comes from ApplyTheme rather than
// SetStyle. TabApp re-applies the active theme to such modals when it changes.
func (m *Modal) FollowsTheme() bool {
	return m.followsTheme
}

// SetTitle sets the modal title.
func (m *Modal) SetTitle(title string) *Modal {
	m.Overlay.SetTitle(title)
	return m
}

// SetContent sets the modal content via a builder function.
// The function receives a flexbox to populate.
func (m *Modal) SetContent(fn func(fb *flexbox.FlexBox)) *Modal {
	m.SetContentFunc(fn)
	return m
}

// SetSize sets the modal size percentage.
func (m *Modal) SetSize(percent int) *Modal {
	m.Overlay.SetSize(percent)
	return m
}

// SetPosition sets where the modal appears.
func (m *Modal) SetPosition(pos Position) *Modal {
	m.Overlay.SetPosition(pos)
	return m
}

// SetExitKey sets the key to dismiss the modal.
func (m *Modal) SetExitKey(key string) *Modal {
	m.Overlay.SetExitKey(key)
	return m
}

// SetBlocking sets whether the modal blocks page keys.
func (m *Modal) SetBlocking(blocking bool) *Modal {
	m.Overlay.SetBlocking(blocking)
	return m
}

// SetStyle sets the modal style. A modal with its own style no longer
// follows theme changes.
func (m *Modal) SetStyle(style lipgloss.Style) *Modal {
	m.Overlay.SetStyle(style)
	m.followsTheme = false
	return m
}

// RegisterKey registers a message to dispatch when key is pressed (Elm pattern).
func (m *Modal) RegisterKey(key string, msg tea.Msg) *Modal {
	m.Overlay.RegisterKey(key, msg)
	return m
}

// RegisterKeys registers multiple messages to dispatch when keys are pressed (Elm pattern).
func (m *Modal) RegisterKeys(bindings map[string]tea.Msg) *Modal {
	m.Overlay.RegisterKeys(bindings)
	return m
}

// SetZIndex sets the z-axis stacking order (CSS-style).
// Higher values render on top. Default is 0.
func (m *Modal) SetZIndex(z int) *Modal {
	m.Overlay.SetZIndex(z)
	return m
}
