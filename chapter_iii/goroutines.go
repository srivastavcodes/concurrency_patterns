package main

import (
	"fmt"
	"sync"
)

func goroutine1() {
	var wg sync.WaitGroup
	sayHello := func() {
		defer wg.Done()
		fmt.Println("hello world")
	}
	wg.Add(1)
	go sayHello()
	wg.Wait()
}

func goroutine2() {
	var wg sync.WaitGroup
	salutation := "Say Hello!"
	wg.Add(1)
	go func() {
		defer wg.Done()
		salutation = "welcome"
	}()
	wg.Wait()
	fmt.Println(salutation)
}

// nondeterministic results
func goroutine3() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"Hello", "Greetings", "Good day!"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println(salutation)
		}()
	}
	wg.Wait()
}

// correct way of writing above code
func goroutine4() {
	var wg sync.WaitGroup
	for _, salutation := range []string{"Hello", "Greetings", "Good day!"} {
		wg.Add(1)
		go func(salutation string) {
			defer wg.Done()
			fmt.Println(salutation)
		}(salutation) // copy is passed to the function so the access still remains
	}
	wg.Wait()
}

func main() {
	fmt.Println("before goroutine1")
	goroutine1()
	fmt.Println("after goroutine1")
	goroutine2()
	fmt.Println("after goroutine2")
	goroutine3()
	fmt.Println("after goroutine3")
	goroutine4()
	fmt.Println("after goroutine4")
}
