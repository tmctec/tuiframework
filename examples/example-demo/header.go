package main

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/tabapp"
)

// renderHeader draws the header band: the app name on the left and the
// current theme and nav bar settings on the right, centered in the band's
// height. A filled header draws on the theme's HeaderFill, so its text is
// reversed to suit the fill. It reads the Model on every render, so it
// follows the Settings popup live.
func renderHeader(fb *flexbox.FlexBox, model interface{}) {
	m := model.(*tabapp.Model)
	theme := m.Theme

	base, heading, muted := theme.BaseStyle, theme.TextHeading, theme.MutedStyle
	if m.HeaderFrame == tabapp.FrameFill {
		base = theme.HeaderFill
		heading, muted = base.Bold(true), base
	}
	fb.UseStyle(base)

	nav := positionName(m.NavBarPosition) + ", " + modeName(m.NavBarMode)
	if m.NavBarPosition == tabapp.NavBarTop && !m.NavBarBelowHeader {
		nav = "above, " + modeName(m.NavBarMode)
	}
	status := fmt.Sprintf("theme %s  ·  nav %s  ·  autohide %s", theme.Name(), nav, onOff(m.NavBarAutoHide))

	// The one line of text sits in the middle of a taller band.
	row := fb.NewRow().StylePassing(true)
	title := flexbox.NewCell(1, 1)
	title.SetStyle(heading.AlignVertical(lipgloss.Center).PaddingLeft(1))
	title.SetContent("tuiframework demo")
	info := flexbox.NewCell(2, 1)
	info.SetStyle(muted.Align(lipgloss.Right, lipgloss.Center).PaddingRight(1))
	info.SetContent(status)
	row.AddCells(title, info)
	fb.AddRows([]*flexbox.Row{row})
}
