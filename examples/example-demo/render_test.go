package main

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/tabapp"
)

var ansiSequence = regexp.MustCompile("\x1b\\[[0-9;?]*[A-Za-z]")

// allPaths walks the section tree and returns one navigation path per page.
func allPaths(root *tabapp.Section, depth int) [][]int {
	var paths [][]int
	var walk func(s *tabapp.Section, prefix []int)
	walk = func(s *tabapp.Section, prefix []int) {
		for i, child := range s.Sections() {
			walk(child, append(append([]int{}, prefix...), i))
		}
		for i := range s.Pages() {
			p := append(append([]int{}, prefix...), len(s.Sections())+i)
			for len(p) < depth {
				p = append(p, 0)
			}
			paths = append(paths, p)
		}
	}
	walk(root, nil)
	return paths
}

// TestEveryCellHasBackground renders every page and the Settings popup with
// a real color profile and checks that no cell falls through to the terminal
// background. Page content and overlay content use plain lipgloss styles with
// no background handling of their own; the framework's fills cover them.
func TestEveryCellHasBackground(t *testing.T) {
	const w, h = 80, 24
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: w, Height: h})
	for _, path := range allPaths(app.Root(), app.Root().Depth()) {
		if err := app.NavigateToPath(path); err != nil {
			t.Fatal(err)
		}
		assertCellsHaveBackground(t, fmt.Sprint(path), app.Render(), w)
	}
	press(app, tabapp.KeyPress("s"))
	assertCellsHaveBackground(t, "settings popup", app.Render(), w)
}

func assertCellsHaveBackground(t *testing.T, label, view string, w int) {
	t.Helper()
	for row, line := range strings.Split(view, "\n") {
		for col, c := range overlays.ParseCells(line, w) {
			if !c.HasBackground() {
				t.Fatalf("%s: row %d col %d has no background", label, row, col)
			}
		}
	}
}

// TestSettingsRowsAlign checks the popup's first and last rows start in the
// same column: the last line of a centered block keeps its trailing padding.
func TestSettingsRowsAlign(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	press(app, tabapp.KeyPress("s"))
	press(app, tabapp.KeyPress("down")) // move the pill off the first row

	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	first, last := -1, -1
	for _, line := range strings.Split(plain, "\n") {
		if i := strings.Index(line, "Theme "); i >= 0 && first < 0 {
			first = i
		}
		if i := strings.Index(line, "Key glow"); i >= 0 {
			last = i
		}
	}
	if first < 0 || last < 0 {
		t.Fatalf("rows not found (first=%d last=%d)", first, last)
	}
	if first != last {
		t.Errorf("first row starts at column %d, last row at %d", first, last)
	}
}

// frameFits checks that a rendered view is exactly w x h cells.
func frameFits(t *testing.T, label, view string, w, h int) string {
	t.Helper()
	plain := ansiSequence.ReplaceAllString(view, "")
	lines := strings.Split(plain, "\n")
	if len(lines) != h {
		t.Errorf("%s: %d lines, want %d", label, len(lines), h)
	}
	for i, l := range strings.Split(view, "\n") {
		if lw := lipgloss.Width(l); lw != w {
			t.Errorf("%s line %d: width %d, want %d", label, i, lw, w)
			break
		}
	}
	return plain
}

// TestRenderSettingsPopup opens the popup over a Catalog page at two sizes.
func TestRenderSettingsPopup(t *testing.T) {
	for _, size := range []struct{ w, h int }{{80, 24}, {120, 40}} {
		app := buildApp()
		app.Update(tea.WindowSizeMsg{Width: size.w, Height: size.h})
		if err := app.NavigateToPath([]int{1, 1, 0, 0}); err != nil {
			t.Fatal(err)
		}
		press(app, tabapp.KeyPress("s"))
		press(app, tabapp.KeyPress("down"))
		press(app, tabapp.KeyPress("down"))

		press(app, tabapp.KeyPress("down"))
		plain := frameFits(t, "popup", app.Render(), size.w, size.h)
		if !strings.Contains(plain, "Nav bar mode") || !strings.Contains(plain, "◂ Compact ▸") {
			t.Error("popup should show the selected Nav bar mode row")
		}
		if testing.Verbose() {
			t.Logf("\n%dx%d settings popup\n%s", size.w, size.h, plain)
		}
	}
}

// TestRenderWithAllBorders turns every border on and renders a page and the
// popup at 80x24.
func TestRenderWithAllBorders(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	app.SetShowHeaderBorder(true).SetShowNavBarBorder(true).SetShowFooterBorder(true).SetShowContentBorder(true).SetNavBarAutoHide(false)
	if err := app.NavigateToPath([]int{1, 1, 0, 0}); err != nil {
		t.Fatal(err)
	}
	plain := frameFits(t, "bordered page", app.Render(), 80, 24)
	if testing.Verbose() {
		t.Logf("\n80x24 all borders\n%s", plain)
	}
	press(app, tabapp.KeyPress("s"))
	plain = frameFits(t, "bordered popup", app.Render(), 80, 24)
	if testing.Verbose() {
		t.Logf("\n80x24 all borders + popup\n%s", plain)
	}
}

// TestRenderEveryPage renders each page at a small and a large terminal size
// and checks that the frame fills the screen and names the page.
func TestRenderEveryPage(t *testing.T) {
	for _, size := range []struct{ w, h int }{{80, 24}, {120, 40}} {
		app := buildApp()
		app.Update(tea.WindowSizeMsg{Width: size.w, Height: size.h})
		for _, path := range allPaths(app.Root(), app.Root().Depth()) {
			if err := app.NavigateToPath(path); err != nil {
				t.Fatalf("%v: %v", path, err)
			}
			view := ansiSequence.ReplaceAllString(app.Render(), "")
			lines := strings.Split(view, "\n")
			if len(lines) != size.h {
				t.Errorf("%dx%d %v: %d lines, want %d", size.w, size.h, path, len(lines), size.h)
			}
			if testing.Verbose() {
				t.Logf("\n%dx%d path %v\n%s", size.w, size.h, path, view)
			}
		}
	}
}

// TestRenderBandFramesAndHeights renders every header and footer frame at
// several band heights and checks that the frame fills the screen.
func TestRenderBandFramesAndHeights(t *testing.T) {
	const w, h = 80, 24
	for _, frame := range bandFrames {
		for _, height := range []int{0, 1, 3, bandHeightMax} {
			app := buildApp()
			app.Update(tea.WindowSizeMsg{Width: w, Height: h})
			app.Update(tabapp.SetHeaderFrameMsg{Frame: frame})
			app.Update(tabapp.SetFooterFrameMsg{Frame: frame})
			app.Update(tabapp.SetHeaderHeightMsg{Height: height})
			app.Update(tabapp.SetFooterHeightMsg{Height: height})

			label := fmt.Sprintf("%v frame, height %d", frame, height)
			plain := frameFits(t, label, app.Render(), w, h)
			if shown := strings.Contains(plain, "tuiframework demo"); shown != (height > 0) {
				t.Errorf("%s: header shown = %v, want %v", label, shown, height > 0)
			}
			if shown := strings.Contains(plain, "q Quit"); shown != (height > 0) {
				t.Errorf("%s: footer keys shown = %v, want %v", label, shown, height > 0)
			}
			if testing.Verbose() {
				t.Logf("\n%dx%d %s\n%s", w, h, label, plain)
			}
		}
	}
}

// TestFilledBandsReverseTheirText renders the header and footer as filled
// bars and checks every cell of both: the fill color behind it and, where
// there is text, the fill's text color.
func TestFilledBandsReverseTheirText(t *testing.T) {
	const w, h = 80, 24
	for _, theme := range demoThemes {
		app := buildApp()
		app.Update(tea.WindowSizeMsg{Width: w, Height: h})
		app.Update(tabapp.SetThemeMsg{Theme: theme})
		app.Update(tabapp.SetHeaderFrameMsg{Frame: tabapp.FrameFill})
		app.Update(tabapp.SetFooterFrameMsg{Frame: tabapp.FrameFill})

		lines := strings.Split(app.Render(), "\n")
		bands := []struct {
			name string
			rows []string
			fill lipgloss.Style
		}{
			{"header", lines[:headerHeight], theme.HeaderFill},
			{"footer", lines[h-footerHeight:], theme.FooterFill},
		}
		for _, band := range bands {
			bg := overlays.ParseCells(band.fill.Render(" "), 1)[0]
			for row, line := range band.rows {
				for col, c := range overlays.ParseCells(line, w) {
					if c.BgColor != bg.BgColor {
						t.Fatalf("%s %s row %d col %d: background %q, want the fill %q", theme.Name(), band.name, row, col, c.BgColor, bg.BgColor)
					}
					if c.Char != ' ' && c.Char != 0 && c.FgColor != bg.FgColor {
						t.Fatalf("%s %s row %d col %d (%q): foreground %q, want the fill's %q", theme.Name(), band.name, row, col, c.Char, c.FgColor, bg.FgColor)
					}
				}
			}
		}
		assertCellsHaveBackground(t, theme.Name()+" filled bands", app.Render(), w)
	}
}
