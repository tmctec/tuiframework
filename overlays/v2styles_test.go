package overlays

import (
	"image/color"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// readBack parses a frame with Lip Gloss's own canvas, so a test compares
// what a terminal built on the same libraries sees, not what this package's
// parser believes it wrote.
func readBack(frame string, width, height int) *lipgloss.Canvas {
	return lipgloss.NewCanvas(width, height).Compose(lipgloss.NewLayer(frame))
}

// sameCells fails unless the two frames hold the same content, style and
// hyperlink in every cell.
func sameCells(t *testing.T, got, want string, width, height int) {
	t.Helper()
	g, w := readBack(got, width, height), readBack(want, width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			gc, wc := g.CellAt(x, y), w.CellAt(x, y)
			if (gc == nil) != (wc == nil) {
				t.Fatalf("cell %d,%d: got %v, want %v", x, y, gc, wc)
			}
			if gc != nil && !gc.Equal(wc) {
				t.Errorf("cell %d,%d = %+v, want %+v", x, y, *gc, *wc)
			}
		}
	}
}

// composite draws the layers on a canvas without fill colors and returns the
// frame.
func composite(width, height int, layers ...layer) string {
	c := NewCanvas().Reset(width, height, nil)
	for _, l := range layers {
		lines, w := GetLines(l.text)
		c.Draw(lines, l.x, l.y, w, len(lines), l.fill)
	}
	return c.String()
}

type layer struct {
	text string
	x, y int
	fill color.Color
}

const exampleURL = "https://example.com/docs"

// cv-1: page content with a hyperlink keeps it on exactly the link's cells.
func TestCanvasKeepsHyperlinks(t *testing.T) {
	link := lipgloss.NewStyle().Hyperlink(exampleURL, "id=doc").Foreground(lipgloss.Color("#00ff00")).Render("docs")
	source := "see " + link + " now"
	width := lipgloss.Width(source)

	frame := composite(width, 1, layer{text: source})
	sameCells(t, frame, source, width, 1)

	canvas := readBack(frame, width, 1)
	for x := 0; x < width; x++ {
		url := canvas.CellAt(x, 0).Link.URL
		inLink := x >= 4 && x < 8
		if inLink && url != exampleURL {
			t.Errorf("cell %d link = %q, want %q", x, url, exampleURL)
		}
		if !inLink && url != "" {
			t.Errorf("cell %d outside the link carries %q", x, url)
		}
	}
	if canvas.CellAt(4, 0).Link.Params != "id=doc" {
		t.Errorf("link params = %q, want id=doc", canvas.CellAt(4, 0).Link.Params)
	}
}

// cv-2: an overlay drawn across the middle of a hyperlink ends the link at
// its left edge, carries none itself, and lets the link resume after it.
func TestCanvasOverlayCutsAHyperlink(t *testing.T) {
	source := lipgloss.NewStyle().Hyperlink(exampleURL).Render("0123456789")
	frame := composite(10, 1,
		layer{text: source},
		layer{text: "XX", x: 4, fill: lipgloss.Color("#202020")},
	)

	canvas := readBack(frame, 10, 1)
	for x := 0; x < 10; x++ {
		cell := canvas.CellAt(x, 0)
		covered := x == 4 || x == 5
		switch {
		case covered && (cell.Link.URL != "" || cell.Content != "X"):
			t.Errorf("cell %d = %q with link %q, want X without a link", x, cell.Content, cell.Link.URL)
		case !covered && cell.Link.URL != exampleURL:
			t.Errorf("cell %d link = %q, want the link on both sides of the overlay", x, cell.Link.URL)
		}
	}
	if strings.Count(frame, hyperlinkPrefix+";\a") != 2 {
		t.Errorf("frame %q should close the link before the overlay and at the end", frame)
	}
}

// cv-3: underline shapes and the underline color survive compositing, and so
// do the attributes beside bold.
func TestCanvasKeepsUnderlinesAndAttributes(t *testing.T) {
	red := lipgloss.Color("#ff0000")
	sources := map[string]string{
		"curly with color": lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineCurly).UnderlineColor(red).Render("wavy"),
		"single":           lipgloss.NewStyle().Underline(true).Render("under"),
		"double":           lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDouble).Render("twice"),
		"dotted and dashed": lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDotted).Render("dots") +
			lipgloss.NewStyle().UnderlineStyle(lipgloss.UnderlineDashed).Render("dash"),
		"attributes": lipgloss.NewStyle().Italic(true).Faint(true).Strikethrough(true).Reverse(true).Blink(true).Render("all"),
		"mixed": lipgloss.NewStyle().Bold(true).Render("b") + lipgloss.NewStyle().Italic(true).Render("i") +
			lipgloss.NewStyle().Underline(true).Foreground(red).Render("u") + "plain",
	}
	for name, source := range sources {
		t.Run(name, func(t *testing.T) {
			width := lipgloss.Width(source)
			sameCells(t, composite(width, 1, layer{text: source}), source, width, 1)
		})
	}

	cells := ParseCells(sources["curly with color"], 4)
	if cells[0].Underline != UnderlineCurly || cells[0].UnderlineColor != "58;2;255;0;0" {
		t.Errorf("cell = %+v, want a curly underline in red", cells[0])
	}
}

// cv-4: plain styled text gains no link and no underline sequence.
func TestCanvasAddsNoLinkOrUnderline(t *testing.T) {
	source := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00ff00")).Render("plain") + " text"
	frame := composite(lipgloss.Width(source), 1, layer{text: source})

	if strings.Contains(frame, hyperlinkPrefix) {
		t.Errorf("frame %q carries a hyperlink sequence", frame)
	}
	for _, cell := range ParseCells(frame, lipgloss.Width(source)) {
		if cell.Underline != UnderlineNone || cell.UnderlineColor != "" || cell.Link != "" {
			t.Errorf("cell %+v gained an underline or a link", cell)
		}
	}
	sameCells(t, frame, source, lipgloss.Width(source), 1)
}

// cv-6: a modal with a blended border keeps each border cell's own color
// when it is drawn over page content.
func TestCanvasKeepsBlendedBorders(t *testing.T) {
	base := strings.TrimSuffix(strings.Repeat(strings.Repeat(".", 12)+"\n", 5), "\n")
	modal := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForegroundBlend(lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")).
		Width(8).Height(3).
		Render("hi")
	lines, width := GetLines(modal)

	frame := composite(12, 5, layer{text: base}, layer{text: modal, x: 2, y: 1})
	got, want := readBack(frame, 12, 5), readBack(modal, width, len(lines))

	colors := map[color.Color]bool{}
	for y := range lines {
		for x := 0; x < width; x++ {
			g, w := got.CellAt(x+2, y+1), want.CellAt(x, y)
			if !g.Equal(w) {
				t.Errorf("modal cell %d,%d = %+v, want %+v", x, y, *g, *w)
			}
			if y == 0 && w.Style.Fg != nil {
				colors[w.Style.Fg] = true
			}
		}
	}
	if len(colors) < 3 {
		t.Errorf("the top border has %d colors, want a gradient", len(colors))
	}
	if got.CellAt(0, 0).Content != "." || got.CellAt(11, 4).Content != "." {
		t.Error("the page content outside the modal should stay")
	}
}

// TestParseCellsReadsEveryAttribute covers the SGR parameters the parser
// tracks, in the forms terminals accept them.
func TestParseCellsReadsEveryAttribute(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want Cell
	}{
		{"faint italic", "\x1b[2;3mx", Cell{Char: 'x', Attrs: AttrFaint | AttrItalic}},
		{"blink, both forms", "\x1b[5;6mx", Cell{Char: 'x', Attrs: AttrBlink}},
		{"reverse conceal strike", "\x1b[7;8;9mx", Cell{Char: 'x', Attrs: AttrReverse | AttrConceal | AttrStrikethrough}},
		{"each attribute turns off", "\x1b[1;2;3;5;7;8;9m\x1b[22;23;25;27;28;29mx", Cell{Char: 'x'}},
		{"double underline, legacy form", "\x1b[21mx", Cell{Char: 'x', Underline: UnderlineDouble}},
		{"underline off", "\x1b[4:3m\x1b[24mx", Cell{Char: 'x'}},
		{"underline shape 4:0 is none", "\x1b[4m\x1b[4:0mx", Cell{Char: 'x'}},
		{"underline shape out of range is ignored", "\x1b[4:9mx", Cell{Char: 'x'}},
		{"underline color and its reset", "\x1b[58;5;9m\x1b[59mx", Cell{Char: 'x'}},
		{"palette underline color", "\x1b[58;5;9mx", Cell{Char: 'x', UnderlineColor: "58;5;9"}},
		{"colon colors", "\x1b[38:2::1:2:3;48:5:7;58:2::4:5:6mx", Cell{Char: 'x', FgColor: "38:2::1:2:3", BgColor: "48:5:7", UnderlineColor: "58:2::4:5:6"}},
		{"a color with an empty component sets nothing", "\x1b[38;2;;1;2mx", Cell{Char: 'x'}},
		{"an unknown colon parameter is ignored", "\x1b[99:1mx", Cell{Char: 'x'}},
		{"an empty parameter resets", "\x1b[1m\x1b[;3mx", Cell{Char: 'x', Attrs: AttrItalic}},
		{"a parameter that is no number is ignored", "\x1b[1;?;1234mx", Cell{Char: 'x', Attrs: AttrBold}},
		{"reset keeps the hyperlink", "\x1b]8;;u\a\x1b[1m\x1b[0mx", Cell{Char: 'x', Link: ";u"}},
		{"link ended by the string terminator", "\x1b]8;id=a;u\x1b\\x", Cell{Char: 'x', Link: "id=a;u"}},
		{"closing link with parameters", "\x1b]8;;u\a\x1b]8;id=a;\ax", Cell{Char: 'x'}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCells(tt.in, 1)[0]; got != tt.want {
				t.Errorf("cell = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestWriteCellsStyleTransitions covers what the writer emits between
// neighbouring cells: only what changes, and a reset before anything turns
// off.
func TestWriteCellsStyleTransitions(t *testing.T) {
	tests := []struct {
		name  string
		cells []Cell
		want  string
	}{
		{
			"attributes add without a reset",
			[]Cell{{Char: 'a', Attrs: AttrBold}, {Char: 'b', Attrs: AttrBold | AttrItalic}},
			"\x1b[1ma\x1b[3mb\x1b[0m",
		},
		{
			"an attribute that ends resets and restates the rest",
			[]Cell{{Char: 'a', Attrs: AttrBold | AttrItalic, FgColor: "31"}, {Char: 'b', Attrs: AttrItalic, FgColor: "31"}},
			"\x1b[1;3;31ma\x1b[0m\x1b[3;31mb\x1b[0m",
		},
		{
			"an underline that ends resets",
			[]Cell{{Char: 'a', Underline: UnderlineCurly, UnderlineColor: "58;5;9"}, {Char: 'b'}},
			"\x1b[4:3;58;5;9ma\x1b[0mb\x1b[0m",
		},
		{
			"an underline changes shape in place",
			[]Cell{{Char: 'a', Underline: UnderlineSingle}, {Char: 'b', Underline: UnderlineDashed}},
			"\x1b[4ma\x1b[4:5mb\x1b[0m",
		},
		{
			"an underline color that ends resets and keeps the underline",
			[]Cell{{Char: 'a', Underline: UnderlineSingle, UnderlineColor: "58;5;9"}, {Char: 'b', Underline: UnderlineSingle}},
			"\x1b[4;58;5;9ma\x1b[0m\x1b[4mb\x1b[0m",
		},
		{
			"a link opens, changes and closes around the style",
			[]Cell{{Char: 'a', Link: ";u"}, {Char: 'b', Link: ";v", Attrs: AttrBold}, {Char: 'c'}},
			"\x1b]8;;u\aa\x1b]8;;v\a\x1b[1mb\x1b]8;;\a\x1b[0mc\x1b[0m",
		},
		{
			"a link open at the end of the row is closed",
			[]Cell{{Char: 'a', Link: ";u"}},
			"\x1b]8;;u\aa\x1b]8;;\a\x1b[0m",
		},
		{
			"an underline shape the writer does not know is left out",
			[]Cell{{Char: 'a', Underline: Underline(9)}},
			"a\x1b[0m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := renderCells(tt.cells); got != tt.want {
				t.Errorf("renderCells = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestNextUnitNeverStalls feeds the unit reader bytes that form no grapheme
// cluster; every call still consumes at least one byte.
func TestNextUnitNeverStalls(t *testing.T) {
	for _, s := range []string{"\xff", "\xc3", "a\xff\xfeb", "\x80\x80"} {
		for i := 0; i < len(s); {
			unit, _, _ := nextUnit(s, i)
			if unit == "" {
				t.Fatalf("nextUnit(%q, %d) consumed nothing", s, i)
			}
			i += len(unit)
		}
		if w := lineWidth(s); w < 0 {
			t.Errorf("lineWidth(%q) = %d", s, w)
		}
	}
}
