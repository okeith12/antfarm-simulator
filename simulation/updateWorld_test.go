package logic

import (
	"antfarm/random"
	"antfarm/types"
	"testing"
)

func TestUpdateWorld(t *testing.T) {
	world := types.NewWorld(40, 30, random.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	UpdateWorld(world)

	if world.Ticks != 1 {
		t.Errorf("Expected Ticks 1, got %d", world.Ticks)
	}
}

func TestUpdateWorldMultipleTicks(t *testing.T) {
	world := types.NewWorld(40, 30, random.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	for i := 0; i < 10; i++ {
		UpdateWorld(world)
	}

	if world.Ticks != 10 {
		t.Errorf("Expected Ticks 10, got %d", world.Ticks)
	}
}

func TestQueenLaysEggs(t *testing.T) {
	world := types.NewWorld(40, 30, random.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
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
	world := types.NewWorld(40, 30, random.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
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
	world := types.NewWorld(40, 30, random.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
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
	// Maturing rolls for a role: 20% nurse, 15% soldier, 65% worker. Seed 2 rolls
	// a worker, so this asserts the exact outcome instead of a 65% coin flip.
	world := types.NewWorld(40, 30, random.New(2))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
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
	world := types.NewWorld(40, 30, random.New(1))

	// Should not panic
	UpdateWorld(world)

	if world.Ticks != 1 {
		t.Errorf("Expected Ticks 1, got %d", world.Ticks)
	}
}

// TestSameSeedProducesSameColony is the property the Go-to-C++ parity gate
// depends on: one seed, many ticks, identical colony state every time.
func TestSameSeedProducesSameColony(t *testing.T) {
	run := func(seed uint32) (int, int, int, int, int) {
		world := types.NewWorld(60, 30, random.New(seed))
		colony := types.NewColony("Red", 15, 10, types.ColonyRed)
		AddColony(world, colony)
		for i := 0; i < 500; i++ {
			UpdateWorld(world)
		}
		return len(colony.Workers), len(colony.Nurses), len(colony.Soldiers),
			len(colony.Larvae), colony.Food
	}

	w1, n1, s1, l1, f1 := run(777)
	w2, n2, s2, l2, f2 := run(777)

	if w1 != w2 || n1 != n2 || s1 != s2 || l1 != l2 || f1 != f2 {
		t.Errorf("same seed diverged after 500 ticks:\n  run 1: workers=%d nurses=%d soldiers=%d larvae=%d food=%d\n  run 2: workers=%d nurses=%d soldiers=%d larvae=%d food=%d",
			w1, n1, s1, l1, f1, w2, n2, s2, l2, f2)
	}
}
