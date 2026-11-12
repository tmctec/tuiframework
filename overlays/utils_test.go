package overlays

import "testing"

func TestGetLines(t *testing.T) {
	tests := []struct {
		name       string
		in         string
		wantLines  int
		wantWidest int
	}{
		{"single line", "abc", 1, 3},
		{"empty", "", 1, 0},
		{"widest wins", "ab\ncde\nf", 3, 3},
		{"ANSI codes do not count", "\x1b[31mab\x1b[0m\nc", 2, 2},
		{"wide characters count twice", "日本\nabc", 2, 4},
		{"trailing newline adds an empty line", "ab\n", 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines, widest := GetLines(tt.in)
			if len(lines) != tt.wantLines {
				t.Errorf("len(lines) = %d, want %d", len(lines), tt.wantLines)
			}
			if widest != tt.wantWidest {
				t.Errorf("widest = %d, want %d", widest, tt.wantWidest)
			}
		})
	}
}
