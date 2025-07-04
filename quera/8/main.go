package main

import (
	"fmt"
	"math"
)

type FilterFunc func(int) bool
type MapperFunc func(int) int

func IsSquare(x int) bool {
	if x < 0 {
		return false
	}
	n := int(math.Sqrt(float64(x)))
	return n*n == x
}

func IsPalindrome(x int) bool {
	s := fmt.Sprintf("%d", x)
	if len(s) > 0 && s[0] == '-' {
		s = s[1:]
	}
	for i := 0; i < len(s)/2; i++ {
		if s[i] != s[len(s)-1-i] {
			return false
		}
	}
	return true
}

func Abs(num int) int {
	//TODO
	if num < 0 {
		return -num
	}
	return num
}

func Cube(num int) int {
	//TODO
	return num * num * num
}

func Filter(input []int, f FilterFunc) []int {
	result := make([]int, 0, len(input))
	for _, v := range input {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map(input []int, m MapperFunc) []int {
	//TODO
	var result []int
	for i := 0; i < len(input); i++ {
		result = append(result, m(input[i]))
	}
	return result
}
