package handlers

import (
	"log"
	"strconv"
	"strings"

	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

var unknownCommand = []byte("unknown command")
var operationError = []byte("operation error")
var quit = []byte("BYE")
var ok = []byte("OK")
var pong = []byte("PONG")

func HandleCommand(command *string, storage storage.Storage) ([]byte, error) {
	var err error
	var response []byte = ok
	var commandParts []string
	if command != nil {
		commandParts = strings.Split(*command, " ")
	} else {
		return unknownCommand, nil
	}
	action := strings.ToUpper(commandParts[0])

	if action == "QUIT" {
		return quit, nil
	}

	var value *[]byte
	key := commandParts[1]

	if len(commandParts) > 2 {
		valueBytes := []byte(commandParts[2])
		value = &valueBytes
	}

	switch action {
	case "PING":
		response, err = pong, nil
	case "SET":
		err = storage.Set(key, value)
	case "GET":
		response, err = storage.Get(key)
	case "DEL":
		err = storage.Delete(key)
	case "INC":
		if value == nil || len(*value) == 0 {
			value = &[]byte{'1'}
		}

		incValue, err := strconv.ParseInt(string(*value), 10, 64)
		if err != nil {
			return operationError, err
		}

		response, err = storage.Increment(key, incValue)
	case "DEC":
		if value == nil || len(*value) == 0 {
			value = &[]byte{'1'}
		}

		incValue, err := strconv.ParseInt(string(*value), 10, 64)
		if err != nil {
			return operationError, err
		}
		response, err = storage.Increment(key, -incValue)
	case "CONFIG":
		return []byte("*0\r\n"), nil
	default:
		response, err = unknownCommand, nil
	}

	log.Printf("HANDLE - %s, { \"%s\": \"%v\" } - %s", action, key, value, string(response))
	return response, err
}
