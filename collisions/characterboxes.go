package collisions

type HurtBox struct {
	width  float64
	height float64
	//define properties of the hurtbox in the future
}

func (s *HurtBox) Init(width float64, height float64) {
	s.width = width
	s.height = height
}

func (s *HurtBox) GetWidth() float64 {
	return s.width
}

func (s *HurtBox) GetHeight() float64 {
	return s.height
}
