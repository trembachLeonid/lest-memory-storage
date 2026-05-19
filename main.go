package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"net"
	"strconv"
	"time"

	"github.com/trembachLeonid/lest-memory-storage/env"
	"github.com/trembachLeonid/lest-memory-storage/handlers"
	"github.com/trembachLeonid/lest-memory-storage/helpers"
	"github.com/trembachLeonid/lest-memory-storage/logging"
	"github.com/trembachLeonid/lest-memory-storage/storage"
)

func main() {
	logger := logging.InitLogger()
	defer logging.CloseLogger()

	_ = env.Load()

	logger.Info("Service started", "port", 8090, "env", "local")
	listener, err := net.Listen("tcp", ":8090")

	if err != nil {
		logger.Error("Error listening", "error", err)
	}

	defer listener.Close()

	var shards, _ = strconv.ParseInt(env.ShardCount.Get(), 10, 32)
	store := storage.NewInMemoryStorage(uint32(shards))
	list := storage.LinkedList[int64]{}
	bag := handlers.HandleBag{
		S:       store,
		ExpireL: &list,
	}

	rm := handlers.RoutineManager{
		Routines: []handlers.Routine{
			{Ticker: time.NewTicker(10 * time.Second), RoutineHandler: handlers.NewExpireRoutine(&bag)},
			{Ticker: time.NewTicker(12 * time.Second), RoutineHandler: &handlers.TestRoutine{}},
		},
	}

	go rm.Handle()

	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.Error("Error accepting conn", "error", err)
			continue
		}

		logger = logging.AddConnectionId(logger)

		logger.Info("CONNECTION ACCEPTED", "address", conn.RemoteAddr().String())

		ctx := context.WithValue(context.Background(), "logger", logger)
		go handleConnection(&ctx, conn, &bag)
	}
}

func handleConnection(ctx *context.Context, conn net.Conn, bag *handlers.HandleBag) {
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
			response := append(helpers.OK, helpers.CRLF...)
			writer.Write(response)
			writer.Flush()
			break
		}
		if err != nil && err != io.EOF {
			logger.Error("Read error", "error", err)
			return
		}

		logger = logging.AddCorrelationId(logger)
		*ctx = context.WithValue(*ctx, "logger", logger)

		logger.Info("Message received", "message", params[0])

		response, err := handlers.HandleCommand(ctx, params, bag)

		if err != nil {
			logger.Error("Error retrieving value", "error", err)
			response = helpers.OPERATION_ERROR
		}

		response = append(response, helpers.CRLF...)
		_, err = writer.Write(response)
		if err != nil {
			logger.Error("Server write error", "error", err)
		}
		writer.Flush()
	}
}
