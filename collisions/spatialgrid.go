package collisions

import (
	"errors"
	"fmt"
)

type SpatialGrid struct {
	width  int
	height int
	rows   int
	cols   int

	grid [][]GridSquare

	entityList []Entity
}

type GridSquare struct {
	Entities []Entity
}

type Entity interface {
	Step()
	GetPos() (float64, float64)
	IsPurged() bool
	GetBoundingBox() (int, int)
	GetType() string
}

func (s *SpatialGrid) InitGrid(swidth int, sheight int, srows int, scols int, entityList []Entity) error {

	//Sanity check the input parameters

	if swidth <= 200 || swidth > 1800 {
		return fmt.Errorf("width of spatial grid is invalid, recieved %v", swidth)
	}

	if sheight <= 200 || sheight > 1800 {
		return fmt.Errorf("height of spatial grid is invalid, recieved %v", sheight)
	}

	if srows <= 2 {
		return fmt.Errorf("rows of spatial grid is invalid, recieved %v", swidth)
	}
	if scols <= 2 {
		return fmt.Errorf("cols of spatial grid is invalid, recieved %v", swidth)
	}

	s.width = swidth
	s.height = sheight
	s.rows = srows
	s.cols = scols

	s.grid = [][]GridSquare{}

	for i := 0; i < s.rows; i++ {
		gridList := []GridSquare{}

		for j := 0; j < s.cols; j++ {
			gridList = append(gridList, GridSquare{})
		}

		s.grid = append(s.grid, gridList)
	}

	s.entityList = entityList

	return nil
}

func (s *SpatialGrid) PurgeDeadEntities() {

	aliveEntities := []Entity{}
	for i, _ := range s.entityList {
		if !s.entityList[i].IsPurged() {
			aliveEntities = append(aliveEntities, s.entityList[i])
		}
	}

	s.entityList = aliveEntities
}

func (s *SpatialGrid) ClearGrid() {
	for i := 0; i < s.rows; i++ {
		for j := 0; j < s.cols; j++ {
			s.grid[i][j].Entities = []Entity{}
		}
	}
}

func (s *SpatialGrid) UpdateEntityPositions() {
	for i, ent := range s.entityList {
		if ent.GetType() == "static" {
			continue
		} else {
			s.entityList[i].Step()
		}
	}
}

func (s *SpatialGrid) Rehash() error {

	for i, ent := range s.entityList {
		xpos, ypos := ent.GetPos()

		width, height := ent.GetBoundingBox()

		sqList, err := findHashPositions(xpos, ypos, width, height, s.width, s.height, s.rows, s.cols)

		if err != nil {
			return err
		}

		for _, square := range sqList {
			if square.row < s.rows && square.row >= 0 && square.col < s.cols && square.col >= 0 {
				s.grid[square.row][square.col].Entities = append(s.grid[square.row][square.col].Entities, s.entityList[i])
			}
		}
	}

	return nil

}

type Square struct {
	row int
	col int
}

func findHashPositions(xpos float64, ypos float64, ewidth int, eheight int, gridWidth int, gridHeight int, rows int, columns int) ([]Square, error) {

	findSquares := func(sq1 Square, sq2 Square) ([]Square, error) {
		if sq1.row > sq2.row {
			return nil, errors.New("sq1 is more down than sq2!")
		} else if sq1.col > sq2.col {
			return nil, errors.New("sq1 is more right than sq2!")
		}

		squareList := []Square{}

		for i := sq1.row; i <= sq2.row; i++ {
			for j := sq1.col; j <= sq2.col; j++ {
				squareList = append(squareList, Square{row: i, col: j})
			}
		}

		return squareList, nil
	}

	xgrid := int(xpos) / (gridWidth / columns)
	ygrid := int(ypos) / (gridHeight / rows)

	xgridtwo := (int(xpos) + ewidth) / (gridWidth / columns)
	ygridtwo := (int(ypos) + eheight) / (gridHeight / rows)

	sqList, err := findSquares(Square{row: ygrid, col: xgrid}, Square{row: ygridtwo, col: xgridtwo})

	return sqList, err

}

func (s *SpatialGrid) GetEntities() []Entity {

	return s.entityList

}

func (s *SpatialGrid) GetEntityAmounts() [][]int {

	entCount := [][]int{}

	for i := 0; i < s.rows; i++ {
		entCountRow := []int{}
		for j := 0; j < s.cols; j++ {
			entCountRow = append(entCountRow, len(s.grid[i][j].Entities))
		}
		entCount = append(entCount, entCountRow)
	}

	return entCount

}

//potential optimization: seperating the entities with static entities and moving entities/static entities
//do an integration test to make sure everything is working
//make entity dummys that are able to move and make a step by step debugging of the hashing and rehashing
//practice writing some formalized unit tests to prove functionality of
