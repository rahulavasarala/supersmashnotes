package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rahulavasarala/supersmashnotes/bones"
	"github.com/rahulavasarala/supersmashnotes/graphics"
)

var BASEX float64 = 200
var BASEY float64 = 200

type Game struct {
	debugMode bool
	wireframe *bones.WireFrame
	currFrame int
	animation bones.Animation
}

func (s *Game) InitPlayer(debugMode bool, boneConfig string, animationConfig string) {
	wireFrame := bones.WireFrame{}
	wireFrame.InitWireFrame(boneConfig)
	s.wireframe = &wireFrame
	s.debugMode = debugMode
	s.currFrame = 0

	animation := bones.Animation{}
	animation.InitAnimation(animationConfig)
	s.animation = animation
}

func (g *Game) Update() error {

	g.wireframe.ApplyAnimation(g.animation, g.currFrame, 200, 200)

	g.currFrame++

	if g.currFrame > 20 {
		g.currFrame = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	frameMap := g.wireframe.FindGlobalBoneFrames(g.wireframe.GetBone(0))

	for key := range frameMap {
		bone := g.wireframe.GetBone(key)

		graphics.DrawBone(screen, color.Black, 2, 500, false, true, bone, frameMap[key])
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
}

func main() {
	ebiten.SetWindowSize(500, 500)
	ebiten.SetWindowTitle("Super Smash Notes")

	game := Game{}

	game.InitPlayer(true, "../bones/boneconfig1.yaml", "../animationtool/testanimationconfig.yaml")

	ebiten.SetTPS(60)

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}

//I will make a simple animation rig that reads from a config file, constantly
//polling it for changes in angles, then I will render the object

//this will happen for a while

//change the value of one thing to save the angles in a timeline

//you have something called an animation map
//it will be like the state map, each animation will have certain bone positions

//position of the character will be the torso mid point

//First what is going to happen is, basically, you have a frameList, which has the saves of all the frames

//then you have a map that is keyed by a string(bone_bone), which corresponds to a list of values per save frame

//then you basically create a new datastructure that holds bone, bone, double property,  for each joint

//then you populate this datastructure, adn then you can implement a read method for the wire frame, and then you are don
