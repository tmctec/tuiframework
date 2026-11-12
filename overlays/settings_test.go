package overlays

import (
	"testing"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/styles"
)

// allLevels lists every notification level in severity order.
var allLevels = []NotificationLevel{
	LevelInfo, LevelNote, LevelTip, LevelImportant,
	LevelSuccess, LevelWarning, LevelCaution, LevelError,
}

// themeLevelStyle returns the theme's entry for a level.
func themeLevelStyle(t *styles.Theme, level NotificationLevel) styles.NotificationStyle {
	switch level {
	case LevelNote:
		return t.NotifNote
	case LevelTip:
		return t.NotifTip
	case LevelImportant:
		return t.NotifImportant
	case LevelSuccess:
		return t.NotifSuccess
	case LevelWarning:
		return t.NotifWarning
	case LevelCaution:
		return t.NotifCaution
	case LevelError:
		return t.NotifError
	default:
		return t.NotifInfo
	}
}

func TestNotificationLevelString(t *testing.T) {
	tests := []struct {
		level NotificationLevel
		want  string
	}{
		{LevelInfo, "Info"},
		{LevelNote, "Note"},
		{LevelTip, "Tip"},
		{LevelImportant, "Important"},
		{LevelSuccess, "Success"},
		{LevelWarning, "Warning"},
		{LevelCaution, "Caution"},
		{LevelError, "Error"},
		{NotificationLevel(-1), "Info"},
		{levelCount, "Info"},
	}
	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("NotificationLevel(%d).String() = %q, want %q", int(tt.level), got, tt.want)
		}
	}
}

func TestNotificationSettingsSettersRoundTrip(t *testing.T) {
	c := &NotificationSettings{}

	if r := c.SetPosition(PositionBottomLeft); r != c || c.Position() != PositionBottomLeft {
		t.Errorf("SetPosition: got %v (self %v)", c.Position(), r == c)
	}
	if r := c.SetWidth(42); r != c || c.Width() != 42 {
		t.Errorf("SetWidth: got %d (self %v)", c.Width(), r == c)
	}
	if r := c.SetHeight(5); r != c || c.Height() != 5 {
		t.Errorf("SetHeight: got %d (self %v)", c.Height(), r == c)
	}
	if r := c.SetWidthPct(60); r != c || c.WidthPct() != 60 {
		t.Errorf("SetWidthPct: got %d (self %v)", c.WidthPct(), r == c)
	}
	if r := c.SetHeightPct(15); r != c || c.HeightPct() != 15 {
		t.Errorf("SetHeightPct: got %d (self %v)", c.HeightPct(), r == c)
	}
	r := c.SetPadding(1, 4)
	if v, h := c.Padding(); r != c || v != 1 || h != 4 {
		t.Errorf("SetPadding: got %d,%d (self %v)", v, h, r == c)
	}
	if r := c.SetBorder(lipgloss.ThickBorder()); r != c || c.Border() != lipgloss.ThickBorder() {
		t.Errorf("SetBorder: got %v (self %v)", c.Border(), r == c)
	}
	if r := c.SetBlocking(true); r != c || !c.Blocking() {
		t.Errorf("SetBlocking: got %v (self %v)", c.Blocking(), r == c)
	}
	if r := c.SetTimeout(9 * time.Second); r != c || c.Timeout() != 9*time.Second {
		t.Errorf("SetTimeout: got %v (self %v)", c.Timeout(), r == c)
	}
	if r := c.SetExitKey("q"); r != c || c.ExitKey() != "q" {
		t.Errorf("SetExitKey: got %q (self %v)", c.ExitKey(), r == c)
	}
}

func TestNotificationSettingsApplyTheme(t *testing.T) {
	c := &NotificationSettings{}
	if r := c.ApplyTheme(styles.Retro); r != c {
		t.Error("ApplyTheme should return the receiver")
	}

	for _, level := range allLevels {
		if got, want := c.LevelStyle(level), themeLevelStyle(styles.Retro, level); got != want {
			t.Errorf("LevelStyle(%s) = %+v, want %+v", level, got, want)
		}
	}
	if c.Border() != styles.Retro.NotifBorder || c.Border() != lipgloss.DoubleBorder() {
		t.Errorf("Border() = %v, want the Retro theme's double border", c.Border())
	}
	if c.LevelStyle(LevelInfo).Icon != "[i]" {
		t.Errorf("Info icon = %q, want the Retro theme's [i]", c.LevelStyle(LevelInfo).Icon)
	}

	// Re-applying a different theme replaces every entry.
	c.ApplyTheme(styles.Dracula)
	if c.LevelStyle(LevelError) != styles.Dracula.NotifError || c.Border() != styles.Dracula.NotifBorder {
		t.Error("ApplyTheme should replace the previous theme's styles")
	}
}

func TestNotificationSettingsApplyThemeLeavesLayoutAlone(t *testing.T) {
	c := (&NotificationSettings{}).
		SetPosition(PositionBottomCenter).
		SetWidth(33).
		SetTimeout(time.Second).
		SetExitKey("x").
		ApplyTheme(styles.Nord)

	if c.Position() != PositionBottomCenter || c.Width() != 33 || c.Timeout() != time.Second || c.ExitKey() != "x" {
		t.Error("ApplyTheme should only touch colors, icons and the border")
	}
}

func TestNotificationSettingsLevelStyleOutOfRange(t *testing.T) {
	c := (&NotificationSettings{}).ApplyTheme(styles.Dracula)
	info := c.LevelStyle(LevelInfo)

	for _, level := range []NotificationLevel{-1, levelCount, 100} {
		if got := c.LevelStyle(level); got != info {
			t.Errorf("LevelStyle(%d) = %+v, want the Info style", int(level), got)
		}
	}
}

func TestNotificationSettingsSetLevelStyle(t *testing.T) {
	c := (&NotificationSettings{}).ApplyTheme(styles.Dracula)
	custom := styles.NotificationStyle{
		Foreground: lipgloss.Color("#111111"),
		Background: lipgloss.Color("#222222"),
		Border:     lipgloss.Color("#333333"),
		Icon:       "[custom]",
	}

	if r := c.SetLevelStyle(LevelWarning, custom); r != c {
		t.Error("SetLevelStyle should return the receiver")
	}
	if c.LevelStyle(LevelWarning) != custom {
		t.Errorf("LevelStyle(Warning) = %+v, want %+v", c.LevelStyle(LevelWarning), custom)
	}
	if c.LevelStyle(LevelCaution) != styles.Dracula.NotifCaution {
		t.Error("SetLevelStyle should leave other levels alone")
	}
}

func TestNotificationSettingsSetLevelStyleIgnoresOutOfRange(t *testing.T) {
	c := (&NotificationSettings{}).ApplyTheme(styles.Dracula)
	before := c.levels
	custom := styles.NotificationStyle{Icon: "[nope]"}

	c.SetLevelStyle(-1, custom)
	c.SetLevelStyle(levelCount, custom)
	c.SetLevelStyle(100, custom)

	if c.levels != before {
		t.Error("an out-of-range level should not change any entry")
	}
}

func TestNotificationConfigDefaults(t *testing.T) {
	c := NotificationConfig
	if c.Position() != PositionTopRight {
		t.Errorf("Position() = %s, want RT", c.Position())
	}
	if c.Width() != 50 || c.Height() != 3 {
		t.Errorf("size = %dx%d, want 50x3", c.Width(), c.Height())
	}
	if c.WidthPct() != 0 || c.HeightPct() != 0 {
		t.Errorf("percentages = %dx%d, want 0x0", c.WidthPct(), c.HeightPct())
	}
	if v, h := c.Padding(); v != 0 || h != 1 {
		t.Errorf("Padding() = %d,%d, want 0,1", v, h)
	}
	if c.Blocking() {
		t.Error("notifications should not block by default")
	}
	if c.Timeout() != 3*time.Second {
		t.Errorf("Timeout() = %v, want 3s", c.Timeout())
	}
	if c.ExitKey() != "esc" {
		t.Errorf("ExitKey() = %q, want esc", c.ExitKey())
	}
	if c.Border() != styles.Dracula.NotifBorder {
		t.Error("Border() should come from the Dracula theme")
	}
	for _, level := range allLevels {
		if c.LevelStyle(level) != themeLevelStyle(styles.Dracula, level) {
			t.Errorf("LevelStyle(%s) should come from the Dracula theme", level)
		}
	}
}
