package properties

import (
	"testing"
)

// func TestParsersNegative(t *testing.T) {
// 	ran := ParseDoubleRange("-0.5--0.5")

// 	var dpInterface interface{} = x

// 	_, ok := dpInterface.(*statemachinery.DoublePair)

// 	t.Errorf("ok status: %v", ok)

// }

func TestParsers(t *testing.T) {

	pair := ParseIntRange("-5-5")

	t.Errorf("%v", pair)
}

//man the hypothesis was correct, you will get full access to the character version of the thing that
//you get from the thing based collision map

//parseDouble is working, so now I can work on parseInt
