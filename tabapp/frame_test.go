package tabapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// fillCell returns the cell a fill style paints: its text and background colors.
func fillCell(fill lipgloss.Style) overlays.Cell {
	return overlays.ParseCells(fill.Render("x"), 1)[0]
}

func TestFrameString(t *testing.T) {
	names := map[Frame]string{
		FrameNone:   "None",
		FrameTheme:  "Theme",
		FrameSingle: "Single",
		FrameDouble: "Double",
		FrameFill:   "Fill",
		Frame(99):   "None",
	}
	for f, want := range names {
		if got := f.String(); got != want {
			t.Errorf("Frame(%d).String() = %q, want %q", int(f), got, want)
		}
	}
}

func TestFrameSettersSyncModel(t *testing.T) {
	app := New()
	if got := app.SetHeaderFrame(FrameDouble).SetFooterFrame(FrameFill); got != app {
		t.Error("the frame setters should return the app for chaining")
	}
	m := app.Model()
	if app.HeaderFrame() != FrameDouble || m.HeaderFrame != FrameDouble {
		t.Errorf("header frame = %v / %v, want Double", app.HeaderFrame(), m.HeaderFrame)
	}
	if app.FooterFrame() != FrameFill || m.FooterFrame != FrameFill {
		t.Errorf("footer frame = %v / %v, want Fill", app.FooterFrame(), m.FooterFrame)
	}

	// Only the three border frames count as a border.
	bordered := map[Frame]bool{FrameNone: false, FrameTheme: true, FrameSingle: true, FrameDouble: true, FrameFill: false}
	for f, want := range bordered {
		app.SetHeaderFrame(f).SetFooterFrame(f)
		if app.ShowHeaderBorder() != want || app.ShowFooterBorder() != want {
			t.Errorf("%v: ShowHeaderBorder=%v ShowFooterBorder=%v, want %v", f, app.ShowHeaderBorder(), app.ShowFooterBorder(), want)
		}
	}

	// The border switches are shorthand for the theme's border or no frame.
	app.SetShowHeaderBorder(true).SetShowFooterBorder(true)
	if m.HeaderFrame != FrameTheme || m.FooterFrame != FrameTheme {
		t.Errorf("frames after the switches on = %v %v, want Theme", m.HeaderFrame, m.FooterFrame)
	}
	app.SetShowHeaderBorder(false).SetShowFooterBorder(false)
	if m.HeaderFrame != FrameNone || m.FooterFrame != FrameNone {
		t.Errorf("frames after the switches off = %v %v, want None", m.HeaderFrame, m.FooterFrame)
	}
}

func TestBandHeightSettersSyncModel(t *testing.T) {
	app := New().SetHeaderHeight(4).SetNavBarHeight(5).SetFooterHeight(6)
	m := app.Model()
	if m.HeaderHeight != 4 || m.NavBarHeight != 5 || m.FooterHeight != 6 {
		t.Errorf("model heights = %d %d %d, want 4 5 6", m.HeaderHeight, m.NavBarHeight, m.FooterHeight)
	}
	if app.HeaderHeight() != 4 || app.NavBarHeight() != 5 || app.FooterHeight() != 6 {
		t.Errorf("getters = %d %d %d, want 4 5 6", app.HeaderHeight(), app.NavBarHeight(), app.FooterHeight())
	}

	app.SetHeaderHeight(-1).SetNavBarHeight(-2).SetFooterHeight(-3)
	if m.HeaderHeight != 0 || m.NavBarHeight != 0 || m.FooterHeight != 0 {
		t.Errorf("negative heights = %d %d %d, want 0 0 0", m.HeaderHeight, m.NavBarHeight, m.FooterHeight)
	}
}

func TestBandMessagesUpdateModel(t *testing.T) {
	app := sizedApp(rootMixApp())
	m := app.Model()

	app.Update(SetHeaderHeightMsg{Height: 3})
	app.Update(SetFooterHeightMsg{Height: 4})
	app.Update(SetHeaderFrameMsg{Frame: FrameSingle})
	app.Update(SetFooterFrameMsg{Frame: FrameFill})
	if m.HeaderHeight != 3 || m.FooterHeight != 4 {
		t.Errorf("heights = %d %d, want 3 4", m.HeaderHeight, m.FooterHeight)
	}
	if m.HeaderFrame != FrameSingle || m.FooterFrame != FrameFill {
		t.Errorf("frames = %v %v, want Single Fill", m.HeaderFrame, m.FooterFrame)
	}

	// A toggle clears any frame, then switches the theme's border on.
	for _, want := range []Frame{FrameNone, FrameTheme, FrameNone} {
		app.Update(ToggleHeaderBorderMsg{})
		app.Update(ToggleFooterBorderMsg{})
		if m.HeaderFrame != want || m.FooterFrame != want {
			t.Errorf("frames after a toggle = %v %v, want %v", m.HeaderFrame, m.FooterFrame, want)
		}
	}
}

// TestHeaderFrames renders the header band in every frame and checks its
// size, its border glyphs and the width the header function is given.
func TestHeaderFrames(t *testing.T) {
	needsNarrowBorders(t)
	cases := []struct {
		frame      Frame
		corner     string // top-left glyph of a border; "" for none
		height     int    // rendered lines for two lines of content
		innerWidth int    // width of the header function's flexbox
	}{
		{FrameNone, "", 2, 80},
		{FrameTheme, styles.Dracula.ContentBorder.TopLeft, 4, 78},
		{FrameSingle, "┌", 4, 78},
		{FrameDouble, "╔", 4, 78},
		{FrameFill, "", 2, 78},
	}
	for _, c := range cases {
		t.Run(c.frame.String(), func(t *testing.T) {
			app := sizedApp(rootMixApp())
			gotWidth := 0
			app.SetHeaderHeight(2).SetHeaderFrame(c.frame).SetHeader(func(fb *flexbox.FlexBox, model interface{}) {
				gotWidth = fb.GetWidth()
				testHeader(fb, model)
			})

			band := ansiSequence.ReplaceAllString(app.renderHeaderSection(), "")
			if lipgloss.Height(band) != c.height || lipgloss.Width(band) != 80 {
				t.Errorf("band is %dx%d, want 80x%d", lipgloss.Width(band), lipgloss.Height(band), c.height)
			}
			if gotWidth != c.innerWidth {
				t.Errorf("header flexbox width = %d, want %d", gotWidth, c.innerWidth)
			}
			if !strings.Contains(band, "HEADER BAND") {
				t.Errorf("header text missing:\n%s", band)
			}
			if c.corner != "" && !strings.HasPrefix(band, c.corner) {
				t.Errorf("band should start with %q:\n%s", c.corner, band)
			}
			if c.corner == "" && strings.ContainsAny(band, "┌╔╭│║") {
				t.Errorf("band should have no border:\n%s", band)
			}
			checkFrame(t, "view", app.Render(), 80, 24)
		})
	}
}

// TestFooterFrames renders the footer in every frame and checks its size,
// its border glyphs and the column its keys start in.
func TestFooterFrames(t *testing.T) {
	needsNarrowBorders(t)
	cases := []struct {
		frame   Frame
		corner  string
		height  int // rendered lines for two lines of content
		textCol int // column of the first key on the last content line
	}{
		{FrameNone, "", 2, 0},
		{FrameTheme, styles.Dracula.FooterBorder.TopLeft, 4, 1},
		{FrameSingle, "┌", 4, 1},
		{FrameDouble, "╔", 4, 1},
		{FrameFill, "", 2, 1},
	}
	for _, c := range cases {
		t.Run(c.frame.String(), func(t *testing.T) {
			app := sizedApp(rootMixApp())
			app.SetFooterHeight(2).SetFooterFrame(c.frame)

			// A framed or filled band is the full width. A bare one is only
			// its text; the canvas fills the rest of the row.
			band := ansiSequence.ReplaceAllString(app.renderFooterSection(), "")
			if w := lipgloss.Width(band); lipgloss.Height(band) != c.height || w > 80 || (c.frame != FrameNone && w != 80) {
				t.Errorf("band is %dx%d, want 80x%d", w, lipgloss.Height(band), c.height)
			}
			if c.corner != "" && !strings.HasPrefix(band, c.corner) {
				t.Errorf("band should start with %q:\n%s", c.corner, band)
			}
			if c.corner == "" && strings.ContainsAny(band, "┌╔╭│║") {
				t.Errorf("band should have no border:\n%s", band)
			}

			lines := strings.Split(band, "\n")
			keys := lines[len(lines)-1]
			if c.frame.bordered() {
				keys = lines[len(lines)-2] // above the bottom border
			}
			before, _, found := strings.Cut(keys, "q Quit")
			if col := lipgloss.Width(before); !found || col != c.textCol {
				t.Errorf("keys start in column %d (found %v), want %d:\n%s", col, found, c.textCol, band)
			}
			checkFrame(t, "view", app.Render(), 80, 24)
		})
	}
}

func TestFooterHeightZeroHidesFooter(t *testing.T) {
	app := sizedApp(rootMixApp())
	app.SetFooterHeight(0).SetFooterFrame(FrameDouble)

	if got := app.renderFooterSection(); got != "" {
		t.Errorf("footer at height 0 = %q, want empty", got)
	}
	lines := checkFrame(t, "no footer", app.Render(), 80, 24)
	if plain := ansiSequence.ReplaceAllString(strings.Join(lines, "\n"), ""); strings.Contains(plain, "Quit") {
		t.Errorf("the footer keys should be gone:\n%s", plain)
	}
}

// TestFilledHeaderTakesFillColors checks the bar behind a header function
// that sets no colors of its own: every cell has the fill background.
func TestFilledHeaderTakesFillColors(t *testing.T) {
	useTrueColor(t)
	app := sizedApp(rootMixApp())
	app.SetHeader(testHeader).SetHeaderHeight(2).SetHeaderFrame(FrameFill)

	fill := fillCell(app.Model().Theme.HeaderFill)
	rows := strings.Split(app.renderHeaderSection(), "\n")
	if len(rows) != 2 {
		t.Fatalf("filled header is %d lines, want 2", len(rows))
	}
	for r, line := range rows {
		for col, c := range overlays.ParseCells(line, 80) {
			if c.BgColor != fill.BgColor {
				t.Fatalf("row %d col %d: background %q, want the fill %q", r, col, c.BgColor, fill.BgColor)
			}
		}
	}

	// A header function that draws on Theme.HeaderFill gets the reversed text.
	app.SetHeader(func(fb *flexbox.FlexBox, model interface{}) {
		fb.UseStyle(model.(*Model).Theme.HeaderFill)
		row := fb.NewRow().StylePassing(true)
		cell := flexbox.NewCell(1, 1)
		cell.SetContent("TITLE")
		row.AddCells(cell)
		fb.AddRows([]*flexbox.Row{row})
	})
	cells := overlays.ParseCells(strings.Split(app.renderHeaderSection(), "\n")[0], 80)
	if cells[0].Char != ' ' || cells[1].Char != 'T' {
		t.Errorf("header text should be inset one column, got %q %q", cells[0].Char, cells[1].Char)
	}
	if cells[1].FgColor != fill.FgColor || cells[1].BgColor != fill.BgColor {
		t.Errorf("title cell = fg %q bg %q, want the fill's %q on %q", cells[1].FgColor, cells[1].BgColor, fill.FgColor, fill.BgColor)
	}

	// The fill follows the theme.
	app.UseTheme(styles.Nord)
	fill = fillCell(styles.Nord.HeaderFill)
	if c := overlays.ParseCells(strings.Split(app.renderHeaderSection(), "\n")[0], 80)[1]; c.BgColor != fill.BgColor {
		t.Errorf("after a theme change the fill is %q, want %q", c.BgColor, fill.BgColor)
	}
}

// TestFilledFooterReversesKeys checks a filled footer: fill behind every
// cell, keys and labels in the fill's text color, keys bold, and the glowing
// key as a chip of the content background.
func TestFilledFooterReversesKeys(t *testing.T) {
	useTrueColor(t)
	app := New()
	app.MustAddPage(page.New("Only").BindKey("r", "Refresh", NoOpMsg{}))
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 12})
	app.SetFooterHeight(1).SetFooterFrame(FrameFill)

	theme := app.Model().Theme
	fill := fillCell(theme.FooterFill)
	footerCells := func() []overlays.Cell {
		lines := strings.Split(app.Render(), "\n")
		return overlays.ParseCells(lines[len(lines)-1], 60)
	}

	cells := footerCells()
	var text strings.Builder
	for col, c := range cells {
		if c.BgColor != fill.BgColor {
			t.Fatalf("col %d: background %q, want the fill %q", col, c.BgColor, fill.BgColor)
		}
		if c.Char != ' ' && c.FgColor != fill.FgColor {
			t.Fatalf("col %d (%q): foreground %q, want the fill's %q", col, c.Char, c.FgColor, fill.FgColor)
		}
		text.WriteRune(c.Char)
	}
	plain := text.String()
	keyCol := strings.Index(plain, "q Quit")
	if keyCol != 1 {
		t.Fatalf("footer = %q, want the keys inset one column", plain)
	}
	if cells[keyCol].Attrs&overlays.AttrBold == 0 || (cells[keyCol+2].Attrs&overlays.AttrBold != 0) {
		t.Errorf("the key should be bold and its label not: key bold=%v label bold=%v", (cells[keyCol].Attrs&overlays.AttrBold != 0), (cells[keyCol+2].Attrs&overlays.AttrBold != 0))
	}

	// Pressing r makes its key glow: glow color on the content background.
	runeKey(app, 'r')
	cells = footerCells()
	before, _, _ := strings.Cut(plain, "r Refresh")
	glowCol := lipgloss.Width(before)
	glow := overlays.ParseCells(lipgloss.NewStyle().Foreground(app.KeyGlowColor()).Background(theme.Palette.Bg0).Render("r"), 1)[0]
	if got := cells[glowCol]; got.FgColor != glow.FgColor || got.BgColor != glow.BgColor || got.Attrs&overlays.AttrBold == 0 {
		t.Errorf("glowing key = fg %q bg %q bold %v, want fg %q bg %q bold", got.FgColor, got.BgColor, (got.Attrs&overlays.AttrBold != 0), glow.FgColor, glow.BgColor)
	}
	if got := cells[glowCol+2]; got.BgColor != fill.BgColor {
		t.Errorf("the label beside a glowing key keeps the fill, got bg %q", got.BgColor)
	}
}

// TestFilledFooterColorsPlainText covers the footer text the footer style
// renders, the default line of a page without labelled keys: the fill sets
// its colors and the footer style keeps its other attributes.
func TestFilledFooterColorsPlainText(t *testing.T) {
	useTrueColor(t)
	app := New()
	app.UnbindKey("q")
	app.MustAddPage(page.New("Bare"))
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 12})
	app.SetFooterStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")).Bold(true)).
		SetFooterHeight(1).
		SetFooterFrame(FrameFill)

	fill := fillCell(app.Model().Theme.FooterFill)
	cells := overlays.ParseCells(app.renderFooterSection(), 60)
	if cells[1].Char != 'T' { // "Tab: Navigate • Ctrl+Q: Quit"
		t.Fatalf("footer starts with %q, want the default text inset one column", cells[1].Char)
	}
	if cells[1].FgColor != fill.FgColor || cells[1].BgColor != fill.BgColor {
		t.Errorf("text = fg %q bg %q, want the fill's %q on %q", cells[1].FgColor, cells[1].BgColor, fill.FgColor, fill.BgColor)
	}
	if cells[1].Attrs&overlays.AttrBold == 0 {
		t.Error("the footer style's bold should survive the fill")
	}
}

// TestBandsAtHeightZeroKeepTheFrame renders with every band at height 0:
// the header and footer drop out and the frame still fills the screen.
func TestBandsAtHeightZeroKeepTheFrame(t *testing.T) {
	app := rootMixApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 30})
	app.SetHeader(testHeader).SetNavBarAutoHide(false).
		SetHeaderHeight(0).SetNavBarHeight(0).SetFooterHeight(0)

	lines := checkFrame(t, "zero heights", app.Render(), 100, 30)
	plain := ansiSequence.ReplaceAllString(strings.Join(lines, "\n"), "")
	if strings.Contains(plain, "HEADER BAND") || strings.Contains(plain, "Quit") {
		t.Errorf("header and footer should be hidden:\n%s", plain)
	}
}

// TestFramesRenderBeforeTheTerminalSizeArrives calls View with no size set,
// as Bubble Tea may before its first WindowSizeMsg, and at a width too
// narrow for a frame. Every frame must render without panicking.
func TestFramesRenderBeforeTheTerminalSizeArrives(t *testing.T) {
	for _, size := range []tea.WindowSizeMsg{{}, {Width: 1, Height: 1}, {Width: 3, Height: 4}} {
		for _, frame := range []Frame{FrameNone, FrameTheme, FrameSingle, FrameDouble, FrameFill} {
			app := rootMixApp()
			app.SetHeader(testHeader).SetHeaderFrame(frame).SetFooterFrame(frame)
			if size.Width > 0 {
				app.Update(size)
			}
			if got := len(strings.Split(app.Render(), "\n")); size.Height > 0 && got != size.Height {
				t.Errorf("%v at %dx%d: %d lines, want %d", frame, size.Width, size.Height, got, size.Height)
			}
		}
	}
}

// TestBandOf covers the band helper: it clips from the far end, pads the
// near end with empty lines, and always returns exactly the height.
func TestBandOf(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		height   int
		atBottom bool
		want     string
	}{
		{"top, exact", []string{"a", "b"}, 2, false, "a\nb"},
		{"top, short is padded below", []string{"a"}, 3, false, "a\n\n"},
		{"top, long keeps the first lines", []string{"a", "b", "c"}, 2, false, "a\nb"},
		{"bottom, short is padded above", []string{"a"}, 3, true, "\n\na"},
		{"bottom, long keeps the last lines", []string{"a", "b", "c"}, 2, true, "b\nc"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bandOf(tt.lines, tt.height, tt.atBottom)
			if got != tt.want {
				t.Errorf("bandOf = %q, want %q", got, tt.want)
			}
			if blockHeight(got) != tt.height {
				t.Errorf("band is %d lines, want %d", blockHeight(got), tt.height)
			}
		})
	}
}
