package overlays

import (
	"image/color"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/tuiframework/internal/keyname"
	"github.com/tmctec/tuiframework/styles"
)

// NotificationLevel is the kind of message a notification carries. It picks
// the notification's colors and icon from NotificationConfig.
type NotificationLevel int

// Notification levels, from least to most severe.
const (
	LevelInfo      NotificationLevel = iota // General information
	LevelNote                               // Highlighted information
	LevelTip                                // Optional helpful info
	LevelImportant                          // Crucial info for success
	LevelSuccess                            // Operation succeeded
	LevelWarning                            // Potential risks
	LevelCaution                            // Negative consequences
	LevelError                              // Something went wrong

	levelCount // number of levels; keep last
)

// String returns the level's name, such as "Info" or "Error".
func (l NotificationLevel) String() string {
	switch l {
	case LevelNote:
		return "Note"
	case LevelTip:
		return "Tip"
	case LevelImportant:
		return "Important"
	case LevelSuccess:
		return "Success"
	case LevelWarning:
		return "Warning"
	case LevelCaution:
		return "Caution"
	case LevelError:
		return "Error"
	default:
		return "Info"
	}
}

// NotificationSettings holds the defaults every new Notification starts
// from: where it appears, how big it is, whether it blocks, how long it
// stays, which key dismisses it, and the colors and icon of each level.
type NotificationSettings struct {
	position  Position
	width     int // fixed width in cells (0 = use widthPct)
	height    int // fixed height in lines (0 = use heightPct, or auto)
	widthPct  int // width as a percentage of the screen (0 = use width)
	heightPct int // height as a percentage of the screen (0 = auto)
	paddingV  int // padding above and below the text
	paddingH  int // padding left and right of the text
	border    lipgloss.Border
	blocking  bool
	timeout   time.Duration // auto-dismiss delay; 0 keeps the notification up
	exitKey   string

	levels [levelCount]styles.NotificationStyle
}

// NotificationConfig holds the process-wide defaults for new notifications:
// top right, 50 cells by 3 lines, non-blocking, dismissed by esc or after
// three seconds, with the Dracula theme's colors and icons. TabApp.UseTheme
// replaces the colors and icons with the active theme's.
var NotificationConfig = (&NotificationSettings{
	position: PositionTopRight,
	width:    50,
	height:   3,
	paddingH: 1,
	timeout:  3 * time.Second,
	exitKey:  "esc",
}).ApplyTheme(styles.Dracula)

// ApplyTheme takes the colors and icon of every level, and the border, from t.
func (c *NotificationSettings) ApplyTheme(t *styles.Theme) *NotificationSettings {
	c.levels = [levelCount]styles.NotificationStyle{
		LevelInfo:      t.NotifInfo,
		LevelNote:      t.NotifNote,
		LevelTip:       t.NotifTip,
		LevelImportant: t.NotifImportant,
		LevelSuccess:   t.NotifSuccess,
		LevelWarning:   t.NotifWarning,
		LevelCaution:   t.NotifCaution,
		LevelError:     t.NotifError,
	}
	c.border = t.NotifBorder
	return c
}

// LevelStyle returns the colors and icon of a level. An unknown level
// returns the Info style.
func (c *NotificationSettings) LevelStyle(level NotificationLevel) styles.NotificationStyle {
	if level < 0 || level >= levelCount {
		level = LevelInfo
	}
	return c.levels[level]
}

// SetLevelStyle sets the colors and icon of a level. An unknown level is ignored.
func (c *NotificationSettings) SetLevelStyle(level NotificationLevel, s styles.NotificationStyle) *NotificationSettings {
	if level >= 0 && level < levelCount {
		c.levels[level] = s
	}
	return c
}

// Position returns where new notifications appear.
func (c *NotificationSettings) Position() Position { return c.position }

// SetPosition sets where new notifications appear.
func (c *NotificationSettings) SetPosition(p Position) *NotificationSettings {
	c.position = p
	return c
}

// Width returns the fixed width in cells; 0 means WidthPct applies.
func (c *NotificationSettings) Width() int { return c.width }

// SetWidth sets a fixed width in cells. 0 makes WidthPct apply.
func (c *NotificationSettings) SetWidth(w int) *NotificationSettings {
	c.width = w
	return c
}

// Height returns the fixed height in lines; 0 means HeightPct applies.
func (c *NotificationSettings) Height() int { return c.height }

// SetHeight sets a fixed height in lines. 0 makes HeightPct apply.
func (c *NotificationSettings) SetHeight(h int) *NotificationSettings {
	c.height = h
	return c
}

// WidthPct returns the width as a percentage of the screen, used when Width is 0.
func (c *NotificationSettings) WidthPct() int { return c.widthPct }

// SetWidthPct sets the width as a percentage of the screen, used when Width is 0.
func (c *NotificationSettings) SetWidthPct(p int) *NotificationSettings {
	c.widthPct = p
	return c
}

// HeightPct returns the height as a percentage of the screen, used when Height is 0.
func (c *NotificationSettings) HeightPct() int { return c.heightPct }

// SetHeightPct sets the height as a percentage of the screen, used when Height is 0.
func (c *NotificationSettings) SetHeightPct(p int) *NotificationSettings {
	c.heightPct = p
	return c
}

// Padding returns the vertical and horizontal padding inside the border.
func (c *NotificationSettings) Padding() (vertical, horizontal int) {
	return c.paddingV, c.paddingH
}

// SetPadding sets the vertical and horizontal padding inside the border.
func (c *NotificationSettings) SetPadding(vertical, horizontal int) *NotificationSettings {
	c.paddingV, c.paddingH = vertical, horizontal
	return c
}

// Border returns the border drawn around new notifications.
func (c *NotificationSettings) Border() lipgloss.Border { return c.border }

// SetBorder sets the border drawn around new notifications.
func (c *NotificationSettings) SetBorder(b lipgloss.Border) *NotificationSettings {
	c.border = b
	return c
}

// Blocking reports whether new notifications block the page's keys.
func (c *NotificationSettings) Blocking() bool { return c.blocking }

// SetBlocking sets whether new notifications block the page's keys.
func (c *NotificationSettings) SetBlocking(b bool) *NotificationSettings {
	c.blocking = b
	return c
}

// Timeout returns how long a notification stays before it dismisses itself;
// 0 keeps it up until its exit key.
func (c *NotificationSettings) Timeout() time.Duration { return c.timeout }

// SetTimeout sets how long a notification stays before it dismisses itself.
// 0 keeps it up until its exit key.
func (c *NotificationSettings) SetTimeout(t time.Duration) *NotificationSettings {
	c.timeout = t
	return c
}

// ExitKey returns the key that dismisses new notifications.
func (c *NotificationSettings) ExitKey() string { return c.exitKey }

// SetExitKey sets the key that dismisses new notifications.
func (c *NotificationSettings) SetExitKey(k string) *NotificationSettings {
	c.exitKey = keyname.Normalize(k)
	return c
}

// Notification is a lightweight overlay for transient messages. It starts
// from NotificationConfig and styles itself from its level.
type Notification struct {
	*Overlay

	timeout time.Duration     // auto-dismiss delay; 0 keeps it up
	level   NotificationLevel // picks the colors and icon
	message string            // the text without the level's icon prefix
}

// NewNotification creates an Info notification with the given message,
// using NotificationConfig for position, size, timeout and exit key.
func NewNotification(message string) *Notification {
	n := &Notification{
		Overlay: New(),
		level:   LevelInfo,
		message: message,
	}

	cfg := NotificationConfig
	n.Overlay.SetPosition(cfg.position)
	n.Overlay.SetBlocking(cfg.blocking)
	n.Overlay.SetExitKey(cfg.exitKey)
	n.timeout = cfg.timeout

	if cfg.width > 0 {
		n.Overlay.SetWidth(cfg.width)
	} else if cfg.widthPct > 0 {
		n.Overlay.SetWidthPct(cfg.widthPct)
	}

	if cfg.height > 0 {
		n.Overlay.SetHeight(cfg.height)
	} else if cfg.heightPct > 0 {
		n.Overlay.SetHeightPct(cfg.heightPct)
	}

	n.applyLevelStyle()
	return n
}

// SetLevel sets the notification level and restyles the notification from it.
func (n *Notification) SetLevel(level NotificationLevel) *Notification {
	n.level = level
	n.applyLevelStyle()
	return n
}

// SetTimeout sets the auto-dismiss delay. 0 keeps the notification up until
// its exit key.
func (n *Notification) SetTimeout(d time.Duration) *Notification {
	n.timeout = d
	return n
}

// SetContent replaces the message. The level's icon is kept in front of it.
func (n *Notification) SetContent(content string) *Notification {
	n.message = content
	n.applyLevelStyle()
	return n
}

// SetPosition sets where the notification appears.
func (n *Notification) SetPosition(pos Position) *Notification {
	n.Overlay.SetPosition(pos)
	return n
}

// SetSize sets the notification size percentage (both width and height).
func (n *Notification) SetSize(percent int) *Notification {
	n.Overlay.SetSize(percent)
	return n
}

// SetWidthPct sets the notification width percentage.
func (n *Notification) SetWidthPct(percent int) *Notification {
	n.Overlay.SetWidthPct(percent)
	return n
}

// SetHeightPct sets the notification height percentage.
func (n *Notification) SetHeightPct(percent int) *Notification {
	n.Overlay.SetHeightPct(percent)
	return n
}

// SetWidth sets the notification fixed width in characters.
func (n *Notification) SetWidth(chars int) *Notification {
	n.Overlay.SetWidth(chars)
	return n
}

// SetHeight sets the notification fixed height in lines.
func (n *Notification) SetHeight(lines int) *Notification {
	n.Overlay.SetHeight(lines)
	return n
}

// SetExitKey sets the key to dismiss the notification.
func (n *Notification) SetExitKey(key string) *Notification {
	n.Overlay.SetExitKey(key)
	return n
}

// SetBlocking sets whether the notification blocks page keys.
func (n *Notification) SetBlocking(blocking bool) *Notification {
	n.Overlay.SetBlocking(blocking)
	return n
}

// SetStyle replaces the level-based style. Use it to customize borders,
// colors and padding in one go.
func (n *Notification) SetStyle(style lipgloss.Style) *Notification {
	n.Overlay.SetStyle(style)
	return n
}

// SetBorder sets the border style, such as lipgloss.NormalBorder() or
// lipgloss.HiddenBorder().
func (n *Notification) SetBorder(border lipgloss.Border) *Notification {
	n.Overlay.SetBorder(border)
	return n
}

// SetBorderColor sets the border foreground color.
func (n *Notification) SetBorderColor(c color.Color) *Notification {
	n.SetBorderForeground(c)
	return n
}

// SetBorderBackground sets the border background color.
func (n *Notification) SetBorderBackground(c color.Color) *Notification {
	n.Overlay.SetBorderBackground(c)
	return n
}

// SetPadding sets the padding (vertical, horizontal).
func (n *Notification) SetPadding(vertical, horizontal int) *Notification {
	n.Overlay.SetPadding(vertical, horizontal)
	return n
}

// SetBackground sets the background color.
func (n *Notification) SetBackground(c color.Color) *Notification {
	n.Overlay.SetBackground(c)
	return n
}

// SetForeground sets the text color.
func (n *Notification) SetForeground(c color.Color) *Notification {
	n.Overlay.SetForeground(c)
	return n
}

// RegisterKey registers a message to dispatch when key is pressed.
func (n *Notification) RegisterKey(key string, msg tea.Msg) *Notification {
	n.Overlay.RegisterKey(key, msg)
	return n
}

// SetZIndex sets the z-axis stacking order (CSS-style).
// Higher values render on top. Default is 0.
func (n *Notification) SetZIndex(z int) *Notification {
	n.Overlay.SetZIndex(z)
	return n
}

// Timeout returns the auto-dismiss delay.
func (n *Notification) Timeout() time.Duration {
	return n.timeout
}

// Level returns the notification level.
func (n *Notification) Level() NotificationLevel {
	return n.level
}

// Message returns the message without the level's icon prefix.
func (n *Notification) Message() string {
	return n.message
}

// applyLevelStyle sets the content and style from the level's entry in
// NotificationConfig.
func (n *Notification) applyLevelStyle() {
	cfg := NotificationConfig
	ls := cfg.LevelStyle(n.level)

	if ls.Icon != "" {
		n.Overlay.SetContent(ls.Icon + " " + n.message)
	} else {
		n.Overlay.SetContent(n.message)
	}

	n.Overlay.SetStyle(lipgloss.NewStyle().
		Background(ls.Background).
		Foreground(ls.Foreground).
		Border(cfg.border).
		BorderForeground(ls.Border).
		BorderBackground(ls.Background).
		Padding(cfg.paddingV, cfg.paddingH))
}

// Info creates an info-level notification.
func Info(message string) *Notification {
	return NewNotification(message).SetLevel(LevelInfo)
}

// Note creates a note-level notification (highlighted information).
func Note(message string) *Notification {
	return NewNotification(message).SetLevel(LevelNote)
}

// Tip creates a tip-level notification (optional helpful info).
func Tip(message string) *Notification {
	return NewNotification(message).SetLevel(LevelTip)
}

// Important creates an important-level notification (crucial info).
func Important(message string) *Notification {
	return NewNotification(message).SetLevel(LevelImportant)
}

// Success creates a success-level notification.
func Success(message string) *Notification {
	return NewNotification(message).SetLevel(LevelSuccess)
}

// Warning creates a warning-level notification.
func Warning(message string) *Notification {
	return NewNotification(message).SetLevel(LevelWarning)
}

// Caution creates a caution-level notification (negative consequences).
func Caution(message string) *Notification {
	return NewNotification(message).SetLevel(LevelCaution)
}

// Error creates an error-level notification.
func Error(message string) *Notification {
	return NewNotification(message).SetLevel(LevelError)
}
