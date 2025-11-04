package main

import (
	"fmt"
	"strings"
	"sync"

	"gorm.io/gorm"
)

type Survey struct {
	db *gorm.DB
	mu sync.RWMutex
}

type Flight struct {
	gorm.Model
	Name     string    `gorm:"type:varchar(255);unique;not null"`
	Tickets  []Ticket  `gorm:"foreignKey:FlightName;references:Name;constraint:OnDelete:CASCADE"`
	Comments []Comment `gorm:"foreignKey:FlightName;references:Name;constraint:OnDelete:CASCADE"`
}

type Ticket struct {
	gorm.Model
	FlightName    string `gorm:"type:varchar(255);not null;uniqueIndex:idx_flight_passenger"`
	PassengerName string `gorm:"type:varchar(255);not null;uniqueIndex:idx_flight_passenger"`
}

type Comment struct {
	gorm.Model
	FlightName    string `gorm:"type:varchar(255);not null;uniqueIndex:idx_comment_unique"`
	PassengerName string `gorm:"type:varchar(255);not null;uniqueIndex:idx_comment_unique"`
	Score         int    `gorm:"not null"`
	Text          string `gorm:"type:text"`
}

func NewSurvey() *Survey {
	db := GetConnection()

	// Clean up existing data
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM tickets")
	db.Exec("DELETE FROM flights")

	err := db.AutoMigrate(&Flight{}, &Ticket{}, &Comment{})
	if err != nil {
		panic(err)
	}

	return &Survey{
		db: db,
		mu: sync.RWMutex{},
	}
}

func (s *Survey) CheckFlight(flightName string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var existing Flight
	err := s.db.Where("name = ?", flightName).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Survey) CheckPassenger(flightName, passengerName string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var existing Ticket
	err := s.db.Where("flight_name = ? AND passenger_name = ?", flightName, passengerName).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Survey) CheckComment(flightName, passengerName string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var existing Comment
	err := s.db.Where("flight_name = ? AND passenger_name = ?", flightName, passengerName).First(&existing).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (s *Survey) GetFlightComments(flightName string) ([]Comment, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var comments []Comment
	res := s.db.Where("flight_name = ?", flightName).Find(&comments)
	if res.Error != nil {
		return nil, res.Error
	}
	return comments, nil
}

func AvgScore(c []Comment) int {
	if len(c) == 0 {
		return 0
	}
	n := 0
	for i := 0; i < len(c); i++ {
		n += c[i].Score
	}
	return n / len(c)
}

func isUniqueConstraintError(err error) bool {
	return strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key")
}

func (s *Survey) AddFlight(flightName string) error {

	flight := Flight{Name: flightName}
	if err := s.db.Create(&flight).Error; err != nil {
		if isUniqueConstraintError(err) {
			return fmt.Errorf("flight already exists: %s", flightName)
		}
		return err
	}
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	var flight Flight
	if err := s.db.Where("name = ?", flightName).First(&flight).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("flight does not exist: %s", flightName)
		}
		return err
	}

	ticket := Ticket{
		FlightName:    flightName,
		PassengerName: passengerName,
	}
	if err := s.db.Create(&ticket).Error; err != nil {
		if isUniqueConstraintError(err) {
			return fmt.Errorf("duplicate ticket for %s %s", flightName, passengerName)
		}
		return err
	}
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {

	var flight Flight
	if err := s.db.Where("name = ?", flightName).First(&flight).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("flight does not exist: %s", flightName)
		}
		return err
	}

	var ticket Ticket
	if err := s.db.Where("flight_name = ? AND passenger_name = ?", flightName, passengerName).First(&ticket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("invalid passenger for %s %s", flightName, passengerName)
		}
		return err
	}

	newComment := Comment{
		FlightName:    flightName,
		PassengerName: passengerName,
		Score:         comment.Score,
		Text:          comment.Text,
	}
	if err := s.db.Create(&newComment).Error; err != nil {
		if isUniqueConstraintError(err) {
			return fmt.Errorf("duplicate comment for %s %s", flightName, passengerName)
		}
		return err
	}
	return nil
}
func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if flight exists
	var flight Flight
	if err := s.db.Where("name = ?", flightName).First(&flight).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, fmt.Errorf("flight does not exist: %s", flightName)
		}
		return 0, err
	}

	// Get comments
	var comments []Comment
	if err := s.db.Where("flight_name = ?", flightName).Find(&comments).Error; err != nil {
		return 0, err
	}

	return float64(AvgScore(comments)), nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var flights []Flight
	s.db.Find(&flights)

	results := make(map[string]float64)

	for _, f := range flights {
		flightName := f.Name
		comments, err := s.GetFlightComments(flightName)
		if err != nil {
			results[flightName] = 0
			continue
		}
		results[flightName] = float64(AvgScore(comments))
	}

	return results
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if flight exists
	var flight Flight
	if err := s.db.Where("name = ?", flightName).First(&flight).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("flight does not exist: %s", flightName)
		}
		return nil, err
	}

	// Get comments
	var comments []Comment
	if err := s.db.Where("flight_name = ?", flightName).Find(&comments).Error; err != nil {
		return nil, err
	}

	res := make([]string, 0, len(comments))
	for _, c := range comments {
		res = append(res, c.Text)
	}
	return res, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var flights []Flight
	s.db.Find(&flights)

	results := make(map[string][]string)

	for _, f := range flights {
		flightName := f.Name
		comments, err := s.GetComments(flightName)
		if err != nil {
			results[flightName] = []string{}
			continue
		}
		results[flightName] = comments
	}

	return results
}
