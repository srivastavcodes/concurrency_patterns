package main

import (
	"testing"
	"time"
)

func doWorkIntervals(done <-chan any, pulseInterval time.Duration, nums ...int) (<-chan any, <-chan int) {
	heartbeat := make(chan any)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(time.Second * 2)
		pulse := time.Tick(pulseInterval)

	outer:
		for _, num := range nums {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					select {
					case heartbeat <- struct{}{}:
					default:
					}
				case intStream <- num:
					continue outer
				}
			}
		}
	}()
	return heartbeat, intStream
}

func TestDoWork_HeartbeatIntervals(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	const timeout = 2 * time.Second
	heartbeat, results := doWorkIntervals(done, timeout/2, intSlice...)

	<-heartbeat // <4>

	i := 0
	for {
		select {
		case r, ok := <-results:
			if ok == false {
				return
			} else if expected := intSlice[i]; r != expected {
				t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
			}
			i++
		case <-heartbeat: // <5>
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
}
