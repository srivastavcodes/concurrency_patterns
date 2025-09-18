package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sendAtWorkStart() {
	doWork := func(done <-chan any) (<-chan any, <-chan int) {
		heartbeatStream := make(chan any, 1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for i := 0; i < 10; i++ {
				select {
				case heartbeatStream <- struct{}{}:
				default:
				}
				select {
				case <-done:
					return
				case workStream <- rand.Intn(10):
				}
			}
		}()
		return heartbeatStream, workStream
	}
	done := make(chan any)
	go func() {
		time.Sleep(time.Second * 3)
		close(done)
	}()
	heartbeat, results := doWork(done)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case value, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("result: %v\n", value)
		}
	}
}
