package tabapp

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

// startTabNavTimer starts the timer for tab navigation mode timeout.
// When the timer expires, tab navigation mode is deactivated.
func (a *TabApp) startTabNavTimer() tea.Cmd {
	a.tabNavTimerID++
	id := a.tabNavTimerID
	return tea.Tick(a.model.TabNavTimeout, func(_ time.Time) tea.Msg {
		return tabNavTimeoutMsg{id: id}
	})
}

// moveAtLevel moves the highlight at the focused level by delta items,
// wrapping at both ends, and selects the item it lands on.
func (a *TabApp) moveAtLevel(delta int) {
	count := a.itemCountAtLevel(a.tabNavLevel)
	if count == 0 || a.tabNavLevel >= len(a.navPath) {
		return
	}
	a.selectAtLevel(a.tabNavLevel, ((a.navPath[a.tabNavLevel]+delta)%count+count)%count)
}

// focusLevelUp moves the focus to the level above, if there is one.
func (a *TabApp) focusLevelUp() {
	if a.tabNavLevel > 0 {
		a.tabNavLevel--
	}
}

// enterTabNav starts tab navigation at the active page's level. It is a
// no-op while navigation is already active.
func (a *TabApp) enterTabNav() {
	if a.model.TabNavActive {
		return
	}
	a.model.TabNavActive = true
	a.tabNavLevel = a.currentPageLevel()
	a.navFocusCrumb = -1
}

// leaveTabNav ends tab navigation and, with autohide on, hides the bar at once.
func (a *TabApp) leaveTabNav() {
	a.model.TabNavActive = false
	a.navFocusCrumb = -1
	if a.model.NavBarAutoHide {
		a.navBarHidden = true
	}
}

// navTimers restarts the tab-navigation timeout and the autohide timer.
func (a *TabApp) navTimers() tea.Cmd {
	return tea.Batch(a.startTabNavTimer(), a.startAutohideTimer())
}

// handleTabNavKey applies one key while tab navigation is active. handled is
// false for keys navigation does not use. done is true when the key finished
// navigation, which happens when enter opens a page on the compact strip.
func (a *TabApp) handleTabNavKey(key string) (handled, done bool) {
	switch a.navShape() {
	case navShapeStrip:
		switch key {
		case "left":
			a.navStripMove(-1)
		case "right":
			a.navStripMove(1)
		case "enter", "down":
			return true, a.navStripOpen()
		case "backspace", "up":
			a.navStripBack()
		default:
			return false, false
		}
	case navShapeRows:
		switch key {
		case "left":
			a.moveAtLevel(-1)
		case "right":
			a.moveAtLevel(1)
		case "up":
			a.focusLevelUp()
		case "down":
			a.navDescend()
		default:
			return false, false
		}
	default:
		switch key {
		case "up":
			a.moveAtLevel(-1)
		case "down":
			a.moveAtLevel(1)
		case "left", "backspace":
			a.focusLevelUp()
		case "right":
			a.navDescend()
		case "enter":
			return true, a.navOpen()
		default:
			return false, false
		}
	}
	return true, false
}

// navShape is the shape the navigation keys follow.
type navShape int

const (
	navShapeStrip navShape = iota // one strip of crumbs and items: the compact top bar
	navShapeRows                  // one row per level: the expanded top bar
	navShapeLists                 // vertical lists: the sidebar, Path and Columns
)

// navShape returns the key shape of the bar in its current position and mode.
func (a *TabApp) navShape() navShape {
	if a.model.NavBarPosition == NavBarLeft {
		return navShapeLists
	}
	switch a.model.NavBarMode {
	case NavBarCompact:
		return navShapeStrip
	case NavBarExpanded:
		return navShapeRows
	default:
		return navShapeLists
	}
}

// focusedIsSection reports whether the highlighted item at the focus level
// is a section with a level beneath it.
func (a *TabApp) focusedIsSection() bool {
	parent := a.sectionAtLevel(a.tabNavLevel)
	if parent == nil || a.tabNavLevel >= a.hierarchyLevels()-1 {
		return false
	}
	return a.selectedAt(a.tabNavLevel) < len(parent.Sections())
}

// navDescend moves the focus one level down when the highlighted item is a
// section. It reports whether the focus moved.
func (a *TabApp) navDescend() bool {
	if !a.focusedIsSection() {
		return false
	}
	a.tabNavLevel++
	a.navFocusCrumb = -1
	return true
}

// navOpen opens the highlighted item: a section shows its children, a page
// ends navigation. It reports whether navigation ended.
func (a *TabApp) navOpen() (done bool) {
	return !a.navDescend()
}

// navStripMove moves the strip highlight by delta positions, crumbs included,
// wrapping at both ends. Landing on an item selects it; landing on a crumb
// leaves the selection alone.
func (a *TabApp) navStripMove(delta int) {
	items := a.itemCountAtLevel(a.tabNavLevel)
	if items == 0 || a.tabNavLevel >= len(a.navPath) {
		return
	}
	crumbs := a.tabNavLevel
	total := crumbs + items
	pos := a.navFocusCrumb
	if pos < 0 {
		pos = crumbs + a.navPath[a.tabNavLevel]
	}
	pos = ((pos+delta)%total + total) % total
	if pos < crumbs {
		a.navFocusCrumb = pos
		return
	}
	a.navFocusCrumb = -1
	a.selectAtLevel(a.tabNavLevel, pos-crumbs)
}

// navStripOpen acts on the highlighted strip entry: a crumb returns to its
// level, a section opens to show its children, and a page ends navigation.
// It reports whether navigation ended.
func (a *TabApp) navStripOpen() (done bool) {
	if a.navFocusCrumb >= 0 {
		a.tabNavLevel = a.navFocusCrumb
		a.navFocusCrumb = -1
		return false
	}
	return a.navOpen()
}

// navStripBack moves the strip up one level and drops any crumb focus.
func (a *TabApp) navStripBack() {
	a.navFocusCrumb = -1
	a.focusLevelUp()
}

// TabNavLevel returns the current tab navigation level.
func (a *TabApp) TabNavLevel() int {
	return a.tabNavLevel
}
