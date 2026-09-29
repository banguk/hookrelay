package main

import "fmt"

type Webhook struct {
	ID       string
	Attempts int
}

// Receives a copy. Changes stay inside this function.
func bumpCopy(n int) {
	n++
}

// Receives an address. Changes reach the caller's variable.
func bumpPointer(n *int) {
	*n++
}

// Same idea with a struct.
func retryCopy(w Webhook) {
	w.Attempts++
}

func retryPointer(w *Webhook) {
	// Go lets you write w.Attempts instead of (*w).Attempts.
	w.Attempts++
}

func main() {
	// 1. & gives an address, * follows it.
	count := 1
	p := &count
	fmt.Println("count:", count, "p points to:", *p)
	*p = 10
	fmt.Println("after *p = 10, count:", count)

	// 2. Passing a value vs passing an address.
	bumpCopy(count)
	fmt.Println("after bumpCopy:", count)
	bumpPointer(&count)
	fmt.Println("after bumpPointer:", count)

	// 3. The same with a struct.
	w := Webhook{ID: "wh-1"}
	retryCopy(w)
	fmt.Println("after retryCopy:", w.Attempts)
	retryPointer(&w)
	fmt.Println("after retryPointer:", w.Attempts)

	queue := []*Webhook{
		{ID: "wh-1"},
		{ID: "wh-2"},
	}

	for _, item := range queue {
		item.Attempts++
	}
	fmt.Println("queue attempts:", queue[0].Attempts, queue[1].Attempts)

	var missing *Webhook
	fmt.Println("missing is nil", missing == nil)
}
