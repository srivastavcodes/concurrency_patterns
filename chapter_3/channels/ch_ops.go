package main

import "fmt"

func main() {
	readOps1()
}

func readOps1() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i < 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}
	resultStream := chanOwner()
	for integer := range resultStream {
		fmt.Printf("received: %d\n", integer)
	}
	fmt.Println("done receiving")
}
