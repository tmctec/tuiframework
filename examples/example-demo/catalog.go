package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

const gridSize = 3 // product grid is gridSize x gridSize

var gridPrices = []string{"$29", "$49", "$79"}

const catalogHelp = `Catalog

Arrows      move the product cursor
Tab         tab navigation: the bar shows its own keys
            between tabs and levels, Esc leaves
1-9         jump to a page at this level

Electronics and Home nest one level deeper
than Fashion. Deals, Kitchen and Garden are
pages beside subsections (mixed children).`

// buildCatalog builds a mixed-depth tree:
//
//	Catalog
//	  Electronics
//	    Phones      Browse, New, Used
//	    PCs         Laptop, Desktop, Tablet
//	    Deals                      (page beside the subsections)
//	  Fashion       Women, Men, Kids
//	  Home
//	    Furniture   Outdoor, Indoor, Office
//	    Kitchen, Garden            (pages beside the subsection)
//
// Names are short so the compact nav bar fits one line at 80 columns.
func buildCatalog() *tabapp.Section {
	// leaf builds a section of product pages; its pages sit at the end of
	// trail, one level below the section itself.
	leaf := func(trail ...string) *tabapp.Section {
		name := trail[len(trail)-1]
		s := tabapp.NewSection(name)
		switch name {
		case "Phones":
			s.AddPages(productPages(trail, "Browse", "New", "Used")...)
		case "PCs":
			s.AddPages(productPages(trail, "Laptop", "Desktop", "Tablet")...)
		case "Fashion":
			s.AddPages(productPages(trail, "Women", "Men", "Kids")...)
		case "Furniture":
			s.AddPages(productPages(trail, "Outdoor", "Indoor", "Office")...)
		}
		return s
	}

	// Deals sits beside Phones and PCs, whose pages are one level deeper.
	electronics := tabapp.NewSection("Electronics").
		AddSection(leaf("Catalog", "Electronics", "Phones")).
		AddSection(leaf("Catalog", "Electronics", "PCs")).
		AddPage(buildProductPage(4, "Catalog", "Electronics", "Deals"))

	// Kitchen and Garden sit beside Furniture, whose pages are one level deeper.
	home := tabapp.NewSection("Home").
		AddSection(leaf("Catalog", "Home", "Furniture")).
		AddPage(buildProductPage(4, "Catalog", "Home", "Kitchen")).
		AddPage(buildProductPage(4, "Catalog", "Home", "Garden"))

	return tabapp.NewSection("Catalog").
		AddSection(electronics).
		AddSection(leaf("Catalog", "Fashion")).
		AddSection(home)
}

// productPages builds the pages under trail; each page's branch is exactly
// one level deeper than the trail.
func productPages(trail []string, names ...string) []*page.Page {
	pages := make([]*page.Page, 0, len(names))
	for _, n := range names {
		full := append(append([]string{}, trail...), n)
		pages = append(pages, buildProductPage(len(full), full...))
	}
	return pages
}

// categoryColor gives each category an accent from base colors every palette defines.
func categoryColor(category string, p styles.Palette) color.Color {
	switch category {
	case "Electronics":
		return p.Blue
	case "Fashion":
		return p.Purple
	default:
		return p.Green
	}
}

// buildProductPage renders a product grid under its breadcrumb trail.
// trail starts at "Catalog" and ends with the page name; branchDepth is how
// deep this branch of the catalog goes, so the title can say "level 3 of 4"
// for a page that sits beside deeper subsections.
func buildProductPage(branchDepth int, trail ...string) *page.Page {
	name := trail[len(trail)-1]
	category := trail[1]
	p := page.New(name).
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			st := getState(m)
			accent := categoryColor(category, theme.Palette)
			themed(fb, theme)

			crumbs := fixedRow(fb, centered(theme.MutedStyle), 1, strings.Join(trail, "  ›  "))
			title := fixedRow(fb, centered(theme.TextHeading).Foreground(accent), 1,
				fmt.Sprintf("%s: level %d of %d", name, len(trail), branchDepth))

			status := "Arrows move the cursor"
			if m.TabNavActive {
				status = navHint(m)
			}
			grid := textRow(fb, centered(theme.ContentStyle),
				fmt.Sprintf("Products in %s\n\n%s\n\n%s", name, renderGrid(theme, st), theme.MutedStyle.Render(status)))

			fb.AddRows([]*flexbox.Row{crumbs, title, spacer(fb), grid})
		}).
		BindKey("up", "", MoveCursorMsg{DX: 0, DY: -1}).
		BindKey("down", "", MoveCursorMsg{DX: 0, DY: 1}).
		BindKey("left", "", MoveCursorMsg{DX: -1, DY: 0}).
		BindKey("right", "", MoveCursorMsg{DX: 1, DY: 0}).
		SetFooter(map[string]string{
			"←→↑↓": "Browse",
			"Tab":  "Menu",
		})
	return withHelp(p, catalogHelp)
}

// renderGrid draws the product grid with the cursor cell in LabelStyle.
func renderGrid(theme *styles.Theme, st demoState) string {
	lines := make([]string, 0, gridSize)
	for y := 0; y < gridSize; y++ {
		cells := make([]string, 0, gridSize)
		for x := 0; x < gridSize; x++ {
			cell := fmt.Sprintf(" %s ", gridPrices[x])
			if x == st.CursorX && y == st.CursorY {
				cell = theme.LabelStyle.Bold(true).Render(cell)
			} else {
				cell = theme.MutedStyle.Render(cell)
			}
			cells = append(cells, cell)
		}
		lines = append(lines, strings.Join(cells, " "))
	}
	return strings.Join(lines, "\n")
}

// clampGrid keeps a cursor coordinate inside the grid.
func clampGrid(v int) int {
	if v < 0 {
		return 0
	}
	if v >= gridSize {
		return gridSize - 1
	}
	return v
}
