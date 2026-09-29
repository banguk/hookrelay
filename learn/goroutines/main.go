package main

import (
	"fmt"
	"time"
)

type Result struct {
	ID       string
	Status   string
	Duration time.Duration
}

func deliver(id string, delay time.Duration, results chan Result) {
	time.Sleep(delay)
	results <- Result{ID: id, Status: "delivered", Duration: delay}
}

func main() {
	start := time.Now()

	done := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		done <- "background work finished"
		time.Sleep(200 * time.Millisecond) // main은 이걸 기다리지 않음
		fmt.Println("goroutine really done")
	}()

	msg := <-done
	fmt.Println(msg)

	results := make(chan Result)

	go deliver("wh-1", 300*time.Millisecond, results)
	go deliver("wh-2", 100*time.Millisecond, results)
	go deliver("wh-3", 200*time.Millisecond, results)

	for range 3 {
		r := <-results
		fmt.Printf("%s %s in %v\n", r.ID, r.Status, r.Duration)
	}

	fmt.Println("total elapsed:", time.Since(start).Round(10*time.Millisecond))
}
