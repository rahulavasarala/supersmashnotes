package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rahulavasarala/supersmashnotes/collisions"
	"github.com/rahulavasarala/supersmashnotes/graphics"
)

type Game struct {
	SpatialGrid      collisions.SpatialGrid
	CollisionFinder  collisions.CollisionFinder
	CollisionHandler collisions.CollisionHandler
	debugMode        bool
	width            int
	height           int
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	graphics.DrawRotatedRectangle(screen, color.RGBA{0, 0, 0, 30}, 10, 25, 10, 50, 0, 500, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Super Smash Notes")
	game := Game{width: 500, height: 500, debugMode: true}
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
