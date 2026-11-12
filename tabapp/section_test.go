package tabapp

import (
	"testing"

	"github.com/tmctec/tuiframework/page"
)

func TestSectionNewAndRename(t *testing.T) {
	s := NewSection("Catalog")
	if s.Sections() == nil || s.Pages() == nil {
		t.Error("a new section has empty, non-nil child slices")
	}
	if s.SetName("Shop") != s || s.Name() != "Shop" {
		t.Errorf("SetName should rename in place and return the section, got %q", s.Name())
	}
}

func TestSectionAddSectionsVariadic(t *testing.T) {
	s := NewSection("Root").AddSections(NewSection("A"), NewSection("B"), NewSection("C"))
	if len(s.Sections()) != 3 || s.Sections()[2].Name() != "C" {
		t.Errorf("Sections = %d, want A, B, C", len(s.Sections()))
	}
	s.AddSection(NewSection("D"))
	if len(s.Sections()) != 4 {
		t.Errorf("Sections = %d after AddSection, want 4", len(s.Sections()))
	}
}

func TestSectionChildPredicates(t *testing.T) {
	empty := NewSection("Empty")
	pagesOnly := NewSection("Pages").AddPages(page.New("p1"), page.New("p2"))
	sectionsOnly := NewSection("Sections").AddSection(NewSection("Child"))
	mixed := NewSection("Mixed").AddSection(NewSection("Child")).AddPage(page.New("p"))

	cases := []struct {
		name                                     string
		s                                        *Section
		hasChildren, hasSections, hasPages, leaf bool
	}{
		{"empty", empty, false, false, false, true},
		{"pages only", pagesOnly, true, false, true, true},
		{"sections only", sectionsOnly, true, true, false, false},
		{"mixed", mixed, true, true, true, false},
	}
	for _, c := range cases {
		if c.s.HasChildren() != c.hasChildren {
			t.Errorf("%s: HasChildren = %v", c.name, c.s.HasChildren())
		}
		if c.s.HasSections() != c.hasSections {
			t.Errorf("%s: HasSections = %v", c.name, c.s.HasSections())
		}
		if c.s.HasPages() != c.hasPages {
			t.Errorf("%s: HasPages = %v", c.name, c.s.HasPages())
		}
		if c.s.IsLeaf() != c.leaf {
			t.Errorf("%s: IsLeaf = %v", c.name, c.s.IsLeaf())
		}
	}
}

func TestSectionDepthOfMixedTrees(t *testing.T) {
	if d := NewSection("Empty").Depth(); d != 1 {
		t.Errorf("empty section depth = %d, want 1", d)
	}

	shallow := NewSection("Shallow").AddPage(page.New("p"))
	deep := NewSection("Deep").AddSection(NewSection("Deeper").AddPage(page.New("p")))
	root := NewSection("Root").AddSections(shallow, deep).AddPage(page.New("root page"))
	if d := root.Depth(); d != 3 {
		t.Errorf("depth = %d, want 3 (the deepest branch wins)", d)
	}

	mixed := NewSection("Mixed").AddSection(NewSection("Child")).AddPages(page.New("a"), page.New("b"))
	if d := mixed.Depth(); d != 2 {
		t.Errorf("depth = %d, want 2 (pages beside a section add no level)", d)
	}
}

func TestSectionItemCount(t *testing.T) {
	if n := NewSection("Empty").ItemCount(); n != 0 {
		t.Errorf("empty ItemCount = %d, want 0", n)
	}
	pages := NewSection("Pages").AddPages(page.New("a"), page.New("b"), page.New("c"))
	if n := pages.ItemCount(); n != 3 {
		t.Errorf("pages ItemCount = %d, want 3", n)
	}
	mixed := NewSection("Mixed").AddSections(NewSection("x"), NewSection("y")).AddPages(page.New("a"), page.New("b"), page.New("c"))
	if n := mixed.ItemCount(); n != 5 {
		t.Errorf("mixed ItemCount = %d, want 5: two sections and three pages", n)
	}
}
