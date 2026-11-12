package tabapp

import (
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/page"
)

// pageName returns the active page's name, or "" when there is none.
func pageName(app *TabApp) string {
	if p := app.currentPage(); p != nil {
		return p.Name()
	}
	return ""
}

// deepMixApp builds root > {S1 > {S11 > {x1}}, Top}: depth 3 with the page
// Top beside the section S1 at the root.
func deepMixApp() *TabApp {
	s11 := NewSection("S11").AddPage(page.New("x1"))
	app := New()
	app.MustAddSection(NewSection("S1").AddSection(s11)).MustAddPage(page.New("Top"))
	return app
}

func TestActiveAndSetActiveSyncTheModel(t *testing.T) {
	flat := New()
	flat.MustAddPage(page.New("p1")).MustAddPage(page.New("p2")).MustAddPage(page.New("p3"))
	if flat.SetActive(2) != flat {
		t.Error("SetActive should return the app for chaining")
	}
	if flat.Active() != 2 || flat.Model().ActivePage != 2 {
		t.Errorf("Active %d model %d, want 2 and 2", flat.Active(), flat.Model().ActivePage)
	}

	mix := rootMixApp()
	mix.SetActive(1)
	if pageName(mix) != "a2" || mix.Model().ActivePage != 1 {
		t.Errorf("page %q model %d, want a2 and 1", pageName(mix), mix.Model().ActivePage)
	}

	bare := New()
	bare.navPath = nil
	bare.SetActive(3)
	if bare.Active() != 0 {
		t.Errorf("Active without a path = %d, want 0", bare.Active())
	}
}

func TestCurrentSectionIndexRange(t *testing.T) {
	app := rootMixApp()
	if err := app.SwitchSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if got := app.CurrentSectionIndex(0); got != 1 {
		t.Errorf("CurrentSectionIndex(0) = %d, want 1", got)
	}
	if got := app.CurrentSectionIndex(1); got != 0 {
		t.Errorf("CurrentSectionIndex(1) = %d, want 0", got)
	}
	for _, level := range []int{-1, 2, 9} {
		if got := app.CurrentSectionIndex(level); got != 0 {
			t.Errorf("CurrentSectionIndex(%d) = %d, want 0", level, got)
		}
	}
}

func TestEnterSectionPushesBreadcrumbAndResetsDeeperLevels(t *testing.T) {
	app := mixedDepthApp()
	if err := app.NavigateToPath([]int{0, 0, 1}); err != nil { // Phones > p2
		t.Fatal(err)
	}

	if err := app.EnterSection(1, 1); err != nil { // Computers
		t.Fatal(err)
	}

	if got := app.NavPath(); !slices.Equal(got, []int{0, 1, 0}) {
		t.Errorf("NavPath = %v, want [0 1 0]", got)
	}
	if pageName(app) != "c1" || app.Model().ActivePage != 0 {
		t.Errorf("page %q model %d, want c1 and 0", pageName(app), app.Model().ActivePage)
	}
	crumbs := app.Breadcrumbs()
	if len(crumbs) != 1 {
		t.Fatalf("breadcrumbs = %d, want 1", len(crumbs))
	}
	if crumbs[0].SectionIndex() != 0 || crumbs[0].PageIndex() != 1 {
		t.Errorf("breadcrumb = section %d page %d, want 0 and 1", crumbs[0].SectionIndex(), crumbs[0].PageIndex())
	}
}

func TestEnterSectionErrorsLeavePositionAlone(t *testing.T) {
	app := deepMixApp()
	cases := []struct {
		name         string
		level, index int
	}{
		{"negative level", -1, 0},
		{"page level", 2, 0},
		{"index too large", 0, 5},
		{"negative index", 0, -1},
	}
	for _, c := range cases {
		before := app.NavPath()
		if err := app.EnterSection(c.level, c.index); err == nil {
			t.Errorf("%s: want an error", c.name)
		}
		if got := app.NavPath(); !slices.Equal(got, before) || len(app.Breadcrumbs()) != 0 {
			t.Errorf("%s: path %v crumbs %d, want %v and 0", c.name, got, len(app.Breadcrumbs()), before)
		}
	}

	if err := app.NavigateToPath([]int{1, 0, 0}); err != nil { // Top: a page at the root
		t.Fatal(err)
	}
	if err := app.EnterSection(1, 0); err == nil {
		t.Error("entering below a page should error: there is no section at that level")
	}
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0, 0}) {
		t.Errorf("path %v, want [1 0 0]", got)
	}
}

func TestPopToBreadcrumbRestoresAndTruncates(t *testing.T) {
	app := rootMixApp()
	if err := app.NavigateToPath([]int{0, 1}); err != nil { // a2
		t.Fatal(err)
	}
	if err := app.EnterSection(0, 1); err != nil { // B > b1
		t.Fatal(err)
	}
	if err := app.EnterSection(0, 0); err != nil { // A > a1
		t.Fatal(err)
	}
	if len(app.Breadcrumbs()) != 2 {
		t.Fatalf("breadcrumbs = %d, want 2", len(app.Breadcrumbs()))
	}

	// Popping the newest breadcrumb drops it from the trail.
	if err := app.PopToBreadcrumb(1); err != nil {
		t.Fatal(err)
	}
	if len(app.Breadcrumbs()) != 1 {
		t.Errorf("pop 1: crumbs %d, want 1", len(app.Breadcrumbs()))
	}

	if err := app.PopToBreadcrumb(0); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("pop 0: path %v, want [0 1]", got)
	}
	if pageName(app) != "a2" || app.Model().ActivePage != 1 {
		t.Errorf("pop 0: page %q model %d, want a2 and 1", pageName(app), app.Model().ActivePage)
	}
	if len(app.Breadcrumbs()) != 0 {
		t.Errorf("breadcrumbs = %d after popping to the first, want 0", len(app.Breadcrumbs()))
	}

	if err := app.PopToBreadcrumb(0); err == nil {
		t.Error("popping with no breadcrumbs should error")
	}
	if err := app.PopToBreadcrumb(-1); err == nil {
		t.Error("a negative breadcrumb index should error")
	}
}

func TestBreadcrumbsReturnsCopy(t *testing.T) {
	app := rootMixApp()
	if err := app.EnterSection(0, 1); err != nil {
		t.Fatal(err)
	}
	got := app.Breadcrumbs()
	got[0] = Breadcrumb{level: 0, path: []int{9, 9}}
	if app.Breadcrumbs()[0].SectionIndex() == 9 {
		t.Error("Breadcrumbs must return a copy")
	}
}

func TestSwitchSectionResetsDeeperLevels(t *testing.T) {
	app := rootMixApp()
	if err := app.NavigateToPath([]int{0, 1}); err != nil { // a2
		t.Fatal(err)
	}

	if err := app.SwitchSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || pageName(app) != "b1" || app.Model().ActivePage != 0 {
		t.Errorf("path %v page %q model %d, want [1 0] b1 0", got, pageName(app), app.Model().ActivePage)
	}

	if err := app.NavigateToPath([]int{0, 1}); err != nil {
		t.Fatal(err)
	}
	if err := app.SwitchSection(0, 0); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("switching to the selected section moved to %v", got)
	}
}

func TestSwitchSectionErrors(t *testing.T) {
	app := rootMixApp()
	for _, c := range []struct {
		name         string
		level, index int
	}{
		{"negative level", -1, 0},
		{"page level", 1, 0},
		{"index too large", 0, 5},
		{"negative index", 0, -1},
	} {
		if err := app.SwitchSection(c.level, c.index); err == nil {
			t.Errorf("%s: want an error", c.name)
		}
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0}) {
		t.Errorf("errors moved the path to %v", got)
	}

	flat := New()
	flat.MustAddPage(page.New("p"))
	if err := flat.SwitchSection(0, 0); err == nil {
		t.Error("a flat app has no section level")
	}

	deep := deepMixApp()
	if err := deep.NavigateToPath([]int{1, 0, 0}); err != nil {
		t.Fatal(err)
	}
	if err := deep.SwitchSection(1, 0); err == nil {
		t.Error("switching below a page should error: there is no section at that level")
	}
}

func TestNavigateToPathErrorKeepsBreadcrumbs(t *testing.T) {
	app := rootMixApp()
	if err := app.NavigateToPath([]int{0, 1}); err != nil {
		t.Fatal(err)
	}
	if err := app.EnterSection(0, 1); err != nil {
		t.Fatal(err)
	}

	if err := app.NavigateToPath([]int{0, -1}); err == nil {
		t.Error("a negative index should error")
	}
	if err := app.NavigateToRootItem(5); err == nil {
		t.Error("an out-of-range root item should error")
	}
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || len(app.Breadcrumbs()) != 1 || app.Model().ActivePage != 0 {
		t.Errorf("path %v crumbs %d model %d, want unchanged [1 0], 1, 0", got, len(app.Breadcrumbs()), app.Model().ActivePage)
	}
}

// Tab navigation on the expanded bar: one row per level.

func TestExpandedTabKeysWrapAndResetDeeperLevels(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarMode(NavBarExpanded).SetNavBarAutoHide(false)

	key(app, "tab")
	if app.TabNavLevel() != 1 || pageName(app) != "a2" || app.Model().ActivePage != 1 {
		t.Fatalf("after tab: level %d page %q model %d", app.TabNavLevel(), pageName(app), app.Model().ActivePage)
	}
	key(app, "shift+tab")
	key(app, "shift+tab")
	if pageName(app) != "a2" {
		t.Errorf("shift+tab twice from a2: %q, want a2 (wrapped)", pageName(app))
	}

	key(app, "up")
	if app.TabNavLevel() != 0 {
		t.Fatalf("level = %d after up, want 0", app.TabNavLevel())
	}
	key(app, "right") // B: level 1 resets to its first item
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || pageName(app) != "b1" || app.Model().ActivePage != 0 {
		t.Errorf("right at root: path %v page %q model %d, want [1 0] b1 0", got, pageName(app), app.Model().ActivePage)
	}
	key(app, "left")
	key(app, "left") // wraps onto Settings
	if got := app.NavPath(); !slices.Equal(got, []int{2, 0}) || pageName(app) != "Settings" {
		t.Errorf("left twice: path %v page %q, want [2 0] Settings", got, pageName(app))
	}
	key(app, "right") // wraps back to A
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0}) {
		t.Errorf("right from the last item: path %v, want [0 0]", got)
	}

	key(app, "down")
	key(app, "down")
	if app.TabNavLevel() != 1 {
		t.Errorf("down clamps at the deepest level, got %d", app.TabNavLevel())
	}
	key(app, "up")
	key(app, "up")
	if app.TabNavLevel() != 0 {
		t.Errorf("up clamps at the root level, got %d", app.TabNavLevel())
	}

	flat := New()
	flat.MustAddPage(page.New("p1")).MustAddPage(page.New("p2"))
	flat.SetNavBarMode(NavBarExpanded)
	key(flat, "tab")
	key(flat, "up")
	key(flat, "down")
	if flat.TabNavLevel() != 0 {
		t.Errorf("a flat app has one level, got %d", flat.TabNavLevel())
	}
}

func TestSidebarUsesListKeys(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarPosition(NavBarLeft).SetNavBarAutoHide(false)
	app.MustBindKey("n", "Nav border", ToggleNavBarBorderMsg{})
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	key(app, "tab")
	if app.TabNavLevel() != 1 || pageName(app) != "a2" {
		t.Fatalf("after tab: level %d page %q", app.TabNavLevel(), pageName(app))
	}
	key(app, "down")
	if pageName(app) != "a1" {
		t.Errorf("down wraps to a1, got %q", pageName(app))
	}
	key(app, "up")
	if pageName(app) != "a2" {
		t.Errorf("up returns to a2, got %q", pageName(app))
	}
	key(app, "right")
	if app.TabNavLevel() != 1 {
		t.Errorf("right on a page stays at level 1, got %d", app.TabNavLevel())
	}
	key(app, "left")
	if app.TabNavLevel() != 0 {
		t.Errorf("left goes up a level, got %d", app.TabNavLevel())
	}
	key(app, "right")
	if app.TabNavLevel() != 1 {
		t.Errorf("right on a section descends, got %d", app.TabNavLevel())
	}
	key(app, "backspace")
	if app.TabNavLevel() != 0 {
		t.Errorf("backspace goes up a level, got %d", app.TabNavLevel())
	}
	key(app, "enter")
	if app.TabNavLevel() != 1 || !app.TabNavActive() {
		t.Errorf("enter on a section descends: level %d active %v", app.TabNavLevel(), app.TabNavActive())
	}

	runeKey(app, 'n')
	if !app.Model().ShowNavBarBorder || !app.TabNavActive() {
		t.Error("a key navigation does not use reaches the bindings and keeps navigation active")
	}

	key(app, "enter")
	if app.TabNavActive() {
		t.Error("enter on a page ends navigation")
	}
}

func TestTabNavGuards(t *testing.T) {
	app := rootMixApp()
	app.tabNavLevel = 5
	app.moveAtLevel(1)
	app.moveAtLevel(-1)
	app.navStripMove(1)
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0}) {
		t.Errorf("an out-of-range level moved the path to %v", got)
	}

	empty := New()
	empty.MustAddSection(NewSection("Empty"))
	empty.SetNavBarMode(NavBarExpanded)
	key(empty, "tab")
	if got := empty.NavPath(); !slices.Equal(got, []int{0, 0}) || !empty.TabNavActive() {
		t.Errorf("tab on an empty section: path %v active %v", got, empty.TabNavActive())
	}
	if empty.currentPage() != nil {
		t.Error("an empty section has no current page")
	}
}

func TestNavigationHelpersOutOfRange(t *testing.T) {
	bare := New()
	if bare.Pages() != nil || bare.sectionAtLevel(0) != nil {
		t.Error("without a root there is no section at any level")
	}
	if bare.itemCountAtLevel(0) != 0 || bare.currentPage() != nil {
		t.Error("without a root there are no items and no page")
	}

	app := rootMixApp()
	if app.sectionAtLevel(-1) != nil {
		t.Error("a negative level has no section")
	}
	if app.itemCountAtLevel(-1) != 0 || app.itemCountAtLevel(2) != 0 {
		t.Error("levels outside the hierarchy hold no items")
	}
	if app.itemCountAtLevel(0) != 3 || app.itemCountAtLevel(1) != 2 {
		t.Errorf("item counts = %d, %d, want 3 and 2", app.itemCountAtLevel(0), app.itemCountAtLevel(1))
	}
	root := app.Root()
	if childName(root, 0) != "A" || childName(root, 2) != "Settings" || childName(root, 3) != "" {
		t.Error("childName should name sections, then pages, then nothing")
	}
	if maxChildren(nil) != 0 || maxChildren(root) != 3 {
		t.Errorf("maxChildren = %d, %d, want 0 and 3", maxChildren(nil), maxChildren(root))
	}

	if err := app.NavigateToPath([]int{2, 0}); err != nil { // Settings
		t.Fatal(err)
	}
	if app.sectionAtLevel(1) != nil || app.itemCountAtLevel(1) != 0 {
		t.Error("below a root page there is no section and no items")
	}
	if app.sectionAtLevel(0) != root {
		t.Error("the section at level 0 is the root")
	}
}

func TestNavigateToPathRepairsPathLength(t *testing.T) {
	app := rootMixApp()
	app.navPath = nil
	if err := app.NavigateToPath([]int{1, 0}); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) || pageName(app) != "b1" {
		t.Errorf("path %v page %q, want [1 0] b1", got, pageName(app))
	}
}

// TestActiveCountsPagesAtTheActiveLevel covers a page that sits beside a
// section: Active is its index among the pages there, not its place in the
// navigation path, and SetActive moves among those pages.
func TestActiveCountsPagesAtTheActiveLevel(t *testing.T) {
	app := New()
	app.MustAddSection(NewSection("S").AddPage(page.New("s1")))
	app.MustAddPage(page.New("P0")).MustAddPage(page.New("P1"))

	if err := app.NavigateToPath([]int{2, 0}); err != nil { // S, P0, then P1
		t.Fatal(err)
	}
	if pageName(app) != "P1" || app.Active() != 1 || app.Model().ActivePage != 1 {
		t.Fatalf("page %q active %d model %d, want P1, 1 and 1", pageName(app), app.Active(), app.Model().ActivePage)
	}

	app.SetActive(0)
	if pageName(app) != "P0" || app.Active() != 0 {
		t.Errorf("after SetActive(0): page %q active %d, want P0 and 0", pageName(app), app.Active())
	}
	for _, bad := range []int{-1, 2, 99} {
		app.SetActive(bad)
		if pageName(app) != "P0" || app.Active() != 0 {
			t.Errorf("SetActive(%d) moved to %q, want an out-of-range index ignored", bad, pageName(app))
		}
	}
}

// TestExpandedDownStopsOnAPage checks the Expanded bar's down key: it opens
// the level under a highlighted section and stays put on a page, which has
// no level beneath it.
func TestExpandedDownStopsOnAPage(t *testing.T) {
	app := electronicsApp(t, NavBarExpanded) // Electronics > Phones > p2
	key(app, "tab")                          // focus on the page level
	key(app, "up")                           // up to Phones, Computers, Deals, Clearance
	if app.TabNavLevel() != 1 {
		t.Fatalf("focus level = %d, want 1", app.TabNavLevel())
	}

	key(app, "right") // Computers, a section
	key(app, "down")
	if app.TabNavLevel() != 2 {
		t.Errorf("down on a section: focus level = %d, want 2", app.TabNavLevel())
	}

	key(app, "up")
	key(app, "right") // Deals, a page beside the sections
	if pageName(app) != "Deals" {
		t.Fatalf("page = %q, want Deals", pageName(app))
	}
	key(app, "down")
	if app.TabNavLevel() != 1 {
		t.Errorf("down on a page: focus level = %d, want it to stay on 1", app.TabNavLevel())
	}
}

// TestFocusStaysWithinThePath moves the position to a shallower page while
// the focus is on a deeper level: the focus follows the page up.
func TestFocusStaysWithinThePath(t *testing.T) {
	app := mixedDepthApp() // Electronics{Phones{p1,p2}, Computers{c1}, Deals, Clearance}
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if err := app.NavigateToPath([]int{0, 2, 0}); err != nil { // Deals, level 1
		t.Fatal(err)
	}
	if err := app.EnterSection(1, 0); err != nil { // into Phones, leaving a crumb at Deals
		t.Fatal(err)
	}
	app.tabNavLevel = 2 // focus on the pages of Phones
	if pageName(app) != "p1" {
		t.Fatalf("page = %q, want p1", pageName(app))
	}

	if err := app.PopToBreadcrumb(0); err != nil {
		t.Fatal(err)
	}
	if pageName(app) != "Deals" || app.TabNavLevel() != 1 {
		t.Errorf("after the pop: page %q focus %d, want Deals and 1", pageName(app), app.TabNavLevel())
	}
}

// TestNavbarRendersWithAFocusOutsideThePath forces the focus below the
// levels the path reaches and renders every layout. Navigation never leaves
// the focus there; the bar must still draw.
func TestNavbarRendersWithAFocusOutsideThePath(t *testing.T) {
	for _, pos := range []NavBarPosition{NavBarTop, NavBarLeft} {
		for _, mode := range []NavBarMode{NavBarCompact, NavBarExpanded, NavBarPath, NavBarColumns} {
			app := deepMixApp() // root > {S1 > {S11 > {x1}}, Top}
			app.SetNavBarPosition(pos).SetNavBarMode(mode).SetNavBarAutoHide(false)
			app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
			if err := app.NavigateToPath([]int{1, 0, 0}); err != nil { // Top, a root page
				t.Fatal(err)
			}
			app.model.TabNavActive = true
			app.tabNavLevel = 2

			plain := ansiSequence.ReplaceAllString(app.Render(), "")
			if lines := strings.Split(plain, "\n"); len(lines) != 24 {
				t.Errorf("pos=%v mode=%v: %d lines, want 24", pos, mode, len(lines))
			}
			if !strings.Contains(plain, "Top") {
				t.Errorf("pos=%v mode=%v: the bar should still name the selected root page:\n%s", pos, mode, plain)
			}
		}
	}
}
