package main

import (
	"errors"
	"slices"
	"sync"
)

type Ticket struct {
	Passenger string
	Flight    string
}

type Comment struct {
	Score int
	Text  string
}

type Survey struct {
	Flights  []string
	Tickets  []Ticket
	Comments map[Ticket]Comment
	mu       sync.Mutex
}

func NewSurvey() *Survey {
	return &Survey{Flights: make([]string, 0), Tickets: make([]Ticket, 0), Comments: make(map[Ticket]Comment)}
}

func NewTicket(passengerName, flightName string) *Ticket {
	return &Ticket{
		Passenger: passengerName,
		Flight:    flightName,
	}
}

func (s *Survey) FlightCheck(flightName string) bool {
	return slices.Contains(s.Flights, flightName)
}

func (s *Survey) PassengerCheck(passengerName, flightName string) bool {
	for _, ticket := range s.Tickets {
		if ticket.Passenger == passengerName && ticket.Flight == flightName {
			return true
		}
	}
	return false
}

func (s *Survey) CommentCheck(passengerName, flightName string) bool {
	for ticket, comment := range s.Comments {
		if ticket.Passenger == passengerName && ticket.Flight == flightName && comment.Score != 0 {
			return true
		}
	}
	return false
}

func (s *Survey) AddFlight(flightName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.FlightCheck(flightName) {
		return errors.New("flight already exists")
	}
	s.Flights = append(s.Flights, flightName)
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.FlightCheck(flightName) {
		return errors.New("flight does not exist")
	}
	if s.PassengerCheck(passengerName, flightName) {
		return errors.New("ticket already exists")
	}
	s.Tickets = append(s.Tickets, *NewTicket(passengerName, flightName))
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if comment.Score <= 0 || comment.Score > 10 {
		return errors.New("invalid score")
	}
	if !s.FlightCheck(flightName) {
		return errors.New("flight does not exist")
	}
	if !s.PassengerCheck(passengerName, flightName) {
		return errors.New("ticket does not exist")
	}
	if s.CommentCheck(passengerName, flightName) {
		return errors.New("duplicate comment")
	}
	ticket := NewTicket(passengerName, flightName)
	s.Comments[*ticket] = comment
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	if s.FlightCheck(flightName) {
		rate := 0
		count := 0
		for key, val := range s.Comments {
			if key.Flight == flightName {
				rate += val.Score
				count++
			}
		}
		if count > 0 {
			return float64(rate) / float64(count), nil
		}
	}
	return 0, errors.New("flight not exists or no comment submitted for flight")
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	flightsAverage := make(map[string]float64)
	for _, flightName := range s.Flights {
		flightAverage, _ := s.GetCommentsAverage(flightName)
		if flightAverage != 0 {
			flightsAverage[flightName] = flightAverage
		}
	}
	return flightsAverage
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	comments := make([]string, 0)
	if !s.FlightCheck(flightName) {
		return comments, errors.New("flight does not exist")
	}
	for ticket, comment := range s.Comments {
		if ticket.Flight == flightName {
			comments = append(comments, comment.Text)
		}
	}
	return comments, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	flightComments := make(map[string][]string)
	for _, flightName := range s.Flights {
		flightComment, _ := s.GetComments(flightName)
		flightComments[flightName] = flightComment
	}
	return flightComments
}
