package main

import (
	"fmt"

	"github.com/rahulavasarala/supersmashnotes/bones"
)

func main() {

	animation_name := "bounce"

	dest_file := fmt.Sprintf("./%vpoints.csv", animation_name)
	dpi := 2

	wireFrame := bones.WireFrame{}
	wireFrame.InitWireFrame("../bones/boneconfig1.yaml")

	animation := bones.Animation{}
	animation.InitAnimation(fmt.Sprintf("../animationtool/%v.yaml", animation_name))

	largeSamplePoints := [][]string{}

	for i := 0; i < 26; i++ {
		wireFrame.ApplyAnimation(animation, i, 200, 200)

		frameMap := wireFrame.FindGlobalBoneFrames(wireFrame.GetBone(0))

		samplePoints := wireFrame.SampleHurtboxPoints(dpi, frameMap, i)
		largeSamplePoints = append(largeSamplePoints, samplePoints...)
	}

	bones.WriteHurtBoxPointsToCSV(largeSamplePoints, dest_file)
}
