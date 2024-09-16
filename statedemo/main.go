package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rahulavasarala/supersmashnotes/collisions"
	"github.com/rahulavasarala/supersmashnotes/statemachinery"
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
	g.SpatialGrid.ClearGrid()
	g.SpatialGrid.UpdateEntityPositions()
	err := g.SpatialGrid.Rehash()
	if err != nil {
		return err
	}

	//handle the collisions over here

	collisionMap := g.CollisionFinder.FindPotentialCollisions(g.SpatialGrid)
	//fmt.Println("Collision map: ", collisionMap)
	err2 := g.CollisionHandler.HandleCollisions(g.SpatialGrid, collisionMap)

	if err2 != nil {
		return err
	}

	g.SpatialGrid.PurgeDeadEntities()

	return nil
}

func draw(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, gameHeight int, antialias bool) {
	vector.DrawFilledRect(screen, float32(xpos), float32(gameHeight)-float32(ypos), float32(width), -1*float32(height), color, antialias)
}

func drawStroked(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, thickness float64, gameHeight int, antialias bool) {
	vector.StrokeRect(screen, float32(xpos), float32(gameHeight)-float32(ypos), float32(width), -1*float32(height), float32(thickness), color, antialias)
}

func printAt(screen *ebiten.Image, message string, xpos float64, ypos float64, gameHeight int) {
	ebitenutil.DebugPrintAt(screen, message, int(xpos), gameHeight-int(ypos))
}

func drawCharacter(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, gameHeight int, antialias bool) {
	xpos = xpos - width/2
	ypos = ypos - height/2

	draw(screen, color, xpos, ypos, width, height, gameHeight, antialias)
}

func drawLine(screen *ebiten.Image, color color.Color, x1 float64, y1 float64, x2 float64, y2 float64, thickness float64, gameHeight int, antialias bool) {
	vector.StrokeLine(screen, float32(x1), float32(gameHeight)-float32(y1), float32(x2), float32(gameHeight)-float32(y2), float32(thickness), color, antialias)
}

func drawEcb(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, thickness float64, gameHeight int, antialias bool) {
	//this method should trace 4 red diamond lines
	x1 := xpos - width/2
	x2 := xpos + width/2
	y1 := ypos - height/2
	y2 := ypos + height/2

	drawLine(screen, color, x1, ypos, xpos, y1, thickness, gameHeight, antialias)
	drawLine(screen, color, x1, ypos, xpos, y2, thickness, gameHeight, antialias)
	drawLine(screen, color, x2, ypos, xpos, y1, thickness, gameHeight, antialias)
	drawLine(screen, color, x2, ypos, xpos, y2, thickness, gameHeight, antialias)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	entities := g.SpatialGrid.GetEntities()
	immovables := g.SpatialGrid.GetImmovables()

	for i, ent := range entities {
		xpos, ypos := ent.GetPos()
		width, height := ent.GetBoundingBox()
		drawCharacter(screen, color.RGBA{R: 0, G: 255, B: 0, A: 255}, xpos, ypos, width, height, g.height, false)
		if g.debugMode {
			drawEcb(screen, color.RGBA{R: 255, G: 0, B: 0, A: 255}, xpos, ypos, width, height, 2, g.height, false)
			var stateDudeInterface interface{} = entities[i]
			sc, ok := stateDudeInterface.(collisions.StateCharacter)
			if ok {
				printAt(screen, sc.GetState(), xpos, ypos, g.height)
			}

		}

	}

	for _, wall := range immovables {
		xpos, ypos := wall.GetPos()
		width, height := wall.GetBoundingBox()
		draw(screen, color.RGBA{0, 0, 255, 100}, xpos, ypos, width, height, g.height, false)
	}

	if g.debugMode {
		g.displayDebugOverlay(screen)
	}
}

func (g *Game) displayDebugOverlay(screen *ebiten.Image) {
	entCount := g.SpatialGrid.GetEntityAmounts()
	xaxis, yaxis := g.SpatialGrid.GetDimensions()

	firstCountX := float64(g.width/xaxis) / 2
	firstCountY := float64(g.height/yaxis) / 2
	widthIncrement := float64(g.width / xaxis)
	heightIncrement := float64(g.height / yaxis)

	for i := 0; i < xaxis; i++ {
		for j := 0; j < yaxis; j++ {
			printAt(screen, fmt.Sprintf("gc: %v,%v \nct: %v", i, j, entCount[i][j]), firstCountX+float64(i)*widthIncrement, firstCountY+float64(j)*heightIncrement, g.height)
			drawStroked(screen, color.RGBA{0, 0, 0, 30}, float64(i)*widthIncrement, float64(j)*heightIncrement, widthIncrement, heightIncrement, 1, g.height, false)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Super Smash Notes")
	game := Game{width: 500, height: 500, debugMode: true}
	game.CollisionFinder = &collisions.EcbCollisionFinder{}
	game.CollisionHandler = &collisions.StateCollisionHandler{}
	ebiten.SetTPS(60)

	//create the collision map with all of the entities
	spatialGrid := collisions.SpatialGrid{}

	characterList := []collisions.Character{}
	sd := collisions.StateDude{}

	smBuilder := statemachinery.StateMachineBuilder{}

	sm := smBuilder.Build("../statemachinery/foxschema.yaml")

	sd.Init(sm, []string{"nil", "nil", "up", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left", "left"}, "stateDude")
	sd.SetPos(250, 300)
	characterList = append(characterList, &sd)
	wallList := []collisions.Thing{}

	wall1 := collisions.Wall{}
	wall1.InitWall(150, 200, 200, 20, "wall1")

	wallList = append(wallList, &wall1)

	spatialGrid.InitGrid(game.width, game.height, 5, 5, characterList, wallList)

	game.SpatialGrid = spatialGrid

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
