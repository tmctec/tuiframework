package tabapp

import (
	"slices"
	"testing"

	"github.com/tmctec/tuiframework/page"
)

// rootMixApp builds: sections A{a1,a2}, B{b1} and a root-level page Settings.
// Depth 2; Settings is root item 2 (after the two sections).
func rootMixApp() *TabApp {
	app := New()
	a := NewSection("A")
	a.AddPage(page.New("a1"))
	a.AddPage(page.New("a2"))
	b := NewSection("B")
	b.AddPage(page.New("b1"))
	app.MustAddSection(a).MustAddSection(b).MustAddPage(page.New("Settings"))
	return app
}

func TestNavigateToPathPageBesideSections(t *testing.T) {
	app := rootMixApp()

	if err := app.NavigateToPath([]int{2, 0}); err != nil {
		t.Fatal(err)
	}

	if got := app.NavPath(); !slices.Equal(got, []int{2, 0}) {
		t.Errorf("NavPath = %v, want [2 0]", got)
	}
	if app.currentPage().Name() != "Settings" {
		t.Errorf("current page = %q, want Settings", app.currentPage().Name())
	}
	if app.TabNavLevel() != 0 {
		t.Errorf("TabNavLevel = %d, want 0", app.TabNavLevel())
	}
}

func TestNavigateToRootItemSectionLandsOnFirstPage(t *testing.T) {
	app := rootMixApp()

	if err := app.NavigateToRootItem(1); err != nil {
		t.Fatal(err)
	}

	if app.currentPage().Name() != "b1" {
		t.Errorf("current page = %q, want b1", app.currentPage().Name())
	}
	if app.TabNavLevel() != 1 {
		t.Errorf("TabNavLevel = %d, want 1", app.TabNavLevel())
	}
}

func TestNavigateToPathResetsBreadcrumbs(t *testing.T) {
	app := rootMixApp()
	if err := app.EnterSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if len(app.Breadcrumbs()) != 1 {
		t.Fatalf("precondition: want 1 breadcrumb, got %d", len(app.Breadcrumbs()))
	}

	if err := app.NavigateToPath([]int{0, 1}); err != nil {
		t.Fatal(err)
	}

	if len(app.Breadcrumbs()) != 0 {
		t.Errorf("breadcrumbs = %d, want 0", len(app.Breadcrumbs()))
	}
	if app.currentPage().Name() != "a2" {
		t.Errorf("current page = %q, want a2", app.currentPage().Name())
	}
}

func TestNavigateToPathRejectsBadInput(t *testing.T) {
	app := rootMixApp()
	before := app.NavPath()

	if err := app.NavigateToPath([]int{0}); err == nil {
		t.Error("wrong path length should error")
	}
	if err := app.NavigateToPath([]int{7, 0}); err == nil {
		t.Error("out-of-range index should error")
	}
	if err := app.NavigateToPath([]int{0, 5}); err == nil {
		t.Error("out-of-range page index should error")
	}
	if got := app.NavPath(); !slices.Equal(got, before) {
		t.Errorf("NavPath changed to %v on error", got)
	}
}

func TestNavigateToPathIgnoresIndicesBelowPage(t *testing.T) {
	app := rootMixApp()
	// Settings is a page at level 0, so the level-1 index is forced to 0.
	if err := app.NavigateToPath([]int{2, 1}); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{2, 0}) {
		t.Errorf("NavPath = %v, want [2 0]", got)
	}
}

func TestNavigateToPathFlatPages(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("p1")).MustAddPage(page.New("p2")).MustAddPage(page.New("p3"))

	if err := app.NavigateToRootItem(2); err != nil {
		t.Fatal(err)
	}
	if app.Active() != 2 {
		t.Errorf("Active = %d, want 2", app.Active())
	}
}

func TestNavigateToPathNoContent(t *testing.T) {
	if err := New().NavigateToPath([]int{0}); err == nil {
		t.Error("navigating with no content should error")
	}
}

func TestNavigateToPathMsgViaUpdate(t *testing.T) {
	app := rootMixApp() // autohide on by default

	_, cmd := app.Update(NavigateToPathMsg{Path: []int{2, 0}})

	if app.currentPage().Name() != "Settings" {
		t.Errorf("current page = %q, want Settings", app.currentPage().Name())
	}
	if cmd == nil {
		t.Error("jump should re-arm the autohide timer")
	}
	if app.navBarHidden {
		t.Error("jump should show the nav bar")
	}
}

func TestNavPathReturnsCopy(t *testing.T) {
	app := rootMixApp()
	got := app.NavPath()
	got[0] = 99
	if app.NavPath()[0] == 99 {
		t.Error("NavPath must return a copy")
	}
}

// mixedDepthApp builds root > Electronics > {Phones{p1,p2}, Computers{c1}} plus
// the pages Deals and Clearance beside the two subsections. Depth 3.
func mixedDepthApp() *TabApp {
	phones := NewSection("Phones")
	phones.AddPage(page.New("p1"))
	phones.AddPage(page.New("p2"))
	computers := NewSection("Computers")
	computers.AddPage(page.New("c1"))
	electronics := NewSection("Electronics")
	electronics.AddSection(phones)
	electronics.AddSection(computers)
	electronics.AddPage(page.New("Deals"))
	electronics.AddPage(page.New("Clearance"))
	app := New()
	app.MustAddSection(electronics)
	return app
}

func TestDigitKeySelectsPageAtCurrentLevel(t *testing.T) {
	app := mixedDepthApp()
	if err := app.NavigateToPath([]int{0, 2, 0}); err != nil { // Deals
		t.Fatal(err)
	}

	app.Update(KeyPress("2"))

	if app.currentPage().Name() != "Clearance" {
		t.Errorf("current page = %q, want Clearance", app.currentPage().Name())
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 3, 0}) {
		t.Errorf("NavPath = %v, want [0 3 0]", got)
	}
}

func TestDigitKeyOutOfRangeAtCurrentLevel(t *testing.T) {
	app := mixedDepthApp()
	if err := app.NavigateToPath([]int{0, 2, 0}); err != nil {
		t.Fatal(err)
	}
	before := app.NavPath()

	app.Update(KeyPress("3"))

	if got := app.NavPath(); !slices.Equal(got, before) {
		t.Errorf("NavPath = %v, want unchanged %v", got, before)
	}
}

func TestNextPrevPageCycleWithinCurrentLevel(t *testing.T) {
	app := mixedDepthApp()
	if err := app.NavigateToPath([]int{0, 2, 0}); err != nil { // Deals
		t.Fatal(err)
	}

	app.updateModel(NavigateNextPageMsg{})
	if app.currentPage().Name() != "Clearance" {
		t.Fatalf("after next: %q, want Clearance", app.currentPage().Name())
	}
	app.updateModel(NavigateNextPageMsg{})
	if app.currentPage().Name() != "Deals" {
		t.Fatalf("after wrap: %q, want Deals", app.currentPage().Name())
	}
	app.updateModel(NavigatePrevPageMsg{})
	if app.currentPage().Name() != "Clearance" {
		t.Fatalf("after prev: %q, want Clearance", app.currentPage().Name())
	}
	if app.Model().ActivePage != 1 {
		t.Errorf("ActivePage = %d, want 1", app.Model().ActivePage)
	}
}

func TestNextPageAtDeepestLevel(t *testing.T) {
	app := mixedDepthApp() // starts at Electronics > Phones > p1
	app.updateModel(NavigateNextPageMsg{})
	if app.currentPage().Name() != "p2" {
		t.Errorf("current page = %q, want p2", app.currentPage().Name())
	}
}
