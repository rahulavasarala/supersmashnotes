package main

import (
	"fmt"

	"github.com/rahulavasarala/supersmashnotes/bones"
)

func main() {
	wireFrame := bones.WireFrame{}

	wireFrame.InitWireFrame("../bones/boneconfig1.yaml")

	animation := bones.Animation{}
	animation.InitAnimation("../animationtool/jab.yaml")

	groups := [][]int{{0, 1, 2, 3}, {4, 5}, {6, 7}, {8, 9}}
	groups2 := [][]int{{0, 1, 2, 7}, {4, 3}, {6, 9}, {8, 5}}

	area1 := wireFrame.CalculateAverageArea(animation, groups, 2)

	area2 := wireFrame.CalculateAverageArea(animation, groups2, 2)

	fmt.Println(area1, area2)

	// sum := 0.

	// for _, area := range areas {
	// 	sum += area
	// }

	// fmt.Println("Sum: ", sum)

	// fmt.Println(areas)

	// _, _, w, h := collisions.FindEncompassingRectangle({}{}float64{{180, 200}, {220, 200}})

	// fmt.Printf("area found: %v, %v", w, h)

}
