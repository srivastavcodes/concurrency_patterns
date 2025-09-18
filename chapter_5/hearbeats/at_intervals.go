package main

import (
	"fmt"
	"time"
)

func main() {
	simulateBadDesign()
}

func simulateBadDesign() {
	doWork := func(done <-chan any, pulseInterval time.Duration) (<-chan any, <-chan time.Time) {
		results := make(chan time.Time)
		heartbeat := make(chan any)
		go func() {
			pulse := time.Tick(pulseInterval)
			workgen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}
			for i := 0; i < 2; i++ {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workgen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}
	done := make(chan any)
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results: %d\n", r.Second())
		// no results above -> system realizes something is wrong -> EXITS!
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			return
		}
	}
}

func sendAtIntervals() {
	doWork := func(done <-chan any, pulseInterval time.Duration) (<-chan any, <-chan time.Time) {
		results := make(chan time.Time)
		heartbeat := make(chan any)
		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)
			workgen := time.Tick(2 * pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct{}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results <- r:
						return
					}
				}
			}
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workgen:
					sendResult(r)
				}
			}
		}()
		return heartbeat, results
	}
	done := make(chan any)
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results: %d\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}
}
