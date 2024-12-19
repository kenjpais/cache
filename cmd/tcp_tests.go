package main

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func startServer() {
	go func() {
		main() // Call the server's main function
	}()
	time.Sleep(1 * time.Second) // Give the server time to start
}

func TestTCPServer(t *testing.T) {
	// Start the server in a separate goroutine
	startServer()

	// Connect to the server
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Helper function to send a command and get the response
	sendCommand := func(command string) string {
		_, err := conn.Write([]byte(command + "\n"))
		if err != nil {
			t.Fatalf("Failed to send command: %v", err)
		}

		response, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("Failed to read response: %v", err)
		}
		return strings.TrimSpace(response)
	}

	// Test SET command
	t.Run("SET Command", func(t *testing.T) {
		response := sendCommand("SET key1 value1")
		if response != "SET key1 OK" {
			t.Errorf("Expected 'SET key1 OK', got '%s'", response)
		}
	})

	// Test GET command
	t.Run("GET Command", func(t *testing.T) {
		response := sendCommand("GET key1")
		if response != "value1" {
			t.Errorf("Expected 'value1', got '%s'", response)
		}
	})

	// Test GET for non-existent key
	t.Run("GET Non-Existent Key", func(t *testing.T) {
		response := sendCommand("GET key_non_existent")
		if response != "Key not found" {
			t.Errorf("Expected 'Key not found', got '%s'", response)
		}
	})

	// Test DELETE command
	t.Run("DELETE Command", func(t *testing.T) {
		response := sendCommand("DELETE key1")
		if response != "DELETED key1" {
			t.Errorf("Expected 'DELETED key1', got '%s'", response)
		}

		// Verify key is deleted
		response = sendCommand("GET key1")
		if response != "Key not found" {
			t.Errorf("Expected 'Key not found' after delete, got '%s'", response)
		}
	})
}
