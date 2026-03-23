package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"github.com/trembachLeonid/lest-memory-storage/internal/handlers"
	"github.com/trembachLeonid/lest-memory-storage/internal/models"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

func main() {

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	var storage = storage.NewInMemoryStorage()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go handleConnection(conn, storage)
	}
}

func handleConnection(conn net.Conn, storage *storage.InMemoryStorage) {

	defer conn.Close()

	for {
		reader := bufio.NewReader(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Read error: %v", err)
			return
		}

		ackMsg := strings.TrimSpace(message)
		commandParts := strings.Split(ackMsg, " ")

		var response []byte

		if commandParts[0] == "PING" {
			response = []byte("PONG")
		} else if commandParts[0] == "SET" && len(commandParts) == 3 {
			response, err = handlers.HandleCommand(models.ActionType(commandParts[0]), commandParts[1], []byte(commandParts[2]), storage)
		} else {
			response = []byte("UNKNOWN COMMAND")
		}

		_, err = conn.Write(response)
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}
