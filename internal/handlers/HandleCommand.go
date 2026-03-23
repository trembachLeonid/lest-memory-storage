package handlers

import (
	"github.com/trembachLeonid/lest-memory-storage/internal/models"
	"github.com/trembachLeonid/lest-memory-storage/internal/storage"
)

func HandleCommand(action models.ActionType, key string, value string, storage *storage.InMemoryStorage) (any, error) {
	var err error
	switch action {
	case "PING":
		return "PONG", nil
	case "SET":
		err = storage.Set(key, value)
		return "OK", nil
	case "GET":
		res, err := storage.Get(key)
		return res, err
	}

	return "unknown command", err
}
