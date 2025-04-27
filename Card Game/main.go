package main

import "fmt"

func main() {
	// var card string = "Ace of Spades"
	cards := deck{"A man from moon", newCard()}
	fmt.Println(cards)
	cards = append(cards, "After a bad day")

	for i, card := range cards {
		fmt.Println(i, card)
	}

}

func newCard() string {
	return "Five of Dimonds"
}
