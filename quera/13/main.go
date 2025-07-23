package main

import (
	"sync"
)

func readData(wg *sync.WaitGroup, out chan string, input chan string, decipherer func(encrypted string) string) {
	defer wg.Done()
	out <- decipherer(<-input)
}

func StartDecipher(senderChan chan string, decipherer func(encrypted string) string) chan string {
	out := make(chan string, 5)
	var wg sync.WaitGroup
	wg.Add(1)
	go readData(&wg, out, senderChan, decipherer)

	wg.Wait()

	return out
}
