package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Order struct {
	ID     int
	Status string
	mu     sync.Mutex
}

var (
	updateCount int
	updateMutex sync.Mutex
)

func main() {
	var wg sync.WaitGroup

	orderChan := make(chan *Order, 20)
	processChan := make(chan *Order, 20)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, order := range generateOrders(20) {
			orderChan <- order
		}
		close(orderChan)
		fmt.Println("orders generated")
	}()

	wg.Add(1)
	go processOrders(orderChan, processChan, &wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case order, ok := <-processChan:
				if !ok {
					fmt.Println("Channel closed")
					return
				}
				fmt.Printf("Processed order: %d with status: %s\n", order.ID, order.Status)

			case <-time.After(10 * time.Second):
				fmt.Println("Timeout waiting for order")
				return
			}
		}
	}()

	/*
		for i := 0; i < 3; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for order := range orderChan {
					updateOrderStatus(order)
				}
			}()
		}
	*/

	wg.Wait()

	fmt.Println(updateCount)
	fmt.Println("All operations completed.")
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count)
	for i := 0; i < count; i++ {
		orders[i] = &Order{ID: i + 1, Status: "Pending"}
	}
	return orders
}

func updateOrderStatus(order *Order) {
	order.mu.Lock()

	time.Sleep(time.Duration(rand.Intn(300)) * time.Millisecond)
	status := []string{
		"Processing", "Shipped", "Delivered",
	}[rand.Intn(3)]
	order.Status = status
	fmt.Printf("Order %d: updated to: %s\n", order.ID, order.Status)

	order.mu.Unlock()

	updateMutex.Lock()
	defer updateMutex.Unlock()

	currUpdates := updateCount
	time.Sleep(5 * time.Millisecond)
	updateCount = currUpdates + 1
}

func processOrders(inchan <-chan *Order, outchan chan<- *Order, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		close(outchan)
	}()
	for order := range inchan {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		order.Status = "Processed"
		outchan <- order
	}
}

func reportOrderStatus(orders []*Order) {
	for i := 0; i < 5; i++ {
		time.Sleep(1 * time.Second)
		fmt.Println("--------Order Status Report--------")

		for _, order := range orders {
			fmt.Printf("Order %d: %v\n", order.ID, order.Status)
		}
	}
	fmt.Println("--------Orders Reported--------")
}
