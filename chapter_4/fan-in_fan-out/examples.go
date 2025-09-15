package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
)

func main() {

}

func FanIn(done <-chan any, channels ...<-chan any) <-chan any {
	var wg sync.WaitGroup
	multiplexedStream := make(chan any)

	multiplex := func(ch <-chan any) {
		defer wg.Done()
		for val := range ch {
			select {
			case <-done:
				return
			case multiplexedStream <- val:
			}
		}
	}
	// select from all channels
	wg.Add(len(channels))
	for _, ch := range channels {
		go multiplex(ch)
	}

	// wait for all reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()
	return multiplexedStream
}

func withFanout() {
	toInt := func(done <-chan any, valueStream <-chan any) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for val := range valueStream {
				select {
				case <-done:
					return
				case intStream <- val.(int):
				}
			}
		}()
		return intStream
	}
	take := func(done <-chan any, valueStream <-chan any, limit int) <-chan any {
		takeStream := make(chan any)
		go func() {
			defer close(takeStream)
			for i := 0; i < limit; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	repeatFunc := func(done <-chan any, fn func() any) <-chan any {
		valueStream := make(chan any)
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	random := func() any { return rand.Intn(50000000) }

	done := make(chan any)
	defer close(done)

	randIntStream := toInt(done, repeatFunc(done, random))

	numFinders := runtime.NumCPU()
	finders := make([]<-chan any, numFinders)

	for i := 0; i < numFinders; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}
	for prime := range take(done, FanIn(done, finders...), 10) {
		fmt.Printf("print some numbers: %d", prime)
	}
}

func withoutFanout() {
	toInt := func(done <-chan any, valueStream <-chan any) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for val := range valueStream {
				select {
				case <-done:
					return
				case intStream <- val.(int):
				}
			}
		}()
		return intStream
	}
	take := func(done <-chan any, valueStream <-chan any, limit int) <-chan any {
		takeStream := make(chan any)
		go func() {
			defer close(takeStream)
			for i := 0; i < limit; i++ {
				select {
				case <-done:
					return
				case takeStream <- <-valueStream:
				}
			}
		}()
		return takeStream
	}
	repeatFunc := func(done <-chan any, fn func() any) <-chan any {
		valueStream := make(chan any)
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}
	random := func() any { return rand.Intn(50000000) }

	done := make(chan any)
	defer close(done)

	randIntStream := toInt(done, repeatFunc(done, random))

	for prime := range take(done, primeFinder(done, randIntStream), 10) {
		fmt.Printf("print some numbers: %d", prime)
	}
}

func primeFinder(done <-chan any, intStream <-chan int) <-chan any {
	return nil
}
