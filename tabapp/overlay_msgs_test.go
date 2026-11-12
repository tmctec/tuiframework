package tabapp

import (
	"slices"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
)

// msgApp is a sized single-page app for driving overlay messages.
func msgApp() *TabApp {
	return sizedApp(New().MustAddPage(page.New("Test")))
}

// overlayLevel reads a notification's level back from the icon that
// prefixes its content, or returns -1 for an overlay that is not one.
func overlayLevel(o *overlays.Overlay) overlays.NotificationLevel {
	for l := overlays.LevelInfo; l <= overlays.LevelError; l++ {
		if strings.HasPrefix(o.Content(), overlays.NotificationConfig.LevelStyle(l).Icon+" ") {
			return l
		}
	}
	return -1
}

// topLevel returns the level of the top overlay, failing when it is not a notification.
func topLevel(t *testing.T, app *TabApp) overlays.NotificationLevel {
	t.Helper()
	top := app.TopOverlay()
	if top == nil {
		t.Fatal("no overlay shown")
	}
	l := overlayLevel(top)
	if l < 0 {
		t.Fatalf("top overlay %q is not a notification", top.Content())
	}
	return l
}

// namedIDs returns the registry of the modals ToggleModalMsg has open.
func namedIDs(t *testing.T, app *TabApp) map[string]int64 {
	t.Helper()
	if app.Model().OpenModals == nil {
		t.Fatal("Model.OpenModals is nil, want a map")
	}
	return app.Model().OpenModals
}

// TestOverlayMsgShowNotification covers ShowNotificationMsg at every level:
// the level styles the notification, the z-index stacks it, and the command
// dismisses it.
func TestOverlayMsgShowNotification(t *testing.T) {
	for level := overlays.LevelInfo; level <= overlays.LevelError; level++ {
		app := msgApp()
		_, cmd := app.Update(ShowNotificationMsg{Level: level, Message: "done", ZIndex: 7})
		if cmd == nil {
			t.Errorf("%v: no timeout command; notifications dismiss themselves", level)
		}
		if got := topLevel(t, app); got != level {
			t.Errorf("level %v, want %v", got, level)
		}
		top := app.TopOverlay()
		if !strings.HasSuffix(top.Content(), " done") || top.ZIndex() != 7 || app.OverlayCount() != 1 {
			t.Errorf("%v: content %q z %d count %d", level, top.Content(), top.ZIndex(), app.OverlayCount())
		}
	}

	// The zero message is an info notification at z-index 0.
	app := msgApp()
	app.Update(ShowNotificationMsg{Message: "plain"})
	if topLevel(t, app) != overlays.LevelInfo || app.TopOverlay().ZIndex() != 0 {
		t.Errorf("zero message: level %v z %d, want Info at 0", topLevel(t, app), app.TopOverlay().ZIndex())
	}
	if len(app.Model().Custom) != 0 {
		t.Errorf("Custom = %v, want it left to the application", app.Model().Custom)
	}
}

func TestOverlayMsgToggleModalOpensThenCloses(t *testing.T) {
	app := msgApp()
	msg := ToggleModalMsg{Title: "Settings", Content: "body text", ZIndex: 3, SizePct: 40}

	app.Update(msg)
	ids := namedIDs(t, app)
	id, ok := ids["Settings"]
	if !ok || !app.hasOverlayID(id) {
		t.Fatalf("OpenModals = %v, want Settings mapped to a shown overlay", ids)
	}
	top := app.TopOverlay()
	if top.Title() != "Settings" || top.Content() != "body text" || top.ZIndex() != 3 || top.Size() != 40 {
		t.Errorf("modal title %q content %q z %d size %d", top.Title(), top.Content(), top.ZIndex(), top.Size())
	}
	if top.ExitKey() != "esc" || top.IsBlocking() {
		t.Errorf("modal exit %q blocking %v, want esc and non-blocking", top.ExitKey(), top.IsBlocking())
	}
	if e := app.overlayStack[0]; e.modal == nil || !e.modal.FollowsTheme() {
		t.Error("a toggled modal follows the theme")
	}

	app.Update(msg)
	if app.HasOverlay() || app.Model().OverlayCount != 0 {
		t.Error("toggling again should close the modal")
	}
	if _, ok := namedIDs(t, app)["Settings"]; ok {
		t.Error("the closed modal should leave the registry")
	}
}

func TestOverlayMsgBringOverlayToFront(t *testing.T) {
	app := msgApp()
	app.Update(ToggleModalMsg{Title: "A", ZIndex: 1})
	app.Update(ToggleModalMsg{Title: "B", ZIndex: 5})
	ids := namedIDs(t, app)

	_, cmd := app.Update(BringOverlayToFrontMsg{Name: "A"})
	if cmd != nil || app.OverlayCount() != 2 {
		t.Errorf("bringing a modal forward shows nothing new: cmd nil = %v, count %d", cmd == nil, app.OverlayCount())
	}
	if top := app.TopOverlay(); top.Title() != "A" || top.ZIndex() != 6 {
		t.Errorf("top = %q z %d, want A at z 6", top.Title(), top.ZIndex())
	}
	if ids["A"] != app.overlayStack[1].id {
		t.Error("the registry should still name A")
	}

	// A title that is not open changes nothing.
	before := stackIDs(app)
	if _, cmd = app.Update(BringOverlayToFrontMsg{Name: "Nope"}); cmd != nil {
		t.Error("an unknown name should return no command")
	}
	if got := stackIDs(app); !slices.Equal(got, before) {
		t.Errorf("an unknown name changed the stack: %v, want %v", got, before)
	}
}

func TestOverlayMsgClearAllOverlays(t *testing.T) {
	app := msgApp()
	app.Update(ToggleModalMsg{Title: "A"})
	app.Update(ShowNotificationMsg{Level: overlays.LevelError, Message: "boom"})

	_, cmd := app.Update(ClearAllOverlaysMsg{})
	if cmd != nil || app.HasOverlay() {
		t.Errorf("clearing hides everything and shows nothing: cmd nil = %v, count %d", cmd == nil, app.OverlayCount())
	}
	if ids := namedIDs(t, app); len(ids) != 0 {
		t.Errorf("registry = %v, want empty", ids)
	}
}

// TestOverlayMsgShowModalSizeAndZIndex covers the optional fields of
// ShowModalMsg: a z-index, and a size that replaces the overlay default.
func TestOverlayMsgShowModalSizeAndZIndex(t *testing.T) {
	app := msgApp()
	app.Update(ShowModalMsg{Title: "Help", Content: "help text", ZIndex: 9, SizePct: 55})
	top := app.TopOverlay()
	if top == nil {
		t.Fatal("no overlay shown")
	}
	if top.Title() != "Help" || top.Size() != 55 || top.ExitKey() != "esc" || top.ZIndex() != 9 || top.Content() != "help text" {
		t.Errorf("modal title %q size %d exit %q z %d content %q", top.Title(), top.Size(), top.ExitKey(), top.ZIndex(), top.Content())
	}
	if !top.IsBlocking() {
		t.Error("the modal blocks page keys")
	}
	if e := app.overlayStack[0]; e.modal == nil || !e.modal.FollowsTheme() {
		t.Error("the modal follows the theme")
	}

	app.HideAllOverlays()
	app.Update(ShowModalMsg{Title: "Plain"})
	if got, want := app.TopOverlay().Size(), overlays.New().Size(); got != want {
		t.Errorf("size without SizePct = %d, want the overlay default %d", got, want)
	}
}

func TestOverlayMsgDismissOverlay(t *testing.T) {
	app := msgApp()
	app.Update(ToggleModalMsg{Title: "A", ZIndex: 1})
	app.Update(ToggleModalMsg{Title: "B", ZIndex: 2})

	app.Update(DismissOverlayMsg{})
	if app.OverlayCount() != 1 || app.TopOverlay().Title() != "A" {
		t.Errorf("DismissOverlayMsg should remove the top: count %d top %q", app.OverlayCount(), app.TopOverlay().Title())
	}
	if _, ok := namedIDs(t, app)["B"]; ok {
		t.Error("the dismissed modal should leave the registry")
	}
	app.Update(DismissOverlayMsg{})
	app.Update(DismissOverlayMsg{}) // nothing left to dismiss
	if app.HasOverlay() {
		t.Error("the stack should be empty")
	}
}

func TestOverlayMsgCreateModalBorderColor(t *testing.T) {
	app := msgApp()
	app.UseTheme(styles.Gruvbox)

	themed := app.createModalFromMsg(ToggleModalMsg{Title: "T", SizePct: 50})
	if got := themed.Style().GetBorderTopForeground(); got != styles.Gruvbox.ModalBorderColor {
		t.Errorf("border = %v, want the Gruvbox modal border %v", got, styles.Gruvbox.ModalBorderColor)
	}
	if themed.Content() != "" {
		t.Errorf("content = %q, want none without Content", themed.Content())
	}

	custom := app.createModalFromMsg(ToggleModalMsg{Title: "C", Content: "x", BorderColor: lipgloss.Color("#ff0000")})
	if got := custom.Style().GetBorderTopForeground(); got != lipgloss.Color("#ff0000") {
		t.Errorf("border = %v, want the override #ff0000", got)
	}
	if custom.Style().GetBackground() != styles.Gruvbox.ModalBackground {
		t.Error("the override keeps the theme's modal background")
	}
}

func TestOverlayMsgShowModal(t *testing.T) {
	app := msgApp()
	app.UseTheme(styles.Nord)

	app.Update(ShowModalMsg{Title: "About", Content: "version 1"})

	top := app.TopOverlay()
	if top == nil {
		t.Fatal("ShowModalMsg should show a modal")
	}
	if top.Title() != "About" || top.Content() != "version 1" || top.ExitKey() != "esc" || !top.IsBlocking() {
		t.Errorf("modal title %q content %q exit %q blocking %v", top.Title(), top.Content(), top.ExitKey(), top.IsBlocking())
	}
	if got := top.Style().GetBorderTopForeground(); got != styles.Nord.ModalBorderColor {
		t.Errorf("border = %v, want the active theme's %v", got, styles.Nord.ModalBorderColor)
	}
	if len(app.Model().OpenModals) != 0 {
		t.Errorf("OpenModals = %v, want only ToggleModalMsg modals listed", app.Model().OpenModals)
	}
}
