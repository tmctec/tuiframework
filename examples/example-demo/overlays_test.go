package main

import (
	"strings"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/tabapp"
)

// overlaysApp is the demo on its Overlays page at 100x40.
func overlaysApp(t *testing.T) *tabapp.TabApp {
	t.Helper()
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if err := app.NavigateToPath([]int{0, 3, 0, 0}); err != nil { // Showcase > Overlays
		t.Fatal(err)
	}
	return app
}

func rune1(r rune) tea.KeyMsg {
	return tabapp.KeyPress(string(r))
}

// topContent returns the text of the topmost overlay, "" when there is none.
func topContent(app *tabapp.TabApp) string {
	if top := app.TopOverlay(); top != nil {
		return top.Content()
	}
	return ""
}

// TestCycleNotificationWalksTheLevels presses n on the Themes page nine
// times: one sample notification per level, then the first level again.
func TestCycleNotificationWalksTheLevels(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	if err := app.NavigateToPath([]int{0, 1, 0, 0}); err != nil { // Showcase > Themes
		t.Fatal(err)
	}

	for i := 0; i <= len(sampleNotifications); i++ {
		app.HideAllOverlays()
		press(app, rune1('n'))
		level := overlays.NotificationLevel(i % len(sampleNotifications))
		icon := overlays.NotificationConfig.LevelStyle(level).Icon
		if got, want := topContent(app), icon+" "+sampleNotifications[level]; got != want {
			t.Fatalf("press %d: notification %q, want %q", i+1, got, want)
		}
	}
}

// TestStackKeyTogglesThreeModals presses z twice on the Overlays page.
func TestStackKeyTogglesThreeModals(t *testing.T) {
	app := overlaysApp(t)

	press(app, rune1('z'))
	if app.Model().OverlayCount != 3 {
		t.Fatalf("after z: %d overlays, want 3", app.Model().OverlayCount)
	}
	for _, title := range []string{"Stack 1", "Stack 2", "Stack 3"} {
		if _, open := app.Model().OpenModals[title]; !open {
			t.Errorf("%q should be open", title)
		}
	}
	if top := app.TopOverlay(); top.Title() != "Stack 3" || top.ZIndex() != 3 {
		t.Errorf("top = %q at z %d, want Stack 3 at z 3", top.Title(), top.ZIndex())
	}

	press(app, rune1('z'))
	if app.Model().OverlayCount != 0 || len(app.Model().OpenModals) != 0 {
		t.Errorf("after a second z: %d overlays, registry %v, want none", app.Model().OverlayCount, app.Model().OpenModals)
	}
}

// TestFrontKeyBringsModalAForward covers f: a note while Modal A is closed,
// and Modal A above Modal B with a tip once it is open.
func TestFrontKeyBringsModalAForward(t *testing.T) {
	app := overlaysApp(t)

	press(app, rune1('f'))
	if got := topContent(app); !strings.HasSuffix(got, "Modal A not open") {
		t.Errorf("f with Modal A closed: %q, want the not-open note", got)
	}

	app.HideAllOverlays()
	press(app, rune1('a')) // z 5
	press(app, rune1('b')) // z 15
	press(app, rune1('f'))

	// Walk the stack from the top down.
	var aZ, bZ int
	var tip bool
	for app.HasOverlay() {
		top := app.TopOverlay()
		switch {
		case top.Title() == "Modal A":
			aZ = top.ZIndex()
		case top.Title() == "Modal B":
			bZ = top.ZIndex()
		case strings.HasSuffix(top.Content(), "Modal A brought to front"):
			tip = true
		}
		app.HideOverlay()
	}
	if aZ <= bZ {
		t.Errorf("Modal A z %d, Modal B z %d: want A above B", aZ, bZ)
	}
	if !tip {
		t.Error("f should say Modal A was brought to the front")
	}
}

// TestClearKeyLeavesOnlyItsNotification covers x: every overlay goes, and
// the notification that says so is shown after the clear, not taken by it.
func TestClearKeyLeavesOnlyItsNotification(t *testing.T) {
	app := overlaysApp(t)
	press(app, rune1('a'))
	press(app, rune1('i'))
	if app.Model().OverlayCount != 2 {
		t.Fatalf("precondition: %d overlays, want a modal and a notification", app.Model().OverlayCount)
	}

	press(app, rune1('x'))
	if app.Model().OverlayCount != 1 || !strings.HasSuffix(topContent(app), "All overlays cleared") {
		t.Errorf("after x: %d overlays, top %q, want only the cleared notification", app.Model().OverlayCount, topContent(app))
	}
	if len(app.Model().OpenModals) != 0 {
		t.Errorf("OpenModals = %v, want none", app.Model().OpenModals)
	}
}

// TestNotifyKeysStackByZIndex covers i, w and e: three levels at z 10, 20
// and 30.
func TestNotifyKeysStackByZIndex(t *testing.T) {
	app := overlaysApp(t)
	for _, k := range []rune{'e', 'i', 'w'} {
		press(app, rune1(k))
	}
	want := []struct {
		level overlays.NotificationLevel
		z     int
	}{{overlays.LevelError, 30}, {overlays.LevelWarning, 20}, {overlays.LevelInfo, 10}}
	for _, w := range want {
		top := app.TopOverlay()
		icon := overlays.NotificationConfig.LevelStyle(w.level).Icon
		if top == nil || top.ZIndex() != w.z || !strings.HasPrefix(top.Content(), icon) {
			t.Fatalf("top overlay = %+v, want the %v notification at z %d", top, w.level, w.z)
		}
		app.HideOverlay()
	}
}

// TestHelpKeyShowsTheHelpModal covers ? on a page: a blocking modal titled
// Help above every other overlay.
func TestHelpKeyShowsTheHelpModal(t *testing.T) {
	app := overlaysApp(t)
	press(app, rune1('c')) // Modal C at z 25
	press(app, rune1('?'))

	top := app.TopOverlay()
	if top == nil || top.Title() != "Help" || top.ZIndex() != 100 || top.Size() != 55 || !top.IsBlocking() {
		t.Fatalf("top overlay = %+v, want the blocking Help modal at z 100 and 55%%", top)
	}
	if !strings.Contains(top.Content(), "bring Modal A to the front") {
		t.Errorf("help content = %q, want the page's keys", top.Content())
	}
	press(app, tabapp.KeyPress("esc"))
	if app.TopOverlay().Title() != "Modal C" {
		t.Error("Esc should close the help modal and leave Modal C")
	}
}

// TestEmitInOrderKeepsTheOrder checks the sequence helper delivers its
// messages one after the other.
func TestEmitInOrderKeepsTheOrder(t *testing.T) {
	app := buildApp()
	app.Update(tea.WindowSizeMsg{Width: 100, Height: 40})
	drain(app, emitInOrder(
		tabapp.SetNavBarTimeoutMsg{Duration: 5 * time.Second},
		tabapp.SetNavBarTimeoutMsg{Duration: 9 * time.Second},
	))
	if got := app.Model().NavBarTimeout; got != 9*time.Second {
		t.Errorf("NavBarTimeout = %v, want the last message's 9s", got)
	}
}
