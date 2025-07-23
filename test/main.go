package main

import (
	"fmt"
	"sync"
	"time"
)

func say(text string, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	for i := 0; i < 3; i++ {
		fmt.Println(text)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go say("سلام از goroutine", &wg)

	say("سلام از main", nil)

	wg.Wait()
}
