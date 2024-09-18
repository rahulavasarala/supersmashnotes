package collisions

import (
	"github.com/rahulavasarala/supersmashnotes/controllers"
	"github.com/rahulavasarala/supersmashnotes/statemachinery"
)

//so the state dude will have thing, and character properties

type StateDude struct {
	xpos        float64
	ypos        float64
	ecbwidth    float64
	ecbheight   float64
	lifeSpan    int
	isGrounded  bool
	xbound1     float64
	xbound2     float64
	xvel        float64
	yvel        float64
	terminalvel float64
	gravity     float64
	maxdrift    float64
	airdrift    float64
	airdrag     float64
	id          string
	sm          *statemachinery.StateMachine
	controller  controllers.Controller
	orientation int
}

func (s *StateDude) Init(sm *statemachinery.StateMachine, id string, controller controllers.Controller) {
	s.xpos = 250
	s.ypos = 250
	s.ecbwidth = 50
	s.ecbheight = 50
	s.lifeSpan = 10000
	s.isGrounded = false
	s.xvel = 0
	s.yvel = 0
	s.id = id
	s.gravity = 0.2
	s.terminalvel = -1
	s.orientation = 1
	s.maxdrift = 2
	s.airdrift = 0.4
	s.airdrag = 0.4
	s.controller = controller

	s.sm = sm

}

func (s *StateDude) Step() {
	//for the sake of convenience let us just implement gravity in the ecbdude step

	//what will the step function be? I think it will be obviously, apply gravity, which is a helper function, it sees the state
	//you are in and gets

	//so it is just applying the properties on the map, this is a helper function that is called by step, you first get the properties
	//from the state machine, and then you apply the properties, basically parsing out the properties per state, not that innefficient
	if s.lifeSpan <= 0 {
		return
	}

	control := s.controller.GetInputs()

	s.sm.Tick(control)

	s.editProperties(control)

	propertyValues := s.sm.GetProperties()

	s.applyProperties(propertyValues, control) //this will apply the isGrounded property, which is pretty important
	s.applyGravity()
	s.applyAerialDrift()
	s.move()

	if s.isGrounded && (s.xpos < s.xbound1 || s.xpos > s.xbound2) {
		s.sm.SetState("freefall")
		s.sm.SetFrame(0)
	}

	s.lifeSpan--

}

func (s *StateDude) GetBoundingBox() (float64, float64) {
	return s.ecbwidth, s.ecbheight
}

func (s *StateDude) GetHitbox() []HitBox {
	return nil
}

func (s *StateDude) GetHurtbox() []HurtBox {
	return nil
}

func (s *StateDude) GetPos() (float64, float64) {

	return s.xpos, s.ypos

}

func (s *StateDude) SetPos(xp float64, yp float64) {
	s.xpos = xp
	s.ypos = yp
}

func (s *StateDude) GetVel() (float64, float64) {
	return s.xvel, s.yvel
}

func (s *StateDude) SetVel(xv float64, yv float64) {
	s.xvel = xv
	s.yvel = yv
}

func (s *StateDude) GetType() string {
	return "blubbula"
}

func (s *StateDude) GetEcb() (float64, float64) {
	return s.ecbwidth, s.ecbheight
}

func (s *StateDude) IsPurged() bool {
	return false
}

func (s *StateDude) GetId() string {
	return s.id
}

func (s *StateDude) SetBounds(bound1 float64, bound2 float64) {
	if bound2 <= bound1 {
		return
	}

	s.xbound1 = bound1
	s.xbound2 = bound2
}

func (s *StateDude) GetBounds() (float64, float64) {
	return s.xbound1, s.xbound2
}

func (s *StateDude) editProperties(control string) {
	if s.GetState() == "firefox" && s.sm.GetFrame() == 61 {
		s.sm.GetPropertyByName("xvel").Reset()
		s.sm.GetPropertyByName("yvel").Reset()
		if control == "left" {
			s.orientation = -1
			s.sm.GetPropertyByName("xvel").Alter(61, statemachinery.NewDoublePair(3, 3))
			s.sm.GetPropertyByName("yvel").Alter(61, statemachinery.NewDoublePair(0, 0))
		} else if control == "right" {
			s.orientation = 1
			s.sm.GetPropertyByName("xvel").Alter(61, statemachinery.NewDoublePair(3, 3))
			s.sm.GetPropertyByName("yvel").Alter(61, statemachinery.NewDoublePair(0, 0))
		}
	}
}

func (s *StateDude) applyProperties(properties map[string]any, control string) {

	if s.sm.GetState() == "idle" {
		s.applyXvel(&properties)
		s.applyYvel(&properties)
	} else if s.sm.GetState() == "freefall" {
		s.applyIsGrounded(&properties)
	} else if s.sm.GetState() == "dash" {
		if s.sm.GetFrame() == 0 {
			if control == "left" {
				s.orientation = -1
			} else if control == "right" {
				s.orientation = 1
			}
		}
		s.applyXvel(&properties)
		s.applyYvel(&properties)
	} else if s.sm.GetState() == "js" {
		s.applyYvel(&properties)
	} else if s.sm.GetState() == "firefox" {
		s.applyXvel(&properties)
		s.applyYvel(&properties)
	}
}

func (s *StateDude) applyGravity() {

	if s.GetState() == "freefall" {
		if s.yvel > s.terminalvel {
			s.yvel -= s.gravity
		}

		if s.yvel < s.terminalvel {
			s.yvel = s.terminalvel
		}
	}
}

func (s *StateDude) applyAerialDrift() {

	if s.GetState() == "freefall" {
		direction := s.controller.GetDirection()

		if direction == "left" {
			s.xvel -= s.airdrift
			if s.xvel < -1*s.maxdrift {
				s.xvel = -1 * s.maxdrift
			}
		} else if direction == "right" {
			s.xvel += s.airdrift
			if s.xvel > s.maxdrift {
				s.xvel = s.maxdrift
			}
		} else {
			if s.xvel < 0 {
				s.xvel += s.airdrag
				if s.xvel > 0 {
					s.xvel = 0
				}
			} else if s.xvel > 0 {
				s.xvel -= s.airdrag
				if s.xvel < 0 {
					s.xvel = 0
				}
			}
		}
	}

}

func (s *StateDude) move() {
	s.xpos = s.xpos + s.xvel
	s.ypos = s.ypos + s.yvel
}

func (s *StateDude) applyXvel(properties *map[string]any) {

	xvelInterface, ok := (*properties)["xvel"]

	if !ok {
		return
	}

	xvel, ok2 := xvelInterface.(float64)

	if !ok2 {
		return
	}

	s.xvel = float64(s.orientation) * xvel

	//in the future add orientation ignore moves
}

func (s *StateDude) applyYvel(properties *map[string]any) {

	yvelInterface, ok := (*properties)["yvel"]

	if !ok {
		return
	}

	yvel, ok2 := yvelInterface.(float64)

	if !ok2 {
		return
	}

	s.yvel = yvel
}

func (s *StateDude) applyIsGrounded(properties *map[string]any) {

	isGroundedInterface, ok := (*properties)["isGrounded"]

	if !ok {
		return
	}

	isGrounded, ok2 := isGroundedInterface.(bool)

	if !ok2 {
		return
	}

	s.isGrounded = isGrounded
}

func (s *StateDude) GetState() string {
	return s.sm.GetState()
}

func (s *StateDude) SetState(newState string) {
	s.sm.SetState(newState)
	s.sm.SetFrame(0)
}

func (s *StateDude) GetGrounded() bool {
	return s.isGrounded
}

func (s *StateDude) SetGrounded(val bool) {
	s.isGrounded = val
}
