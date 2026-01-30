package types

// soldier.go
// Soldiers defend the colony and fight enemy ants

// Soldier represents a combat-focused ant
type SoldierAnt struct {
	*Ant                     // Embedded base ant
	AttackPower    int       // Damage dealt per attack
	DefenseBonus   int       // Damage reduction
	IsPatrolling   bool      // Is this soldier on patrol duty?
	TargetPosition *Position // Where the soldier is heading
}

// NewSoldier creates a new soldier ant at the specified position
func NewSoldier(id int, x, y int, colonyID string) *SoldierAnt {
	ant := NewAnt(id, Soldier, x, y, colonyID, SoldierMaxHealth, SoldierMaxTick)
	ant.Health = 150 // Soldiers have more health than workers

	return &SoldierAnt{
		Ant:            ant,
		AttackPower:    20,
		DefenseBonus:   10,
		IsPatrolling:   false,
		TargetPosition: nil,
	}
}

// GetAnt returns the embedded base Ant
func (s *SoldierAnt) GetAnt() *Ant {
	return s.Ant
}

// GetIcon returns the display icon for the soldier
func (s *SoldierAnt) GetIcon() rune {
	return '⚔'
}

// GetRole returns the soldier's role
func (s *SoldierAnt) GetRole() Role {
	return Soldier
}
