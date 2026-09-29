package main

import "fmt"

type Webhook struct {
	ID       string
	URL      string
	Attempts int
	Status   string
}

func (w Webhook) Summary() string {
	return fmt.Sprintf("[%s] %s attempts=%d status=%s", w.ID, w.URL, w.Attempts, w.Status)
}

func (w *Webhook) RecordAttempt() {
	w.Attempts++
	w.Status = "delivering"
}

func main() {
	first := Webhook{ID: "wh-1", URL: "http://localhost:9000/hook", Status: "pending"}

	var empty Webhook
	fmt.Println("empty:", empty.Summary())

	var queue []Webhook
	queue = append(queue, first)
	queue = append(queue, Webhook{ID: "wh-2", URL: "http://localhost:9000/hook", Status: "pending"})
	queue = append(queue, Webhook{ID: "wh-3", URL: "http://localhost:9001/hook", Status: "pending"})
	fmt.Println("queue length:", len(queue))

	for i, w := range queue {
		fmt.Println(i, w.Summary())
	}

	queue[0].RecordAttempt()
	queue[0].RecordAttempt()
	fmt.Println("after attempts:", queue[0].Summary())

	byID := make(map[string]Webhook)
	for _, w := range queue {
		byID[w.ID] = w
	}

	if w, ok := byID["wh-2"]; ok {
		fmt.Println("found:", w.Summary())
	}
	if _, ok := byID["wh-99"]; !ok {
		fmt.Println("wh-99 not found")
	}

	missing := byID["wh-99"]
	fmt.Println("missing id is empty string:", missing.ID == "")

	delete(byID, "wh-3")
	fmt.Println("map size after delete:", len(byID))
}
