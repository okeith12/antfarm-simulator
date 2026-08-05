package logic

import (
	"antfarm/types"
)

// spawn.go -  ant spawning logic
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

// SpawnQueenWithID creates a new queen with a specific ID (used when a larvae
// matures into a queen). If the colony has no reigning queen she takes the
// throne immediately; otherwise she joins the spares as an heir.
func SpawnQueenWithID(colony *types.Colony, id int, x, y int) *types.QueenAnt {
	queen := types.NewQueen(id, x, y, colony.Name)
	if colony.Queen == nil {
		colony.Queen = queen
		colony.QueenPosition = types.Position{X: x, Y: y}
	} else {
		colony.Queens = append(colony.Queens, queen)
	}
	return queen
}

// RemoveQueen removes a spare queen from the colony's heirs
func RemoveQueen(colony *types.Colony, queen *types.QueenAnt) {
	for i, q := range colony.Queens {
		if q.ID == queen.ID {
			colony.Queens = append(colony.Queens[:i], colony.Queens[i+1:]...)
			return
		}
	}
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
