package main

import (
	"bytes"
	"fmt"
	"sync"
)

func main() {
	lexicalConfinementBuffer()
}

func lexicalConfinementBuffer() {
	printData := func(wg *sync.WaitGroup, data []byte) {
		defer wg.Done()

		var buf bytes.Buffer
		for _, b := range data {
			fmt.Fprintf(&buf, "%c", b)
		}
		fmt.Println(buf.String())
	}
	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("stickynotes")
	go printData(&wg, data[:6])
	go printData(&wg, data[6:])
	wg.Wait()
	// data remains as it, copy is passed around
	fmt.Println(string(data))
}

func lexicalConfinement() {
	chanOwner := func() <-chan int {
		results := make(chan int)
		go func() {
			defer close(results)
			for i := 0; i < 5; i++ {
				results <- i
			}
		}()
		return results
	}
	consumer := func(results <-chan int) {
		for data := range results {
			fmt.Printf("received: %d\n", data)
		}
		fmt.Println("done receiving")
	}
	results := chanOwner()
	consumer(results)
}

func adhocConfinement() {
	data := make([]int, 4)

	loopData := func(handleData chan<- int) {
		defer close(handleData)
		for i := range data {
			handleData <- data[i]
		}
	}
	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}
