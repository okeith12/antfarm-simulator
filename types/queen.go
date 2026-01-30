package types

// queen.go
// The queen is the heart of the colony, responsible for laying eggs

type QueenAnt struct {
	*Ant                  // Embedded base ant
	EggLayingCooldown int // Ticks until queen can lay another egg
	TotalEggsLaid     int // Lifetime egg count
}

// NewQueen creates a new queen ant at the specified position
func NewQueen(id int, x, y int, colonyID string) *QueenAnt {
	ant := NewAnt(id, Queen, x, y, colonyID, QueenMaxHealth, QueenMaxTick)
	ant.Health = 200

	return &QueenAnt{
		Ant:               ant,
		EggLayingCooldown: 0,
		TotalEggsLaid:     0,
	}
}

// GetAnt returns the Ant
func (q *QueenAnt) GetAnt() *Ant {
	return q.Ant
}

// GetIcon returns the display icon for the queen
func (q *QueenAnt) GetIcon() rune {
	return '♛'
}

// GetRole returns the queen's role
func (q *QueenAnt) GetRole() Role {
	return Queen
}
