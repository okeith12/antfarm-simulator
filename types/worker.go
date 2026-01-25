package types

// worker.go
// Workers dig tunnels, gather food, and expand the colony

// Worker represents a laboring ant that digs and forages
type WorkerAnt struct {
	*Ant
	CarryingFood   bool      // Is this worker carrying food?
	FoodAmount     int       // How much food is being carried
	DiggingPower   int       // How fast this worker digs (1-10)
	TargetPosition *Position // Where the worker is trying to go
}

// NewWorker creates a new worker ant at the specified position
func NewWorker(id int, x, y int, colonyID string) *WorkerAnt {
	ant := NewAnt(id, Worker, x, y, colonyID)
	ant.Health = 100

	return &WorkerAnt{
		Ant:            ant,
		CarryingFood:   false,
		FoodAmount:     0,
		DiggingPower:   1,
		TargetPosition: nil,
	}
}

// GetAnt returns the Ant
func (w *WorkerAnt) GetAnt() *Ant {
	return w.Ant
}

// GetIcon returns the display icon for the worker
func (w *WorkerAnt) GetAntIcon() rune {
	return '●'
}

// GetRole returns the worker's role
func (w *WorkerAnt) GetRole() Role {
	return Worker
}
