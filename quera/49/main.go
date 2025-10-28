package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type ticket struct {
	Passenger string
	Flight    string
}

var (
	flights = []string{}
	tickets = []ticket{}
)

//func checkFlight() bool {
//	//
//}

func getList(n int) []string {
	// get list of anything
	list1 := []string{}
	scanner := bufio.NewScanner(os.Stdin)
	for i := 0; i < n; i++ {
		scanner.Scan()
		list1 = append(list1, scanner.Text())
	}
	return list1
}

func getInput() int {
	// get input from user and return int
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		panic(err)
	}
	res, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}
	return res
}

func engine() {
	// a func for run all func
	nFlights := getInput()
	flights := getList(nFlights)
	fmt.Println(flights)
	nTickets := getInput()
	ListOfTicket := getList(nTickets)
	fmt.Println(ListOfTicket)
}

func main() {
	// call engine and run program
	engine()
}
