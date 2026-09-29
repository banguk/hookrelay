package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		fmt.Println("attempt", i)
	}

	backoff := 1
	for backoff < 10 {
		fmt.Println("backoff", backoff)
		backoff *= 2
	}

	// Infinite loop with break, and continue to skip.
	n := 0
	for {
		n++
		if n%2 == 0 {
			continue
		}
		if n > 5 {
			break
		}
		fmt.Println("odd", n)
	}

	// Range over an integer (Go 1.22+).
	for i := range 3 {
		fmt.Println("range int", i)
	}

	// Range over a slice: index and copy.
	statuses := []string{"pending", "delivering", "delivered"}
	for i, s := range statuses {
		fmt.Println(i, s)
	}

	// if with an init statement; v is scoped to the if/else.
	if v := len(statuses); v > 2 {
		fmt.Println("more than two statuses:", v)
	} else {
		fmt.Println("two or fewer:", v)
	}

	// switch on a value, no break needed.
	for _, s := range statuses {
		switch s {
		case "pending", "delivering":
			fmt.Println(s, "-> in progress")
		case "delivered":
			fmt.Println(s, "-> done")
		default:
			fmt.Println(s, "-> unknown")
		}
	}

	attempts := 3
	switch {
	case attempts == 0:
		fmt.Println("never tried")
	case attempts < 5:
		fmt.Println("retrying")
	default:
		fmt.Println("give up")
	}

}
