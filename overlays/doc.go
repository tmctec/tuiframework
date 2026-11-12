// Package overlays provides the transient components that float above a
// page: Overlay is the base, Notification is a timed message styled by its
// level, and Modal is a blocking dialog that follows the theme.
//
// An overlay sizes itself as a percentage of the screen or in fixed cells,
// anchors to one of nine positions, and must have an exit key. Overlays
// stack by z-index; tabapp.TabApp draws them over the page on a Canvas, the
// grid of cells this package composites a frame with. A cell keeps the
// colors, attributes, underline and hyperlink of the text drawn on it.
//
//	n := overlays.Success("Saved").SetPosition(overlays.PositionBottomRight)
//	cmd := app.ShowNotification(n) // dismisses itself after its timeout
//
//	m := overlays.NewModal().SetTitle("Confirm").SetExitKey("esc").
//	    RegisterKey("y", DeleteMsg{})
//	id, err := app.ShowModal(m)
//
// NotificationConfig holds the defaults new notifications start from; a
// TabApp keeps its colors and icons in step with the active theme.
package overlays
