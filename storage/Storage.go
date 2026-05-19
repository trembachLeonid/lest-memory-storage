package storage

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

type Storage interface {
	Set(key string, value []byte)
	Get(key string) *StorageValue
	Delete(key string)
	Increment(key string, incValue int64) ([]byte, error)
	Expire(key string, expireTime int64) error
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
	ExpireTime int64
}

type StorageShard struct {
	mu   sync.RWMutex
	data map[string]StorageValue
}

type InMemoryStorage struct {
	shardCount uint32
	shards     []StorageShard
}

func NewInMemoryStorage(shardCount uint32) *InMemoryStorage {
	s := InMemoryStorage{
		shardCount: shardCount,
		shards:     make([]StorageShard, shardCount),
	}

	for i := 0; i < int(s.shardCount); i++ {
		s.shards[i] = StorageShard{
			data: make(map[string]StorageValue),
		}
	}
	return &s
}

func (s *InMemoryStorage) Set(key string, value []byte) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = StorageValue{Type: STRING, Value: value}
}

func (s *InMemoryStorage) Get(key string) *StorageValue {
	shard := s.getShard(key)
	shard.mu.RLock()

	value, ok := shard.data[key]
	if !ok {
		shard.mu.RUnlock()
		return nil
	}

	if value.ExpireTime != 0 && value.ExpireTime < time.Now().Unix() {
		shard.mu.RUnlock()
		shard.mu.Lock()

		delete(shard.data, key)

		shard.mu.Unlock()
		return nil
	}

	shard.mu.RUnlock()
	return &value
}

func (s *InMemoryStorage) Delete(key string) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	delete(shard.data, key)
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
	shard.data[key] = obj

	return []byte(strconv.FormatInt(obj.Value.(int64), 10)), nil
}

func (s *InMemoryStorage) Expire(key string, expireTime int64) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	obj, ok := shard.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	obj.ExpireTime = expireTime
	shard.data[key] = obj
	return nil
}

func (s *InMemoryStorage) getShard(key string) *StorageShard {
	keyHash := murmur3.Sum32([]byte(key))
	shardKey := keyHash % s.shardCount
	shard := &s.shards[shardKey]
	return shard
}
