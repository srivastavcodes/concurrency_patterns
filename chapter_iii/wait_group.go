package main

import (
	"fmt"
	"sync"
	"time"
)

func waitGroup1() {
	var wg sync.WaitGroup

	wg.Add(1) // adds that a goroutine has been added
	go func() {
		defer wg.Done() // defers to make sure we let WaitGroup know we've exited
		fmt.Println("1st goroutine is sleeping")
		time.Sleep(1)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine is sleeping")
		time.Sleep(2)
	}()
	wg.Wait() // blocks the main goroutine till all the goroutines indicate they've exited
	fmt.Println("All goroutines complete")
}

func waitGroup2() {
	concurrentHello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello from %v!\n", id)
	}
	const numGreeters = 5

	var wg sync.WaitGroup
	wg.Add(numGreeters)

	for i := 0; i < numGreeters; i++ {
		go concurrentHello(&wg, i+1)
	}
	wg.Wait()
}
