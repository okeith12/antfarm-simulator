package types

import "testing"

func TestNewAnt(t *testing.T) {
	ant := NewAnt(1, Worker, 10, 20, "TestColony")

	if ant.ID != 1 {
		t.Errorf("Expected ID 1, got %d", ant.ID)
	}
	if ant.Role != Worker {
		t.Errorf("Expected Role Worker, got %d", ant.Role)
	}
	if ant.Position.X != 10 || ant.Position.Y != 20 {
		t.Errorf("Expected Position (10, 20), got (%d, %d)", ant.Position.X, ant.Position.Y)
	}
	if ant.Health != 100 {
		t.Errorf("Expected Health 100, got %d", ant.Health)
	}
	if ant.ColonyID != "TestColony" {
		t.Errorf("Expected ColonyID 'TestColony', got '%s'", ant.ColonyID)
	}
	if ant.Age != 0 {
		t.Errorf("Expected Age 0, got %d", ant.Age)
	}
	if ant.CurrentAction != "idle" {
		t.Errorf("Expected CurrentAction 'idle', got '%s'", ant.CurrentAction)
	}
}
