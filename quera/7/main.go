package main

import (
	"fmt"
)

func convert(n int) string {
	if n >= 10 {
		return fmt.Sprintf("%d", n)
	} else if n < 10 {
		return fmt.Sprintf("0%d", n)
	}

	return ""
}

func ConvertToDigitalFormat(hour, minute, second int) string {
	// TODO
	return convert(hour) + ":" + convert(minute) + ":" + convert(second)
}

func ExtractTimeUnits(seconds int) (int, int, int) {
	// TODO
	h, m := 0, 0

	if seconds >= 3600 {
		h = seconds / 3600
		seconds = seconds - (h * 3600)
	}

	if seconds >= 60 {
		m = seconds / 60
		seconds = seconds - (m * 60)
	}

	return h, m, seconds
}
