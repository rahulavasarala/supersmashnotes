package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/rahulavasarala/supersmashnotes/controllers"
)

type Game struct {
	debugMode  bool
	controller controllers.Controller
}

func (g *Game) Update() error {
	return nil
}

func printAt(screen *ebiten.Image, message string, xpos float64, ypos float64, gameHeight int) {
	ebitenutil.DebugPrintAt(screen, message, int(xpos), gameHeight-int(ypos))
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	if g.debugMode {
		output := g.controller.GetInputs()
		if output == "" {
			output = "nil"
		}

		printAt(screen, output, 250, 250, 500)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Super Smash Notes")
	controller := controllers.SimpleController{}
	controller.Init("../controllers/samplebuttonmap.yaml")
	game := Game{debugMode: true, controller: &controller}
	ebiten.SetTPS(60)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
