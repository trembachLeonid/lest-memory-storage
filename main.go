package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/trembachLeonid/lest-memory-storage/internal/handlers"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

func main() {

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	var storage storage.Storage = storage.NewInMemoryStorage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage storage.Storage) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Printf("Read error: %v", err)
			return
		}

		var entryType = message[0]
		if entryType != '*' {
			log.Printf("Unsupported entry type: %v", entryType)
			continue
		}

		paramCount, err := strconv.Atoi(strings.TrimSpace(message[1:]))
		if err != nil {
			log.Printf("Error parsing param count: %v", err)
		}

		for i := 0; i < paramCount; i++ {
			log.Printf("Parameter %d: %v", i, message[1:])
		}
		continue

		ackMsg := strings.TrimSpace(message)
		log.Printf("Message received: %s", ackMsg)

		response, err := handlers.HandleCommand(&ackMsg, storage)

		if err != nil {
			log.Printf("Error retrieving value - %v", err)
			response = []byte("ERROR RETRIEVING VALUE")
		}

		response = append(response, '\r', '\n')
		_, err = conn.Write(response)
		if string(response) == "QUIT" {
			break
		}
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}
