package logic

import (
	"antfarm/types"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestSpawnWorker(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialID := colony.NextAntID

	worker := SpawnWorker(colony, 15, 15)

	if worker == nil {
		t.Fatal("SpawnWorker returned nil")
	}
	if worker.ID != initialID {
		t.Errorf("Expected ID %d, got %d", initialID, worker.ID)
	}
	if colony.NextAntID != initialID+1 {
		t.Error("NextAntID should increment")
	}
	if len(colony.Workers) != 1 {
		t.Errorf("Expected 1 worker, got %d", len(colony.Workers))
	}
}

func TestSpawnWorkerWithID(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialNextID := colony.NextAntID

	worker := SpawnWorkerWithID(colony, 999, 15, 15)

	if worker.ID != 999 {
		t.Errorf("Expected ID 999, got %d", worker.ID)
	}
	if colony.NextAntID != initialNextID {
		t.Error("NextAntID should not change with SpawnWorkerWithID")
	}
}

func TestSpawnSoldier(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialID := colony.NextAntID

	soldier := SpawnSoldier(colony, 12, 12)

	if soldier == nil {
		t.Fatal("SpawnSoldier returned nil")
	}
	if soldier.ID != initialID {
		t.Errorf("Expected ID %d, got %d", initialID, soldier.ID)
	}
	if len(colony.Soldiers) != 1 {
		t.Errorf("Expected 1 soldier, got %d", len(colony.Soldiers))
	}
}

func TestSpawnSoliderWithID(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialNextID := colony.NextAntID

	soldier := SpawnSoldierWithID(colony, 999, 15, 15)

	if soldier.ID != 999 {
		t.Errorf("Expected ID 999, got %d", soldier.ID)
	}
	if colony.NextAntID != initialNextID {
		t.Error("NextAntID should not change with SpawnWorkerWithID")
	}
}

func TestSpawnNurse(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialID := colony.NextAntID

	nurse := SpawnNurse(colony, 11, 11)

	if nurse == nil {
		t.Fatal("SpawnNurse returned nil")
	}
	if nurse.ID != initialID {
		t.Errorf("Expected ID %d, got %d", initialID, nurse.ID)
	}
	if len(colony.Nurses) != 1 {
		t.Errorf("Expected 1 nurse, got %d", len(colony.Nurses))
	}
}
func TestNurseWithID(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialNextID := colony.NextAntID

	nurse := SpawnNurseWithID(colony, 999, 15, 15)

	if nurse.ID != 999 {
		t.Errorf("Expected ID 999, got %d", nurse.ID)
	}
	if colony.NextAntID != initialNextID {
		t.Error("NextAntID should not change with SpawnNurseWithID")
	}
}

func TestSpawnLarvae(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)
	initialID := colony.NextAntID

	larvae := SpawnLarvae(colony, 10, 11)

	if larvae == nil {
		t.Fatal("SpawnLarvae returned nil")
	}
	if larvae.ID != initialID {
		t.Errorf("Expected ID %d, got %d", initialID, larvae.ID)
	}
	if len(colony.Larvae) != 1 {
		t.Errorf("Expected 1 larvae, got %d", len(colony.Larvae))
	}
}

func TestRemoveLarvae(t *testing.T) {
	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)

	larvae1 := SpawnLarvae(colony, 10, 11)
	larvae2 := SpawnLarvae(colony, 10, 12)
	SpawnLarvae(colony, 10, 13)

	if len(colony.Larvae) != 3 {
		t.Fatalf("Expected 3 larvae, got %d", len(colony.Larvae))
	}

	RemoveLarvae(colony, larvae2)

	if len(colony.Larvae) != 2 {
		t.Errorf("Expected 2 larvae after removal, got %d", len(colony.Larvae))
	}

	// Verify correct one removed
	for _, l := range colony.Larvae {
		if l.ID == larvae2.ID {
			t.Error("larvae2 should have been removed")
		}
	}

	// Verify larvae1 still exists
	found := false
	for _, l := range colony.Larvae {
		if l.ID == larvae1.ID {
			found = true
		}
	}
	if !found {
		t.Error("larvae1 should still exist")
	}
}
