package statemachinery

import (
	"testing"
)

func TestBuilderUnfurl(t *testing.T) {

	smBuilder := StateMachineBuilder{}

	blob := smBuilder.unfurl("foxschema.yaml")

	if true {
		t.Errorf("Could not unfurl foxschema.yaml properly %v", blob)
	}

}

func TestStateMachineGeneration(t *testing.T) {

	smBuilder := StateMachineBuilder{}

	sm := smBuilder.Build("foxschema.yaml")

	if sm.states["firefox"].properties["xvel"].Read(2).(float64) != 0 {
		t.Errorf("Wrong property value read for xvel on frame 2 of firefox animation, read %v", sm.states["firefox"].properties["xvel"].Read(2))
	}

	if sm.states["firefox"].properties["yvel"].Read(7).(float64) != 3 {
		t.Errorf("Wrong property value read for xvel on frame 7 of firefox animation, read %v", sm.states["firefox"].properties["yvel"].Read(7))
	}

	if sm.states["firefox"].properties["yvel"].Read(4).(float64) != 4 {
		t.Errorf("Wrong property value read for xvel on frame 4 of firefox animation, read %v", sm.states["firefox"].properties["yvel"].Read(4))
	}
}

func TestStateMachinAlter(t *testing.T) {

	smBuilder := StateMachineBuilder{}

	sm := smBuilder.Build("foxschema.yaml")

	sm.currState = "firefox"
	sm.currFrame = 0

	dp := NewDoublePair(10, 10)

	sm.GetPropertyByName("xvel").Alter(5, dp)

	if sm.states["firefox"].properties["xvel"].Read(5).(float64) != 10 {
		t.Errorf("Wrong property value read for xvel on frame 5 of firefox animation, read %v", sm.states["firefox"].properties["xvel"].Read(5))
	}
}
