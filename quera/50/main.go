package main

import (
	"errors"
)

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

func CheckTicket(FlightName string, passengerName string, ticket []ticket) bool {
	for _, v := range ticket {
		if passengerName == v.Passenger && FlightName == v.Flight.Name {
			return true
		}
	}
	return false
}

func CheckComment(FlightName string, passengerName string, comments []Comment) bool {
	for _, v := range comments {
		if v.Passenger == passengerName && v.Flight == FlightName {
			return true
		}
	}
	return false
}

func CheckScore(n int, minN int, maxN int) bool {
	return n >= minN && n <= maxN
}

func GetAverage(FlightName string, comments []Comment) average {
	if len(comments) == 0 {
		return average{Flight: FlightName, Value: 0}
	}
	n := 0
	for _, v := range comments {
		n += v.Score
	}
	avg := float64(n) / float64(len(comments))
	return average{Flight: FlightName, Value: avg}
}

func GetAllAverage(comments []Comment) []average {
	if len(comments) == 0 {
		return []average{}
	}

	uniqueFlights := map[string]bool{}
	for _, c := range comments {
		uniqueFlights[c.Flight] = true
	}

	var averages []average

	for flightName := range uniqueFlights {
		var flightComments []Comment
		for _, c := range comments {
			if c.Flight == flightName {
				flightComments = append(flightComments, c)
			}
		}

		avg := GetAverage(flightName, flightComments)
		averages = append(averages, avg)
	}

	return averages
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
	FoundFlight := CheckFlight(flightName, s.flights)
	if !FoundFlight {
		return errors.New("Flight does not exist")
	}
	FoundTicket := CheckTicket(flightName, passengerName, s.tickets)
	if FoundTicket {
		return errors.New("Ticket already exists")
	}
	s.tickets = append(s.tickets, ticket{
		Passenger: passengerName,
		Flight:    flight{Name: flightName},
	})
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	FoundFlight := CheckFlight(flightName, s.flights)
	if !FoundFlight {
		return errors.New("Flight does not exist")
	}
	FoundTicket := CheckTicket(flightName, passengerName, s.tickets)
	if !FoundTicket {
		return errors.New("Ticket does not exist")
	}
	FoundComment := CheckComment(flightName, passengerName, s.comments)
	if FoundComment {
		return errors.New("Comment already exists")
	}
	FoundScore := CheckScore(comment.Score, 1, 10)
	if !FoundScore {
		return errors.New("Score is not supported")
	}

	comment.Flight = flightName
	comment.Passenger = passengerName
	comment.Flight = flightName

	s.comments = append(s.comments, comment)
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	if !CheckFlight(flightName, s.flights) {
		return 0, errors.New("Flight does not exist")
	}

	var flightComments []Comment
	for _, c := range s.comments {
		if c.Flight == flightName {
			flightComments = append(flightComments, c)
		}
	}

	avg := GetAverage(flightName, flightComments)
	return avg.Value, nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	result := make(map[string]float64)
	allAvg := GetAllAverage(s.comments)
	for _, a := range allAvg {
		result[a.Flight] = a.Value
	}
	return result
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	if !CheckFlight(flightName, s.flights) {
		return nil, errors.New("Flight does not exist")
	}

	var texts []string
	for _, c := range s.comments {
		if c.Flight == flightName {
			texts = append(texts, c.Text)
		}
	}

	return texts, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	result := make(map[string][]string)
	for _, c := range s.comments {
		result[c.Flight] = append(result[c.Flight], c.Text)
	}
	return result
}
