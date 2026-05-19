package handlers

import (
	"fmt"
	"time"

	"github.com/trembachLeonid/lest-memory-storage/storage"
)

type RoutineHandler interface {
	execute()
}

type Routine struct {
	Ticker         *time.Ticker
	RoutineHandler RoutineHandler
}

func (sr *Routine) runRoutine() {
	defer sr.Ticker.Stop()
	for range sr.Ticker.C {
		go sr.RoutineHandler.execute()
	}
}

type ExpireRoutine struct {
	Bag *HandleBag
}

type TestRoutine struct{}

func NewExpireRoutine(bag *HandleBag) *ExpireRoutine {
	return &ExpireRoutine{Bag: bag}
}

func (r *ExpireRoutine) execute() {
	l := r.Bag.ExpireL
	s := r.Bag.S
	currentTimestamp := time.Now().UTC().Unix()
	node := l.GetHead()
	var prev *storage.LinkedNode[int64]

	for node != nil {
		if node.Value < currentTimestamp {
			s.Delete(node.Key)
			if node == l.GetHead() {
				l.RemoveHead()
			} else {
				l.RemoveNext(prev)
			}
		} else {
			prev = node
		}
		node = node.Next
	}
}

func (r *TestRoutine) execute() {
	fmt.Println("Test routine executed at", time.Now().UTC())
}
