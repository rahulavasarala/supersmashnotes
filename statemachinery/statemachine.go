package statemachinery

import "log"

type StateMachine struct {
	currState string
	currFrame int

	states map[string]State
}

func (s *StateMachine) Init(states map[string]State) {
	s.currState = "freefall"
	s.currFrame = 0
	s.states = states
}

func (s *StateMachine) obtainStateChangeInfo(control string) (StringProperty, bool) {
	stateChange, exists := s.states[s.currState].controlToState[control]
	return stateChange, exists
}

//The tick function of the state machine uses info on the current frame, and the input provided
//to change the current state or not

func (s *StateMachine) Tick(control string) {
	stateChange, exists := s.obtainStateChangeInfo(control)

	if exists {
		for i, window := range stateChange.timeline {
			if window.first <= s.currFrame && s.currFrame <= window.second {
				s.currState = stateChange.values[i]
				s.currFrame = 0
				return
			}
		}

	}

	//basically checking for a cancellation to another state in the code above

	s.currFrame++

	if s.currFrame > s.states[s.currState].endingFrame {
		s.currState = s.states[s.currState].loopState
		s.currFrame = 0
	}
}

func (s *StateMachine) GetState() string {
	return s.currState
}

func (s *StateMachine) SetState(newState string) {
	_, ok := s.states[newState]

	if !ok {
		return
	}

	s.currState = newState
}

func (s *StateMachine) GetFrame() int {
	return s.currFrame
}

func (s *StateMachine) SetFrame(frame int) {
	s.currFrame = frame
}

func (s *StateMachine) GetProperties() map[string]any {
	state := s.states[s.currState]

	return state.ExtractProperties(s.currFrame)
}

func (s *StateMachine) GetPropertyByName(name string) Property {
	state := s.states[s.currState]

	property, ok := state.properties[name]

	if !ok {
		log.Fatalf("so such property by name for state %v exists", s.currState)
	}

	return property
}
