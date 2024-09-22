package properties

//The properties package is amazing because it will have very key functions that will be used
//to generate properties and simultaneously

import "fmt"

type Property interface {
	Read(frame int) any
	GetName() string
	Alter(frame int, val any)
	Reset()
}

type IntPair struct {
	First  int
	Second int
}

func SearchTimeline(windowList []IntPair, frame int) int {

	for i, window := range windowList {
		if frame >= window.First && frame <= window.Second {
			return i
		}
	}

	return -1

}

type DoublePair struct {
	First  float64
	Second float64
}

type DoubleProperty struct { //Double property is the gold standard, make other properties based on this
	name          string
	timeline      []IntPair
	values        []DoublePair
	alteredValues []*DoublePair
}

func ValidateTimeline(timeline []IntPair) error { //This logic supports having a timeline of 5-5 5-7, which is perfect for all
	for i, window := range timeline {
		if window.First > window.Second {
			return fmt.Errorf("timeline is not ascending")
		}
		if i < len(timeline)-1 {
			if timeline[i].Second > timeline[i+1].First {
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

	index := SearchTimeline(s.timeline, frame)

	if index == -1 {
		return nil
	}

	defaultVal := s.values[index]

	if s.alteredValues[index] != nil {
		defaultVal = *s.alteredValues[index]
	}

	if defaultVal.First == defaultVal.Second {
		return defaultVal.First
	}

	slope := (defaultVal.Second - defaultVal.First) / (float64(s.timeline[index].Second - s.timeline[index].First))
	val := defaultVal.First + slope*float64(frame-s.timeline[index].First)

	return val
}

func (s *DoubleProperty) Alter(frame int, val any) {

	var valInterface interface{} = val
	dp, ok := valInterface.(*DoublePair)

	if !ok {
		return
	}

	index := SearchTimeline(s.timeline, frame)

	if index == -1 {
		return
	}

	s.alteredValues[index] = dp
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

	index := SearchTimeline(s.timeline, frame)

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

	index := SearchTimeline(s.timeline, frame)

	if index == -1 {
		return
	}

	s.alteredValues[index] = dp
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

	index := SearchTimeline(s.timeline, frame)

	if index == -1 {
		return nil
	}

	var defaultVal string = s.values[index]

	if s.alteredValues[index] != nil {
		defaultVal = *s.alteredValues[index]
	}

	return defaultVal
}

func (s *StringProperty) Alter(frame int, val any) {

	var valInterface interface{} = val

	dp, ok := valInterface.(*string)

	if !ok {
		return
	}

	index := SearchTimeline(s.timeline, frame)

	if index == -1 {
		return
	}

	s.alteredValues[index] = dp
}
