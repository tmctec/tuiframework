package overlays

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/styles"
)

func TestOverlayRenderIncludesTitle(t *testing.T) {
	o := New().SetTitle("Hello Title").SetContent("body").SetExitKey("esc")

	rendered := o.Render(100, 50)

	if !strings.Contains(rendered, "Hello Title") {
		t.Error("rendered overlay should contain the title")
	}
	if !strings.Contains(rendered, "body") {
		t.Error("rendered overlay should contain the content")
	}
	if o.Title() != "Hello Title" {
		t.Errorf("Title() = %q", o.Title())
	}
}

func TestOverlayRenderWithoutTitleKeepsHeight(t *testing.T) {
	with := New().SetTitle("T").SetContent("x").SetExitKey("esc").Render(60, 30)
	without := New().SetContent("x").SetExitKey("esc").Render(60, 30)
	if lipgloss.Height(with) != lipgloss.Height(without) {
		t.Errorf("title changes overlay height: %d vs %d", lipgloss.Height(with), lipgloss.Height(without))
	}
}

func TestModalApplyTheme(t *testing.T) {
	m := NewModal().ApplyTheme(styles.Gruvbox)

	if got := m.Style().GetBorderTopForeground(); got != styles.Gruvbox.ModalBorderColor {
		t.Errorf("border = %v, want %v", got, styles.Gruvbox.ModalBorderColor)
	}
	if got := m.Style().GetBackground(); got != styles.Gruvbox.ModalBackground {
		t.Errorf("background = %v, want %v", got, styles.Gruvbox.ModalBackground)
	}
	if got := m.Style().GetForeground(); got != styles.Gruvbox.Palette.Fg1 {
		t.Errorf("foreground = %v, want %v", got, styles.Gruvbox.Palette.Fg1)
	}
}

func TestNewModalDefaultsToDracula(t *testing.T) {
	m := NewModal()
	if got := m.Style().GetBorderTopForeground(); got != styles.Dracula.ModalBorderColor {
		t.Errorf("border = %v, want %v", got, styles.Dracula.ModalBorderColor)
	}
	if !m.IsBlocking() {
		t.Error("modal should block by default")
	}
}

func TestOverlayFixedSizeClampsToContainer(t *testing.T) {
	o := New().SetWidth(200).SetHeight(100).SetExitKey("esc")
	w, h := o.CalculateDimensions(80, 24)
	if w != 80 || h != 24 {
		t.Errorf("dimensions = %dx%d, want 80x24", w, h)
	}
}

func TestOverlayContentFuncRunsEveryRender(t *testing.T) {
	calls := 0
	o := New().SetExitKey("esc").SetContentFunc(func(fb *flexbox.FlexBox) {
		calls++
		row := fb.NewRow()
		cell := flexbox.NewCell(1, 1)
		cell.SetContent("call")
		row.AddCells(cell)
		fb.AddRows([]*flexbox.Row{row})
	})
	if calls != 1 {
		t.Fatalf("SetContentFunc should build once up front, got %d calls", calls)
	}
	o.Render(80, 24)
	o.Render(80, 24)
	if calls != 3 {
		t.Errorf("content func ran %d times, want 3 (build + 2 renders)", calls)
	}
	// 80% of 80 columns is 64; the border takes 2 and the padding 4.
	if o.Flexbox() == nil || o.Flexbox().GetWidth() != 58 {
		t.Errorf("render should hand the func a flexbox sized to the content area, got width %d", o.Flexbox().GetWidth())
	}
}

func TestOverlayRenderMatchesComputedDimensions(t *testing.T) {
	needsNarrowBorders(t)
	o := New().SetWidth(40).SetHeight(10).SetContent("x").SetExitKey("esc")
	rendered := o.Render(80, 24)
	if w, h := lipgloss.Width(rendered), lipgloss.Height(rendered); w != 40 || h != 10 {
		t.Errorf("rendered %dx%d, want the configured 40x10 including the border", w, h)
	}
	titled := New().SetWidth(40).SetHeight(10).SetTitle("T").SetContent("x").SetExitKey("esc").Render(80, 24)
	if w, h := lipgloss.Width(titled), lipgloss.Height(titled); w != 40 || h != 10 {
		t.Errorf("titled rendered %dx%d, want 40x10", w, h)
	}
}

func TestModalApplyThemeThenSetStyle(t *testing.T) {
	m := NewModal()
	if !m.FollowsTheme() {
		t.Error("NewModal should follow the theme")
	}
	m.SetStyle(lipgloss.NewStyle())
	if m.FollowsTheme() {
		t.Error("SetStyle should stop the modal from following the theme")
	}
	m.ApplyTheme(styles.Nord)
	if !m.FollowsTheme() {
		t.Error("ApplyTheme should make the modal follow the theme again")
	}
}

func TestModalTitleRendersThroughOverlay(t *testing.T) {
	m := NewModal().SetTitle("Modal Title").SetExitKey("esc")
	m.Overlay.SetContent("content")

	rendered := m.Render(80, 24)

	if !strings.Contains(rendered, "Modal Title") {
		t.Error("the embedded Overlay should render the modal title")
	}
}
