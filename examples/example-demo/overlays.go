package main

import (
	"fmt"
	"image/color"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/overlays"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/styles"
	"github.com/tmctec/tuiframework/tabapp"
)

// demoModalSpec describes one of the toggle modals. Border picks the border
// color from the active palette at the time the key is pressed, so modals
// follow the theme.
type demoModalSpec struct {
	Title   string
	Content string
	ZIndex  int
	SizePct int
	Border  func(p styles.Palette) color.Color
}

var demoModals = map[string]demoModalSpec{
	"A": {
		Title:   "Modal A",
		Content: "z-index 5, 70% of the screen.\n\nPress a again to close it, f to bring it to the front,\nor Esc to close whichever overlay is on top.",
		ZIndex:  5, SizePct: 70,
		Border: func(p styles.Palette) color.Color { return p.Primary },
	},
	"B": {
		Title:   "Modal B",
		Content: "z-index 15, 50% of the screen.\n\nIt renders above Modal A.",
		ZIndex:  15, SizePct: 50,
		Border: func(p styles.Palette) color.Color { return p.Secondary },
	},
	"C": {
		Title:   "Modal C",
		Content: "z-index 25, 35% of the screen.\n\nIt renders above A and B.",
		ZIndex:  25, SizePct: 35,
		Border: func(p styles.Palette) color.Color { return p.Success },
	},
}

// modalOrder lists the toggle modals in display order.
var modalOrder = []string{"A", "B", "C"}

func (s demoModalSpec) toggleMsg(p styles.Palette) tabapp.ToggleModalMsg {
	return tabapp.ToggleModalMsg{
		Title:       s.Title,
		Content:     s.Content,
		ZIndex:      s.ZIndex,
		SizePct:     s.SizePct,
		BorderColor: s.Border(p),
	}
}

// stackMsgs toggles three modals at z 1, 2 and 3: the first press opens
// them, and a press while they are open closes them.
func stackMsgs(p styles.Palette) []tea.Msg {
	return []tea.Msg{
		tabapp.ToggleModalMsg{Title: "Stack 1", Content: "z-index 1, 70%. Esc closes the top one first.", ZIndex: 1, SizePct: 70, BorderColor: p.Fg3},
		tabapp.ToggleModalMsg{Title: "Stack 2", Content: "z-index 2, 50%.", ZIndex: 2, SizePct: 50, BorderColor: p.Warning},
		tabapp.ToggleModalMsg{Title: "Stack 3", Content: "z-index 3, 30%.", ZIndex: 3, SizePct: 30, BorderColor: p.Error},
	}
}

// sampleNotifications is the text CycleNotificationMsg shows for each level.
var sampleNotifications = [...]string{
	overlays.LevelInfo:      "This is an info notification",
	overlays.LevelNote:      "This is a note notification",
	overlays.LevelTip:       "This is a tip notification",
	overlays.LevelImportant: "This is an important notification",
	overlays.LevelSuccess:   "Operation completed successfully!",
	overlays.LevelWarning:   "Check your input carefully",
	overlays.LevelCaution:   "This action may have consequences",
	overlays.LevelError:     "Something went wrong!",
}

// openModals reports which toggle modals are open, from Model.OpenModals.
func openModals(m *tabapp.Model) map[string]bool {
	open := make(map[string]bool, len(demoModals))
	for id, spec := range demoModals {
		_, open[id] = m.OpenModals[spec.Title]
	}
	return open
}

const overlaysHelp = `Overlays

i / w / e   info, warning, error notification
            at z-index 10 / 20 / 30
a / b / c   toggle Modal A (z 5), B (z 15), C (z 25)
z           stack three modals at z 1, 2, 3; again closes them
f           bring Modal A to the front
x           clear every overlay
Esc         close the top overlay

Higher z-index renders on top. Colors come
from the active theme when the modal opens.`

const overlaysKeys = `i / w / e   notification at z 10 / 20 / 30
a / b / c   toggle Modal A (z 5, 70%), B (z 15, 50%), C (z 25, 35%)
z           stack three modals at z 1 / 2 / 3, or close them
f           bring Modal A to the front
x           clear all overlays`

// buildOverlaysPage demonstrates z-index stacking of notifications and modals.
func buildOverlaysPage() *page.Page {
	p := page.New("Overlays").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			theme := m.Theme
			themed(fb, theme)

			title := fixedRow(fb, centered(theme.TextHeading), 1, "Overlay stacking")

			open := openModals(m)
			parts := make([]string, 0, len(modalOrder))
			for _, id := range modalOrder {
				label := "Modal " + id + " closed"
				style := theme.MutedStyle
				if open[id] {
					label = "Modal " + id + " open"
					style = theme.StatusSuccess
				}
				parts = append(parts, style.Render(label))
			}
			status := fixedRow(fb, centered(theme.BaseStyle), 3,
				theme.ContentStyle.Render(fmt.Sprintf("Overlays on screen: %d", m.OverlayCount))+"\n"+
					strings.Join(parts, theme.MutedStyle.Render("  ·  ")))

			keys := fixedRow(fb, card(theme).Align(lipgloss.Left), 7, overlaysKeys)

			explain := fixedRow(fb, centered(theme.MutedStyle), 2,
				"Higher z-index renders on top, like CSS. Esc closes only the top overlay.\nModal colors are chosen from the active theme when they open.")

			fb.AddRows([]*flexbox.Row{title, status, keys, spacer(fb), explain, filler(fb)})
		}).
		BindKey("i", "", tabapp.ShowNotificationMsg{Level: overlays.LevelInfo, Message: "Info notification at z-index 10", ZIndex: 10}).
		BindKey("w", "", tabapp.ShowNotificationMsg{Level: overlays.LevelWarning, Message: "Warning notification at z-index 20", ZIndex: 20}).
		BindKey("e", "", tabapp.ShowNotificationMsg{Level: overlays.LevelError, Message: "Error notification at z-index 30", ZIndex: 30}).
		BindKey("a", "", ToggleDemoModalMsg{ID: "A"}).
		BindKey("b", "", ToggleDemoModalMsg{ID: "B"}).
		BindKey("c", "", ToggleDemoModalMsg{ID: "C"}).
		BindKey("z", "Stack", StackDemoModalsMsg{}).
		BindKey("f", "Front A", FrontDemoModalMsg{ID: "A"}).
		BindKey("x", "Clear", ClearOverlaysMsg{}).
		SetFooter(map[string]string{
			"i/w/e": "Notify",
			"a/b/c": "Modals",
			"z":     "Stack",
			"f":     "Front A",
			"x":     "Clear",
		})
	return withHelp(p, overlaysHelp)
}
