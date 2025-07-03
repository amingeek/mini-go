package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var names []string
	// var numbers []int
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		var nums []int
		var fasele []int

		scanner.Scan()
		list := strings.Fields(scanner.Text())
		names = append(names, list[0])

		for k := 1; k < len(list); k++ {

			num, err := strconv.Atoi(list[k])
			if err != nil {
				panic(err)
			}

			nums = append(nums, num)
		}

		for k := 1; k < len(nums); k++ {
			fasele = append(fasele, nums[k]-nums[k-1])
		}

		fmt.Println(names, nums, fasele)
	}
}
