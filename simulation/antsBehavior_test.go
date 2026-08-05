package logic

import (
	"antfarm/rng"
	"antfarm/types"
	"testing"
)

func TestUpdateWorkerAges(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	worker := SpawnWorker(colony, 20, 14)
	world.GetCell(20, 14).IsTunnel = true
	PlaceAnt(world, worker)

	initialAge := worker.Age

	UpdateWorld(world)

	if worker.Age != initialAge+1 {
		t.Errorf("Worker age should increment, got %d", worker.Age)
	}
}

func TestUpdateSoldierAges(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	soldier := SpawnSoldier(colony, 20, 14)
	world.GetCell(20, 14).IsTunnel = true
	PlaceAnt(world, soldier)

	initialAge := soldier.Age

	UpdateWorld(world)

	if soldier.Age != initialAge+1 {
		t.Errorf("Soldier age should increment, got %d", soldier.Age)
	}
}

func TestUpdateNurseAges(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	initialAge := colony.HeadNurse.Age

	UpdateWorld(world)

	if colony.HeadNurse.Age != initialAge+1 {
		t.Errorf("Nurse age should increment, got %d", colony.HeadNurse.Age)
	}
}

func TestWorkerPicksUpFood(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	worker := SpawnWorker(colony, 10, 1)
	world.GetCell(10, 1).Food = 5
	PlaceAnt(world, worker)

	UpdateWorld(world)

	if !worker.CarryingFood {
		t.Error("Worker should pick up food")
	}
}

func TestWorkerDepositsFood(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	// Place worker adjacent to queen with food
	worker := SpawnWorker(colony, 21, 15)
	worker.CarryingFood = true
	worker.FoodAmount = 10
	world.GetCell(21, 15).IsTunnel = true
	PlaceAnt(world, worker)

	initialFood := colony.Food

	UpdateWorld(world)

	if worker.CarryingFood {
		t.Error("Worker should have deposited food")
	}
	if colony.Food != initialFood+10 {
		t.Errorf("Colony food should increase by 10, got %d", colony.Food)
	}
}

func TestLarvaeAges(t *testing.T) {
	world := types.NewWorld(40, 30, rng.New(1))
	colony := types.NewColony("Red", 20, 15, types.ColonyRed)
	AddColony(world, colony)

	larvae := SpawnLarvae(colony, 21, 15)
	world.GetCell(21, 15).IsTunnel = true
	PlaceAnt(world, larvae)

	initialAge := larvae.Age

	UpdateWorld(world)

	if larvae.Age != initialAge+1 {
		t.Errorf("Larvae age should increment, got %d", larvae.Age)
	}
}
