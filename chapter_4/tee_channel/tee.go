package main

import "fmt"

func main() {
	splitWithTee()
}

func splitWithTee() {
	orDone := func(done <-chan any, ch <-chan any) <-chan any {
		valStream := make(chan any)
		go func() {
			defer close(valStream)
			for {
				select {
				case <-done:
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					select {
					case valStream <- v:
					case <-done:
					}
				}
			}
		}()
		return valStream
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
	tee := func(done <-chan any, instream <-chan any) (_, _ <-chan any) {
		out1 := make(chan any)
		out2 := make(chan any)
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done, instream) {
				var out1, out2 = out1, out2
				for i := 0; i < 2; i++ {
					select {
					case <-done:
						return
					case out1 <- val:
						out1 = nil
					case out2 <- val:
						out2 = nil
					}
				}
			}
		}()
		return out1, out2
	}
	done := make(chan any)
	defer close(done)

	out1, out2 := tee(done, take(done, repeat(done, 1, 2), 4))
	for val := range out1 {
		fmt.Printf("out1: %v, out2: %v\n", val, <-out2)
	}
}
