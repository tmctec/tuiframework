package main

import (
	"fmt"
	"slices"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// settingsRow identifies one row of the Settings page.
type settingsRow int

const (
	rowTheme settingsRow = iota
	rowPosition
	rowNavOrder
	rowMode
	rowAutoHide
	rowAutoHideTimeout
	rowTabNavTimeout
	rowHeaderHeight
	rowHeaderFrame
	rowNavBorder
	rowContentBorder
	rowFooterHeight
	rowFooterFrame
	rowKeyGlow
	rowTerminalColors
	rowFollowBackground
	rowCount
)

// settingsSnapshot is the Settings page's view of library state, read from the
// Model and never written by the demo. The reducer works on snapshots and
// emits the library messages that move the app from one to the next.
type settingsSnapshot struct {
	ThemeIndex        int
	Position          tabapp.NavBarPosition
	NavBarBelowHeader bool
	Mode              tabapp.NavBarMode
	AutoHide          bool
	AutoHideTimeout   time.Duration
	TabNavTimeout     time.Duration
	HeaderHeight      int
	HeaderFrame       tabapp.Frame
	NavBorder         bool
	ContentBorder     bool
	FooterHeight      int
	FooterFrame       tabapp.Frame
	KeyGlow           bool
	TerminalColors    bool
	FollowBackground  bool
	DarkBackground    bool // what the terminal reported; decides which theme following picks
}

// defaultSettings matches the configuration buildApp starts with.
var defaultSettings = settingsSnapshot{
	ThemeIndex:        0,
	Position:          tabapp.NavBarTop,
	NavBarBelowHeader: true,
	Mode:              tabapp.NavBarCompact,
	AutoHide:          true,
	AutoHideTimeout:   2 * time.Second,
	TabNavTimeout:     2 * time.Second,
	HeaderHeight:      headerHeight,
	FooterHeight:      footerHeight,
	KeyGlow:           true,
	TerminalColors:    true,
	DarkBackground:    true,
}

const (
	autoHideStep   = time.Second
	autoHideMin    = time.Second
	autoHideMax    = 30 * time.Second
	tabNavStep     = 500 * time.Millisecond
	tabNavMin      = 500 * time.Millisecond
	tabNavMax      = 10 * time.Second
	bandHeightMax  = 5  // tallest header or footer, in lines of content; 0 hides the band
	settingsName   = 20 // width of the name column
	settingsValue  = 28 // width of the value column
	settingsIndent = 2  // room for the selection pill's padding
)

func snapshotOf(m *tabapp.Model) settingsSnapshot {
	return settingsSnapshot{
		ThemeIndex:        themeIndex(m),
		Position:          m.NavBarPosition,
		NavBarBelowHeader: m.NavBarBelowHeader,
		Mode:              m.NavBarMode,
		AutoHide:          m.NavBarAutoHide,
		AutoHideTimeout:   m.NavBarTimeout,
		TabNavTimeout:     m.TabNavTimeout,
		HeaderHeight:      m.HeaderHeight,
		HeaderFrame:       m.HeaderFrame,
		NavBorder:         m.ShowNavBarBorder,
		ContentBorder:     m.ShowContentBorder,
		FooterHeight:      m.FooterHeight,
		FooterFrame:       m.FooterFrame,
		KeyGlow:           m.KeyGlowEnabled,
		TerminalColors:    m.TerminalColors,
		FollowBackground:  m.FollowsBackground,
		DarkBackground:    m.DarkBackground,
	}
}

// adaptivePair returns the light and dark themes that following the
// terminal's background switches between: the theme and its counterpart,
// or the first theme's family for a theme that has none.
func adaptivePair(t *styles.Theme) (light, dark *styles.Theme) {
	other := styles.Counterpart(t)
	if other == nil {
		t = demoThemes[0]
		other = styles.Counterpart(t)
	}
	if t.IsDark() {
		return other, t
	}
	return t, other
}

// moveCursor steps cur by delta inside [0, n), wrapping at both ends.
func moveCursor(cur, delta, n int) int {
	return ((cur+delta)%n + n) % n
}

// clampHeight keeps a header or footer height inside [0, bandHeightMax].
func clampHeight(h int) int {
	return min(max(h, 0), bandHeightMax)
}

func clampDuration(d, lo, hi time.Duration) time.Duration {
	if d < lo {
		return lo
	}
	if d > hi {
		return hi
	}
	return d
}

// adjustRow applies delta steps to row. It returns the snapshot the app will
// reach and the library messages that take it there.
func adjustRow(s settingsSnapshot, row settingsRow, delta int) (settingsSnapshot, []tea.Msg) {
	switch row {
	case rowTheme:
		s.ThemeIndex = moveCursor(s.ThemeIndex, delta, len(demoThemes))
		s.FollowBackground = false // a theme picked by hand stays
		return s, []tea.Msg{tabapp.SetThemeMsg{Theme: demoThemes[s.ThemeIndex]}}
	case rowPosition:
		if s.Position == tabapp.NavBarTop {
			s.Position = tabapp.NavBarLeft
		} else {
			s.Position = tabapp.NavBarTop
		}
		return s, []tea.Msg{tabapp.ToggleNavBarPositionMsg{LeftWidth: sideBarWidth}}
	case rowNavOrder:
		if s.Position != tabapp.NavBarTop {
			return s, nil // a side nav bar has no above/below
		}
		s.NavBarBelowHeader = !s.NavBarBelowHeader
		return s, []tea.Msg{tabapp.ToggleNavBarBelowHeaderMsg{}}
	case rowMode:
		s.Mode = cycleMode(s.Mode, delta)
		return s, []tea.Msg{tabapp.SetNavBarModeMsg{Mode: s.Mode}}
	case rowAutoHide:
		s.AutoHide = !s.AutoHide
		return s, []tea.Msg{tabapp.ToggleNavBarAutoHideMsg{}}
	case rowAutoHideTimeout:
		step := time.Duration(delta) * autoHideStep
		s.AutoHideTimeout = clampDuration(s.AutoHideTimeout+step, autoHideMin, autoHideMax)
		return s, []tea.Msg{tabapp.AdjustNavBarTimeoutMsg{Delta: step}}
	case rowTabNavTimeout:
		step := time.Duration(delta) * tabNavStep
		s.TabNavTimeout = clampDuration(s.TabNavTimeout+step, tabNavMin, tabNavMax)
		return s, []tea.Msg{tabapp.AdjustTabNavTimeoutMsg{Delta: step}}
	case rowHeaderHeight:
		s.HeaderHeight = clampHeight(s.HeaderHeight + delta)
		return s, []tea.Msg{tabapp.SetHeaderHeightMsg{Height: s.HeaderHeight}}
	case rowHeaderFrame:
		s.HeaderFrame = cycle(bandFrames, s.HeaderFrame, delta)
		return s, []tea.Msg{tabapp.SetHeaderFrameMsg{Frame: s.HeaderFrame}}
	case rowNavBorder:
		s.NavBorder = !s.NavBorder
		return s, []tea.Msg{tabapp.ToggleNavBarBorderMsg{}}
	case rowContentBorder:
		s.ContentBorder = !s.ContentBorder
		return s, []tea.Msg{tabapp.ToggleContentBorderMsg{}}
	case rowFooterHeight:
		s.FooterHeight = clampHeight(s.FooterHeight + delta)
		return s, []tea.Msg{tabapp.SetFooterHeightMsg{Height: s.FooterHeight}}
	case rowFooterFrame:
		s.FooterFrame = cycle(bandFrames, s.FooterFrame, delta)
		return s, []tea.Msg{tabapp.SetFooterFrameMsg{Frame: s.FooterFrame}}
	case rowKeyGlow:
		s.KeyGlow = !s.KeyGlow
		return s, []tea.Msg{tabapp.ToggleKeyGlowMsg{}}
	case rowTerminalColors:
		s.TerminalColors = !s.TerminalColors
		return s, []tea.Msg{tabapp.ToggleTerminalColorsMsg{}}
	case rowFollowBackground:
		if s.FollowBackground {
			// Choosing the theme that is active ends the following.
			s.FollowBackground = false
			return s, []tea.Msg{tabapp.SetThemeMsg{Theme: demoThemes[s.ThemeIndex]}}
		}
		light, dark := adaptivePair(demoThemes[s.ThemeIndex])
		s.FollowBackground = true
		s.ThemeIndex = slices.Index(demoThemes, light)
		if s.DarkBackground {
			s.ThemeIndex = slices.Index(demoThemes, dark)
		}
		return s, []tea.Msg{tabapp.UseAdaptiveThemeMsg{Light: light, Dark: dark}}
	}
	return s, nil
}

// resetMsgs returns the messages that take s to def: one per differing row.
func resetMsgs(s, def settingsSnapshot) []tea.Msg {
	var msgs []tea.Msg
	if s.ThemeIndex != def.ThemeIndex {
		msgs = append(msgs, tabapp.SetThemeMsg{Theme: demoThemes[def.ThemeIndex]})
	}
	if s.Position != def.Position {
		msgs = append(msgs, tabapp.ToggleNavBarPositionMsg{LeftWidth: sideBarWidth})
	}
	if s.NavBarBelowHeader != def.NavBarBelowHeader {
		msgs = append(msgs, tabapp.ToggleNavBarBelowHeaderMsg{})
	}
	if s.Mode != def.Mode {
		msgs = append(msgs, tabapp.SetNavBarModeMsg{Mode: def.Mode})
	}
	if s.AutoHide != def.AutoHide {
		msgs = append(msgs, tabapp.ToggleNavBarAutoHideMsg{})
	}
	if s.AutoHideTimeout != def.AutoHideTimeout {
		msgs = append(msgs, tabapp.SetNavBarTimeoutMsg{Duration: def.AutoHideTimeout})
	}
	if s.TabNavTimeout != def.TabNavTimeout {
		msgs = append(msgs, tabapp.SetTabNavTimeoutMsg{Duration: def.TabNavTimeout})
	}
	if s.HeaderHeight != def.HeaderHeight {
		msgs = append(msgs, tabapp.SetHeaderHeightMsg{Height: def.HeaderHeight})
	}
	if s.HeaderFrame != def.HeaderFrame {
		msgs = append(msgs, tabapp.SetHeaderFrameMsg{Frame: def.HeaderFrame})
	}
	if s.NavBorder != def.NavBorder {
		msgs = append(msgs, tabapp.ToggleNavBarBorderMsg{})
	}
	if s.ContentBorder != def.ContentBorder {
		msgs = append(msgs, tabapp.ToggleContentBorderMsg{})
	}
	if s.FooterHeight != def.FooterHeight {
		msgs = append(msgs, tabapp.SetFooterHeightMsg{Height: def.FooterHeight})
	}
	if s.FooterFrame != def.FooterFrame {
		msgs = append(msgs, tabapp.SetFooterFrameMsg{Frame: def.FooterFrame})
	}
	if s.KeyGlow != def.KeyGlow {
		msgs = append(msgs, tabapp.ToggleKeyGlowMsg{})
	}
	if s.TerminalColors != def.TerminalColors {
		msgs = append(msgs, tabapp.ToggleTerminalColorsMsg{})
	}
	if s.FollowBackground != def.FollowBackground {
		// The default is a theme picked by hand, which ends the following.
		msgs = append(msgs, tabapp.SetThemeMsg{Theme: demoThemes[def.ThemeIndex]})
	}
	return msgs
}

func rowLabel(row settingsRow) string {
	switch row {
	case rowTheme:
		return "Theme"
	case rowPosition:
		return "Nav bar position"
	case rowNavOrder:
		return "Nav bar order"
	case rowMode:
		return "Nav bar mode"
	case rowAutoHide:
		return "Nav bar autohide"
	case rowAutoHideTimeout:
		return "Autohide timeout"
	case rowTabNavTimeout:
		return "Tab-nav timeout"
	case rowHeaderHeight:
		return "Header height"
	case rowHeaderFrame:
		return "Header border"
	case rowNavBorder:
		return "Nav bar border"
	case rowContentBorder:
		return "Content border"
	case rowFooterHeight:
		return "Footer height"
	case rowFooterFrame:
		return "Footer border"
	case rowKeyGlow:
		return "Key glow"
	case rowTerminalColors:
		return "Terminal colors"
	case rowFollowBackground:
		return "Follow terminal"
	}
	return ""
}

func rowValue(s settingsSnapshot, row settingsRow) string {
	switch row {
	case rowTheme:
		return fmt.Sprintf("%s (%d/%d)", demoThemes[s.ThemeIndex].Name(), s.ThemeIndex+1, len(demoThemes))
	case rowPosition:
		return positionName(s.Position)
	case rowNavOrder:
		switch {
		case s.Position != tabapp.NavBarTop:
			return "Sidebar"
		case s.NavBarBelowHeader:
			return "Below header"
		default:
			return "Above header"
		}
	case rowMode:
		return modeName(s.Mode)
	case rowAutoHide:
		return onOff(s.AutoHide)
	case rowAutoHideTimeout:
		return fmt.Sprintf("%.1f s", s.AutoHideTimeout.Seconds())
	case rowTabNavTimeout:
		return fmt.Sprintf("%.1f s", s.TabNavTimeout.Seconds())
	case rowHeaderHeight:
		return heightName(s.HeaderHeight)
	case rowHeaderFrame:
		return frameName(s.HeaderFrame)
	case rowNavBorder:
		return onOff(s.NavBorder)
	case rowContentBorder:
		return onOff(s.ContentBorder)
	case rowFooterHeight:
		return heightName(s.FooterHeight)
	case rowFooterFrame:
		return frameName(s.FooterFrame)
	case rowKeyGlow:
		return onOff(s.KeyGlow)
	case rowTerminalColors:
		return onOff(s.TerminalColors)
	case rowFollowBackground:
		return onOff(s.FollowBackground)
	}
	return ""
}

func rowDescription(row settingsRow) string {
	switch row {
	case rowTheme:
		return "Colors for every page, tab and overlay. Cyberpunk is the custom one."
	case rowPosition:
		return "Top draws the nav bar across the top; Left draws it as a sidebar."
	case rowNavOrder:
		return "A top nav bar sits below or above the header. A sidebar ignores this."
	case rowMode:
		return "Compact: one strip with crumbs. Expanded: a row per level. Path: the route with a drop-down. Columns: one per level."
	case rowAutoHide:
		return "Hide the nav bar after the timeout. Tab or Shift+Tab brings it back."
	case rowAutoHideTimeout:
		return "Seconds before the nav bar hides again, from 1 to 30."
	case rowTabNavTimeout:
		return "Seconds the arrows keep steering tabs after Tab, from 0.5 to 10."
	case rowHeaderHeight:
		return "Lines of header content, from 0 (hidden) to 5. A border adds a row above and below."
	case rowHeaderFrame:
		return "None, a single or double line around the header, or a solid fill with reversed text."
	case rowNavBorder:
		return "Border around the nav bar, across the top or down the side."
	case rowContentBorder:
		return "Border around the page content."
	case rowFooterHeight:
		return "Lines of footer content, from 0 (hidden) to 5. A border adds a row above and below."
	case rowFooterFrame:
		return "None, a single or double line around the footer, or a solid fill with reversed keys."
	case rowKeyGlow:
		return "Pressed keys briefly light up in the footer."
	case rowTerminalColors:
		return "The terminal's own background and text color take the theme's, so the edge around the screen matches."
	case rowFollowBackground:
		return "Switch between the light and the dark theme of a family as the terminal's background says."
	}
	return ""
}

const (
	settingsModalWidth  = 66                     // 50 columns of rows plus frame and padding
	settingsModalHeight = 26                     // title, 16 rows, a blank line, description and hint plus frame
	settingsTextWidth   = settingsModalWidth - 6 // inside the frame and its padding
	settingsDescLines   = 2                      // lines a row description may take
	settingsZIndex      = 50                     // above the Overlays page modals, below help
)

// settingsHint names the popup's keys and the cursor's place among the rows,
// which tells where the list is when a short terminal shows only part of it.
func settingsHint(cursor int) string {
	return fmt.Sprintf("↑/↓ row %d/%d · ←/→ or Enter change · r reset · Esc close", cursor+1, int(rowCount))
}

// rowWindow returns the first row to show and how many, so that a list of
// height lines keeps the cursor row in view.
func rowWindow(cursor, height int) (first, count int) {
	count = min(max(height, 1), int(rowCount))
	first = min(max(cursor-count/2, 0), int(rowCount)-count)
	return first, count
}

// renderSettingsRows draws the rows that fit height lines, the cursor row
// as a pill.
func renderSettingsRows(theme *styles.Theme, snap settingsSnapshot, cursor, height int) string {
	name := lipgloss.NewStyle().Width(settingsName + settingsIndent)
	value := lipgloss.NewStyle().Width(settingsValue)
	indent := strings.Repeat(" ", settingsIndent)

	first, count := rowWindow(cursor, height)
	lines := make([]string, 0, count)
	for row := settingsRow(first); row < settingsRow(first+count); row++ {
		var line string
		if int(row) == cursor {
			line = name.Render(theme.LabelStyle.Render(rowLabel(row))) +
				value.Render(theme.AccentStyle.Bold(true).Render("◂ "+rowValue(snap, row)+" ▸"))
		} else {
			line = name.Render(theme.ContentStyle.Render(indent+rowLabel(row))) +
				value.Render(theme.MutedStyle.Render(indent+rowValue(snap, row)))
		}
		lines = append(lines, line)
	}
	return strings.Join(lines, "\n")
}

// buildSettingsModal creates the centered Settings popup. Its content
// function reads the Model on every render, so the rows stay current while
// settings change; the frame follows the theme through TabApp.
func buildSettingsModal(app *tabapp.TabApp) *overlays.Modal {
	m := app.Model()
	modal := overlays.NewModal().
		ApplyTheme(m.Theme).
		SetTitle("Settings").
		SetExitKey("esc").
		SetBlocking(true).
		SetZIndex(settingsZIndex).
		RegisterKeys(map[string]tea.Msg{
			"up":    SettingsMoveMsg{Delta: -1},
			"down":  SettingsMoveMsg{Delta: 1},
			"left":  SettingsAdjustMsg{Delta: -1},
			"right": SettingsAdjustMsg{Delta: 1},
			"enter": SettingsAdjustMsg{Delta: 1},
			" ":     SettingsAdjustMsg{Delta: 1},
			"r":     SettingsResetMsg{},
			"s":     LeaveSettingsMsg{},
			"q":     tabapp.QuitMsg{},
		})
	modal.Overlay.SetWidth(settingsModalWidth).SetHeight(settingsModalHeight)
	modal.SetContent(func(fb *flexbox.FlexBox) {
		theme := m.Theme
		st := getState(m)
		snap := snapshotOf(m)
		themed(fb, theme)

		// The list takes the lines the description and the hint leave.
		listHeight := fb.GetHeight() - settingsDescLines - 1
		list := textRow(fb, centered(theme.BaseStyle), renderSettingsRows(theme, snap, st.SettingsCursor, listHeight))
		desc := fixedRow(fb, centered(theme.MutedStyle), settingsDescLines, rowDescription(settingsRow(st.SettingsCursor)))
		hint := fixedRow(fb, centered(theme.DisabledStyle), 1, settingsHint(st.SettingsCursor))
		fb.AddRows([]*flexbox.Row{list, desc, hint})
	})
	return modal
}
