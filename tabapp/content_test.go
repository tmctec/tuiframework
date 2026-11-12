package tabapp

import (
	"slices"
	"strings"
	"testing"

	"github.com/tmctec/tuiframework/page"
)

// pageNames returns the names of pages in order.
func pageNames(pages []*page.Page) []string {
	names := make([]string, 0, len(pages))
	for _, p := range pages {
		names = append(names, p.Name())
	}
	return names
}

func TestAddSectionRejectsConflictingPageKeys(t *testing.T) {
	app := New()

	err := app.AddSection(NewSection("S").AddPage(page.New("p").BindKey("q", "Quit", NoOpMsg{})))
	if err == nil || !strings.Contains(err.Error(), "conflicts with universal") {
		t.Errorf("universal conflict error = %v", err)
	}

	nested := NewSection("Outer").AddSection(
		NewSection("Inner").AddPage(page.New("p").BindKey("ctrl+q", "Quit", NoOpMsg{})))
	err = app.AddSection(nested)
	if err == nil || !strings.Contains(err.Error(), "hardcoded") {
		t.Errorf("hardcoded conflict error = %v", err)
	}

	if app.HasContent() {
		t.Error("a rejected section must not be added")
	}
}

func TestMustAddSectionPanicsOnConflict(t *testing.T) {
	app := New()
	defer func() {
		if recover() == nil {
			t.Error("MustAddSection should panic on a key conflict")
		}
	}()
	app.MustAddSection(NewSection("S").AddPage(page.New("p").BindKey("q", "Quit", NoOpMsg{})))
}

func TestSetRootSizesNavPathToDepth(t *testing.T) {
	app := New()
	root := NewSection("").AddSection(
		NewSection("L1").AddSection(
			NewSection("L2").AddPage(page.New("deep"))))

	if app.MustSetRoot(root) != app {
		t.Error("MustSetRoot should return the app for chaining")
	}
	if app.Root() != root {
		t.Error("Root should return the section given to SetRoot")
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0, 0}) {
		t.Errorf("NavPath = %v, want [0 0 0]", got)
	}
	if app.currentPage().Name() != "deep" {
		t.Errorf("current page = %q, want deep", app.currentPage().Name())
	}
	if !app.HasSections() || !app.HasContent() {
		t.Error("a root with sections has sections and content")
	}

	if err := app.SetRoot(nil); err != nil {
		t.Fatalf("SetRoot(nil) = %v, want nil", err)
	}
	if app.Root() != nil || app.HasContent() {
		t.Error("SetRoot(nil) clears the content")
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0}) {
		t.Errorf("NavPath = %v after SetRoot(nil), want [0]", got)
	}
	if !strings.HasPrefix(app.Render(), "No content.") {
		t.Errorf("View without content = %q", app.Render())
	}

	app.MustSetRoot(NewSection(""))
	if app.HasContent() || app.HasSections() {
		t.Error("an empty root has no content and no sections")
	}
}

func TestPagesAtTheCurrentLevel(t *testing.T) {
	if New().Pages() != nil {
		t.Error("Pages without a root should be nil")
	}

	flat := New()
	flat.MustAddPage(page.New("p1")).MustAddPage(page.New("p2"))
	if got := pageNames(flat.Pages()); !slices.Equal(got, []string{"p1", "p2"}) {
		t.Errorf("flat Pages = %v", got)
	}

	mix := rootMixApp()
	if got := pageNames(mix.Pages()); !slices.Equal(got, []string{"a1", "a2"}) {
		t.Errorf("Pages in A = %v", got)
	}
	if err := mix.SwitchSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if got := pageNames(mix.Pages()); !slices.Equal(got, []string{"b1"}) {
		t.Errorf("Pages in B = %v", got)
	}
	if err := mix.NavigateToPath([]int{2, 0}); err != nil {
		t.Fatal(err)
	}
	if got := pageNames(mix.Pages()); !slices.Equal(got, []string{"Settings"}) {
		t.Errorf("Pages with a root page active = %v, want the root's pages", got)
	}
}

func TestAddPageAndAddSectionKeepNavPathInShape(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("p"))
	if got := app.NavPath(); !slices.Equal(got, []int{0}) {
		t.Errorf("NavPath = %v after AddPage, want [0]", got)
	}
	if app.HasSections() || !app.HasContent() {
		t.Error("flat pages: content but no sections")
	}

	app.MustAddSection(NewSection("S").AddPage(page.New("s1")))
	if got := app.NavPath(); len(got) != 2 {
		t.Errorf("NavPath = %v after AddSection, want two levels", got)
	}
	if !app.HasSections() {
		t.Error("HasSections should be true after AddSection")
	}
	// A second section of the same depth keeps the path.
	app.MustAddSection(NewSection("T").AddPage(page.New("t1")))
	if got := app.NavPath(); len(got) != 2 {
		t.Errorf("NavPath = %v after a same-depth section, want two levels", got)
	}
}

func TestHasContentWithoutRoot(t *testing.T) {
	app := New()
	if app.HasContent() || app.HasSections() {
		t.Error("a new app has neither content nor sections")
	}
	if err := app.AddPage(page.New("p")); err != nil {
		t.Fatal(err)
	}
	if !app.HasContent() || app.HasSections() {
		t.Error("after AddPage: content without sections")
	}
}

// TestSetRootValidatesPageKeys checks SetRoot holds a tree to the same key
// rules as AddSection: a conflict is an error and the old tree stays.
func TestSetRootValidatesPageKeys(t *testing.T) {
	app := New().MustAddPage(page.New("keep"))
	if err := app.NavigateToPath([]int{0}); err != nil {
		t.Fatal(err)
	}
	old := app.Root()

	bad := NewSection("").AddSection(NewSection("S").AddPage(page.New("p").BindKey("q", "Quit", NoOpMsg{})))
	err := app.SetRoot(bad)
	if err == nil || !strings.Contains(err.Error(), `"q"`) {
		t.Fatalf("SetRoot with a page bound to q: err = %v, want a conflict naming q", err)
	}
	if app.Root() != old || pageName(app) != "keep" {
		t.Error("a rejected root should leave the tree alone")
	}

	defer func() {
		if recover() == nil {
			t.Error("MustSetRoot should panic on a key conflict")
		}
	}()
	app.MustSetRoot(bad)
}

// TestSetRootStartsOver checks a new tree starts at its first page with the
// position, the breadcrumb trail and the Model in step.
func TestSetRootStartsOver(t *testing.T) {
	app := rootMixApp() // A{a1,a2}, B{b1}, Settings
	if err := app.EnterSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if err := app.NavigateToPath([]int{0, 1}); err != nil { // a2
		t.Fatal(err)
	}
	if err := app.EnterSection(0, 1); err != nil {
		t.Fatal(err)
	}
	if len(app.Breadcrumbs()) == 0 {
		t.Fatal("precondition: a breadcrumb should be on the trail")
	}

	app.MustSetRoot(NewSection("").AddPages(page.New("x"), page.New("y")))

	if got := app.NavPath(); !slices.Equal(got, []int{0}) {
		t.Errorf("NavPath = %v, want [0]", got)
	}
	if pageName(app) != "x" || app.Active() != 0 || app.Model().ActivePage != 0 {
		t.Errorf("page %q active %d model %d, want x, 0 and 0", pageName(app), app.Active(), app.Model().ActivePage)
	}
	if len(app.Breadcrumbs()) != 0 || app.TabNavLevel() != 0 {
		t.Errorf("breadcrumbs %d, nav level %d, want none and 0", len(app.Breadcrumbs()), app.TabNavLevel())
	}
}
