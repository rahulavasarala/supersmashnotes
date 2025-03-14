package collisions

func inDiamond(x float64, y float64, dx float64, dy float64, dw float64, dh float64) bool {

	//bmws
	m1 := -1 * dh / dw
	m2 := dh / dw
	m3 := -1 * dh / dw
	m4 := dh / dw

	p1 := -m1*dx + (dy + dh/2)
	p2 := -m2*dx + (dy + dh/2)
	p3 := -m3*dx + (dy - dh/2)
	p4 := -m4*dx + (dy - dh/2)

	if y < m1*x+p1 && y < m2*x+p2 && y > m3*x+p3 && y > m4*x+p4 {
		return true
	}

	return false
}

func inRect(x float64, y float64, rx float64, ry float64, rw float64, rh float64) bool {
	if x < rx+rw && x > rx && y > ry && y < ry+rh {
		return true
	}

	return false
}

func FindEncompassingRectangle(points [][]float64) (float64, float64, float64, float64) {
	//this rectangle will be returned in centered format!
	x_min := points[0][0]
	x_max := points[0][0]
	y_min := points[0][1]
	y_max := points[0][1]

	for _, point := range points {
		x := point[0]
		y := point[1]

		if x > x_max {
			x_max = x
		} else if x < x_min {
			x_min = x
		}

		if y < y_min {
			y_min = y
		} else if y > y_max {
			y_max = y
		}

	}

	return (x_max + x_min) / 2., (y_min + y_max) / 2., x_max - x_min, y_max - y_min

}
