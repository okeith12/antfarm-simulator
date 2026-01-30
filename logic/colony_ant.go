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

// SpawnSoldierWithID creates a new soldier ant with a specific ID (used when larvae becomes soldier)
func SpawnSoldierWithID(colony *types.Colony, id int, x, y int) *types.SoldierAnt {
	soldier := types.NewSoldier(id, x, y, colony.Name)
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

// SpawnNurseWithID creates a new nurse ant with a specific ID (used when larvae becomes nurse)
func SpawnNurseWithID(colony *types.Colony, id int, x, y int) *types.NurseAnt {
	nurse := types.NewNurse(id, x, y, colony.Name)
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

// RemoveWorker removes a worker from the colony's worker list
func RemoveWorker(colony *types.Colony, worker *types.WorkerAnt) {
	for i, w := range colony.Workers {
		if w.ID == worker.ID {
			colony.Workers = append(colony.Workers[:i], colony.Workers[i+1:]...)
			return
		}
	}
}

// RemoveSoldier removes a soldier from the colony's soldier list
func RemoveSoldier(colony *types.Colony, soldier *types.SoldierAnt) {
	for i, s := range colony.Soldiers {
		if s.ID == soldier.ID {
			colony.Soldiers = append(colony.Soldiers[:i], colony.Soldiers[i+1:]...)
			return
		}
	}
}

// RemoveNurse removes a nurse from the colony's nurse list
func RemoveNurse(colony *types.Colony, nurse *types.NurseAnt) {
	for i, n := range colony.Nurses {
		if n.ID == nurse.ID {
			colony.Nurses = append(colony.Nurses[:i], colony.Nurses[i+1:]...)
			return
		}
	}
}
