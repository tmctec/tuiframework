package overlays

import "strings"

// GetLines splits a block of rendered text into lines and returns the width
// in cells of the widest one, escape sequences excluded.
func GetLines(s string) (lines []string, widest int) {
	lines = strings.Split(s, "\n")
	for _, l := range lines {
		widest = max(widest, lineWidth(l))
	}
	return lines, widest
}
