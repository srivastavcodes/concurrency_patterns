package main

import "fmt"

func main() {

}

func orDoneChannel() {
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
	done := make(chan any)
	defer close(done)

	ch := make(chan any)
	defer close(ch)

	// the or-done pattern allows for simple for-loops like this
	for val := range orDone(done, ch) {
		// do something with the channel
		fmt.Printf("%v", val)
	}
}
