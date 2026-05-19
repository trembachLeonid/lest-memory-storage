package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/trembachLeonid/lest-memory-storage/helpers"
	"github.com/trembachLeonid/lest-memory-storage/storage"
)

type commandHandler func(command [][]byte, bag *HandleBag) ([]byte, error)

var commandHandlers = map[string]commandHandler{
	string(helpers.PING):   handlePing,
	string(helpers.QUIT):   func(_ [][]byte, _ *HandleBag) ([]byte, error) { return helpers.QUIT, nil },
	string(helpers.SET):    handleSet,
	string(helpers.GET):    handleGet,
	string(helpers.DEL):    handleDelete,
	string(helpers.INC):    handleInc,
	string(helpers.DEC):    handleDec,
	string(helpers.CONFIG): handleConfig,
	string(helpers.EXPIRE): handleExpire,
}

func HandleCommand(ctx *context.Context, command [][]byte, bag *HandleBag) ([]byte, error) {
	action := string(command[0])
	handler, ok := commandHandlers[action]
	if !ok {
		return helpers.UNKNOWN_COMMAND, nil
	}

	response, err := handler(command, bag)
	return response, err
}

func handlePing(_ [][]byte, _ *HandleBag) ([]byte, error) {
	return helpers.PONG, nil
}

func handleSet(command [][]byte, bag *HandleBag) ([]byte, error) {
	key := string(command[1])
	fmt.Printf("Setting key: %s\n", key)
	value := command[2]
	bag.S.Set(key, value)
	return helpers.OK, nil
}

func handleGet(command [][]byte, bag *HandleBag) ([]byte, error) {
	storageValue := bag.S.Get(string(command[1]))
	if storageValue == nil {
		return []byte("$-1"), nil
	}
	switch storageValue.Type {
	case storage.STRING:
		return storageValue.Value.([]byte), nil
	case storage.INTEGER:
		return []byte(strconv.FormatInt(storageValue.Value.(int64), 10)), nil
	default:
		return []byte("$-1"), nil
	}
}

func handleDelete(command [][]byte, bag *HandleBag) ([]byte, error) {
	bag.S.Delete(string(command[1]))
	return helpers.OK, nil
}

func handleExpire(command [][]byte, bag *HandleBag) ([]byte, error) {
	seconds, err := strconv.ParseInt(string(command[2]), 10, 64)
	if err != nil {
		return helpers.OPERATION_ERROR, err
	}
	key := string(command[1])
	expiryTime := time.Now().UTC().Unix() + seconds

	err = bag.S.Expire(key, expiryTime)
	if err != nil {
		return helpers.OPERATION_ERROR, err
	}

	bag.ExpireL.Append(key, expiryTime)
	return helpers.OK, err
}

// func handleSetExpire(command [][]byte, bag *HandleBag) ([]byte, error) {
// 	value := command[2]
// 	bag.S.Set(string(command[1]), value)
// 	seconds, err := strconv.ParseInt(string(command[3]), 10, 64)

// 	key := string(command[1])
// 	expiryTime := time.Now().UTC().Unix() + seconds
// 	err = bag.S.Expire(key, expiryTime)
// 	if err != nil {
// 		return helpers.OPERATION_ERROR, err
// 	}

// 	bag.ExpireL.Append(key, expiryTime)
// 	return helpers.OK, err
// }

// func handleGetTtl(command [][]byte, bag *HandleBag) ([]byte, error) {
// 	currTime := time.Now().UTC().Unix()

// }

func handleInc(command [][]byte, bag *HandleBag) ([]byte, error) {
	by := []byte{'1'}
	if len(command) > 2 && len(command[2]) > 0 {
		by = command[2]
	}
	n, err := strconv.ParseInt(string(by), 10, 64)
	if err != nil {
		return helpers.OPERATION_ERROR, err
	}
	return bag.S.Increment(string(command[1]), n)
}

func handleDec(command [][]byte, bag *HandleBag) ([]byte, error) {
	by := []byte{'1'}
	if len(command) > 2 && len(command[2]) > 0 {
		by = command[2]
	}
	n, err := strconv.ParseInt(string(by), 10, 64)
	if err != nil {
		return helpers.OPERATION_ERROR, err
	}
	return bag.S.Increment(string(command[1]), -n)
}

func handleConfig(_ [][]byte, _ *HandleBag) ([]byte, error) {
	return []byte("*0"), nil
}
