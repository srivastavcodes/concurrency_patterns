package main

import (
	"fmt"
	"sync"
)

func main() {
	onceExample()
}

func onceExample() {
	var count int

	increment := func() {
		count++
	}
	var once sync.Once

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			once.Do(increment)
		}()
	}
	wg.Wait()
	fmt.Printf("count is: %d\n", count)
}
