package collisions

import (
	"reflect"
	"testing"
)

func TestFindHashPositions(t *testing.T) {
	tests := []struct {
		xpos, ypos            float64
		charWidth, charHeight float64
		gridWidth, gridHeight int
		xaxis                 int
		yaxis                 int
		expectedSquares       []Square
		expectedError         error
	}{
		// Test case 1: Add test case scenarios here
		{
			xpos:       50,
			ypos:       50,
			charWidth:  20,
			charHeight: 20,
			gridWidth:  500,
			gridHeight: 500,
			xaxis:      5,
			yaxis:      5,
			expectedSquares: []Square{
				{xcoord: 0, ycoord: 0},
			},
			expectedError: nil,
		},
		{
			xpos:       99,
			ypos:       99,
			charWidth:  20,
			charHeight: 20,
			gridWidth:  500,
			gridHeight: 500,
			xaxis:      5,
			yaxis:      5,
			expectedSquares: []Square{
				{xcoord: 0, ycoord: 0},
				{xcoord: 0, ycoord: 1},
				{xcoord: 1, ycoord: 0},
				{xcoord: 1, ycoord: 1},
			},
			expectedError: nil,
		},
		{
			xpos:       -5,
			ypos:       150,
			charWidth:  20,
			charHeight: 20,
			gridWidth:  500,
			gridHeight: 500,
			xaxis:      5,
			yaxis:      5,
			expectedSquares: []Square{
				{xcoord: 0, ycoord: 1},
			},
			expectedError: nil,
		},
		{
			xpos:            0,
			ypos:            -50,
			charWidth:       20,
			charHeight:      20,
			gridWidth:       500,
			gridHeight:      500,
			xaxis:           5,
			yaxis:           5,
			expectedSquares: []Square{},
			expectedError:   nil,
		},
	}

	for _, test := range tests {
		squares, err := findHashPositions(test.xpos, test.ypos, test.charWidth, test.charHeight, test.gridWidth, test.gridHeight, test.xaxis, test.yaxis)

		if err != test.expectedError {
			t.Errorf("For (%f, %f), expected error '%v' but got '%v'", test.xpos, test.ypos, test.expectedError, err)
		}

		if !reflect.DeepEqual(squares, test.expectedSquares) {
			t.Errorf("For (%f, %f), expected squares %+v but got %+v", test.xpos, test.ypos, test.expectedSquares, squares)
		}

	}
}
