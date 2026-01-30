package types

import (
	"testing"
)

func TestNewAnt(t *testing.T) {
	ant := NewAnt(1, Worker, 10, 20, "TestColony", 100, 500)

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
	if ant.MaxHealth != 100 {
		t.Errorf("Expected MaxHealth 100, got %d", ant.MaxHealth)
	}
	if ant.ColonyID != "TestColony" {
		t.Errorf("Expected ColonyID 'TestColony', got '%s'", ant.ColonyID)
	}
	if ant.Age != 0 {
		t.Errorf("Expected Age 0, got %d", ant.Age)
	}
	if ant.MaxAge != 500 {
		t.Errorf("Expected MaxAge 500, got %d", ant.MaxAge)
	}
	if ant.CurrentAction != "idle" {
		t.Errorf("Expected CurrentAction 'idle', got '%s'", ant.CurrentAction)
	}
}

func TestAntIsDead(t *testing.T) {
	tests := []struct {
		name     string
		health   int
		age      int
		maxAge   int
		expected bool
	}{
		{"Healthy young ant", 100, 0, 500, false},
		{"Healthy old ant at limit", 100, 500, 500, true},
		{"Healthy ant past limit", 100, 600, 500, true},
		{"Dead ant (0 health)", 0, 0, 500, true},
		{"Dead ant (negative health)", -10, 0, 500, true},
		{"Dying ant (1 health)", 1, 0, 500, false},
		{"Old ant (age = maxAge - 1)", 100, 499, 500, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ant := NewAnt(1, Worker, 0, 0, "test", 100, tt.maxAge)
			ant.Health = tt.health
			ant.Age = tt.age

			if ant.IsDead() != tt.expected {
				t.Errorf("IsDead() = %v, expected %v", ant.IsDead(), tt.expected)
			}
		})
	}
}

func TestLifespanConstants(t *testing.T) {
	// Verify queens live longest
	if QueenMaxTick <= WorkerMaxTick {
		t.Error("Queens should live longer than workers")
	}
	if QueenMaxTick <= NurseMaxTick {
		t.Error("Queens should live longer than nurses")
	}
	if QueenMaxTick <= SoldierMaxTick {
		t.Error("Queens should live longer than soldiers")
	}

	// Verify larvae have shortest lifespan (they must mature or die)
	if LarvaeMaxTick >= WorkerMaxTick {
		t.Error("Larvae should have shorter max age than workers")
	}

	// Verify soldiers are more durable than workers
	if SoldierMaxHealth <= WorkerMaxHealth {
		t.Error("Soldiers should have more health than workers")
	}

	// Verify queen has most health
	if QueenMaxHealth <= SoldierMaxHealth {
		t.Error("Queen should have more health than soldiers")
	}
}
