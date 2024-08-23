package collisions

type EcbCollisionFinder struct {
}

func (s *EcbCollisionFinder) FindPotentialCollisions(sg SpatialGrid) map[string]map[string]bool {

	globalCollisionMap := map[string]map[string]bool{}

	for i := 0; i < sg.xaxis; i++ {
		for j := 0; j < sg.yaxis; j++ {
			//for each of these grid squares, let us find all the collisions between characters and walls

			gridSquare := sg.grid[i][j]

			localCollisionMap := map[string]map[string]bool{}

			for _, char := range gridSquare.Characters {
				for _, wall := range gridSquare.Immovables {
					if bbCollision(char, wall) {
						if _, exists := localCollisionMap[char.GetId()]; !exists {
							localCollisionMap[char.GetId()] = map[string]bool{}
						}
						localCollisionMap[char.GetId()][wall.GetId()] = true
					}
				}
			}

			for id, _ := range localCollisionMap {
				if _, exists := globalCollisionMap[id]; !exists {
					globalCollisionMap[id] = map[string]bool{}
				}

				for collisionId, _ := range localCollisionMap[id] {
					globalCollisionMap[id][collisionId] = true
				}
			}

		}
	}

	return globalCollisionMap

}

//The standard format to put entities in to find collisions will be top left point  + width and height

func bbCollision(e1 Thing, e2 Thing) bool {
	e1x, e1y := e1.GetPos()
	e1w, e1h := e1.GetBoundingBox()

	char := ConvertToCharacter(e1)

	if char != nil {
		e1x = e1x - e1w/2
		e1y = e1y - e1h/2
	}

	e2x, e2y := e2.GetPos()
	e2w, e2h := e2.GetBoundingBox()

	char2 := ConvertToCharacter(e2)

	if char2 != nil {
		e2x = e2x - e2w/2
		e2y = e2y - e2h/2
	}

	if inRect(e1x, e1y, e2x, e2y, e2w, e2h) {
		return true
	} else if inRect(e1x+e1w, e1y, e2x, e2y, e2w, e2h) {
		return true
	} else if inRect(e1x, e1y+e1h, e2x, e2y, e2w, e2h) {
		return true
	} else if inRect(e1x+e1w, e1y+e1h, e2x, e2y, e2w, e2h) {
		return true
	}

	return false
}

//let us make the collision handler and the collision finder combined to be the same in this set up
//for diamonds intersecting with squares, there is only 4 cases in which you need to solve. Diamond poke
//or square poke

//will do this from 4-5

//the intersectionhas to be solved on the left or right side

//basically the skill is picking the side that you want the diagonal to slide to the point
//top points should be resolved with first

//since walls will be drawn with their edges for easy hashing, i will compute the center point

//find the potential collisions and at the end develop a loop function that can resolve them

//dynamic ecb collision is going to be pretty cool, I think that I will implement that
//But having priority over the intersections would be a cleaner thing to do, and you can definitively
//solve the problem

//now

//you can assemble the local collisions here, and append it to a global map of character collisions
//not that bad, few extra lines of code

//there are some assumptions about the size of walls and the size of characters. Walls >> characters
//platforms

//there is no special logic that has to be implemented in this case
