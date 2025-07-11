package main

import (
	"fmt"
	"testing"
)

func TestNewStore(t *testing.T) {
	store := NewStore()
	if store == nil {
		t.Fatal("NewStore() returned nil")
	}
	if store.Entries == nil {
		t.Fatal("Store.Entries is nil")
	}
	if len(store.Entries) != 0 {
		t.Fatalf("Expected empty store, got %d entries", len(store.Entries))
	}
}

func TestStore_SetAndGet_String(t *testing.T) {
	store := NewStore()

	// Test string key and value
	store.Set("name", "Obi-Wan Kenobi")
	value := store.Get("name")

	if value != "Obi-Wan Kenobi" {
		t.Fatalf("Expected 'Obi-Wan Kenobi', got '%v'", value)
	}
}

func TestStore_SetAndGet_Int(t *testing.T) {
	store := NewStore()

	// Test int key and value
	store.Set(123, 456)
	value := store.Get(123)

	if value != "456" {
		t.Fatalf("Expected '456', got '%v'", value)
	}
}

func TestStore_SetAndGet_Bool(t *testing.T) {
	store := NewStore()

	// Test bool key and value
	store.Set(true, false)
	value := store.Get(true)

	if value != "false" {
		t.Fatalf("Expected 'false', got '%v'", value)
	}
}

func TestStore_Get_NonExistent(t *testing.T) {
	store := NewStore()

	value := store.Get("Alderaan")
	if value != nil {
		t.Fatalf("Expected nil for non-existent key, got '%v'", value)
	}
}

func TestStore_Overwrite(t *testing.T) {
	store := NewStore()

	// Set initial value
	store.Set("key", "old_value")

	// Overwrite with new value
	store.Set("key", "new_value")

	value := store.Get("key")
	if value != "new_value" {
		t.Fatalf("Expected 'new_value', got '%v'", value)
	}
}

func TestStore_ConcurrentAccess(t *testing.T) {
	store := NewStore()
	done := make(chan bool, 10)

	// Start multiple goroutines setting values
	for i := range 10 {
		go func(id int) {
			key := fmt.Sprintf("key_%d", id)
			value := fmt.Sprintf("value_%d", id)
			store.Set(key, value)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for range 10 {
		<-done
	}

	// Verify all values were set correctly
	for i := range 10 {
		key := fmt.Sprintf("key_%d", i)
		expected := fmt.Sprintf("value_%d", i)
		value := store.Get(key)
		if value != expected {
			t.Fatalf("Expected '%s', got '%v' for key '%s'", expected, value, key)
		}
	}
}
