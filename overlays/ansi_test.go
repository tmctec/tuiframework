package overlays

import (
	"testing"

	"charm.land/lipgloss/v2"
)

// cellsEqual compares two cell slices and reports the first difference.
func cellsEqual(t *testing.T, got, want []Cell) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("got %d cells, want %d", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("cell %d = %+v, want %+v", i, got[i], want[i])
		}
	}
}

// cell builds a plain cell for a character.
func cell(r rune) Cell { return Cell{Char: r} }

func TestParseCellsPlainText(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		width int
		want  []Cell
	}{
		{"pads with spaces", "ab", 4, []Cell{cell('a'), cell('b'), cell(' '), cell(' ')}},
		{"truncates at width", "abcdef", 3, []Cell{cell('a'), cell('b'), cell('c')}},
		{"empty input", "", 2, []Cell{cell(' '), cell(' ')}},
		{"zero width", "abc", 0, []Cell{}},
		{"exact fit", "abc", 3, []Cell{cell('a'), cell('b'), cell('c')}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cellsEqual(t, ParseCells(tt.in, tt.width), tt.want)
		})
	}
}

func TestParseCellsColors(t *testing.T) {
	red := "38;2;255;0;0"
	bg := "48;2;40;42;54"
	tests := []struct {
		name  string
		in    string
		width int
		want  []Cell
	}{
		{
			"24-bit foreground",
			"\x1b[" + red + "mab\x1b[0m",
			3,
			[]Cell{{Char: 'a', FgColor: red}, {Char: 'b', FgColor: red}, cell(' ')},
		},
		{
			"24-bit background",
			"\x1b[" + bg + "mx\x1b[0my",
			2,
			[]Cell{{Char: 'x', BgColor: bg}, cell('y')},
		},
		{
			"256-color foreground and background",
			"\x1b[38;5;196;48;5;21mx",
			1,
			[]Cell{{Char: 'x', FgColor: "38;5;196", BgColor: "48;5;21"}},
		},
		{
			"combined bold foreground background",
			"\x1b[1;" + red + ";" + bg + "mx\x1b[0m",
			1,
			[]Cell{{Char: 'x', Attrs: AttrBold, FgColor: red, BgColor: bg}},
		},
		{
			"color persists across characters until reset",
			"\x1b[" + bg + "mab\x1b[0mc",
			3,
			[]Cell{
				{Char: 'a', BgColor: bg},
				{Char: 'b', BgColor: bg},
				cell('c'),
			},
		},
		{
			"background colors with zero components",
			"\x1b[48;2;0;0;0mx",
			1,
			[]Cell{{Char: 'x', BgColor: "48;2;0;0;0"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cellsEqual(t, ParseCells(tt.in, tt.width), tt.want)
		})
	}
}

func TestParseCellsBoldAndResets(t *testing.T) {
	fg := "38;2;1;2;3"
	bg := "48;2;4;5;6"
	tests := []struct {
		name  string
		in    string
		width int
		want  []Cell
	}{
		{
			"bold on then off",
			"\x1b[1ma\x1b[22mb",
			2,
			[]Cell{{Char: 'a', Attrs: AttrBold}, cell('b')},
		},
		{
			"short reset clears everything",
			"\x1b[1;" + fg + ";" + bg + "ma\x1b[mb",
			2,
			[]Cell{{Char: 'a', Attrs: AttrBold, FgColor: fg, BgColor: bg}, cell('b')},
		},
		{
			"full reset clears everything",
			"\x1b[1;" + fg + ";" + bg + "ma\x1b[0mb",
			2,
			[]Cell{{Char: 'a', Attrs: AttrBold, FgColor: fg, BgColor: bg}, cell('b')},
		},
		{
			"foreground reset keeps background",
			"\x1b[" + fg + ";" + bg + "ma\x1b[39mb",
			2,
			[]Cell{
				{Char: 'a', FgColor: fg, BgColor: bg},
				{Char: 'b', BgColor: bg},
			},
		},
		{
			"background reset keeps foreground",
			"\x1b[" + fg + ";" + bg + "ma\x1b[49mb",
			2,
			[]Cell{
				{Char: 'a', FgColor: fg, BgColor: bg},
				{Char: 'b', FgColor: fg},
			},
		},
		{
			"reset prefix clears bold and background before new foreground",
			"\x1b[1;" + bg + "ma\x1b[0;" + fg + "mb",
			2,
			[]Cell{
				{Char: 'a', Attrs: AttrBold, BgColor: bg},
				{Char: 'b', FgColor: fg},
			},
		},
		{
			"unrelated sequence keeps state",
			"\x1b[1ma\x1b[10m\x1b[2Kb",
			2,
			[]Cell{{Char: 'a', Attrs: AttrBold}, {Char: 'b', Attrs: AttrBold}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cellsEqual(t, ParseCells(tt.in, tt.width), tt.want)
		})
	}
}

func TestParseCellsMalformedColorSequence(t *testing.T) {
	// A sequence ended by a letter other than m is not a style, and a 38 or
	// 48 without its color components sets no color; the character is still
	// placed, unstyled.
	for _, in := range []string{
		"\x1b[48;2;1;2x\x1b[38;2;3;4xa", // cut short by a non-SGR terminator
		"\x1b[48;2;1;2m\x1b[38;9;3ma",   // too few components, unknown color kind
	} {
		cells := ParseCells(in, 1)
		if cells[0].Char != 'a' {
			t.Fatalf("%q: cell 0 = %+v, want the character placed", in, cells[0])
		}
		if cells[0] != (Cell{Char: 'a'}) {
			t.Errorf("%q: cell 0 = %+v, want an unstyled cell", in, cells[0])
		}
	}
}

func TestParseCellsResetCodesInsideColorsAreColors(t *testing.T) {
	// 39, 49 and 22 inside a color's components are components, and a 0
	// before 1 in one sequence resets and then sets bold.
	tests := []struct {
		name string
		in   string
		want Cell
	}{
		{"49 in a foreground keeps the background", "\x1b[48;2;1;2;3m\x1b[38;2;49;0;0ma", Cell{Char: 'a', BgColor: "48;2;1;2;3", FgColor: "38;2;49;0;0"}},
		{"256-color 49 keeps the background", "\x1b[48;2;1;2;3m\x1b[38;5;49ma", Cell{Char: 'a', BgColor: "48;2;1;2;3", FgColor: "38;5;49"}},
		{"22 in a color keeps bold", "\x1b[1m\x1b[38;2;22;0;0ma", Cell{Char: 'a', FgColor: "38;2;22;0;0", Attrs: AttrBold}},
		{"39 in a background keeps the foreground", "\x1b[38;2;1;2;3m\x1b[48;2;39;0;0ma", Cell{Char: 'a', BgColor: "48;2;39;0;0", FgColor: "38;2;1;2;3"}},
		{"reset then bold in one sequence", "\x1b[31m\x1b[0;1ma", Cell{Char: 'a', Attrs: AttrBold}},
		{"basic colors are kept as their code", "\x1b[31;44ma", Cell{Char: 'a', BgColor: "44", FgColor: "31"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseCells(tt.in, 1)[0]; got != tt.want {
				t.Errorf("cell = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestParseCellsWideCharacters(t *testing.T) {
	fg := "38;2;1;2;3"
	tests := []struct {
		name  string
		in    string
		width int
		want  []Cell
	}{
		{
			"CJK takes two cells with a continuation marker",
			"日本",
			4,
			[]Cell{cell('日'), {Char: 0}, cell('本'), {Char: 0}},
		},
		{
			"emoji continuation inherits the style",
			"\x1b[" + fg + "m🎉\x1b[0m",
			2,
			[]Cell{{Char: '🎉', FgColor: fg}, {Char: 0, FgColor: fg}},
		},
		{
			"wide character that does not fit is dropped",
			"a日",
			2,
			[]Cell{cell('a'), cell(' ')},
		},
		{
			"wide then narrow",
			"日a",
			3,
			[]Cell{cell('日'), {Char: 0}, cell('a')},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cellsEqual(t, ParseCells(tt.in, tt.width), tt.want)
		})
	}
}

func TestParseCellsRoundTripsLipglossOutput(t *testing.T) {
	styled := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#ff0000")).
		Background(lipgloss.Color("#282A36")).
		Render("ok")

	cells := ParseCells(styled, 3)
	want := Cell{Char: 'o', Attrs: AttrBold, FgColor: "38;2;255;0;0", BgColor: "48;2;40;42;54"}
	if cells[0] != want {
		t.Errorf("cell 0 = %+v, want %+v", cells[0], want)
	}
	if cells[1].Char != 'k' || cells[1] == cell('k') {
		t.Errorf("cell 1 = %+v, want a styled 'k'", cells[1])
	}
	if cells[2] != cell(' ') {
		t.Errorf("cell 2 = %+v, want an unstyled padding space", cells[2])
	}
}

func TestFillBackgroundPadsAndFills(t *testing.T) {
	bg := "48;2;40;42;54"

	out := FillBackground([]string{"ab", ""}, 4, lipgloss.Color("#282A36"))
	if len(out) != 2 {
		t.Fatalf("got %d lines, want 2", len(out))
	}
	if want := "\x1b[" + bg + "mab  \x1b[0m"; out[0] != want {
		t.Errorf("line 0 = %q, want %q", out[0], want)
	}
	for i, line := range out {
		if w := lipgloss.Width(line); w != 4 {
			t.Errorf("line %d width = %d, want 4", i, w)
		}
		for c, cl := range ParseCells(line, 4) {
			if !cl.HasBackground() || cl.BgColor != bg {
				t.Errorf("line %d cell %d = %+v, want background %q", i, c, cl, bg)
			}
		}
	}
}

func TestFillBackgroundKeepsExistingColors(t *testing.T) {
	own := "48;2;1;2;3"
	fg := "38;2;255;0;0"

	out := FillBackground([]string{"\x1b[" + own + "mX\x1b[0my", "\x1b[" + fg + "mz\x1b[0m"}, 2, lipgloss.Color("#282A36"))

	first := ParseCells(out[0], 2)
	if first[0].BgColor != own {
		t.Errorf("cell with its own background = %q, want %q kept", first[0].BgColor, own)
	}
	if first[1].BgColor != "48;2;40;42;54" {
		t.Errorf("cell without background = %q, want the fill color", first[1].BgColor)
	}

	second := ParseCells(out[1], 2)
	if second[0].FgColor != fg || second[0].BgColor != "48;2;40;42;54" {
		t.Errorf("styled cell = %+v, want foreground kept and background filled", second[0])
	}
}

func TestRenderCells(t *testing.T) {
	red := "38;2;255;0;0"
	blue := "38;2;0;0;255"
	bg := "48;2;1;2;3"
	tests := []struct {
		name  string
		cells []Cell
		want  string
	}{
		{"empty ends with a reset", []Cell{}, "\x1b[0m"},
		{"plain text", []Cell{cell('a'), cell('b')}, "ab\x1b[0m"},
		{
			"unchanged style is reused",
			[]Cell{{Char: 'a', FgColor: red}, {Char: 'b', FgColor: red}},
			"\x1b[" + red + "mab\x1b[0m",
		},
		{
			"changed color needs no reset",
			[]Cell{{Char: 'a', FgColor: red}, {Char: 'b', FgColor: blue}},
			"\x1b[" + red + "ma\x1b[" + blue + "mb\x1b[0m",
		},
		{
			"bold turning off emits a reset",
			[]Cell{{Char: 'a', Attrs: AttrBold}, cell('b')},
			"\x1b[1ma\x1b[0mb\x1b[0m",
		},
		{
			"foreground turning off emits a reset",
			[]Cell{{Char: 'a', FgColor: red}, cell('b')},
			"\x1b[" + red + "ma\x1b[0mb\x1b[0m",
		},
		{
			"background turning off emits a reset",
			[]Cell{{Char: 'a', BgColor: bg}, cell('b')},
			"\x1b[" + bg + "ma\x1b[0mb\x1b[0m",
		},
		{
			"reset then new style on the same cell",
			[]Cell{{Char: 'a', Attrs: AttrBold, FgColor: red}, {Char: 'b', FgColor: blue}},
			"\x1b[1;" + red + "ma\x1b[0m\x1b[" + blue + "mb\x1b[0m",
		},
		{
			"bold foreground and background in one sequence",
			[]Cell{{Char: 'a', Attrs: AttrBold, FgColor: red, BgColor: bg}},
			"\x1b[1;" + red + ";" + bg + "ma\x1b[0m",
		},
		{
			"continuation markers are skipped",
			[]Cell{cell('日'), {Char: 0}, cell('本'), {Char: 0}},
			"日本\x1b[0m",
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

func TestRenderCellsRoundTrip(t *testing.T) {
	in := "\x1b[1;38;2;255;0;0ma\x1b[0m\x1b[48;2;1;2;3m日\x1b[0mb"
	cells := ParseCells(in, 5)
	again := ParseCells(renderCells(cells), 5)
	cellsEqual(t, again, cells)
}
