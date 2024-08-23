package collisions

// import "fmt"

// type CollisionRefiner interface {
// 	RefineCollisions(gs GridSquare, collisionMap [][]Thing) [][]Thing
// }

// type SimpleCollisionRefiner struct {
// }

// // This method will take the existing collision map, and basically check to see if the collisions existing
// //algorithm: you have a list of characters in the grid square, and you have a collision map, which says that
// //this character is colliding with the following things. If the character is colliding with a character
// //store the hitboxes that intersect with the hurt boxes of the character, if the thing is a wall(thing), then just
// //create an artificial hitbox and put it on the player

// type ConnectedHitboxes struct {
// 	hitBoxMap map[HitBox]bool
// }

// func (s *SimpleCollisionRefiner) RefineCollisions(gs GridSquare, collisionMap [][]Thing) ([][]ConnectedHitboxes, error) {

// 	if len(collisionMap) == 0 || len(collisionMap[0]) == 0 {
// 		return nil, fmt.Errorf("no collisions to refine in the first place")
// 	}

// 	refinedCollisionMap := [][]ConnectedHitboxes{}
// 	for i := 0; i < len(collisionMap); i++ {
// 		refinedCollisionMap = append(refinedCollisionMap, []ConnectedHitboxes{})
// 	}

// 	for i, player := range gs.Characters {
// 		hurtboxes := player.GetHurtbox()

// 		for _, thing := range collisionMap[i] {
// 			if _, ok := thing.(Character); ok { //thing is a Character
// 				hitboxes := thing.(Character).GetHitbox()
// 				conn := ConnectedHitboxes{}
// 				conn.hitBoxMap = map[HitBox]bool{}

// 				for i := 0; i < len(hurtboxes); i++ {
// 					for j := 0; j < len(hitboxes); j++ {
// 						if HurtBoxHitBoxIntersect(hurtboxes[i], hitboxes[j]) {
// 							conn.hitBoxMap[hitboxes[j]] = true
// 						}
// 					}
// 				}

// 				refinedCollisionMap[i] = append(refinedCollisionMap[i], conn)

// 			} else {
// 				xpos, ypos := thing.GetPos()
// 				width, height := thing.GetBoundingBox()

// 				conn := ConnectedHitboxes{}
// 				conn.hitBoxMap = map[HitBox]bool{}
// 				conn.hitBoxMap[HitBox{xpos: xpos, ypos: ypos, width: width, height: height}] = true

// 				refinedCollisionMap[i] = append(refinedCollisionMap[i],
// 					conn)

// 			}
// 		}
// 	}

// 	return refinedCollisionMap, nil
// }
