package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	usingWard()
}

func usingWard() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan any)
	defer close(done)

	doWork, intStream := doWorkFunc(done, 1, 2, -1, 3, 4, 5)
	doWorkWithSteward := newSteward(1*time.Millisecond, doWork)
	doWorkWithSteward(done, 1*time.Minute)

	for intVal := range take(done, intStream, 6) {
		fmt.Printf("Received: %v\n", intVal)
	}
}

func newSteward(timeout time.Duration, startGoroutine startGoroutineFunc) startGoroutineFunc {
	return func(done <-chan any, pulseInterval time.Duration) <-chan any {
		heartbeat := make(chan any)
		go func() {
			defer close(heartbeat)

			var wardDone chan any
			var wardHeartbeat <-chan any

			startWard := func() {
				wardDone = make(chan any)
				wardHeartbeat = startGoroutine(or(wardDone, done), timeout/2)
			}
			startWard()
			pulse := time.Tick(pulseInterval)
		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						log.Println("steward: ward unhealthy; restarting")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

func doWorkFunc(done <-chan any, intList ...int) (startGoroutineFunc, <-chan any) {
	intchanStream := make(chan (<-chan any))
	intStream := bridge(done, intchanStream)

	doWork := func(done <-chan any, pulseInterval time.Duration) <-chan any {
		heartbeat := make(chan any)
		intStream := make(chan any)
		go func() {
			defer close(intStream)
			select {
			case intchanStream <- intStream:
			case <-done:
				return
			}
			pulse := time.Tick(pulseInterval)
			for {
			valueLoop:
				for _, intVal := range intList {
					if intVal < 0 {
						log.Printf("negative value: %v\n", intVal)
						return
					}
					for {
						select {
						case <-pulse:
							select {
							case heartbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return heartbeat
	}
	return doWork, intStream
}

func bridge(done <-chan any, chanStream <-chan <-chan any) <-chan any {
	valStream := make(chan any)
	go func() {
		defer close(valStream)
		for {
			var stream <-chan any
			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}
			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func orDone(done <-chan any, ch <-chan any) <-chan any {
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

func take(done <-chan any, valueStream <-chan any, num int) <-chan any {
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
