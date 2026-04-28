package handlers

import (
	"bytes"
	"context"
	"log"
	"strconv"

	"github.com/trembachLeonid/lest-memory-storage/helpers"
	"github.com/trembachLeonid/lest-memory-storage/logging"
	"github.com/trembachLeonid/lest-memory-storage/storage"
)

func HandleCommand(ctx *context.Context, command [][]byte, storage storage.Storage) ([]byte, error) {
	logger := logging.FromContext(*ctx)
	logger.Info("HANDLE COMMAND", "command", command)

	var err error
	var response []byte = helpers.OK

	action := command[0]

	var value *[]byte
	key := string(command[1])

	if len(command) > 2 {
		valueBytes := command[2]
		value = &valueBytes
	}

	if bytes.Equal(action, helpers.PING) {
		response, err = helpers.PONG, nil
	} else if bytes.Equal(action, helpers.SET) {
		err = storage.Set(key, value)
	} else if bytes.Equal(action, helpers.GET) {
		response, err = storage.Get(key)
	} else if bytes.Equal(action, helpers.DEL) {
		err = storage.Delete(key)
	} else if bytes.Equal(action, helpers.INC) {
		if value == nil || len(*value) == 0 {
			value = &[]byte{'1'}
		}

		incValue, err := strconv.ParseInt(string(*value), 10, 64)
		if err != nil {
			return helpers.OPERATION_ERROR, err
		}

		response, err = storage.Increment(key, incValue)
	} else if bytes.Equal(action, helpers.DEC) {
		if value == nil || len(*value) == 0 {
			value = &[]byte{'1'}
		}

		incValue, err := strconv.ParseInt(string(*value), 10, 64)
		if err != nil {
			return helpers.OPERATION_ERROR, err
		}
		response, err = storage.Increment(key, -incValue)
	} else if bytes.Equal(action, helpers.CONFIG) {
		return []byte("*0\r\n"), nil
	} else {
		response, err = helpers.UNKNOWN_COMMAND, nil
	}

	log.Printf("HANDLE - %s, { \"%s\": \"%v\" } - %s", action, key, value, string(response))
	return response, err
}
