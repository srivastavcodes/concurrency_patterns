package main

import (
	"fmt"
	"net/http"
)

func main() {
	goodHandlingMadeBetter()
}

type Result struct {
	Error    error
	Response *http.Response
}

func goodHandlingMadeBetter() {
	checkStatus := func(done <-chan any, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				res, err := http.Get(url)

				result := Result{
					Response: res,
					Error:    err,
				}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}
	done := make(chan any)
	defer close(done)

	errCount := 0
	urls := []string{"a", "https://google.com", "b", "c", "d", "e"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!!!")
				break
			}
			continue
		}
		fmt.Printf("response: %v\n", result.Response.Status)
	}
}

func goodHandling() {
	checkStatus := func(done <-chan any, urls ...string) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				res, err := http.Get(url)

				result := Result{
					Response: res,
					Error:    err,
				}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()
		return results
	}
	done := make(chan any)
	defer close(done)

	urls := []string{"https://google.com", "https://badhost"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			continue
		}
		fmt.Printf("response: %v\n", result.Response.Status)
	}
}
