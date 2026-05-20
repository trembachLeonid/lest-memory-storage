package storage

import "sync"

type StorageShard struct {
	mu         sync.RWMutex
	data       map[string]*StorageNode
	expiryList List[*StorageNode]
}
