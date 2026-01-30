package types

// nurse.go - Nurse ant structure
// Nurses care for larvae and help them grow into full ants

// Nurse represents a caretaking ant that tends to larvae
type NurseAnt struct {
	*Ant                        // Embedded base ant
	CurrentlyNursing *LarvaeAnt // The larvae this nurse is currently tending
	NursingSpeed     int        // How fast this nurse helps larvae grow (1-10)
	LarvaeNursed     int        // Lifetime count of larvae successfully raised
}

// NewNurse creates a new nurse ant at the specified position
func NewNurse(id int, x, y int, colonyID string) *NurseAnt {
	ant := NewAnt(id, Nurse, x, y, colonyID, NurseMaxHealth, NurseMaxTick)
	ant.Health = 100

	return &NurseAnt{
		Ant:              ant,
		CurrentlyNursing: nil,
		NursingSpeed:     1,
		LarvaeNursed:     0,
	}
}

// GetAnt returns the  Ant
func (n *NurseAnt) GetAnt() *Ant {
	return n.Ant
}

// GetIcon returns the display icon for the nurse
func (n *NurseAnt) GetIcon() rune {
	return '○'
}

// GetRole returns the nurse's role
func (n *NurseAnt) GetRole() Role {
	return Nurse
}
