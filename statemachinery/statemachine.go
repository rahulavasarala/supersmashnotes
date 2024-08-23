package statemachinery

type StateMachine struct {
	currState string
	currFrame int

	states map[string]State
}

func (s *StateMachine) init(states map[string]State) {
	s.currState = "idle"
	s.currFrame = 0
	s.states = states
}

func (s *StateMachine) tick(control string) {
	//stay in the same state if you cannot cancel into anything, if you can
	//cancel into something, change to it, if on the last frame and nothing else to do, loop back to the loopstate

	//implement the code to see whether character should cancel to another state

	windowList, exists := s.states[s.currState].controlToState[control]

	if exists {
		for _, window := range windowList {
			if window.startFrame <= s.currFrame && s.currFrame <= window.endFrame {
				s.currState = window.stateToChange
				s.currFrame = 0
				return
			}
		}

	}

	s.currFrame++

	if s.currFrame > s.states[s.currState].endingFrame {
		s.currState = s.states[s.currState].loopState
		s.currFrame = 0
	}
}

//This dude has a dictionary of states, and the current state
//I will make the initial design kind of simple
//will make the state machine be keyed by a string and the string represents the current state
//current state and current frame
//Then I will be victorious

//Create a state filler later in the game, where you can just create states from .txt files

//everyone's starting state is 'idle'
//then the states get better and better
