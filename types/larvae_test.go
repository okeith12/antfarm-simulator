package types

import "testing"

func TestNewLarvae(t *testing.T) {
	larvae := NewLarvae(6, 10, 20, "Red")

	if larvae.ID != 6 {
		t.Errorf("Expected ID 6, got %d", larvae.ID)
	}
	if larvae.Role != Larvae {
		t.Errorf("Expected Role Larvae, got %d", larvae.Role)
	}
	if larvae.Health != 50 {
		t.Errorf("Expected Health 50, got %d", larvae.Health)
	}
	if larvae.HasNurseCare {
		t.Error("Expected HasNurseCare false")
	}
	if larvae.GrowthProgress != 0 {
		t.Errorf("Expected GrowthProgress 0, got %d", larvae.GrowthProgress)
	}
	if larvae.DestinedRole != Worker {
		t.Errorf("Expected DestinedRole Worker, got %d", larvae.DestinedRole)
	}
}

func TestLarvaeGetAnt(t *testing.T) {
	larvae := NewLarvae(1, 0, 0, "Red")
	if larvae.GetAnt() == nil {
		t.Error("GetAnt() should not return nil")
	}
}

func TestLarvaeGetIcon(t *testing.T) {
	larvae := NewLarvae(1, 0, 0, "Red")
	if larvae.GetIcon() != '◦' {
		t.Errorf("Expected icon '◦', got '%c'", larvae.GetIcon())
	}
}

func TestLarvaeGetRole(t *testing.T) {
	larvae := NewLarvae(1, 0, 0, "Red")
	if larvae.GetRole() != Larvae {
		t.Errorf("Expected Larvae role, got %d", larvae.GetRole())
	}
}
