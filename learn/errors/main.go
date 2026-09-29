package main

import (
	"errors"
	"fmt"
)

type Webhook struct {
	ID  string
	URL string
}

// Store is the behaviour we need from any storage backend.
// No type declares that it implements this; matching methods is enou
type Store interface {
	Save(W Webhook) error
	Get(id string) (Webhook, error)
}

// Sentinel error: one shared value callers can compare against.
var ErrNotFound = errors.New("webhook not found")

// Custom error type: carries data, so it is a struct with an Error method.
type ValidationError struct {
	Field string
}

func (e ValidationError) Error() string {
	return "invalid " + e.Field
}

type MemoryStore struct {
	items map[string]Webhook
}

// Constructor: Go has no constructors, so a NewX function is the convention.
// & takes the address; pointers get their own session later.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{items: make(map[string]Webhook)}
}

func (s *MemoryStore) Save(w Webhook) error {
	if w.ID == "" {
		return ValidationError{Field: "ID"}
	}
	if w.URL == "" {
		return ValidationError{Field: "URL"}
	}
	s.items[w.ID] = w
	return nil
}

func (s *MemoryStore) Get(id string) (Webhook, error) {
	w, ok := s.items[id]

	if !ok {
		return Webhook{}, fmt.Errorf("get %q: %w", id, ErrNotFound)
	}

	return w, nil
}

// describe only knows about Store. It never mentions MemoryStore.
func describe(s Store, id string) {
	w, err := s.Get(id)
	if err != nil {
		fmt.Println("describe failed:", err)
		return
	}
	fmt.Println("describe", w.ID, "->", w.URL)
}

func main() {
	store := NewMemoryStore()

	// The most common shape in Go: call, check err, move on.
	if err := store.Save(Webhook{ID: "wh-1", URL: "http://localhost:9000/hook"}); err != nil {
		fmt.Println("unexpected:", err)
	}

	err := store.Save(Webhook{ID: "", URL: "http://localhost:9000/hook"})
	fmt.Println("save with empty id:", err)

	var ve ValidationError
	if errors.As(err, &ve) {
		fmt.Println("validation failed on field:", ve.Field)
	}

	describe(store, "wh-1")
	describe(store, "wh-9")

	_, err = store.Get("wh-9")
	fmt.Println("message:", err)
	fmt.Println("errors.Is ErrNotFound?", errors.Is(err, ErrNotFound))
	fmt.Println("== ErrNotFound?", err == ErrNotFound)
}
