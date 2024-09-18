package controllers

import "strings"

type ControlList struct {
	controls []string
	iter     int
}

func (s *ControlList) GetInputs() string {
	s.iter++
	if s.iter >= len(s.controls) {
		s.iter = 0
	}
	control := s.controls[s.iter%len(s.controls)]
	return control
}

func (s *ControlList) GetDirection() string {
	control := s.controls[s.iter%len(s.controls)]
	if control == "left" || control == "right" || control == "up" || control == "down" {
		return control
	}

	return ""
}

func (s *ControlList) GetSpecial() bool {
	control := s.controls[s.iter%len(s.controls)]
	return strings.Contains(control, "special")
}

func (s *ControlList) GetNormal() bool {
	control := s.controls[s.iter%len(s.controls)]
	return strings.Contains(control, "normal")
}

func (s *ControlList) GetShield() bool {
	control := s.controls[s.iter%len(s.controls)]

	return control == "shield"
}

func (s *ControlList) Init(controls []string) {
	s.iter = -1
	s.controls = controls
}
