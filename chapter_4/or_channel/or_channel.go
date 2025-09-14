package main

import (
	"fmt"
	"time"
)

func main() {
	example()
}

func example() {
	var or func(channels ...<-chan any) <-chan any

	or = func(channels ...<-chan any) <-chan any {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan any)
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	sig := func(after time.Duration) <-chan any {
		ch := make(chan any)
		go func() {
			defer close(ch)
			time.Sleep(after)
		}()
		return ch
	}
	start := time.Now()
	<-or(
		sig(3*time.Second),
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(2*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}
