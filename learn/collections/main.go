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
	// Struct literal with field names. Fields you skip get zero values.
	first := Webhook{ID: "wh-1", URL: "http://localhost:9000/hook", Status: "pending"}

	// Zero-value struct: every field is its zero value.
	var empty Webhook
	fmt.Println("empty:", empty.Summary())

	// Slice: starts nil, append grows it.
	var queue []Webhook
	queue = append(queue, first)
	queue = append(queue, Webhook{ID: "wh-2", URL: "http://localhost:9000/hook", Status: "pending"})
	queue = append(queue, Webhook{ID: "wh-3", URL: "http://localhost:9001/hook", Status: "pending"})
	fmt.Println("queue length:", len(queue))

	// range gives index and a COPY of the element.
	for i, w := range queue {
		fmt.Println(i, w.Summary())
	}

	// Mutate through the slice index, not the range copy.
	queue[0].RecordAttempt()
	queue[0].RecordAttempt()

	for _, w := range queue {
		w.RecordAttempt()
	}

	fmt.Println("after attempts:", queue[0].Summary())

	byID := make(map[string]Webhook)
	for _, w := range queue {
		byID[w.ID] = w
	}

	// Lookup with the comma-ok form.
	if w, ok := byID["wh-2"]; ok {
		fmt.Println("found:", w.Summary())
	}
	if _, ok := byID["wh-99"]; !ok {
		fmt.Println("wh-99 not found")
	}

	// Missing key returns the zero value, not an error.
	missing := byID["wh-99"]
	fmt.Println("missing id is empty string:", missing.ID == "")

	delete(byID, "wh-3")
	fmt.Println("map size after delete:", len(byID))
}
