package main

import (
	"errors"
	"math"
)

type Comment struct {
	Score int
	Text  string
}

type Survey struct {
	flights  map[string]struct{}
	tickets  map[string]map[string]struct{}
	comments map[string]map[string]Comment
}

func NewSurvey() *Survey {
	return &Survey{
		flights:  make(map[string]struct{}),
		tickets:  make(map[string]map[string]struct{}),
		comments: make(map[string]map[string]Comment),
	}
}

// تابع کمکی برای گرد کردن به 2 رقم اعشار
func round2(x float64) float64 {
	return math.Round(x*100) / 100
}

func (s *Survey) AddFlight(flightName string) error {
	if _, exists := s.flights[flightName]; exists {
		return errors.New("Flight already exists")
	}
	s.flights[flightName] = struct{}{}
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	if _, exists := s.flights[flightName]; !exists {
		return errors.New("Flight does not exist")
	}
	if s.tickets[flightName] == nil {
		s.tickets[flightName] = make(map[string]struct{})
	}
	if _, exists := s.tickets[flightName][passengerName]; exists {
		return errors.New("Ticket already exists")
	}
	s.tickets[flightName][passengerName] = struct{}{}
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	if _, exists := s.flights[flightName]; !exists {
		return errors.New("Flight does not exist")
	}
	if s.tickets[flightName] == nil {
		return errors.New("Ticket does not exist")
	}
	if _, hasTicket := s.tickets[flightName][passengerName]; !hasTicket {
		return errors.New("Ticket does not exist")
	}
	if s.comments[flightName] == nil {
		s.comments[flightName] = make(map[string]Comment)
	}
	if _, exists := s.comments[flightName][passengerName]; exists {
		return errors.New("Comment already exists")
	}
	if comment.Score < 1 || comment.Score > 10 {
		return errors.New("Score is not supported")
	}
	s.comments[flightName][passengerName] = comment
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	if _, exists := s.flights[flightName]; !exists {
		return 0, errors.New("Flight does not exist")
	}
	commentsForFlight := s.comments[flightName]
	if len(commentsForFlight) == 0 {
		return 0, errors.New("No comments for this flight")
	}
	sum := 0
	count := 0
	for _, c := range commentsForFlight {
		sum += c.Score
		count++
	}
	avg := float64(sum) / float64(count)
	return round2(avg), nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	result := make(map[string]float64)

	for flightName := range s.flights {
		commentsForFlight := s.comments[flightName]
		if len(commentsForFlight) == 0 {
			continue // پرواز بدون کامنت در خروجی نمی‌اید
		}
		sum := 0
		count := 0
		for _, c := range commentsForFlight {
			sum += c.Score
			count++
		}
		if count > 0 {
			avg := float64(sum) / float64(count)
			result[flightName] = round2(avg)
		}
	}

	return result
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	if _, exists := s.flights[flightName]; !exists {
		return nil, errors.New("Flight does not exist")
	}
	commentsForFlight := s.comments[flightName]
	if len(commentsForFlight) == 0 {
		return []string{}, nil
	}
	result := make([]string, 0, len(commentsForFlight))
	for _, c := range commentsForFlight {
		result = append(result, c.Text)
	}
	return result, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	result := make(map[string][]string)

	// اینجا از کل پروازها عبور می‌کنیم
	for flightName := range s.flights {
		commentsForFlight := s.comments[flightName]
		if len(commentsForFlight) == 0 {
			// اگر کامنت نیست، اسلایس خالی بگذار
			result[flightName] = []string{}
		} else {
			texts := make([]string, 0, len(commentsForFlight))
			for _, c := range commentsForFlight {
				texts = append(texts, c.Text)
			}
			result[flightName] = texts
		}
	}

	return result
}
