package collisions

// import (
// 	"reflect"
// 	"testing"
// )

// func NewConnectedHitBoxes(hitList []HitBox) ConnectedHitboxes {
// 	conn := ConnectedHitboxes{}
// 	conn.hitBoxMap = map[HitBox]bool{}

// 	for _, hitbox := range hitList {
// 		conn.hitBoxMap[hitbox] = true
// 	}

// 	return conn
// }

// func TestRefineCollisions(t *testing.T) {
// 	tests := []struct {
// 		characterList       []Character
// 		wallList            []Thing
// 		expectedRefinedList [][]ConnectedHitboxes
// 	}{
// 		// Test case 1: Add test case scenarios here
// 		{
// 			characterList: []Character{NewCharacter(50, 50, 30, 30), NewCharacter(60, 60, 30, 30)},
// 			wallList:      []Thing{NewWall(40, 40, 100, 20)},
// 			expectedRefinedList: [][]ConnectedHitboxes{
// 				{
// 					NewConnectedHitBoxes([]HitBox{
// 						{60, 60, 30, 30},
// 					}),
// 					NewConnectedHitBoxes([]HitBox{
// 						{40, 40, 100, 20},
// 					}),
// 				},
// 				{
// 					NewConnectedHitBoxes([]HitBox{
// 						{50, 50, 30, 30},
// 					}),
// 					NewConnectedHitBoxes([]HitBox{
// 						{40, 40, 100, 20},
// 					}),
// 				},
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		gridSquare := GridSquare{}

// 		gridSquare.Characters = test.characterList
// 		gridSquare.Immovables = test.wallList

// 		collisionWolf := CollisionWolf{}
// 		collisionList := collisionWolf.FindPotentialCollisions(gridSquare)

// 		simpleCollisionRefiner := SimpleCollisionRefiner{}

// 		refinedCollisionList, err := simpleCollisionRefiner.RefineCollisions(gridSquare, collisionList)

// 		if err != nil {
// 			t.Error("Failed")
// 		}

// 		if !reflect.DeepEqual(test.expectedRefinedList, refinedCollisionList) {
// 			t.Errorf("For gridSquare %v, expected refined collision list %v but got %v",
// 				gridSquare, test.expectedRefinedList, refinedCollisionList)
// 		}
// 	}
// }
