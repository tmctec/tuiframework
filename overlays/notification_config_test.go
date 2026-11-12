package overlays

import (
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/styles"
)

// restoreNotificationConfig saves NotificationConfig and restores it when the
// test ends, so a test can change the defaults freely.
func restoreNotificationConfig(t *testing.T) {
	t.Helper()
	saved := *NotificationConfig
	t.Cleanup(func() { *NotificationConfig = saved })
}

// testNotifMsg is a message for testing notification key registration.
type testNotifMsg struct{}

func TestNewNotificationAppliesConfig(t *testing.T) {
	restoreNotificationConfig(t)
	NotificationConfig.
		SetPosition(PositionBottomLeft).
		SetBlocking(true).
		SetExitKey("q").
		SetTimeout(7 * time.Second).
		SetWidth(20).
		SetHeight(4)

	n := NewNotification("hi")
	if n.Position() != PositionBottomLeft {
		t.Errorf("Position() = %s, want LB", n.Position())
	}
	if !n.IsBlocking() {
		t.Error("IsBlocking() = false, want the configured true")
	}
	if n.ExitKey() != "q" {
		t.Errorf("ExitKey() = %q, want q", n.ExitKey())
	}
	if n.Timeout() != 7*time.Second {
		t.Errorf("Timeout() = %v, want 7s", n.Timeout())
	}
	if n.WidthFixed() != 20 || n.HeightFixed() != 4 {
		t.Errorf("fixed size = %dx%d, want 20x4", n.WidthFixed(), n.HeightFixed())
	}
	if n.WidthPct() != 0 || n.HeightPct() != 0 {
		t.Errorf("percentages = %dx%d, want 0x0 when fixed sizes apply", n.WidthPct(), n.HeightPct())
	}
	if n.Message() != "hi" || n.Level() != LevelInfo {
		t.Errorf("Message() = %q Level() = %s, want hi and Info", n.Message(), n.Level())
	}
}

func TestNewNotificationUsesPercentagesWhenFixedSizeIsZero(t *testing.T) {
	restoreNotificationConfig(t)
	NotificationConfig.SetWidth(0).SetHeight(0).SetWidthPct(40).SetHeightPct(10)

	n := NewNotification("hi")
	if n.WidthPct() != 40 || n.HeightPct() != 10 {
		t.Errorf("percentages = %dx%d, want 40x10", n.WidthPct(), n.HeightPct())
	}
	if n.WidthFixed() != 0 || n.HeightFixed() != 0 {
		t.Errorf("fixed size = %dx%d, want 0x0", n.WidthFixed(), n.HeightFixed())
	}
	if w, h := n.CalculateDimensions(100, 50); w != 40 || h != 5 {
		t.Errorf("CalculateDimensions = %dx%d, want 40x5", w, h)
	}
}

func TestNewNotificationFallsBackToSizeWithoutAnySizing(t *testing.T) {
	restoreNotificationConfig(t)
	NotificationConfig.SetWidth(0).SetHeight(0).SetWidthPct(0).SetHeightPct(0)

	n := NewNotification("hi")
	if n.WidthPct() != 0 || n.HeightPct() != 0 || n.WidthFixed() != 0 || n.HeightFixed() != 0 {
		t.Error("no sizing should be applied when the config has none")
	}
	if w, h := n.CalculateDimensions(100, 50); w != 80 || h != 40 {
		t.Errorf("CalculateDimensions = %dx%d, want the default 80%% size", w, h)
	}
}

func TestNotificationConstructors(t *testing.T) {
	tests := []struct {
		name        string
		constructor func(string) *Notification
		level       NotificationLevel
	}{
		{"Info", Info, LevelInfo},
		{"Note", Note, LevelNote},
		{"Tip", Tip, LevelTip},
		{"Important", Important, LevelImportant},
		{"Success", Success, LevelSuccess},
		{"Warning", Warning, LevelWarning},
		{"Caution", Caution, LevelCaution},
		{"Error", Error, LevelError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := tt.constructor("msg")
			ls := NotificationConfig.LevelStyle(tt.level)
			if n.Level() != tt.level {
				t.Errorf("Level() = %s, want %s", n.Level(), tt.level)
			}
			if want := ls.Icon + " msg"; n.Content() != want {
				t.Errorf("Content() = %q, want %q", n.Content(), want)
			}
			if n.Message() != "msg" {
				t.Errorf("Message() = %q, want msg", n.Message())
			}
			s := n.Style()
			if s.GetBackground() != ls.Background || s.GetForeground() != ls.Foreground {
				t.Errorf("colors = %v/%v, want %v/%v", s.GetForeground(), s.GetBackground(), ls.Foreground, ls.Background)
			}
			if s.GetBorderTopForeground() != ls.Border || s.GetBorderTopBackground() != ls.Background {
				t.Error("border colors should come from the level style")
			}
			if s.GetBorderStyle() != NotificationConfig.Border() {
				t.Error("border style should come from the config")
			}
			if v, h := NotificationConfig.Padding(); s.GetPaddingTop() != v || s.GetPaddingLeft() != h {
				t.Errorf("padding = %d,%d, want %d,%d", s.GetPaddingTop(), s.GetPaddingLeft(), v, h)
			}
		})
	}
}

func TestNotificationSetContentKeepsIcon(t *testing.T) {
	n := Success("first").SetContent("second")
	if want := NotificationConfig.LevelStyle(LevelSuccess).Icon + " second"; n.Content() != want {
		t.Errorf("Content() = %q, want %q", n.Content(), want)
	}
	if n.Message() != "second" {
		t.Errorf("Message() = %q, want second", n.Message())
	}
}

func TestNotificationSetLevelSwapsIcon(t *testing.T) {
	n := Info("msg").SetLevel(LevelError)
	if want := NotificationConfig.LevelStyle(LevelError).Icon + " msg"; n.Content() != want {
		t.Errorf("Content() = %q, want %q", n.Content(), want)
	}
	if n.Style().GetBorderTopForeground() != NotificationConfig.LevelStyle(LevelError).Border {
		t.Error("SetLevel should restyle from the new level")
	}
}

func TestNotificationWithoutIconShowsMessageOnly(t *testing.T) {
	restoreNotificationConfig(t)
	NotificationConfig.SetLevelStyle(LevelTip, styles.NotificationStyle{
		Foreground: lipgloss.Color("#ffffff"),
		Background: lipgloss.Color("#000000"),
		Border:     lipgloss.Color("#888888"),
	})

	n := Tip("plain")
	if n.Content() != "plain" {
		t.Errorf("Content() = %q, want the bare message", n.Content())
	}
}

func TestNotificationSetStyleOverridesLevelStyle(t *testing.T) {
	custom := lipgloss.NewStyle().Background(lipgloss.Color("#123456"))
	n := Error("msg")
	if r := n.SetStyle(custom); r != n {
		t.Error("SetStyle should return the receiver")
	}
	if n.Style().GetBackground() != lipgloss.Color("#123456") {
		t.Errorf("background = %v, want the custom style", n.Style().GetBackground())
	}

	// SetLevel restyles from the level again.
	n.SetLevel(LevelError)
	if n.Style().GetBackground() != NotificationConfig.LevelStyle(LevelError).Background {
		t.Error("SetLevel should apply the level style over a custom style")
	}
}

func TestNotificationSizeSetters(t *testing.T) {
	n := NewNotification("msg")

	if r := n.SetWidthPct(30); r != n || n.WidthPct() != 30 || n.WidthFixed() != 0 {
		t.Errorf("SetWidthPct: pct %d fixed %d", n.WidthPct(), n.WidthFixed())
	}
	if r := n.SetHeightPct(20); r != n || n.HeightPct() != 20 || n.HeightFixed() != 0 {
		t.Errorf("SetHeightPct: pct %d fixed %d", n.HeightPct(), n.HeightFixed())
	}
	if r := n.SetWidth(12); r != n || n.WidthFixed() != 12 || n.WidthPct() != 0 {
		t.Errorf("SetWidth: fixed %d pct %d", n.WidthFixed(), n.WidthPct())
	}
	if r := n.SetHeight(2); r != n || n.HeightFixed() != 2 || n.HeightPct() != 0 {
		t.Errorf("SetHeight: fixed %d pct %d", n.HeightFixed(), n.HeightPct())
	}
	if r := n.SetSize(25); r != n || n.Size() != 25 {
		t.Errorf("SetSize: size %d", n.Size())
	}
}

func TestNotificationStyleSetters(t *testing.T) {
	border := lipgloss.Color("#ff0000")
	borderBg := lipgloss.Color("#00ff00")
	bg := lipgloss.Color("#0000ff")
	fg := lipgloss.Color("#ffffff")

	n := NewNotification("msg")
	results := []*Notification{
		n.SetBorder(lipgloss.DoubleBorder()),
		n.SetBorderColor(border),
		n.SetBorderBackground(borderBg),
		n.SetPadding(1, 2),
		n.SetBackground(bg),
		n.SetForeground(fg),
	}
	for i, r := range results {
		if r != n {
			t.Errorf("setter %d did not return the receiver", i)
		}
	}

	s := n.Style()
	if s.GetBorderStyle() != lipgloss.DoubleBorder() {
		t.Error("SetBorder did not apply")
	}
	if s.GetBorderTopForeground() != border {
		t.Errorf("border color = %v, want %v", s.GetBorderTopForeground(), border)
	}
	if s.GetBorderTopBackground() != borderBg {
		t.Errorf("border background = %v, want %v", s.GetBorderTopBackground(), borderBg)
	}
	if s.GetPaddingTop() != 1 || s.GetPaddingLeft() != 2 {
		t.Errorf("padding = %d,%d, want 1,2", s.GetPaddingTop(), s.GetPaddingLeft())
	}
	if s.GetBackground() != bg || s.GetForeground() != fg {
		t.Errorf("colors = %v/%v, want %v/%v", s.GetForeground(), s.GetBackground(), fg, bg)
	}
}

func TestNotificationKeysAndZIndex(t *testing.T) {
	n := NewNotification("msg")
	if r := n.RegisterKey("y", testNotifMsg{}); r != n {
		t.Error("RegisterKey should return the receiver")
	}
	if !n.HasKey("y") {
		t.Error("HasKey(y) = false, want true")
	}
	if _, ok := n.Keys()["y"].(testNotifMsg); !ok {
		t.Error("Keys()[y] should hold the registered message")
	}

	if r := n.SetZIndex(3); r != n || n.ZIndex() != 3 {
		t.Errorf("SetZIndex: ZIndex() = %d, want 3", n.ZIndex())
	}
}

func TestNotificationSetContentReturnsSelf(t *testing.T) {
	n := NewNotification("a")
	if n.SetContent("b") != n {
		t.Error("SetContent should return the receiver")
	}
}
