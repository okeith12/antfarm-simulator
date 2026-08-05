package logic

import "antfarm/types"

//matureLarvaeToAnt is a helper function to decide the next step in the lifecylce for the larvae

// Role spawn chances (out of 100)
// Adjust these values to change the distribution of ant types
const (
	queenSpawnChance   = 1  // 1% chance to become a queen (0)
	nurseSpawnChance   = 20 // 20% chance to become nurse (1-20)
	soldierSpawnChance = 15 // 15% chance to become soldier (21-35)
	// Remaining 64% become workers (36-99)
)

// larvaeToAnt creates the appropriate adult ant based on a random roll
// roll should be 0-99, larvae provides ID and position
func matureLarvaeToAnt(colony *types.Colony, larvae *types.LarvaeAnt, roll int) types.AntInterface {
	var newAnt types.AntInterface

	switch {
	case roll < queenSpawnChance:
		// Become a queen: takes the throne if vacant, otherwise an heir
		queen := SpawnQueenWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
		queen.CurrentAction = "newly hatched"
		newAnt = queen

	case roll < queenSpawnChance+nurseSpawnChance:
		// Become a nurse
		nurse := SpawnNurseWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
		nurse.CurrentAction = "newly hatched"
		newAnt = nurse

	case roll < queenSpawnChance+nurseSpawnChance+soldierSpawnChance:
		// Become a soldier
		soldier := SpawnSoldierWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
		soldier.CurrentAction = "newly hatched"
		newAnt = soldier

	default:
		// Become a worker (most common)
		worker := SpawnWorkerWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
		worker.CurrentAction = "newly hatched"
		newAnt = worker
	}

	return newAnt
}
