package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	writeLeakFixed()
}

func writeLeakFixed() {
	randomStreamFunc := func(done <-chan any) <-chan int {
		randomStream := make(chan int)
		go func() {
			defer close(randomStream)
			defer fmt.Println("randomStreamFunc closure exited")
			for {
				select {
				case randomStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randomStream
	}
	done := make(chan any)
	randomStream := randomStreamFunc(done)

	fmt.Println("3 random ints:")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randomStream)
	}
	close(done)
	time.Sleep(time.Millisecond)
}

func writeLeak() {
	randStreamFunc := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("randStreamFunc closure exited")
			defer close(randStream)
			// forever blocks - no signal to exit
			for {
				randStream <- rand.Int()
			}
		}()
		return randStream
	}
	randStream := randStreamFunc()
	fmt.Println("3 random ints:")
	for i := 0; i < 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}
