package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"unicode"
)

func tavan(n int) string {
	st := strconv.Itoa(n)
	x := 0
	for i := 0; i < len(st); i++ {
		b, err := strconv.Atoi(string(st[i]))
		if err == nil {
			x += int(math.Pow(float64(b), float64(len(st))))
		}
	}

	if x == n {
		return "YES"
	}

	return "NO"
}

func findNumbers(s string) int {
	n := 0
	current := ""

	for _, v := range s {
		if unicode.IsDigit(v) {
			current += string(v)
		} else if !unicode.IsDigit(v) {
			b, err := strconv.Atoi(current)
			if err == nil {
				n += b
				current = ""
			}
		}
	}
	if current != "" {
		b, err := strconv.Atoi(current)
		if err == nil {
			n += b
			current = ""
		}
	}

	return n

}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	text := scanner.Text()

	fmt.Println(tavan(findNumbers(text)))

}
