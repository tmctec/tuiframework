package tabapp

import (
	"regexp"
	"slices"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var manySpaces = regexp.MustCompile(` {2,}`)

// separatorColumns returns the rune columns holding a column separator.
func separatorColumns(row string) []int {
	var cols []int
	for i, r := range []rune(row) {
		if r == '│' {
			cols = append(cols, i)
		}
	}
	return cols
}

// plainView returns the frame without ANSI sequences, runs of spaces collapsed.
func plainView(app *TabApp) string {
	return manySpaces.ReplaceAllString(ansiSequence.ReplaceAllString(app.Render(), ""), " ")
}

// electronicsApp is mixedDepthApp on Phones > p2 with the given layout.
func electronicsApp(t *testing.T, mode NavBarMode) *TabApp {
	t.Helper()
	app := mixedDepthApp() // Electronics{Phones{p1,p2}, Computers{c1}, Deals, Clearance}
	app.SetNavBarMode(mode).SetNavBarAutoHide(false)
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if err := app.NavigateToPath([]int{0, 0, 1}); err != nil {
		t.Fatal(err)
	}
	return app
}

func TestNavBarModeCycles(t *testing.T) {
	app := rootMixApp()
	for _, want := range []NavBarMode{NavBarExpanded, NavBarPath, NavBarColumns, NavBarCompact} {
		app.Update(ToggleNavBarModeMsg{})
		if got := app.Model().NavBarMode; got != want {
			t.Fatalf("mode = %v, want %v", got, want)
		}
	}
}

func TestPathShowsRouteAndDropDown(t *testing.T) {
	app := electronicsApp(t, NavBarPath)
	if plain := plainView(app); !strings.Contains(plain, "Electronics › Phones › p2") || strings.Contains(plain, "▾ Phones") {
		t.Fatalf("idle Path bar should show the route only:\n%s", plain)
	}
	if h := lipgloss.Height(app.renderNavSection()); h != 2 {
		t.Errorf("idle band height = %d, want 2 (route and hint)", h)
	}

	key(app, "tab") // wraps from p2 to p1 and opens the drop-down
	plain := plainView(app)
	if !strings.Contains(plain, "Electronics › Phones › p1") {
		t.Errorf("route should follow the selection:\n%s", plain)
	}
	if strings.Count(plain, "p2") != 1 || strings.Count(plain, "p1") != 2 {
		t.Errorf("drop-down should list p1 and p2 once each:\n%s", plain)
	}
	if h := lipgloss.Height(app.renderNavSection()); h != 4 {
		t.Errorf("band height = %d, want 4 (route, two siblings, hint)", h)
	}

	key(app, "left") // focus the Phones segment
	if app.tabNavLevel != 1 {
		t.Fatalf("level = %d, want 1", app.tabNavLevel)
	}
	if plain := plainView(app); !strings.Contains(plain, "▾ Computers") || !strings.Contains(plain, "Clearance") {
		t.Errorf("drop-down should list Phones' siblings:\n%s", plain)
	}
	key(app, "down") // Computers
	if plain := plainView(app); !strings.Contains(plain, "Electronics › Computers › c1") {
		t.Errorf("route should cascade into Computers:\n%s", plain)
	}
	key(app, "right") // Computers is a section: focus its children
	if app.tabNavLevel != 2 {
		t.Fatalf("level = %d, want 2", app.tabNavLevel)
	}
	key(app, "right") // c1 is a page: nothing deeper
	if app.tabNavLevel != 2 {
		t.Errorf("right on a page moved to level %d", app.tabNavLevel)
	}
	key(app, "backspace")
	if app.tabNavLevel != 1 {
		t.Errorf("backspace should go back to level 1, got %d", app.tabNavLevel)
	}
	key(app, "right")
	key(app, "enter") // open page c1
	if app.model.TabNavActive {
		t.Error("enter on a page should end navigation")
	}
	if h := lipgloss.Height(app.renderNavSection()); h != 2 {
		t.Errorf("band height after navigation = %d, want 2", h)
	}
}

func TestColumnsListEveryLevel(t *testing.T) {
	app := electronicsApp(t, NavBarColumns)
	plain := plainView(app)
	for _, want := range []string{"▾ Electronics", "▾ Phones", "▾ Computers", "Deals", "Clearance", "p1", "p2", "│"} {
		if !strings.Contains(plain, want) {
			t.Errorf("%q missing from the Columns bar:\n%s", want, plain)
		}
	}
	if h := lipgloss.Height(app.renderNavSection()); h != 5 {
		t.Errorf("band height = %d, want 5 (four rows for Electronics' children, plus the hint)", h)
	}
	// The separators sit in the same columns on every row, empty slots included.
	rows := strings.Split(ansiSequence.ReplaceAllString(app.renderNavSection(), ""), "\n")
	want := separatorColumns(rows[0])
	for i, row := range rows[:4] {
		if got := separatorColumns(row); !slices.Equal(got, want) {
			t.Errorf("row %d separators at %v, want %v:\n%s", i, got, want, strings.Join(rows, "\n"))
		}
	}

	key(app, "tab") // level 2: p2 wraps to p1
	if app.tabNavLevel != 2 || app.NavPath()[2] != 0 {
		t.Fatalf("after tab: level %d path %v", app.tabNavLevel, app.NavPath())
	}
	key(app, "left") // Phones' column
	key(app, "down") // Computers
	if p := app.NavPath(); p[1] != 1 || p[2] != 0 {
		t.Fatalf("path %v, want Computers > c1", p)
	}
	if plain := plainView(app); strings.Contains(plain, "p1") || !strings.Contains(plain, "c1") {
		t.Errorf("third column should now list Computers' children:\n%s", plain)
	}
	key(app, "enter") // Computers is a section
	if app.tabNavLevel != 2 || !app.model.TabNavActive {
		t.Fatalf("enter on a section should focus its column: level %d active %v", app.tabNavLevel, app.model.TabNavActive)
	}
	key(app, "enter") // c1 is a page
	if app.model.TabNavActive {
		t.Error("enter on a page should end navigation")
	}
}

func TestHintNamesTheStartKeyWhileNavigationIsOff(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarAutoHide(false)
	app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	if !strings.Contains(plainView(app), "Tab: nav mode") {
		t.Error("idle hint should say how to start navigation")
	}
	key(app, "tab")
	if strings.Contains(plainView(app), "Tab: nav mode") {
		t.Error("active hint should drop the start key")
	}
}
