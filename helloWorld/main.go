package main

import (
	"fmt"
)

func main() {
	var inumber float64
	var fnumber int
	fmt.Scanf("%f %d", &inumber, &fnumber)
	if int(inumber) >= 2 && fnumber >= 2 && int(inumber) <= 100 && fnumber <= 100000 {
		for i := 1; i <= fnumber; i++ {
			if i%int(inumber) == 0 {
				counter := i / int(inumber)
				for k := 0; k < counter; k++ {
					fmt.Print("Hope ")
				}
				fmt.Print("\n")
			} else {
				fmt.Println(i)
			}
		}
	} else {
		fmt.Println("Invalid input. Please ensure 2 <= inumber <= 100 and fnumber <= 100000.")
	}
}
