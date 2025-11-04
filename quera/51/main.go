package main

import (
	"fmt"
	"sync"

	"github.com/go-sql-driver/mysql"
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
	Text          string `gorm:"type:longtext"`
}

func NewSurvey() *Survey {
	db := GetConnection()

	err := db.AutoMigrate(&Flight{}, &Ticket{}, &Comment{})
	if err != nil {
		panic(err)
	}

	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM tickets")
	db.Exec("DELETE FROM flights")

	return &Survey{
		db: db,
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

func (s *Survey) AddFlight(flightName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Transaction(func(tx *gorm.DB) error {
		var existing Flight
		if err := tx.Where("name = ?", flightName).First(&existing).Error; err == nil {
			return fmt.Errorf("flight already exists: %s", flightName)
		} else if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		flight := Flight{Name: flightName}
		if err := tx.Create(&flight).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				return fmt.Errorf("flight already exists: %s", flightName)
			}
			return fmt.Errorf("failed to create flight: %v", err)
		}
		return nil
	})
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var flight Flight
		if err := tx.Where("name = ?", flightName).First(&flight).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("flight does not exist: %s", flightName)
			}
			return err
		}

		ticket := Ticket{
			FlightName:    flightName,
			PassengerName: passengerName,
		}
		if err := tx.Create(&ticket).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				return fmt.Errorf("duplicate ticket for %s %s", flightName, passengerName)
			}
			return fmt.Errorf("failed to create ticket: %v", err)
		}
		return nil
	})
}

func (s *Survey) AddComment(flightName, passengerName string, comment Comment) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var flight Flight
		if err := tx.Where("name = ?", flightName).First(&flight).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("flight does not exist: %s", flightName)
			}
			return err
		}

		var ticket Ticket
		if err := tx.Where("flight_name = ? AND passenger_name = ?", flightName, passengerName).First(&ticket).Error; err != nil {
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
		if err := tx.Create(&newComment).Error; err != nil {
			if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
				return fmt.Errorf("duplicate comment for %s %s", flightName, passengerName)
			}
			return fmt.Errorf("failed to create comment: %v", err)
		}

		return nil
	})
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	check, _ := s.CheckFlight(flightName)
	if !check {
		return 0, fmt.Errorf("flight does not exist: %s", flightName)
	}
	comments, err := s.GetFlightComments(flightName)
	if err != nil {
		return 0, err
	}
	return float64(AvgScore(comments)), nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	type result struct {
		flightName string
		avg        float64
	}

	var flights []Flight
	s.db.Find(&flights)

	results := make(map[string]float64)
	resultCh := make(chan result, len(flights))
	var wg sync.WaitGroup

	for _, f := range flights {
		wg.Add(1)
		flightName := f.Name
		go func(fn string) {
			defer wg.Done()
			avg, err := s.GetCommentsAverage(fn)
			if err != nil {
				avg = 0
			}
			resultCh <- result{flightName: fn, avg: avg}
		}(flightName)
	}

	wg.Wait()
	close(resultCh)

	for r := range resultCh {
		results[r.flightName] = r.avg
	}

	return results
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	check, _ := s.CheckFlight(flightName)
	if !check {
		return nil, fmt.Errorf("flight does not exist: %s", flightName)
	}
	comments, err := s.GetFlightComments(flightName)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for _, c := range comments {
		res = append(res, c.Text)
	}
	return res, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	type result struct {
		flightName string
		comments   []string
	}

	var flights []Flight
	s.db.Find(&flights)

	results := make(map[string][]string)
	resultCh := make(chan result, len(flights))
	var wg sync.WaitGroup

	for _, f := range flights {
		wg.Add(1)
		flightName := f.Name
		go func(fn string) {
			defer wg.Done()
			comments, err := s.GetComments(fn)
			if err != nil {
				comments = []string{}
			}
			resultCh <- result{flightName: fn, comments: comments}
		}(flightName)
	}

	wg.Wait()
	close(resultCh)

	for r := range resultCh {
		results[r.flightName] = r.comments
	}

	return results
}
