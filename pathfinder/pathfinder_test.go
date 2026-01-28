package pathfinder

import (
	"antfarm/types"
	"testing"
)

func TestDirectionToOffset(t *testing.T) {
	tests := []struct {
		dir       Direction
		expectedX int
		expectedY int
	}{
		{DirIdle, 0, 0},
		{DirUp, 0, -1},
		{DirDown, 0, 1},
		{DirLeft, -1, 0},
		{DirRight, 1, 0},
		{DirUpLeft, -1, -1},
		{DirUpRight, 1, -1},
		{DirDownLeft, -1, 1},
		{DirDownRight, 1, 1},
	}

	for _, tt := range tests {
		x, y := DirectionToOffset(tt.dir)
		if x != tt.expectedX || y != tt.expectedY {
			t.Errorf("DirectionToOffset(%d) = (%d,%d), want (%d,%d)", tt.dir, x, y, tt.expectedX, tt.expectedY)
		}
	}
}

func TestGetCardinalDirections(t *testing.T) {
	dirs := GetCardinalDirections()
	if len(dirs) != 4 {
		t.Errorf("Expected 4 cardinal directions, got %d", len(dirs))
	}
}

func TestGetAllDirections(t *testing.T) {
	dirs := GetAllDirections()
	if len(dirs) != 8 {
		t.Errorf("Expected 8 directions, got %d", len(dirs))
	}
}

func TestIsAdjacent(t *testing.T) {
	tests := []struct {
		pos1     types.Position
		pos2     types.Position
		expected bool
	}{
		{types.Position{X: 5, Y: 5}, types.Position{X: 6, Y: 5}, true},
		{types.Position{X: 5, Y: 5}, types.Position{X: 5, Y: 6}, true},
		{types.Position{X: 5, Y: 5}, types.Position{X: 6, Y: 6}, true},
		{types.Position{X: 5, Y: 5}, types.Position{X: 5, Y: 5}, false},
		{types.Position{X: 5, Y: 5}, types.Position{X: 7, Y: 5}, false},
	}

	for _, tt := range tests {
		result := IsAdjacent(tt.pos1, tt.pos2)
		if result != tt.expected {
			t.Errorf("IsAdjacent(%v, %v) = %v, want %v", tt.pos1, tt.pos2, result, tt.expected)
		}
	}
}

func TestIsAdjacentOrSame(t *testing.T) {
	tests := []struct {
		pos1     types.Position
		pos2     types.Position
		expected bool
	}{
		{types.Position{X: 5, Y: 5}, types.Position{X: 5, Y: 5}, true},
		{types.Position{X: 5, Y: 5}, types.Position{X: 6, Y: 5}, true},
		{types.Position{X: 5, Y: 5}, types.Position{X: 7, Y: 5}, false},
	}

	for _, tt := range tests {
		result := IsAdjacentOrSame(tt.pos1, tt.pos2)
		if result != tt.expected {
			t.Errorf("IsAdjacentOrSame(%v, %v) = %v, want %v", tt.pos1, tt.pos2, result, tt.expected)
		}
	}
}

func TestManhattanDistance(t *testing.T) {
	tests := []struct {
		pos1     types.Position
		pos2     types.Position
		expected int
	}{
		{types.Position{X: 0, Y: 0}, types.Position{X: 0, Y: 0}, 0},
		{types.Position{X: 0, Y: 0}, types.Position{X: 5, Y: 0}, 5},
		{types.Position{X: 0, Y: 0}, types.Position{X: 3, Y: 4}, 7},
	}

	for _, tt := range tests {
		result := ManhattanDistance(tt.pos1, tt.pos2)
		if result != tt.expected {
			t.Errorf("ManhattanDistance(%v, %v) = %d, want %d", tt.pos1, tt.pos2, result, tt.expected)
		}
	}
}

func TestCanMoveTo(t *testing.T) {
	world := types.NewWorld(20, 20)
	world.Grid[5][5].IsTunnel = true

	if !CanMoveTo(world, 5, 5) {
		t.Error("Should be able to move to empty tunnel")
	}
	if CanMoveTo(world, 10, 10) {
		t.Error("Should not be able to move to non-tunnel")
	}
	if CanMoveTo(world, -1, 5) {
		t.Error("Should not be able to move out of bounds")
	}
}

func TestCanDigTo(t *testing.T) {
	world := types.NewWorld(20, 20)
	world.Grid[10][10].Soil = types.Rock

	if !CanDigTo(world, 5, 5) {
		t.Error("Should be able to dig sand")
	}
	if CanDigTo(world, 10, 10) {
		t.Error("Should not be able to dig rock")
	}
	if CanDigTo(world, 5, 0) {
		t.Error("Should not be able to dig tunnel (surface)")
	}
}

func TestMove(t *testing.T) {
	world := types.NewWorld(20, 20)
	world.Grid[5][5].IsTunnel = true
	world.Grid[5][6].IsTunnel = true

	worker := types.NewWorker(1, 5, 5, "Red")
	world.Grid[5][5].Occupant = worker

	Move(world, worker, 6, 5)

	if worker.Position.X != 6 || worker.Position.Y != 5 {
		t.Errorf("Worker should be at (6,5), got (%d,%d)", worker.Position.X, worker.Position.Y)
	}
	if world.Grid[5][5].Occupant != nil {
		t.Error("Old cell should be empty")
	}
}

func TestDigAndMove(t *testing.T) {
	world := types.NewWorld(20, 20)
	world.Grid[5][5].IsTunnel = true

	worker := types.NewWorker(1, 5, 5, "Red")
	world.Grid[5][5].Occupant = worker

	dm := DigAndMove(world, worker, 6, 5)

	if !dm {
		t.Error("DigAndMove should happen on a diggable cell")
	}
	if !world.Grid[5][6].IsTunnel {
		t.Error("Cell should now be a tunnel")
	}
	if worker.Position.X != 6 {
		t.Error("Worker should have moved")
	}
}

func TestDigAndMoveFail(t *testing.T) {
	world := types.NewWorld(20, 20)
	world.Grid[5][5].IsTunnel = true
	world.Grid[5][6].Soil = types.Rock

	worker := types.NewWorker(1, 5, 5, "Red")

	success := DigAndMove(world, worker, 6, 5)

	if success {
		t.Error("DigAndMove should fail on rock")
	}
}
