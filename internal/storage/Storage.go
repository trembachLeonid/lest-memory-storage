package storage

import (
	"fmt"
)

type Storage interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Delete(key string) error
	Increment(key string) (int, error)
}

type InMemoryStorage struct {
	data map[string]string
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Set(key string, value string) error {
	s.data[key] = value
	return nil
}

func (s *InMemoryStorage) Get(key string) (string, error) {
	if value, ok := s.data[key]; ok {
		return value, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}

func (s *InMemoryStorage) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *InMemoryStorage) Increment(key string) (int, error) {

	return 1, nil
}
