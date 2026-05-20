package routine

import (
	"fmt"
	"time"

	"github.com/trembachLeonid/lest-memory-storage/storage"
)

type RoutineHandler interface {
	ExecuteRoutine()
}

type Routine struct {
	Ticker         *time.Ticker
	RoutineHandler RoutineHandler
}

func (sr *Routine) runRoutine() {
	defer sr.Ticker.Stop()
	for range sr.Ticker.C {
		go sr.RoutineHandler.ExecuteRoutine()
	}
}

type ExpireRoutine struct {
	Storage *storage.InMemoryStorage
}

type TestRoutine struct{}

func (r *TestRoutine) ExecuteRoutine() {
	fmt.Println("Test routine executed at", time.Now().UTC())
}
