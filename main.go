package main

import (
	"bufio"
	"log"
	"net"
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

	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		ackMsg := strings.TrimSpace(message)

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
