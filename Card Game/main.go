package main

func main() {
	// var card string = "Ace of Spades"
	cards := deck{"A man from moon", newCard()}
	cards = append(cards, "After a bad day")

	cards.print()
}

func newCard() string {
	return "Five of Dimonds"
}
