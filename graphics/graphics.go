package graphics

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/rahulavasarala/supersmashnotes/bones"
	"gonum.org/v1/gonum/mat"
)

func Draw(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, gameHeight int, antialias bool) {
	vector.DrawFilledRect(screen, float32(xpos), float32(gameHeight)-float32(ypos), float32(width), -1*float32(height), color, antialias)
}

func DrawStroked(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, thickness float64, gameHeight int, antialias bool) {
	vector.StrokeRect(screen, float32(xpos), float32(gameHeight)-float32(ypos), float32(width), -1*float32(height), float32(thickness), color, antialias)
}

func PrintAt(screen *ebiten.Image, message string, xpos float64, ypos float64, gameHeight int) {
	ebitenutil.DebugPrintAt(screen, message, int(xpos), gameHeight-int(ypos))
}

func DrawCharacter(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, gameHeight int, antialias bool) {
	xpos = xpos - width/2
	ypos = ypos - height/2

	Draw(screen, color, xpos, ypos, width, height, gameHeight, antialias)
}

func DrawLine(screen *ebiten.Image, color color.Color, x1 float64, y1 float64, x2 float64, y2 float64, thickness float64, gameHeight int, antialias bool) {
	vector.StrokeLine(screen, float32(x1), float32(gameHeight)-float32(y1), float32(x2), float32(gameHeight)-float32(y2), float32(thickness), color, antialias)
}

func DrawEcb(screen *ebiten.Image, color color.Color, xpos float64, ypos float64, width float64, height float64, thickness float64, gameHeight int, antialias bool) {
	//this method should trace 4 red diamond lines
	x1 := xpos - width/2
	x2 := xpos + width/2
	y1 := ypos - height/2
	y2 := ypos + height/2

	DrawLine(screen, color, x1, ypos, xpos, y1, thickness, gameHeight, antialias)
	DrawLine(screen, color, x1, ypos, xpos, y2, thickness, gameHeight, antialias)
	DrawLine(screen, color, x2, ypos, xpos, y1, thickness, gameHeight, antialias)
	DrawLine(screen, color, x2, ypos, xpos, y2, thickness, gameHeight, antialias)
}

func DrawBone(screen *ebiten.Image, color color.Color, thickness float64, gameHeight int, antialias bool, debugMode bool, bone *bones.Bone, frame *mat.Dense) {
	//so you use the frame which is ove

	width := bone.GetWidth()

	rows, cols := frame.Dims()

	if rows != 3 && cols != 3 {
		return
	}

	xbase := frame.At(0, 2)
	ybase := frame.At(1, 2)

	xv, yv := frame.At(0, 0), frame.At(1, 0)

	x1, y1 := xbase+xv*width/2, ybase+yv*width/2

	x2, y2 := xbase-xv*width/2, ybase-yv*width/2

	DrawLine(screen, color, x1, y1, x2, y2, thickness, gameHeight, antialias)

	if debugMode {
		PrintAt(screen, fmt.Sprintf("%v", bone.GetId()), xbase, ybase, gameHeight)
	}

}
