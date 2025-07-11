package main

import (
	"fmt"
	"sync"
)

// Store is a simple key-value store that uses a map to store entries.
// Note: The key and value are both converted to strings in order to keep
// the implementation of the write-ahead log as simple as possible.
type Store struct {
	Entries map[string]string
	mu      sync.Mutex
}

func NewStore() *Store {
	return &Store{
		Entries: make(map[string]string),
		mu:      sync.Mutex{},
	}
}

func (s *Store) Set(key, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := fmt.Sprintf("%v", key)
	valueStr := fmt.Sprintf("%v", value)
	s.Entries[keyStr] = valueStr
}

func (s *Store) Get(key any) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	keyStr := fmt.Sprintf("%v", key)
	value, exists := s.Entries[keyStr]
	if !exists {
		return nil
	}
	return value
}
