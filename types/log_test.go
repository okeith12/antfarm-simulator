package types

import "testing"

func TestGetActivityString(t *testing.T) {
	ant := NewAnt(5, Worker, 10, 20, "Red")

	result := ant.GetActivityString(1)
	expected := "Red_Worker_Ant_5 is currently at (10,20) in World_1"

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestGetActivityStringDifferentRoles(t *testing.T) {
	tests := []struct {
		role     Role
		roleName string
	}{
		{Worker, "Worker"},
		{Soldier, "Soldier"},
		{Nurse, "Nurse"},
		{Queen, "Queen"},
		{Larvae, "Larvae"},
	}

	for _, tt := range tests {
		ant := NewAnt(1, tt.role, 0, 0, "Test")
		result := ant.GetActivityString(1)

		if result == "" {
			t.Errorf("GetActivityString returned empty for role %s", tt.roleName)
		}
	}
}
