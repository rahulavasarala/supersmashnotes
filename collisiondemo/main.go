package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rahulavasarala/supersmashnotes/collisions"
)

type Game struct {
	SpatialGrid collisions.SpatialGrid
}

func (g *Game) Update() error {
	g.SpatialGrid.ClearGrid()
	g.SpatialGrid.UpdateEntityPositions()
	err := g.SpatialGrid.Rehash()
	if err != nil {
		return err
	}
	g.SpatialGrid.PurgeDeadEntities()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	entities := g.SpatialGrid.GetEntities()

	entCount := g.SpatialGrid.GetEntityAmounts()

	for i := 0; i < len(entCount); i++ {
		for j := 0; j < len(entCount[0]); j++ {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%v", entCount[i][j]), 125*i, 125*j)
		}
	}

	for _, ent := range entities {
		xpos, ypos := ent.GetPos()
		width, height := ent.GetBoundingBox()
		vector.DrawFilledRect(screen, float32(xpos), float32(ypos), float32(width), float32(height), color.RGBA{R: 0, G: 255, B: 0, A: 255}, false)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Hello, World!")
	game := Game{}

	//create the collision map with all of the entities
	spatialGrid := collisions.SpatialGrid{}

	entityList := []collisions.Entity{}

	for i := 0; i < 20; i++ {
		dummyEnt := collisions.DummyEntity{}
		dummyEnt.InitDummyEntity()
		entityList = append(entityList, &dummyEnt)
	}

	spatialGrid.InitGrid(500, 500, 4, 4, entityList)

	game.SpatialGrid = spatialGrid

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
