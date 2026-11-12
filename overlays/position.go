package overlays

// Position defines where an overlay is anchored on screen.
// Uses two-axis justification: Horizontal (L/C/R) + Vertical (T/C/B).
type Position string

// The nine anchors: horizontal L, C or R followed by vertical T, C or B.
const (
	PositionTopLeft      Position = "LT"
	PositionTopCenter    Position = "CT"
	PositionTopRight     Position = "RT"
	PositionCenterLeft   Position = "LC"
	PositionCenter       Position = "CC" // Default
	PositionCenterRight  Position = "RC"
	PositionBottomLeft   Position = "LB"
	PositionBottomCenter Position = "CB"
	PositionBottomRight  Position = "RB"
)

// CalculatePosition returns the x, y offset for placing an overlay of the given
// dimensions within a container of the given size.
func CalculatePosition(pos Position, overlayWidth, overlayHeight, containerWidth, containerHeight int) (x, y int) {
	if len(pos) != 2 {
		pos = PositionCenter
	}
	// Horizontal position
	switch pos[0] {
	case 'L':
		x = 0
	case 'C':
		x = (containerWidth - overlayWidth) / 2
	case 'R':
		x = containerWidth - overlayWidth
	}

	// Vertical position
	switch pos[1] {
	case 'T':
		y = 0
	case 'C':
		y = (containerHeight - overlayHeight) / 2
	case 'B':
		y = containerHeight - overlayHeight
	}

	// Ensure non-negative
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}

	return x, y
}
