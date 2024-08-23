package collisions

import (
	"testing"
)

func TestInDiamond(t *testing.T) {
	tests := []struct {
		x, y     float64
		dx, dy   float64
		dw, dh   float64
		expected bool
	}{
		// Test case 1: Add test case scenarios here
		{
			x:        51,
			y:        50,
			dx:       60,
			dy:       60,
			dw:       40,
			dh:       40,
			expected: true,
		},
		{
			x:        25,
			y:        25,
			dx:       60,
			dy:       60,
			dw:       40,
			dh:       40,
			expected: false,
		},
		{
			x:        64,
			y:        33,
			dx:       60,
			dy:       60,
			dw:       60,
			dh:       60,
			expected: false,
		},
		{
			x:        63,
			y:        34,
			dx:       60,
			dy:       60,
			dw:       60,
			dh:       60,
			expected: true,
		},
	}

	for _, test := range tests {
		result := inDiamond(test.x, test.y, test.dx, test.dy, test.dw, test.dh)

		if result != test.expected {
			t.Errorf("For point (%v, %v) with diamond (%v, %v, %v, %v), expected %v but got %v", test.x, test.y, test.dx, test.dy, test.dw, test.dh, test.expected, result)
		}

	}
}

func TestInRect(t *testing.T) {
	tests := []struct {
		x, y     float64
		rx, ry   float64
		rw, rh   float64
		expected bool
	}{
		// Test case 1: Add test case scenarios here
		{
			x:        50,
			y:        50,
			rx:       30,
			ry:       30,
			rw:       40,
			rh:       40,
			expected: true,
		},
		{
			x:        25,
			y:        25,
			rx:       60,
			ry:       60,
			rw:       40,
			rh:       40,
			expected: false,
		},
	}

	for _, test := range tests {
		result := inRect(test.x, test.y, test.rx, test.ry, test.rw, test.rh)

		if result != test.expected {
			t.Errorf("For point (%v, %v) with diamond (%v, %v, %v, %v), expected %v but got %v", test.x, test.y, test.rx, test.ry, test.rw, test.rh, test.expected, result)
		}

	}
}
