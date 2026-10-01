package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

func deliver(ctx context.Context, id string, delay time.Duration) error {
	select {
	case <-time.After(delay):
		fmt.Println(id, "delivered after", delay)
		return nil
	case <-ctx.Done():
		fmt.Println(id, "gave up:", ctx.Err())
		return ctx.Err()
	}
}

func main() {
	// 1. Timeout: the deadline passes before the work finishes.
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	err := deliver(ctx, "wh-1", 500*time.Millisecond)
	fmt.Println("wh-1 error:", err)
	fmt.Println("is DeadlineExceeded?", errors.Is(err, context.DeadlineExceeded))

	// 2. Same context, but the work is fast enough.
	ctx2, cancel2 := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel2()

	err = deliver(ctx2, "wh-2", 50*time.Millisecond)
	fmt.Println("wh-2 error:", err)

	// 3. Manual cancel: one call stops every goroutine sharing this ctx.
	ctx3, cancel3 := context.WithCancel(context.Background())

	for _, id := range []string{"wh-3", "wh-4", "wh-5"} {
		go deliver(ctx3, id, 1*time.Second)
	}

	time.Sleep(100 * time.Millisecond)
	fmt.Println("operator pressed cancel")
	cancel3()

	time.Sleep(100 * time.Millisecond)

}
