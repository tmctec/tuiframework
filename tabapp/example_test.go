package tabapp_test

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/tmctec/stickers/v2/flexbox"
	"github.com/tmctec/tuiframework/page"
	"github.com/tmctec/tuiframework/tabapp"
)

// CounterMsg is an application message; the custom handler applies it.
type CounterMsg struct{ Delta int }

func counterPage() *page.Page {
	return page.New("Counter").
		SetContent(func(fb *flexbox.FlexBox, model interface{}, _, _ int) {
			m := model.(*tabapp.Model)
			row := fb.NewRow()
			cell := flexbox.NewCell(1, 1)
			cell.SetContent(fmt.Sprintf("count: %d", m.GetCustomInt("count")))
			row.AddCells(cell)
			fb.AddRows([]*flexbox.Row{row})
		}).
		BindKey("+", "Add", CounterMsg{Delta: 1})
}

// Pages are grouped into sections; the nav bar reflects the tree.
func Example() {
	app := tabapp.New()
	phones := tabapp.NewSection("Phones").AddPages(page.New("Browse"), page.New("Deals"))
	app.MustAddSection(tabapp.NewSection("Catalog").AddSection(phones).AddPage(page.New("Overview")))
	app.MustAddPage(page.New("Settings"))

	fmt.Println("depth:", app.Root().Depth())
	fmt.Println("start:", app.NavPath())

	// Jump to Catalog > Phones > Deals: index 0 (Catalog), 0 (Phones), 1 (Deals).
	if err := app.NavigateToPath([]int{0, 0, 1}); err != nil {
		fmt.Println(err)
	}
	fmt.Println("after jump:", app.NavPath(), "page", app.Model().ActivePage)
	// Output:
	// depth: 3
	// start: [0 0 0]
	// after jump: [0 0 1] page 1
}

// A key bound on a page dispatches its message; the custom handler keeps
// application state in Model.Custom and the page renders from it.
func ExampleTabApp_SetCustomMessageHandler() {
	app := tabapp.New()
	app.MustAddPage(counterPage())
	app.SetCustomMessageHandler(func(msg tea.Msg) (bool, tea.Cmd) {
		if m, ok := msg.(CounterMsg); ok {
			app.Model().SetCustom("count", app.Model().GetCustomInt("count")+m.Delta)
			return true, nil
		}
		return false, nil
	})

	app.Update(tabapp.KeyPress("+"))
	app.Update(tabapp.KeyPress("+"))
	fmt.Println(app.Model().GetCustomInt("count"))
	// Output: 2
}

// Tab starts navigation and moves the highlight; the Model reports it.
func ExampleTabApp_Update() {
	app := tabapp.New()
	app.MustAddPage(page.New("One")).MustAddPage(page.New("Two"))
	app.Update(tea.WindowSizeMsg{Width: 60, Height: 20})

	app.Update(tabapp.KeyPress("tab"))
	fmt.Println("navigating:", app.Model().TabNavActive, "active page:", app.Model().ActivePage)

	app.Update(tabapp.KeyPress("esc"))
	fmt.Println("navigating:", app.Model().TabNavActive)
	// Output:
	// navigating: true active page: 1
	// navigating: false
}

// Settings are written by setters and read back from the Model.
func ExampleTabApp_SetNavBarMode() {
	app := tabapp.New().SetNavBarMode(tabapp.NavBarColumns).SetNavBarAutoHide(false)
	fmt.Println(app.Model().NavBarMode, app.Model().NavBarAutoHide)
	// Output: Columns false
}

// A frame sets the header or footer off from the content: a border around
// it, or a solid fill in the theme's fill colors.
func ExampleTabApp_SetHeaderFrame() {
	app := tabapp.New().
		SetHeaderHeight(1).
		SetHeaderFrame(tabapp.FrameFill).
		SetFooterHeight(2).
		SetFooterFrame(tabapp.FrameDouble)

	fmt.Println(app.HeaderFrame(), app.FooterFrame())
	fmt.Println(app.ShowHeaderBorder(), app.ShowFooterBorder())

	// Messages change the bands from key bindings.
	app.Update(tabapp.SetFooterFrameMsg{Frame: tabapp.FrameNone})
	app.Update(tabapp.SetFooterHeightMsg{Height: 0}) // 0 hides the footer
	fmt.Println(app.Model().FooterFrame, app.Model().FooterHeight)
	// Output:
	// Fill Double
	// false true
	// None 0
}
