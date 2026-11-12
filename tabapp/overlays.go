package tabapp

import (
	"fmt"
	"slices"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
)

// overlayEntry is an overlay on the stack. Its z-index is the overlay's own.
type overlayEntry struct {
	overlay *overlays.Overlay
	modal   *overlays.Modal // Set when shown through ShowModal, so theme changes can reach it
	id      int64
}

// ShowOverlay displays an overlay on top of the current page.
// Overlays are stacked by z-index (CSS-style: higher = on top).
// Returns the overlay's unique ID for later removal. An overlay without an
// exit key fails with overlays.ErrNoExitKey, and one whose keys include a
// hardcoded quit key with an error that wraps overlays.ErrKeyConflict.
func (a *TabApp) ShowOverlay(o *overlays.Overlay) (int64, error) {
	if err := o.Validate(); err != nil {
		return 0, err
	}
	for key := range o.Keys() {
		if _, isHardcoded := a.hardcodedKeys[key]; isHardcoded {
			return 0, fmt.Errorf("%w: key %q", overlays.ErrKeyConflict, key)
		}
	}
	if _, isHardcoded := a.hardcodedKeys[o.ExitKey()]; isHardcoded {
		return 0, fmt.Errorf("%w: exit key %q", overlays.ErrKeyConflict, o.ExitKey())
	}

	a.overlaySeq++
	id := a.overlaySeq
	a.insertOverlayByZIndex(overlayEntry{overlay: o, id: id})
	o.Show()
	a.syncOverlays()
	return id, nil
}

// syncOverlays mirrors the overlay stack into the Model: the count, and the
// OpenModals registry, which drops a modal however it closed (exit key,
// timeout, or Hide*).
func (a *TabApp) syncOverlays() {
	a.model.OverlayCount = len(a.overlayStack)
	for title, id := range a.model.OpenModals {
		if !a.hasOverlayID(id) {
			delete(a.model.OpenModals, title)
		}
	}
}

// hasOverlayID reports whether an overlay with this id is on the stack.
func (a *TabApp) hasOverlayID(id int64) bool {
	for _, e := range a.overlayStack {
		if e.id == id {
			return true
		}
	}
	return false
}

// insertOverlayByZIndex inserts an overlay entry maintaining z-index sort order.
// Same z-index: insertion order (last added = on top of peers).
func (a *TabApp) insertOverlayByZIndex(entry overlayEntry) {
	// After every entry whose z-index is not higher.
	z := entry.overlay.ZIndex()
	pos := len(a.overlayStack)
	for i, e := range a.overlayStack {
		if e.overlay.ZIndex() > z {
			pos = i
			break
		}
	}
	a.overlayStack = slices.Insert(a.overlayStack, pos, entry)
}

// MustShowOverlay displays an overlay, panicking on error.
// Returns the overlay's unique ID.
func (a *TabApp) MustShowOverlay(o *overlays.Overlay) int64 {
	id, err := a.ShowOverlay(o)
	if err != nil {
		panic(err)
	}
	return id
}

// ShowNotification displays a notification overlay.
// Returns a tea.Cmd that handles auto-dismiss if timeout is configured.
// The caller should return this command from their key handler.
func (a *TabApp) ShowNotification(n *overlays.Notification) tea.Cmd {
	id, err := a.ShowOverlay(n.Overlay)
	if err != nil {
		return nil
	}

	// If timeout is configured, return a command that will dismiss after duration
	if n.Timeout() > 0 {
		return tea.Tick(n.Timeout(), func(_ time.Time) tea.Msg {
			return overlayTimeoutMsg{overlayID: id}
		})
	}
	return nil
}

// ShowModal displays a modal overlay.
// A modal styled by ApplyTheme keeps following the theme while it is shown.
// Returns the overlay's unique ID for later removal, or error if invalid.
func (a *TabApp) ShowModal(m *overlays.Modal) (int64, error) {
	id, err := a.ShowOverlay(m.Overlay)
	if err != nil {
		return 0, err
	}
	for i := range a.overlayStack {
		if a.overlayStack[i].id == id {
			a.overlayStack[i].modal = m
		}
	}
	return id, nil
}

// HideOverlay dismisses the top overlay, the one with the highest z-index.
func (a *TabApp) HideOverlay() {
	if len(a.overlayStack) > 0 {
		top := a.overlayStack[len(a.overlayStack)-1]
		top.overlay.Hide()
		a.overlayStack = a.overlayStack[:len(a.overlayStack)-1]
		a.syncOverlays()
	}
}

// HideOverlayByID removes a specific overlay by its unique ID.
// Returns true if found and removed, false if not found.
func (a *TabApp) HideOverlayByID(id int64) bool {
	for i, entry := range a.overlayStack {
		if entry.id == id {
			entry.overlay.Hide()
			a.overlayStack = slices.Delete(a.overlayStack, i, i+1)
			a.syncOverlays()
			return true
		}
	}
	return false
}

// HideAllOverlays removes all overlays from the stack.
func (a *TabApp) HideAllOverlays() {
	for _, entry := range a.overlayStack {
		entry.overlay.Hide()
	}
	a.overlayStack = nil
	a.syncOverlays()
}

// HasOverlay returns true if any overlay is currently displayed.
func (a *TabApp) HasOverlay() bool {
	return len(a.overlayStack) > 0
}

// ActiveOverlay returns the top overlay (highest z-index), or nil if none.
//
// Deprecated: use TopOverlay.
func (a *TabApp) ActiveOverlay() *overlays.Overlay {
	return a.TopOverlay()
}

// TopOverlay returns the overlay with highest z-index, or nil if none.
func (a *TabApp) TopOverlay() *overlays.Overlay {
	if len(a.overlayStack) == 0 {
		return nil
	}
	return a.overlayStack[len(a.overlayStack)-1].overlay
}

// OverlayCount returns the number of overlays in the stack.
func (a *TabApp) OverlayCount() int {
	return len(a.overlayStack)
}

// handleOverlayKeys processes a key press against the overlay stack.
// Overlays are checked from top (highest z-index) to bottom.
// Returns (true, cmd) if a key was handled, (false, nil) otherwise.
func (a *TabApp) handleOverlayKeys(key string) (bool, tea.Cmd) {
	// Process from highest z-index to lowest (reverse iteration)
	for i := len(a.overlayStack) - 1; i >= 0; i-- {
		entry := a.overlayStack[i]
		overlay := entry.overlay

		// Check if this is the overlay's exit key
		if key == overlay.ExitKey() {
			a.HideOverlayByID(entry.id)
			return true, nil
		}

		// Check overlay's registered keys (Elm pattern: dispatch message)
		if msg, ok := overlay.Keys()[key]; ok {
			return true, func() tea.Msg { return msg }
		}

		// If this overlay is blocking, stop here (don't pass to lower overlays or page)
		if overlay.IsBlocking() {
			return true, nil
		}
	}
	return false, nil
}

// BringToFront moves an overlay to the top of the stack.
// Sets its z-index to max+1 of current stack.
func (a *TabApp) BringToFront(id int64) {
	// Find the overlay
	idx := -1
	for i, entry := range a.overlayStack {
		if entry.id == id {
			idx = i
			break
		}
	}
	if idx < 0 {
		return // Not found
	}

	// Calculate max z-index
	maxZ := 0
	for _, entry := range a.overlayStack {
		maxZ = max(maxZ, entry.overlay.ZIndex())
	}

	// Move the entry to the end with the highest z-index
	entry := a.overlayStack[idx]
	a.overlayStack = slices.Delete(a.overlayStack, idx, idx+1)
	entry.overlay.SetZIndex(maxZ + 1)
	a.overlayStack = append(a.overlayStack, entry)
}

// textModal creates a modal in the active theme with a title and text
// content, closed by esc.
func (a *TabApp) textModal(title, content string) *overlays.Modal {
	modal := overlays.NewModal().
		ApplyTheme(a.model.Theme).
		SetTitle(title).
		SetExitKey("esc")
	modal.Overlay.SetContent(content)
	return modal
}

// createModalFromMsg creates a modal from ToggleModalMsg parameters.
// The modal takes the active theme; a BorderColor that is set overrides the
// theme's modal border color.
func (a *TabApp) createModalFromMsg(msg ToggleModalMsg) *overlays.Modal {
	modal := a.textModal(msg.Title, msg.Content).
		SetSize(msg.SizePct).
		SetZIndex(msg.ZIndex).
		SetBlocking(false)
	if msg.BorderColor != nil {
		modal.SetBorderColor(msg.BorderColor)
	}

	return modal
}
