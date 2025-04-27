package main

import "fmt"

func main() {
	// var card string = "Ace of Spades"
	cards := []string{"A man from moon", newCard()}
	fmt.Println(cards)
	cards = append(cards, "After a bad day")
	fmt.Println(cards)

}

func newCard() string {
	return "Five of Dimonds"
}
