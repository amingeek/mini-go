package main

// Create a deck

type deck []string

func (d deck) print() {

	for i, card := range d {
		println(i, card)
	}
}
