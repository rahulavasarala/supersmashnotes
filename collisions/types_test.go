package collisions

import (
	"testing"
)

func TestTypeCheck(t *testing.T) {

	pikaraw := &EcbDude{}

	var pikachu interface{} = pikaraw

	if _, ok := pikachu.(Thing); !ok {
		t.Error("Pikachu was not detected as a thing ")
	} else if _, ok := pikachu.(Character); !ok {
		t.Error("Pikachu was not detected as a character")
	}

}

func TestConvertToCharacter(t *testing.T) {
	wall := &Wall{}

	dude := &EcbDude{}

	char := ConvertToCharacter(wall)

	if char != nil {
		t.Error("wall was detected as a character")
	}

	char2 := ConvertToCharacter(dude)

	if char2 == nil {
		t.Error("dude was not detected as a character")
	}
}

func TestTypeCheck2(t *testing.T) {

	ecbDude := &EcbDude{xpos: 1, ypos: 1}

	thingList := []Thing{ecbDude}

	checkCharacter := func(charAsThing Thing) bool {
		var charInterface interface{} = charAsThing

		char, ok := charInterface.(Character)

		if !ok {
			return false
		}

		char.SetPos(-2, -2)
		return true
	}

	result := checkCharacter(thingList[0])

	xpos, ypos := ecbDude.GetPos()

	t.Errorf("xpos: %v, ypos: %v, result: %v", xpos, ypos, result)

}

//man the hypothesis was correct, you will get full access to the character version of the thing that
//you get from the thing based collision map
