package statemachinery

import "fmt"

type Property interface {
	Read(frame int) any
	GetName() string
	Alter(frame int, val any)
	Reset()
}

type IntPair struct {
	first  int
	second int
}

type DoublePair struct {
	first  float64
	second float64
}

func NewDoublePair(f float64, s float64) *DoublePair {
	return &DoublePair{first: f, second: s}
}

type DoubleProperty struct { //Double property is the gold standard, make other properties based on this
	name          string
	timeline      []IntPair
	values        []DoublePair
	alteredValues []*DoublePair
}

func ValidateTimeline(timeline []IntPair) error {
	for i, window := range timeline {
		if window.first > window.second {
			return fmt.Errorf("timeline is not ascending")
		}
		if i < len(timeline)-1 {
			if timeline[i].second > timeline[i+1].first {
				return fmt.Errorf("timeline is not ascending")
			}
		}
	}

	return nil
}

func (s *DoubleProperty) Init(name string, timeline []IntPair, values []DoublePair) error {
	s.name = name
	s.timeline = timeline
	s.values = values

	//make sure that timeline is in ascending order
	err := ValidateTimeline(timeline)

	if err != nil {
		return err
	}

	if len(values) != len(timeline) {
		return fmt.Errorf("values is not same length as timeline")
	}

	s.alteredValues = make([]*DoublePair, len(values))

	return nil
}

func (s *DoubleProperty) GetName() string {
	return s.name
}

func (s *DoubleProperty) Read(frame int) any {

	index := -1

	for i, window := range s.timeline {
		if frame <= window.second && frame >= window.first {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	defaultVal := s.values[index]

	if s.alteredValues[index] != nil {
		defaultVal = *s.alteredValues[index]
	}

	if defaultVal.first == defaultVal.second {
		return defaultVal.first
	}

	slope := (defaultVal.second - defaultVal.first) / (float64(s.timeline[index].second - s.timeline[index].first))
	val := defaultVal.first + slope*float64(frame-s.timeline[index].first)

	return val
}

func (s *DoubleProperty) Alter(frame int, val any) {

	var valInterface interface{} = val
	dp, ok := valInterface.(*DoublePair)

	if !ok {
		return
	}

	for i, window := range s.timeline {
		if frame >= window.first && frame <= window.second {
			s.alteredValues[i] = dp
		}
	}
}

func (s *DoubleProperty) Reset() {
	for i := range s.alteredValues {
		s.alteredValues[i] = nil
	}
}

type BoolProperty struct {
	name          string
	timeline      []IntPair
	values        []bool
	alteredValues []*bool
}

func (s *BoolProperty) GetName() string {
	return s.name
}

func (s *BoolProperty) Init(name string, timeline []IntPair, values []bool) error {
	s.name = name
	s.timeline = timeline
	s.values = values

	//make sure that timeline is in ascending order
	err := ValidateTimeline(timeline)

	if err != nil {
		return err
	}

	if len(values) != len(timeline) {
		return fmt.Errorf("values is not same length as timeline")
	}

	s.alteredValues = make([]*bool, len(values))

	return nil
}

func (s *BoolProperty) Read(frame int) any {

	index := -1

	for i, window := range s.timeline {
		if frame <= window.second && frame >= window.first {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	var defaultVal bool = s.values[index]

	if s.alteredValues[index] != nil {
		defaultVal = *s.alteredValues[index]
	}

	return defaultVal
}

func (s *BoolProperty) Alter(frame int, val any) {

	var valInterface interface{} = val

	dp, ok := valInterface.(*bool)

	if !ok {
		return
	}

	for i, window := range s.timeline {
		if frame >= window.first && frame <= window.second {
			s.alteredValues[i] = dp
		}
	}
}

func (s *BoolProperty) Reset() {
	for i := range s.alteredValues {
		s.alteredValues[i] = nil
	}
}

type StringProperty struct {
	name          string
	timeline      []IntPair
	values        []string
	alteredValues []*string
}

func (s *StringProperty) GetName() string {
	return s.name
}

func (s *StringProperty) Init(name string, timeline []IntPair, values []string) error {
	s.name = name
	s.timeline = timeline
	s.values = values

	if len(timeline) != len(values) {
		return fmt.Errorf("values length is not same as timeline")
	}

	err := ValidateTimeline(timeline)

	if err != nil {
		return err
	}

	s.alteredValues = make([]*string, len(values))

	return nil
}

func (s *StringProperty) Read(frame int) any {

	index := -1

	for i, window := range s.timeline {
		if frame <= window.second && frame >= window.first {
			index = i
		}
	}

	if index == -1 {
		return nil
	}

	var defaultVal string = s.values[index]

	if s.alteredValues[index] != nil {
		defaultVal = *s.alteredValues[index]
	}

	if frame == s.timeline[index].second {
		s.alteredValues[index] = nil
	}

	return defaultVal
}

func (s *StringProperty) Alter(frame int, val any) {

	var valInterface interface{} = val

	dp, ok := valInterface.(*string)

	if !ok {
		return
	}

	for i, window := range s.timeline {
		if frame >= window.first && frame <= window.second {
			s.alteredValues[i] = dp
		}
	}
}

//for the alter method, you need to put the * to the value you want to alter

//The property list will be part of a state, as states will have access to changing the character's properties
//the properties that are seen by the state are later applied to the character

//Create many different types of property types that can be read, such as bool property, slide double property, constant property
//

//So the recap is: What if we want to control the xspeed and yspeed of a character based on certain frames of the state?
//We would need a system to be  able to model graphs of double values for properties - well any value
//so, the solution is, create an object called a property, and it has the ability to read any value

//the different types of properties will just return their types in interface form, so they will still
//in the future, hitboxes and hurtboxes can also have properties as well, such as intangibility and

//timelines for the properties should be inclusive inclusive
//if you have a frame windows 1-5 6-7..., you can represent it as the start frame of every window, 1 6
//if there is no n+1 to compute the window, then use the end frame
//you cannot use this method to represent gaps, where reading is nil
