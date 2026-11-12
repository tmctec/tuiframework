package tabapp

import (
	"slices"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/page"
)

// useTrueColor makes lipgloss emit color sequences for the test, so styled
// and unstyled renders of the same text differ.
func useTrueColor(t *testing.T) {
	t.Helper()
}

// plainLines strips ANSI sequences and collapses runs of spaces on every line.
func plainLines(s string) []string {
	return strings.Split(manySpaces.ReplaceAllString(ansiSequence.ReplaceAllString(s, ""), " "), "\n")
}

// leadingSpaces returns the number of spaces a line starts with.
func leadingSpaces(s string) int {
	return len(s) - len(strings.TrimLeft(s, " "))
}

func TestNavbarCompactFlatStrip(t *testing.T) {
	useTrueColor(t)
	app := New()
	app.MustAddPage(page.New("One")).MustAddPage(page.New("Two")).MustAddPage(page.New("Three"))
	app.SetActive(1)
	sizedApp(app)

	tabs := app.renderNavBar()
	for _, want := range []string{app.tabStyle.Render("One"), app.activeTabStyle.Render("Two"), app.tabStyle.Render("Three")} {
		if !strings.Contains(tabs, want) {
			t.Errorf("%q missing from the strip:\n%s", want, tabs)
		}
	}
	if got := strings.Contains(tabs, app.activeTabStyle.Render("One")); got {
		t.Error("an inactive page should not take the active style")
	}
	lines := plainLines(tabs)
	if len(lines) != 2 || strings.TrimSpace(lines[0]) != "One Two Three" {
		t.Errorf("strip lines = %q, want the page names then the hint", lines)
	}
}

func TestNavbarCompactCrumbsAtDepth(t *testing.T) {
	app := electronicsApp(t, NavBarCompact) // Electronics > Phones > p2
	lines := plainLines(app.renderNavBar())
	if got := strings.TrimSpace(lines[0]); got != "▴ Electronics ▴ Phones p1 p2" {
		t.Errorf("strip = %q, want the two crumbs then the pages", got)
	}
	if !strings.Contains(lines[1], "⌫ back") {
		t.Errorf("hint = %q, want the back key with a level above", lines[1])
	}
}

// TestNavbarCompactCrumbFocusStyles moves the strip highlight onto a crumb:
// the crumb takes the active style and the selected item dims.
func TestNavbarCompactCrumbFocusStyles(t *testing.T) {
	useTrueColor(t)
	app := stripApp()
	key(app, "tab") // level 1, a2
	tabs := app.compactTabs()
	if len(tabs) != 3 {
		t.Fatalf("%d tabs, want crumb A, a1, a2", len(tabs))
	}
	if tabs[0] != app.selectedInactiveStyle.Render("▴ A") || tabs[2] != app.activeTabStyle.Render("a2") {
		t.Errorf("idle crumb and selected item: %q %q", tabs[0], tabs[2])
	}

	key(app, "left") // a1
	key(app, "left") // crumb A
	if app.navFocusCrumb != 0 {
		t.Fatalf("focus = %d, want crumb 0", app.navFocusCrumb)
	}
	tabs = app.compactTabs()
	if tabs[0] != app.activeTabStyle.Render("▴ A") {
		t.Errorf("focused crumb = %q, want the active style", tabs[0])
	}
	if tabs[1] != app.selectedInactiveStyle.Render("a1") {
		t.Errorf("selected item under a focused crumb = %q, want the muted style", tabs[1])
	}
}

func TestNavbarExpandedRows(t *testing.T) {
	app := electronicsApp(t, NavBarExpanded) // Electronics > Phones > p2, focus level 2
	rows := plainLines(app.renderExpandedTabs())
	if len(rows) != 4 {
		t.Fatalf("idle rows = %q, want three levels and the hint", rows)
	}
	for i, r := range rows[:3] {
		if strings.HasPrefix(r, "▸") {
			t.Errorf("idle row %d carries the level marker: %q", i, r)
		}
	}

	key(app, "tab") // navigation at level 2
	rows = plainLines(app.renderExpandedTabs())
	if !strings.HasPrefix(rows[2], "▸ ") || strings.HasPrefix(rows[0], "▸") || strings.HasPrefix(rows[1], "▸") {
		t.Errorf("marker should sit on the focused level 2:\n%s", strings.Join(rows, "\n"))
	}

	key(app, "up") // level 1: level 2 stays visible, unmarked
	rows = plainLines(app.renderExpandedTabs())
	if len(rows) != 4 || !strings.HasPrefix(rows[1], "▸ ") || !strings.Contains(rows[2], "p1") {
		t.Errorf("at level 1 the row below stays:\n%s", strings.Join(rows, "\n"))
	}

	key(app, "up") // level 0: level 2 is hidden
	all := strings.Join(rows, "\n")
	rows = plainLines(app.renderExpandedTabs())
	if len(rows) != 3 || !strings.HasPrefix(rows[0], "▸ ") {
		t.Errorf("at level 0 two levels and the hint remain:\n%s", strings.Join(rows, "\n"))
	}
	if joined := strings.Join(rows, "\n"); strings.Contains(joined, "p1") || strings.Contains(joined, "p2") {
		t.Errorf("pages two levels below the focus should be hidden:\n%s\nbefore:\n%s", joined, all)
	}
}

// TestNavbarExpandedSkipsLevelsUnderAPage selects a root-level page: the
// level beneath it has no parent section, so only the root row renders.
func TestNavbarExpandedSkipsLevelsUnderAPage(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarMode(NavBarExpanded).SetNavBarAutoHide(false)
	if err := app.NavigateToPath([]int{2, 0}); err != nil { // Settings
		t.Fatal(err)
	}
	rows := plainLines(app.renderExpandedTabs())
	if len(rows) != 2 || !strings.Contains(rows[0], "Settings") {
		t.Errorf("rows = %q, want the root row and the hint", rows)
	}

	if rows := plainLines(New().renderExpandedTabs()); len(rows) != 1 || !strings.Contains(rows[0], "▴▾ level") {
		t.Errorf("an empty app renders only the hint, got %q", rows)
	}
}

func TestNavbarPathRouteAndDropDown(t *testing.T) {
	app := electronicsApp(t, NavBarPath) // Electronics > Phones > p2
	lines := plainLines(app.renderPathTabs())
	if len(lines) != 2 || strings.TrimSpace(lines[0]) != "Electronics › Phones › p2" {
		t.Fatalf("idle Path lines = %q, want the route and the hint", lines)
	}

	sepW := lipgloss.Width(app.footerStyle.Render(" › "))
	segW := func(name string) int { return lipgloss.Width(app.selectedInactiveStyle.UnsetMargins().Render(name)) }

	key(app, "tab") // focus level 2, p2 wraps to p1
	lines = plainLines(app.renderPathTabs())
	if len(lines) != 4 || !strings.Contains(lines[1], "p1") || !strings.Contains(lines[2], "p2") {
		t.Fatalf("navigating Path lines = %q, want route, p1, p2, hint", lines)
	}
	// Page labels carry two spaces so they line up with the ▾ of sections.
	wantIndent := segW("Electronics") + sepW + segW("Phones") + sepW + app.activeTabStyle.GetPaddingLeft() + 2
	raw := strings.Split(ansiSequence.ReplaceAllString(app.renderPathTabs(), ""), "\n")
	if got := leadingSpaces(raw[1]); got != wantIndent {
		t.Errorf("drop-down indent = %d, want %d under the page segment:\n%s", got, wantIndent, raw[1])
	}

	key(app, "left") // focus the Phones segment
	raw = strings.Split(ansiSequence.ReplaceAllString(app.renderPathTabs(), ""), "\n")
	wantIndent = segW("Electronics") + sepW + app.tabStyle.GetPaddingLeft()
	if got := leadingSpaces(raw[1]); got != wantIndent {
		t.Errorf("drop-down indent = %d, want %d under the Phones segment:\n%s", got, wantIndent, raw[1])
	}
	if lines = plainLines(strings.Join(raw, "\n")); len(lines) != 6 {
		t.Errorf("Phones has four siblings: lines = %q", lines)
	}
}

// TestNavbarPathClampsFocus renders the Path bar with a focus level outside
// the route: it is clamped into the route instead of panicking.
func TestNavbarPathClampsFocus(t *testing.T) {
	app := electronicsApp(t, NavBarPath)
	app.model.TabNavActive = true

	app.tabNavLevel = 7
	lines := plainLines(app.renderPathTabs())
	if len(lines) != 4 || !strings.Contains(lines[1], "p1") {
		t.Errorf("focus past the page level should list the page's siblings: %q", lines)
	}

	app.tabNavLevel = -1
	lines = plainLines(app.renderPathTabs())
	if len(lines) != 3 || !strings.Contains(lines[1], "▾ Electronics") {
		t.Errorf("negative focus should list the root's children: %q", lines)
	}
}

func TestNavbarColumnsPadBlankSlots(t *testing.T) {
	app := stripApp()
	app.SetNavBarMode(NavBarColumns)
	raw := strings.Split(ansiSequence.ReplaceAllString(app.renderColumnTabs(), ""), "\n")
	if len(raw) != 4 {
		t.Fatalf("%d lines, want three rows (A, B, Settings) and the hint", len(raw))
	}
	want := separatorColumns(raw[0])
	if len(want) != 1 {
		t.Fatalf("row 0 separators at %v, want one between the two columns", want)
	}
	for i, row := range raw[:3] {
		if got := separatorColumns(row); !slices.Equal(got, want) {
			t.Errorf("row %d separators at %v, want %v", i, got, want)
		}
		if lipgloss.Width(row) != lipgloss.Width(raw[0]) {
			t.Errorf("row %d width %d, want %d: blank slots must be padded", i, lipgloss.Width(row), lipgloss.Width(raw[0]))
		}
	}
	if !strings.Contains(raw[2], "Settings") || strings.TrimSpace(strings.SplitN(raw[2], "│", 2)[1]) != "" {
		t.Errorf("third row should list Settings and a blank second slot: %q", raw[2])
	}

	// A page at the root has nothing beneath it: one column only.
	if err := app.NavigateToPath([]int{2, 0}); err != nil {
		t.Fatal(err)
	}
	raw = strings.Split(ansiSequence.ReplaceAllString(app.renderColumnTabs(), ""), "\n")
	if len(raw) != 4 || strings.Contains(strings.Join(raw, "\n"), "│") {
		t.Errorf("Settings selected: %q, want one column and no separator", raw)
	}

	// A selected section without children ends the cascade too.
	empty := New()
	empty.MustAddSection(NewSection("Empty")).MustAddSection(NewSection("Full").AddPage(page.New("f1")))
	empty.SetNavBarMode(NavBarColumns)
	raw = strings.Split(ansiSequence.ReplaceAllString(empty.renderColumnTabs(), ""), "\n")
	if strings.Contains(strings.Join(raw, "\n"), "│") || !strings.Contains(raw[0], "▾ Empty") {
		t.Errorf("empty section selected: %q, want one column", raw)
	}
}

func TestNavbarSideBarCompactLevels(t *testing.T) {
	app := rootMixApp()
	app.SetNavBarPosition(NavBarLeft).SetNavBarAutoHide(false)
	sizedApp(app)

	// New starts the focus at the root level, so the sidebar lists the
	// sections with their arrows and the root page indented.
	plain := plainLines(app.renderSideNavBar(18))
	joined := strings.Join(plain, "\n")
	for _, want := range []string{"▸ A", "▸ B", "Settings", "nav mode"} {
		if !strings.Contains(joined, want) {
			t.Errorf("%q missing from the root-level sidebar:\n%s", want, joined)
		}
	}

	key(app, "tab") // level 1: a2
	joined = strings.Join(plainLines(app.renderSideNavBar(18)), "\n")
	for _, want := range []string{"◂ A", "a1", "a2"} {
		if !strings.Contains(joined, want) {
			t.Errorf("%q missing from the level-1 sidebar:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "▸ B") {
		t.Errorf("compact sidebar at level 1 should not list the root sections:\n%s", joined)
	}

	// Back at the root level A is highlighted. A cell spans the bar: all 20
	// columns without the border, 18 inside it.
	key(app, "left")
	useTrueColor(t)
	side := app.renderSideNavBar(18)
	if !strings.Contains(side, app.activeTabStyle.UnsetMargins().Width(20).MaxWidth(20).Render("▸ A")) {
		t.Errorf("A should take the active style at the focused root level:\n%s", side)
	}
	app.SetShowNavBarBorder(true)
	side = app.renderSideNavBar(18)
	if !strings.Contains(side, app.activeTabStyle.UnsetMargins().Width(18).MaxWidth(18).Render("▸ A")) {
		t.Errorf("inside the border A's cell should be two columns narrower:\n%s", side)
	}
}

// TestNavbarSideBarTreeExpandsAlongTheFocus renders the tree sidebar: the
// selected branch opens down to the focused level and the focused item takes
// the active style.
func TestNavbarSideBarTreeExpandsAlongTheFocus(t *testing.T) {
	for _, mode := range []NavBarMode{NavBarExpanded, NavBarColumns} {
		app := mixedDepthApp() // Electronics > Phones > p1
		app.SetNavBarPosition(NavBarLeft).SetNavBarMode(mode).SetNavBarAutoHide(false).SetNavBarWidth(30)
		sizedApp(app)

		raw := strings.Split(ansiSequence.ReplaceAllString(app.renderSideNavBar(18), ""), "\n")
		joined := strings.Join(raw, "\n")
		// indentOf returns the indent of the line holding label, border excluded.
		indentOf := func(label string) int {
			for _, l := range raw {
				if strings.Contains(l, label) {
					return leadingSpaces(strings.TrimLeft(l, "│"))
				}
			}
			t.Errorf("%v: %q missing from the tree:\n%s", mode, label, joined)
			return -1
		}
		root := indentOf("▸ Electronics")
		for label, depth := range map[string]int{"▸ Phones": 1, "▸ Computers": 1, "Deals": 1, "Clearance": 1, "p1": 2, "p2": 2} {
			if got := indentOf(label); got != root+2*depth {
				t.Errorf("%v: %q indented %d, want %d (depth %d)", mode, label, got, root+2*depth, depth)
			}
		}
		if strings.Contains(joined, "c1") {
			t.Errorf("%v: Computers is not selected, so c1 stays folded:\n%s", mode, joined)
		}

		key(app, "tab")  // navigation at level 2
		key(app, "left") // level 1: Phones' children fold
		joined = strings.Join(plainLines(app.renderSideNavBar(18)), "\n")
		if strings.Contains(joined, "p1") || !strings.Contains(joined, "▸ Phones") {
			t.Errorf("%v: focus on level 1 should fold level 2:\n%s", mode, joined)
		}
	}
}

// TestNavbarSideBarTruncatesLongNames renders long and multibyte names in a
// narrow sidebar: nothing panics and the names end with an ellipsis.
func TestNavbarSideBarTruncatesLongNames(t *testing.T) {
	long := NewSection("Café Très Longue Section")
	long.AddSection(NewSection("Sous-section Étendue").AddPage(page.New("Página Larga Çà")))
	long.AddPage(page.New("Résumé Général"))
	for _, mode := range []NavBarMode{NavBarCompact, NavBarExpanded, NavBarColumns} {
		for _, width := range []int{8, 4} {
			app := New()
			app.MustAddSection(long).MustAddSection(NewSection("B").AddPage(page.New("b1")))
			app.SetNavBarPosition(NavBarLeft).SetNavBarMode(mode).SetNavBarAutoHide(false).SetNavBarWidth(width)
			sizedApp(app)
			if err := app.NavigateToPath([]int{0, 0, 0}); err != nil {
				t.Fatal(err)
			}

			view := app.Render()
			side := plainLines(app.renderSideNavBar(18))
			var truncated int
			for _, l := range side {
				if strings.HasSuffix(strings.TrimRight(strings.Trim(l, "│"), " "), "…") {
					truncated++
				}
			}
			if truncated == 0 {
				t.Errorf("mode %v width %d: no sidebar line ends with …:\n%s", mode, width, strings.Join(side, "\n"))
			}
			if strings.Contains(strings.Join(side, "\n"), "Café Très Longue Section") {
				t.Errorf("mode %v width %d: the long name was not truncated", mode, width)
			}
			if len(strings.Split(view, "\n")) != 24 {
				t.Errorf("mode %v width %d: frame has %d lines", mode, width, len(strings.Split(view, "\n")))
			}
		}
	}
}

func TestNavbarHintNamesTheForwardKey(t *testing.T) {
	app := rootMixApp()
	if got, want := app.navHint("keys"), app.footerStyle.Render("Tab: nav mode • keys"); got != want {
		t.Errorf("navHint = %q, want %q", got, want)
	}

	app.SetTabNavKeyForward("f1")
	if got, want := app.navHint("keys"), app.footerStyle.Render("F1: nav mode • keys"); got != want {
		t.Errorf("navHint with f1 = %q, want %q", got, want)
	}

	app.model.TabNavActive = true
	if got, want := app.navHint("keys"), app.footerStyle.Render("keys"); got != want {
		t.Errorf("navHint while navigating = %q, want %q", got, want)
	}

	app.model.TabNavActive = false
	app.SetTabNavKeyForward("")
	if got, want := app.navHint("keys"), app.footerStyle.Render("keys"); got != want {
		t.Errorf("navHint without a forward key = %q, want %q", got, want)
	}
}

func TestNavbarListCellWidthUsesWidestPadding(t *testing.T) {
	app := New()
	base := app.listCellWidth([]string{"ab", "abcd"})
	if want := 4 + app.tabStyle.GetHorizontalPadding(); base != want {
		t.Errorf("listCellWidth = %d, want %d", base, want)
	}
	app.SetActiveTabStyle(lipgloss.NewStyle().Padding(0, 6))
	if got := app.listCellWidth([]string{"ab", "abcd"}); got != 4+12 {
		t.Errorf("listCellWidth with a wider active padding = %d, want 16", got)
	}
}
