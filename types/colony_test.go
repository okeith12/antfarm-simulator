package types

import (
	"testing"
)

func TestNewColony(t *testing.T) {
	colony := NewColony("Red", 10, 20, ColonyRed)

	if colony.Name != "Red" {
		t.Errorf("Expected Name 'Red', got '%s'", colony.Name)
	}
	if colony.Color != ColonyRed {
		t.Errorf("Expected Color Red, got %v", colony.Color)
	}
	if colony.Queen == nil {
		t.Error("Expected Queen to exist")
	}
	if colony.HeadNurse == nil {
		t.Error("Expected HeadNurse to exist")
	}
	// Every colony opens with the same trio: queen, nurse, worker.
	if len(colony.Workers) != 1 {
		t.Errorf("Expected 1 founding worker, got %d", len(colony.Workers))
	}
	if colony.Food != 50*FoodScale {
		t.Errorf("Expected Food %d, got %d", 50*FoodScale, colony.Food)
	}
	if colony.NextAntID != 3 {
		t.Errorf("Expected NextAntID 3, got %d", colony.NextAntID)
	}
}

func TestGetAllAnts(t *testing.T) {
	colony := NewColony("Red", 10, 20, ColonyRed)

	// Initial: queen + head nurse + founding worker
	if len(colony.GetAllAnts()) != 3 {
		t.Errorf("Expected 3 ants, got %d", len(colony.GetAllAnts()))
	}

	// Add another worker
	colony.Workers = append(colony.Workers, NewWorker(3, 5, 5, "Red"))
	if len(colony.GetAllAnts()) != 4 {
		t.Errorf("Expected 4 ants, got %d", len(colony.GetAllAnts()))
	}
}

func TestGetAntCount(t *testing.T) {
	colony := NewColony("Red", 10, 20, ColonyRed)

	if colony.GetAntCount() != 3 {
		t.Errorf("Expected count 3, got %d", colony.GetAntCount())
	}

	colony.Workers = append(colony.Workers, NewWorker(3, 5, 5, "Red"))
	colony.Soldiers = append(colony.Soldiers, NewSoldier(4, 6, 6, "Red"))

	if colony.GetAntCount() != 5 {
		t.Errorf("Expected count 5, got %d", colony.GetAntCount())
	}
}
