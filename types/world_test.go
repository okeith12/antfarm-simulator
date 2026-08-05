package types

import (
	"antfarm/random"
	"testing"
)

func TestNewWorld(t *testing.T) {
	world := NewWorld(80, 40, random.New(1))

	if world.Width != 80 {
		t.Errorf("Expected Width 80, got %d", world.Width)
	}
	if world.Height != 40 {
		t.Errorf("Expected Height 40, got %d", world.Height)
	}
	if world.Ticks != 0 {
		t.Errorf("Expected Ticks 0, got %d", world.Ticks)
	}
	if len(world.Cells) != 80*40 {
		t.Errorf("Expected %d cells, got %d", 80*40, len(world.Cells))
	}
}

func TestWorldIndexIsRowMajor(t *testing.T) {
	world := NewWorld(80, 40, random.New(1))

	// Row-major: stepping x moves one slot, stepping y moves a whole row.
	if got := world.Index(0, 0); got != 0 {
		t.Errorf("Index(0,0) = %d, want 0", got)
	}
	if got := world.Index(1, 0); got != 1 {
		t.Errorf("Index(1,0) = %d, want 1", got)
	}
	if got := world.Index(0, 1); got != 80 {
		t.Errorf("Index(0,1) = %d, want 80", got)
	}
	if got := world.Index(79, 39); got != 80*40-1 {
		t.Errorf("Index(79,39) = %d, want %d", got, 80*40-1)
	}

	// Every position must map to a distinct slot. A formula like x*y would
	// collide here: (2,3) and (3,2) and (6,1) would all land on 6.
	seen := make(map[int]bool, world.Width*world.Height)
	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			i := world.Index(x, y)
			if seen[i] {
				t.Fatalf("Index collision at (%d,%d) -> %d", x, y, i)
			}
			seen[i] = true
		}
	}
	if len(seen) != len(world.Cells) {
		t.Errorf("covered %d slots, expected %d", len(seen), len(world.Cells))
	}
}

func TestGetCellAliasesTheSameStorage(t *testing.T) {
	world := NewWorld(80, 40, random.New(1))

	// GetCell hands back a pointer into the flat slice, not a copy, so writes
	// through it must be visible to the next reader.
	world.GetCell(5, 7).IsTunnel = true
	if !world.Cells[world.Index(5, 7)].IsTunnel {
		t.Error("write through GetCell did not reach the backing slice")
	}
	if !world.GetCell(5, 7).IsTunnel {
		t.Error("second GetCell did not observe the write")
	}
}

func TestWorldIsValidPosition(t *testing.T) {
	world := NewWorld(80, 40, random.New(1))

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
	world := NewWorld(80, 40, random.New(1))

	cell := world.GetCell(10, 5)
	if cell == nil {
		t.Error("Expected cell at (10,5)")
	}

	cell = world.GetCell(-1, 5)
	if cell != nil {
		t.Error("Expected nil for invalid position")
	}
}
