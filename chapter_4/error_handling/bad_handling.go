package main

import (
	"fmt"
	"net/http"
)

func main() {
	badHandling()
}

func badHandling() {
	checkStatus := func(done <-chan any, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				res, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- res:
				}
			}
		}()
		return responses
	}
	done := make(chan any)
	defer close(done)

	urls := []string{"https://google.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("response: %v\n", response.Status)
	}
}
