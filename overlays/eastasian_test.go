package overlays

import (
	"os"
	"strconv"
	"testing"
)

// eastAsianWidths reports whether RUNEWIDTH_EASTASIAN is on. The switch
// makes x/ansi, and with it Lip Gloss and Bubble Tea, measure ambiguous
// characters such as box-drawing lines two cells wide.
func eastAsianWidths() bool {
	on, _ := strconv.ParseBool(os.Getenv("RUNEWIDTH_EASTASIAN"))
	return on
}

// needsNarrowBorders skips a test that asserts the exact size of a bordered
// box. With RUNEWIDTH_EASTASIAN on, Lip Gloss sizes a border as one cell and
// measures it as two, so a bordered box is wider than the width it was given.
func needsNarrowBorders(t *testing.T) {
	t.Helper()
	if eastAsianWidths() {
		t.Skip("RUNEWIDTH_EASTASIAN is on: Lip Gloss draws borders wider than it sizes them")
	}
}
