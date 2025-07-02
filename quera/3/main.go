package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func sumSlice(numbers []int) int {
	total := 0
	for _, num := range numbers {
		total += num
	}
	return total / len(numbers)
}

func check_score(n int) string {
	if n >= 80 {
		return "Excellent"
	} else if n >= 60 && n < 80 {
		return "Very Good"
	} else if n >= 40 && n < 60 {
		return "Good"
	} else {
		return "Fair"
	}
}

func main() {
	var names []string
	var numbers []int
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		var t2 []int

		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())

		scanner.Scan()
		words := strings.Fields(scanner.Text())

		for _, word := range words {
			num, err := strconv.Atoi(word)
			if err != nil {
				panic(err)
			}
			t2 = append(t2, num)
		}

		names = append(names, name)
		numbers = append(numbers, sumSlice(t2))
	}
	for i := 0; i < len(names); i++ {
		fmt.Println(names[i], check_score(numbers[i]))
	}

}
