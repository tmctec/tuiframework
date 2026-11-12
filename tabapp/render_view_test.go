package tabapp

import (
	"fmt"
	"image/color"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
)

// sizedApp returns app after a 80x24 window size message.
func sizedApp(app *TabApp) *TabApp {
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return app
}

// hideNavBar arms the autohide timer and fires it, so the bar is hidden.
func hideNavBar(t *testing.T, app *TabApp) {
	t.Helper()
	if cmd := app.Init(); cmd == nil {
		t.Fatal("Init should arm the autohide timer while autohide is on")
	}
	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID})
	if !app.navBarHidden {
		t.Fatal("the autohide timeout should hide the nav bar")
	}
}

// checkFrame reports whether view is exactly w cells by h lines.
func checkFrame(t *testing.T, label, view string, w, h int) []string {
	t.Helper()
	lines := strings.Split(view, "\n")
	if len(lines) != h {
		t.Errorf("%s: %d lines, want %d", label, len(lines), h)
	}
	for i, l := range lines {
		if lw := lipgloss.Width(l); lw != w {
			t.Errorf("%s line %d: width %d, want %d", label, i, lw, w)
			break
		}
	}
	return lines
}

// sizeRecorder is a page content function that records the flexbox size it
// receives on every render.
func sizeRecorder(w, h *int) page.ContentFunc {
	return func(fb *flexbox.FlexBox, _ interface{}, _, _ int) {
		*w, *h = fb.GetWidth(), fb.GetHeight()
	}
}

func TestViewNoContentMessage(t *testing.T) {
	app := sizedApp(New())
	got := app.Render()
	if !strings.Contains(got, "No content") || strings.Contains(got, "\n") {
		t.Errorf("View() = %q, want the single-line no-content message", got)
	}
}

// TestViewFrameWithNavBarHidden hides the bar through the autohide timeout in
// every position and mode: the frame keeps the terminal size, the bar renders
// nothing, and the content takes the bar's room.
func TestViewFrameWithNavBarHidden(t *testing.T) {
	const w, h = 80, 24
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		for _, mode := range []NavBarMode{NavBarCompact, NavBarExpanded, NavBarPath, NavBarColumns} {
			label := fmt.Sprintf("pos=%v mode=%v", pos, mode)
			var cw, ch int
			app := rootMixApp()
			app.Pages()[0].SetContent(sizeRecorder(&cw, &ch))
			app.SetNavBarPosition(pos).SetNavBarMode(mode).SetNavBarAutoHide(true)
			sizedApp(app)

			app.Render()
			shownW, shownH := cw, ch
			hideNavBar(t, app)

			lines := checkFrame(t, label, app.Render(), w, h)
			if pos == NavBarTop {
				if nav := app.renderNavSection(); nav != "" {
					t.Errorf("%s: hidden nav section = %q, want empty", label, nav)
				}
				if ch <= shownH {
					t.Errorf("%s: content height %d with the bar hidden, want more than %d", label, ch, shownH)
				}
			} else {
				if side := app.renderSideNavBar(h); side != "" {
					t.Errorf("%s: hidden side bar = %q, want empty", label, side)
				}
				if cw <= shownW {
					t.Errorf("%s: content width %d with the bar hidden, want more than %d", label, cw, shownW)
				}
			}
			plain := ansiSequence.ReplaceAllString(strings.Join(lines, "\n"), "")
			for _, tab := range []string{"▾ A", "▸ A", "a1", "nav mode"} {
				if strings.Contains(plain, tab) {
					t.Errorf("%s: %q visible while the bar is hidden", label, tab)
				}
			}
		}
	}
}

func TestViewHeaderBand(t *testing.T) {
	needsNarrowBorders(t)
	app := sizedApp(rootMixApp())

	if got := app.renderHeaderSection(); got != "" {
		t.Errorf("no header function: %q, want empty", got)
	}
	app.SetHeader(testHeader).SetHeaderHeight(0)
	if got := app.renderHeaderSection(); got != "" {
		t.Errorf("header height 0: %q, want empty", got)
	}

	app.SetHeaderHeight(2)
	plainHeader := ansiSequence.ReplaceAllString(app.renderHeaderSection(), "")
	if lipgloss.Height(plainHeader) != 2 || lipgloss.Width(plainHeader) != 80 {
		t.Errorf("header is %dx%d, want 80x2", lipgloss.Width(plainHeader), lipgloss.Height(plainHeader))
	}
	if !strings.Contains(plainHeader, "HEADER BAND") {
		t.Errorf("header text missing:\n%s", plainHeader)
	}

	app.SetShowHeaderBorder(true)
	bordered := ansiSequence.ReplaceAllString(app.renderHeaderSection(), "")
	if lipgloss.Height(bordered) != 4 || lipgloss.Width(bordered) != 80 {
		t.Errorf("bordered header is %dx%d, want 80x4", lipgloss.Width(bordered), lipgloss.Height(bordered))
	}
	if !strings.HasPrefix(bordered, app.contentBorderStyle.TopLeft) || !strings.Contains(bordered, "HEADER BAND") {
		t.Errorf("bordered header should start with the border and keep the text:\n%s", bordered)
	}

	// A header function that renders more lines than the band is clipped.
	app.SetShowHeaderBorder(false).SetHeaderHeight(1).SetHeader(func(fb *flexbox.FlexBox, _ interface{}) {
		fb.SetHeight(3)
		for _, text := range []string{"ROW ONE", "ROW TWO", "ROW THREE"} {
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetContent(text)
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		}
	})
	clipped := ansiSequence.ReplaceAllString(app.renderHeaderSection(), "")
	if lipgloss.Height(clipped) != 1 {
		t.Errorf("header with three rows is %d lines, want 1:\n%s", lipgloss.Height(clipped), clipped)
	}
	if strings.Contains(clipped, "ROW THREE") {
		t.Errorf("clipped header should drop the rows past its height:\n%s", clipped)
	}
}

func TestViewContentSection(t *testing.T) {
	needsNarrowBorders(t)
	app := New()
	app.MustAddPage(page.New("Body").SetContent(func(fb *flexbox.FlexBox, _ interface{}, _, _ int) {
		row := fb.NewRow()
		cell := flexbox.NewCell(1, 1)
		cell.SetContent("PAGE BODY")
		row.AddCells(cell)
		fb.AddRows([]*flexbox.Row{row})
	}))
	sizedApp(app)

	plain := ansiSequence.ReplaceAllString(app.renderContentSection(80, 18), "")
	if lipgloss.Width(plain) != 80 || lipgloss.Height(plain) != 18 {
		t.Errorf("content is %dx%d, want 80x18", lipgloss.Width(plain), lipgloss.Height(plain))
	}
	if !strings.Contains(plain, "PAGE BODY") {
		t.Errorf("page body missing:\n%s", plain)
	}
	if strings.Contains(plain, app.contentBorderStyle.TopLeft) {
		t.Error("content without a border should have no border rune")
	}

	app.SetShowContentBorder(true)
	bordered := ansiSequence.ReplaceAllString(app.renderContentSection(80, 18), "")
	if lipgloss.Width(bordered) != 80 || lipgloss.Height(bordered) != 18 {
		t.Errorf("bordered content is %dx%d, want 80x18", lipgloss.Width(bordered), lipgloss.Height(bordered))
	}
	if !strings.HasPrefix(bordered, app.contentBorderStyle.TopLeft) || !strings.Contains(bordered, "PAGE BODY") {
		t.Errorf("bordered content should start with the border and keep the body:\n%s", bordered)
	}
}

// TestViewEmptySectionShowsNoContent renders a section without pages: the
// content area says so and the frame still fills the screen.
func TestViewEmptySectionShowsNoContent(t *testing.T) {
	app := New()
	app.MustAddSection(NewSection("Empty"))
	app.SetNavBarAutoHide(false)
	sizedApp(app)

	if got := app.renderContentSection(80, 18); got != "No content" {
		t.Errorf("renderContentSection = %q, want \"No content\"", got)
	}
	lines := checkFrame(t, "empty section", app.Render(), 80, 24)
	if !strings.Contains(ansiSequence.ReplaceAllString(strings.Join(lines, "\n"), ""), "No content") {
		t.Error("the frame should carry the no-content text")
	}
}

func TestViewFooterClipsAndSitsAtBottom(t *testing.T) {
	needsNarrowBorders(t)
	app := New()
	app.MustAddPage(page.New("Only").SetFooter(map[string]string{
		"a": strings.Repeat("alpha ", 5),
		"z": strings.Repeat("zulu ", 5),
	}))
	app.Update(tea.WindowSizeMsg{Width: 40, Height: 20})

	app.SetFooterHeight(1)
	footer := ansiSequence.ReplaceAllString(app.renderFooterSection(), "")
	if lipgloss.Height(footer) != 1 {
		t.Fatalf("footer is %d lines, want 1:\n%s", lipgloss.Height(footer), footer)
	}
	if !strings.Contains(footer, "zulu") || strings.Contains(footer, "alpha") {
		t.Errorf("a clipped footer keeps its last line:\n%s", footer)
	}

	app.SetFooterHeight(4)
	footer = ansiSequence.ReplaceAllString(app.renderFooterSection(), "")
	lines := strings.Split(footer, "\n")
	if len(lines) != 4 {
		t.Fatalf("footer is %d lines, want 4:\n%s", len(lines), footer)
	}
	if strings.TrimSpace(lines[0]) != "" || !strings.Contains(lines[3], "zulu") {
		t.Errorf("footer text should sit on the bottom lines:\n%s", footer)
	}

	app.SetShowFooterBorder(true)
	footer = ansiSequence.ReplaceAllString(app.renderFooterSection(), "")
	if lipgloss.Height(footer) != 6 || lipgloss.Width(footer) != 40 {
		t.Errorf("bordered footer is %dx%d, want 40x6", lipgloss.Width(footer), lipgloss.Height(footer))
	}
}

// TestViewPadsAndColorsEveryCell checks the frame is exactly the terminal
// size and, with a real color profile, that every cell has a background.
func TestViewPadsAndColorsEveryCell(t *testing.T) {
	app := sizedApp(rootMixApp())
	app.SetNavBarAutoHide(false)
	checkFrame(t, "frame", app.Render(), 80, 24)

	useTrueColor(t)
	for row, line := range strings.Split(app.Render(), "\n") {
		for col, cell := range overlays.ParseCells(line, 80) {
			if !cell.HasBackground() {
				t.Fatalf("row %d col %d has no background", row, col)
			}
		}
	}
}

// TestDrawOverlayPlacement draws an overlay over a frame of marker cells and
// checks the frame shows on both sides and on every other row.
func TestDrawOverlayPlacement(t *testing.T) {
	needsNarrowBorders(t)
	app := sizedApp(rootMixApp())
	base := make([]string, 24)
	for i := range base {
		base[i] = strings.Repeat("B", 80)
	}

	for _, pos := range []overlays.Position{overlays.PositionCenter, overlays.PositionTopLeft, overlays.PositionBottomRight} {
		o := overlays.New().SetExitKey("esc").SetWidth(20).SetHeight(5).SetContent("MID").SetPosition(pos)
		fgLines, fgWidth := overlays.GetLines(o.Render(80, 24))
		x, y := overlays.CalculatePosition(pos, fgWidth, len(fgLines), 80, 24)

		app.canvas.Reset(80, 24, nil).Draw(base, 0, 0, 80, 24, nil)
		app.drawOverlay(o)
		lines := strings.Split(app.canvas.String(), "\n")
		if len(lines) != 24 {
			t.Fatalf("%s: %d lines, want 24", pos, len(lines))
		}
		for row, line := range lines {
			cells := []rune(ansiSequence.ReplaceAllString(line, ""))
			if len(cells) != 80 {
				t.Fatalf("%s row %d: %d cells, want 80", pos, row, len(cells))
			}
			inside := row >= y && row < y+len(fgLines)
			for col, r := range cells {
				overlayCol := col >= x && col < x+fgWidth
				if inside && overlayCol {
					continue
				}
				if r != 'B' {
					t.Fatalf("%s row %d col %d = %q, want the base", pos, row, col, r)
				}
			}
			if inside && string(cells[x:x+fgWidth]) != ansiSequence.ReplaceAllString(fgLines[row-y], "") {
				t.Errorf("%s row %d: overlay cells %q", pos, row, string(cells[x:x+fgWidth]))
			}
		}
	}
}

// TestOverlayFillUsesOverlayOrThemeBackground checks the fill under an
// overlay: its own background when it has one, the theme Bg0 otherwise.
func TestOverlayFillUsesOverlayOrThemeBackground(t *testing.T) {
	bgSeq := func(c color.Color) string {
		return overlays.ParseCells(lipgloss.NewStyle().Background(c).Render(" "), 1)[0].BgColor
	}
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	content := red.Render("x") + strings.Repeat(" ", 6) + red.Render("y")

	app := sizedApp(New().MustAddPage(page.New("Test")))
	own := lipgloss.Color("#123456")
	cases := []struct {
		name    string
		overlay *overlays.Overlay
		want    string
	}{
		{"bare overlay", overlays.New().SetExitKey("esc").SetWidth(20).SetHeight(5).SetContent(content), bgSeq(app.Model().Theme.Palette.Bg0)},
		{"own background", overlays.New().SetExitKey("esc").SetWidth(20).SetHeight(5).SetContent(content).SetBackground(own), bgSeq(own)},
	}
	for _, tc := range cases {
		app.HideAllOverlays()
		app.MustShowOverlay(tc.overlay)
		fgLines, fgWidth := overlays.GetLines(tc.overlay.Render(80, 24))
		x, y := overlays.CalculatePosition(tc.overlay.Position(), fgWidth, len(fgLines), 80, 24)

		lines := strings.Split(app.Render(), "\n")
		for row := y; row < y+len(fgLines); row++ {
			cells := overlays.ParseCells(lines[row], 80)
			for col := x; col < x+fgWidth; col++ {
				if cells[col].BgColor != tc.want {
					t.Fatalf("%s: row %d col %d background %q, want %q", tc.name, row, col, cells[col].BgColor, tc.want)
				}
			}
		}
	}
}

func TestOverlaysCompositeInZOrder(t *testing.T) {
	app := sizedApp(New().MustAddPage(page.New("Test")))
	low := overlays.New().SetExitKey("1").SetWidth(20).SetHeight(5).SetContent("LOW").SetZIndex(1)
	high := overlays.New().SetExitKey("2").SetWidth(20).SetHeight(5).SetContent("HIGH").SetZIndex(10)
	app.MustShowOverlay(high)
	lowID := app.MustShowOverlay(low)

	plain := plainView(app)
	if !strings.Contains(plain, "HIGH") || strings.Contains(plain, "LOW") {
		t.Errorf("the higher z-index overlay should cover the lower one:\n%s", plain)
	}

	app.BringToFront(lowID)
	plain = plainView(app)
	if !strings.Contains(plain, "LOW") || strings.Contains(plain, "HIGH") {
		t.Errorf("after BringToFront the former lower overlay should cover:\n%s", plain)
	}
}

// TestViewNavSectionClipsToBarHeight renders a bar taller than its band: the
// band keeps its configured height and drops the lines past it.
func TestViewNavSectionClipsToBarHeight(t *testing.T) {
	app := sizedApp(rootMixApp())
	app.SetNavBarAutoHide(false).SetNavBarHeight(1)
	nav := ansiSequence.ReplaceAllString(app.renderNavSection(), "")
	if lipgloss.Height(nav) != 1 {
		t.Errorf("nav band is %d lines, want 1:\n%s", lipgloss.Height(nav), nav)
	}
	if strings.Contains(nav, "nav mode") || !strings.Contains(nav, "▾ A") {
		t.Errorf("the band should keep the tabs and drop the hint:\n%s", nav)
	}
	checkFrame(t, "nav height 1", app.Render(), 80, 24)
}

// TestPageBorderOverridesTheContentBorder checks a page's own border color
// and style replace the app's while that page is active, and only then.
func TestPageBorderOverridesTheContentBorder(t *testing.T) {
	useTrueColor(t)
	pink := lipgloss.Color("#ff00aa")
	app := New()
	app.MustAddPage(page.New("Own").SetBorderColor(pink).SetBorderStyle(lipgloss.DoubleBorder()))
	app.MustAddPage(page.New("Plain"))
	app.SetShowContentBorder(true).SetNavBarAutoHide(false)
	sizedApp(app)

	corner := func() overlays.Cell {
		return overlays.ParseCells(strings.Split(app.renderContentSection(40, 6), "\n")[0], 40)[0]
	}
	fgOf := func(c color.Color) string {
		return overlays.ParseCells(lipgloss.NewStyle().Foreground(c).Render("x"), 1)[0].FgColor
	}

	if c := corner(); string(c.Char) != lipgloss.DoubleBorder().TopLeft || c.FgColor != fgOf(pink) {
		t.Errorf("the page's border: corner %q color %q, want %q in %q", c.Char, c.FgColor, lipgloss.DoubleBorder().TopLeft, fgOf(pink))
	}

	app.SetActive(1)
	theme := app.Model().Theme
	if c := corner(); string(c.Char) != theme.ContentBorder.TopLeft || c.FgColor != fgOf(theme.ContentBorderColor) {
		t.Errorf("a page without overrides: corner %q color %q, want the app's %q in %q", c.Char, c.FgColor, theme.ContentBorder.TopLeft, fgOf(theme.ContentBorderColor))
	}

	// Without the content border a page's border settings draw nothing.
	app.SetActive(0).SetShowContentBorder(false)
	if plain := ansiSequence.ReplaceAllString(app.Render(), ""); strings.Contains(plain, lipgloss.DoubleBorder().TopLeft) {
		t.Errorf("the page's border should not show with the content border off:\n%s", plain)
	}
}

// TestContentBackgroundFillsTheContentArea checks the color behind the page:
// the content background where it is set, the theme's base otherwise, and
// in both cases only under the content and the header.
func TestContentBackgroundFillsTheContentArea(t *testing.T) {
	useTrueColor(t)
	bgOf := func(c color.Color) string {
		return overlays.ParseCells(lipgloss.NewStyle().Background(c).Render(" "), 1)[0].BgColor
	}
	app := New()
	app.MustAddPage(page.New("Bare")) // no content: every content cell is a fill cell
	app.SetHeader(testHeader).SetHeaderHeight(1).SetNavBarHeight(2).SetFooterHeight(1).SetNavBarAutoHide(false)
	sizedApp(app)
	base := bgOf(app.Model().Theme.Palette.Bg0)

	rowBg := func(row int) string {
		cells := overlays.ParseCells(strings.Split(app.Render(), "\n")[row], 80)
		return cells[79].BgColor // the last column: padding in every band
	}
	if got := rowBg(10); got != base {
		t.Errorf("content row background = %q, want the theme's base %q", got, base)
	}

	custom := lipgloss.Color("#123456")
	app.SetContentBackground(custom)
	if got := rowBg(0); got != bgOf(custom) {
		t.Errorf("header row background = %q, want the content background %q", got, bgOf(custom))
	}
	if got := rowBg(10); got != bgOf(custom) {
		t.Errorf("content row background = %q, want the content background %q", got, bgOf(custom))
	}
	if got := rowBg(1); got != base {
		t.Errorf("navigation row background = %q, want the theme's base %q", got, base)
	}
	if got := rowBg(23); got != base {
		t.Errorf("footer row background = %q, want the theme's base %q", got, base)
	}

	app.SetContentBackground(nil)
	if got := rowBg(10); got != base {
		t.Errorf("content row background with none set = %q, want the theme's base %q", got, base)
	}
}

// TestSideBarTallerThanTheContentKeepsTheFooter gives the side bar more
// items than the content area has lines: the bar is cut, the footer stays.
func TestSideBarTallerThanTheContentKeepsTheFooter(t *testing.T) {
	app := New()
	for i := 0; i < 30; i++ {
		app.MustAddPage(page.New(fmt.Sprintf("page-%02d", i)))
	}
	app.SetNavBarPosition(NavBarLeft).SetNavBarAutoHide(false).SetFooterHeight(1)
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 12})

	lines := checkFrame(t, "tall side bar", app.Render(), 60, 12)
	last := ansiSequence.ReplaceAllString(lines[11], "")
	if !strings.Contains(last, "Quit") {
		t.Errorf("the last line should be the footer, got %q", last)
	}
}

// TestViewAtDegenerateSizes renders every position and layout, with every
// border on and an overlay shown, at sizes too small to hold them. The frame
// must come out at the terminal size without a panic.
func TestViewAtDegenerateSizes(t *testing.T) {
	sizes := []tea.WindowSizeMsg{{Width: 1, Height: 1}, {Width: 3, Height: 4}, {Width: 10, Height: 2}, {Width: 25, Height: 6}}
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		for _, mode := range []NavBarMode{NavBarCompact, NavBarExpanded, NavBarPath, NavBarColumns} {
			for _, size := range sizes {
				app := mixedDepthApp()
				app.SetHeader(testHeader).SetHeaderFrame(FrameDouble).SetFooterFrame(FrameFill).
					SetShowNavBarBorder(true).SetShowContentBorder(true).
					SetNavBarPosition(pos).SetNavBarMode(mode).SetNavBarAutoHide(false)
				app.Update(size)
				app.MustShowOverlay(overlays.NewModal().SetTitle("Tiny").SetExitKey("esc").Overlay)
				app.ShowNotification(overlays.Warning("squeezed"))

				label := fmt.Sprintf("pos=%v mode=%v %dx%d", pos, mode, size.Width, size.Height)
				checkFrame(t, label, app.Render(), size.Width, size.Height)
			}
		}
	}

	// With no size at all the frame is empty.
	app := mixedDepthApp()
	if got := app.Render(); got != "" {
		t.Errorf("View before any size = %q, want an empty frame", got)
	}
}
