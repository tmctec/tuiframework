package tabapp

import (
	"slices"
	"testing"
)

// TestPopToBreadcrumbReturnsToTheRecordedPosition enters sections at two
// different levels of a depth-3 tree and checks each pop lands exactly where
// the app was before the matching EnterSection.
func TestPopToBreadcrumbReturnsToTheRecordedPosition(t *testing.T) {
	app := mixedDepthApp()
	if err := app.NavigateToPath([]int{0, 0, 1}); err != nil { // Electronics > Phones > p2
		t.Fatal(err)
	}

	if err := app.EnterSection(1, 1); err != nil { // Electronics > Computers > c1
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1, 0}) {
		t.Fatalf("after EnterSection NavPath = %v, want [0 1 0]", got)
	}
	bc := app.Breadcrumbs()[0]
	if bc.Level() != 1 || bc.SectionIndex() != 0 || bc.PageIndex() != 1 || !slices.Equal(bc.Path(), []int{0, 0, 1}) {
		t.Errorf("breadcrumb = level %d section %d page %d path %v, want level 1 section 0 page 1 path [0 0 1]",
			bc.Level(), bc.SectionIndex(), bc.PageIndex(), bc.Path())
	}

	if err := app.PopToBreadcrumb(0); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0, 1}) {
		t.Errorf("after pop NavPath = %v, want [0 0 1]", got)
	}
	if p := app.currentPage(); p == nil || p.Name() != "p2" {
		t.Errorf("current page after pop = %v, want p2", p)
	}
	if app.Model().ActivePage != 1 {
		t.Errorf("Model.ActivePage = %d, want 1", app.Model().ActivePage)
	}
	if len(app.Breadcrumbs()) != 0 {
		t.Errorf("breadcrumbs = %d, want 0", len(app.Breadcrumbs()))
	}
}

// TestPopToBreadcrumbRestoresTheNewestCrumb pushes two crumbs at the root
// level and pops the newest one, which returns to the section entered just
// before it.
func TestPopToBreadcrumbRestoresTheNewestCrumb(t *testing.T) {
	app := rootMixApp()
	app.SetActive(1) // A > a2

	if err := app.EnterSection(0, 1); err != nil { // B > b1
		t.Fatal(err)
	}
	if err := app.EnterSection(0, 0); err != nil { // A > a1
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 0}) {
		t.Fatalf("NavPath = %v, want [0 0]", got)
	}

	if err := app.PopToBreadcrumb(1); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{1, 0}) {
		t.Errorf("after popping the newest crumb NavPath = %v, want [1 0]", got)
	}
	if len(app.Breadcrumbs()) != 1 {
		t.Fatalf("breadcrumbs = %d, want 1", len(app.Breadcrumbs()))
	}

	if err := app.PopToBreadcrumb(0); err != nil {
		t.Fatal(err)
	}
	if got := app.NavPath(); !slices.Equal(got, []int{0, 1}) {
		t.Errorf("after popping the first crumb NavPath = %v, want [0 1]", got)
	}
}
