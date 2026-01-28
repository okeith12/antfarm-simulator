package types

import "testing"

func TestNewQueen(t *testing.T) {
	queen := NewQueen(0, 10, 20, "Red")

	if queen.ID != 0 {
		t.Errorf("Expected ID 0, got %d", queen.ID)
	}
	if queen.Role != Queen {
		t.Errorf("Expected Role Queen, got %d", queen.Role)
	}
	if queen.Health != 200 {
		t.Errorf("Expected Health 200, got %d", queen.Health)
	}
	if queen.EggLayingCooldown != 0 {
		t.Errorf("Expected EggLayingCooldown 0, got %d", queen.EggLayingCooldown)
	}
	if queen.TotalEggsLaid != 0 {
		t.Errorf("Expected TotalEggsLaid 0, got %d", queen.TotalEggsLaid)
	}
}

func TestQueenGetAnt(t *testing.T) {
	queen := NewQueen(0, 0, 0, "Red")
	if queen.GetAnt() == nil {
		t.Error("GetAnt() should not return nil")
	}
}

func TestQueenGetAntIcon(t *testing.T) {
	queen := NewQueen(0, 0, 0, "Red")
	if queen.GetIcon() != '♛' {
		t.Errorf("Expected icon '♛', got '%c'", queen.GetIcon())
	}
}

func TestQueenGetRole(t *testing.T) {
	queen := NewQueen(0, 0, 0, "Red")
	if queen.GetRole() != Queen {
		t.Errorf("Expected Queen role, got %d", queen.GetRole())
	}
}
