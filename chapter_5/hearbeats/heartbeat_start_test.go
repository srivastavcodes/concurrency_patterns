package main

import (
	"testing"
	"time"
)

func doWorkAtStart(done <-chan any, nums ...int) (<-chan any, <-chan int) {
	heartbeat := make(chan any)
	intStream := make(chan int)
	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2 * time.Second)

		for _, val := range nums {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
			select {
			case <-done:
				return
			case intStream <- val:
			}
		}
	}()
	return heartbeat, intStream
}

func TestDoWork_GoodTestExample(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	heartbeat, results := doWorkAtStart(done, intSlice...)

	<-heartbeat // <1>

	i := 0
	for r := range results {
		if expected := intSlice[i]; r != expected {
			t.Errorf("index %v: expected %v, but received %v,", i, expected, r)
		}
		i++
	}
}

func TestDoWork_BadTestExample(t *testing.T) {
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0, 1, 2, 3, 5}
	_, results := doWorkAtStart(done, intSlice...)

	for i, expected := range intSlice {
		select {
		case r := <-results:
			if r != expected {
				t.Errorf(
					"index %v: expected %v, but received %v,",
					i,
					expected,
					r,
				)
			}
		case <-time.After(1 * time.Second): // <1>
			t.Fatal("test timed out")
		}
	}
}
