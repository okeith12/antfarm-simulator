package logic

import (
	"antfarm/types"
	"antfarm/util"

	"math/rand"
)

// world.go - Main World simulation update logic
// Handles per-tick updates for the entire world including all colonies

// UpdateWorld advances the simulation by one tick
// Updates all colonies (ants, queens, eggs) and world physics
func UpdateWorld(world *types.World) {
	world.Ticks++

	// Update each colony's ants and resources
	for _, colony := range world.Colonies {
		updateColony(world, colony)
	}
}

// Timing constants
var (
	eggLayingInterval = 50 // Queen lays eggs every 50 ticks
	eggHatchTime      = 30 // Eggs become larvae after 30 ticks
	larvaeGrowTime    = 50 // Larvae become workers after 50 ticks with nurse care
)

// updateColony handles all updates for a single colony
func updateColony(world *types.World, colony *types.Colony) {
	// Set default queen action
	if colony.Queen != nil {
		colony.Queen.CurrentAction = "resting"
	}

	// Queen lays 1-5 eggs periodically
	if world.Ticks > 0 && world.Ticks%eggLayingInterval == 0 && colony.Food >= 10 {
		eggsToLay := rand.Intn(5) + 1 // Random 1-5 eggs

		// Only lay as many eggs as we can afford
		for i := 0; i < eggsToLay && colony.Food >= 10; i++ {
			colony.Eggs++
			colony.Food -= 10
			if colony.Queen != nil {
				colony.Queen.TotalEggsLaid++
			}
		}

		if colony.Queen != nil {
			colony.Queen.CurrentAction = "laying eggs"
		}
	}

	// Eggs hatch into larvae
	if world.Ticks > 0 && world.Ticks%eggHatchTime == 0 && colony.Eggs > 0 {
		colony.Eggs--

		// Find an empty cell near queen to spawn larvae
		spawnX, spawnY := findEmptySpawnPosition(world, colony.QueenPosition)
		if spawnX != -1 && spawnY != -1 {
			larvae := SpawnLarvae(colony, spawnX, spawnY)

			// Place larvae in world
			cell := world.GetCell(spawnX, spawnY)
			if cell != nil {
				cell.IsTunnel = true
				cell.Occupant = larvae
			}
		}
	}

	// Check larvae that have been nursed - they become workers
	for i := len(colony.Larvae) - 1; i >= 0; i-- {
		larvae := colony.Larvae[i]

		if larvae.HasNurseCare && larvae.Age >= larvaeGrowTime {
			// Remove larvae from world
			RemoveAnt(world, larvae)

			// Create new worker at same position, KEEPING THE SAME ID
			worker := SpawnWorkerWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
			worker.CurrentAction = "newly hatched"

			// Place worker in world
			PlaceAnt(world, worker)

			// Clear the nurse's CurrentlyNursing if it was this larvae
			if colony.HeadNurse != nil && colony.HeadNurse.CurrentlyNursing != nil &&
				colony.HeadNurse.CurrentlyNursing.ID == larvae.ID {
				colony.HeadNurse.LarvaeNursed++
				colony.HeadNurse.CurrentlyNursing = nil
			}
			for _, nurse := range colony.Nurses {
				if nurse.CurrentlyNursing != nil && nurse.CurrentlyNursing.ID == larvae.ID {
					nurse.LarvaeNursed++
					nurse.CurrentlyNursing = nil
				}
			}

			// Remove from larvae list
			RemoveLarvae(colony, larvae)
		}
	}

	// Update head nurse
	if colony.HeadNurse != nil {
		updateNurse(world, colony, colony.HeadNurse)
	}

	// Update other nurses
	for _, nurse := range colony.Nurses {
		updateNurse(world, colony, nurse)
	}

	// Update workers
	for _, worker := range colony.Workers {
		updateWorker(world, colony, worker)
	}

	// Update soldiers
	for _, soldier := range colony.Soldiers {
		updateSoldier(world, soldier)
	}

	// Age all larvae and set their action based on whether a nurse is actively caring for them
	for _, larvae := range colony.Larvae {
		larvae.Age++

		// Check if this specific larvae is being nursed (nurse is adjacent)
		isBeingNursed := false

		if colony.HeadNurse != nil && colony.HeadNurse.CurrentlyNursing != nil &&
			colony.HeadNurse.CurrentlyNursing.ID == larvae.ID {
			// Check nurse is actually adjacent
			xDist := util.Abs(colony.HeadNurse.Position.X - larvae.Position.X)
			yDist := util.Abs(colony.HeadNurse.Position.Y - larvae.Position.Y)
			if xDist <= 1 && yDist <= 1 {
				isBeingNursed = true
			}
		}

		// Also check all other nurses
		if !isBeingNursed {
			for _, nurse := range colony.Nurses {
				if nurse.CurrentlyNursing != nil && nurse.CurrentlyNursing.ID == larvae.ID {
					xDist := util.Abs(nurse.Position.X - larvae.Position.X)
					yDist := util.Abs(nurse.Position.Y - larvae.Position.Y)
					if xDist <= 1 && yDist <= 1 {
						isBeingNursed = true
						break
					}
				}
			}
		}

		if isBeingNursed {
			larvae.CurrentAction = "getting care"
		} else {
			larvae.CurrentAction = "waiting for care"
		}
	}
}

// findEmptySpawnPosition finds an empty tunnel cell near the queen to spawn larvae
func findEmptySpawnPosition(world *types.World, queenPos types.Position) (int, int) {
	// Check positions around queen in expanding rings
	offsets := [][2]int{
		{1, 0}, {0, 1}, {-1, 0}, {0, -1}, // Adjacent
		{1, 1}, {-1, 1}, {1, -1}, {-1, -1}, // Diagonal
		{2, 0}, {0, 2}, {-2, 0}, {0, -2}, // Two away
		{2, 1}, {1, 2}, {-1, 2}, {-2, 1}, // Further out
		{-2, -1}, {-1, -2}, {1, -2}, {2, -1},
	}

	for _, offset := range offsets {
		x := queenPos.X + offset[0]
		y := queenPos.Y + offset[1]

		cell := world.GetCell(x, y)
		if cell != nil && cell.Occupant == nil {
			return x, y
		}
	}

	return -1, -1 // No empty position found
}
