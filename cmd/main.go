package main

import (
	"bufio"
	"cache/pkg"
	"fmt"
	"net"
	"strings"
)

func main() {
	cache := db.NewCache(1024, db.NewLRU())

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn, cache)
	}
}

func handleConnection(conn net.Conn, cache *db.Cache) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read command from client.
		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Connection closed.")
			return
		}

		// Process the command.
		command = strings.TrimSpace(command)
		parts := strings.SplitN(command, " ", 3)

		if len(parts) < 2 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}

		action, key := strings.ToUpper(parts[0]), parts[1]

		switch action {
		case "SET":
			if len(parts) != 3 {
				conn.Write([]byte("Usage: SET <key> <value>\n"))
				continue
			}
			value := parts[2]
			cache.AddEntry(key, value)
			conn.Write([]byte(fmt.Sprintf("SET %s OK\n", key)))

		case "GET":
			value, exists := cache.Get(key)
			if exists {
				conn.Write([]byte(fmt.Sprintf("%s\n", value)))
			} else {
				conn.Write([]byte("Key not found\n"))
			}

		case "DELETE":
			cache.Delete(key)
			conn.Write([]byte(fmt.Sprintf("DELETED %s\n", key)))

		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}
