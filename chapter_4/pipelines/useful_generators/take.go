package main

import (
	"fmt"
)

func main() {
	takeGenerator()
}

// works when paired with repeat, takes only the first (num) values from repeat
func takeGenerator() {
	// any should be replaced with the relevant data type
	take := func(done <-chan any, valueStream <-chan any, num int) <-chan any {
		takeStream := make(chan any)
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	repeat := func(done <-chan any, values ...any) <-chan any {
		valueStream := make(chan any)
		go func() {
			defer close(valueStream)
			for {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}
	done := make(chan any)
	defer close(done)

	for val := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%d ", val)
	}
}
