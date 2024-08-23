package collisions

type EcbDude struct {
	xpos        float64
	ypos        float64
	ecbwidth    float64
	ecbheight   float64
	xvel        float64
	yvel        float64
	terminalvel float64
	gravity     float64
	lifespan    int
	entType     string
	id          string

	//let us just do a gravity based simulation, where the characters are constantly being kept to the ground and pushed up
	//by their ecbs
}

func (s *EcbDude) Step() {
	//for the sake of convenience let us just implement gravity in the ecbdude step

	if s.yvel > s.terminalvel {
		s.yvel -= s.gravity
	}

	if s.yvel < s.terminalvel {
		s.yvel = s.terminalvel
	}

	s.xpos += s.xvel
	s.ypos += s.yvel
}

func (s *EcbDude) GetBoundingBox() (float64, float64) {
	return s.ecbwidth, s.ecbheight
}

func (s *EcbDude) GetHitbox() []HitBox {
	return nil
}

func (s *EcbDude) GetHurtbox() []HurtBox {
	return nil
}

func (s *EcbDude) GetPos() (float64, float64) {

	return s.xpos, s.ypos

}

func (s *EcbDude) SetPos(xp float64, yp float64) {
	s.xpos = xp
	s.ypos = yp
}

func (s *EcbDude) GetVel() (float64, float64) {
	return s.xvel, s.yvel
}

func (s *EcbDude) SetVel(xv float64, yv float64) {
	s.xvel = xv
	s.yvel = yv
}

func (s *EcbDude) GetType() string {
	return s.entType
}

func (s *EcbDude) GetEcb() (float64, float64) {
	return s.ecbwidth, s.ecbheight
}

func (s *EcbDude) IsPurged() bool {
	return false
}

func (s *EcbDude) GetId() string {
	return s.id
}

func (s *EcbDude) InitEcbDude(xpos float64, ypos float64, ecbwidth float64, ecbheight float64, id string) {
	s.xpos = xpos
	s.ypos = ypos
	s.ecbwidth = ecbwidth
	s.ecbheight = ecbheight
	s.id = id
	s.terminalvel = -2
	s.gravity = 0.2
	s.xvel = 0
	s.yvel = 0
}
