package tabapp

import (
	"regexp"
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// Setters own both the TabApp field and its Model mirror.

func TestSettersSyncModel(t *testing.T) {
	app := New()
	app.SetNavBarPosition(NavBarLeft).
		SetNavBarAutoHideTimeout(7 * time.Second).
		SetTabNavKeyTimeout(3 * time.Second).
		SetShowHeaderBorder(true).
		SetShowNavBarBorder(true).
		SetShowFooterBorder(true).
		SetShowContentBorder(true).
		SetNavBarBelowHeader(false).
		SetKeyGlowEnabled(false)

	m := app.Model()
	if m.NavBarPosition != NavBarLeft {
		t.Errorf("NavBarPosition = %v, want Left", m.NavBarPosition)
	}
	if m.NavBarTimeout != 7*time.Second {
		t.Errorf("NavBarTimeout = %v, want 7s", m.NavBarTimeout)
	}
	if m.TabNavTimeout != 3*time.Second {
		t.Errorf("TabNavTimeout = %v, want 3s", m.TabNavTimeout)
	}
	if m.HeaderFrame != FrameTheme || m.FooterFrame != FrameTheme {
		t.Errorf("frames = %v %v, want the theme's border on both", m.HeaderFrame, m.FooterFrame)
	}
	if !m.ShowNavBarBorder || !m.ShowContentBorder {
		t.Errorf("border flags = %v %v, want both true", m.ShowNavBarBorder, m.ShowContentBorder)
	}
	if m.NavBarBelowHeader {
		t.Error("NavBarBelowHeader should be false")
	}
	if m.KeyGlowEnabled {
		t.Error("KeyGlowEnabled should be false")
	}
	if !app.ShowContentBorder() {
		t.Error("ShowContentBorder() should be true")
	}
}

func TestTogglePositionAfterSetter(t *testing.T) {
	app := New()
	app.SetNavBarPosition(NavBarLeft)

	app.updateModel(ToggleNavBarPositionMsg{})

	if app.NavBarPosition() != NavBarTop || app.Model().NavBarPosition != NavBarTop {
		t.Errorf("after toggle: field=%v model=%v, want Top", app.NavBarPosition(), app.Model().NavBarPosition)
	}
}

func TestAdjustTimeoutStartsFromConfiguredValue(t *testing.T) {
	app := New()
	app.SetNavBarAutoHideTimeout(5 * time.Second)

	app.updateModel(AdjustNavBarTimeoutMsg{Delta: time.Second})

	if app.NavBarAutoHideTimeout() != 6*time.Second {
		t.Errorf("timeout = %v, want 6s", app.NavBarAutoHideTimeout())
	}
	if app.Model().NavBarTimeout != 6*time.Second {
		t.Errorf("model timeout = %v, want 6s", app.Model().NavBarTimeout)
	}
}

func TestAdjustTimeoutsClamp(t *testing.T) {
	app := New()
	app.updateModel(AdjustNavBarTimeoutMsg{Delta: time.Hour})
	if app.Model().NavBarTimeout != 30*time.Second {
		t.Errorf("nav bar timeout = %v, want 30s", app.Model().NavBarTimeout)
	}
	app.updateModel(AdjustTabNavTimeoutMsg{Delta: -time.Hour})
	if app.Model().TabNavTimeout != 500*time.Millisecond {
		t.Errorf("tab nav timeout = %v, want 500ms", app.Model().TabNavTimeout)
	}
}

func TestSetNavBarAutoHideOffShowsBar(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.Init()
	app.Update(navBarAutoHideMsg{id: app.navBarAutoHideID})
	if !app.navBarHidden {
		t.Fatal("precondition: bar should be hidden after the autohide timer fires")
	}

	app.SetNavBarAutoHide(false)

	if app.navBarHidden {
		t.Error("disabling autohide should show the bar")
	}
}

func TestEnablingAutohideArmsTheTimer(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.SetNavBarAutoHide(false)

	_, cmd := app.Update(ToggleNavBarAutoHideMsg{})

	if cmd == nil {
		t.Fatal("enabling autohide should return the timer command")
	}
	app.Update(cmd()) // the timer fires
	if !app.navBarHidden {
		t.Error("the bar should hide when the timer armed by the toggle fires")
	}
}

func TestToggleContentBorderMsg(t *testing.T) {
	app := New()
	app.updateModel(ToggleContentBorderMsg{})
	if !app.ShowContentBorder() || !app.Model().ShowContentBorder {
		t.Error("ToggleContentBorderMsg should enable the content border in field and model")
	}
}

func TestToggleKeyGlowMsg(t *testing.T) {
	app := New()
	app.keyGlowKey = "x"
	app.updateModel(ToggleKeyGlowMsg{})
	if app.KeyGlowEnabled() || app.Model().KeyGlowEnabled {
		t.Error("ToggleKeyGlowMsg should disable key glow in field and model")
	}
	if app.keyGlowKey != "" {
		t.Error("disabling key glow should clear the glowing key")
	}
}

// Update routes every non-key message through updateModel.

func TestUpdateAppliesBuiltinMsgDirectly(t *testing.T) {
	app := New()
	before := app.Model().NavBarAutoHide

	app.Update(ToggleNavBarAutoHideMsg{})

	if app.Model().NavBarAutoHide == before {
		t.Error("a built-in message delivered to Update should be applied")
	}
}

func TestOverlayKeyMsgRoundTrip(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc").
		RegisterKey("y", ToggleNavBarAutoHideMsg{})
	app.MustShowOverlay(overlay)
	before := app.Model().NavBarAutoHide

	_, cmd := app.Update(KeyPress("y"))
	if cmd == nil {
		t.Fatal("overlay key should return a command")
	}
	app.Update(cmd())

	if app.Model().NavBarAutoHide == before {
		t.Error("the message emitted by an overlay key should be applied by Update")
	}
}

func TestUpdateForwardsUnknownMsgToCustomHandler(t *testing.T) {
	app := New()
	var got []tea.Msg
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		got = append(got, msg)
		return true, nil
	})

	app.Update(TestPageMsg{Value: "x"})

	if len(got) != 1 {
		t.Fatalf("custom handler called %d times, want 1", len(got))
	}
	if m, ok := got[0].(TestPageMsg); !ok || m.Value != "x" {
		t.Errorf("custom handler got %#v", got[0])
	}
}

func TestUpdateSetThemeMsg(t *testing.T) {
	app := New()
	app.Update(SetThemeMsg{Theme: styles.Gruvbox})
	if app.Model().Theme != styles.Gruvbox {
		t.Errorf("theme = %s, want Gruvbox", app.Model().Theme.Name())
	}
}

// CycleThemeMsg steps from the active theme.

func TestCycleThemeStartsFromActiveTheme(t *testing.T) {
	app := New()
	app.UseTheme(styles.Gruvbox)

	app.updateModel(CycleThemeMsg{Direction: 1, Themes: styles.AllThemes})

	if app.Model().Theme != styles.GruvboxLight {
		t.Errorf("theme = %s, want GruvboxLight", app.Model().Theme.Name())
	}
}

func TestCycleThemeWraps(t *testing.T) {
	app := New() // Dracula is first in AllThemes
	app.updateModel(CycleThemeMsg{Direction: -1, Themes: styles.AllThemes})
	if app.Model().Theme != styles.AllThemes[len(styles.AllThemes)-1] {
		t.Errorf("theme = %s, want last theme", app.Model().Theme.Name())
	}
}

func TestCycleThemeUnknownStartsAtFirst(t *testing.T) {
	custom := styles.DraculaPalette
	custom.Name = "Custom"
	app := New()
	app.UseTheme(styles.NewTheme(custom))

	app.updateModel(CycleThemeMsg{Direction: 1, Themes: styles.AllThemes})

	if app.Model().Theme != styles.AllThemes[0] {
		t.Errorf("theme = %s, want %s", app.Model().Theme.Name(), styles.AllThemes[0].Name())
	}
	if len(app.Model().Custom) != 0 {
		t.Errorf("Custom = %v, want it left to the application", app.Model().Custom)
	}
}

// Model.OverlayCount mirrors the overlay stack.

func TestModelOverlayCountTracksStack(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	id1 := app.MustShowOverlay(overlays.New().SetContent("1").SetExitKey("esc"))
	app.MustShowOverlay(overlays.New().SetContent("2").SetExitKey("esc"))
	if app.Model().OverlayCount != 2 {
		t.Fatalf("OverlayCount = %d, want 2", app.Model().OverlayCount)
	}
	app.HideOverlayByID(id1)
	if app.Model().OverlayCount != 1 {
		t.Errorf("OverlayCount = %d, want 1", app.Model().OverlayCount)
	}
	app.HideAllOverlays()
	if app.Model().OverlayCount != 0 {
		t.Errorf("OverlayCount = %d, want 0", app.Model().OverlayCount)
	}
}

func TestOpenModalsPrunedWhenModalClosesByExitKey(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.updateModel(ToggleModalMsg{Title: "Modal A", Content: "a", SizePct: 50})
	if _, open := app.Model().OpenModals["Modal A"]; !open {
		t.Fatal("precondition: Modal A should be registered")
	}

	app.Update(KeyPress("esc")) // the modal's exit key

	if app.HasOverlay() {
		t.Fatal("Esc should close the modal")
	}
	if _, open := app.Model().OpenModals["Modal A"]; open {
		t.Error("OpenModals should forget a modal closed by its exit key")
	}
	// Toggling again opens it instead of trying to close a vanished overlay.
	app.updateModel(ToggleModalMsg{Title: "Modal A", Content: "a", SizePct: 50})
	if !app.HasOverlay() {
		t.Error("ToggleModalMsg after Esc should open the modal again")
	}
}

// Modals created by TabApp take the active theme and draw their title.

func TestToggleModalMsgUsesCurrentTheme(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.UseTheme(styles.Gruvbox)

	app.updateModel(ToggleModalMsg{Title: "X", Content: "c", SizePct: 50})

	top := app.TopOverlay()
	if top == nil {
		t.Fatal("modal should be shown")
	}
	if top.Style().GetBorderTopForeground() != styles.Gruvbox.ModalBorderColor {
		t.Errorf("border = %v, want %v", top.Style().GetBorderTopForeground(), styles.Gruvbox.ModalBorderColor)
	}
	if top.Style().GetBackground() != styles.Gruvbox.ModalBackground {
		t.Errorf("background = %v, want %v", top.Style().GetBackground(), styles.Gruvbox.ModalBackground)
	}
}

func TestToggleModalMsgBorderColorOverride(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	app.updateModel(ToggleModalMsg{Title: "X", SizePct: 50, BorderColor: lipgloss.Color("#123456")})

	if got := app.TopOverlay().Style().GetBorderTopForeground(); got != lipgloss.Color("#123456") {
		t.Errorf("border = %v, want #123456", got)
	}
}

var ansiSequence = regexp.MustCompile("\x1b\\[[0-9;?]*[A-Za-z]")

func TestShowModalRendersTitle(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})

	modal := overlays.NewModal().SetTitle("Hello Title").SetExitKey("esc")
	if _, err := app.ShowModal(modal); err != nil {
		t.Fatal(err)
	}

	plain := ansiSequence.ReplaceAllString(app.Render(), "")
	if !strings.Contains(plain, "Hello Title") {
		t.Error("View() should contain the modal title")
	}
}
