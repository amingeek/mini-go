package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
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

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	n := 0
	scanner.Scan()
	// n, err := strconv.Atoi(scanner.Text())
	text := scanner.Text()
	for i := 0; i < len(text); i++ {
		b, err := strconv.Atoi(string(text[i]))
		if err == nil {
			n += b
		}
	}

	fmt.Println(tavan(n))

}
