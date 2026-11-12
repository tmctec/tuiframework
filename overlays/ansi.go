package overlays

import (
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// Attr is a set of text attributes.
type Attr uint8

// The text attributes a cell can carry.
const (
	AttrBold Attr = 1 << iota
	AttrFaint
	AttrItalic
	AttrBlink
	AttrReverse
	AttrConceal
	AttrStrikethrough
)

// Underline is the shape of the line under a cell's character.
type Underline uint8

// The underline shapes, numbered as in the SGR parameter 4:n.
const (
	UnderlineNone Underline = iota
	UnderlineSingle
	UnderlineDouble
	UnderlineCurly
	UnderlineDotted
	UnderlineDashed
)

// Cell is one screen column of a line: a character and the styling in effect
// where it stands. A character is a whole grapheme cluster, so an accented
// letter or a joined emoji is one cell's worth of text.
//
// Colors are kept as the SGR parameters that select them, for example
// "38;2;255;255;255" or "31" for a foreground and "48;5;201" for a
// background. An empty color is the terminal's default.
type Cell struct {
	Char           rune      // The character's first rune; 0 marks the second cell of a wide character
	Attrs          Attr      // Bold, italic and the other attributes
	Underline      Underline // The underline shape
	FgColor        string    // The foreground color
	BgColor        string    // The background color
	UnderlineColor string    // The underline color, such as "58;2;255;0;0"
	Link           string    // The hyperlink as "params;uri", the payload of an OSC 8 sequence

	rest string // the character's other runes: combining marks, joiners, variation selectors
}

// HasBackground reports whether the cell has a background color of its own.
func (c Cell) HasBackground() bool {
	return c.BgColor != ""
}

// ParseCells parses a string into cells, one per screen column, tracking
// the styling so every cell carries its colors, attributes, underline and
// hyperlink. The result has exactly width cells: text beyond the width is
// dropped and short lines are padded with unstyled spaces. A wide character
// takes two cells, the second with Char 0. Combining marks stay with the
// character they follow, and widths are measured per grapheme cluster, the
// way Lip Gloss measures them.
func ParseCells(s string, width int) []Cell {
	cells := make([]Cell, max(width, 0))
	parseCellsInto(cells, s, Cell{Char: ' '})
	return cells
}

// parseCellsInto is ParseCells writing into cells, whose length is the
// width. A cell the text does not reach is the blank, and a character
// without a background color of its own takes the blank's. Callers that
// parse many lines reuse one buffer.
func parseCellsInto(cells []Cell, s string, blank Cell) {
	var st sgrState
	style := st.cell(blank) // the styling of the next character, with Char 0
	pos := 0                // the next free cell
	last := -1              // the cell of the last character, which marks that follow join

parse:
	for i := 0; i < len(s); {
		// Printable ASCII with no mark after it is most of any line.
		if c := s[i]; c >= ' ' && c < 0x7f && (i+1 == len(s) || s[i+1] < utf8.RuneSelf) {
			if pos == len(cells) {
				break
			}
			cells[pos] = style
			cells[pos].Char = rune(c)
			last = pos
			pos++
			i++
			continue
		}

		unit, w, esc := nextUnit(s, i)
		i += len(unit)
		switch {
		case esc:
			st.apply(unit)
			style = st.cell(blank)
		case w == 0:
			// A zero-width cluster is a mark or a joiner for the character
			// before it; a control character is dropped.
			if r, _ := utf8.DecodeRuneInString(unit); last >= 0 && !unicode.IsControl(r) {
				cells[last].rest += unit
			}
		case pos+w > len(cells):
			break parse // a character that does not fit ends the line
		default:
			for k := 1; k < w; k++ {
				cells[pos+k] = style // Char 0: the rest of a wide character
			}
			r, size := utf8.DecodeRuneInString(unit)
			cells[pos] = style
			cells[pos].Char, cells[pos].rest = r, unit[size:]
			last = pos
			pos += w
		}
	}

	for ; pos < len(cells); pos++ {
		cells[pos] = blank
	}
}

// nextUnit reads the unit of s that starts at byte i: an escape sequence,
// reported with esc true, or one character and the number of cells it
// takes. A character is a grapheme cluster with the width Lip Gloss gives
// it, so the canvas and the layout that fills it always agree. The locale
// does not change a width; RUNEWIDTH_EASTASIAN=1 makes ambiguous characters
// such as box-drawing lines two cells wide here as it does in Lip Gloss.
func nextUnit(s string, i int) (unit string, cells int, esc bool) {
	c := s[i]
	switch {
	case c == '\x1b':
		return escapeAt(s, i), 0, true
	case c < utf8.RuneSelf && (i+1 == len(s) || s[i+1] < utf8.RuneSelf):
		// An ASCII character with no mark after it needs no segmentation.
		if c < ' ' || c == 0x7f {
			return s[i : i+1], 0, false
		}
		return s[i : i+1], 1, false
	}
	unit, cells = ansi.FirstGraphemeCluster(s[i:], ansi.GraphemeWidth)
	return unit, cells, false
}

// escapeAt returns the escape sequence that starts at byte i of s, which
// holds ESC: a CSI sequence up to its final byte, an OSC sequence up to its
// terminator, or ESC and the byte after it.
func escapeAt(s string, i int) string {
	if i+1 >= len(s) {
		return s[i:]
	}
	j := i + 2
	switch s[i+1] {
	case '[': // CSI: parameters, then a final byte from @ to ~
		for j < len(s) && (s[j] < 0x40 || s[j] > 0x7e) {
			j++
		}
		if j < len(s) {
			j++
		}
	case ']': // OSC: ends at BEL or at the string terminator, ESC \
		for j < len(s) {
			if s[j] == '\a' {
				j++
				break
			}
			if s[j] == '\x1b' && j+1 < len(s) && s[j+1] == '\\' {
				j += 2
				break
			}
			j++
		}
	}
	return s[i:j]
}

// lineWidth returns the number of cells a line of text takes, escape
// sequences excluded.
func lineWidth(s string) int {
	width := 0
	for i := 0; i < len(s); {
		if c := s[i]; c >= ' ' && c < 0x7f && (i+1 == len(s) || s[i+1] < utf8.RuneSelf) {
			width++
			i++
			continue
		}
		unit, w, _ := nextUnit(s, i)
		i += len(unit)
		width += w
	}
	return width
}

// sgrState is the styling in effect at a point in a line.
type sgrState struct {
	fg, bg, ul string // color SGR parameters; "" when the default color applies
	attrs      Attr
	underline  Underline
	link       string // the open hyperlink, as "params;uri"
}

// cell returns a cell with the state's styling and no character. Without a
// background of its own it takes the blank's.
func (st *sgrState) cell(blank Cell) Cell {
	c := Cell{
		Attrs: st.attrs, Underline: st.underline,
		FgColor: st.fg, BgColor: st.bg, UnderlineColor: st.ul, Link: st.link,
	}
	if c.BgColor == "" {
		c.BgColor = blank.BgColor
	}
	return c
}

// reset clears the colors, attributes and underline. A hyperlink is not SGR
// state and stays open.
func (st *sgrState) reset() {
	*st = sgrState{link: st.link}
}

// hyperlinkPrefix starts an OSC 8 sequence, which opens or closes a hyperlink.
const hyperlinkPrefix = "\x1b]8;"

// apply updates the state from one escape sequence: an SGR sequence (a CSI
// sequence ending in m) or an OSC 8 hyperlink. Parameters are read one at a
// time, and a 38, 48 or 58 consumes the components of its color, so a
// component equal to 0, 1 or 22 is never mistaken for an attribute. The
// colors and the link it keeps are substrings of seq, so applying a sequence
// allocates nothing.
func (st *sgrState) apply(seq string) {
	if strings.HasPrefix(seq, hyperlinkPrefix) {
		st.link = linkPayload(seq)
		return
	}
	if !strings.HasPrefix(seq, "\x1b[") || !strings.HasSuffix(seq, "m") {
		return
	}
	params := seq[2 : len(seq)-1]
	if params == "" {
		st.reset()
		return
	}
	for params != "" {
		p, rest, _ := strings.Cut(params, ";")
		switch p {
		case "0", "":
			st.reset()
		case "1":
			st.attrs |= AttrBold
		case "2":
			st.attrs |= AttrFaint
		case "3":
			st.attrs |= AttrItalic
		case "4":
			st.underline = UnderlineSingle
		case "5", "6":
			st.attrs |= AttrBlink
		case "7":
			st.attrs |= AttrReverse
		case "8":
			st.attrs |= AttrConceal
		case "9":
			st.attrs |= AttrStrikethrough
		case "21":
			st.underline = UnderlineDouble
		case "22":
			st.attrs &^= AttrBold | AttrFaint
		case "23":
			st.attrs &^= AttrItalic
		case "24":
			st.underline = UnderlineNone
		case "25":
			st.attrs &^= AttrBlink
		case "27":
			st.attrs &^= AttrReverse
		case "28":
			st.attrs &^= AttrConceal
		case "29":
			st.attrs &^= AttrStrikethrough
		case "39":
			st.fg = ""
		case "49":
			st.bg = ""
		case "59":
			st.ul = ""
		case "38", "48", "58":
			n := extendedColorLen(rest)
			if n == 0 {
				return // malformed color: ignore the rest of the sequence
			}
			st.setColor(p, params[:len(p)+1+n])
			rest = strings.TrimPrefix(rest[n:], ";")
		default:
			st.applyOther(p)
		}
		params = rest
	}
}

// setColor stores a color given by its SGR parameters under the kind that
// starts them: 38 foreground, 48 background, 58 underline.
func (st *sgrState) setColor(kind, params string) {
	switch kind {
	case "38":
		st.fg = params
	case "48":
		st.bg = params
	case "58":
		st.ul = params
	}
}

// applyOther handles the parameters apply has no case for: the basic colors
// 30 to 37, 40 to 47, 90 to 97 and 100 to 107, and the parameters written
// with colons, an underline shape such as 4:3 or a color such as 38:2::1:2:3.
func (st *sgrState) applyOther(p string) {
	if kind, sub, colon := strings.Cut(p, ":"); colon {
		switch kind {
		case "4":
			if len(sub) == 1 && sub[0] >= '0' && sub[0] <= '5' {
				st.underline = Underline(sub[0] - '0')
			}
		case "38", "48", "58":
			st.setColor(kind, p)
		}
		return
	}

	code := 0
	for i := 0; i < len(p); i++ {
		if p[i] < '0' || p[i] > '9' || i == 3 {
			return
		}
		code = code*10 + int(p[i]-'0')
	}
	switch {
	case (code >= 30 && code <= 37) || (code >= 90 && code <= 97):
		st.fg = p
	case (code >= 40 && code <= 47) || (code >= 100 && code <= 107):
		st.bg = p
	}
}

// linkPayload returns the "params;uri" of an OSC 8 sequence, or "" for the
// sequence that closes a hyperlink, which has no URI.
func linkPayload(seq string) string {
	p := strings.TrimPrefix(seq, hyperlinkPrefix)
	p = strings.TrimSuffix(p, "\a")
	p = strings.TrimSuffix(p, "\x1b\\")
	if _, uri, _ := strings.Cut(p, ";"); uri == "" {
		return ""
	}
	return p
}

// extendedColorLen returns the length of the color components at the start
// of s, the parameters that follow a 38, 48 or 58: "2;r;g;b" for a true
// color and "5;n" for a palette color. It returns 0 when they do not form a
// color.
func extendedColorLen(s string) int {
	kind, _, _ := strings.Cut(s, ";")
	components := 0
	switch kind {
	case "2":
		components = 3
	case "5":
		components = 1
	default:
		return 0
	}
	n := len(kind)
	for ; components > 0; components-- {
		if n >= len(s) || s[n] != ';' {
			return 0
		}
		n++
		start := n
		for n < len(s) && s[n] != ';' {
			n++
		}
		if n == start {
			return 0
		}
	}
	return n
}

// blankCell returns an empty cell on the background bg: a space that
// carries the color's SGR parameters as Lip Gloss renders them. A nil color
// or lipgloss.NoColor gives a space without a background.
func blankCell(bg color.Color) Cell {
	blank := Cell{Char: ' '}
	if bg == nil {
		return blank
	}
	var st sgrState
	rendered := lipgloss.NewStyle().Background(bg).Render(" ")
	for i := 0; i < len(rendered); {
		unit, _, esc := nextUnit(rendered, i)
		i += len(unit)
		if !esc {
			break // the sequences in front of the space are the color
		}
		st.apply(unit)
	}
	blank.BgColor = st.bg
	return blank
}

// FillBackground pads every line to width cells and gives each cell that has
// no background color the color bg. Lines keep their own colors otherwise.
func FillBackground(lines []string, width int, bg color.Color) []string {
	blank := blankCell(bg)
	out := make([]string, len(lines))
	buf := make([]Cell, max(width, 0))
	for i, line := range lines {
		parseCellsInto(buf, line, blank)
		out[i] = renderCells(buf)
	}
	return out
}

// renderCells converts cells back to a string with ANSI codes.
func renderCells(cells []Cell) string {
	var b strings.Builder
	b.Grow(len(cells) * 2)
	writeCells(&b, cells)
	return b.String()
}

// writeCells writes a row of cells as text, with a style sequence only where
// the style changes between neighbouring cells, a hyperlink sequence only
// where the link changes, and a reset at the end.
func writeCells(b *strings.Builder, cells []Cell) {
	var last Cell // the styling written so far

	for i := range cells {
		cell := &cells[i]
		if cell.Char == 0 {
			continue // the second cell of a wide character
		}

		if cell.Link != last.Link {
			writeLink(b, cell.Link)
			last.Link = cell.Link
		}
		if cell.Attrs != last.Attrs || cell.Underline != last.Underline || cell.FgColor != last.FgColor ||
			cell.BgColor != last.BgColor || cell.UnderlineColor != last.UnderlineColor {
			writeStyle(b, &last, cell)
		}

		if cell.Char >= 0 && cell.Char < utf8.RuneSelf {
			b.WriteByte(byte(cell.Char)) // #nosec G115 -- the bounds check keeps the rune within one byte
		} else {
			b.WriteRune(cell.Char)
		}
		b.WriteString(cell.rest)
	}

	if last.Link != "" {
		writeLink(b, "")
	}
	b.WriteString("\x1b[0m")
}

// writeLink writes the OSC 8 sequence that opens link, or closes the open
// hyperlink when link is "".
func writeLink(b *strings.Builder, link string) {
	b.WriteString(hyperlinkPrefix)
	if link == "" {
		b.WriteByte(';')
	} else {
		b.WriteString(link)
	}
	b.WriteByte('\a')
}

// attrCodes pairs each attribute with the SGR parameter that turns it on.
var attrCodes = [...]struct {
	attr Attr
	code string
}{
	{AttrBold, "1"}, {AttrFaint, "2"}, {AttrItalic, "3"}, {AttrBlink, "5"},
	{AttrReverse, "7"}, {AttrConceal, "8"}, {AttrStrikethrough, "9"},
}

// underlineCodes holds the SGR parameter of each underline shape.
var underlineCodes = [...]string{
	UnderlineSingle: "4", UnderlineDouble: "4:2", UnderlineCurly: "4:3",
	UnderlineDotted: "4:4", UnderlineDashed: "4:5",
}

// writeStyle writes the SGR sequence that takes the styling from last to
// cell, and records the new styling in last. Styling that turns off needs a
// reset first, after which everything the cell has is written again.
func writeStyle(b *strings.Builder, last, cell *Cell) {
	if last.Attrs&^cell.Attrs != 0 ||
		(last.Underline != UnderlineNone && cell.Underline == UnderlineNone) ||
		(last.FgColor != "" && cell.FgColor == "") ||
		(last.BgColor != "" && cell.BgColor == "") ||
		(last.UnderlineColor != "" && cell.UnderlineColor == "") {
		b.WriteString("\x1b[0m")
		*last = Cell{Link: last.Link}
	}

	open := false
	if on := cell.Attrs &^ last.Attrs; on != 0 {
		for _, a := range attrCodes {
			if on&a.attr != 0 {
				writeParam(b, &open, a.code)
			}
		}
	}
	if cell.Underline != last.Underline && int(cell.Underline) < len(underlineCodes) {
		writeParam(b, &open, underlineCodes[cell.Underline])
	}
	if cell.FgColor != last.FgColor {
		writeParam(b, &open, cell.FgColor)
	}
	if cell.BgColor != last.BgColor {
		writeParam(b, &open, cell.BgColor)
	}
	if cell.UnderlineColor != last.UnderlineColor {
		writeParam(b, &open, cell.UnderlineColor)
	}
	if open {
		b.WriteByte('m')
	}

	last.Attrs, last.Underline = cell.Attrs, cell.Underline
	last.FgColor, last.BgColor, last.UnderlineColor = cell.FgColor, cell.BgColor, cell.UnderlineColor
}

// writeParam writes one SGR parameter, opening the sequence before the first
// and separating the others.
func writeParam(b *strings.Builder, open *bool, param string) {
	if *open {
		b.WriteByte(';')
	} else {
		b.WriteString("\x1b[")
		*open = true
	}
	b.WriteString(param)
}
