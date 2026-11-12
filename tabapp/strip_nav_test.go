package tabapp

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/page"
)

func key(app *TabApp, k string) {
	app.Update(KeyPress(k))
}

// stripApp is rootMixApp (A{a1,a2}, B{b1}, Settings) on the compact top bar
// with autohide off, so the strip is always visible.
func stripApp() *TabApp {
	app := rootMixApp()
	app.SetNavBarAutoHide(false)
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	return app
}

func TestStripLeftFromFirstItemFocusesCrumb(t *testing.T) {
	app := stripApp()
	key(app, "tab") // enters navigation at the page level and moves to a2
	if app.tabNavLevel != 1 || app.NavPath()[1] != 1 {
		t.Fatalf("after tab: level %d path %v, want level 1 path [0 1]", app.tabNavLevel, app.NavPath())
	}
	key(app, "left") // a1
	key(app, "left") // crumb A
	if app.navFocusCrumb != 0 {
		t.Fatalf("focus = %d, want crumb 0", app.navFocusCrumb)
	}
	if app.NavPath()[1] != 0 {
		t.Errorf("moving onto the crumb changed the selection: %v", app.NavPath())
	}

	key(app, "enter") // opens the crumb: back to the root level
	if app.tabNavLevel != 0 || app.navFocusCrumb != -1 {
		t.Fatalf("after enter on crumb: level %d focus %d, want 0 and -1", app.tabNavLevel, app.navFocusCrumb)
	}
	key(app, "right") // B
	if p := app.NavPath(); p[0] != 1 || p[1] != 0 {
		t.Fatalf("after right at root: path %v, want [1 0]", p)
	}
	key(app, "enter") // opens section B
	if app.tabNavLevel != 1 || !app.model.TabNavActive {
		t.Fatalf("after enter on a section: level %d active %v, want 1 and true", app.tabNavLevel, app.model.TabNavActive)
	}
	key(app, "enter") // b1 is a page: navigation ends
	if app.model.TabNavActive || app.Model().TabNavActive {
		t.Error("enter on a page should end tab navigation")
	}
}

func TestStripWrapsOntoCrumb(t *testing.T) {
	app := stripApp()
	key(app, "tab")   // a2, the last item
	key(app, "right") // wraps onto crumb A
	if app.navFocusCrumb != 0 {
		t.Fatalf("focus = %d, want crumb 0", app.navFocusCrumb)
	}
	key(app, "right") // a1
	if app.navFocusCrumb != -1 || app.NavPath()[1] != 0 {
		t.Errorf("focus %d path %v, want -1 and [0 0]", app.navFocusCrumb, app.NavPath())
	}
}

func TestStripBackspaceGoesUpALevel(t *testing.T) {
	app := stripApp()
	key(app, "tab")
	key(app, "backspace")
	if app.tabNavLevel != 0 || !app.model.TabNavActive {
		t.Fatalf("after backspace: level %d active %v, want 0 and true", app.tabNavLevel, app.model.TabNavActive)
	}
	key(app, "backspace")
	if app.tabNavLevel != 0 {
		t.Errorf("backspace at the root moved to level %d", app.tabNavLevel)
	}
	key(app, "down") // alias for open: A is a section
	if app.tabNavLevel != 1 {
		t.Errorf("down should open the section, level = %d", app.tabNavLevel)
	}
	key(app, "up") // alias for back
	if app.tabNavLevel != 0 {
		t.Errorf("up should go back a level, level = %d", app.tabNavLevel)
	}
}

func TestStripCrumbRoundTripKeepsDeepSelection(t *testing.T) {
	app := mixedDepthApp() // Electronics{Phones{p1,p2}, Computers{c1}, Deals, Clearance}
	app.SetNavBarAutoHide(false)
	if err := app.NavigateToPath([]int{0, 0, 1}); err != nil { // Phones > p2
		t.Fatal(err)
	}
	key(app, "tab") // the strip is Electronics, Phones, p1, p2: from p2 it wraps onto Electronics
	if app.navFocusCrumb != 0 {
		t.Fatalf("focus = %d, want crumb 0 (Electronics)", app.navFocusCrumb)
	}
	key(app, "right") // crumb Phones
	key(app, "right") // p1: the selection changes
	if p := app.NavPath(); p[2] != 0 {
		t.Fatalf("path %v, want p1 selected", p)
	}
	key(app, "left") // back onto crumb Phones: selection stays p1
	key(app, "enter")
	if app.tabNavLevel != 1 {
		t.Fatalf("enter on the Phones crumb should show level 1, got %d", app.tabNavLevel)
	}
	key(app, "enter") // Phones is still selected: reopen it
	if p := app.NavPath(); app.tabNavLevel != 2 || p[2] != 0 {
		t.Errorf("reopening the same section changed the path: level %d path %v", app.tabNavLevel, p)
	}
}

func TestExpandedBarKeepsLevelArrows(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarMode(NavBarExpanded).SetNavBarAutoHide(false)
	key(app, "tab")
	if app.tabNavLevel != 1 {
		t.Fatalf("level = %d, want 1", app.tabNavLevel)
	}
	key(app, "up")
	if app.tabNavLevel != 0 {
		t.Errorf("up in expanded mode should move to level 0, got %d", app.tabNavLevel)
	}
	key(app, "down")
	if app.tabNavLevel != 1 {
		t.Errorf("down in expanded mode should move to level 1, got %d", app.tabNavLevel)
	}
	key(app, "enter")
	if !app.model.TabNavActive {
		t.Error("enter is not a navigation key in expanded mode")
	}
}

func TestStripHintNamesItsKeys(t *testing.T) {
	app := stripApp()
	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	for _, want := range []string{"◂▸ move", "⏎ open", "⌫ back"} {
		if !strings.Contains(plain, want) {
			t.Errorf("%q missing from the compact bar", want)
		}
	}

	flat := New()
	flat.MustAddPage(page.New("Only")).MustAddPage(page.New("Two"))
	flat.SetNavBarAutoHide(false)
	flat.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	plain = ansiSequence.ReplaceAllString(flat.Render(), "")
	if strings.Contains(plain, "⌫ back") {
		t.Error("a flat page list has no level to go back to")
	}
}
