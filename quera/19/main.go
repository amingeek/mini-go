package main

import (
	"sync/atomic"
	"time"
)

type FutureResult struct {
	Done       atomic.Bool
	ResultChan chan string
	// TODO
}

type Task func() string

func Async(t Task) *FutureResult {
	ch := make(chan string, 1)
	f := &FutureResult{ResultChan: ch}

	go func() {
		tasks := t()
		f.ResultChan <- tasks
		f.Done.Store(true)
	}()

	return f

}
func AsyncWithTimeout(t Task, timeout time.Duration) *FutureResult {
	f := &FutureResult{ResultChan: make(chan string, 1)}

	go func() {
		resChan := make(chan string, 1)

		go func() {
			resChan <- t()
		}()

		timer := time.NewTimer(timeout)
		select {
		case res := <-resChan:
			f.ResultChan <- res
			f.Done.Store(true)
		case <-timer.C:
			f.ResultChan <- "timeout"
		}
	}()

	return f

}

func (fResult *FutureResult) Await() string {
	res := <-fResult.ResultChan
	return res
}

func CombineFutureResults(fResults ...*FutureResult) *FutureResult {
	combined := &FutureResult{
		ResultChan: make(chan string, len(fResults)),
	}

	go func() {
		for _, fr := range fResults {
			res := <-fr.ResultChan
			combined.ResultChan <- res
		}
		combined.Done.Store(true)
	}()
	return combined
}
