package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var numbers []string
	var check string
	var nn int
	count := make(map[string]string)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		scanner.Scan()
		countrys := strings.Fields(scanner.Text())
		count[countrys[1]] = countrys[0]
	}

	scanner.Scan()
	nn, err = strconv.Atoi(scanner.Text())

	if err != nil {
		panic(err)
	}

	for i := 0; i < nn; i++ {
		scanner.Scan()
		numbers = append(numbers, scanner.Text())
	}

	for i := 0; i < len(numbers); i++ {
		check = numbers[i][:3]
		value, exists := count[check]
		if exists {
			fmt.Println(value)
		} else {
			fmt.Println("Invalid Number")
		}
	}
}
