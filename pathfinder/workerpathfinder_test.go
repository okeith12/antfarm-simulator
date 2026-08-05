package pathfinder

import (
	"antfarm/rng"
	"antfarm/types"
	"testing"
)

func TestNewWorkerPathfinder(t *testing.T) {
	wp := NewWorkerPathfinder()
	if wp == nil {
		t.Error("NewWorkerPathfinder should not return nil")
	}
}

func TestWorkerMoveRandomly(t *testing.T) {
	world := types.NewWorld(20, 20, rng.New(1))
	wp := NewWorkerPathfinder()

	worker := types.NewWorker(1, 10, 1, "Red")
	world.GetCell(10, 1).Occupant = worker

	moved := wp.MoveRandomly(world, worker)
	if !moved {
		t.Error("Worker should be able to move on surface")
	}
}

func TestWorkerBringFoodToQueen(t *testing.T) {
	world := types.NewWorld(20, 20, rng.New(1))
	wp := NewWorkerPathfinder()

	colony := types.NewColony("Red", 10, 10, types.ColonyRed)

	// Create tunnel path
	for x := 8; x <= 12; x++ {
		world.GetCell(x, 10).IsTunnel = true
	}
	world.GetCell(10, 10).Occupant = colony.Queen

	worker := types.NewWorker(2, 8, 10, "Red")
	worker.CarryingFood = true
	world.GetCell(8, 10).Occupant = worker

	success := wp.BringFoodToQueen(world, colony, worker)

	if !success {
		t.Error("Worker should be able to move toward queen")
	}
	if worker.Position.X <= 8 {
		t.Error("Worker should have moved closer to queen")
	}
}

func TestWorkerIsAdjacentToTarget(t *testing.T) {
	wp := NewWorkerPathfinder()

	worker := types.NewWorker(1, 10, 10, "Red")

	if !wp.IsAdjacentToTarget(worker, types.Position{X: 11, Y: 10}) {
		t.Error("Worker should be adjacent to (11,10)")
	}
	if wp.IsAdjacentToTarget(worker, types.Position{X: 15, Y: 10}) {
		t.Error("Worker should not be adjacent to (15,10)")
	}
}

func TestWorkerMoveTowardTarget(t *testing.T) {
	world := types.NewWorld(20, 20, rng.New(1))
	wp := NewWorkerPathfinder()

	// Create tunnel path
	for x := 5; x <= 15; x++ {
		world.GetCell(x, 10).IsTunnel = true
	}

	worker := types.NewWorker(1, 5, 10, "Red")
	world.GetCell(5, 10).Occupant = worker

	target := types.Position{X: 15, Y: 10}
	success := wp.MoveTowardTarget(world, worker, target)

	if !success {
		t.Error("Worker should be able to move toward target")
	}
	if worker.Position.X <= 5 {
		t.Error("Worker should have moved right toward target")
	}
}
