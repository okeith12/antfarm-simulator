package logic

import "antfarm/types"

//matureLarvaeToAnt is a helper function to decide the next step in the lifecylce for the larvae

// Role spawn chances (out of 100)
// Adjust these values to change the distribution of ant types
const (
	nurseSpawnChance   = 20 // 20% chance to become nurse (0-19)
	soldierSpawnChance = 15 // 15% chance to become soldier (20-34)
	// Remaining 65% become workers (35-99)
)

// larvaeToAnt creates the appropriate adult ant based on a random roll
// roll should be 0-99, larvae provides ID and position
func matureLarvaeToAnt(colony *types.Colony, larvae *types.LarvaeAnt, roll int) types.AntInterface {
	var newAnt types.AntInterface

	switch {
	case roll < nurseSpawnChance:
		// Become a nurse
		nurse := SpawnNurseWithID(colony, larvae.ID, larvae.Position.X, larvae.Position.Y)
		nurse.CurrentAction = "newly hatched"
		newAnt = nurse

	case roll < nurseSpawnChance+soldierSpawnChance:
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
