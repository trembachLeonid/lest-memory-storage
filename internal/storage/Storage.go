package storage

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

var val, _ = strconv.ParseUint(os.Getenv("SHARD_COUNT"), 10, 32)
var shardCount = uint32(1024)

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
	Type       ValueType
	Value      any
	ExpireTime time.Time
}

type StorageShard struct {
	mu   sync.RWMutex
	data map[string]*StorageValue
}

type InMemoryStorage struct {
	shards map[uint32]*StorageShard
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		shards: make(map[uint32]*StorageShard, shardCount),
	}
}

func (s *InMemoryStorage) Set(key string, value *[]byte) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = &StorageValue{Type: STRING, Value: *value}
	return nil
}

func (s *InMemoryStorage) Get(key string) ([]byte, error) {
	shard := s.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()

	value, ok := shard.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key not found: %s", key)
	}
	if !value.ExpireTime.IsZero() && value.ExpireTime.Before(time.Now()) {
		s.Delete(key)
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
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.data, key)

	return nil
}

func (s *InMemoryStorage) Increment(key string, incValue int64) ([]byte, error) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	obj, ok := shard.data[key]
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

	return []byte(strconv.FormatInt(obj.Value.(int64), 10)), nil
}

func (s *InMemoryStorage) Expire(key string, expireTime time.Time) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	obj, ok := shard.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	obj.ExpireTime = expireTime
	return nil
}

func (s *InMemoryStorage) getShard(key string) *StorageShard {
	keyHash := murmur3.Sum32([]byte(key))
	shardKey := keyHash % shardCount
	shard, ok := s.shards[shardKey]
	if !ok {
		s.shards[shardKey] = &StorageShard{
			data: make(map[string]*StorageValue),
		}
		shard = s.shards[shardKey]
	}
	log.Printf("getShard - key = %s , keyHash = %v , shard = %v", key, keyHash, shardKey)
	return shard
}
