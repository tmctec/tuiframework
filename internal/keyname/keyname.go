// Package keyname holds the one rule every key binding entry point shares:
// how a key name written by an application maps to the name a key press
// reports.
package keyname

// Normalize returns the name a key press reports for a key name as an
// application may write it. The space bar is "space"; a binding written " "
// means the same key. Every other name is returned unchanged.
func Normalize(name string) string {
	if name == " " {
		return "space"
	}
	return name
}
