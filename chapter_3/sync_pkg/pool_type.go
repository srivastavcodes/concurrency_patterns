package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func main() {
	poolPattern()
}

func poolExample() {
	pool := &sync.Pool{
		New: func() any {
			fmt.Println("creating new one")
			return struct{}{}
		},
	}
	one := pool.Get()
	two := pool.Get()

	pool.Put(one)
	pool.Put(two)

	pool.Get()
	pool.Get()
}

func poolPattern() {
	var creationCount int

	pool := &sync.Pool{
		New: func() any {
			creationCount++
			buf := make([]byte, 1024)
			return &buf
		},
	}
	// created a pool with 4kb
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())
	pool.Put(pool.New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			buf := pool.Get().(*[]byte)
			defer pool.Put(buf)
			// assume something happens here
		}()
	}
	wg.Wait()
	fmt.Printf("%d calcualtors were created", creationCount)
}

func connectToService() any {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

// --------------- with caching(creating pools, etc...)

func warmServiceConnCache() *sync.Pool {
	pool := &sync.Pool{
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		pool.Put(pool.New())
	}
	return pool
}

func startNetworkDaemon2() *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", "localhost:4000")
		if err != nil {
			log.Fatalf("khatam, samapt, done: %v", err)
		}
		defer server.Close()
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v", err)
				continue
			}
			srvconn := connPool.Get()

			fmt.Fprintln(conn, "")
			connPool.Put(srvconn)
			conn.Close()
		}
	}()
	return &wg
}

// --------------- without caching(creating pools, etc...)

func startNetworkDaemon1() *sync.WaitGroup {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", "localhost:4000")
		if err != nil {
			log.Fatalf("khatam, samapt, done: %v", err)
		}
		defer server.Close()
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("failed to accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Println(conn, "")
			_ = conn.Close()
		}
	}()
	return &wg
}
