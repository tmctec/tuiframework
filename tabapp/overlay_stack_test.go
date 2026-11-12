package tabapp

import (
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// overlayApp is a sized single-page app with the page key "a" toggling autohide.
func overlayApp() *TabApp {
	app := New()
	app.MustAddPage(page.New("Test").BindKey("a", "Toggle autohide", ToggleNavBarAutoHideMsg{}))
	return sizedApp(app)
}

// stackIDs returns the overlay ids from the bottom of the stack to the top.
func stackIDs(app *TabApp) []int64 {
	ids := make([]int64, 0, len(app.overlayStack))
	for _, e := range app.overlayStack {
		ids = append(ids, e.id)
	}
	return ids
}

func TestOverlayStackIDsIncrease(t *testing.T) {
	app := overlayApp()
	var last int64
	for i := 0; i < 3; i++ {
		id, err := app.ShowOverlay(overlays.New().SetExitKey("esc"))
		if err != nil {
			t.Fatal(err)
		}
		if id <= last {
			t.Errorf("id %d after %d, want increasing", id, last)
		}
		last = id
	}
	if app.OverlayCount() != 3 || app.Model().OverlayCount != 3 {
		t.Errorf("OverlayCount = %d, model %d, want 3", app.OverlayCount(), app.Model().OverlayCount)
	}
}

func TestOverlayStackRejectsInvalidOverlays(t *testing.T) {
	app := overlayApp()

	if _, err := app.ShowOverlay(overlays.New()); !errors.Is(err, overlays.ErrNoExitKey) {
		t.Errorf("no exit key: err = %v, want ErrNoExitKey", err)
	}
	if _, err := app.ShowOverlay(overlays.New().SetExitKey("esc").RegisterKey("ctrl+c", NoOpMsg{})); err == nil || !strings.Contains(err.Error(), "ctrl+c") {
		t.Errorf("hardcoded overlay key: err = %v, want a conflict naming ctrl+c", err)
	}
	if _, err := app.ShowOverlay(overlays.New().SetExitKey("ctrl+q")); err == nil || !strings.Contains(err.Error(), "ctrl+q") {
		t.Errorf("hardcoded exit key: err = %v, want a conflict naming ctrl+q", err)
	}
	if app.HasOverlay() || app.Model().OverlayCount != 0 {
		t.Error("a rejected overlay must not reach the stack")
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustShowOverlay should panic on a rejected overlay")
		}
	}()
	app.MustShowOverlay(overlays.New())
}

func TestOverlayStackKeepsZOrderThenInsertionOrder(t *testing.T) {
	app := overlayApp()
	for _, z := range []int{5, 1, 5, 3} {
		app.MustShowOverlay(overlays.New().SetExitKey("esc").SetZIndex(z))
	}
	wantIDs := []int64{2, 4, 1, 3} // z1, z3, then the two z5 in insertion order
	wantZ := []int{1, 3, 5, 5}
	for i, e := range app.overlayStack {
		if e.id != wantIDs[i] || e.overlay.ZIndex() != wantZ[i] {
			t.Errorf("stack[%d] = id %d z %d, want id %d z %d", i, e.id, e.overlay.ZIndex(), wantIDs[i], wantZ[i])
		}
	}
}

func TestOverlayStackNotificationTimeoutCommand(t *testing.T) {
	app := overlayApp()

	if cmd := app.ShowNotification(overlays.Info("stays").SetTimeout(0)); cmd != nil {
		t.Error("a notification without a timeout returns no command")
	}
	if cmd := app.ShowNotification(overlays.Info("invalid").SetExitKey("")); cmd != nil || app.OverlayCount() != 1 {
		t.Errorf("an invalid notification returns no command and is not shown: cmd %v count %d", cmd, app.OverlayCount())
	}

	timed := overlays.Info("goes").SetTimeout(time.Millisecond)
	cmd := app.ShowNotification(timed)
	if cmd == nil {
		t.Fatal("a notification with a timeout returns a command")
	}
	msg, ok := cmd().(overlayTimeoutMsg)
	if !ok {
		t.Fatalf("command produced %T, want overlayTimeoutMsg", msg)
	}
	app.Update(msg)
	if app.OverlayCount() != 1 || app.TopOverlay() == timed.Overlay {
		t.Errorf("the timeout should hide the timed notification only: count %d", app.OverlayCount())
	}
	if app.TopOverlay().Content() != overlays.NotificationConfig.LevelStyle(overlays.LevelInfo).Icon+" stays" {
		t.Errorf("remaining overlay = %q, want the untimed notification", app.TopOverlay().Content())
	}
}

func TestOverlayStackShowModalRecordsAndRestyles(t *testing.T) {
	app := overlayApp()
	themed := overlays.NewModal().SetTitle("Themed").SetExitKey("esc")
	styled := overlays.NewModal().SetTitle("Styled").SetExitKey("esc").
		SetStyle(lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#123456")))

	if _, err := app.ShowModal(overlays.NewModal()); !errors.Is(err, overlays.ErrNoExitKey) {
		t.Errorf("modal without an exit key: err = %v, want ErrNoExitKey", err)
	}
	idThemed, err := app.ShowModal(themed)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := app.ShowModal(styled); err != nil {
		t.Fatal(err)
	}
	for _, e := range app.overlayStack {
		switch e.id {
		case idThemed:
			if e.modal != themed {
				t.Error("the themed modal should be recorded on its entry")
			}
		default:
			if e.modal != styled {
				t.Error("the styled modal should be recorded on its entry")
			}
		}
	}

	app.UseTheme(styles.Gruvbox)
	if got := themed.Style().GetBorderTopForeground(); got != styles.Gruvbox.ModalBorderColor {
		t.Errorf("themed modal border = %v, want Gruvbox %v", got, styles.Gruvbox.ModalBorderColor)
	}
	if got := styled.Style().GetBorderTopForeground(); got != lipgloss.Color("#123456") {
		t.Errorf("styled modal border = %v, want its own #123456", got)
	}
}

func TestOverlayStackHideOperations(t *testing.T) {
	app := overlayApp()
	if app.TopOverlay() != nil || app.ActiveOverlay() != nil {
		t.Error("an empty stack has no top overlay")
	}
	app.HideOverlay() // nothing to hide

	low := overlays.New().SetExitKey("1").SetZIndex(1)
	high := overlays.New().SetExitKey("2").SetZIndex(9)
	idLow := app.MustShowOverlay(low)
	app.MustShowOverlay(high)

	app.HideOverlay()
	if app.TopOverlay() != low || high.IsVisible() || app.Model().OverlayCount != 1 {
		t.Errorf("HideOverlay should remove the top: top %v visible %v model %d", app.TopOverlay() == low, high.IsVisible(), app.Model().OverlayCount)
	}
	if app.HideOverlayByID(999) {
		t.Error("HideOverlayByID should report false for an unknown id")
	}
	if !app.HideOverlayByID(idLow) || low.IsVisible() || app.HasOverlay() {
		t.Error("HideOverlayByID should remove the overlay and mark it hidden")
	}

	app.MustShowOverlay(overlays.New().SetExitKey("1"))
	app.MustShowOverlay(overlays.New().SetExitKey("2"))
	app.HideAllOverlays()
	if app.HasOverlay() || app.OverlayCount() != 0 || app.Model().OverlayCount != 0 || app.TopOverlay() != nil {
		t.Error("HideAllOverlays should clear the stack and the model count")
	}
}

type stackKeyMsg struct{ From string }

// TestOverlayStackKeyRouting drives keys through a non-blocking overlay on
// top of a blocking one and then the page.
func TestOverlayStackKeyRouting(t *testing.T) {
	app := overlayApp()
	lower := overlays.New().SetExitKey("1").SetZIndex(1).SetBlocking(true).RegisterKey("y", stackKeyMsg{From: "lower"})
	upper := overlays.New().SetExitKey("2").SetZIndex(5).SetBlocking(false)
	app.MustShowOverlay(lower)
	app.MustShowOverlay(upper)

	// The top overlay is non-blocking, so the lower overlay's key is reached.
	handled, cmd := app.handleOverlayKeys("y")
	if !handled || cmd == nil {
		t.Fatal("a registered key on a lower overlay should be handled")
	}
	if msg, ok := cmd().(stackKeyMsg); !ok || msg.From != "lower" {
		t.Errorf("command produced %#v, want the lower overlay's message", cmd())
	}

	// The lower overlay blocks: unknown keys stop there and never reach the page.
	autohide := app.Model().NavBarAutoHide
	app.Update(KeyPress("a"))
	if handled, cmd := app.handleOverlayKeys("zzz"); !handled || cmd != nil {
		t.Error("a blocking overlay swallows unknown keys")
	}
	if app.Model().NavBarAutoHide != autohide {
		t.Error("the page key must not fire under a blocking overlay")
	}

	// The lower overlay's exit key reaches it through the non-blocking top.
	app.Update(KeyPress("1"))
	if app.OverlayCount() != 1 || app.TopOverlay() != upper || lower.IsVisible() {
		t.Errorf("exit key 1 should hide the lower overlay only: count %d", app.OverlayCount())
	}

	// With only the non-blocking overlay left, unknown keys reach the page.
	if handled, _ := app.handleOverlayKeys("a"); handled {
		t.Error("a non-blocking overlay passes unknown keys on")
	}
	app.Update(KeyPress("a"))
	if app.Model().NavBarAutoHide == autohide {
		t.Error("the page key should fire under a non-blocking overlay")
	}

	app.Update(KeyPress("2"))
	if app.HasOverlay() {
		t.Error("exit key 2 should hide the last overlay")
	}
}

func TestOverlayStackBringToFront(t *testing.T) {
	app := overlayApp()
	a := overlays.New().SetExitKey("1").SetZIndex(1)
	b := overlays.New().SetExitKey("2").SetZIndex(4)
	c := overlays.New().SetExitKey("3").SetZIndex(9)
	idA := app.MustShowOverlay(a)
	app.MustShowOverlay(b)
	app.MustShowOverlay(c)
	before := stackIDs(app)

	app.BringToFront(999)
	if got := stackIDs(app); !slices.Equal(got, before) || c.ZIndex() != 9 {
		t.Errorf("unknown id changed the stack: %v, want %v", got, before)
	}

	app.BringToFront(idA)
	if app.TopOverlay() != a || a.ZIndex() != 10 {
		t.Errorf("BringToFront: top is a %v, z %d, want z 10", app.TopOverlay() == a, a.ZIndex())
	}
	if got := stackIDs(app); !slices.Equal(got, []int64{2, 3, idA}) {
		t.Errorf("stack order = %v, want b, c, a", got)
	}
}

func TestOverlayStackSyncDropsStaleNamedIDs(t *testing.T) {
	app := overlayApp()
	live := app.MustShowOverlay(overlays.New().SetExitKey("1"))
	ids := app.Model().OpenModals
	ids["Live"], ids["Stale"] = live, 999

	app.MustShowOverlay(overlays.New().SetExitKey("2")) // syncs the registry
	if _, ok := ids["Stale"]; ok {
		t.Error("a name whose overlay is gone should be dropped")
	}
	if ids["Live"] != live {
		t.Error("a name whose overlay is shown should stay")
	}

	app.HideOverlayByID(live)
	if len(ids) != 0 {
		t.Errorf("registry = %v, want empty after hiding the live overlay", ids)
	}

	// The library keeps its state out of the application's Custom map.
	if len(app.Model().Custom) != 0 {
		t.Errorf("Custom = %v, want it left to the application", app.Model().Custom)
	}
}

// TestShowOverlayKeyConflictWrapsErrKeyConflict checks that a caller can
// tell a key conflict from other failures with errors.Is.
func TestShowOverlayKeyConflictWrapsErrKeyConflict(t *testing.T) {
	app := overlayApp()

	_, err := app.ShowOverlay(overlays.New().SetExitKey("esc").RegisterKey("ctrl+c", NoOpMsg{}))
	if !errors.Is(err, overlays.ErrKeyConflict) || !strings.Contains(err.Error(), `"ctrl+c"`) {
		t.Errorf("registered quit key: err = %v, want ErrKeyConflict naming ctrl+c", err)
	}

	_, err = app.ShowOverlay(overlays.New().SetExitKey("ctrl+q"))
	if !errors.Is(err, overlays.ErrKeyConflict) || !strings.Contains(err.Error(), `"ctrl+q"`) {
		t.Errorf("quit key as exit key: err = %v, want ErrKeyConflict naming ctrl+q", err)
	}

	_, err = app.ShowOverlay(overlays.New())
	if !errors.Is(err, overlays.ErrNoExitKey) || errors.Is(err, overlays.ErrKeyConflict) {
		t.Errorf("no exit key: err = %v, want ErrNoExitKey", err)
	}
	if app.HasOverlay() {
		t.Error("a rejected overlay is not shown")
	}
}
