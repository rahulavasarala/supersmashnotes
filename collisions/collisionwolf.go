package collisions

// type CollisionWolf struct {
// }

// func (s *CollisionWolf) FindPotentialCollisions(gs GridSquare) [][]Thing {

// 	collisionMap := [][]Thing{}

// 	for i := 0; i < len(gs.Characters); i++ {
// 		collisionMap = append(collisionMap, []Thing{})
// 	}

// 	for i, _ := range gs.Characters {
// 		for j := i + 1; j < len(gs.Characters); j++ {
// 			if HashBBIntersec(gs.Characters[i], gs.Characters[j]) {
// 				collisionMap[i] = append(collisionMap[i], gs.Characters[j])
// 				collisionMap[j] = append(collisionMap[j], gs.Characters[i])
// 			}
// 		}
// 	}

// 	for i, _ := range gs.Characters {
// 		for j, _ := range gs.Immovables {
// 			if HashBBIntersec(gs.Characters[i], gs.Immovables[j]) {
// 				collisionMap[i] = append(collisionMap[i], gs.Immovables[j])
// 			}
// 		}
// 	}

// 	return collisionMap
// }

// func HashBBIntersec(e1 Thing, e2 Thing) bool {

// 	xpos, ypos := e1.GetPos()
// 	bbw, bbh := e1.GetBoundingBox()

// 	xpos2, ypos2 := e2.GetPos()
// 	bbw2, bbh2 := e2.GetBoundingBox()

// 	if int(xpos) > int(xpos2)+bbw2 || int(xpos2) > int(xpos)+bbw || int(ypos)-bbh > int(ypos2) || int(ypos2)-bbh2 > int(ypos) {
// 		return false
// 	}

// 	return true
// }

// func HurtBoxHitBoxIntersect(hurt HurtBox, hit HitBox) bool {
// 	if int(hurt.xpos) > int(hit.xpos)+hit.width || int(hit.xpos) > int(hurt.xpos)+hurt.width || int(hurt.ypos)-hurt.height > int(hit.ypos) || int(hit.ypos)-hit.height > int(hurt.ypos) {
// 		return false
// 	}
// 	return true
// }

// //this is the part where we get into the collision handler algorithms
