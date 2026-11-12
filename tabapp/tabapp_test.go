package tabapp

import (
	"strings"
	"testing"
	"time"

	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
)

func TestNew(t *testing.T) {
	app := New()

	// Check hardcoded keys are registered
	if !app.IsHardcoded("ctrl+c") {
		t.Error("ctrl+c should be hardcoded")
	}
	if !app.IsHardcoded("ctrl+q") {
		t.Error("ctrl+q should be hardcoded")
	}

	// Check default universal keys are registered
	if !app.IsUniversalKey("q") {
		t.Error("q should be universal by default")
	}
	if !app.IsUniversalKey("tab") {
		t.Error("tab should be universal by default")
	}
	if !app.IsUniversalKey("shift+tab") {
		t.Error("shift+tab should be universal by default")
	}
}

// TestPageMsg is a test message type for page key bindings
type TestPageMsg struct {
	Value string
}

func TestAddPageNoConflict(t *testing.T) {
	app := New()
	p := page.New("Test").
		BindKey("a", "Action A", TestPageMsg{Value: "a"}).
		BindKey("b", "Action B", TestPageMsg{Value: "b"})

	err := app.AddPage(p)
	if err != nil {
		t.Errorf("should not error: %v", err)
	}
	if len(app.Pages()) != 1 {
		t.Errorf("expected 1 page, got %d", len(app.Pages()))
	}
}

func TestAddPageUniversalConflict(t *testing.T) {
	app := New()
	p := page.New("Test").
		BindKey("q", "Quit", TestPageMsg{Value: "q"}) // 'q' is universal by default

	err := app.AddPage(p)
	if err == nil {
		t.Error("should error on universal key conflict")
	}
	if !strings.Contains(err.Error(), "conflicts with universal") {
		t.Errorf("error should mention universal conflict: %v", err)
	}
}

func TestAddPageHardcodedConflict(t *testing.T) {
	app := New()
	p := page.New("Test").
		BindKey("ctrl+c", "Force quit", TestPageMsg{Value: "ctrl+c"})

	err := app.AddPage(p)
	if err == nil {
		t.Error("should error on hardcoded key conflict")
	}
	if !strings.Contains(err.Error(), "hardcoded") {
		t.Errorf("error should mention hardcoded: %v", err)
	}
}

func TestMustAddPagePanic(t *testing.T) {
	app := New()
	p := page.New("Test").
		BindKey("q", "Quit", TestPageMsg{Value: "q"})

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustAddPage should panic on conflict")
		}
	}()

	app.MustAddPage(p)
}

func TestMustAddPageSuccess(t *testing.T) {
	app := New()
	p := page.New("Test").
		BindKey("a", "Action", TestPageMsg{Value: "a"})

	// Should not panic
	result := app.MustAddPage(p)
	if result != app {
		t.Error("MustAddPage should return self for chaining")
	}
}

func TestUnbindKeyAllowsPageKey(t *testing.T) {
	app := New()

	// Remove 'q' from universal
	app.UnbindKey("q")

	if app.IsUniversalKey("q") {
		t.Error("q should no longer be universal")
	}

	// Now page can use 'q'
	p := page.New("Test").
		BindKey("q", "Quit", TestPageMsg{Value: "q"})

	err := app.AddPage(p)
	if err != nil {
		t.Errorf("should allow 'q' after unbinding: %v", err)
	}
}

// TestBindKeyMsg is a message for testing universal key bindings
type TestBindKeyMsg struct {
	Triggered bool
}

func TestBindKey(t *testing.T) {
	app := New()

	// Use message-based binding (Elm pattern)
	err := app.BindKey("?", "Help", TestBindKeyMsg{Triggered: true})

	if err != nil {
		t.Errorf("should not error: %v", err)
	}
	if !app.IsUniversalKey("?") {
		t.Error("'?' should be registered as universal")
	}

	// Verify the stored message
	if msg, ok := app.universalBindings["?"].msg.(TestBindKeyMsg); !ok || !msg.Triggered {
		t.Error("message should be stored correctly")
	}
}

func TestBindKeyHardcodedConflict(t *testing.T) {
	app := New()

	// Use message-based binding (Elm pattern)
	err := app.BindKey("ctrl+c", "Quit", QuitMsg{})
	if err == nil {
		t.Error("should error when registering hardcoded key as universal")
	}
}

func TestMustBindKeyPanic(t *testing.T) {
	app := New()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustBindKey should panic on hardcoded conflict")
		}
	}()

	// Use message-based binding (Elm pattern)
	app.MustBindKey("ctrl+c", "Quit", QuitMsg{})
}

func TestKeyDispatchHardcoded(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	// Simulate ctrl+c keypress
	msg := KeyPress("ctrl+c")
	_, cmd := app.Update(msg)

	// Should return quit command
	if cmd == nil {
		t.Error("hardcoded ctrl+c should return a command")
	}
}

func TestKeyDispatchUniversal(t *testing.T) {
	app := New()

	// Use message-based binding with observable side effect (Elm pattern)
	// We'll use ToggleNavBarAutoHideMsg since it modifies model state
	app.MustBindKey("x", "Test", ToggleNavBarAutoHideMsg{})

	app.MustAddPage(page.New("Test"))

	initialAutohide := app.Model().NavBarAutoHide

	// Simulate 'x' keypress
	msg := KeyPress("x")
	app.Update(msg)

	// The message should have been dispatched and toggled autohide
	if app.Model().NavBarAutoHide == initialAutohide {
		t.Error("universal message should be dispatched and toggle autohide")
	}
}

func TestKeyDispatchPageScoped(t *testing.T) {
	app := New()

	// Use NavigateToPageMsg which has observable side effect
	p := page.New("Test").
		BindKey("a", "Go to page 2", NavigateToPageMsg{Index: 1})

	app.MustAddPage(p)
	app.MustAddPage(page.New("Test2")) // Add second page

	// Simulate 'a' keypress - should navigate to page 1
	msg := KeyPress("a")
	app.Update(msg)

	if app.Active() != 1 {
		t.Errorf("page message should trigger navigation, got active=%d, want 1", app.Active())
	}
}

func TestKeyDispatchPageScopedNotOnOtherPage(t *testing.T) {
	app := New()

	// Page1 has 'a' key that toggles autohide
	// Page2 has 'b' key that toggles tab position
	p1 := page.New("Page1").
		BindKey("a", "Toggle autohide", ToggleNavBarAutoHideMsg{})
	p2 := page.New("Page2").
		BindKey("b", "Toggle position", ToggleNavBarPositionMsg{})

	app.MustAddPage(p1)
	app.MustAddPage(p2)

	initialPosition := app.Model().NavBarPosition

	// Active page is 0 (Page1)
	// Press 'b' which is only on Page2 - should not change position
	msg := KeyPress("b")
	app.Update(msg)

	if app.Model().NavBarPosition != initialPosition {
		t.Error("page2 message should not be dispatched when page2 is not active")
	}

	// Switch to Page2
	app.SetActive(1)
	app.Update(msg)

	if app.Model().NavBarPosition == initialPosition {
		t.Error("page2 message should be dispatched when page2 is active")
	}
}

func TestKeyDispatchPriority(t *testing.T) {
	// This test verifies the dispatch order:
	// 1. Hardcoded > 2. Universal > 3. Page-scoped

	app := New()

	// Use message-based binding for universal key (Elm pattern)
	// ToggleNavBarModeMsg for universal 'x'
	app.MustBindKey("x", "Test", ToggleNavBarModeMsg{})

	// Test that page key works when not universal
	// Use ToggleNavBarAutoHideMsg for observable side effect
	p := page.New("Test").
		BindKey("y", "Toggle autohide", ToggleNavBarAutoHideMsg{})
	app.MustAddPage(p)

	initialAutohide := app.Model().NavBarAutoHide
	initialMode := app.Model().NavBarMode

	// Press 'y' - should trigger page message
	msg := KeyPress("y")
	app.Update(msg)
	if app.Model().NavBarAutoHide == initialAutohide {
		t.Error("page message for 'y' should be dispatched")
	}

	// Press 'x' - should trigger universal message
	msg = KeyPress("x")
	app.Update(msg)
	if app.Model().NavBarMode == initialMode {
		t.Error("universal message for 'x' should be dispatched and toggle mode")
	}
}

func TestTabNavigation(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Page1"))
	app.MustAddPage(page.New("Page2"))
	app.MustAddPage(page.New("Page3"))

	if app.Active() != 0 {
		t.Errorf("initial active should be 0, got %d", app.Active())
	}

	// Press tab
	msg := KeyPress("tab")
	app.Update(msg)
	if app.Active() != 1 {
		t.Errorf("after tab, active should be 1, got %d", app.Active())
	}

	// Press tab again
	app.Update(msg)
	if app.Active() != 2 {
		t.Errorf("after second tab, active should be 2, got %d", app.Active())
	}

	// Press tab wraps around
	app.Update(msg)
	if app.Active() != 0 {
		t.Errorf("after third tab, active should wrap to 0, got %d", app.Active())
	}
}

func TestShiftTabNavigation(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Page1"))
	app.MustAddPage(page.New("Page2"))
	app.MustAddPage(page.New("Page3"))

	// Press shift+tab from page 0 should wrap to page 2
	msg := KeyPress("shift+tab")
	app.Update(msg)
	if app.Active() != 2 {
		t.Errorf("after shift+tab from 0, active should be 2, got %d", app.Active())
	}
}

func TestNumberKeyNavigation(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Page1"))
	app.MustAddPage(page.New("Page2"))
	app.MustAddPage(page.New("Page3"))

	// Press '2' to jump to page 2 (index 1)
	msg := KeyPress("2")
	app.Update(msg)
	if app.Active() != 1 {
		t.Errorf("pressing '2' should go to page index 1, got %d", app.Active())
	}

	// Press '3' to jump to page 3 (index 2)
	msg = KeyPress("3")
	app.Update(msg)
	if app.Active() != 2 {
		t.Errorf("pressing '3' should go to page index 2, got %d", app.Active())
	}
}

// Overlay Tests

func TestShowOverlay(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().
		SetContent("Test overlay").
		SetExitKey("esc")

	_, err := app.ShowOverlay(overlay)
	if err != nil {
		t.Errorf("ShowOverlay should not error: %v", err)
	}
	if !app.HasOverlay() {
		t.Error("HasOverlay should return true")
	}
	if app.ActiveOverlay() != overlay {
		t.Error("ActiveOverlay should return the shown overlay")
	}
}

func TestShowOverlayNoExitKey(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().SetContent("Test")
	// No exit key set

	_, err := app.ShowOverlay(overlay)
	if err == nil {
		t.Error("ShowOverlay should error without exit key")
	}
}

func TestShowOverlayHardcodedConflict(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().
		SetExitKey("ctrl+c") // conflicts with hardcoded

	_, err := app.ShowOverlay(overlay)
	if err == nil {
		t.Error("ShowOverlay should error on hardcoded exit key conflict")
	}
}

func TestHideOverlay(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc")

	app.MustShowOverlay(overlay)
	app.HideOverlay()

	if app.HasOverlay() {
		t.Error("HasOverlay should return false after HideOverlay")
	}
	if app.ActiveOverlay() != nil {
		t.Error("ActiveOverlay should return nil after HideOverlay")
	}
}

func TestOverlayExitKey(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc")

	app.MustShowOverlay(overlay)

	// Press escape
	msg := KeyPress("esc")
	app.Update(msg)

	if app.HasOverlay() {
		t.Error("overlay should be hidden after pressing exit key")
	}
}

// testOverlayKeyMsg is a message for testing overlay key dispatch
type testOverlayKeyMsg struct{}

func TestOverlayKeyHandler(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	// Register a message to dispatch when 'y' is pressed (Elm pattern)
	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc").
		RegisterKey("y", testOverlayKeyMsg{})

	app.MustShowOverlay(overlay)

	// Press 'y' - should return a command that produces our test message
	keyMsg := KeyPress("y")
	_, cmd := app.Update(keyMsg)

	// The command should produce our test message when executed
	if cmd == nil {
		t.Fatal("overlay key should return a command")
	}

	resultMsg := cmd()
	if _, ok := resultMsg.(testOverlayKeyMsg); !ok {
		t.Errorf("expected testOverlayKeyMsg, got %T", resultMsg)
	}
}

func TestOverlayBlocksPageKeys(t *testing.T) {
	app := New()

	// Use ToggleNavBarAutoHideMsg to test
	p := page.New("Test").
		BindKey("a", "Toggle autohide", ToggleNavBarAutoHideMsg{})
	app.MustAddPage(p)

	initialAutohide := app.Model().NavBarAutoHide

	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc").
		SetBlocking(true)

	app.MustShowOverlay(overlay)

	// Press 'a' - should be blocked
	msg := KeyPress("a")
	app.Update(msg)

	if app.Model().NavBarAutoHide != initialAutohide {
		t.Error("page message should NOT be dispatched when blocking overlay is active")
	}
}

func TestOverlayNonBlockingAllowsPageKeys(t *testing.T) {
	app := New()

	// Use ToggleNavBarAutoHideMsg to test
	p := page.New("Test").
		BindKey("a", "Toggle autohide", ToggleNavBarAutoHideMsg{})
	app.MustAddPage(p)

	initialAutohide := app.Model().NavBarAutoHide

	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc").
		SetBlocking(false)

	app.MustShowOverlay(overlay)

	// Press 'a' - should pass through
	msg := KeyPress("a")
	app.Update(msg)

	if app.Model().NavBarAutoHide == initialAutohide {
		t.Error("page message SHOULD be dispatched when non-blocking overlay is active")
	}
}

func TestHardcodedKeysWorkWithOverlay(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay := overlays.New().
		SetContent("Test").
		SetExitKey("esc").
		SetBlocking(true)

	app.MustShowOverlay(overlay)

	// Press ctrl+c - hardcoded keys should always work
	msg := KeyPress("ctrl+c")
	_, cmd := app.Update(msg)

	if cmd == nil {
		t.Error("hardcoded ctrl+c should work even with blocking overlay")
	}
}

func TestShowNotification(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	// Explicitly set timeout to 0 to disable auto-dismiss
	notif := overlays.Success("Success!").SetExitKey("enter").SetTimeout(0)
	cmd := app.ShowNotification(notif)

	// No timeout set, so cmd should be nil
	if cmd != nil {
		t.Error("ShowNotification should return nil cmd when no timeout")
	}
	if !app.HasOverlay() {
		t.Error("HasOverlay should return true after ShowNotification")
	}
}

func TestShowNotificationWithTimeout(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	notif := overlays.Success("Success!").
		SetExitKey("enter").
		SetTimeout(100 * time.Millisecond)

	cmd := app.ShowNotification(notif)

	// Should return a timeout command
	if cmd == nil {
		t.Error("ShowNotification should return cmd when timeout is set")
	}
	if !app.HasOverlay() {
		t.Error("HasOverlay should return true after ShowNotification")
	}

	// The cmd is a tea.Tick - when executed, it will eventually send overlayTimeoutMsg
	// We can't easily test the full tick cycle in unit tests, but we can verify
	// that sending the timeout message dismisses the overlay
}

func TestOverlayTimeoutMsg(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	notif := overlays.Success("Success!").
		SetExitKey("enter").
		SetTimeout(100 * time.Millisecond)

	app.ShowNotification(notif)

	if !app.HasOverlay() {
		t.Error("HasOverlay should return true before timeout")
	}
	if app.OverlayCount() != 1 {
		t.Errorf("OverlayCount should be 1, got %d", app.OverlayCount())
	}

	// Get the overlay ID from the stack for testing timeout
	overlayID := app.overlayStack[0].id

	// Simulate the timeout message
	msg := overlayTimeoutMsg{overlayID: overlayID}
	app.Update(msg)

	if app.HasOverlay() {
		t.Error("HasOverlay should return false after timeout message")
	}
}

func TestOverlayTimeoutMsgWrongID(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	notif := overlays.Success("Success!").SetExitKey("enter")
	app.ShowNotification(notif)

	// Simulate a timeout message with wrong ID (stale timeout)
	msg := overlayTimeoutMsg{overlayID: 999}
	app.Update(msg)

	// Overlay should NOT be dismissed because IDs don't match
	if !app.HasOverlay() {
		t.Error("HasOverlay should still be true - wrong timeout ID")
	}
}

func TestShowModal(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	modal := overlays.NewModal().
		SetTitle("Confirm").
		SetExitKey("esc")
	_, err := app.ShowModal(modal)

	if err != nil {
		t.Errorf("ShowModal should not error: %v", err)
	}
	if !app.HasOverlay() {
		t.Error("HasOverlay should return true after ShowModal")
	}
}

// Z-axis Overlay Stacking Tests

func TestOverlayZIndexStacking(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	// Create overlays with different z-indices
	overlay1 := overlays.New().SetContent("Low").SetExitKey("1").SetZIndex(1)
	overlay2 := overlays.New().SetContent("High").SetExitKey("2").SetZIndex(10)
	overlay3 := overlays.New().SetContent("Mid").SetExitKey("3").SetZIndex(5)

	// Show in non-z-order (should sort by z-index)
	app.MustShowOverlay(overlay1) // z=1
	app.MustShowOverlay(overlay2) // z=10
	app.MustShowOverlay(overlay3) // z=5

	if app.OverlayCount() != 3 {
		t.Errorf("OverlayCount = %d, want 3", app.OverlayCount())
	}

	// Top overlay should be the highest z-index (overlay2)
	if app.TopOverlay() != overlay2 {
		t.Error("TopOverlay should be the highest z-index overlay")
	}

	// Stack order should be: overlay1 (z=1), overlay3 (z=5), overlay2 (z=10)
	if app.overlayStack[0].overlay.ZIndex() != 1 {
		t.Errorf("First overlay zIndex = %d, want 1", app.overlayStack[0].overlay.ZIndex())
	}
	if app.overlayStack[1].overlay.ZIndex() != 5 {
		t.Errorf("Second overlay zIndex = %d, want 5", app.overlayStack[1].overlay.ZIndex())
	}
	if app.overlayStack[2].overlay.ZIndex() != 10 {
		t.Errorf("Third overlay zIndex = %d, want 10", app.overlayStack[2].overlay.ZIndex())
	}
}

func TestOverlaySameZIndexOrder(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	// Create overlays with same z-index
	overlay1 := overlays.New().SetContent("First").SetExitKey("1").SetZIndex(5)
	overlay2 := overlays.New().SetContent("Second").SetExitKey("2").SetZIndex(5)
	overlay3 := overlays.New().SetContent("Third").SetExitKey("3").SetZIndex(5)

	app.MustShowOverlay(overlay1)
	app.MustShowOverlay(overlay2)
	app.MustShowOverlay(overlay3)

	// Same z-index: insertion order (last added = on top)
	if app.TopOverlay() != overlay3 {
		t.Error("TopOverlay should be the last added overlay when z-index is same")
	}
}

func TestHideOverlayByID(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay1 := overlays.New().SetContent("One").SetExitKey("1")
	overlay2 := overlays.New().SetContent("Two").SetExitKey("2")

	id1, _ := app.ShowOverlay(overlay1)
	app.MustShowOverlay(overlay2)

	if app.OverlayCount() != 2 {
		t.Errorf("OverlayCount = %d, want 2", app.OverlayCount())
	}

	// Hide the first overlay by ID
	removed := app.HideOverlayByID(id1)
	if !removed {
		t.Error("HideOverlayByID should return true")
	}
	if app.OverlayCount() != 1 {
		t.Errorf("OverlayCount after hide = %d, want 1", app.OverlayCount())
	}
	if app.TopOverlay() != overlay2 {
		t.Error("Remaining overlay should be overlay2")
	}
}

func TestHideAllOverlays(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	app.MustShowOverlay(overlays.New().SetContent("One").SetExitKey("1"))
	app.MustShowOverlay(overlays.New().SetContent("Two").SetExitKey("2"))
	app.MustShowOverlay(overlays.New().SetContent("Three").SetExitKey("3"))

	if app.OverlayCount() != 3 {
		t.Errorf("OverlayCount = %d, want 3", app.OverlayCount())
	}

	app.HideAllOverlays()

	if app.HasOverlay() {
		t.Error("HasOverlay should return false after HideAllOverlays")
	}
	if app.OverlayCount() != 0 {
		t.Errorf("OverlayCount = %d, want 0", app.OverlayCount())
	}
}

func TestBringToFront(t *testing.T) {
	app := New()
	app.MustAddPage(page.New("Test"))

	overlay1 := overlays.New().SetContent("One").SetExitKey("1").SetZIndex(1)
	overlay2 := overlays.New().SetContent("Two").SetExitKey("2").SetZIndex(10)

	id1, _ := app.ShowOverlay(overlay1)
	app.MustShowOverlay(overlay2)

	// overlay2 should be on top (z=10)
	if app.TopOverlay() != overlay2 {
		t.Error("overlay2 should be on top initially")
	}

	// Bring overlay1 to front
	app.BringToFront(id1)

	// Now overlay1 should be on top (z=11)
	if app.TopOverlay() != overlay1 {
		t.Error("overlay1 should be on top after BringToFront")
	}
	if overlay1.ZIndex() != 11 {
		t.Errorf("overlay1 zIndex = %d, want 11", overlay1.ZIndex())
	}
}

// Section Tests

func TestNewSection(t *testing.T) {
	section := NewSection("Settings")
	if section.Name() != "Settings" {
		t.Errorf("Name() = %q, want %q", section.Name(), "Settings")
	}
	if len(section.Pages()) != 0 {
		t.Errorf("Pages() length = %d, want 0", len(section.Pages()))
	}
	if len(section.Sections()) != 0 {
		t.Errorf("Sections() length = %d, want 0", len(section.Sections()))
	}
}

func TestSectionAddPage(t *testing.T) {
	section := NewSection("Settings").
		AddPage(page.New("General")).
		AddPage(page.New("Advanced"))

	if len(section.Pages()) != 2 {
		t.Errorf("Pages() length = %d, want 2", len(section.Pages()))
	}
	if section.Pages()[0].Name() != "General" {
		t.Errorf("Pages()[0].Name() = %q, want %q", section.Pages()[0].Name(), "General")
	}
}

func TestSectionAddPages(t *testing.T) {
	section := NewSection("Settings").
		AddPages(page.New("General"), page.New("Advanced"), page.New("Debug"))

	if len(section.Pages()) != 3 {
		t.Errorf("Pages() length = %d, want 3", len(section.Pages()))
	}
}

func TestSectionAddSection(t *testing.T) {
	parent := NewSection("Parent").
		AddSection(NewSection("Child1")).
		AddSection(NewSection("Child2"))

	if len(parent.Sections()) != 2 {
		t.Errorf("Sections() length = %d, want 2", len(parent.Sections()))
	}
	if parent.Sections()[0].Name() != "Child1" {
		t.Errorf("Sections()[0].Name() = %q, want %q", parent.Sections()[0].Name(), "Child1")
	}
}

func TestSectionDepth(t *testing.T) {
	// Single level section (pages only)
	leaf := NewSection("Leaf").AddPage(page.New("Page1"))
	if leaf.Depth() != 1 {
		t.Errorf("leaf.Depth() = %d, want 1", leaf.Depth())
	}

	// Two level section
	parent := NewSection("Parent").
		AddSection(NewSection("Child").AddPage(page.New("Page1")))
	if parent.Depth() != 2 {
		t.Errorf("parent.Depth() = %d, want 2", parent.Depth())
	}

	// Three level section
	root := NewSection("Root").
		AddSection(NewSection("L1").
			AddSection(NewSection("L2").AddPage(page.New("Page1"))))
	if root.Depth() != 3 {
		t.Errorf("root.Depth() = %d, want 3", root.Depth())
	}
}

func TestSectionIsLeaf(t *testing.T) {
	leaf := NewSection("Leaf").AddPage(page.New("Page1"))
	if !leaf.IsLeaf() {
		t.Error("IsLeaf() should return true for section with only pages")
	}

	parent := NewSection("Parent").AddSection(NewSection("Child"))
	if parent.IsLeaf() {
		t.Error("IsLeaf() should return false for section with child sections")
	}
}

func TestAddSection(t *testing.T) {
	app := New()

	mainSection := NewSection("Main").AddPage(page.New("Home"))
	err := app.AddSection(mainSection)

	if err != nil {
		t.Errorf("AddSection should not error: %v", err)
	}
	if !app.HasSections() {
		t.Error("HasSections should return true")
	}
}

func TestAddSectionEmpty(t *testing.T) {
	app := New()

	// Empty sections are allowed (flexibility for later population)
	emptySection := NewSection("Empty") // No pages or children
	err := app.AddSection(emptySection)

	if err != nil {
		t.Errorf("AddSection should allow empty sections: %v", err)
	}
}

func TestMustAddSection(t *testing.T) {
	app := New()

	mainSection := NewSection("Main").AddPage(page.New("Home"))
	result := app.MustAddSection(mainSection)

	if result != app {
		t.Error("MustAddSection should return self for chaining")
	}
}

func TestMustAddSectionWithEmptySection(t *testing.T) {
	app := New()

	// Empty sections are allowed - should not panic
	result := app.MustAddSection(NewSection("Empty"))
	if result != app {
		t.Error("MustAddSection should return self for chaining")
	}
}

func TestSwitchSection(t *testing.T) {
	app := New()

	mainSection := NewSection("Main").AddPage(page.New("Home"))
	settingsSection := NewSection("Settings").AddPage(page.New("General"))

	app.MustAddSection(mainSection)
	app.MustAddSection(settingsSection)

	// Start at section 0
	if app.CurrentSectionIndex(0) != 0 {
		t.Errorf("CurrentSectionIndex(0) = %d, want 0", app.CurrentSectionIndex(0))
	}

	// Switch to section 1
	err := app.SwitchSection(0, 1)
	if err != nil {
		t.Errorf("SwitchSection should not error: %v", err)
	}

	if app.CurrentSectionIndex(0) != 1 {
		t.Errorf("CurrentSectionIndex(0) = %d, want 1", app.CurrentSectionIndex(0))
	}
}

func TestSwitchSectionInvalid(t *testing.T) {
	app := New()

	mainSection := NewSection("Main").AddPage(page.New("Home"))
	app.MustAddSection(mainSection)

	err := app.SwitchSection(0, 99) // Invalid index
	if err == nil {
		t.Error("SwitchSection should error for invalid index")
	}
}

func TestBreadcrumb(t *testing.T) {
	bc := Breadcrumb{level: 0, path: []int{1, 2}}

	if bc.SectionIndex() != 1 {
		t.Errorf("SectionIndex() = %d, want 1", bc.SectionIndex())
	}
	if bc.PageIndex() != 2 {
		t.Errorf("PageIndex() = %d, want 2", bc.PageIndex())
	}
}

func TestSectionSetName(t *testing.T) {
	section := NewSection("Old").SetName("New")
	if section.Name() != "New" {
		t.Errorf("Name() = %q, want %q", section.Name(), "New")
	}
}

func TestFlatPagesWithSections(t *testing.T) {
	app := New()

	// Add flat pages first
	app.MustAddPage(page.New("Page1"))
	app.MustAddPage(page.New("Page2"))

	// Should NOT have sections
	if app.HasSections() {
		t.Error("HasSections should return false when using flat pages")
	}
}

func TestNestedSections(t *testing.T) {
	app := New()

	// Create 3-level structure: (Root) > Main > Settings > Pages
	// Root (anonymous) -> Main -> Settings -> Pages
	// Depth: Root=3, Main=2, Settings=1
	mainSection := NewSection("Main").
		AddSection(NewSection("Settings").
			AddPage(page.New("General")).
			AddPage(page.New("Advanced")))

	app.MustAddSection(mainSection)

	// hierarchyLevels = root.Depth() = 3
	if app.hierarchyLevels() != 3 {
		t.Errorf("hierarchyLevels() = %d, want 3", app.hierarchyLevels())
	}
}
