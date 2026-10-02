package main

import (
	"encoding/json"
	"fmt"
)

type Webhook struct {
	ID       string `json:"id"`
	URL      string `json:"url"`
	Attempts int    `json:"attempts"`
	// omitempty: leave the key out when the value is the zero value.
	LastError string `json:"last_error,omitempty"`
	// Lowercase field: invisible to encoding/json, silently skipped.
	secret string
}

func main() {
	// 1. Struct -> JSON. Marshal returns bytes, so convert to string to print.
	w := Webhook{ID: "wh-1", URL: "http://localhost:9000/hook", secret: "hidden"}
	data, err := json.Marshal(w)
	if err != nil {
		fmt.Println("marshal failed:", err)
		return
	}
	fmt.Println("compact:", string(data))

	// 2. Pretty-printed version, handy for debugging.
	w.Attempts = 3
	w.LastError = "connection refused"
	pretty, _ := json.MarshalIndent(w, "", " ")
	fmt.Println(string(pretty))

	// 3. JSON -> struct. Pass a pointer so Unmarshal can fill it in.
	// Unknown keys are ignored; missing keys leave the zero value.
	input := []byte(`{"id":"wh-2","url":"http://example.com","extra":"ignored"}`)
	var got Webhook
	if err := json.Unmarshal(input, &got); err != nil {
		fmt.Println("unmarshal failed:", err)
		return
	}
	fmt.Printf("parsed: %+v\n", got)

	// 4. Broken JSON gives an error instead of a half-filled struct.
	var bad Webhook
	err = json.Unmarshal([]byte(`{"id": "wh-3", "attempts": "three"}`), &bad)
	fmt.Println("bad input error:", err)
}
