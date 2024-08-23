package collisions

type Wall struct {
	xpos   float64
	ypos   float64
	width  float64
	height float64
	id     string
}

func (s *Wall) InitWall(xpos float64, ypos float64, width float64, height float64, id string) {
	s.xpos = xpos
	s.ypos = ypos
	s.width = width
	s.height = height
	s.id = id
}

func (s *Wall) GetType() string {
	return "wall"
}

func (s *Wall) GetId() string {
	return s.id
}

func (s *Wall) GetPos() (float64, float64) {
	return s.xpos, s.ypos
}

func (s *Wall) GetBoundingBox() (float64, float64) {
	return s.width, s.height
}
