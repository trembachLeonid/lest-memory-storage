package storage

import (
	"fmt"
	"strconv"
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
			data:       make(map[string]*StorageNode),
			expiryList: &LinkedList[*StorageNode]{},
		}
	}
	return &s
}

func (s *InMemoryStorage) Set(key string, value []byte) {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.data[key] = &StorageNode{Key: key, SV: StorageValue{Type: STRING, Value: value}}
}

func (s *InMemoryStorage) Get(key string) *StorageValue {
	shard := s.getShard(key)
	shard.mu.RLock()

	node, ok := shard.data[key]
	if !ok {
		shard.mu.RUnlock()
		return nil
	}

	if node.ExpireTime != 0 && node.ExpireTime < time.Now().Unix() {
		shard.mu.RUnlock()
		shard.mu.Lock()

		delete(shard.data, key)

		shard.mu.Unlock()
		return nil
	}

	shard.mu.RUnlock()
	return &node.SV
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

	node, ok := shard.data[key]
	if !ok {
		return []byte{}, fmt.Errorf("key not found: %s", key)
	}

	switch node.SV.Type {
	case STRING:
		converted, err := strconv.ParseInt(string(node.SV.Value.([]byte)), 10, 64)
		if err != nil {
			return []byte{}, fmt.Errorf("value is not a number: %s", key)
		}
		converted += incValue
		node.SV.Type = INTEGER
		node.SV.Value = converted
	case INTEGER:
		newVal := node.SV.Value.(int64) + incValue
		node.SV.Value = newVal
	default:
		return []byte{}, fmt.Errorf("unsupported type: %v", node.SV.Type)
	}
	shard.data[key] = node

	return []byte(strconv.FormatInt(node.SV.Value.(int64), 10)), nil
}

func (s *InMemoryStorage) Expire(key string, expireTime int64) error {
	shard := s.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()

	node, ok := shard.data[key]
	if !ok {
		return fmt.Errorf("key not found: %s", key)
	}

	// TODO: Ensure that the same key is not added multiple times to the expiry list
	if node.ExpireTime == 0 {
		shard.expiryList.Append(node)
	}
	node.ExpireTime = expireTime
	return nil
}

func (s *InMemoryStorage) GetExpiryList(shardIndex int) List[*StorageNode] {
	return s.shards[shardIndex].expiryList
}

func (s *InMemoryStorage) getShard(key string) *StorageShard {
	keyHash := murmur3.Sum32([]byte(key))
	shardKey := keyHash % s.shardCount
	shard := &s.shards[shardKey]
	return shard
}
