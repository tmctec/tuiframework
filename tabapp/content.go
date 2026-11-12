package tabapp

import (
	"fmt"

	"github.com/tmctec/tuiframework/page"
)

// Pages returns the pages at the active page's level: the active page and
// its siblings.
func (a *TabApp) Pages() []*page.Page {
	if parent, _, _ := a.pageLevel(); parent != nil {
		return parent.Pages()
	}
	return nil
}

// Root returns the root section.
func (a *TabApp) Root() *Section {
	return a.root
}

// SetRoot replaces the whole navigation tree with s and returns to its first
// page. A nil root clears the content. Returns an error, and leaves the tree
// alone, if a page key conflicts with a universal or hardcoded binding.
func (a *TabApp) SetRoot(s *Section) error {
	if s != nil {
		if err := a.validateSectionPages(s); err != nil {
			return err
		}
	}
	a.root = s
	a.resetNavigation()
	return nil
}

// MustSetRoot replaces the navigation tree, panicking on error.
// Returns self for chaining.
func (a *TabApp) MustSetRoot(s *Section) *TabApp {
	if err := a.SetRoot(s); err != nil {
		panic(err)
	}
	return a
}

// resetNavigation points the navigation at the first item of every level of
// the current tree, with the bar on the root level, and clears the
// breadcrumb trail.
func (a *TabApp) resetNavigation() {
	a.navPath = make([]int, a.hierarchyLevels())
	a.breadcrumbs = a.breadcrumbs[:0]
	a.navFocusCrumb = -1
	a.tabNavLevel = 0
	a.syncActivePage()
}

// validatePageKeys checks if any page keys conflict with universal/hardcoded keys.
func (a *TabApp) validatePageKeys(p *page.Page) error {
	for key := range p.Bindings() {
		if _, isHardcoded := a.hardcodedKeys[key]; isHardcoded {
			return fmt.Errorf("key %q is hardcoded and cannot be overridden", key)
		}
		if _, isUniversal := a.universalBindings[key]; isUniversal {
			return fmt.Errorf("key %q conflicts with universal binding (use UnbindKey to allow)", key)
		}
	}
	return nil
}

// validateSectionPages recursively validates all pages in a section tree.
func (a *TabApp) validateSectionPages(s *Section) error {
	// Validate pages in this section
	for _, p := range s.Pages() {
		if err := a.validatePageKeys(p); err != nil {
			return err
		}
	}
	// Validate child sections
	for _, child := range s.Sections() {
		if err := a.validateSectionPages(child); err != nil {
			return err
		}
	}
	return nil
}

// AddSection adds a section to the root.
// If root is nil, creates a root section to hold added sections.
// Returns error if any page keys conflict.
func (a *TabApp) AddSection(s *Section) error {
	if err := a.validateSectionPages(s); err != nil {
		return err
	}
	if a.root == nil {
		// Create anonymous root
		a.root = NewSection("")
	}
	a.root.AddSection(s)
	// A deeper tree needs a longer path; the position starts over.
	if len(a.navPath) != a.root.Depth() {
		a.resetNavigation()
	}
	return nil
}

// MustAddSection adds a section, panicking on error.
// Returns self for chaining.
func (a *TabApp) MustAddSection(s *Section) *TabApp {
	if err := a.AddSection(s); err != nil {
		panic(err)
	}
	return a
}

// AddPage adds a page directly to the root section (flat page mode).
// Creates an anonymous root section if needed.
// Returns error if any page keys conflict.
func (a *TabApp) AddPage(p *page.Page) error {
	if err := a.validatePageKeys(p); err != nil {
		return err
	}
	if a.root == nil {
		a.root = NewSection("")
	}
	a.root.AddPage(p)
	return nil
}

// MustAddPage adds a page, panicking on error.
// Returns self for chaining.
func (a *TabApp) MustAddPage(p *page.Page) *TabApp {
	if err := a.AddPage(p); err != nil {
		panic(err)
	}
	return a
}

// HasSections returns true if the app has section-based navigation configured.
// This distinguishes from flat pages mode (where pages are added directly to root).
func (a *TabApp) HasSections() bool {
	return a.root != nil && a.root.HasSections()
}

// HasContent returns true if the app has any content (sections or pages).
func (a *TabApp) HasContent() bool {
	return a.root != nil && (a.root.HasSections() || a.root.HasPages())
}
