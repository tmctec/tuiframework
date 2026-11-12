package overlays

import (
	"image/color"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/styles"
)

func TestNewModalFollowsTheme(t *testing.T) {
	m := NewModal()
	if !m.FollowsTheme() {
		t.Error("NewModal should follow the theme")
	}
	if !m.IsBlocking() {
		t.Error("NewModal should block")
	}
	if m.ZIndex() != 0 {
		t.Errorf("ZIndex() = %d, want 0", m.ZIndex())
	}
}

func TestModalApplyThemeSetsFrameAndTitle(t *testing.T) {
	for _, theme := range []*styles.Theme{styles.Gruvbox, styles.Retro, styles.Nord} {
		m := NewModal().ApplyTheme(theme)
		s := m.Style()

		if s.GetBackground() != theme.ModalBackground {
			t.Errorf("%s: background = %v, want %v", theme.Name(), s.GetBackground(), theme.ModalBackground)
		}
		if s.GetForeground() != theme.Palette.Fg1 {
			t.Errorf("%s: foreground = %v, want %v", theme.Name(), s.GetForeground(), theme.Palette.Fg1)
		}
		if s.GetBorderStyle() != theme.ModalBorder {
			t.Errorf("%s: border style does not match the theme", theme.Name())
		}
		if s.GetBorderTopForeground() != theme.ModalBorderColor {
			t.Errorf("%s: border color = %v, want %v", theme.Name(), s.GetBorderTopForeground(), theme.ModalBorderColor)
		}
		if s.GetBorderTopBackground() != theme.ModalBackground {
			t.Errorf("%s: border background = %v, want %v", theme.Name(), s.GetBorderTopBackground(), theme.ModalBackground)
		}
		if s.GetPaddingTop() != 1 || s.GetPaddingLeft() != 2 {
			t.Errorf("%s: padding = %d,%d, want 1,2", theme.Name(), s.GetPaddingTop(), s.GetPaddingLeft())
		}
		if m.titleStyle.GetForeground() != theme.ModalTitle.GetForeground() || m.titleStyle.GetBold() != theme.ModalTitle.GetBold() {
			t.Errorf("%s: title style should come from the theme", theme.Name())
		}
	}
}

func TestModalSetStyleReplacesFrame(t *testing.T) {
	custom := lipgloss.NewStyle().Background(lipgloss.Color("#123456")).Border(lipgloss.ThickBorder())
	m := NewModal().SetStyle(custom)

	if m.FollowsTheme() {
		t.Error("SetStyle should stop the modal from following the theme")
	}
	if m.Style().GetBackground() != lipgloss.Color("#123456") || m.Style().GetBorderStyle() != lipgloss.ThickBorder() {
		t.Error("SetStyle should replace the frame")
	}
}

func TestModalSettersReturnSelf(t *testing.T) {
	m := NewModal()
	results := []*Modal{
		m.SetTitle("t"),
		m.SetContent(func(_ *flexbox.FlexBox) {}),
		m.SetSize(50),
		m.SetPosition(PositionTopLeft),
		m.SetExitKey("esc"),
		m.SetBlocking(false),
		m.SetStyle(lipgloss.NewStyle()),
		m.RegisterKey("y", testModalKeyMsg{key: "y"}),
		m.RegisterKeys(map[string]tea.Msg{"n": testModalKeyMsg{key: "n"}}),
		m.SetZIndex(4),
		m.ApplyTheme(styles.Dracula),
	}
	for i, r := range results {
		if r != m {
			t.Errorf("setter %d did not return the receiver", i)
		}
	}

	if m.Title() != "t" || m.Size() != 50 || m.Position() != PositionTopLeft || m.ExitKey() != "esc" {
		t.Error("setters should apply to the embedded overlay")
	}
	if m.IsBlocking() {
		t.Error("SetBlocking(false) should apply")
	}
	if !m.HasKey("y") || !m.HasKey("n") {
		t.Error("RegisterKey and RegisterKeys should apply")
	}
	if m.ZIndex() != 4 {
		t.Errorf("ZIndex() = %d, want 4", m.ZIndex())
	}
}

func TestModalSetZIndex(t *testing.T) {
	m := NewModal().SetZIndex(10)
	if m.ZIndex() != 10 {
		t.Errorf("ZIndex() = %d, want 10", m.ZIndex())
	}
	base := m.Overlay
	if base.ZIndex() != 10 {
		t.Error("SetZIndex should write through to the embedded overlay")
	}
}

// TestModalBorderColorSurvivesApplyTheme covers the override SetBorderColor
// keeps: a later theme restyles the modal and leaves the border color.
func TestModalBorderColorSurvivesApplyTheme(t *testing.T) {
	red := lipgloss.Color("#ff0000")
	m := NewModal().SetBorderColor(red)
	if got := m.Style().GetBorderTopForeground(); got != red {
		t.Fatalf("border = %v, want the color set", got)
	}

	m.ApplyTheme(styles.Nord)
	if got := m.Style().GetBorderTopForeground(); got != red {
		t.Errorf("border after ApplyTheme = %v, want the override kept", got)
	}
	if got := m.Style().GetBackground(); got != styles.Nord.ModalBackground {
		t.Errorf("background = %v, want the theme's", got)
	}
}

// TestModalTakesTheThemeBorderBlend covers a theme with a border gradient:
// the modal's border takes it, a border color of the modal's own replaces
// it, and a theme without a gradient leaves none.
func TestModalTakesTheThemeBorderBlend(t *testing.T) {
	p := styles.NordPalette
	p.BorderBlend = [2]color.Color{lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")}
	blend := styles.NewTheme(p)

	m := NewModal().ApplyTheme(blend)
	if got := m.Style().GetBorderForegroundBlend(); len(got) != 2 {
		t.Fatalf("blend = %v, want the theme's two colors", got)
	}

	m.ApplyTheme(styles.Nord)
	if got := m.Style().GetBorderForegroundBlend(); len(got) != 0 {
		t.Errorf("a theme without a gradient leaves blend = %v", got)
	}

	own := NewModal().SetBorderColor(lipgloss.Color("#00ff00")).ApplyTheme(blend)
	if got := own.Style().GetBorderForegroundBlend(); len(got) != 0 {
		t.Errorf("the modal's own border color should replace the gradient, blend = %v", got)
	}
}
