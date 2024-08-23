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
