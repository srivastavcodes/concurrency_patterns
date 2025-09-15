package main

import (
	"fmt"
	"math/rand"
)

func main() {
	repeatFunction()
}

func repeatFunction() {
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
	done := make(chan any)
	defer close(done)

	random := func() any { return rand.Int() }

	for val := range take(done, repeatFunc(done, random), 10) {
		fmt.Println(val)
	}
}
