package bones

import (
	"fmt"
	"testing"

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

func matPrint(m *mat.Dense) {
	fc := mat.Formatted(m, mat.Prefix("    "), mat.Squeeze())
	fmt.Printf("Matrix:\n%v\n", fc)
}

func TestFindFrames(t *testing.T) {

}
