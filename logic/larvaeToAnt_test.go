package logic

import (
	"fmt"
	"testing"

	"antfarm/types"
)

func TestLarvaeToAnt_RoleSelection(t *testing.T) {
	colony := &types.Colony{}
	larvae := types.NewLarvae(1, 0, 0, "Red")

	tests := []struct {
		name        string
		roll        int
		expectedTyp string
	}{
		{"Nurse lower bound", 0, "*types.NurseAnt"},
		{"Nurse upper bound", 19, "*types.NurseAnt"},
		{"Soldier lower bound", 20, "*types.SoldierAnt"},
		{"Soldier upper bound", 34, "*types.SoldierAnt"},
		{"Worker lower bound", 35, "*types.WorkerAnt"},
		{"Worker upper bound", 99, "*types.WorkerAnt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ant := larvaeToAnt(colony, larvae, tt.roll)

			if ant == nil {
				t.Fatalf("expected ant, got nil")
			}

			actualType := getTypeName(ant)
			if actualType != tt.expectedTyp {
				t.Errorf("expected type %s, got %s", tt.expectedTyp, actualType)
			}
		})
	}
}

func getTypeName(a types.AntInterface) string {
	return fmt.Sprintf("%T", a)
}
