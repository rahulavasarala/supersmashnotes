package statemachinery

// import (
// 	"fmt"
// 	"testing"
// )

// func createIdleRunDashStates() map[string]State {
// 	chicken := map[string]State{}
// 	chicken["dash"] = State{"dash", 3, map[string][]changeStateWindow{"side": {changeStateWindow{3, 3, "run"}}}, "idle"}
// 	chicken["idle"] = State{"idle", 3, map[string][]changeStateWindow{"side": {changeStateWindow{0, 3, "dash"}}}, "idle"}
// 	chicken["run"] = State{"run", 3, map[string][]changeStateWindow{"side": {changeStateWindow{0, 3, "run"}}}, "idle"}

// 	return chicken
// }

// func TestStateMachineTransitioning(t *testing.T) {
// 	tests := []struct {
// 		states         map[string]State
// 		controls       []string
// 		expectedStates []string
// 	}{
// 		{
// 			states:         createIdleRunDashStates(),
// 			controls:       []string{"nut", "nut", "nut", "nut", "side", "nut", "nut", "nut", "nut"},
// 			expectedStates: []string{"idle", "idle", "idle", "idle", "dash", "dash", "dash", "dash", "idle"},
// 		},
// 		{
// 			states:         createIdleRunDashStates(),
// 			controls:       []string{"nut", "nut", "side", "nut", "nut", "side", "side", "side", "side"},
// 			expectedStates: []string{"idle", "idle", "dash", "dash", "dash", "dash", "run", "run", "run"},
// 		},
// 	}

// 	for testnum, test := range tests {
// 		sm := StateMachine{}
// 		sm.init(test.states)

// 		fmt.Println("output for test case: ", testnum)

// 		for i, event := range test.controls {

// 			prevState := sm.currState
// 			prevFrame := sm.currFrame

// 			sm.tick(event)

// 			fmt.Printf("%v on frame %v + %v = %v on frame %v \n", prevState, prevFrame, event, sm.currState, sm.currFrame)

// 			if sm.currState != test.expectedStates[i] {

// 				t.Errorf("%v expected state and %v current state are not matching on iteration %v", test.expectedStates[i], sm.currState, i)
// 			}

// 		}
// 	}
// }
