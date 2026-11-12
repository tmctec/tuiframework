package tabapp

import (
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
)

// View returns the frame and the terminal features the application wants
// (Bubble Tea interface): the alternate screen, and from the Model the
// window title, terminal colors, focus reporting, mouse mode and progress
// bar. It reads the Model and declares; it changes nothing.
func (a *TabApp) View() tea.View {
	v := tea.NewView(a.Render())
	a.declare(&v)
	return v
}

// Render returns the frame as text: exactly the terminal's width and height,
// one line per row. View wraps it for Bubble Tea; tests and programs that
// embed a TabApp read it directly.
//
// The frame is a canvas of cells. The bands are drawn onto it as layers,
// each in its rectangle: header, nav bar, content and footer. The
// overlays follow in z order, lowest first. Every layer fills the cells it
// leaves without a background, so the terminal's own color never shows.
func (a *TabApp) Render() string {
	if !a.HasContent() {
		return "No content. Use AddPage(), AddSection(), or SetRoot() to add content."
	}

	a.canvas.Reset(a.width, a.height, a.model.Theme.Palette.Bg0)
	if a.model.NavBarPosition == NavBarLeft {
		a.drawSideBarLayout()
	} else {
		a.drawTopBarLayout()
	}
	for _, entry := range a.overlayStack {
		a.drawOverlay(entry.overlay)
	}
	return a.canvas.String()
}

// blockHeight returns the number of lines a rendered block takes; an empty
// block takes none.
func blockHeight(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

// drawBand draws a rendered band on the canvas at line y, full width, and
// returns the line below it. An empty band takes no room.
func (a *TabApp) drawBand(band string, y int, bg color.Color) int {
	if band == "" {
		return y
	}
	lines := strings.Split(band, "\n")
	a.canvas.Draw(lines, 0, y, a.width, len(lines), bg)
	return y + len(lines)
}

// contentFill returns the color behind the header and the page content: the
// content background, or the theme's base background when there is none.
func (a *TabApp) contentFill() color.Color {
	if a.contentBgColor == nil {
		return a.model.Theme.Palette.Bg0
	}
	return a.contentBgColor
}

// drawTopBarLayout draws the layout with the nav bar across the top:
// header and nav bar in the configured order, content, footer.
// The content takes whatever height the fixed bands leave.
func (a *TabApp) drawTopBarLayout() {
	base := a.model.Theme.Palette.Bg0
	header := a.renderHeaderSection()
	nav := a.renderNavSection()
	footer := a.renderFooterSection()
	contentHeight := a.height - blockHeight(header) - blockHeight(nav) - blockHeight(footer)

	y := 0
	if a.model.NavBarBelowHeader {
		y = a.drawBand(header, y, a.contentFill())
		y = a.drawBand(nav, y, base)
	} else {
		y = a.drawBand(nav, y, base)
		y = a.drawBand(header, y, a.contentFill())
	}
	a.drawContent(0, y, a.width, contentHeight)
	a.drawBand(footer, y+max(contentHeight, 0), base)
}

// drawSideBarLayout draws the layout with the nav bar down the left:
// header across the top, then the side bar beside the content, then footer.
func (a *TabApp) drawSideBarLayout() {
	base := a.model.Theme.Palette.Bg0
	header := a.renderHeaderSection()
	footer := a.renderFooterSection()
	contentHeight := a.height - blockHeight(header) - blockHeight(footer)

	y := a.drawBand(header, 0, a.contentFill())
	sidebarWidth := 0
	if sidebar := a.renderSideNavBar(contentHeight); sidebar != "" {
		sidebarWidth = a.navBarWidth
		a.canvas.Draw(strings.Split(sidebar, "\n"), 0, y, sidebarWidth, contentHeight, base)
	}
	a.drawContent(sidebarWidth, y, a.width-sidebarWidth, contentHeight)
	a.drawBand(footer, y+max(contentHeight, 0), base)
}

// drawContent draws the active page in the rectangle the bands leave.
func (a *TabApp) drawContent(x, y, width, height int) {
	content := a.renderContentSection(width, height)
	a.canvas.Draw(strings.Split(content, "\n"), x, y, width, height, a.contentFill())
}

// drawOverlay draws one overlay over the frame, at the anchor of its
// position. The overlay's cells without a background take its own
// background color, or the theme's base background when it has none, so the
// page never shows through and overlay content needs no color handling of
// its own.
func (a *TabApp) drawOverlay(overlay *overlays.Overlay) {
	fill := overlay.Style().GetBackground()
	if _, none := fill.(lipgloss.NoColor); none {
		fill = a.model.Theme.Palette.Bg0
	}

	// The rendered size, frame included, places the overlay.
	lines, width := overlays.GetLines(overlay.Render(a.width, a.height))
	x, y := overlays.CalculatePosition(overlay.Position(), width, len(lines), a.width, a.height)
	a.canvas.Draw(lines, x, y, width, len(lines), fill)
}

// renderHeaderSection renders the header band from the header function,
// inside the header frame. Returns "" when the app has no header.
func (a *TabApp) renderHeaderSection() string {
	height := a.model.HeaderHeight
	if a.headerFunc == nil || height <= 0 {
		return ""
	}
	frame := a.model.HeaderFrame
	width := a.width - frame.columns()

	fb := flexbox.New(width, height)
	a.headerFunc(fb, a.model)
	// The flexbox renders exactly its size. A header function that resizes
	// it gets clipped to the band.
	lines := strings.Split(fb.Render(), "\n")
	if len(lines) > height {
		lines = lines[:height]
	}
	content := strings.Join(lines, "\n")

	if frame == FrameFill {
		return a.fillBand(content, height, a.model.Theme.HeaderFill)
	}
	if !frame.bordered() {
		// The flexbox is already the band's size, and the canvas gives its
		// cells the content background.
		return content
	}

	style := borderedBox(width, height).
		Border(frame.border(a.contentBorderStyle)).
		BorderForeground(a.headerBorderColor)
	if a.contentBgColor != nil {
		style = style.Background(a.contentBgColor).BorderBackground(a.contentBgColor)
	}
	return style.Render(content)
}

// fitLines returns the lines of text for an area width cells wide. Text that
// fits is only split; wider text is wrapped by Lip Gloss, which costs far
// more than the split.
func fitLines(text string, width int) []string {
	lines, widest := overlays.GetLines(text)
	if widest <= width {
		return lines
	}
	return strings.Split(lipgloss.NewStyle().Width(width).Render(text), "\n")
}

// bandOf joins lines into a band exactly height lines tall: lines past the
// height are dropped from the given end, and missing ones are left empty for
// the canvas to fill. atBottom keeps the last lines and sets them at the
// bottom; otherwise the first lines are kept at the top.
func bandOf(lines []string, height int, atBottom bool) string {
	if len(lines) > height {
		if atBottom {
			lines = lines[len(lines)-height:]
		} else {
			lines = lines[:height]
		}
	}
	text := strings.Join(lines, "\n")
	pad := strings.Repeat("\n", max(height-len(lines), 0))
	if atBottom {
		return pad + text
	}
	return text + pad
}

// borderColors colors the border of style: a gradient through the border
// blend when the app has one, the single color c otherwise.
func (a *TabApp) borderColors(style lipgloss.Style, c color.Color) lipgloss.Style {
	if len(a.borderBlend) >= 2 {
		return style.BorderForegroundBlend(a.borderBlend...)
	}
	return style.BorderForeground(c)
}

// borderedBox returns a style for a bordered box whose content area is
// width by height cells. Lip Gloss counts a border inside a style's Width
// and Height, so the box is two cells larger each way; the caller sets the
// border.
func borderedBox(width, height int) lipgloss.Style {
	return lipgloss.NewStyle().Width(width + 2).Height(height + 2)
}

// fillBand renders a band as a solid bar: the full width of the screen in
// the colors of fill, with the content inset one column on each side. Cells
// the content left without a background take the fill color.
func (a *TabApp) fillBand(content string, height int, fill lipgloss.Style) string {
	bar := fill.
		Width(a.width).
		Height(height).
		Padding(0, 1).
		Render(content)
	lines := overlays.FillBackground(strings.Split(bar, "\n"), a.width, fill.GetBackground())
	return strings.Join(lines, "\n")
}

// renderNavSection renders the fixed-height nav bar with the tabs.
// Returns "" while autohide has the bar hidden.
func (a *TabApp) renderNavSection() string {
	if a.model.NavBarAutoHide && a.navBarHidden {
		return ""
	}

	width := a.width
	if a.model.ShowNavBarBorder {
		width -= 2 // account for left + right border
	}

	// Path and Columns size the band themselves and are cut at the width;
	// the other layouts wrap to the width and clip to the configured bar
	// height (top justified).
	bar := a.renderNavBar()
	height := a.model.NavBarHeight
	sizesItself := a.model.NavBarMode == NavBarPath || a.model.NavBarMode == NavBarColumns

	if !a.model.ShowNavBarBorder {
		// The canvas cuts each line at the width and fills what is short.
		if sizesItself {
			return bar
		}
		return bandOf(fitLines(bar, width), height, false)
	}

	var lines []string
	if sizesItself {
		lines = strings.Split(lipgloss.NewStyle().MaxWidth(width).Render(bar), "\n")
		height = len(lines)
	} else {
		lines = fitLines(bar, width)
		if len(lines) > height {
			lines = lines[:height]
		}
	}
	navStyle := a.borderColors(borderedBox(width, height).Border(a.navBarBorderStyle), a.navBarBorderColor)
	return navStyle.Render(strings.Join(lines, "\n"))
}

// renderContentSection renders the active page for a content area of width
// by height cells: the page's flexbox, inside the content border when it is
// on. A page's own border color and style override the app's. Without a
// border the canvas fits the text to the area, so it is returned as the
// flexbox rendered it.
func (a *TabApp) renderContentSection(width, height int) string {
	currentPage := a.currentPage()
	if currentPage == nil {
		return "No content"
	}
	if !a.model.ShowContentBorder {
		return a.renderPageContent(currentPage, width, height).Render()
	}

	// The frame takes 2 columns and 2 rows of the section, so both the
	// style's inner size and the page get width-2 x height-2.
	innerWidth, innerHeight := width-2, height-2
	border := a.contentBorderStyle
	if b := currentPage.BorderStyle(); b != nil {
		border = *b
	}
	style := borderedBox(innerWidth, innerHeight).Border(border)
	if c := currentPage.BorderColor(); c != nil {
		style = style.BorderForeground(c) // the page's own color beats a gradient
	} else {
		style = a.borderColors(style, a.contentBorderColor)
	}
	if a.contentBgColor != nil {
		style = style.Background(a.contentBgColor).BorderBackground(a.contentBgColor)
	}
	return style.Render(a.renderPageContent(currentPage, innerWidth, innerHeight).Render())
}

// renderPageContent renders the page content (Elm architecture: View receives Model).
func (a *TabApp) renderPageContent(p *page.Page, width, height int) *flexbox.FlexBox {
	fb := flexbox.New(width, height)
	if p.Content() != nil {
		p.Content()(fb, a.model, a.model.NavBarHeight, a.model.FooterHeight)
	}
	return fb
}

// renderFooterSection renders the fixed-height footer inside the footer
// frame. Returns "" when the footer height is 0.
func (a *TabApp) renderFooterSection() string {
	height := a.model.FooterHeight
	if height <= 0 {
		return ""
	}
	frame := a.model.FooterFrame
	width := a.width - frame.columns()

	// Wrap the keys to the width and keep the last lines, so the footer is
	// bottom justified.
	lines := fitLines(a.renderFooterContent(), width)
	if len(lines) > height {
		lines = lines[len(lines)-height:]
	}

	// A filled footer keeps the shape of the footer style and takes the
	// colors of the fill.
	text := a.footerStyle
	fill := a.model.Theme.FooterFill
	if frame == FrameFill {
		text = text.Foreground(fill.GetForeground()).Background(fill.GetBackground())
	}
	body := text.Render(strings.Join(lines, "\n"))
	if frame == FrameNone {
		// The canvas fills the lines above the keys and the cells beside them.
		return bandOf(strings.Split(body, "\n"), height, true)
	}

	placed := lipgloss.Place(width, height, lipgloss.Left, lipgloss.Bottom, body)
	if frame == FrameFill {
		return a.fillBand(placed, height, fill)
	}
	style := borderedBox(width, height).
		Border(frame.border(a.footerBorderStyle)).
		BorderForeground(a.footerBorderColor)
	return style.Render(placed)
}
