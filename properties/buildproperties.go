package properties

import (
	"log"
	"strconv"
	"strings"
)

//This is the file where we build properties based on a list of range values

type RangeValue struct {
	Range string `yaml:"range"`
	Value string `yaml:"value"`
}

func BuildDoubleProperty(name string, rangeVals []RangeValue) *DoubleProperty {
	dProperty := DoubleProperty{}
	values := []DoublePair{}

	timeline := BuildTimeline(rangeVals)

	for _, rangeVal := range rangeVals {
		val := rangeVal.Value

		dp := ParseDoubleRange(val)
		values = append(values, dp)
	}

	err := dProperty.Init(name, timeline, values)

	if err != nil {
		return nil
	}

	return &dProperty

}

func BuildStringProperty(name string, rangeVals []RangeValue) *StringProperty {
	sProperty := StringProperty{}
	values := []string{}

	timeline := BuildTimeline(rangeVals)

	for _, rangeVal := range rangeVals {
		val := rangeVal.Value

		values = append(values, val)
	}

	err := sProperty.Init(name, timeline, values)

	if err != nil {
		return nil
	}

	return &sProperty
}

func BuildBoolProperty(name string, rangeVals []RangeValue) *BoolProperty {
	bProperty := BoolProperty{}
	values := []bool{}

	timeline := BuildTimeline(rangeVals)

	for _, rangeVal := range rangeVals {
		val := rangeVal.Value

		bVal := ParseBool(val)

		values = append(values, bVal)
	}

	err := bProperty.Init(name, timeline, values)

	if err != nil {
		return nil
	}

	return &bProperty

}

func ParseIntRange(window string) IntPair { //make this support negative integers too

	if window == "" {
		log.Fatalf("parseintrange: cannot parse an empty string")
	}

	//Check if an individual int is in the range
	val, ok := strconv.Atoi(window)

	if ok == nil {
		return IntPair{First: val, Second: val}
	}

	splitIndex := -1

	if window[0] == '-' {
		splitIndex = findSecondOccurence(window, "-")
	} else {
		splitIndex = strings.Index(window, "-")
	}

	if splitIndex == -1 {
		log.Fatalf("parseintrange: could not find split index")
	}

	first := window[:splitIndex]
	second := window[splitIndex+1:]

	firstInt, ok := strconv.Atoi(first)
	if ok != nil {
		log.Fatalf("ParseIntRange: could not convert range string to int")
	}

	secondInt, ok := strconv.Atoi(second)
	if ok != nil {
		log.Fatalf("ParseIntRange: could not convert range string to int")
	}

	if second < first {
		log.Fatalf("ParseIntRange: invalid range inputted in yaml")
	}

	pair := IntPair{}
	pair.First = firstInt
	pair.Second = secondInt

	return pair
}

func findSecondOccurence(s string, a string) int {
	first := strings.Index(s, a)

	if first == -1 {
		return -1
	}

	second := strings.Index(s[first+1:], a)

	return first + second + 1
}

func ParseDoubleRange(doubleRange string) DoublePair {
	if doubleRange == "" {
		log.Fatalf("parseDoubleRange: cannot parse an empty string")
	}
	//check to see if the string is a single range value
	val, ok := strconv.ParseFloat(doubleRange, 64)

	if ok == nil {
		return DoublePair{First: val, Second: val}
	}

	splitIndex := -1

	if doubleRange[0] == '-' {
		splitIndex = findSecondOccurence(doubleRange, "-")
	} else {
		splitIndex = strings.Index(doubleRange, "-")
	}

	if splitIndex == -1 {
		log.Fatalf("parseDoubleRange: could not find an index to split the expression")
	}

	first := doubleRange[:splitIndex]
	second := doubleRange[splitIndex+1:]

	firstDouble, ok := strconv.ParseFloat(first, 64)
	if ok != nil {
		log.Fatalf("parseDoubleRange: could not convert first part of rangeval to double")
	}

	secondDouble, ok := strconv.ParseFloat(second, 64)
	if ok != nil {
		log.Fatalf("parseDoubleRange: could not convert second part of rangeval to double")
	}

	pair := DoublePair{}
	pair.First = firstDouble
	pair.Second = secondDouble

	return pair
}

func ParseBool(boolVal string) bool {
	val, ok := strconv.ParseBool(boolVal)

	if ok != nil {
		log.Fatalf("parseBool: bool provided is invalid")
	}

	return val
}

func BuildTimeline(rangeVals []RangeValue) []IntPair {
	timeline := []IntPair{}
	for _, rangeVal := range rangeVals {
		window := ParseIntRange(rangeVal.Range)
		timeline = append(timeline, window)
	}

	return timeline
}

func BuildRangeVals(rangeList []string, valList []string) []RangeValue {

	rangeVals := []RangeValue{}

	for i := range rangeList {
		rangeVal := RangeValue{Range: rangeList[i], Value: valList[i]}
		rangeVals = append(rangeVals, rangeVal)
	}

	return rangeVals
}
