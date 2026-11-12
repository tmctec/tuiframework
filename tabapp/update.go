package tabapp

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/styles"
)

// updateModel handles Elm messages and updates the centralized model.
// This is the single place where state transitions happen.
func (a *TabApp) updateModel(msg Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case QuitMsg:
		return a, tea.Quit

	case NoOpMsg:
		// Do nothing - used for keys handled specially elsewhere
		return a, nil

	// Nav bar and chrome settings: the setters own the Model, so every
	// handler computes the new value and delegates.
	case ToggleNavBarPositionMsg:
		next := NavBarLeft
		if a.model.NavBarPosition == NavBarLeft {
			next = NavBarTop
		}
		if next == NavBarLeft && msg.LeftWidth > 0 {
			a.SetNavBarWidth(msg.LeftWidth)
		}
		a.SetNavBarPosition(next)

	case SetNavBarWidthMsg:
		if msg.Width > 0 {
			a.SetNavBarWidth(msg.Width)
		}

	case ToggleNavBarAutoHideMsg:
		a.SetNavBarAutoHide(!a.model.NavBarAutoHide)
		// Switching autohide on arms the timer, so the bar hides without
		// waiting for the next tab-nav key.
		return a, a.startAutohideTimer()

	case SetNavBarTimeoutMsg:
		a.SetNavBarAutoHideTimeout(msg.Duration)

	case AdjustNavBarTimeoutMsg:
		a.SetNavBarAutoHideTimeout(clampDuration(a.model.NavBarTimeout+msg.Delta, time.Second, 30*time.Second))

	case SetTabNavTimeoutMsg:
		a.SetTabNavKeyTimeout(msg.Duration)

	case AdjustTabNavTimeoutMsg:
		a.SetTabNavKeyTimeout(clampDuration(a.model.TabNavTimeout+msg.Delta, 500*time.Millisecond, 10*time.Second))

	case SetNavBarModeMsg:
		a.SetNavBarMode(msg.Mode)

	case ToggleNavBarModeMsg:
		a.SetNavBarMode(a.model.NavBarMode.next())

	case SetHeaderHeightMsg:
		a.SetHeaderHeight(msg.Height)

	case SetFooterHeightMsg:
		a.SetFooterHeight(msg.Height)

	case SetHeaderFrameMsg:
		a.SetHeaderFrame(msg.Frame)

	case SetFooterFrameMsg:
		a.SetFooterFrame(msg.Frame)

	case ToggleHeaderBorderMsg:
		a.SetShowHeaderBorder(a.model.HeaderFrame == FrameNone)

	case ToggleFooterBorderMsg:
		a.SetShowFooterBorder(a.model.FooterFrame == FrameNone)

	case ToggleContentBorderMsg:
		a.SetShowContentBorder(!a.model.ShowContentBorder)

	case ToggleNavBarBorderMsg:
		a.SetShowNavBarBorder(!a.model.ShowNavBarBorder)

	case ToggleNavBarBelowHeaderMsg:
		a.SetNavBarBelowHeader(!a.model.NavBarBelowHeader)

	case ToggleKeyGlowMsg:
		a.SetKeyGlowEnabled(!a.model.KeyGlowEnabled)

	// Section messages move the highlight at the focused level while it is a
	// section level; at the page level they do nothing.
	case NavigateNextSectionMsg:
		if a.tabNavLevel < a.hierarchyLevels()-1 {
			a.moveAtLevel(1)
		}

	case NavigatePrevSectionMsg:
		if a.tabNavLevel < a.hierarchyLevels()-1 {
			a.moveAtLevel(-1)
		}

	case NavigateToPageMsg:
		a.selectPageAtCurrentLevel(msg.Index)

	case NavigateNextPageMsg:
		if idx, count := a.currentPageIndex(); count > 0 {
			a.selectPageAtCurrentLevel((idx + 1) % count)
		}

	case NavigatePrevPageMsg:
		if idx, count := a.currentPageIndex(); count > 0 {
			a.selectPageAtCurrentLevel((idx - 1 + count) % count)
		}

	case NavigateToPathMsg:
		_ = a.NavigateToPath(msg.Path)
		return a, a.startAutohideTimer()

	case DismissOverlayMsg:
		a.HideOverlay()

	// Navigation messages with invalid targets leave the position alone.
	case EnterSectionMsg:
		_ = a.EnterSection(msg.Level, msg.Index)

	case PopToBreadcrumbMsg:
		_ = a.PopToBreadcrumb(msg.BreadcrumbIndex)

	case SetThemeMsg:
		a.UseTheme(msg.Theme)
		return a, nil

	case CycleThemeMsg:
		// The step starts from the active theme's place in the list.
		n := len(msg.Themes)
		if n == 0 {
			return a, nil
		}
		newIdx := 0
		if cur := themeIndex(msg.Themes, a.model.Theme); cur >= 0 {
			newIdx = ((cur+msg.Direction)%n + n) % n
		} else if msg.Direction < 0 {
			newIdx = n - 1
		}
		a.UseTheme(msg.Themes[newIdx])
		return a, nil

	case ShowNotificationMsg:
		n := overlays.NewNotification(msg.Message).SetLevel(msg.Level).SetZIndex(msg.ZIndex)
		return a, a.ShowNotification(n)

	case ShowModalMsg:
		modal := a.textModal(msg.Title, msg.Content).SetZIndex(msg.ZIndex)
		if msg.SizePct > 0 {
			modal.SetSize(msg.SizePct)
		}
		_, _ = a.ShowModal(modal)
		return a, nil

	// A modal that is open closes; syncOverlays then drops it from OpenModals.
	case ToggleModalMsg:
		if id, ok := a.model.OpenModals[msg.Title]; ok {
			a.HideOverlayByID(id)
			return a, nil
		}
		if id, err := a.ShowModal(a.createModalFromMsg(msg)); err == nil {
			a.model.OpenModals[msg.Title] = id
		}
		return a, nil

	case BringOverlayToFrontMsg:
		if id, ok := a.model.OpenModals[msg.Name]; ok {
			a.BringToFront(id)
		}

	case ClearAllOverlaysMsg:
		a.HideAllOverlays()

	// Terminal features: the setters own the Model and View declares them.
	case SetWindowTitleMsg:
		a.SetWindowTitle(msg.Title)

	case ToggleTerminalColorsMsg:
		a.SetTerminalColors(!a.model.TerminalColors)

	case ToggleReportFocusMsg:
		a.SetReportFocus(!a.model.ReportFocus)

	case SetMouseModeMsg:
		a.SetMouseMode(msg.Mode)

	case SetProgressMsg:
		a.SetProgress(msg.State, msg.Value)

	case UseAdaptiveThemeMsg:
		if msg.Light != nil && msg.Dark != nil {
			a.UseAdaptiveTheme(msg.Light, msg.Dark)
		}

	// Reports from the terminal: the Model records them, and the
	// application's handler receives them as well.
	case tea.BackgroundColorMsg:
		a.model.DarkBackground, a.model.BackgroundKnown = msg.IsDark(), true
		a.followBackground()
		return a, a.forward(msg)

	case tea.ColorProfileMsg:
		a.model.ColorProfile = msg.Profile
		return a, a.forward(msg)

	case tea.KeyboardEnhancementsMsg:
		a.model.KeyDisambiguation = msg.SupportsKeyDisambiguation()
		return a, a.forward(msg)

	case tea.FocusMsg:
		a.model.Focused = true
		return a, a.forward(msg)

	case tea.BlurMsg:
		a.model.Focused = !a.model.ReportFocus
		return a, a.forward(msg)

	default:
		// Unknown message type - forward to custom message handler if set
		if a.customMsgHandler != nil {
			if handled, cmd := a.customMsgHandler(msg); handled {
				return a, cmd
			}
		}
	}

	return a, nil
}

// clampDuration limits d to the inclusive range [lo, hi].
func clampDuration(d, lo, hi time.Duration) time.Duration {
	if d < lo {
		return lo
	}
	if d > hi {
		return hi
	}
	return d
}

// themeIndex returns the position of t in themes, matching by pointer first
// and then by name, or -1 when t is absent.
func themeIndex(themes []*styles.Theme, t *styles.Theme) int {
	if t == nil {
		return -1
	}
	for i, th := range themes {
		if th == t {
			return i
		}
	}
	for i, th := range themes {
		if th != nil && th.Name() == t.Name() {
			return i
		}
	}
	return -1
}
