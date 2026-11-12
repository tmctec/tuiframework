package styles

import (
	"encoding/json"
	"image/color"
	"os"
	"reflect"
	"strings"
	"testing"

	"charm.land/lipgloss/v2"
)

func TestHex(t *testing.T) {
	tests := []struct {
		name string
		c    color.Color
		want string
	}{
		{"hex color", lipgloss.Color("#BD93F9"), "#bd93f9"},
		{"rgba", color.RGBA{R: 1, G: 2, B: 3, A: 0xff}, "#010203"},
		{"nil", nil, ""},
		{"no color", lipgloss.NoColor{}, ""},
		{"unparsable string", lipgloss.Color("nope"), ""},
	}
	for _, tc := range tests {
		if got := Hex(tc.c); got != tc.want {
			t.Errorf("Hex(%s) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestColorOrTreatsNilAndNoColorAsUnset(t *testing.T) {
	fallback := lipgloss.Color("#111111")
	set := lipgloss.Color("#222222")

	if got := colorOr(nil, fallback); got != fallback {
		t.Errorf("colorOr(nil) = %v, want the fallback", got)
	}
	if got := colorOr(lipgloss.NoColor{}, fallback); got != fallback {
		t.Errorf("colorOr(NoColor) = %v, want the fallback", got)
	}
	if got := colorOr(set, fallback); got != set {
		t.Errorf("colorOr(set) = %v, want the color itself", got)
	}
}

// TestBuiltInThemesMatchGolden compares every color of every resolved
// built-in palette with testdata/palettes.golden.json, the values the themes
// had when the file was recorded. A theme added later is not in the file and
// is covered by the registry and contrast tests instead.
func TestBuiltInThemesMatchGolden(t *testing.T) {
	raw, err := os.ReadFile("testdata/palettes.golden.json")
	if err != nil {
		t.Fatal(err)
	}
	var golden map[string]map[string]string
	if err := json.Unmarshal(raw, &golden); err != nil {
		t.Fatal(err)
	}
	if len(golden) != 21 {
		t.Fatalf("golden file holds %d themes, want 21", len(golden))
	}

	colorType := reflect.TypeOf((*color.Color)(nil)).Elem()
	checked := 0
	for _, th := range AllThemes {
		want, ok := golden[th.Name()]
		if !ok {
			continue
		}
		checked++
		v := reflect.ValueOf(th.Palette)
		fields := 0
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Type() != colorType {
				continue
			}
			name := v.Type().Field(i).Name
			hex, recorded := want[name]
			if !recorded {
				continue // a field newer than the golden file
			}
			fields++
			c, _ := v.Field(i).Interface().(color.Color)
			if got := Hex(c); got != strings.ToLower(hex) {
				t.Errorf("%s.%s = %q, want %q", th.Name(), name, got, strings.ToLower(hex))
			}
		}
		if fields != len(want) {
			t.Errorf("%s: compared %d fields, golden file has %d", th.Name(), fields, len(want))
		}
	}
	if checked != len(golden) {
		t.Errorf("compared %d themes, golden file has %d", checked, len(golden))
	}
}
