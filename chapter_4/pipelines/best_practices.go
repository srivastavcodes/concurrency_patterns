package main

import (
	"fmt"
	"time"
)

func main() {
	concurrentChannelPipelines()
}

func concurrentChannelPipelines() {
	generator := func(done <-chan any, integers ...int) <-chan int {
		instream := make(chan int, len(integers))
		go func() {
			defer close(instream)
			for _, value := range integers {
				select {
				case <-done:
					return
				case instream <- value:
				}
			}
		}()
		return instream
	}
	multiply := func(done <-chan any, instream <-chan int, multiplier int) <-chan int {
		multiplyStream := make(chan int)
		go func() {
			defer close(multiplyStream)
			for value := range instream {
				select {
				case <-done:
					return
				case multiplyStream <- value * multiplier:
				}
			}
		}()
		return multiplyStream
	}
	add := func(done <-chan any, instream <-chan int, additive int) <-chan int {
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for value := range instream {
				select {
				case <-done:
					return
				case addStream <- value + additive:
				}
			}
		}()
		return addStream
	}
	done := make(chan any)
	go func() {
		time.Sleep(time.Microsecond * 50)
		defer close(done)
	}()
	instream := generator(done, 1, 2, 3, 4)
	pipeline := multiply(done, add(done, multiply(done, instream, 2), 1), 2)

	for val := range pipeline {
		fmt.Printf("%d ", val)
	}
}
