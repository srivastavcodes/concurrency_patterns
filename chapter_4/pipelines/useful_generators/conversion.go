package main

import "fmt"

func main() {
	convertGeneratorToString()
}

func convertGeneratorToString() {
	toString := func(done <-chan any, valueStream <-chan any) <-chan string {
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)

			for val := range valueStream {
				select {
				case <-done:
					return
				case stringStream <- val.(string):
				}
			}
		}()
		return stringStream
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

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
}
