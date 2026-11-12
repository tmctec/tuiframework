package tabapp

import (
	"os"
	"strconv"
	"testing"
)

// needsNarrowBorders skips a test that asserts the exact size of a bordered
// box. With RUNEWIDTH_EASTASIAN on, x/ansi measures box-drawing characters
// two cells wide while Lip Gloss sizes a border as one, so a bordered box is
// wider than the width it was given.
func needsNarrowBorders(t *testing.T) {
	t.Helper()
	if on, _ := strconv.ParseBool(os.Getenv("RUNEWIDTH_EASTASIAN")); on {
		t.Skip("RUNEWIDTH_EASTASIAN is on: Lip Gloss draws borders wider than it sizes them")
	}
}
