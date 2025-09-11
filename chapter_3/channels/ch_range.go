package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func main() {
	bufferedChannels()
}

func bufferedChannels() {
	var stdoutBuf bytes.Buffer
	defer stdoutBuf.WriteTo(os.Stdout)

	instream := make(chan int, 4)
	go func() {
		defer close(instream)
		defer fmt.Fprintln(&stdoutBuf, "Producer done.")

		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuf, "Sending: %d\n", i)
			instream <- i
		}
	}()
	for integer := range instream {
		fmt.Fprintf(&stdoutBuf, "Received: %d\n", integer)
	}
}

func unblockGoroutines() {
	begin := make(chan any)
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}()
	}
	fmt.Println("unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func rangeOverChannel() {
	instream := make(chan int)
	go func() {
		defer close(instream)
		for i := 1; i <= 5; i++ {
			instream <- i
		}
	}()
	for integer := range instream {
		fmt.Printf("%v ", integer)
	}
}
