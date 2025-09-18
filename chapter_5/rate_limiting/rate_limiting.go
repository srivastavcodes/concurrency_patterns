package main

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	defer log.Print("Done")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	apiConn := Open()
	var wg sync.WaitGroup

	wg.Add(20)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConn.ReadFile(context.Background())
			if err != nil {
				log.Printf("cannot ReadFile: %v", err)
			}
			log.Printf("ReadFile")
		}()
	}
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			err := apiConn.ResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot ResolveAddress: %v", err)
			}
			log.Printf("ResolveAddress")
		}()
	}
	wg.Wait()
}

type ApiConnection struct {
	rateLimiter RateLimiter
}

func Open() *ApiConnection {
	secondLimit := rate.NewLimiter(Per(2, time.Second), 1)
	minuteLimit := rate.NewLimiter(Per(10, time.Minute), 10)
	return &ApiConnection{
		rateLimiter: NewMultiLimiter(secondLimit, minuteLimit),
	}
}

func (conn *ApiConnection) ReadFile(ctx context.Context) error {
	if err := conn.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// pretend we do work here
	return nil
}

func (conn *ApiConnection) ResolveAddress(ctx context.Context) error {
	if err := conn.rateLimiter.Wait(ctx); err != nil {
		return err
	}
	// pretend we do work here
	return nil
}
