package main

import "fmt"

func main() {
	repeatGenerator()
}

// repeats any value given until done is closed
func repeatGenerator() {
	// any should be replaced with the relevant data type
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
	fmt.Println(repeat(done, 1, 2, 3, 4, 5))
}
