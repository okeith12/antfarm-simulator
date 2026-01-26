package logic

import (
	"antfarm/types"
)

// colony_ant.go - Colony's ant spawning logic
// Handles ant lifecycle management within a colony

// SpawnWorker creates a new worker ant at the given position
func SpawnWorker(colony *types.Colony, x, y int) *types.WorkerAnt {
	worker := types.NewWorker(colony.NextAntID, x, y, colony.Name)
	colony.NextAntID++
	colony.Workers = append(colony.Workers, worker)
	return worker
}

// SpawnWorkerWithID creates a new worker ant with a specific ID (used when larvae becomes worker)
func SpawnWorkerWithID(colony *types.Colony, id int, x, y int) *types.WorkerAnt {
	worker := types.NewWorker(id, x, y, colony.Name)
	colony.Workers = append(colony.Workers, worker)
	return worker
}

// SpawnSoldier creates a new soldier ant at the given position
func SpawnSoldier(colony *types.Colony, x, y int) *types.SoldierAnt {
	soldier := types.NewSoldier(colony.NextAntID, x, y, colony.Name)
	colony.NextAntID++
	colony.Soldiers = append(colony.Soldiers, soldier)
	return soldier
}

// SpawnNurse creates a new nurse ant at the given position
func SpawnNurse(colony *types.Colony, x, y int) *types.NurseAnt {
	nurse := types.NewNurse(colony.NextAntID, x, y, colony.Name)
	colony.NextAntID++
	colony.Nurses = append(colony.Nurses, nurse)
	return nurse
}

// SpawnLarvae creates a new larvae at the given position
func SpawnLarvae(colony *types.Colony, x, y int) *types.LarvaeAnt {
	larvae := types.NewLarvae(colony.NextAntID, x, y, colony.Name)
	colony.NextAntID++
	colony.Larvae = append(colony.Larvae, larvae)
	return larvae
}

// RemoveLarvae removes a larvae from the colony's larvae list
func RemoveLarvae(colony *types.Colony, larvae *types.LarvaeAnt) {
	for i, l := range colony.Larvae {
		if l.ID == larvae.ID {
			colony.Larvae = append(colony.Larvae[:i], colony.Larvae[i+1:]...)
			return
		}
	}
}
