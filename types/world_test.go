package types

import "testing"

func TestNewWorld(t *testing.T) {
	world := NewWorld(80, 40)

	if world.Width != 80 {
		t.Errorf("Expected Width 80, got %d", world.Width)
	}
	if world.Height != 40 {
		t.Errorf("Expected Height 40, got %d", world.Height)
	}
	if world.Ticks != 0 {
		t.Errorf("Expected Ticks 0, got %d", world.Ticks)
	}
	if len(world.Grid) != 40 {
		t.Errorf("Expected Grid height 40, got %d", len(world.Grid))
	}
	if len(world.Grid[0]) != 80 {
		t.Errorf("Expected Grid width 80, got %d", len(world.Grid[0]))
	}
}

func TestWorldIsValidPosition(t *testing.T) {
	world := NewWorld(80, 40)

	if !world.IsValidPosition(0, 0) {
		t.Error("(0,0) should be valid")
	}
	if !world.IsValidPosition(79, 39) {
		t.Error("(79,39) should be valid")
	}
	if world.IsValidPosition(-1, 0) {
		t.Error("(-1,0) should be invalid")
	}
	if world.IsValidPosition(80, 0) {
		t.Error("(80,0) should be invalid")
	}
	if world.IsValidPosition(0, 40) {
		t.Error("(0,40) should be invalid")
	}
}

func TestGetCell(t *testing.T) {
	world := NewWorld(80, 40)

	cell := world.GetCell(10, 5)
	if cell == nil {
		t.Error("Expected cell at (10,5)")
	}

	cell = world.GetCell(-1, 5)
	if cell != nil {
		t.Error("Expected nil for invalid position")
	}
}
