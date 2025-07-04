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
	var numbers []int
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	n, err := strconv.Atoi(scanner.Text())
	if err != nil {
		panic(err)
	}

	for i := 0; i < n; i++ {
		var nums []int
		run := 0
		count := 0

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

		for j := 2; j < len(nums); j++ {

			if nums[j]-nums[j-1] == nums[j-1]-nums[j-2] {
				run++
				count += run
			} else {
				run = 0
			}

		}

		numbers = append(numbers, count)

	}

	for i := 0; i < len(names); i++ {
		fmt.Println(names[i], numbers[i])
	}
}
