package main

import (
	"fmt"
)

func main() {
	var income int
	maliat := 0
	fmt.Scanf("%d", &income)

	if income <= 100 {
		maliat += income * 5 / 100
	}
	if income > 100 && income <= 400 {
		income -= 100
		maliat += income*10/100 + 5
	}
	if income > 500 && income <= 1000 {
		income -= 500
		maliat += income*15/100 + 45
	}
	if income > 1000 {
		income -= 1000
		maliat += income*20/100 + 120
	}

	fmt.Println(maliat)
}
