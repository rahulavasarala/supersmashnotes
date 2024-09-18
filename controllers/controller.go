package controllers

import (
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"gopkg.in/yaml.v2"
)

//controllers need to support reading multiple combinations of inputs
//for now, the simple controller will not have sheild functionality

//you need to make sure there are no collisions in the map

type Controller interface {
	GetInputs() string
	GetDirection() string
	GetSpecial() bool
	GetNormal() bool
	GetShield() bool
}

type SimpleController struct {
	name      string
	buttonMap map[string]string
}

// this is the Yaml struct that is used to collect the button map
type ButtonMapContainer struct {
	Name    string `yaml:"name"`
	Up      string `yaml:"up"`
	Down    string `yaml:"down"`
	Left    string `yaml:"left"`
	Right   string `yaml:"right"`
	Special string `yaml:"special"`
	Normal  string `yaml:"normal"`
	Shield  string `yaml:"shield"`
}

func (s *SimpleController) Init(buttonMappingFile string) error {
	data, err := os.ReadFile(buttonMappingFile)
	if err != nil {
		log.Fatalf("Error reading ButtonMapping YAML: %v", err)
	}

	// Create an instance of Config
	var config ButtonMapContainer

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling ButtonMapping YAML: %v", err)
	}

	if config.Name == "" {
		return fmt.Errorf("button map config is missing a name")
	}

	s.name = config.Name

	s.buttonMap = map[string]string{}
	s.buttonMap["up"] = config.Up
	s.buttonMap["down"] = config.Down
	s.buttonMap["left"] = config.Left
	s.buttonMap["right"] = config.Right
	s.buttonMap["special"] = config.Special
	s.buttonMap["normal"] = config.Normal
	s.buttonMap["shield"] = config.Shield

	uniquenessMap := map[string]bool{}

	for control, val := range s.buttonMap {
		if val == "" {
			return fmt.Errorf("keys for control %v are missing for %v", control, s.name)
		}
		if _, ok := uniquenessMap[val]; ok {
			return fmt.Errorf("multiple controls map to the same key in %v", s.name) //throw the names of the controller in the error
		} else {
			uniquenessMap[val] = true
		}
	}

	return nil
}

func (s *SimpleController) GetInputs() string {

	special := s.GetSpecial()
	normal := s.GetNormal()
	direction := s.GetDirection()
	shield := s.GetShield()

	if special {
		return fmt.Sprintf("%vspecial", direction)
	} else if normal {
		return fmt.Sprintf("%vnormal", direction)
	} else if shield {
		return "shield"
	} else if direction != "" {
		return direction
	}

	return ""
}

func (s *SimpleController) GetSpecial() bool {
	result := ObserveButtonIsPressed(s.buttonMap["special"])
	return result
}

func (s *SimpleController) GetNormal() bool {
	result := ObserveButtonIsPressed(s.buttonMap["normal"])
	return result
}

func (s *SimpleController) GetShield() bool {
	result := ObserveButtonIsPressed(s.buttonMap["shield"])
	return result
}

func (s *SimpleController) GetDirection() string {
	left := ObserveButtonIsPressed(s.buttonMap["left"])
	right := ObserveButtonIsPressed(s.buttonMap["right"])
	up := ObserveButtonIsPressed(s.buttonMap["up"])
	down := ObserveButtonIsPressed(s.buttonMap["down"])

	simultaneosClick := 0

	if left {
		simultaneosClick++
	}

	if right {
		simultaneosClick++
	}

	if up {
		simultaneosClick++
	}

	if down {
		simultaneosClick++
	}

	if simultaneosClick >= 2 {
		return ""
	}

	if left {
		return "left"
	} else if right {
		return "right"
	} else if up {
		return "up"
	} else if down {
		return "down"
	}

	return ""
}

func ObserveButtonIsPressed(button string) bool {
	key, valid := GetEbitenKey(button)

	if !valid {
		return false
	}

	return ebiten.IsKeyPressed(key)
}

func GetEbitenKey(button string) (ebiten.Key, bool) {
	switch button {
	case "A":
		return ebiten.KeyA, true
	case "S":
		return ebiten.KeyS, true
	case "D":
		return ebiten.KeyD, true
	case "F":
		return ebiten.KeyF, true
	case "H":
		return ebiten.KeyH, true
	case "J":
		return ebiten.KeyJ, true
	case "K":
		return ebiten.KeyK, true
	case "L":
		return ebiten.KeyL, true
	case "I":
		return ebiten.KeyI, true
	case "O":
		return ebiten.KeyO, true
	case "Q":
		return ebiten.KeyQ, true
	case "E":
		return ebiten.KeyE, true
	default:
		return ebiten.Key0, false
	}
}

//The thing that I realized is that, basically, having multiple controls being sent at a time is slower, and makes the game more
//forgiving and less manual like melee or any hard skill
