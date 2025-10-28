package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

// ================= STRUCTS =================

type ticket struct {
	Passenger string
	Flight    string
}

type comment struct {
	Passenger string
	Flight    string
	Score     int
	Message   string
}

type average struct {
	Flight string
	Value  float64
}

// ================= GLOBALS =================

var (
	flights  []string
	tickets  []ticket
	comments []comment
)

// ================= UTILITIES =================

func getInputInt(scanner *bufio.Scanner) int {
	scanner.Scan()
	num, _ := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	return num
}
func getList(scanner *bufio.Scanner, n int) []string {
	list := []string{}
	for i := 0; i < n; i++ {
		scanner.Scan()
		item := strings.TrimSpace(scanner.Text())
		list = append(list, item)
	}
	return list
}
func search(str string, ref []string) bool {
	for _, v := range ref {
		if v == str {
			return true
		}
	}
	return false
}

func countTickets(p string, f string, ref []ticket) int {
	count := 0
	for _, v := range ref {
		if v.Passenger == p && v.Flight == f {
			count++
		}
	}
	return count
}

func countComments(p string, f string, ref []comment) int {
	count := 0
	for _, v := range ref {
		if v.Passenger == p && v.Flight == f {
			count++
		}
	}
	return count
}

// ================== VALIDATIONS ==================

func checkFlight(f string) bool {
	return search(f, flights)
}

func statusTicket(t ticket) string {
	if !checkFlight(t.Flight) {
		return fmt.Sprintf("Invalid flight %s", t.Flight)
	}
	if countTickets(t.Passenger, t.Flight, tickets) >= 1 {
		return fmt.Sprintf("Duplicate ticket for %s %s", t.Flight, t.Passenger)
	}
	return "Accepted ticket"
}

func statusComment(c comment) (string, int) {
	if !checkFlight(c.Flight) {
		return fmt.Sprintf("Invalid flight %s", c.Flight), 400
	}
	validTicket := false
	for _, t := range tickets {
		if t.Passenger == c.Passenger && t.Flight == c.Flight {
			validTicket = true
			break
		}
	}
	if !validTicket {
		return fmt.Sprintf("Invalid passenger for %s %s", c.Flight, c.Passenger), 400
	}
	if countComments(c.Passenger, c.Flight, comments) >= 1 {
		// <-- fixed: include "by" to match expected output exactly
		return fmt.Sprintf("Duplicate comment for %s by %s", c.Flight, c.Passenger), 400
	}
	return fmt.Sprintf("Accepted comment for %s by %s", c.Flight, c.Passenger), 200
}

// ================= INPUT HANDLERS =================

func getTickets(scanner *bufio.Scanner, n int) []ticket {
	for i := 0; i < n; i++ {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}
		t := ticket{Passenger: parts[0], Flight: parts[1]}
		msg := statusTicket(t)
		if msg != "Accepted ticket" {
			fmt.Println(msg)
			continue
		}
		tickets = append(tickets, t)
	}
	return tickets
}
func getComments(scanner *bufio.Scanner, n int) []comment {
	for i := 0; i < n; i++ {
		scanner.Scan()
		line := strings.TrimSpace(scanner.Text())
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		score, _ := strconv.Atoi(parts[2])
		msg := strings.Join(parts[3:], " ")
		c := comment{Passenger: parts[0], Flight: parts[1], Score: score, Message: msg}
		status, code := statusComment(c)
		fmt.Println(status)
		if code == 200 {
			comments = append(comments, c)
		}
	}
	return comments
}

// ================= AVERAGE CALCULATOR =================

func calculateAverages(c []comment) []average {
	scores := make(map[string][]int)

	for _, com := range c {
		scores[com.Flight] = append(scores[com.Flight], com.Score)
	}

	flightsSorted := make([]string, 0, len(scores))
	for f := range scores {
		flightsSorted = append(flightsSorted, f)
	}
	sort.Strings(flightsSorted)

	result := []average{}
	for _, f := range flightsSorted {
		total := 0
		for _, s := range scores[f] {
			total += s
		}
		avg := float64(total) / float64(len(scores[f]))
		result = append(result, average{Flight: f, Value: avg})
	}

	return result
}

func printAverages(avgs []average) {
	for _, a := range avgs {
		fmt.Printf("Average score for %s is %.2f\n", a.Flight, a.Value)
	}
}

func getAverage(c []comment) {
	avgs := calculateAverages(c)
	printAverages(avgs)
}

// ================= MAIN ENGINE =================

func engine() {
	scanner := bufio.NewScanner(os.Stdin)

	nFlights := getInputInt(scanner)
	flights = getList(scanner, nFlights)

	nTickets := getInputInt(scanner)
	getTickets(scanner, nTickets)

	nComments := getInputInt(scanner)
	getComments(scanner, nComments)

	getAverage(comments)
}

func main() {
	engine()
}
