package storage

import (
	"fmt"
	"time"
)

func (s *InMemoryStorage) ExecuteRoutine() {
	currentTimestamp := time.Now().UTC().Unix()

	for i := 0; i < len(s.shards); i++ {
		shard := &s.shards[i]
		shard.mu.Lock()

		node := shard.expiryList.GetHead()
		var prev *LinkedNode[*StorageNode]

		for node != nil {
			fmt.Printf("Checking expiry for key: %s, expire time: %d", node.Value.Key, node.Value.ExpireTime)
			if node.Value.ExpireTime < currentTimestamp {
				delete(shard.data, node.Value.Key)
				fmt.Printf("Key expired and removed: %s", node.Value.Key)
				if node == shard.expiryList.GetHead() {
					shard.expiryList.RemoveHead()
				} else {
					shard.expiryList.RemoveNext(prev)
				}
			} else {
				prev = node
			}
			node = node.Next
		}

		shard.mu.Unlock()
	}
}
