package main

import (
	"context"
	"log"
	"os"
	"sync"

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
	rateLimiter *rate.Limiter
}

func Open() *ApiConnection {
	return &ApiConnection{
		rateLimiter: rate.NewLimiter(rate.Limit(1), 1),
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
