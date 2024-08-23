package collisions

import "fmt"

type EcbCollisionHandler struct {
}

func (s *EcbCollisionHandler) HandleCollisions(sg SpatialGrid, collisionMap map[string]map[string]bool) error {

	for id := range collisionMap {
		for collisionId := range collisionMap[id] {
			err := EcbWallResolution(sg.thingDictionary[id], sg.thingDictionary[collisionId])

			if err != nil {
				return err
			}
		}
	}

	return nil

}

func EcbWallResolution(charAsThing Thing, wall Thing) error {

	//first check to see intersections of the diamond in the square, easy resolutions
	//the resolution of this will guarenteed be able to not make 2 bodies intersect
	var charInterface interface{} = charAsThing

	char, ok := charInterface.(Character)

	if !ok {
		return fmt.Errorf("given char to resolve ecb on is not a character")
	}

	xChar, yChar := char.GetPos()
	ecbWidth, ecbHeight := char.GetEcb()

	xWall, yWall := wall.GetPos()
	wallWidth, wallHeight := wall.GetBoundingBox()

	if inWall(xChar, yChar-ecbHeight/2, wall) { //good
		shiftDistance := yWall + wallHeight - (yChar - ecbHeight/2)
		char.SetPos(xChar, yChar+shiftDistance)
		char.SetVel(0, 0)
		return fmt.Errorf("there is a collision")
	} else if inWall(xChar, yChar+ecbHeight/2, wall) { //good
		shiftDistance := (yChar + ecbHeight/2) - yWall
		char.SetPos(xChar, yChar-shiftDistance)
	} else if inWall(xChar+ecbWidth/2, yChar, wall) { //good
		shiftDistance := (xChar + ecbWidth/2) - xWall
		char.SetPos(xChar-shiftDistance, yChar)
		char.SetVel(0, 0)
	} else if inWall(xChar-ecbWidth/2, yChar, wall) { //good
		shiftDistance := (xWall + wallWidth) - (xChar - ecbWidth/2)
		char.SetPos(xChar+shiftDistance, yChar)
		char.SetVel(0, 0)
	} else if inEcb(xWall, yWall, char) { //good
		m := ecbHeight / ecbWidth
		b := -1*m*(xChar+ecbWidth/2) + yChar
		xIntersect := (yWall - b) / m
		shiftDistance := xIntersect - xWall
		char.SetPos(xChar+shiftDistance, yChar)
	} else if inEcb(xWall+wallWidth, yWall, char) { //good
		m := ecbHeight / ecbWidth
		b := -1*m*(xChar-ecbWidth/2) + yChar
		xIntersect := (yWall - b) / m
		shiftDistance := xIntersect - (xWall + wallWidth)
		char.SetPos(xChar+shiftDistance, yChar)
	} else if inEcb(xWall, yWall+wallHeight, char) { //good
		m := ecbHeight / ecbWidth
		b := -1*m*(xChar+ecbWidth/2) + yChar
		xIntersect := (yWall + wallHeight - b) / m
		shiftDistance := xIntersect - xWall
		char.SetPos(xChar-shiftDistance, yChar)
	} else if inEcb(xWall+wallWidth, yWall+wallHeight, char) { //good
		m := -1 * ecbHeight / ecbWidth
		b := -1*m*(xChar-ecbWidth/2) + yChar
		xIntersect := (yWall + wallHeight - b) / m
		shiftDistance := xIntersect - (xWall + wallWidth)
		char.SetPos(xChar-shiftDistance, yChar)
	}

	return nil
}

// looks good
func inEcb(x float64, y float64, char Character) bool {

	//s1,s2,s3,s4 are the slopes of the diamond lines of each quadrant

	ecbWidth, ecbHeight := char.GetEcb()
	charX, charY := char.GetPos()

	return inDiamond(x, y, charX, charY, ecbWidth, ecbHeight)
}

// looks good
func inWall(x float64, y float64, wall Thing) bool {

	xWall, yWall := wall.GetPos()
	wallWidth, wallHeight := wall.GetBoundingBox()

	return inRect(x, y, xWall, yWall, wallWidth, wallHeight)
}
