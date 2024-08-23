package statemachinery

type State struct {
	name           string
	endingFrame    int
	controlToState map[string][]changeStateWindow
	loopState      string
}

type changeStateWindow struct {
	startFrame    int
	endFrame      int
	stateToChange string
}

//You need to have this model of state, because let us say you want to cancel a state or cancel a dash.
//also you need to coordinate the movement of the character's hurt and hitboxes based on the current
//frame of the state
