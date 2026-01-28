package types

import "testing"

func TestNewSoldier(t *testing.T) {
	soldier := NewSoldier(4, 10, 20, "Red")

	if soldier.ID != 4 {
		t.Errorf("Expected ID 4, got %d", soldier.ID)
	}
	if soldier.Role != Soldier {
		t.Errorf("Expected Role Soldier, got %d", soldier.Role)
	}
	if soldier.Health != 150 {
		t.Errorf("Expected Health 150, got %d", soldier.Health)
	}
	if soldier.AttackPower != 20 {
		t.Errorf("Expected AttackPower 20, got %d", soldier.AttackPower)
	}
	if soldier.DefenseBonus != 10 {
		t.Errorf("Expected DefenseBonus 10, got %d", soldier.DefenseBonus)
	}
	if soldier.IsPatrolling {
		t.Error("Expected IsPatrolling false")
	}
}

func TestSoldierGetAnt(t *testing.T) {
	soldier := NewSoldier(1, 0, 0, "Red")
	if soldier.GetAnt() == nil {
		t.Error("GetAnt() should not return nil")
	}
}

func TestSoldierGetAntIcon(t *testing.T) {
	soldier := NewSoldier(1, 0, 0, "Red")
	if soldier.GetIcon() != '⚔' {
		t.Errorf("Expected icon '⚔', got '%c'", soldier.GetIcon())
	}
}

func TestSoldierGetRole(t *testing.T) {
	soldier := NewSoldier(1, 0, 0, "Red")
	if soldier.GetRole() != Soldier {
		t.Errorf("Expected Soldier role, got %d", soldier.GetRole())
	}
}
