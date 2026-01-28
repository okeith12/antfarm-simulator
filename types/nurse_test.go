package types

import "testing"

func TestNewNurse(t *testing.T) {
	nurse := NewNurse(1, 10, 20, "Red")

	if nurse.ID != 1 {
		t.Errorf("Expected ID 1, got %d", nurse.ID)
	}
	if nurse.Role != Nurse {
		t.Errorf("Expected Role Nurse, got %d", nurse.Role)
	}
	if nurse.Health != 100 {
		t.Errorf("Expected Health 100, got %d", nurse.Health)
	}
	if nurse.CurrentlyNursing != nil {
		t.Error("Expected CurrentlyNursing nil")
	}
	if nurse.NursingSpeed != 1 {
		t.Errorf("Expected NursingSpeed 1, got %d", nurse.NursingSpeed)
	}
	if nurse.LarvaeNursed != 0 {
		t.Errorf("Expected LarvaeNursed 0, got %d", nurse.LarvaeNursed)
	}
}

func TestNurseGetAnt(t *testing.T) {
	nurse := NewNurse(1, 0, 0, "Red")
	if nurse.GetAnt() == nil {
		t.Error("GetAnt() should not return nil")
	}
}

func TestNurseGetAntIcon(t *testing.T) {
	nurse := NewNurse(1, 0, 0, "Red")
	if nurse.GetIcon() != '○' {
		t.Errorf("Expected icon '○', got '%c'", nurse.GetIcon())
	}
}

func TestNurseGetRole(t *testing.T) {
	nurse := NewNurse(1, 0, 0, "Red")
	if nurse.GetRole() != Nurse {
		t.Errorf("Expected Nurse role, got %d", nurse.GetRole())
	}
}
