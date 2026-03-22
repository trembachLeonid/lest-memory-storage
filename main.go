package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

func main() {

	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	defer listener.Close()

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

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

		var response string

		if commandParts[0] == "PING" {
			response = "PONG"
		} else if commandParts[0] == "SET" && len(commandParts) == 3 {
			response = "SET command received"
		} else {
			response = "UNKNOWN COMMAND"
		}

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Printf("Server write error: %v", err)
		}
	}
}
