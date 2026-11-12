package tabapp

import (
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/page"
)

// fallbackFooterText is the footer shown when no binding has a label.
const fallbackFooterText = "Tab: Navigate • Ctrl+Q: Quit"

// plainFooter renders a page's footer without ANSI sequences.
func plainFooter(app *TabApp, p *page.Page) string {
	return ansiSequence.ReplaceAllString(app.renderPageFooter(p), "")
}

func TestFooterOrdersFixedKeysFirstThenSorted(t *testing.T) {
	app := New()
	p := page.New("Keys").
		BindKey("x", "Xray", NoOpMsg{}).
		BindKey("a", "Alpha", NoOpMsg{}).
		BindKey("Enter", "Open", NoOpMsg{}).
		BindKey("↑/↓", "Move", NoOpMsg{}).
		BindKey("h", "", NoOpMsg{}).
		SetFooter(map[string]string{"Esc": "Back"})
	app.MustAddPage(p)

	got := plainFooter(app, p)
	want := "↑/↓ Move • Enter Open • Esc Back • q Quit • a Alpha • x Xray"
	if got != want {
		t.Errorf("footer = %q\nwant     %q", got, want)
	}
}

func TestFooterLabelPrecedence(t *testing.T) {
	app := New()
	p := page.New("P").BindKey("h", "Page help", NoOpMsg{})
	app.MustAddPage(p)
	app.MustBindKey("h", "Universal help", NoOpMsg{}) // bound after the page, so both exist

	if got := plainFooter(app, p); !strings.Contains(got, "h Page help") || strings.Contains(got, "Universal") {
		t.Errorf("page labels override universal ones: %q", got)
	}

	p.SetFooter(map[string]string{"h": "Manual"})
	if got := plainFooter(app, p); !strings.Contains(got, "h Manual") || strings.Contains(got, "Page help") {
		t.Errorf("SetFooter overrides page labels: %q", got)
	}

	other := page.New("Other")
	app.MustAddPage(other)
	if got := plainFooter(app, other); got != "q Quit • h Universal help" {
		t.Errorf("universal labels alone: %q", got)
	}
}

func TestFooterHidesEmptyLabelsAndFallsBack(t *testing.T) {
	app := New()
	p := page.New("Quiet").BindKey("x", "", NoOpMsg{})
	app.MustAddPage(p)

	if got := plainFooter(app, p); got != "q Quit" {
		t.Errorf("footer = %q, want only the quit label (tab keys and empty labels hidden)", got)
	}

	app.UnbindKey("q")
	if got := plainFooter(app, p); got != fallbackFooterText {
		t.Errorf("footer = %q, want the fallback text", got)
	}
	if got := app.renderFooterContent(); got != fallbackFooterText {
		t.Errorf("renderFooterContent = %q, want the page's footer", got)
	}

	empty := New()
	empty.MustAddSection(NewSection("Empty"))
	if empty.currentPage() != nil {
		t.Fatal("precondition: an empty section has no page")
	}
	if got := empty.renderFooterContent(); got != fallbackFooterText {
		t.Errorf("renderFooterContent without a page = %q", got)
	}
}

func TestFooterGlowsThePressedKey(t *testing.T) {
	app := New()
	app.MustBindKey("x", "Mark", NoOpMsg{})
	p := page.New("P")
	app.MustAddPage(p)
	if app.KeyGlowColor() == app.FooterKeyColor() {
		t.Fatal("precondition: the glow color must differ from the key color")
	}
	glowX := lipgloss.NewStyle().Foreground(app.KeyGlowColor()).Bold(true).Render("x")
	plainX := lipgloss.NewStyle().Foreground(app.FooterKeyColor()).Bold(true).Render("x")
	plainQ := lipgloss.NewStyle().Foreground(app.FooterKeyColor()).Bold(true).Render("q")

	runeKey(app, 'x')
	footer := app.renderPageFooter(p)
	if !strings.Contains(footer, glowX) {
		t.Error("the pressed key should render in the glow color")
	}
	if strings.Contains(footer, plainX) {
		t.Error("the pressed key should not also render in the normal key color")
	}
	if !strings.Contains(footer, plainQ) {
		t.Error("other keys keep the normal key color")
	}

	app.Update(keyGlowTimeoutMsg{id: app.keyGlowTimerID})
	footer = app.renderPageFooter(p)
	if strings.Contains(footer, glowX) || !strings.Contains(footer, plainX) {
		t.Error("after the glow timeout the key renders normally")
	}
}

func TestKeyBindingAccessorsAndRegistry(t *testing.T) {
	app := New()
	app.MustBindKey("x", "Mark", markerMsg{n: 3})

	kb := app.UniversalBindings()["x"]
	if kb.Key() != "x" || kb.Label() != "Mark" {
		t.Errorf("binding = %q %q, want x Mark", kb.Key(), kb.Label())
	}
	if m, ok := kb.Msg().(markerMsg); !ok || m.n != 3 {
		t.Errorf("binding message = %#v, want markerMsg{n: 3}", kb.Msg())
	}

	bindings := app.UniversalBindings()
	delete(bindings, "q")
	bindings["zz"] = KeyBinding{key: "zz"}
	if !app.IsUniversalKey("q") || app.IsUniversalKey("zz") {
		t.Error("UniversalBindings must return a copy")
	}

	if err := app.AddPage(page.New("P").BindKey("x", "Page x", NoOpMsg{})); err == nil {
		t.Error("a page may not bind a universal key")
	}
	app.UnbindKey("x")
	if err := app.AddPage(page.New("P").BindKey("x", "Page x", NoOpMsg{})); err != nil {
		t.Errorf("after UnbindKey a page may bind the key: %v", err)
	}

	if err := app.BindKey("ctrl+q", "Quit", QuitMsg{}); err == nil || !strings.Contains(err.Error(), "hardcoded") {
		t.Errorf("BindKey on a hardcoded key = %v, want a hardcoded error", err)
	}
	if !app.IsHardcoded("ctrl+q") || app.IsHardcoded("q") {
		t.Error("only ctrl+c and ctrl+q are hardcoded")
	}
}

func TestUnbindKeyChains(t *testing.T) {
	app := New()
	if got := app.UnbindKey("q").UnbindKey("tab"); got != app {
		t.Error("UnbindKey should return the app for chaining")
	}
	if app.IsUniversalKey("q") || app.IsUniversalKey("tab") || !app.IsUniversalKey("shift+tab") {
		t.Error("UnbindKey should remove exactly the keys given")
	}
}
