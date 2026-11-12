package tabapp

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
)

// renderNavBar renders the nav bar in the layout of the nav bar mode.
// A tree with a single level has nothing to expand, so its Expanded bar is
// the Compact one.
func (a *TabApp) renderNavBar() string {
	switch a.model.NavBarMode {
	case NavBarPath:
		return a.renderPathTabs()
	case NavBarColumns:
		return a.renderColumnTabs()
	case NavBarExpanded:
		if a.hierarchyLevels() > 1 {
			return a.renderExpandedTabs()
		}
	}
	return a.renderCompactTabs()
}

// eachChild calls fn for every child of parent in navigation-path order:
// its sections first, then its pages. idx is the child's path index.
func eachChild(parent *Section, fn func(idx int, name string, isSection bool)) {
	for i, s := range parent.Sections() {
		fn(i, s.Name(), true)
	}
	offset := len(parent.Sections())
	for i, p := range parent.Pages() {
		fn(offset+i, p.Name(), false)
	}
}

// renderCompactTabs renders the Compact layout: the strip of tabs, then the
// key hint on its own line. Enter opens the highlight, and backspace goes
// back a level when there is one.
func (a *TabApp) renderCompactTabs() string {
	hint := "◂▸ move • ⏎ open"
	if a.hierarchyLevels() > 1 {
		hint += " • ⌫ back"
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, a.compactTabs()...) + "\n" + a.navHint(hint)
}

// compactTabs returns the tabs of the Compact strip: a crumb for the item
// selected at each level above the focus, then the items at the focused
// level, sections marked ▾.
func (a *TabApp) compactTabs() []string {
	var tabs []string

	// The crumb in focus takes the active style; the others stay dim.
	for level := 0; level < a.tabNavLevel && level < a.hierarchyLevels()-1; level++ {
		parent := a.sectionAtLevel(level)
		if parent == nil {
			continue
		}
		style := a.selectedInactiveStyle
		if a.model.TabNavActive && a.navFocusCrumb == level {
			style = a.activeTabStyle
		}
		if idx := a.selectedAt(level); idx < parent.ItemCount() {
			tabs = append(tabs, style.Render("▴ "+childName(parent, idx)))
		}
	}

	parent := a.sectionAtLevel(a.tabNavLevel)
	if parent == nil {
		return tabs
	}

	// The selected item dims while a crumb holds the focus.
	selectedStyle := a.activeTabStyle
	if a.model.TabNavActive && a.navFocusCrumb >= 0 {
		selectedStyle = a.selectedInactiveStyle
	}
	selected := a.selectedAt(a.tabNavLevel)
	eachChild(parent, func(idx int, name string, isSection bool) {
		style := a.tabStyle
		if idx == selected {
			style = selectedStyle
		}
		if isSection {
			name = "▾ " + name
		}
		tabs = append(tabs, style.Render(name))
	})
	return tabs
}

// renderExpandedTabs renders the Expanded layout: a row of tabs for every
// level on the active path, the focused row marked ▸, then the key hint.
func (a *TabApp) renderExpandedTabs() string {
	markerStyle := lipgloss.NewStyle().
		Foreground(a.footerKeyColor).
		Bold(true)

	var rows []string
	for level := 0; level < a.hierarchyLevels(); level++ {
		// While navigating, the rows more than one level below the focus hide.
		if a.model.TabNavActive && level > a.tabNavLevel+1 {
			continue
		}
		parent := a.sectionAtLevel(level)
		if parent == nil {
			continue
		}
		focused := a.model.TabNavActive && a.tabNavLevel == level
		// The row below the focus shows no selection: it follows the highlight.
		belowFocus := a.model.TabNavActive && a.tabNavLevel < level
		selected := a.selectedAt(level)

		var tabs []string
		eachChild(parent, func(idx int, name string, isSection bool) {
			style := a.tabStyle
			switch {
			case idx != selected:
			case focused:
				style = a.activeTabStyle
			case !belowFocus:
				style = a.selectedInactiveStyle
			}
			if isSection {
				name = "▾ " + name
			}
			tabs = append(tabs, style.Render(name))
		})
		if len(tabs) == 0 {
			continue
		}

		marker := "  "
		if focused {
			marker = markerStyle.Render("▸ ")
		}
		rows = append(rows, marker+lipgloss.JoinHorizontal(lipgloss.Top, tabs...))
	}

	rows = append(rows, a.navHint("◂▸ move • ▴▾ level"))
	return strings.Join(rows, "\n")
}

// renderSideNavBar renders the nav bar for the Left position, height lines
// tall, inside the nav bar border when it is on: the key hint, then a tree
// of the open levels in the Expanded and Columns modes, or the crumbs and
// the focused level in the others. Returns "" while autohide has the bar
// hidden.
func (a *TabApp) renderSideNavBar(height int) string {
	if a.model.NavBarAutoHide && a.navBarHidden {
		return ""
	}

	frame := 0 // columns and lines the border takes
	if a.model.ShowNavBarBorder {
		frame = 2
	}
	contentWidth := a.navBarWidth - frame
	// Names get the cell width minus the tab padding; longer ones are cut
	// with an ellipsis rather than wrapped onto a second line.
	textWidth := contentWidth - a.tabPadding()
	label := func(prefix, name string) string {
		return prefix + ansi.Truncate(name, max(textWidth-lipgloss.Width(prefix), 1), "…")
	}

	// Every cell spans the bar; the tab styles' margins stay outside it.
	fit := func(st lipgloss.Style) lipgloss.Style {
		return st.UnsetMargins().Width(contentWidth).MaxWidth(contentWidth)
	}
	normal, active, dim := fit(a.tabStyle), fit(a.activeTabStyle), fit(a.selectedInactiveStyle)

	var rows []string
	levels := a.hierarchyLevels()
	if a.model.NavBarMode == NavBarExpanded || a.model.NavBarMode == NavBarColumns {
		// A tree: each level indented under the section it belongs to.
		var walk func(section *Section, level int)
		walk = func(section *Section, level int) {
			if section == nil || level >= levels {
				return
			}
			indent := strings.Repeat("  ", level)
			selected := a.selectedAt(level)
			selectedStyle := dim
			if a.tabNavLevel == level {
				selectedStyle = active
			}
			eachChild(section, func(idx int, name string, isSection bool) {
				style := normal
				if idx == selected {
					style = selectedStyle
				}
				if !isSection {
					rows = append(rows, style.Render(label(indent, name)))
					return
				}
				rows = append(rows, style.Render(label(indent+"▸ ", name)))
				// The selected section opens while the focus is below it.
				if idx == selected && a.tabNavLevel > level {
					walk(section.Sections()[idx], level+1)
				}
			})
		}
		walk(a.root, 0)
	} else {
		// A crumb for each level above the focus, then the focused level.
		for level := 0; level < a.tabNavLevel && level < levels-1; level++ {
			parent := a.sectionAtLevel(level)
			if parent == nil {
				continue
			}
			if idx := a.selectedAt(level); idx < parent.ItemCount() {
				rows = append(rows, dim.Render(label("◂ ", childName(parent, idx))))
			}
		}
		if parent := a.sectionAtLevel(a.tabNavLevel); parent != nil {
			selected := a.selectedAt(a.tabNavLevel)
			eachChild(parent, func(idx int, name string, isSection bool) {
				style := normal
				if idx == selected {
					style = active
				}
				prefix := "  "
				if isSection {
					prefix = "▸ "
				}
				rows = append(rows, style.Render(label(prefix, name)))
			})
		}
	}

	// The key hint wraps to the bar's width above the items.
	hint := lipgloss.NewStyle().Width(contentWidth).Render(a.navHint(listNavKeys))
	column := hint + "\n" + strings.Join(rows, "\n")

	if !a.model.ShowNavBarBorder {
		// The canvas cuts the column to the bar and fills what is short.
		return column
	}
	style := a.borderColors(borderedBox(contentWidth, height-frame).Border(a.navBarBorderStyle), a.navBarBorderColor)
	return style.Render(column)
}

// navHint renders the navigation key hint. While navigation is off it starts
// with the key that turns it on.
func (a *TabApp) navHint(keys string) string {
	if !a.model.TabNavActive && a.tabNavKeyFwd != "" {
		k := a.tabNavKeyFwd
		keys = strings.ToUpper(k[:1]) + k[1:] + ": nav mode • " + keys
	}
	return a.footerStyle.Render(keys)
}

const listNavKeys = "▴▾ item • ◂▸ level • ⏎ open • ⌫ back"

// childLabels returns a section's children in navPath order: sections first,
// marked ▾, then pages indented to align with them.
func childLabels(parent *Section) []string {
	labels := make([]string, 0, parent.ItemCount())
	eachChild(parent, func(_ int, name string, isSection bool) {
		prefix := "  "
		if isSection {
			prefix = "▾ "
		}
		labels = append(labels, prefix+name)
	})
	return labels
}

// listCellWidth returns the width that fits every label in each tab style.
// Path and Columns render their cells without the styles' margins.
func (a *TabApp) listCellWidth(labels []string) int {
	w := 0
	for _, l := range labels {
		if lw := lipgloss.Width(l); lw > w {
			w = lw
		}
	}
	return w + a.tabPadding()
}

// tabPadding returns the widest horizontal padding among the tab styles.
func (a *TabApp) tabPadding() int {
	pad := a.tabStyle.GetHorizontalPadding()
	for _, st := range []lipgloss.Style{a.activeTabStyle, a.selectedInactiveStyle} {
		if p := st.GetHorizontalPadding(); p > pad {
			pad = p
		}
	}
	return pad
}

// renderPathTabs renders the Path layout: the active route on one line, a
// drop-down of the focused level's siblings while navigating, then the hint.
func (a *TabApp) renderPathTabs() string {
	pageLevel := a.currentPageLevel()
	focus := max(min(a.tabNavLevel, pageLevel), 0)

	sep := a.footerStyle.Render(" › ")
	var route strings.Builder
	x, menuX := 0, 0
	for level := 0; level <= pageLevel && level < len(a.navPath); level++ {
		parent := a.sectionAtLevel(level)
		if parent == nil {
			break
		}
		style := a.selectedInactiveStyle.UnsetMargins()
		if a.model.TabNavActive && level == focus || !a.model.TabNavActive && level == pageLevel {
			style = a.activeTabStyle.UnsetMargins()
		}
		if level > 0 {
			route.WriteString(sep)
			x += lipgloss.Width(sep)
		}
		if level == focus {
			menuX = x
		}
		seg := style.Render(childName(parent, a.navPath[level]))
		route.WriteString(seg)
		x += lipgloss.Width(seg)
	}
	lines := []string{route.String()}

	if parent := a.sectionAtLevel(focus); a.model.TabNavActive && parent != nil && focus < len(a.navPath) {
		labels := childLabels(parent)
		w := a.listCellWidth(labels)
		indent := strings.Repeat(" ", menuX)
		for i, label := range labels {
			style := a.tabStyle.UnsetMargins()
			if i == a.navPath[focus] {
				style = a.activeTabStyle.UnsetMargins()
			}
			lines = append(lines, indent+style.Width(w).Render(label))
		}
	}

	lines = append(lines, a.navHint(listNavKeys))
	return strings.Join(lines, "\n")
}

// renderColumnTabs renders the Columns layout: one column per level along
// the active path, every item listed, the selection cascading to the right.
// The band is as tall as the widest level in the tree, so it never jumps.
func (a *TabApp) renderColumnTabs() string {
	rows := maxChildren(a.root)
	sep := a.footerStyle.Render(" │ ")

	var cols [][]string
	for level := 0; level < a.hierarchyLevels() && level < len(a.navPath); level++ {
		parent := a.sectionAtLevel(level)
		if parent == nil {
			break
		}
		labels := childLabels(parent)
		if len(labels) == 0 {
			break
		}
		w := a.listCellWidth(labels)
		focused := a.model.TabNavActive && level == a.tabNavLevel
		cells := make([]string, 0, rows)
		for i, label := range labels {
			style := a.tabStyle.UnsetMargins()
			if i == a.navPath[level] {
				style = a.selectedInactiveStyle.UnsetMargins()
				if focused {
					style = a.activeTabStyle.UnsetMargins()
				}
			}
			cells = append(cells, style.Width(w).Render(label))
		}
		// Empty slots are as wide as a rendered cell so the separators line up.
		blank := strings.Repeat(" ", lipgloss.Width(cells[0]))
		for len(cells) < rows {
			cells = append(cells, blank)
		}
		cols = append(cols, cells)
	}

	lines := make([]string, 0, rows+1)
	for r := 0; r < rows; r++ {
		parts := make([]string, 0, len(cols))
		for _, col := range cols {
			parts = append(parts, col[r])
		}
		lines = append(lines, strings.Join(parts, sep))
	}
	lines = append(lines, a.navHint(listNavKeys))
	return strings.Join(lines, "\n")
}
