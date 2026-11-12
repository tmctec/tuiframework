package tabapp

import (
	"fmt"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

func testHeader(fb *flexbox.FlexBox, _ interface{}) {
	row := fb.NewRow()
	cell := flexbox.NewCell(1, 1)
	cell.SetContent("HEADER BAND")
	row.AddCells(cell)
	fb.AddRows([]*flexbox.Row{row})
}

// TestViewFitsTerminalWithAllBorders renders every border, position, mode and
// navigation-order combination with a header and checks the frame is exactly
// the terminal size and every band is on screen.
func TestViewFitsTerminalWithAllBorders(t *testing.T) {
	const w, h = 80, 24
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		for _, mode := range []NavBarMode{NavBarCompact, NavBarExpanded, NavBarPath, NavBarColumns} {
			for _, navOn := range []bool{false, true} {
				below := navOn // also covers both navigation orders
				label := fmt.Sprintf("pos=%v mode=%v navOn=%v navBelow=%v", pos, mode, navOn, below)
				app := rootMixApp()
				app.SetHeader(testHeader).
					SetShowHeaderBorder(true).
					SetShowNavBarBorder(true).
					SetShowFooterBorder(true).
					SetShowContentBorder(true).
					SetNavBarAutoHide(false).
					SetNavBarPosition(pos).
					SetNavBarMode(mode).
					SetNavBarBelowHeader(below)
				app.Update(tea.WindowSizeMsg{Width: w, Height: h})
				if navOn {
					app.Update(KeyPress("tab"))
				}

				lines := strings.Split(app.Render(), "\n")
				if len(lines) != h {
					t.Errorf("%s: %d lines, want %d", label, len(lines), h)
				}
				for i, l := range lines {
					if lw := lipgloss.Width(l); lw != w {
						t.Errorf("%s line %d: width %d, want %d", label, i, lw, w)
						break
					}
				}
				plain := ansiSequence.ReplaceAllString(app.Render(), "")
				for _, want := range []string{"HEADER BAND", "Quit", "a1"} {
					if !strings.Contains(plain, want) {
						t.Errorf("%s: %q missing from view", label, want)
					}
				}
			}
		}
	}
}

func TestNavBarBelowHeaderOrder(t *testing.T) {
	app := rootMixApp()
	app.SetHeader(testHeader).SetNavBarAutoHide(false)
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	const tabs = "▾ A"
	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	if h, n := strings.Index(plain, "HEADER BAND"), strings.Index(plain, tabs); h < 0 || n < 0 || h > n {
		t.Errorf("header should come before the tabs when the nav bar sits below it (header=%d tabs=%d)", h, n)
	}

	app.SetNavBarBelowHeader(false)
	plain = ansiSequence.ReplaceAllString(app.Render(), "")
	if h, n := strings.Index(plain, "HEADER BAND"), strings.Index(plain, tabs); h < 0 || n < 0 || h < n {
		t.Errorf("tabs should come before the header when the nav bar sits above it (header=%d tabs=%d)", h, n)
	}
}

func TestContentBorderKeepsFooterVisible(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Only").SetFooter(map[string]string{"z": "Marker"}))
	app.SetShowContentBorder(true).SetNavBarAutoHide(false)
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 20})

	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	lines := strings.Split(plain, "\n")
	if len(lines) != 20 {
		t.Fatalf("%d lines, want 20", len(lines))
	}
	if !strings.Contains(plain, "Marker") {
		t.Error("footer text should be visible with the content border on")
	}
	// The border's right edge is on screen: the first content row ends with a border rune.
	for _, l := range lines {
		if strings.HasPrefix(strings.TrimSpace(l), "╭") || strings.HasPrefix(strings.TrimSpace(l), "┌") {
			trimmed := strings.TrimRight(l, " ")
			if !strings.HasSuffix(trimmed, "╮") && !strings.HasSuffix(trimmed, "┐") {
				t.Errorf("top border line is cut off: %q", trimmed)
			}
			break
		}
	}
}

func TestModalFollowsThemeChange(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	modal := overlays.NewModal().SetTitle("Live").SetExitKey("esc")
	if _, err := app.ShowModal(modal); err != nil {
		t.Fatal(err)
	}

	app.UseTheme(styles.Gruvbox)

	if got := modal.Style().GetBorderTopForeground(); got != styles.Gruvbox.ModalBorderColor {
		t.Errorf("border = %v, want Gruvbox %v", got, styles.Gruvbox.ModalBorderColor)
	}
}

func TestStyledModalIgnoresThemeChange(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	custom := lipgloss.NewStyle().BorderForeground(lipgloss.Color("#123456")).Border(lipgloss.NormalBorder())
	modal := overlays.NewModal().SetStyle(custom).SetExitKey("esc")
	if _, err := app.ShowModal(modal); err != nil {
		t.Fatal(err)
	}

	app.UseTheme(styles.Gruvbox)

	if got := modal.Style().GetBorderTopForeground(); got != lipgloss.Color("#123456") {
		t.Errorf("border = %v, want the custom #123456", got)
	}
}
