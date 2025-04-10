package main

import (
	"fmt"
	"time"
)

func main() {
	race1()
}

func race1() {
	var data int
	go func() { data++ }()
	time.Sleep(1 * time.Second) // bad practice
	if data == 0 {
		fmt.Printf("the value is %v.\n", data)
	}
}
