package main

import (
	"time"
)

func main() {
	example()
}

func example() {
	done := make(chan any)
	stringStream := make(chan string)

	go func() {
		time.Sleep(3 * time.Second)
		close(done)
	}()
	for _, val := range []string{"a", "b", "c"} {
		select {
		case <-done:
			return
		case stringStream <- val:
		}
	}
}
