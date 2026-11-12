package tabapp

import (
	"slices"

	"github.com/tmctec/tuiframework/page"
)

// Section is a node of the navigation tree. It holds child sections, pages,
// or both, so the tree nests to any depth: pages alone make one level,
// sections of pages two, and so on. The nav bar lists a section's
// child sections first and then its pages.
type Section struct {
	name     string
	sections []*Section   // Child sections (for nesting)
	pages    []*page.Page // Pages (leaf nodes)
}

// NewSection creates a new section with the given name.
func NewSection(name string) *Section {
	return &Section{
		name:     name,
		sections: make([]*Section, 0),
		pages:    make([]*page.Page, 0),
	}
}

// Name returns the section name.
func (s *Section) Name() string {
	return s.name
}

// SetName sets the section name and returns self for chaining.
func (s *Section) SetName(name string) *Section {
	s.name = name
	return s
}

// Sections returns the child sections. The slice is the section's own:
// read it, and add to it with AddSection.
func (s *Section) Sections() []*Section {
	return s.sections
}

// Pages returns the pages in this section. The slice is the section's own:
// read it, and add to it with AddPage.
func (s *Section) Pages() []*page.Page {
	return s.pages
}

// AddSection adds a child section and returns self for chaining.
func (s *Section) AddSection(child *Section) *Section {
	s.sections = append(s.sections, child)
	return s
}

// AddSections adds multiple child sections and returns self for chaining.
func (s *Section) AddSections(children ...*Section) *Section {
	s.sections = append(s.sections, children...)
	return s
}

// AddPage adds a page to this section and returns self for chaining.
func (s *Section) AddPage(p *page.Page) *Section {
	s.pages = append(s.pages, p)
	return s
}

// AddPages adds multiple pages to this section and returns self for chaining.
func (s *Section) AddPages(pages ...*page.Page) *Section {
	s.pages = append(s.pages, pages...)
	return s
}

// HasChildren returns true if section has any content (sections or pages).
func (s *Section) HasChildren() bool {
	return len(s.sections) > 0 || len(s.pages) > 0
}

// HasSections returns true if section has child sections.
func (s *Section) HasSections() bool {
	return len(s.sections) > 0
}

// HasPages returns true if section has pages.
func (s *Section) HasPages() bool {
	return len(s.pages) > 0
}

// IsLeaf returns true if section has no child sections (only pages or empty).
func (s *Section) IsLeaf() bool {
	return len(s.sections) == 0
}

// Depth returns the maximum depth of this section tree.
// Returns 1 for a leaf section (pages only), >1 for nested sections.
// An empty section returns 1.
func (s *Section) Depth() int {
	if len(s.sections) == 0 {
		// Leaf section - if it has pages, depth is 1
		// If empty, still consider it depth 1 (pages level)
		return 1
	}
	maxChildDepth := 0
	for _, child := range s.sections {
		if d := child.Depth(); d > maxChildDepth {
			maxChildDepth = d
		}
	}
	return 1 + maxChildDepth
}

// ItemCount returns the number of items the nav bar shows for this
// section: its child sections and its pages.
func (s *Section) ItemCount() int {
	return len(s.sections) + len(s.pages)
}

// Breadcrumb is an entry in the trail EnterSection pushes: the position the
// app was at before it entered a section, so PopToBreadcrumb can return there.
type Breadcrumb struct {
	level int   // the level EnterSection was called for
	path  []int // the navigation path before the move, one index per level
}

// Level returns the level the section was entered at.
func (b Breadcrumb) Level() int {
	return b.level
}

// SectionIndex returns the item that was selected at the entered level.
func (b Breadcrumb) SectionIndex() int {
	return b.path[b.level]
}

// PageIndex returns the index that was selected at the deepest level.
func (b Breadcrumb) PageIndex() int {
	return b.path[len(b.path)-1]
}

// Path returns a copy of the navigation path before the move.
func (b Breadcrumb) Path() []int {
	return slices.Clone(b.path)
}
