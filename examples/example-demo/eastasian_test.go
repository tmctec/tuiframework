package main

import (
	"os"
	"strconv"
	"testing"
)

// needsNarrowBorders skips a test that counts on the exact size of a
// bordered box. With RUNEWIDTH_EASTASIAN on, Lip Gloss draws a border two
// cells wide and sizes it as one, so a bordered popup is wider than asked.
func needsNarrowBorders(t *testing.T) {
	t.Helper()
	if on, _ := strconv.ParseBool(os.Getenv("RUNEWIDTH_EASTASIAN")); on {
		t.Skip("RUNEWIDTH_EASTASIAN is on: Lip Gloss draws borders wider than it sizes them")
	}
}
