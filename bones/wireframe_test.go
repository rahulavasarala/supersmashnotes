package bones

import (
	"fmt"
	"testing"

	"github.com/rahulavasarala/supersmashnotes/properties"
	"gonum.org/v1/gonum/mat"
)

func TestMatrixMultiplication(t *testing.T) {

	data := []float64{
		1, 2,
		4, 5,
	}

	// Initialize a matrix
	m := mat.NewDense(2, 2, data)

	m2 := mat.NewDense(2, 2, data)
	m3 := new(mat.Dense)

	m3.Mul(m, m2)

	if m3.At(0, 0) != 9 {
		t.Errorf("wrong value found , %v", m3.At(0, 0))
	}

}

func matPrint(m *mat.Dense) string {
	fc := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	return fmt.Sprintf("Matrix:\n%v\n", fc)
}

func TestFindFrames(t *testing.T) {
	wireframe := WireFrame{}

	err := wireframe.InitWireFrame("./boneconfig1.yaml")

	if err != nil {
		t.Errorf("%v", err)
	}

	frameMap := map[int]*mat.Dense{}
	frame := mat.NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})

	visited := map[int]bool{}

	FindFrames(wireframe.boneMap[0], frame, frameMap, visited)

	t.Errorf("%v", matPrint(frameMap[2]))

}

// WireFrame init is working! YAY
func TestWireFrameInit(t *testing.T) {
	wireframe := WireFrame{}

	err := wireframe.InitWireFrame("./boneconfig1.yaml")

	if err != nil {
		t.Errorf("%v", err)
	}

	if result, side := wireframe.boneMap[0].GetLink(1); result == nil || result.id != 1 || side == "left" {
		t.Errorf("bone 0 is not connected to bone 1, it is connected to %v on side %v", result.id, side)
	}

	if result, _ := wireframe.boneMap[1].GetLink(2); result != nil {
		t.Errorf("false connection detected between bones 1 and 2")
	}

	if result, _ := wireframe.boneMap[2].GetLink(1); result != nil {
		t.Errorf("false connection detected between bones 1 and 2")
	}

	if result, side := wireframe.boneMap[1].GetLink(0); result == nil || result.id != 0 || side == "right" {
		t.Errorf("bone 1 is not connected to bone 0, it is connected to %v on side %v", result.id, side)
	}
}

// Ok the write animation config is working
func TestWriteAnimationConfig(t *testing.T) {
	config := AnimationConfig{}

	config.BaseOrientation = []properties.RangeValue{{Range: "chick", Value: "ken"}, {Range: "chick", Value: "ken"}, {Range: "chick", Value: "ken"}}

	config.JointConfigList = []JointConfig{}

	config.JointConfigList = append(config.JointConfigList, JointConfig{J1: 1, J2: 2, ThetaList: config.BaseOrientation})

	WriteAnimationConfigToYAML("chicken.yaml", &config)
}
