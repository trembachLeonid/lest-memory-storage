package storage

import (
	"fmt"
	"log"
	"strconv"
	"sync"
)

type Storage interface {
	Set(key string, value *[]byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Increment(key string, incValue int64) ([]byte, error)
}

type ValueType int

const (
	STRING ValueType = iota
	INTEGER
	FLOAT
	BOOLEAN
)

type StorageValue struct {
	Type  ValueType
	Value any
}

type InMemoryStorage struct {
	mu   sync.RWMutex
	data map[string]*StorageValue
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]*StorageValue, 10000),
	}
}

func (s *InMemoryStorage) Set(key string, value *[]byte) error {
	s.data[key] = &StorageValue{Type: STRING, Value: *value}
	return nil
}

func (s *InMemoryStorage) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key not found: %s", key)
	}

	switch value.Type {
	case STRING:
		return value.Value.([]byte), nil
	case INTEGER:
		return []byte(strconv.FormatInt(value.Value.(int64), 10)), nil
	default:
		log.Printf("Unsupported type for key %s: %v", key, value.Type)
	}

	return []byte{}, fmt.Errorf("key not found: %s", key)
}

func (s *InMemoryStorage) Delete(key string) error {
	delete(s.data, key)
	return nil
}

func (s *InMemoryStorage) Increment(key string, incValue int64) ([]byte, error) {
	s.mu.Lock()

	obj, ok := s.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key not found: %s", key)
	}

	switch obj.Type {
	case STRING:
		converted, err := strconv.ParseInt(string(obj.Value.([]byte)), 10, 64)
		if err != nil {
			return []byte{}, fmt.Errorf("value is not a number: %s", key)
		}
		converted += incValue
		obj.Type = INTEGER
		obj.Value = converted
	case INTEGER:
		newVal := obj.Value.(int64) + incValue
		obj.Value = newVal
	default:
		return []byte{}, fmt.Errorf("unsupported type: %v", obj.Type)
	}

	s.mu.Unlock()

	return s.Get(key)
}
