package types

// larvae.go - Larvae structure
// Larvae are baby ants that need nurse care to grow into workers

// Larvae represents a baby ant that hasn't matured yet
type LarvaeAnt struct {
	*Ant                // Embedded base ant
	HasNurseCare   bool // Has a nurse tended to this larvae?
	GrowthProgress int  // Progress toward becoming a full ant (0-100)
	DestinedRole   Role // What role this larvae will become (usually Worker)
}

// NewLarvae creates a new larvae at the specified position
func NewLarvae(id int, x, y int, colonyID string) *LarvaeAnt {
	ant := NewAnt(id, Larvae, x, y, colonyID)
	ant.Health = 50 // Larvae are fragile

	return &LarvaeAnt{
		Ant:            ant,
		HasNurseCare:   false,
		GrowthProgress: 0,
		DestinedRole:   Worker, // Default to becoming a worker
	}
}

// GetAnt returns the embedded base Ant
func (l *LarvaeAnt) GetAnt() *Ant {
	return l.Ant
}

// GetIcon returns the display icon for the larvae
func (l *LarvaeAnt) GetIcon() rune {
	return '◦'
}

// GetRole returns the larvae's role
func (l *LarvaeAnt) GetRole() Role {
	return Larvae
}
