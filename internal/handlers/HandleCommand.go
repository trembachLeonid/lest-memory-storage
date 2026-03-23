package handlers

import (
	"log"
	"strconv"

	"github.com/trembachLeonid/lest-memory-storage/internal/models"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

var unknownCommand = []byte("unknown command")
var operationError = []byte("operation error")
var ok = []byte("OK")
var pong = []byte("PONG")

func HandleCommand(action models.ActionType, key string, value *[]byte, storage storage.Storage) ([]byte, error) {
	var err error
	var response []byte = ok

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
	default:
		response, err = unknownCommand, nil
	}

	log.Printf("HANDLE - %s, \"%s\", \"%s\" - %s", action, key, string(*value), string(response))
	return response, err
}
