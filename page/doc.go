// Package page defines the pages a tabapp.TabApp shows. A Page has a name
// for the tab, a content function that builds a flexbox from the current
// Model on every render, key bindings that dispatch messages, and footer
// labels for its keys.
//
//	page.New("Orders").
//	    SetContent(func(fb *flexbox.FlexBox, model interface{}, navBarHeight, footerHeight int) {
//	        m := model.(*tabapp.Model)
//	        // add rows and cells to fb from m
//	    }).
//	    BindKey("r", "Refresh", RefreshOrdersMsg{})
//
// The content function receives the model as interface{} because package
// tabapp imports this package; assert it to *tabapp.Model.
//
// A key is named as a key press reports it: "r", "space", "enter",
// "shift+tab". A key written " " is the space bar.
package page
