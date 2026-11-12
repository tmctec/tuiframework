package tabapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
)

// TestOverlayCellsKeepOverlayBackground renders a modal whose content has
// plain spaces after a styled segment, the pattern that otherwise lets the
// terminal background show through, and checks every overlay cell has a
// background color.
func TestOverlayCellsKeepOverlayBackground(t *testing.T) {
	// Two styled segments with plain spaces between them: the inner color
	// reset leaves those spaces without a background.
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000"))
	content := red.Render("x") + strings.Repeat(" ", 10) + red.Render("y")

	modal := overlays.NewModal().SetTitle("Fill").SetExitKey("esc")
	modal.Overlay.SetWidth(40).SetHeight(8).SetContent(content)
	plain := overlays.New().SetExitKey("esc").SetContent(content) // no background of its own

	for name, o := range map[string]*overlays.Overlay{"themed modal": modal.Overlay, "bare overlay": plain} {
		app := New()
		app.MustAddPage(page.New("Test"))
		app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
		if _, err := app.ShowOverlay(o); err != nil {
			t.Fatal(err)
		}

		fgLines, fgWidth := overlays.GetLines(o.Render(80, 24))
		x, y := overlays.CalculatePosition(o.Position(), fgWidth, len(fgLines), 80, 24)

		lines := strings.Split(app.Render(), "\n")
		for row := y; row < y+len(fgLines); row++ {
			cells := overlays.ParseCells(lines[row], 80)
			for col := x; col < x+fgWidth; col++ {
				if !cells[col].HasBackground() {
					t.Fatalf("%s: row %d col %d inside the overlay has no background", name, row, col)
				}
			}
		}
	}
}
