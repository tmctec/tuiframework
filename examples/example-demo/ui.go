package main

import (
	"fmt"
	"slices"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// withHelp binds ? on p to a help modal showing text, above every other
// overlay.
func withHelp(p *page.Page, text string) *page.Page {
	return p.BindKey("?", "Help", tabapp.ShowModalMsg{Title: "Help", Content: text, ZIndex: 100, SizePct: 55})
}

// themed gives fb the theme's base colors and passes them down to every row
// and cell, so unstyled text renders in Fg0 on Bg0.
func themed(fb *flexbox.FlexBox, theme *styles.Theme) {
	fb.UseStyle(theme.BaseStyle)
}

// textRow builds a single-cell row. The row passes the inherited theme style
// on to its cell. The caller adds the row to the flexbox.
func textRow(fb *flexbox.FlexBox, style lipgloss.Style, text string) *flexbox.Row {
	row := fb.NewRow().StylePassing(true)
	cell := flexbox.NewCell(1, 1)
	cell.SetStyle(style)
	cell.SetContent(text)
	row.AddCells(cell)
	return row
}

// fixedRow is textRow with a fixed height in lines.
func fixedRow(fb *flexbox.FlexBox, style lipgloss.Style, height int, text string) *flexbox.Row {
	return textRow(fb, style, text).SetFixedHeight(height)
}

// card is a filled, borderless panel: Fg1 text on Bg1 with padding. Pages
// use it instead of bordered boxes so the only frame around content is the
// content border setting.
func card(theme *styles.Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Background(theme.Palette.Bg1).
		Foreground(theme.Palette.Fg1).
		Padding(1, 2)
}

// spacer is an empty row one line high.
func spacer(fb *flexbox.FlexBox) *flexbox.Row {
	return fixedRow(fb, lipgloss.NewStyle(), 1, "")
}

// filler is an empty row that absorbs the height the fixed rows leave over.
func filler(fb *flexbox.FlexBox) *flexbox.Row {
	return textRow(fb, lipgloss.NewStyle(), "")
}

// block pads every line of text to the width of the widest one, so a cell
// that centers its content keeps the lines' columns aligned.
func block(text string) string {
	lines := strings.Split(text, "\n")
	width := 0
	for _, l := range lines {
		if w := lipgloss.Width(l); w > width {
			width = w
		}
	}
	for i, l := range lines {
		lines[i] = l + strings.Repeat(" ", width-lipgloss.Width(l))
	}
	return strings.Join(lines, "\n")
}

// centered returns style with centered content.
func centered(style lipgloss.Style) lipgloss.Style {
	return style.Align(lipgloss.Center)
}

func onOff(v bool) string {
	if v {
		return "On"
	}
	return "Off"
}

func positionName(p tabapp.NavBarPosition) string {
	if p == tabapp.NavBarLeft {
		return "Left"
	}
	return "Top"
}

func modeName(m tabapp.NavBarMode) string {
	return m.String()
}

// navModes lists the nav bar layouts in picker order.
var navModes = []tabapp.NavBarMode{tabapp.NavBarCompact, tabapp.NavBarExpanded, tabapp.NavBarPath, tabapp.NavBarColumns}

// cycleMode returns the layout delta steps after m in navModes, wrapping.
func cycleMode(m tabapp.NavBarMode, delta int) tabapp.NavBarMode {
	return cycle(navModes, m, delta)
}

// bandFrames lists the header and footer frames in picker order.
var bandFrames = []tabapp.Frame{tabapp.FrameNone, tabapp.FrameSingle, tabapp.FrameDouble, tabapp.FrameFill}

// cycle returns the item delta steps after cur in items, wrapping. An item
// outside the list counts as the first.
func cycle[T comparable](items []T, cur T, delta int) T {
	i := max(slices.Index(items, cur), 0)
	return items[moveCursor(i, delta, len(items))]
}

// frameName names a header or footer frame the way the Settings popup shows it.
func frameName(f tabapp.Frame) string {
	switch f {
	case tabapp.FrameSingle:
		return "Single line"
	case tabapp.FrameDouble:
		return "Double line"
	default:
		return f.String()
	}
}

// heightName names a band height in lines; 0 reads as Hidden.
func heightName(lines int) string {
	switch lines {
	case 0:
		return "Hidden"
	case 1:
		return "1 line"
	default:
		return fmt.Sprintf("%d lines", lines)
	}
}

// navHint names the tab-navigation keys for the bar's current shape.
func navHint(m *tabapp.Model) string {
	switch {
	case m.NavBarPosition == tabapp.NavBarLeft:
		return "Tab navigation: ▴▾ move, ◂▸ change level"
	case m.NavBarMode == tabapp.NavBarExpanded:
		return "Tab navigation: ◂▸ move, ▴▾ change level"
	case m.NavBarMode == tabapp.NavBarCompact:
		return "Tab navigation: ◂▸ move, Enter opens, Backspace goes back"
	default:
		return "Tab navigation: ▴▾ item, ◂▸ level, Enter opens, Backspace goes back"
	}
}
