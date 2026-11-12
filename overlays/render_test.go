package overlays

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
)

// boxWith returns a flexbox holding a single cell with the given text.
func boxWith(text string) *flexbox.FlexBox {
	fb := flexbox.New(0, 0)
	row := fb.NewRow()
	c := flexbox.NewCell(1, 1)
	c.SetContent(text)
	row.AddCells(c)
	fb.AddRows([]*flexbox.Row{row})
	return fb
}

func TestOverlayRenderPercentSizeMatchesDimensions(t *testing.T) {
	needsNarrowBorders(t)
	o := New().SetContent("Hello").SetExitKey("esc")

	rendered := o.Render(80, 24)
	wantW, wantH := o.CalculateDimensions(80, 24)
	if w, h := lipgloss.Width(rendered), lipgloss.Height(rendered); w != wantW || h != wantH {
		t.Errorf("rendered %dx%d, want %dx%d", w, h, wantW, wantH)
	}
	if !strings.Contains(rendered, "Hello") {
		t.Error("rendered overlay should contain the content string")
	}
}

func TestOverlayRenderWithFlexbox(t *testing.T) {
	needsNarrowBorders(t)
	fb := boxWith("Boxed")
	o := New().SetWidth(40).SetHeight(10).SetFlexbox(fb).SetExitKey("esc")

	if o.Flexbox() != fb {
		t.Fatal("Flexbox() should return the flexbox passed to SetFlexbox")
	}

	rendered := o.Render(80, 24)
	if !strings.Contains(rendered, "Boxed") {
		t.Error("rendered overlay should contain the flexbox content")
	}
	// 40 wide less a 2-cell border and 4 cells of padding; 10 high less 2 and 2.
	if fb.GetWidth() != 34 || fb.GetHeight() != 6 {
		t.Errorf("flexbox sized %dx%d, want the content area 34x6", fb.GetWidth(), fb.GetHeight())
	}
	if w, h := lipgloss.Width(rendered), lipgloss.Height(rendered); w != 40 || h != 10 {
		t.Errorf("rendered %dx%d, want 40x10", w, h)
	}
}

func TestOverlayRenderFlexboxBeatsContentString(t *testing.T) {
	o := New().SetContent("plain").SetFlexbox(boxWith("Boxed")).SetExitKey("esc")

	rendered := o.Render(80, 24)
	if !strings.Contains(rendered, "Boxed") || strings.Contains(rendered, "plain") {
		t.Error("a flexbox should take precedence over the content string")
	}
	if o.Content() != "plain" {
		t.Errorf("Content() = %q, want the string kept", o.Content())
	}
}

func TestOverlayRenderContentFuncBeatsFlexbox(t *testing.T) {
	o := New().SetFlexbox(boxWith("Static")).SetExitKey("esc")
	o.SetContentFunc(func(fb *flexbox.FlexBox) {
		row := fb.NewRow()
		c := flexbox.NewCell(1, 1)
		c.SetContent("Dynamic")
		row.AddCells(c)
		fb.AddRows([]*flexbox.Row{row})
	})

	rendered := o.Render(80, 24)
	if !strings.Contains(rendered, "Dynamic") || strings.Contains(rendered, "Static") {
		t.Error("the content func should take precedence over a static flexbox")
	}
}

func TestOverlayRenderTitleTakesTwoContentRows(t *testing.T) {
	build := func(title string) *Overlay {
		return New().SetWidth(40).SetHeight(10).SetTitle(title).SetExitKey("esc").
			SetContentFunc(func(_ *flexbox.FlexBox) {})
	}

	plain := build("")
	plain.Render(80, 24)
	if h := plain.Flexbox().GetHeight(); h != 6 {
		t.Errorf("untitled content height = %d, want 6", h)
	}

	titled := build("Title")
	titled.Render(80, 24)
	if h := titled.Flexbox().GetHeight(); h != 4 {
		t.Errorf("titled content height = %d, want 4 (two rows for the title)", h)
	}
	if w := titled.Flexbox().GetWidth(); w != 34 {
		t.Errorf("titled content width = %d, want 34 (unchanged by the title)", w)
	}
}

func TestOverlayRenderTitleSitsAboveContent(t *testing.T) {
	rendered := New().SetWidth(30).SetHeight(8).SetTitle("Head").SetContent("body").SetExitKey("esc").Render(80, 24)
	lines := strings.Split(rendered, "\n")

	titleRow, bodyRow := -1, -1
	for i, line := range lines {
		if strings.Contains(line, "Head") {
			titleRow = i
		}
		if strings.Contains(line, "body") {
			bodyRow = i
		}
	}
	if titleRow < 0 || bodyRow < 0 {
		t.Fatalf("title row %d, body row %d; both should render", titleRow, bodyRow)
	}
	if bodyRow != titleRow+2 {
		t.Errorf("body on row %d, title on row %d; want one blank row between them", bodyRow, titleRow)
	}
	between := strings.Trim(lines[titleRow+1], "│ ")
	if between != "" {
		t.Errorf("row between title and body = %q, want blank", between)
	}
}

func TestOverlayRenderUsesTitleStyle(t *testing.T) {
	o := New().SetWidth(30).SetHeight(8).SetTitle("Head").SetContent("body").SetExitKey("esc").
		SetTitleStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff0000")))

	if !strings.Contains(o.Render(80, 24), "\x1b[38;2;255;0;0mHead") {
		t.Error("the title should render with the title style")
	}
}

func TestNotificationRenderUsesConfiguredSize(t *testing.T) {
	needsNarrowBorders(t)
	saved := *NotificationConfig
	defer func() { *NotificationConfig = saved }()
	NotificationConfig.SetWidth(50).SetHeight(3)

	rendered := Info("Saved").Render(120, 40)
	if w, h := lipgloss.Width(rendered), lipgloss.Height(rendered); w != 50 || h != 3 {
		t.Errorf("rendered %dx%d, want the configured 50x3", w, h)
	}
	if !strings.Contains(rendered, "Saved") {
		t.Error("rendered notification should contain the message")
	}
}

// TestOverlayRenderCacheFollowsEverySetter renders an overlay, so its text
// is cached, changes it with one setter, and renders again. The result must
// equal that of an overlay configured the same way and never rendered
// before: no setter may leave a stale frame behind.
func TestOverlayRenderCacheFollowsEverySetter(t *testing.T) {
	red, blue := lipgloss.Color("#ff0000"), lipgloss.Color("#0000ff")
	setters := map[string]func(*Overlay){
		"SetContent":          func(o *Overlay) { o.SetContent("other text") },
		"SetSize":             func(o *Overlay) { o.SetSize(40) },
		"SetWidthPct":         func(o *Overlay) { o.SetWidthPct(30) },
		"SetHeightPct":        func(o *Overlay) { o.SetHeightPct(30) },
		"SetWidth":            func(o *Overlay) { o.SetWidth(24) },
		"SetHeight":           func(o *Overlay) { o.SetHeight(7) },
		"SetStyle":            func(o *Overlay) { o.SetStyle(lipgloss.NewStyle().Border(lipgloss.DoubleBorder())) },
		"SetBorder":           func(o *Overlay) { o.SetBorder(lipgloss.ThickBorder()) },
		"SetBorderForeground": func(o *Overlay) { o.SetBorderForeground(red) },
		"SetBorderBlend":      func(o *Overlay) { o.SetBorderBlend(red, blue) },
		"SetBorderBackground": func(o *Overlay) { o.SetBorderBackground(blue) },
		"SetPadding":          func(o *Overlay) { o.SetPadding(0, 5) },
		"SetBackground":       func(o *Overlay) { o.SetBackground(blue) },
		"SetForeground":       func(o *Overlay) { o.SetForeground(red) },
		"SetTitle":            func(o *Overlay) { o.SetTitle("Other") },
		"SetTitleStyle":       func(o *Overlay) { o.SetTitleStyle(lipgloss.NewStyle().Foreground(red)) },
	}
	base := func() *Overlay { return New().SetTitle("Title").SetContent("body").SetExitKey("esc") }

	for name, set := range setters {
		t.Run(name, func(t *testing.T) {
			rendered := base()
			before := rendered.Render(80, 24)
			if rendered.Render(80, 24) != before {
				t.Fatal("an unchanged overlay should render the same text")
			}
			set(rendered)

			fresh := base()
			set(fresh)
			want := fresh.Render(80, 24)
			if got := rendered.Render(80, 24); got != want {
				t.Errorf("the render after %s is stale:\n%s\nwant:\n%s", name, got, want)
			}
			if want == before {
				t.Errorf("%s should change the render; the test would not notice a stale one", name)
			}
		})
	}

	// A different container is a different render.
	o := base()
	if o.Render(80, 24) == o.Render(120, 40) {
		t.Error("the render should follow the container size")
	}

	// Content that can change between frames is rendered every time.
	calls := 0
	dynamic := base().SetContentFunc(func(*flexbox.FlexBox) { calls++ })
	calls = 0
	dynamic.Render(80, 24)
	dynamic.Render(80, 24)
	if calls != 2 {
		t.Errorf("a content function ran %d times for two renders, want 2", calls)
	}
	withBox := base().SetFlexbox(flexbox.New(10, 3))
	if withBox.Render(80, 24) == "" || withBox.cache.valid {
		t.Error("an overlay with a flexbox should render and never cache")
	}
}
