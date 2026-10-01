package main

import (
	"fmt"
	"time"
)

func main() {
	// Experiment 1: a receive waits until someone sends.
	ch := make(chan int)
	go func() {
		time.Sleep(5 * time.Second)
		fmt.Println("  goroutine: sending 42")
		ch <- 42
	}()

	fmt.Println("main: waiting for a value...")
	v := <-ch
	fmt.Println("main: got", v)
	fmt.Println()

	// Experiment 2: closing wakes the receiver, with no value.
	ch2 := make(chan int)
	go func() {
		time.Sleep(1 * time.Second)
		fmt.Println("  goroutine: closing ch2")
		close(ch2)
	}()
	fmt.Println("main: waiting on ch2...")
	v2, ok := <-ch2
	fmt.Println("main: woke up. value:", v2, "ok:", ok)
	fmt.Println()

	ch3 := make(chan int)
	for i := 1; i <= 3; i++ {
		go func() {
			<-ch3
			fmt.Println("  waiter", i, "woke up")
		}()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("main: closing ch3")
	close(ch3)
	time.Sleep(100 * time.Millisecond)
}
