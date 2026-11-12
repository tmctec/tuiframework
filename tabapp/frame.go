package tabapp

import "charm.land/lipgloss/v2"

// Frame selects how the header or footer band is set off from the content.
type Frame int

const (
	// FrameNone draws the band bare (default).
	FrameNone Frame = iota
	// FrameTheme draws a border in the theme's shape, the palette's Border.
	FrameTheme
	// FrameSingle draws a single-line border.
	FrameSingle
	// FrameDouble draws a double-line border.
	FrameDouble
	// FrameFill paints the band as a solid bar in the theme's fill colors,
	// Theme.HeaderFill or Theme.FooterFill, with the text reversed to match.
	FrameFill
)

// String returns the frame's name.
func (f Frame) String() string {
	switch f {
	case FrameTheme:
		return "Theme"
	case FrameSingle:
		return "Single"
	case FrameDouble:
		return "Double"
	case FrameFill:
		return "Fill"
	default:
		return "None"
	}
}

// themeFrame returns the frame of a band whose border is switched on or off:
// the theme's border or none.
func themeFrame(show bool) Frame {
	if show {
		return FrameTheme
	}
	return FrameNone
}

// bordered reports whether the frame draws a border, which takes a row above
// and below the band.
func (f Frame) bordered() bool {
	return f == FrameTheme || f == FrameSingle || f == FrameDouble
}

// columns returns the columns the frame takes from the band's width: one on
// each side for a border, and the same inset for a fill so the text keeps
// clear of the screen edge.
func (f Frame) columns() int {
	if f.bordered() || f == FrameFill {
		return 2
	}
	return 0
}

// border returns the border the frame draws; themed is the shape FrameTheme
// uses.
func (f Frame) border(themed lipgloss.Border) lipgloss.Border {
	switch f {
	case FrameSingle:
		return lipgloss.NormalBorder()
	case FrameDouble:
		return lipgloss.DoubleBorder()
	default:
		return themed
	}
}
