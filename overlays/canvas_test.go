package overlays

import (
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// plainRows returns a canvas's rows without escape sequences.
func plainRows(c *Canvas) []string {
	rows := strings.Split(c.String(), "\n")
	for i, r := range rows {
		rows[i] = stripANSI(r)
	}
	return rows
}

// stripANSI removes the escape sequences from s.
func stripANSI(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		unit, _, esc := nextUnit(s, i)
		i += len(unit)
		if !esc {
			b.WriteString(unit)
		}
	}
	return b.String()
}

func TestCanvasResetSizesAndClears(t *testing.T) {
	c := NewCanvas()
	if c.Width() != 0 || c.Height() != 0 || c.String() != "" {
		t.Errorf("a new canvas is %dx%d %q, want empty", c.Width(), c.Height(), c.String())
	}

	if got := c.Reset(4, 2, nil); got != c {
		t.Error("Reset should return the canvas for chaining")
	}
	if c.Width() != 4 || c.Height() != 2 {
		t.Errorf("size = %dx%d, want 4x2", c.Width(), c.Height())
	}
	if rows := plainRows(c); len(rows) != 2 || rows[0] != "    " || rows[1] != "    " {
		t.Errorf("rows = %q, want two rows of four spaces", rows)
	}

	// A smaller frame after a larger one keeps nothing of the old frame.
	c.Draw([]string{"abcd", "efgh"}, 0, 0, 4, 2, nil)
	c.Reset(2, 1, nil)
	if rows := plainRows(c); len(rows) != 1 || rows[0] != "  " {
		t.Errorf("after a smaller Reset rows = %q, want one row of two spaces", rows)
	}

	c.Reset(-3, -1, nil)
	if c.Width() != 0 || c.Height() != 0 || c.String() != "" {
		t.Errorf("a negative size gives %dx%d %q, want empty", c.Width(), c.Height(), c.String())
	}
}

func TestCanvasDrawClipsAndPads(t *testing.T) {
	c := NewCanvas().Reset(6, 4, nil)
	if got := c.Draw([]string{"abcdefgh", "ij"}, 1, 1, 4, 3, nil); got != c {
		t.Error("Draw should return the canvas for chaining")
	}
	want := []string{
		"      ",
		" abcd ", // cut at the rectangle's width
		" ij   ", // a short line is padded
		"      ", // a missing line is blank
	}
	if rows := plainRows(c); strings.Join(rows, "|") != strings.Join(want, "|") {
		t.Errorf("rows = %q\nwant   %q", rows, want)
	}

	// A rectangle that leaves the canvas is cut at its edge, and a negative
	// corner counts as 0.
	c.Reset(4, 2, nil).Draw([]string{"12345", "67890", "extra"}, 2, 1, 10, 10, nil)
	if rows := plainRows(c); rows[0] != "    " || rows[1] != "  12" {
		t.Errorf("clipped rows = %q, want the layer cut at the right and bottom edges", rows)
	}
	c.Reset(4, 2, nil).Draw([]string{"ab"}, -5, -5, 2, 1, nil)
	if rows := plainRows(c); rows[0] != "ab  " {
		t.Errorf("negative corner: row 0 = %q, want the layer at column 0", rows[0])
	}

	// Rectangles with no area or outside the canvas draw nothing.
	c.Reset(4, 2, nil)
	c.Draw([]string{"x"}, 0, 0, 0, 1, nil).Draw([]string{"x"}, 0, 0, 1, 0, nil)
	c.Draw([]string{"x"}, 4, 0, 1, 1, nil).Draw([]string{"x"}, 0, 2, 1, 1, nil)
	if rows := plainRows(c); rows[0] != "    " || rows[1] != "    " {
		t.Errorf("empty rectangles changed the canvas: %q", rows)
	}
}

// TestCanvasLayersReplaceWhatIsBelow draws one layer over another: the
// upper layer's rectangle shows none of the lower one, blank cells included.
func TestCanvasLayersReplaceWhatIsBelow(t *testing.T) {
	c := NewCanvas().Reset(8, 3, nil)
	c.Draw([]string{"BBBBBBBB", "BBBBBBBB", "BBBBBBBB"}, 0, 0, 8, 3, nil)
	c.Draw([]string{"x"}, 2, 1, 4, 1, nil)

	want := []string{"BBBBBBBB", "BBx   BB", "BBBBBBBB"}
	if rows := plainRows(c); strings.Join(rows, "|") != strings.Join(want, "|") {
		t.Errorf("rows = %q\nwant   %q", rows, want)
	}
	for i, line := range strings.Split(c.String(), "\n") {
		if w := lipgloss.Width(line); w != 8 {
			t.Errorf("row %d is %d cells wide, want 8", i, w)
		}
	}
}

func TestCanvasFillsMissingBackgrounds(t *testing.T) {
	base, fill, own := lipgloss.Color("#101010"), lipgloss.Color("#202020"), lipgloss.Color("#303030")
	bgOf := func(c color.Color) string { return blankCell(c).BgColor }

	c := NewCanvas().Reset(6, 2, base)
	styled := lipgloss.NewStyle().Background(own).Render("ab") + "c"
	c.Draw([]string{styled}, 1, 0, 4, 1, fill)

	rows := strings.Split(c.String(), "\n")
	cells := ParseCells(rows[0], 6)
	want := []string{bgOf(base), bgOf(own), bgOf(own), bgOf(fill), bgOf(fill), bgOf(base)}
	for i, w := range want {
		if cells[i].BgColor != w || !cells[i].HasBackground() {
			t.Errorf("row 0 cell %d background = %q, want %q", i, cells[i].BgColor, w)
		}
	}
	for i, cell := range ParseCells(rows[1], 6) {
		if cell.BgColor != bgOf(base) {
			t.Errorf("row 1 cell %d background = %q, want the canvas color", i, cell.BgColor)
		}
	}
}

// TestCanvasWideCharacterCutByALayer covers a layer whose edge falls inside
// a wide character below it: the half that stays visible becomes a space, so
// the row keeps its width.
func TestCanvasWideCharacterCutByALayer(t *testing.T) {
	cases := []struct {
		name  string
		layer string
		x     int
		want  string
	}{
		{"covers a right half", "x", 1, " x本語"},
		{"covers a left half", "x", 2, "日x 語"},
		{"covers a right half and a left half", "xx", 1, " xx 語"},
		{"covers a whole character", "xx", 2, "日xx語"},
		{"covers the last half of the row", "x", 5, "日本 x"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCanvas().Reset(6, 1, nil).Draw([]string{"日本語"}, 0, 0, 6, 1, nil)
			c.Draw([]string{tc.layer}, tc.x, 0, len(tc.layer), 1, nil)
			if got := plainRows(c)[0]; got != tc.want {
				t.Errorf("row = %q, want %q", got, tc.want)
			}
			if w := lipgloss.Width(c.String()); w != 6 {
				t.Errorf("row is %d cells wide, want 6", w)
			}
		})
	}
}

func TestCanvasReusesItsCells(t *testing.T) {
	c := NewCanvas()
	lines := []string{strings.Repeat("a", 80)}
	c.Reset(80, 24, nil).Draw(lines, 0, 0, 80, 1, nil)
	_ = c.String()

	// A frame on a canvas that already has its cells allocates the text only.
	allocs := testing.AllocsPerRun(20, func() {
		c.Reset(80, 24, nil).Draw(lines, 0, 0, 80, 1, nil)
		_ = c.String()
	})
	if allocs > 2 {
		t.Errorf("a frame took %.0f allocations, want at most 2", allocs)
	}
}
