package collisions

import (
	"errors"
	"fmt"

	"math"
)

type SpatialGrid struct {
	width  int
	height int
	xaxis  int
	yaxis  int

	grid [][]GridSquare

	thingDictionary map[string]Thing

	characterList []Character
	Immovables    []Thing
}

type GridSquare struct {
	Characters []Character
	Immovables []Thing
}

type Character interface {
	Step()
	GetPos() (float64, float64)
	SetPos(ypos float64, xpos float64)
	IsPurged() bool
	GetBoundingBox() (float64, float64)
	GetType() string
	GetVel() (float64, float64)
	SetVel(xv float64, yv float64)
	GetHurtbox() []HurtBox
	GetHitbox() []HitBox
	GetEcb() (float64, float64)
	GetId() string
	//add the logic for set hitstun later, for now just role with this
}

type HurtBox struct {
	xpos   float64
	ypos   float64
	width  int
	height int
}

type HitBox struct {
	xpos   float64
	ypos   float64
	width  int
	height int
}

type Thing interface {
	GetPos() (float64, float64)
	GetBoundingBox() (float64, float64)
	GetType() string
	GetId() string
}

func ConvertToCharacter(e Thing) Character {
	var eInterface interface{} = e

	char, ok := eInterface.(Character)

	if !ok {
		return nil
	}

	return char
}

func (s *SpatialGrid) InitGrid(swidth int, sheight int, xaxis int, yaxis int, characterList []Character, immovables []Thing) error {

	//Sanity check the input parameters

	if swidth <= 200 || swidth > 1800 {
		return fmt.Errorf("width of spatial grid is invalid, recieved %v", swidth)
	}

	if sheight <= 200 || sheight > 1800 {
		return fmt.Errorf("height of spatial grid is invalid, recieved %v", sheight)
	}

	if xaxis < 1 {
		return fmt.Errorf("xaxis of spatial grid is invalid, recieved %v", xaxis)
	}
	if yaxis < 1 {
		return fmt.Errorf("yaxis of spatial grid is invalid, recieved %v", yaxis)
	}

	s.width = swidth
	s.height = sheight
	s.xaxis = xaxis
	s.yaxis = yaxis

	s.grid = [][]GridSquare{}

	for i := 0; i < s.xaxis; i++ {
		gridList := []GridSquare{}

		for j := 0; j < s.yaxis; j++ {
			gridList = append(gridList, GridSquare{})
		}

		s.grid = append(s.grid, gridList)
	}

	s.characterList = characterList
	s.Immovables = immovables

	s.thingDictionary = map[string]Thing{}

	for i := 0; i < len(characterList); i++ {
		s.thingDictionary[characterList[i].GetId()] = characterList[i]
	}

	for i := 0; i < len(immovables); i++ {
		s.thingDictionary[immovables[i].GetId()] = immovables[i]
	}

	for i, wall := range immovables {
		xWall, yWall := wall.GetPos()
		wallWidth, wallHeight := wall.GetBoundingBox()

		sqList, err := findHashPositions(xWall, yWall, wallWidth, wallHeight, s.width, s.height, s.xaxis, s.yaxis)

		if err != nil {
			return err
		}

		for _, sq := range sqList {
			s.grid[sq.xcoord][sq.ycoord].Immovables = append(s.grid[sq.xcoord][sq.ycoord].Immovables, immovables[i])
		}
	}

	return nil
}

// seems good, if a character is purged, then the character will not appear in the character list and thingDictionary
func (s *SpatialGrid) PurgeDeadEntities() {

	aliveCharacters := []Character{}
	for i := range s.characterList {
		if s.characterList[i].IsPurged() {
			delete(s.thingDictionary, s.characterList[i].GetId())
		} else {
			aliveCharacters = append(aliveCharacters, s.characterList[i])
		}
	}

	s.characterList = aliveCharacters
}

func (s *SpatialGrid) ClearGrid() {
	for i := 0; i < s.xaxis; i++ {
		for j := 0; j < s.yaxis; j++ {
			s.grid[i][j].Characters = []Character{}
		}
	}
}

func (s *SpatialGrid) UpdateEntityPositions() {
	for i := range s.characterList {
		s.characterList[i].Step()
	}
}

// seems good
func (s *SpatialGrid) Rehash() error {

	for i, char := range s.characterList {
		xChar, yChar := char.GetPos()

		width, height := char.GetBoundingBox()

		xChar = xChar - width/2
		yChar = yChar - height/2

		sqList, err := findHashPositions(xChar, yChar, width, height, s.width, s.height, s.xaxis, s.yaxis)

		if err != nil {
			return err
		}

		for _, square := range sqList {
			s.grid[square.xcoord][square.ycoord].Characters = append(s.grid[square.xcoord][square.ycoord].Characters, s.characterList[i])
		}
	}

	return nil

}

type Square struct { //rows and columns refer to x and
	xcoord int
	ycoord int
}

// seems good
func findHashPositions(x float64, y float64, w float64, h float64, gridW int, gridH int, xaxis int, yaxis int) ([]Square, error) {

	//this getHashSquares function covers the case in which an object can be outside the map but only hashes into
	//edge squares

	getHashSquares := func(sq1 Square, sq2 Square) ([]Square, error) {
		if sq1.xcoord > sq2.xcoord {
			return nil, errors.New("sq1 is more right than sq2")
		} else if sq1.ycoord > sq2.ycoord {
			return nil, errors.New("sq1 is higher than sq2")
		}

		squareList := []Square{}

		for i := sq1.xcoord; i <= sq2.xcoord; i++ {
			for j := sq1.ycoord; j <= sq2.ycoord; j++ {
				if i < xaxis && i >= 0 && j < yaxis && j >= 0 {
					squareList = append(squareList, Square{xcoord: i, ycoord: j})
				}
			}
		}

		return squareList, nil
	}

	x1 := int(math.Floor(x / float64(gridW/xaxis)))
	y1 := int(math.Floor(y / float64(gridH/yaxis)))
	x2 := int(math.Floor((x + w) / float64(gridW/xaxis)))
	y2 := int(math.Floor((y + h) / float64(gridH/yaxis)))

	squareBottomLeft := Square{xcoord: x1, ycoord: y1}
	squareTopRight := Square{xcoord: x2, ycoord: y2}

	sqList, err := getHashSquares(squareBottomLeft, squareTopRight)

	return sqList, err

}

func (s *SpatialGrid) GetEntities() []Character {

	return s.characterList

}

func (s *SpatialGrid) GetImmovables() []Thing {

	return s.Immovables

}

func (s *SpatialGrid) GetEntityAmounts() [][]int {

	entCount := [][]int{}

	for i := 0; i < s.xaxis; i++ {
		entCountRow := []int{}
		for j := 0; j < s.yaxis; j++ {
			entCountRow = append(entCountRow, len(s.grid[i][j].Characters)+len(s.grid[i][j].Immovables))
		}
		entCount = append(entCount, entCountRow)
	}

	return entCount

}

func (s *SpatialGrid) GetDimensions() (int, int) {
	return s.xaxis, s.yaxis
}

//potential optimization: seperating the entities with static entities and moving entities/static entities
//do an integration test to make sure everything is working
//make entity dummys that are able to move and make a step by step debugging of the hashing and rehashing
//practice writing some formalized unit tests to prove functionality of

//the next feature would be checking colisions with static entities
//for the sake of demo purposes, I should have 4 invisible walls hashed into the linings of the grid
//first check entity entity collisions, and then check wall entity collisions
