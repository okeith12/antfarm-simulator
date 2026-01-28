package logic

import (
	"antfarm/types"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestAddColony(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)

	AddColony(world, colony)

	if len(world.Colonies) != 1 {
		t.Errorf("Expected 1 colony, got %d", len(world.Colonies))
	}

	// Queen should be placed
	queenCell := world.GetCell(20, 15)
	if queenCell.Occupant == nil {
		t.Error("Queen should be placed in world")
	}
	if !queenCell.IsTunnel {
		t.Error("Queen cell should be a tunnel")
	}

	// Head nurse should be placed
	nurseCell := world.GetCell(21, 15)
	if nurseCell.Occupant == nil {
		t.Error("Head nurse should be placed in world")
	}
}

func TestAddMultipleColonies(t *testing.T) {
	world := types.NewWorld(80, 40)
	colony1 := types.NewColony("Red", 20, 20, tcell.ColorRed)
	colony2 := types.NewColony("Blue", 60, 20, tcell.ColorBlue)

	AddColony(world, colony1)
	AddColony(world, colony2)

	if len(world.Colonies) != 2 {
		t.Errorf("Expected 2 colonies, got %d", len(world.Colonies))
	}
}

func TestPlaceAnt(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true

	worker := types.NewWorker(1, 10, 10, "Red")
	success := PlaceAnt(world, worker)

	if !success {
		t.Error("PlaceAnt should succeed for empty tunnel")
	}
	if world.Grid[10][10].Occupant != worker {
		t.Error("Worker should be in cell")
	}
}

func TestPlaceAntFailsOnOccupied(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true

	worker1 := types.NewWorker(1, 10, 10, "Red")
	worker2 := types.NewWorker(2, 10, 10, "Red")

	PlaceAnt(world, worker1)
	success := PlaceAnt(world, worker2)

	if success {
		t.Error("PlaceAnt should fail on occupied cell")
	}
}

func TestPlaceAntFailsOnNonTunnel(t *testing.T) {
	world := types.NewWorld(40, 30)

	worker := types.NewWorker(1, 10, 10, "Red")
	success := PlaceAnt(world, worker)

	if success {
		t.Error("PlaceAnt should fail on non-tunnel")
	}
}

func TestRemoveAnt(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true

	worker := types.NewWorker(1, 10, 10, "Red")
	PlaceAnt(world, worker)

	RemoveAnt(world, worker)

	if world.Grid[10][10].Occupant != nil {
		t.Error("Cell should be empty after RemoveAnt")
	}
}

func TestMoveAnt(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true
	world.Grid[10][11].IsTunnel = true

	worker := types.NewWorker(1, 10, 10, "Red")
	PlaceAnt(world, worker)

	success := MoveAnt(world, worker, 11, 10)

	if !success {
		t.Error("MoveAnt should succeed")
	}
	if worker.Position.X != 11 || worker.Position.Y != 10 {
		t.Errorf("Worker should be at (11,10), got (%d,%d)", worker.Position.X, worker.Position.Y)
	}
	if world.Grid[10][10].Occupant != nil {
		t.Error("Old cell should be empty")
	}
	if world.Grid[10][11].Occupant != worker {
		t.Error("New cell should have worker")
	}
}

func TestMoveAntFailsOnInvalid(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true

	worker := types.NewWorker(1, 10, 10, "Red")
	PlaceAnt(world, worker)

	success := MoveAnt(world, worker, -1, 10)

	if success {
		t.Error("MoveAnt should fail on invalid position")
	}
}

func TestMoveAntFailsOnOccupied(t *testing.T) {
	world := types.NewWorld(40, 30)
	world.Grid[10][10].IsTunnel = true
	world.Grid[10][11].IsTunnel = true

	worker1 := types.NewWorker(1, 10, 10, "Red")
	worker2 := types.NewWorker(2, 11, 10, "Red")
	PlaceAnt(world, worker1)
	PlaceAnt(world, worker2)

	success := MoveAnt(world, worker1, 11, 10)

	if success {
		t.Error("MoveAnt should fail on occupied cell")
	}
}
