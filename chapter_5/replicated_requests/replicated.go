package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	replicatedRequests()
}

func replicatedRequests() {
	doWork := func(done <-chan any, id int, wg *sync.WaitGroup, result chan<- int) {
		started := time.Now()
		defer wg.Done()

		// simulate random load
		loadTime := time.Duration(1+rand.Intn(5)) * time.Second
		select {
		case <-done:
		case <-time.After(loadTime):
		}
		select {
		case <-done:
		case result <- id:
		}
		took := time.Since(started)
		// display how long handlers would take
		if took < loadTime {
			took = loadTime
		}
		fmt.Printf("%v took %d\n", id, took)
	}
	done := make(chan any)

	result := make(chan int)
	var wg sync.WaitGroup

	wg.Add(10)
	for i := 0; i < 10; i++ {
		go doWork(done, i, &wg, result)
	}
	firstReturned := <-result

	close(done)
	wg.Wait()

	fmt.Printf("received an answer from: #%v\n", firstReturned)
}
