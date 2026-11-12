package tabapp

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/tuiframework/overlays"
)

// benchApp returns a sized mixed-depth app with the nav bar shown.
func benchApp(pos NavBarPosition) *TabApp {
	app := rootMixApp()
	app.SetHeader(testHeader).SetNavBarAutoHide(false).SetNavBarPosition(pos)
	app.Update(tea.WindowSizeMsg{Width: 160, Height: 48})
	return app
}

func BenchmarkViewTopBar(b *testing.B) {
	app := benchApp(NavBarTop)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = app.Render()
	}
}

func BenchmarkViewSideBar(b *testing.B) {
	app := benchApp(NavBarLeft)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = app.Render()
	}
}

func BenchmarkViewWithOverlays(b *testing.B) {
	app := benchApp(NavBarTop)
	app.MustShowOverlay(overlays.NewModal().SetTitle("One").SetExitKey("esc").SetZIndex(1).Overlay)
	app.MustShowOverlay(overlays.NewModal().SetTitle("Two").SetExitKey("esc").SetSize(50).SetZIndex(2).Overlay)
	app.MustShowOverlay(overlays.Info("hello").Overlay)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = app.Render()
	}
}

func BenchmarkKeyPress(b *testing.B) {
	app := benchApp(NavBarTop)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		app.Update(KeyPress("tab"))
		app.Update(KeyPress("right"))
	}
}
