package main

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func pool1() {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("creating new instance.")
			return struct{}{}
		},
	}
	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func connectToService() interface{} {
	time.Sleep(1 * time.Second)
	return struct{}{}
}

// without using pool
func startNetworkDaemon1() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		server, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer func(server net.Listener) {
			err := server.Close()
			if err != nil {
				log.Fatalf("cannot close listener: %v", err)
			}
		}(server)
		wg.Done()
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			connectToService()
			fmt.Println(conn, "")
			_ = conn.Close()
		}
	}()
	return &wg
}

func warmServiceConnCache() *sync.Pool {
	p := &sync.Pool{
		New: connectToService,
	}
	for i := 0; i < 10; i++ {
		p.Put(p.New)
	}
	return p
}

// using pool
func startNetworkDaemon2() *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		connPool := warmServiceConnCache()

		server, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("cannot listen: %v", err)
		}
		defer func(server net.Listener) {
			_ = server.Close()
		}(server)
		wg.Done()

		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("cannot accept connection: %v", err)
				continue
			}
			svcConn := connPool.Get()
			fmt.Println(conn, "")
			connPool.Put(svcConn)
			_ = conn.Close()
		}
	}()
	return &wg
}
