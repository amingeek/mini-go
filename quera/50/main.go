package main

import "errors"

type Comment struct {
	// do not modify or remove these fields
	Passenger string
	Flight    string
	Score     int
	Text      string
	// but you can add anything you want
}
type flight struct {
	Name string
}
type average struct {
	Flight string
	Value  float64
}
type ticket struct {
	Passenger string
	Flight    flight
}
type Survey struct {
	comments []Comment
	tickets  []ticket
	flights  []flight
}

func CheckFlight(key string, f []flight) bool {
	for _, v := range f {
		if key == v.Name {
			return true
		}
	}
	return false
}
func NewSurvey() *Survey {
	return &Survey{}
}

func (s *Survey) AddFlight(flightName string) error {
	found := CheckFlight(flightName, s.flights)
	if found {
		return errors.New("Flight already exists")
	}
	s.flights = append(s.flights, flight{Name: flightName})
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	return 0, nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	return nil
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	return nil, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	return nil
}
