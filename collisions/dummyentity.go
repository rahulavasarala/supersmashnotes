package collisions

import "math/rand"

//for the sake of a comprehensive demo, each of the Dummy entities will be 20,20 in a 500,500 field with 4 rows and 4 cols
//the demo will show all of the entities inside the squares at once

type DummyEntity struct {
	xpos     float64
	ypos     float64
	bbwidth  int
	bbheight int
	xvel     float64
	yvel     float64

	entType  string
	lifeSpan int
}

func (s *DummyEntity) InitDummyEntity() {
	s.bbwidth = 20
	s.bbheight = 20
	rN := rand.Intn(2)
	if rN == 1 {
		s.xvel = 1
	} else {
		s.xvel = -1
	}
	rN = rand.Intn(2)
	if rN == 1 {
		s.yvel = 1
	} else {
		s.yvel = -1
	}

	rN = rand.Intn(441) + 20
	s.xpos = float64(rN)

	rN = rand.Intn(441) + 20
	s.ypos = float64(rN)

	s.lifeSpan = rand.Intn(500) + 2
	s.entType = "Player"

}

func (s *DummyEntity) Step() {
	if s.xpos < 0 {
		s.xpos = 1
		s.xvel = s.xvel * -1
	}

	if s.xpos > 500 {
		s.xpos = 499
		s.xvel = s.xvel * -1
	}

	if s.ypos < 0 {
		s.ypos = 1
		s.yvel = s.yvel * -1
	}

	if s.ypos > 480 {
		s.ypos = 479
		s.yvel = s.yvel * -1
	}
	s.xpos = s.xpos + s.xvel
	s.ypos = s.ypos + s.yvel
	s.lifeSpan--
}

func (s *DummyEntity) GetPos() (float64, float64) {

	return s.xpos, s.ypos

}

func (s *DummyEntity) IsPurged() bool {
	if s.lifeSpan == 0 {
		return true
	}

	return false
}

func (s *DummyEntity) GetBoundingBox() (int, int) {
	return s.bbwidth, s.bbheight
}

func (s *DummyEntity) GetType() string {
	return s.entType
}

//for the purposes of the demo, i want to put very thin stationary objects
