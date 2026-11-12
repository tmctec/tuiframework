package overlays

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

// TestParseCellsMeasuresLikeLipgloss checks the cell parser against lipgloss
// on text where a rune-by-rune measure goes wrong: combining marks, emoji
// built from several runes, and characters whose width depends on locale.
func TestParseCellsMeasuresLikeLipgloss(t *testing.T) {
	texts := []string{
		"plain ascii",
		"box ─│╭╮╰╯═║ lines",
		"arrows ▸▾◂▴ • … ›",
		"cafe\u0301 decomposed",         // e + combining acute
		"Thai \u0e17\u0e35\u0e48 marks", // consonant, vowel sign, tone mark
		"flag \U0001F1EF\U0001F1F5 pair",
		"warning \u26a0\ufe0f sign",                             // text symbol + emoji variation selector
		"family \U0001F468\u200d\U0001F469\u200d\U0001F467 zwj", // joined emoji
		"keycap 1\ufe0f\u20e3 one",
		"日本語 wide",
		"tone \U0001F44D\U0001F3FD modifier",
		"हिन्दी क्षि Devanagari", // conjuncts and vowel signs
	}
	for _, text := range texts {
		want := lipgloss.Width(text)
		if got := lineWidth(text); got != want {
			t.Errorf("lineWidth(%q) = %d, want lipgloss's %d", text, got, want)
		}

		// Parsed at its own width and written back, the text is unchanged.
		cells := ParseCells(text, want)
		if got := strings.TrimSuffix(renderCells(cells), "\x1b[0m"); got != text {
			t.Errorf("round trip of %q = %q", text, got)
		}
		// One cell wider, exactly one padding space follows.
		wider := ParseCells(text, want+1)
		if last := wider[want]; last.Char != ' ' {
			t.Errorf("%q: cell %d = %q, want the padding space", text, want, last.Char)
		}
	}
}

func TestParseCellsKeepsMarksWithTheirCharacter(t *testing.T) {
	cells := ParseCells("e\u0301x", 3)
	if cells[0].Char != 'e' || cells[0].rest != "\u0301" {
		t.Errorf("cell 0 = %q + %q, want e with its combining acute", cells[0].Char, cells[0].rest)
	}
	if cells[1].Char != 'x' || cells[2].Char != ' ' {
		t.Errorf("cells 1 and 2 = %q %q, want x and a space", cells[1].Char, cells[2].Char)
	}

	// A mark after a style change still joins the character before it.
	cells = ParseCells("\x1b[31me\x1b[0m\u0301", 1)
	if cells[0].rest != "\u0301" || cells[0].FgColor != "31" {
		t.Errorf("cell = %+v, want the mark joined and the color kept", cells[0])
	}

	// A mark with nothing before it, and control characters, are dropped.
	cells = ParseCells("\u0301a\tb\x07", 3)
	if cells[0].Char != 'a' || cells[0].rest != "" || cells[1].Char != 'b' || cells[2].Char != ' ' {
		t.Errorf("cells = %q %q %q, want a, b and a space", cells[0].Char, cells[1].Char, cells[2].Char)
	}
}

// TestParseCellsSkipsOtherEscapeSequences covers the sequences that are not
// SGR: they take no cells and leave the styling alone.
func TestParseCellsSkipsOtherEscapeSequences(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"cursor movement", "a\x1b[2Kb\x1b[10;5Hc", "abc"},
		{"hyperlink closed by BEL", "\x1b]8;;https://example.com\aLink\x1b]8;;\a!", "Link!"},
		{"hyperlink closed by ST", "\x1b]8;;https://example.com\x1b\\Link\x1b]8;;\x1b\\!", "Link!"},
		{"two-byte escape", "a\x1b7b", "ab"},
		{"escape at the end of the line", "ab\x1b", "ab"},
		{"unterminated CSI", "ab\x1b[31", "ab"},
		{"unterminated OSC", "ab\x1b]0;title", "ab"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lineWidth(tt.in); got != len(tt.want) {
				t.Errorf("lineWidth = %d, want %d", got, len(tt.want))
			}
			cells := ParseCells(tt.in, len(tt.want))
			var b strings.Builder
			for _, c := range cells {
				b.WriteRune(c.Char)
				if c.FgColor != "" || c.HasBackground() || c.Attrs&AttrBold != 0 {
					t.Errorf("cell %q is styled: %+v", c.Char, c)
				}
			}
			if b.String() != tt.want {
				t.Errorf("cells = %q, want %q", b.String(), tt.want)
			}
		})
	}
}

// TestBoxDrawingMeasuresLikeLipgloss guards the East Asian ambiguous
// characters: box-drawing lines and the navigation arrows take the width
// Lip Gloss lays them out with. That is one cell each whatever the locale
// says, and two only when RUNEWIDTH_EASTASIAN asks for it.
func TestBoxDrawingMeasuresLikeLipgloss(t *testing.T) {
	const symbols = "─│┌┐└┘╭╮╰╯═║╔╗╚╝▸▾◂▴•…›"
	n := len([]rune(symbols))
	want := lipgloss.Width(symbols)
	if !eastAsianWidths() && want != n {
		t.Fatalf("lipgloss measures %d cells, want one per symbol (%d)", want, n)
	}
	if got := lineWidth(symbols); got != want {
		t.Errorf("lineWidth = %d, want lipgloss's %d", got, want)
	}

	shown := 0
	for i, c := range ParseCells(symbols, want) {
		switch c.Char {
		case ' ':
			t.Errorf("cell %d is padding, want the symbols to fill the line", i)
		case 0: // the second cell of a wide symbol
		default:
			shown++
		}
	}
	if shown != n {
		t.Errorf("%d symbols parsed, want %d", shown, n)
	}
}
