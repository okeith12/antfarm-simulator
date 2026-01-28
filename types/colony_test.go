package types

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewColony(t *testing.T) {
	colony := NewColony("Red", 10, 20, tcell.ColorRed)

	if colony.Name != "Red" {
		t.Errorf("Expected Name 'Red', got '%s'", colony.Name)
	}
	if colony.Color != tcell.ColorRed {
		t.Errorf("Expected Color Red, got %v", colony.Color)
	}
	if colony.Queen == nil {
		t.Error("Expected Queen to exist")
	}
	if colony.HeadNurse == nil {
		t.Error("Expected HeadNurse to exist")
	}
	if colony.Food != 50 {
		t.Errorf("Expected Food 50, got %d", colony.Food)
	}
	if colony.NextAntID != 2 {
		t.Errorf("Expected NextAntID 2, got %d", colony.NextAntID)
	}
}

func TestGetAllAnts(t *testing.T) {
	colony := NewColony("Red", 10, 20, tcell.ColorRed)

	// Initial: queen + head nurse
	if len(colony.GetAllAnts()) != 2 {
		t.Errorf("Expected 2 ants, got %d", len(colony.GetAllAnts()))
	}

	// Add a worker
	colony.Workers = append(colony.Workers, NewWorker(2, 5, 5, "Red"))
	if len(colony.GetAllAnts()) != 3 {
		t.Errorf("Expected 3 ants, got %d", len(colony.GetAllAnts()))
	}
}

func TestGetAntCount(t *testing.T) {
	colony := NewColony("Red", 10, 20, tcell.ColorRed)

	if colony.GetAntCount() != 2 {
		t.Errorf("Expected count 2, got %d", colony.GetAntCount())
	}

	colony.Workers = append(colony.Workers, NewWorker(2, 5, 5, "Red"))
	colony.Soldiers = append(colony.Soldiers, NewSoldier(3, 6, 6, "Red"))

	if colony.GetAntCount() != 4 {
		t.Errorf("Expected count 4, got %d", colony.GetAntCount())
	}
}
