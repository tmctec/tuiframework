package overlays

import "testing"

func TestCalculatePosition(t *testing.T) {
	// Container: 100x50, Overlay: 40x20
	containerW, containerH := 100, 50
	overlayW, overlayH := 40, 20

	tests := []struct {
		name    string
		pos     Position
		expectX int
		expectY int
	}{
		{"TopLeft", PositionTopLeft, 0, 0},
		{"TopCenter", PositionTopCenter, 30, 0},   // (100-40)/2 = 30
		{"TopRight", PositionTopRight, 60, 0},     // 100-40 = 60
		{"CenterLeft", PositionCenterLeft, 0, 15}, // (50-20)/2 = 15
		{"Center", PositionCenter, 30, 15},
		{"CenterRight", PositionCenterRight, 60, 15},
		{"BottomLeft", PositionBottomLeft, 0, 30}, // 50-20 = 30
		{"BottomCenter", PositionBottomCenter, 30, 30},
		{"BottomRight", PositionBottomRight, 60, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := CalculatePosition(tt.pos, overlayW, overlayH, containerW, containerH)
			if x != tt.expectX {
				t.Errorf("x = %d, want %d", x, tt.expectX)
			}
			if y != tt.expectY {
				t.Errorf("y = %d, want %d", y, tt.expectY)
			}
		})
	}
}

func TestCalculatePositionOversized(t *testing.T) {
	// Overlay larger than container should clamp to 0
	x, y := CalculatePosition(PositionCenter, 200, 100, 100, 50)
	if x != 0 {
		t.Errorf("x = %d, want 0 (clamped)", x)
	}
	if y != 0 {
		t.Errorf("y = %d, want 0 (clamped)", y)
	}
}

func TestCalculatePositionMalformedFallsBackToCenter(t *testing.T) {
	wantX, wantY := CalculatePosition(PositionCenter, 40, 20, 100, 50)
	for _, pos := range []Position{"", "C", "LTX", "center"} {
		x, y := CalculatePosition(pos, 40, 20, 100, 50)
		if x != wantX || y != wantY {
			t.Errorf("CalculatePosition(%q) = (%d, %d), want center (%d, %d)", pos, x, y, wantX, wantY)
		}
	}
}

func TestCalculatePositionClampsEachAxis(t *testing.T) {
	// Only the width overflows: x clamps to 0 while y still centers.
	x, y := CalculatePosition(PositionBottomRight, 200, 20, 100, 50)
	if x != 0 || y != 30 {
		t.Errorf("x, y = %d, %d; want 0, 30", x, y)
	}
	// Only the height overflows.
	x, y = CalculatePosition(PositionBottomRight, 40, 100, 100, 50)
	if x != 60 || y != 0 {
		t.Errorf("x, y = %d, %d; want 60, 0", x, y)
	}
}
