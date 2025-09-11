package main

import (
	"fmt"
	"runtime"
	"sync"
)

func main() {
	memoryStats()
}

func singleVariable() {
	var wg sync.WaitGroup
	salutation := "hello"

	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation) // prints welcome
}

// rangeFunc1 wrong usage
func rangeFunc1() {
	var wg sync.WaitGroup

	for _, val := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(val)
		}()
	}
	wg.Wait()
}

// rangeFunc2 correct usage
func rangeFunc2() {
	var wg sync.WaitGroup

	for _, val := range []string{"hello", "greetings", "good day"} {
		wg.Add(1)
		go func(val string) {
			defer wg.Done()
			fmt.Println(val)
		}(val)
	}
	wg.Wait()
}

func memoryStats() {
	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}
	var ch <-chan any
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-ch }

	const numGoroutines = 1e4
	wg.Add(numGoroutines)

	before := memConsumed()
	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb", float64(after-before)/numGoroutines/1000)
}
