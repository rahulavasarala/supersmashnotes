package collisions

// import (
// 	"reflect"
// 	"testing"
// )

// func NewCharacter(xpos float64, ypos float64, ewidth int, eheight int) Character {
// 	charac := DummyEntity{}

// 	charac.InitCharacter(xpos, ypos, ewidth, eheight)

// 	return &charac
// }

// func NewWall(xpos float64, ypos float64, ewidth int, eheight int) Thing {
// 	wall := DummyEntity{}
// 	wall.InitWall(xpos, ypos, ewidth, eheight)

// 	return &wall
// }

// func TestFindPotentialCollisions(t *testing.T) {
// 	tests := []struct {
// 		characterList         []Character
// 		wallList              []Thing
// 		expectedCollisionList [][]Thing
// 	}{
// 		// Test case 1: Add test case scenarios here
// 		{
// 			characterList:         []Character{NewCharacter(10, 10, 30, 30), NewCharacter(100, 100, 30, 30)},
// 			wallList:              []Thing{NewWall(0, 20, 100, 20)},
// 			expectedCollisionList: [][]Thing{{NewWall(0, 20, 100, 20)}, {}},
// 		},
// 	}

// 	for _, test := range tests {
// 		gridSquare := GridSquare{}

// 		gridSquare.Characters = test.characterList
// 		gridSquare.Immovables = test.wallList

// 		collisionWolf := CollisionWolf{}
// 		collisionList := collisionWolf.FindPotentialCollisions(gridSquare)

// 		if !reflect.DeepEqual(test.expectedCollisionList, collisionList) {
// 			t.Errorf("For gridSquare %v, expected collision list %v but got %v",
// 				gridSquare, test.expectedCollisionList, collisionList)
// 		}
// 	}
// }
