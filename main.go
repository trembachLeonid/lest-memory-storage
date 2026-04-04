package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"

	"github.com/trembachLeonid/lest-memory-storage/internal/handlers"
	"github.com/trembachLeonid/lest-memory-storage/internal/helpers"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
	"github.com/trembachLeonid/lest-memory-storage/logging"
)

func main() {
	logger := logging.InitLogger()
	defer logging.CloseLogger()

	logger.Info("Service started", "port", 8090, "env", "local")
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		logger.Error("Error listening", "error", err)
	}

	defer listener.Close()

	var storage storage.Storage = storage.NewInMemoryStorage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Error accepting conn", "error", err)
			continue
		}

		logger = logging.AddConnectionId(logger)

		logger.Info("CONNECTION ACCEPTED", "address", conn.RemoteAddr().String())

		ctx := context.WithValue(context.Background(), "logger", logger)
		go handleConnection(&ctx, conn, storage)
	}
}

func handleConnection(ctx *context.Context, conn net.Conn, storage storage.Storage) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	parser := helpers.NewCommandParser(reader)
	writer := bufio.NewWriter(conn)

	if tcp, ok := conn.(*net.TCPConn); ok {
		tcp.SetNoDelay(true)
	}

	logger := logging.FromContext(*ctx)

	for {
		params, err := parser.Parse(ctx)
		if bytes.Equal(params[0], helpers.QUIT) {
			break
		}
		if err != nil && err != io.EOF {
			logger.Error("Read error", "error", err)
			return
		}

		logger = logging.AddCorrelationId(logger)
		*ctx = context.WithValue(*ctx, "logger", logger)

		logger.Info("Message received", "message", params[0])

		response, err := handlers.HandleCommand(ctx, params, storage)

		if err != nil {
			logger.Error("Error retrieving value", "error", err)
			response = helpers.OPERATION_ERROR
		}

		response = append(response, '\r', '\n')
		_, err = writer.Write(response)
		if err != nil {
			logger.Error("Server write error", "error", err)
		}
		writer.Flush()
	}
}
