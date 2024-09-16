package statemachinery

import (
	"log"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

type StateMachineBuilder struct {
}

type StateListYaml struct {
	States []StateYaml `yaml:"states"`
}

type StateYaml struct {
	Name           string                      `yaml:"name"`
	End            int                         `yaml:"end"`
	Properties     map[string][]RangeValueYaml `yaml:"properties"`
	ControlToState map[string][]RangeValueYaml `yaml:"controlToState"`
	Loop           string                      `yaml:"loop"`
}

type RangeValueYaml struct {
	Range string `yaml:"range"`
	Value string `yaml:"value"`
}

// Ok, so I succeeded in unfurling the yaml file to build state machines, and have a testing platform
func (s *StateMachineBuilder) unfurl(yamlFile string) StateListYaml {

	data, err := os.ReadFile(yamlFile)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	// Create an instance of Config
	var config StateListYaml

	// Unmarshal the YAML file into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error unmarshaling YAML: %v", err)
	}

	return config

}

func (s *StateMachineBuilder) propertyTypeDecider(propertyName string) string {

	if propertyName == "xvel" || propertyName == "yvel" {
		return "double"
	} else if propertyName == "intangibility" || propertyName == "invincibility" || propertyName == "isGrounded" {
		return "bool"
	}

	return "none"

}

func (s *StateMachineBuilder) convertYamlProperties(yamlProperties map[string][]RangeValueYaml) map[string]Property {

	propertyDict := map[string]Property{}

	for key, rangeValueList := range yamlProperties {
		timeline := []IntPair{}
		propertyName := key

		//build the property timeline
		for _, rangeVal := range rangeValueList {
			window := parseIntRange(rangeVal.Range)
			timeline = append(timeline, window)
		}

		result := s.propertyTypeDecider(propertyName)

		if result == "double" {
			property := DoubleProperty{}
			values := []DoublePair{}

			for _, rangeVal := range rangeValueList {
				values = append(values, parseDoubleRange(rangeVal.Value))
			}

			err := property.Init(propertyName, timeline, values)
			if err != nil {
				log.Fatal(err.Error())
			}

			propertyDict[propertyName] = &property

		} else if result == "bool" {
			property := BoolProperty{}
			values := []bool{}

			for _, rangeVal := range rangeValueList {
				values = append(values, parseBool(rangeVal.Value))
			}

			err := property.Init(propertyName, timeline, values)
			if err != nil {
				log.Fatal(err.Error())
			}

			propertyDict[propertyName] = &property
		} else {
			log.Fatalf("could not create a property type based on name of property: %v", propertyName)
		}

	}

	return propertyDict

}

func parseIntRange(window string) IntPair {
	single, check := strconv.Atoi(window)

	if check == nil {
		return IntPair{first: single, second: single}
	}

	parts := strings.Split(window, "-")

	if len(parts) != 2 {
		log.Fatalf("statemachinebuilder/parseIntRange: range provided in statemachine yaml is invalid")
	}
	first, ok := strconv.Atoi(parts[0])
	if ok != nil {
		log.Fatalf("statemachinebuilder/parseIntRange: could not convert range string to int")
	}

	second, ok2 := strconv.Atoi(parts[1])
	if ok2 != nil {
		log.Fatalf("statemachinebuilder/parseIntRange: could not convert range string to int")
	}

	if second < first {
		log.Fatalf("statemachinebuilder/parseIntRange: invalid range inputted in yaml")
	}

	pair := IntPair{}
	pair.first = first
	pair.second = second

	return pair
}

func parseDoubleRange(doubleRange string) DoublePair {
	single, check := strconv.ParseFloat(doubleRange, 64)

	if check == nil {
		return DoublePair{first: single, second: single}
	}

	parts := strings.Split(doubleRange, "-")

	if len(parts) != 2 {
		log.Fatalf("statemachinebuilder/parseDoubleRange: range provided in statemachine yaml is invalid")
	}
	first, ok := strconv.ParseFloat(parts[0], 64)
	if ok != nil {
		log.Fatalf("statemachinebuilder/parseDoubleRange: could not convert range string to double")
	}

	second, ok2 := strconv.ParseFloat(parts[1], 64)
	if ok2 != nil {
		log.Fatalf("statemachinebuilder/parseDoubleRange: could not convert range string to double")
	}

	pair := DoublePair{}
	pair.first = first
	pair.second = second

	return pair
}

func parseBool(boolVal string) bool {
	val, ok := strconv.ParseBool(boolVal)

	if ok != nil {
		log.Fatalf("statemachinebuilder/parseBool: bool provided in statemachine yaml is invalid")
	}

	return val
}

func (s *StateMachineBuilder) convertYamlControls(yamlControls map[string][]RangeValueYaml) map[string]StringProperty {

	controlToState := map[string]StringProperty{}

	for control, rangeValList := range yamlControls {
		cancel := StringProperty{}
		timeline := []IntPair{}
		values := []string{}

		for _, rangeVal := range rangeValList {
			window := parseIntRange(rangeVal.Range)
			timeline = append(timeline, window)
			values = append(values, rangeVal.Value)
		}

		cancel.Init(control, timeline, values)

		controlToState[control] = cancel
	}

	return controlToState

}

func (s *StateMachineBuilder) Build(yamlFile string) *StateMachine {

	sm := StateMachine{}

	config := s.unfurl(yamlFile)

	stateDict := map[string]State{}

	for _, state := range config.States {
		stateToAdd := State{}
		stateToAdd.endingFrame = state.End
		stateToAdd.loopState = state.Loop
		stateToAdd.name = state.Name
		stateToAdd.properties = s.convertYamlProperties(state.Properties)
		stateToAdd.controlToState = s.convertYamlControls(state.ControlToState)

		stateDict[state.Name] = stateToAdd
	}

	sm.Init(stateDict)

	return &sm
}

//I need to have a parser that can read everything

//ok so I need to create a state machine that is keyed with states, each state will have a name, the name will be inside the
//state
