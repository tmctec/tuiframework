package overlays

import (
	"image/color"
	"strings"
)

// Canvas is a frame of cells that layers of rendered text are drawn onto,
// one rectangle each: the bands of a screen, then the overlays above them.
// A layer replaces what was under it and fills the cells its text leaves
// without a background, so nothing shows through. The frame becomes text
// once, in String, however many layers it took.
//
// A Canvas keeps its cells between frames, so drawing one allocates nothing
// but the final string. Reset starts the next frame.
type Canvas struct {
	width, height int
	cells         []Cell // row after row, width cells each
	size          int    // length of the last frame's text, to size the next
}

// NewCanvas creates an empty canvas. Reset gives it a size.
func NewCanvas() *Canvas {
	return &Canvas{}
}

// Reset starts a frame: it sizes the canvas to width by height cells and
// clears every cell to a space on the color bg. Returns self for chaining.
func (c *Canvas) Reset(width, height int, bg color.Color) *Canvas {
	c.width, c.height = max(width, 0), max(height, 0)
	n := c.width * c.height
	if cap(c.cells) < n {
		c.cells = make([]Cell, n)
	}
	c.cells = c.cells[:n]

	blank := blankCell(bg)
	for i := range c.cells {
		c.cells[i] = blank
	}
	return c
}

// Width returns the width of the canvas in cells.
func (c *Canvas) Width() int {
	return c.width
}

// Height returns the height of the canvas in lines.
func (c *Canvas) Height() int {
	return c.height
}

// Draw draws a layer: lines of rendered text in the rectangle width cells
// wide and height lines tall whose top left corner is at column x, line y.
// A negative x or y counts as 0. Text outside the rectangle is dropped, and
// so is the part of the rectangle outside the canvas. Every cell of the
// rectangle is replaced: where a line is short or missing by a space, and
// where the text has no background color by the color bg. Returns self for
// chaining.
func (c *Canvas) Draw(lines []string, x, y, width, height int, bg color.Color) *Canvas {
	x, y = max(x, 0), max(y, 0)
	width = min(width, c.width-x)
	height = min(height, c.height-y)
	if width <= 0 || height <= 0 {
		return c
	}

	blank := blankCell(bg)
	for i := 0; i < height; i++ {
		row := c.cells[(y+i)*c.width : (y+i+1)*c.width]

		// A wide character the rectangle cuts in two leaves a space for the
		// half that stays visible.
		if x > 0 && row[x].Char == 0 {
			row[x-1].Char, row[x-1].rest = ' ', ""
		}
		if end := x + width; end < c.width && row[end].Char == 0 {
			row[end].Char = ' '
		}

		line := ""
		if i < len(lines) {
			line = lines[i]
		}
		parseCellsInto(row[x:x+width], line, blank)
	}
	return c
}

// String returns the frame as text: one line per row, joined by newlines,
// every line exactly as wide as the canvas.
func (c *Canvas) String() string {
	var b strings.Builder
	b.Grow(c.size + c.size/8)
	for y := 0; y < c.height; y++ {
		if y > 0 {
			b.WriteByte('\n')
		}
		writeCells(&b, c.cells[y*c.width:(y+1)*c.width])
	}
	c.size = b.Len()
	return b.String()
}
