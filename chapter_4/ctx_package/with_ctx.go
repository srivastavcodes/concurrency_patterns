package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printGreetings2(ctx); err != nil {
			fmt.Printf("cannot print greetings: %v\n", err)
			cancel()
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := printFarewell2(ctx); err != nil {
			fmt.Printf("cannot print farewell: %v\n", err)
		}
	}()
	wg.Wait()
}

func printGreetings2(ctx context.Context) error {
	greeting, err := genGreeting2(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func printFarewell2(ctx context.Context) error {
	greeting, err := genFarewell2(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world!\n", greeting)
	return nil
}

func genGreeting2(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale2, err := betterLocale2(ctx); {
	case err != nil:
		return "", err
	case locale2 == "EN/INDIA":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupported locale2")
}

func genFarewell2(ctx context.Context) (string, error) {
	switch locale2, err := betterLocale2(ctx); {
	case err != nil:
		return "", err
	case locale2 == "EN/INDIA":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale2")
}

func locale2(ctx context.Context) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Second):
	}
	return "EN/INDIA", nil
}

func betterLocale2(ctx context.Context) (string, error) {
	// only useful if you have some idea of how long your call-graph will take
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(10*time.Second)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}
	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case <-time.After(10 * time.Second):
	}
	return "EN/INDIA", nil
}
