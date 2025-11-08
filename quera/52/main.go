package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
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

type Server struct {
	portNumber int
	survey     *Survey
}

// ============ Request/Response Structs ============

type AddFlightRequest struct {
	Name string `json:"Name"`
}

type AddTicketRequest struct {
	FlightName    string `json:"FlightName"`
	PassengerName string `json:"PassengerName"`
}

type AddCommentRequest struct {
	FlightName    string `json:"FlightName"`
	PassengerName string `json:"PassengerName"`
	Score         int    `json:"Score"`
	Text          string `json:"Text"`
}

// ============ Survey Methods ============

func NewSurvey() *Survey {
	return &Survey{
		Flights:  make([]string, 0),
		Tickets:  make([]Ticket, 0),
		Comments: make(map[Ticket]Comment),
	}
}

func (s *Survey) AddFlight(flightName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// بررسی تکراری
	for _, f := range s.Flights {
		if f == flightName {
			return errors.New("flight already exists")
		}
	}

	s.Flights = append(s.Flights, flightName)
	return nil
}

func (s *Survey) AddTicket(flightName, passengerName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// بررسی وجود پرواز
	flightExists := false
	for _, f := range s.Flights {
		if f == flightName {
			flightExists = true
			break
		}
	}
	if !flightExists {
		return errors.New("flight does not exist")
	}

	// بررسی تکراری تیکت
	for _, t := range s.Tickets {
		if t.Passenger == passengerName && t.Flight == flightName {
			return errors.New("ticket already exists")
		}
	}

	s.Tickets = append(s.Tickets, Ticket{Passenger: passengerName, Flight: flightName})
	return nil
}

func (s *Survey) AddComment(flightName, passengerName string, score int, text string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// بررسی اعتبار امتیاز
	if score <= 0 || score > 10 {
		return errors.New("invalid score")
	}

	// بررسی وجود پرواز
	flightExists := false
	for _, f := range s.Flights {
		if f == flightName {
			flightExists = true
			break
		}
	}
	if !flightExists {
		return errors.New("flight does not exist")
	}

	// بررسی وجود تیکت
	ticketExists := false
	for _, t := range s.Tickets {
		if t.Passenger == passengerName && t.Flight == flightName {
			ticketExists = true
			break
		}
	}
	if !ticketExists {
		return errors.New("ticket does not exist")
	}

	// بررسی تکراری نظر
	ticket := Ticket{Passenger: passengerName, Flight: flightName}
	if _, exists := s.Comments[ticket]; exists {
		return errors.New("duplicate comment")
	}

	s.Comments[ticket] = Comment{Score: score, Text: text}
	return nil
}

func (s *Survey) GetCommentsAverage(flightName string) (float64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// بررسی وجود پرواز
	flightExists := false
	for _, f := range s.Flights {
		if f == flightName {
			flightExists = true
			break
		}
	}
	if !flightExists {
		return 0, errors.New("flight does not exist")
	}

	// محاسبه میانگین
	sum := 0
	count := 0
	for ticket, comment := range s.Comments {
		if ticket.Flight == flightName {
			sum += comment.Score
			count++
		}
	}

	if count == 0 {
		return 0, errors.New("no comments")
	}

	return float64(sum) / float64(count), nil
}

func (s *Survey) GetAllCommentsAverage() map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string]float64)
	flightScores := make(map[string][]int)

	// جمع‌آوری امتیازات
	for ticket, comment := range s.Comments {
		flightScores[ticket.Flight] = append(flightScores[ticket.Flight], comment.Score)
	}

	// محاسبه میانگین
	for flight, scores := range flightScores {
		if len(scores) > 0 {
			sum := 0
			for _, score := range scores {
				sum += score
			}
			result[flight] = float64(sum) / float64(len(scores))
		}
	}

	return result
}

func (s *Survey) GetComments(flightName string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// بررسی وجود پرواز
	flightExists := false
	for _, f := range s.Flights {
		if f == flightName {
			flightExists = true
			break
		}
	}
	if !flightExists {
		return []string{}, errors.New("flight does not exist")
	}

	// جمع‌آوری نظرات
	texts := make([]string, 0)
	for ticket, comment := range s.Comments {
		if ticket.Flight == flightName {
			texts = append(texts, comment.Text)
		}
	}

	return texts, nil
}

func (s *Survey) GetAllComments() map[string][]string {
	s.mu.Lock()
	defer s.mu.Unlock()

	result := make(map[string][]string)

	for flight := range s.Flights {
		result[s.Flights[flight]] = make([]string, 0)
	}

	for ticket, comment := range s.Comments {
		result[ticket.Flight] = append(result[ticket.Flight], comment.Text)
	}

	return result
}

// ============ Server Implementation ============

func NewServer(port int) *Server {
	return &Server{
		portNumber: port,
		survey:     NewSurvey(),
	}
}

func (server *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"Message": "OK"})
	})

	// Flights
	router.POST("/flights", server.handleAddFlight)

	// Tickets
	router.POST("/tickets", server.handleAddTicket)

	// Comments
	router.POST("/comments", server.handleAddComment)
	router.GET("/comments", server.handleGetAllComments)
	router.GET("/comments/:flightname", server.handleGetComments)

	listenAddress := fmt.Sprintf(":%d", server.portNumber)
	router.Run(listenAddress)
}

// ============ Handler Functions ============

func (server *Server) handleAddFlight(c *gin.Context) {
	var req AddFlightRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := server.survey.AddFlight(req.Name); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "OK"})
}

func (server *Server) handleAddTicket(c *gin.Context) {
	var req AddTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := server.survey.AddTicket(req.FlightName, req.PassengerName); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "OK"})
}

func (server *Server) handleAddComment(c *gin.Context) {
	var req AddCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	if err := server.survey.AddComment(req.FlightName, req.PassengerName, req.Score, req.Text); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"Message": "OK"})
}

func (server *Server) handleGetComments(c *gin.Context) {
	flightName := c.Param("flightname")
	average := c.Query("average")

	if average == "true" {
		avg, err := server.survey.GetCommentsAverage(flightName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Message": "OK",
			"Average": avg,
		})
	} else {
		texts, err := server.survey.GetComments(flightName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"Message": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Message": "OK",
			"Texts":   texts,
		})
	}
}

func (server *Server) handleGetAllComments(c *gin.Context) {
	average := c.Query("average")

	if average == "true" {
		averages := server.survey.GetAllCommentsAverage()
		c.JSON(http.StatusOK, gin.H{
			"Message":  "OK",
			"Averages": averages,
		})
	} else {
		allComments := server.survey.GetAllComments()
		c.JSON(http.StatusOK, gin.H{
			"Message": "OK",
			"Texts":   allComments,
		})
	}
}

func main() {
	port := 8080
	server := NewServer(port)
	server.Start()
}
