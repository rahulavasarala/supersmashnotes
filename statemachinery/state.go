package statemachinery

type State struct {
	name           string
	endingFrame    int
	controlToState map[string]StringProperty
	properties     map[string]Property
	loopState      string
}

func (s *State) ExtractProperties(frame int) map[string]any {

	propertyVals := map[string]any{}

	for propName, property := range s.properties {
		val := property.Read(frame)

		if val != nil {
			propertyVals[propName] = val
		}

	}

	return propertyVals
}

//You need to have this model of state, because let us say you want to cancel a state or cancel a dash.
//also you need to coordinate the movement of the character's hurt and hitboxes based on the current
//frame of the state

//See, if a state has properties, then the controlToState dictionary should also
//assumption is that on each frame, we have a distinct input, and we can observe inputs in such a way that only 1 input
//exists at all times
