package gui

import (
	"antfarm/types"
	"testing"
)

// stats_test.go - Tests for stats rendering functions

func TestGetRoleString(t *testing.T) {
	tests := []struct {
		name     string
		role     types.Role
		expected string
	}{
		{"Worker role", types.Worker, "Worker"},
		{"Soldier role", types.Soldier, "Soldier"},
		{"Queen role", types.Queen, "Queen"},
		{"Nurse role", types.Nurse, "Nurse"},
		{"Larvae role", types.Larvae, "Larvae"},
		{"Unknown role", types.Role(99), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getRoleString(tt.role)
			if result != tt.expected {
				t.Errorf("getRoleString(%d) = %s; want %s", tt.role, result, tt.expected)
			}
		})
	}
}
