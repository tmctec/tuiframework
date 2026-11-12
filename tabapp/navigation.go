package tabapp

import (
	"fmt"
	"slices"

	"github.com/tmctec/tuiframework/page"
)

// Active returns the index of the active page among the pages at its level,
// the value views read as Model.ActivePage.
func (a *TabApp) Active() int {
	return a.model.ActivePage
}

// SetActive activates the page at index among the pages at the active page's
// level. An index out of range leaves the position alone. Returns self for
// chaining.
func (a *TabApp) SetActive(index int) *TabApp {
	a.selectPageAtCurrentLevel(index)
	return a
}

// NavPath returns a copy of the current navigation path: one index per
// hierarchy level, counting sections first and then pages at each level.
func (a *TabApp) NavPath() []int {
	return slices.Clone(a.navPath)
}

// hierarchyLevels returns the number of navigation levels, the depth of the
// tree. An app without content still has one level, for pages.
func (a *TabApp) hierarchyLevels() int {
	if a.root == nil {
		return 1
	}
	return a.root.Depth()
}

// sectionAtLevel returns the section whose children are the items at level:
// the root for level 0, below it the section the path selects at each level
// above. It is nil without content and where the path selects a page before
// reaching the level.
func (a *TabApp) sectionAtLevel(level int) *Section {
	if a.root == nil || level < 0 {
		return nil
	}
	section := a.root
	for i := 0; i < level && i < len(a.navPath); i++ {
		idx := a.navPath[i]
		if idx >= len(section.Sections()) {
			return nil
		}
		section = section.Sections()[idx]
	}
	return section
}

// selectedAt returns the index the path selects at level, counting sections
// first and then pages. A level outside the path reports 0.
func (a *TabApp) selectedAt(level int) int {
	if level < 0 || level >= len(a.navPath) {
		return 0
	}
	return a.navPath[level]
}

// currentPage returns the active page: the first page the path selects on
// its way down the tree, or nil when it selects none.
func (a *TabApp) currentPage() *page.Page {
	parent, level, offset := a.pageLevel()
	if parent == nil {
		return nil
	}
	if i := a.navPath[level] - offset; i >= 0 && i < len(parent.Pages()) {
		return parent.Pages()[i]
	}
	return nil
}

// currentPageLevel returns the level the active page is on: the first level
// where the path selects a page instead of a section, so a page beside
// subsections counts as much as a leaf. A path that selects no page reports
// the deepest level.
func (a *TabApp) currentPageLevel() int {
	for level := range a.navPath {
		parent := a.sectionAtLevel(level)
		if parent == nil {
			break
		}
		if i := a.navPath[level] - len(parent.Sections()); i >= 0 && i < len(parent.Pages()) {
			return level
		}
	}
	return a.hierarchyLevels() - 1
}

// itemCountAtLevel returns the number of items the path offers at level:
// the sections and pages of the section there, 0 where there is none.
func (a *TabApp) itemCountAtLevel(level int) int {
	parent := a.sectionAtLevel(level)
	if parent == nil || level >= a.hierarchyLevels() {
		return 0
	}
	return parent.ItemCount()
}

// Breadcrumbs returns a copy of the breadcrumb trail, oldest first.
func (a *TabApp) Breadcrumbs() []Breadcrumb {
	return slices.Clone(a.breadcrumbs)
}

// CurrentSectionIndex returns the index the navigation path selects at
// level, counting sections first and then pages, and 0 for a level outside
// the path.
func (a *TabApp) CurrentSectionIndex(level int) int {
	return a.selectedAt(level)
}

// EnterSection selects the section at index among the child sections at
// level and pushes the position it leaves onto the breadcrumb trail, for
// PopToBreadcrumb to return to.
func (a *TabApp) EnterSection(level, index int) error {
	if err := a.checkSection(level, index); err != nil {
		return err
	}
	a.breadcrumbs = append(a.breadcrumbs, Breadcrumb{level: level, path: slices.Clone(a.navPath)})
	a.selectAtLevel(level, index)
	return nil
}

// checkSection reports whether index names a child section at level, the
// only items EnterSection and SwitchSection select.
func (a *TabApp) checkSection(level, index int) error {
	if maxLevel := a.hierarchyLevels() - 2; level < 0 || level > maxLevel {
		return fmt.Errorf("level %d out of range for sections (the deepest section level is %d)", level, maxLevel)
	}
	parent := a.sectionAtLevel(level)
	if parent == nil {
		return fmt.Errorf("no section at level %d", level)
	}
	if count := len(parent.Sections()); index < 0 || index >= count {
		return fmt.Errorf("section index %d out of range at level %d (have %d sections)", index, level, count)
	}
	return nil
}

// PopToBreadcrumb returns to the position recorded by the breadcrumb at
// breadcrumbIndex and discards it and every deeper breadcrumb.
func (a *TabApp) PopToBreadcrumb(breadcrumbIndex int) error {
	if breadcrumbIndex < 0 || breadcrumbIndex >= len(a.breadcrumbs) {
		return fmt.Errorf("breadcrumb index %d out of range (have %d breadcrumbs)", breadcrumbIndex, len(a.breadcrumbs))
	}

	// Return to the recorded position, unless the tree has changed depth
	// since, and drop this crumb and every deeper one.
	bc := a.breadcrumbs[breadcrumbIndex]
	if len(bc.path) == len(a.navPath) {
		copy(a.navPath, bc.path)
	}
	a.breadcrumbs = a.breadcrumbs[:breadcrumbIndex]
	a.syncActivePage()
	return nil
}

// SwitchSection selects the section at index among the child sections at
// level. Every deeper level returns to its first item. Unlike EnterSection
// it leaves the breadcrumb trail alone.
func (a *TabApp) SwitchSection(level, index int) error {
	if err := a.checkSection(level, index); err != nil {
		return err
	}
	a.selectAtLevel(level, index)
	return nil
}

// NavigateToPath jumps to an absolute navigation path: one index per
// hierarchy level, counting sections first and then pages at each level.
// Once an index selects a page, deeper indices are forced to 0. The jump
// clears the breadcrumb trail and points the nav bar at the level of the
// resulting page.
func (a *TabApp) NavigateToPath(path []int) error {
	if a.root == nil {
		return fmt.Errorf("no content to navigate")
	}
	levels := a.hierarchyLevels()
	if len(path) != levels {
		return fmt.Errorf("path length %d does not match hierarchy depth %d", len(path), levels)
	}

	resolved := make([]int, levels)
	parent := a.root
	for level, idx := range path {
		if parent == nil {
			break // a page was selected at a shallower level
		}
		count := len(parent.Sections()) + len(parent.Pages())
		if idx < 0 || idx >= count {
			return fmt.Errorf("index %d out of range at level %d (have %d items)", idx, level, count)
		}
		resolved[level] = idx
		if idx < len(parent.Sections()) {
			parent = parent.Sections()[idx]
		} else {
			parent = nil
		}
	}

	if len(a.navPath) != levels {
		a.navPath = make([]int, levels)
	}
	copy(a.navPath, resolved)
	a.breadcrumbs = a.breadcrumbs[:0]
	a.tabNavLevel = a.currentPageLevel()
	a.navFocusCrumb = -1
	a.syncActivePage()
	return nil
}

// NavigateToRootItem jumps to the root-level item at index (sections first,
// then pages) and to the first item at every deeper level.
func (a *TabApp) NavigateToRootItem(index int) error {
	path := make([]int, a.hierarchyLevels())
	path[0] = index
	return a.NavigateToPath(path)
}

// selectAtLevel highlights item idx at level. When that changes the
// selection, every deeper level returns to its first item.
func (a *TabApp) selectAtLevel(level, idx int) {
	if a.navPath[level] == idx {
		return
	}
	a.navPath[level] = idx
	for i := level + 1; i < len(a.navPath); i++ {
		a.navPath[i] = 0
	}
	a.syncActivePage()
}

// syncActivePage mirrors the active page's index at its level into the Model
// and keeps the navigation focus on or above the page's level. Every navPath
// change ends here so views always see the current page.
func (a *TabApp) syncActivePage() {
	a.model.ActivePage, _ = a.currentPageIndex()
	a.tabNavLevel = min(a.tabNavLevel, a.currentPageLevel())
}

// pageLevel returns the section whose pages include the active page, the
// level those pages are on, and the path index of the first of them (the
// section's count of child sections). parent is nil when the path reaches
// no section at that level: without content, or under an empty section.
func (a *TabApp) pageLevel() (parent *Section, level, offset int) {
	level = a.currentPageLevel()
	parent = a.sectionAtLevel(level)
	if parent == nil || level >= len(a.navPath) {
		return nil, level, 0
	}
	return parent, level, len(parent.Sections())
}

// selectPageAtCurrentLevel activates the page at pageIndex (0-based among the
// pages at the current page level) and resets deeper levels. An index out of
// range leaves the position alone.
func (a *TabApp) selectPageAtCurrentLevel(pageIndex int) {
	parent, level, offset := a.pageLevel()
	if parent == nil || pageIndex < 0 || pageIndex >= len(parent.Pages()) {
		return
	}
	a.selectAtLevel(level, offset+pageIndex)
}

// currentPageIndex returns the active page's index among the pages at the
// current page level, and that level's page count (0 when there is none).
func (a *TabApp) currentPageIndex() (index, count int) {
	parent, level, offset := a.pageLevel()
	if parent == nil {
		return 0, 0
	}
	return max(a.navPath[level]-offset, 0), len(parent.Pages())
}

// childName returns the plain name of a section's child at navPath index idx.
func childName(parent *Section, idx int) string {
	if idx < len(parent.Sections()) {
		return parent.Sections()[idx].Name()
	}
	if i := idx - len(parent.Sections()); i >= 0 && i < len(parent.Pages()) {
		return parent.Pages()[i].Name()
	}
	return ""
}

// maxChildren returns the largest number of children any section holds.
func maxChildren(s *Section) int {
	if s == nil {
		return 0
	}
	n := len(s.Sections()) + len(s.Pages())
	for _, sub := range s.Sections() {
		if m := maxChildren(sub); m > n {
			n = m
		}
	}
	return n
}
