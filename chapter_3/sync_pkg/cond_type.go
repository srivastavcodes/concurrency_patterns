package main

import (
	"fmt"
	"sync"
	"time"
)

var IsWhat bool

func main() {
	sendAndReceiveSignalsUsingCond()
}

func sendAndReceiveSignalsUsingCond() {
	c := sync.NewCond(&sync.Mutex{})
	queue := make([]any, 0, 10)

	removeFromQueue := func(delay time.Duration) {
		time.Sleep(delay)

		c.L.Lock()
		queue = queue[1:]
		fmt.Println("removed from queue")
		c.L.Unlock()

		c.Signal() // signals to waiting go-routine
	}
	for i := 0; i < 10; i++ {
		c.L.Lock()
		for len(queue) == 2 {
			c.Wait()
		}
		fmt.Println("append to queue")
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock()
		fmt.Println(cap(queue))
	}
}

func replaceForLoopWithCond() {
	for !conditionTrue() { // consumes too much cpu cycle
		// send signal or something
	}
	// better alternative
	c := sync.NewCond(&sync.Mutex{})
	c.L.Lock()
	for !conditionTrue() {
		c.Wait()
	}
	c.L.Unlock()
}

func conditionTrue() bool { return IsWhat }
