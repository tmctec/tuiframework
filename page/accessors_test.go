package page_test

import (
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
)

type pingMsg struct{}

// TestZeroValuePageIsUsable binds keys and footer items on a Page that was
// not created with New; the maps are created on first use.
func TestZeroValuePageIsUsable(t *testing.T) {
	var p page.Page
	p.BindKey("a", "Act", pingMsg{}).AddFooterItem("q", "Quit")

	if !p.HasKey("a") || p.Bindings()["a"].Label() != "Act" {
		t.Errorf("binding = %+v, want a bound to Act", p.Bindings()["a"])
	}
	if p.Footer()["q"] != "Quit" {
		t.Errorf("footer = %v, want q: Quit", p.Footer())
	}
	if p.Name() != "" || p.Content() != nil {
		t.Errorf("zero page has name %q and content %v, want none", p.Name(), p.Content() != nil)
	}
}

// TestFooterSettersMerge checks SetFooter replaces the map and AddFooterItem
// adds to it, including after SetFooter(nil).
func TestFooterSettersMerge(t *testing.T) {
	p := page.New("Footer").SetFooter(map[string]string{"a": "A"}).AddFooterItem("b", "B")
	if len(p.Footer()) != 2 || p.Footer()["a"] != "A" || p.Footer()["b"] != "B" {
		t.Errorf("footer = %v, want a: A and b: B", p.Footer())
	}

	p.SetFooter(nil).AddFooterItem("c", "C")
	if len(p.Footer()) != 1 || p.Footer()["c"] != "C" {
		t.Errorf("footer after SetFooter(nil) = %v, want only c: C", p.Footer())
	}
}

// TestBorderOverrides checks the per-page border color and style are nil
// until set and returned afterwards.
func TestBorderOverrides(t *testing.T) {
	p := page.New("Borders")
	if p.BorderColor() != nil || p.BorderStyle() != nil {
		t.Fatal("a new page has no border overrides")
	}

	p.SetBorderColor(lipgloss.Color("#123456")).SetBorderStyle(lipgloss.DoubleBorder())
	if c := p.BorderColor(); c != lipgloss.Color("#123456") {
		t.Errorf("BorderColor = %v, want #123456", c)
	}
	if b := p.BorderStyle(); b == nil || b.Top != lipgloss.DoubleBorder().Top {
		t.Errorf("BorderStyle = %v, want the double border", b)
	}
}

// TestBindingLifecycle binds, rebinds, reads and unbinds a key.
func TestBindingLifecycle(t *testing.T) {
	p := page.New("Keys").BindKey("x", "First", pingMsg{})
	kb := p.Bindings()["x"]
	if kb.Key() != "x" || kb.Label() != "First" {
		t.Errorf("binding = key %q label %q, want x and First", kb.Key(), kb.Label())
	}
	if _, ok := kb.Msg().(pingMsg); !ok {
		t.Errorf("Msg = %T, want pingMsg", kb.Msg())
	}

	p.BindKey("x", "Second", pingMsg{})
	if p.Bindings()["x"].Label() != "Second" {
		t.Error("rebinding a key should replace its label")
	}

	p.UnbindKey("x").UnbindKey("never-bound")
	if p.HasKey("x") || len(p.Bindings()) != 0 {
		t.Errorf("bindings after unbind = %v, want none", p.Bindings())
	}
}

// TestContentFunctionReceivesTheArguments passes the model and the band
// heights through to the content function unchanged.
func TestContentFunctionReceivesTheArguments(t *testing.T) {
	var gotModel interface{}
	var gotNav, gotFooter int
	p := page.New("Args").SetContent(func(_ *flexbox.FlexBox, model interface{}, navBarHeight, footerHeight int) {
		gotModel, gotNav, gotFooter = model, navBarHeight, footerHeight
	})

	p.Content()(flexbox.New(10, 2), "the model", 3, 2)
	if gotModel != "the model" || gotNav != 3 || gotFooter != 2 {
		t.Errorf("content received %v %d %d, want the model, 3, 2", gotModel, gotNav, gotFooter)
	}
}
