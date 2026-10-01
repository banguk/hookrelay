package main

import (
	"fmt"
	"sync"
	"time"
)

// work does something without returning a result.
// wg is a pointer: copying a WaitGroup would break the count.
func work(id string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done() // runs when work returns, however it returns
	time.Sleep(delay)
	fmt.Println(id, "finished after", delay)
}

// slowDeliver sends its result only after delay.
func slowDeliver(delay time.Duration, results chan string) {
	time.Sleep(delay)
	results <- "delivered"
}

func main() {
	// 1. WaitGroup: wait until every goroutine has finished.
	var wg sync.WaitGroup

	wg.Add(3)
	go work("wh-1", 300*time.Millisecond, &wg)
	go work("wh-2", 100*time.Millisecond, &wg)
	go work("wh-3", 200*time.Millisecond, &wg)

	wg.Wait() // blocks until Done has been called three times
	fmt.Println("all workers finished")

	// 2. select with a timeout: the response arrives in time.
	fast := make(chan string)
	go slowDeliver(100*time.Millisecond, fast)

	select {
	case r := <-fast:
		fmt.Println("fast:", r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("fast: time out")
	}

	// 3. select with a timeout: the response is too slow.
	slow := make(chan string, 1) // buffer of 1, explained below
	go slowDeliver(2*time.Second, slow)

	select {
	case r := <-slow:
		fmt.Println("slow:", r)
	case <-time.After(500 * time.Millisecond):
		fmt.Println("slow: timed out")
	}
}
