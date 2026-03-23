package handlers

import (
	"github.com/trembachLeonid/lest-memory-storage/internal/models"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

func HandleCommand(action models.ActionType, key string, value []byte, storage *storage.InMemoryStorage) ([]byte, error) {
	var err error
	switch action {
	case "PING":
		return []byte("PONG"), nil
	case "SET":
		err = storage.Set(key, value)
		return []byte("OK"), nil
	case "GET":
		res, err := storage.Get(key)
		return []byte(res), err
	}

	return []byte("unknown command"), err
}
