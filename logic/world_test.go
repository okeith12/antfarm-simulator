package logic

import (
	"antfarm/types"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestUpdateWorld(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	AddColony(world, colony)

	UpdateWorld(world)

	if world.Ticks != 1 {
		t.Errorf("Expected Ticks 1, got %d", world.Ticks)
	}
}

func TestUpdateWorldMultipleTicks(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	AddColony(world, colony)

	for i := 0; i < 10; i++ {
		UpdateWorld(world)
	}

	if world.Ticks != 10 {
		t.Errorf("Expected Ticks 10, got %d", world.Ticks)
	}
}

func TestQueenLaysEggs(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	colony.Food = 200
	AddColony(world, colony)

	// Run until egg laying tick // TODO make it a parameter or part of config
	for i := 0; i < 50; i++ {
		UpdateWorld(world)
	}

	if colony.Eggs == 0 && colony.Queen.TotalEggsLaid == 0 {
		t.Error("Queen should have laid eggs by tick 50")
	}
}

func TestQueenDoesNotLayEggsWithoutFood(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	colony.Food = 5 // Not enough
	AddColony(world, colony)

	for i := 0; i < 50; i++ {
		UpdateWorld(world)
	}

	if colony.Queen.TotalEggsLaid > 0 {
		t.Error("Queen should not lay eggs without enough food")
	}
}

func TestEggsHatchIntoLarvae(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	colony.Eggs = 1
	AddColony(world, colony)

	// Run until hatch tick (30)
	for i := 0; i < 30; i++ {
		UpdateWorld(world)
	}

	if len(colony.Larvae) == 0 {
		t.Error("Egg should have hatched into larvae")
	}
}

func TestLarvaeBecomesWorker(t *testing.T) {
	world := types.NewWorld(40, 30)
	colony := types.NewColony("Red", 20, 15, tcell.ColorRed)
	AddColony(world, colony)

	// Spawn larvae with nurse care
	larvae := SpawnLarvae(colony, 21, 15)
	larvae.HasNurseCare = true
	larvae.Age = 49
	PlaceAnt(world, larvae)

	initialWorkers := len(colony.Workers)

	UpdateWorld(world)

	if len(colony.Workers) != initialWorkers+1 {
		t.Error("Larvae should have become a worker")
	}
}

func TestUpdateWorldEmptyWorld(t *testing.T) {
	world := types.NewWorld(40, 30)

	// Should not panic
	UpdateWorld(world)

	if world.Ticks != 1 {
		t.Errorf("Expected Ticks 1, got %d", world.Ticks)
	}
}
