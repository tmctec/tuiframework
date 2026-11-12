package overlays_test

import (
	"fmt"
	"strings"
	"time"

	"github.com/tmctec/tuiframework/overlays"
)

// A notification takes its level's icon and the defaults in NotificationConfig.
func ExampleNotification() {
	n := overlays.Warning("Disk at 90%").
		SetTimeout(0). // stays until its exit key
		SetPosition(overlays.PositionBottomRight)

	fmt.Println(n.Level(), n.Message())
	fmt.Println(n.Position(), n.ExitKey(), n.IsBlocking(), n.Timeout())
	fmt.Println(overlays.Info("x").Timeout() == overlays.NotificationConfig.Timeout())
	// Output:
	// Warning Disk at 90%
	// RB esc false 0s
	// true
}

// Fixed sizes beat percentages, and nothing exceeds the screen.
func ExampleOverlay_CalculateDimensions() {
	o := overlays.New().SetSize(50)
	w, h := o.CalculateDimensions(100, 40)
	fmt.Println(w, h)

	o.SetWidth(30).SetHeightPct(25)
	w, h = o.CalculateDimensions(100, 40)
	fmt.Println(w, h)

	w, h = overlays.New().SetWidth(500).CalculateDimensions(100, 40)
	fmt.Println(w, h)
	// Output:
	// 50 20
	// 30 10
	// 100 32
}

// Every overlay needs an exit key before it can be shown.
func ExampleOverlay_Validate() {
	fmt.Println(overlays.New().Validate())
	fmt.Println(overlays.New().SetExitKey("esc").Validate())
	// Output:
	// overlay must have an exit key defined: call SetExitKey()
	// <nil>
}

// An anchor gives the offset of an overlay inside the screen.
func ExampleCalculatePosition() {
	x, y := overlays.CalculatePosition(overlays.PositionCenter, 20, 10, 80, 24)
	fmt.Println(x, y)
	x, y = overlays.CalculatePosition(overlays.PositionBottomRight, 20, 10, 80, 24)
	fmt.Println(x, y)
	// Output:
	// 30 7
	// 60 14
}

// NotificationConfig changes the defaults for notifications created afterwards.
func ExampleNotificationSettings() {
	prev := *overlays.NotificationConfig
	defer func() { *overlays.NotificationConfig = prev }()

	overlays.NotificationConfig.SetTimeout(5 * time.Second).SetPosition(overlays.PositionTopCenter)
	n := overlays.Tip("Press ? for help")
	fmt.Println(n.Timeout(), n.Position())
	// Output: 5s CT
}

// A canvas composites a frame from layers. Each layer replaces what is under
// its rectangle, and the frame is exactly as wide as the canvas.
func ExampleCanvas() {
	page := []string{"............", "............", "............"}
	popup := []string{"+----+", "| hi |", "+----+"}

	c := overlays.NewCanvas().Reset(12, 3, nil)
	c.Draw(page, 0, 0, 12, 3, nil)
	c.Draw(popup, 3, 0, 6, 3, nil)

	fmt.Println(strings.ReplaceAll(c.String(), "\x1b[0m", ""))
	// Output:
	// ...+----+...
	// ...| hi |...
	// ...+----+...
}
