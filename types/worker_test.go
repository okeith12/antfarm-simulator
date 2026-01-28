package types

import "testing"

func TestNewWorker(t *testing.T) {
	worker := NewWorker(5, 10, 20, "Red")

	if worker.ID != 5 {
		t.Errorf("Expected ID 5, got %d", worker.ID)
	}
	if worker.Role != Worker {
		t.Errorf("Expected Role Worker, got %d", worker.Role)
	}
	if worker.Position.X != 10 || worker.Position.Y != 20 {
		t.Errorf("Expected Position (10,20), got (%d,%d)", worker.Position.X, worker.Position.Y)
	}
	if worker.Health != 100 {
		t.Errorf("Expected Health 100, got %d", worker.Health)
	}
	if worker.CarryingFood {
		t.Error("Expected CarryingFood false")
	}
	if worker.DiggingPower != 1 {
		t.Errorf("Expected DiggingPower 1, got %d", worker.DiggingPower)
	}
}

func TestWorkerGetAnt(t *testing.T) {
	worker := NewWorker(1, 0, 0, "Red")
	if worker.GetAnt() == nil {
		t.Error("GetAnt() should not return nil")
	}
}

func TestWorkerGetAntIcon(t *testing.T) {
	worker := NewWorker(1, 0, 0, "Red")
	if worker.GetIcon() != '●' {
		t.Errorf("Expected icon '●', got '%c'", worker.GetIcon())
	}
}

func TestWorkerGetRole(t *testing.T) {
	worker := NewWorker(1, 0, 0, "Red")
	if worker.GetRole() != Worker {
		t.Errorf("Expected Worker role, got %d", worker.GetRole())
	}
}
