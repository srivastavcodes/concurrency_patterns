package main

import (
	"fmt"
	"time"
)

func main() {

}

func readLeakFixed() {
	// done being the first param is convention
	doWork := func(done <-chan any, strings <-chan string) <-chan any {
		completed := make(chan any)
		go func() {
			defer fmt.Println("work completed")
			defer close(completed)
			for {
				select {
				case str := <-strings:
					// do something real
					fmt.Println(str)
				case <-done:
					return
				}
			}
		}()
		return completed
	}
	done := make(chan any)

	stringStream := make(chan string, 1)
	stringStream <- "stream deez nuts"

	completed := doWork(done, stringStream)

	go func() {
		time.Sleep(time.Second)
		fmt.Println("cancelling do work go routine")
		close(done)
	}()
	<-completed
	fmt.Println("done")
}

func readLeak() {
	doWork := func(strs <-chan string) <-chan string {
		completed := make(chan string)
		go func() {
			defer fmt.Println("work completed")
			defer close(completed)
			for s := range strs {
				// do something interesting
				fmt.Println(s)
			}
		}()
		return completed
	}
	doWork(nil)
	fmt.Println("done")
}
