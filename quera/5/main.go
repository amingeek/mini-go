package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	fmt.Println(n)
}
