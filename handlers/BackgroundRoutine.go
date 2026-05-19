package handlers

import (
	"fmt"
	"time"

	"github.com/trembachLeonid/lest-memory-storage/storage"
)

type Routine interface {
	ExecuteRoutine()
}

func (bag *HandleBag) ExecuteRoutine() {
	l := bag.ExpireL
	s := bag.S
	currentTimestamp := time.Now().UTC().Unix()
	node := l.GetHead()
	var prev *storage.LinkedNode[int64]

	for node != nil {
		if node.Value < currentTimestamp {
			s.Delete(node.Key)
			if node == l.GetHead() {
				fmt.Println("Yes remove head")
				l.RemoveHead()
			} else {
				l.RemoveNext(prev)
			}
		} else {
			prev = node
		}
		node = node.Next
	}

	fmt.Println(l.(*storage.LinkedList[int64]).Head)
	fmt.Println(l.(*storage.LinkedList[int64]).Tail)
}
