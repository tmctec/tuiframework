package page_test

import (
	"fmt"

	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
)

// RefreshMsg is the message the r key dispatches.
type RefreshMsg struct{}

// A page has a name, a content function that builds a flexbox from the
// model, key bindings that dispatch messages, and footer labels.
func ExampleNew() {
	p := page.New("Orders").
		SetContent(func(fb *flexbox.FlexBox, _ interface{}, _, _ int) {
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetContent("3 open orders")
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		}).
		BindKey("r", "Refresh", RefreshMsg{}).
		BindKey("up", "", RefreshMsg{}). // an empty label keeps the key out of the footer
		SetFooter(map[string]string{"↑/↓": "Select"})

	fmt.Println(p.Name())
	fmt.Println(p.HasKey("r"), p.HasKey("x"))
	fmt.Println(p.Bindings()["r"].Label())
	fmt.Println(p.Footer()["↑/↓"])

	fb := flexbox.New(20, 1)
	p.Content()(fb, nil, 0, 0)
	fmt.Println(fb.Render())
	// Output:
	// Orders
	// true false
	// Refresh
	// Select
	// 3 open orders
}
