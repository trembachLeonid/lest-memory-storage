package handlers

type RoutineManager struct {
	Routines []Routine
}

func (rm *RoutineManager) Handle() {
	for _, sr := range rm.Routines {
		go sr.runRoutine()
	}
}
