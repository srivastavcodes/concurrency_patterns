package main

import (
	"fmt"
	"time"
)

func main() {
	forSelectExample()
}

func forSelectExample() {
	done := make(chan any)
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()
	workCounter := 0
loop:
	for {
		select {
		case <-done:
			break loop
		default:
		}
		// simulate real work
		workCounter++
		time.Sleep(1 * time.Second)
	}
	fmt.Printf("%v cycles of work before stopping.\n", workCounter)
}

func selectExample() {
	start := time.Now()

	ch := make(chan any)
	go func() {
		time.Sleep(5 * time.Second)
		close(ch)
	}()
	fmt.Println("blocking on ch")

	select {
	case <-ch:
		fmt.Printf("unblocked %v later.\n", time.Since(start))
	case <-time.After(1 * time.Second):
		fmt.Printf("timed out")
	}
}
