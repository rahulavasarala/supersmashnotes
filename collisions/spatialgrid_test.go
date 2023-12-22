package collisions

import (
	"reflect"
	"testing"
)

func TestFindHashPositions(t *testing.T) {
	tests := []struct {
		xpos, ypos      float64
		ewidth          int
		eheight         int
		gridWidth       int
		gridHeight      int
		rows            int
		columns         int
		expectedSquares []Square
		expectedError   error
	}{
		// Test case 1: Add test case scenarios here
		{
			xpos:       110,
			ypos:       110,
			ewidth:     150,
			eheight:    150,
			gridWidth:  500,
			gridHeight: 500,
			rows:       5,
			columns:    5,
			expectedSquares: []Square{
				{row: 1, col: 1},
				{row: 1, col: 2},
				{row: 2, col: 1},
				{row: 2, col: 2},
			},
			expectedError: nil,
		},

		{
			xpos:       201,
			ypos:       201,
			ewidth:     98,
			eheight:    98,
			gridWidth:  500,
			gridHeight: 500,
			rows:       5,
			columns:    5,
			expectedSquares: []Square{
				{row: 2, col: 2},
			},
			expectedError: nil,
		},
		// Test case 2: Add more test case scenarios here
	}

	for _, test := range tests {
		squares, err := findHashPositions(test.xpos, test.ypos, test.ewidth, test.eheight, test.gridWidth, test.gridHeight, test.rows, test.columns)

		if err != test.expectedError {
			t.Errorf("For (%f, %f), expected error '%v' but got '%v'", test.xpos, test.ypos, test.expectedError, err)
		}

		if !reflect.DeepEqual(squares, test.expectedSquares) {
			t.Errorf("For (%f, %f), expected squares %+v but got %+v", test.xpos, test.ypos, test.expectedSquares, squares)
		}

	}
}
