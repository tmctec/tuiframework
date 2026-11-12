package keyname

import "testing"

func TestNormalize(t *testing.T) {
	for in, want := range map[string]string{
		" ": "space", "space": "space", "s": "s", "": "", "shift+tab": "shift+tab", "  ": "  ",
	} {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}
