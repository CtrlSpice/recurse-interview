package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewServer(t *testing.T) {
	server := NewServer("localhost:4000")
	if server == nil {
		t.Fatal("NewServer() returned nil")
	}
	if server.store == nil {
		t.Fatal("Server.store is nil")
	}
	if server.server.Addr != "localhost:4000" {
		t.Fatalf("Expected address 'localhost:4000', got '%s'", server.server.Addr)
	}
}

func TestServer_SetHandler_Success(t *testing.T) {
	server := NewServer("localhost:4000")

	req, err := http.NewRequest("GET", "/set?name=Obi-Wan Kenobi", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.setHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	expected := "Set key 'name' to value 'Obi-Wan Kenobi'"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}

	// Verify the value was actually stored for real this time
	value := server.store.Get("name")
	if value != "Obi-Wan Kenobi" {
		t.Fatalf("Expected 'Obi-Wan Kenobi' in store, got '%v'", value)
	}
}

func TestServer_SetHandler_MissingKey(t *testing.T) {
	server := NewServer("localhost:4000")

	req, err := http.NewRequest("GET", "/set", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.setHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	expected := "missing key parameter\n"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestServer_GetHandler_Success(t *testing.T) {
	server := NewServer("localhost:4000")

	server.store.Set("name", "Obi-Wan Kenobi")

	req, err := http.NewRequest("GET", "/get?key=name", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	expected := "Obi-Wan Kenobi"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestServer_GetHandler_MissingKey(t *testing.T) {
	server := NewServer("localhost:4000")

	req, err := http.NewRequest("GET", "/get", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Fatalf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	expected := "missing key parameter\n"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestServer_GetHandler_KeyNotFound(t *testing.T) {
	server := NewServer("localhost:4000")

	req, err := http.NewRequest("GET", "/get?key=nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, status)
	}

	expected := "key not found\n"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestServer_Handler_Routes(t *testing.T) {
	server := NewServer("localhost:4000")
	handler := server.Handler()

	// set route
	req, err := http.NewRequest("GET", "/set?test=value", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	// get route
	req, err = http.NewRequest("GET", "/get?key=test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Expected status %d, got %d", http.StatusOK, status)
	}

	expected := "value"
	if rr.Body.String() != expected {
		t.Fatalf("Expected '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestServer_Handler_NotFound(t *testing.T) {
	server := NewServer("localhost:4000")
	handler := server.Handler()

	req, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Fatalf("Expected status %d, got %d", http.StatusNotFound, status)
	}
}
