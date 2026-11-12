package styles

import (
	"fmt"
	"image/color"

	"charm.land/lipgloss/v2"
)

// isSet reports whether c holds a color. A nil color and lipgloss.NoColor
// both mean "not set".
func isSet(c color.Color) bool {
	if c == nil {
		return false
	}
	_, none := c.(lipgloss.NoColor)
	return !none
}

// colorOr returns c if it is set, otherwise the fallback.
// This enables optional palette fields with documented defaults.
func colorOr(c, fallback color.Color) color.Color {
	if !isSet(c) {
		return fallback
	}
	return c
}

// Hex returns c as "#rrggbb", or "" for a color that is not set. Use it to
// show or compare a palette color; a color.Color prints as its channels.
func Hex(c color.Color) string {
	if !isSet(c) {
		return ""
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}
