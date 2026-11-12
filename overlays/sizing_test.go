package overlays

import "testing"

func TestOverlayPercentSettersClampAndClearFixed(t *testing.T) {
	o := New().SetWidth(30).SetHeight(10)

	o.SetWidthPct(0)
	if o.WidthPct() != 1 || o.WidthFixed() != 0 {
		t.Errorf("after SetWidthPct(0): pct %d fixed %d, want 1 and 0", o.WidthPct(), o.WidthFixed())
	}
	o.SetWidthPct(250)
	if o.WidthPct() != 100 {
		t.Errorf("after SetWidthPct(250): pct %d, want 100", o.WidthPct())
	}

	o.SetHeightPct(-3)
	if o.HeightPct() != 1 || o.HeightFixed() != 0 {
		t.Errorf("after SetHeightPct(-3): pct %d fixed %d, want 1 and 0", o.HeightPct(), o.HeightFixed())
	}
	o.SetHeightPct(101)
	if o.HeightPct() != 100 {
		t.Errorf("after SetHeightPct(101): pct %d, want 100", o.HeightPct())
	}
}

func TestOverlayFixedSettersClampAndClearPercent(t *testing.T) {
	o := New().SetWidthPct(50).SetHeightPct(40)

	o.SetWidth(0)
	if o.WidthFixed() != 1 || o.WidthPct() != 0 {
		t.Errorf("after SetWidth(0): fixed %d pct %d, want 1 and 0", o.WidthFixed(), o.WidthPct())
	}
	o.SetHeight(-7)
	if o.HeightFixed() != 1 || o.HeightPct() != 0 {
		t.Errorf("after SetHeight(-7): fixed %d pct %d, want 1 and 0", o.HeightFixed(), o.HeightPct())
	}

	o.SetWidth(42).SetHeight(7)
	if o.WidthFixed() != 42 || o.HeightFixed() != 7 {
		t.Errorf("fixed = %dx%d, want 42x7", o.WidthFixed(), o.HeightFixed())
	}
}

func TestOverlaySetSizeClearsPercentagesOnly(t *testing.T) {
	o := New().SetWidthPct(50).SetHeightPct(40).SetSize(70)
	if o.Size() != 70 || o.WidthPct() != 0 || o.HeightPct() != 0 {
		t.Errorf("size %d pct %dx%d, want 70 and 0x0", o.Size(), o.WidthPct(), o.HeightPct())
	}

	fixed := New().SetWidth(30).SetSize(50)
	if fixed.WidthFixed() != 30 {
		t.Errorf("SetSize cleared the fixed width: %d, want 30", fixed.WidthFixed())
	}
}

func TestOverlayCalculateDimensionsPrecedence(t *testing.T) {
	tests := []struct {
		name  string
		o     *Overlay
		wantW int
		wantH int
	}{
		{"size alone", New(), 80, 40},
		{"size 100 fills the container", New().SetSize(100), 100, 50},
		{"percentages beat size", New().SetWidthPct(50).SetHeightPct(20), 50, 10},
		{"fixed width beats percentage", New().SetWidthPct(50).SetWidth(30).SetHeightPct(20), 30, 10},
		{"fixed height beats percentage", New().SetWidthPct(50).SetHeightPct(20).SetHeight(7), 50, 7},
		{"fixed beats size", New().SetWidth(30).SetHeight(7).SetSize(50), 30, 7},
		{"fixed is clamped to the container", New().SetWidth(500).SetHeight(500), 100, 50},
		{"mixed axes", New().SetWidth(25).SetHeightPct(50), 25, 25},
		{"percentage rounds down", New().SetWidthPct(33).SetHeightPct(33), 33, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := tt.o.CalculateDimensions(100, 50)
			if w != tt.wantW || h != tt.wantH {
				t.Errorf("CalculateDimensions = %dx%d, want %dx%d", w, h, tt.wantW, tt.wantH)
			}
		})
	}
}

func TestOverlayCalculateDimensionsNeverExceedsContainer(t *testing.T) {
	overlays := []*Overlay{
		New(),
		New().SetSize(100),
		New().SetWidthPct(100).SetHeightPct(100),
		New().SetWidth(1000).SetHeight(1000),
		New().SetWidth(1000).SetHeightPct(100),
	}
	for _, size := range [][2]int{{80, 24}, {1, 1}, {200, 3}} {
		for i, o := range overlays {
			w, h := o.CalculateDimensions(size[0], size[1])
			if w > size[0] || h > size[1] {
				t.Errorf("overlay %d in %dx%d: %dx%d exceeds the container", i, size[0], size[1], w, h)
			}
		}
	}
}

func TestOverlayEffectiveWidthAndHeight(t *testing.T) {
	o := New()
	if o.EffectiveWidth() != 80 || o.EffectiveHeight() != 80 {
		t.Errorf("effective = %dx%d, want the size 80x80", o.EffectiveWidth(), o.EffectiveHeight())
	}

	o.SetWidthPct(30).SetHeightPct(40)
	if o.EffectiveWidth() != 30 || o.EffectiveHeight() != 40 {
		t.Errorf("effective = %dx%d, want 30x40", o.EffectiveWidth(), o.EffectiveHeight())
	}

	// A fixed size clears the percentage, so the size applies again.
	o.SetWidth(10).SetHeight(5)
	if o.EffectiveWidth() != 80 || o.EffectiveHeight() != 80 {
		t.Errorf("effective = %dx%d after fixed sizes, want 80x80", o.EffectiveWidth(), o.EffectiveHeight())
	}
}
