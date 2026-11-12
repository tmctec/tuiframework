package overlays

import (
	"testing"
	"time"
)

func TestNewNotification(t *testing.T) {
	n := NewNotification("Test message")

	// Content includes the icon prefix from the default level (Info)
	expectedContent := NotificationConfig.LevelStyle(LevelInfo).Icon + " Test message"
	if n.Content() != expectedContent {
		t.Errorf("Content() = %q, want %q", n.Content(), expectedContent)
	}
	if n.Level() != LevelInfo {
		t.Errorf("Level() = %d, want LevelInfo", n.Level())
	}
	if n.IsBlocking() != NotificationConfig.Blocking() {
		t.Errorf("IsBlocking() = %v, want %v", n.IsBlocking(), NotificationConfig.Blocking())
	}
	if n.Position() != NotificationConfig.Position() {
		t.Errorf("Position() = %s, want %s", n.Position(), NotificationConfig.Position())
	}
	// WidthPct should be applied from config
	if NotificationConfig.WidthPct() > 0 && n.WidthPct() != NotificationConfig.WidthPct() {
		t.Errorf("WidthPct() = %d, want %d", n.WidthPct(), NotificationConfig.WidthPct())
	}
}

func TestNotificationLevels(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(string) *Notification
		wantLevel   NotificationLevel
	}{
		{"Info", Info, LevelInfo},
		{"Success", Success, LevelSuccess},
		{"Warning", Warning, LevelWarning},
		{"Error", Error, LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.constructor("message")
			if n.Level() != tt.wantLevel {
				t.Errorf("Level() = %d, want %d", n.Level(), tt.wantLevel)
			}
		})
	}
}

func TestNotificationSetLevel(t *testing.T) {
	n := NewNotification("msg").SetLevel(LevelError)

	if n.Level() != LevelError {
		t.Errorf("Level() = %d, want LevelError", n.Level())
	}
}

func TestNotificationTimeout(t *testing.T) {
	n := NewNotification("msg").SetTimeout(3 * time.Second)

	if n.Timeout() != 3*time.Second {
		t.Errorf("Timeout() = %v, want 3s", n.Timeout())
	}
}

func TestNotificationChaining(t *testing.T) {
	n := NewNotification("msg")
	result := n.
		SetLevel(LevelSuccess).
		SetTimeout(5 * time.Second).
		SetPosition(PositionTopCenter).
		SetSize(30).
		SetExitKey("enter").
		SetBlocking(true)

	if result != n {
		t.Error("Chaining broken - setter did not return self")
	}
}

func TestNotificationInheritsOverlay(t *testing.T) {
	n := NewNotification("msg").
		SetExitKey("esc").
		SetPosition(PositionBottomRight)

	// Should have overlay functionality
	if n.ExitKey() != "esc" {
		t.Errorf("ExitKey() = %q, want %q", n.ExitKey(), "esc")
	}
	if n.Position() != PositionBottomRight {
		t.Errorf("Position() = %s, want RB", n.Position())
	}

	// Validate should work
	err := n.Validate()
	if err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
}

func TestNotificationStyleByLevel(t *testing.T) {
	// Different levels should have different styles
	info := Info("msg")
	success := Success("msg")

	// They should have different background colors (styles are applied)
	// Just verify they don't panic and return valid overlays
	if info.Level() != LevelInfo {
		t.Error("Info() should create LevelInfo notification")
	}
	if success.Level() != LevelSuccess {
		t.Error("Success() should create LevelSuccess notification")
	}
}
