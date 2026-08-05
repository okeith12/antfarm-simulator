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
// throne immediately; otherwise she joins the heirs in waiting.
//
// A queen is immortal until she bears an heir. That birth starts her slow
// decline, and from then on she loses health until an heir takes over.
func SpawnQueenWithID(colony *types.Colony, id int, x, y int) *types.QueenAnt {
	queen := types.NewQueen(id, x, y, colony.Name)
	if colony.Queen == nil {
		colony.Queen = queen
		colony.QueenPosition = types.Position{X: x, Y: y}
		return queen
	}

	colony.Queens = append(colony.Queens, queen)
	colony.Queen.Declining = true
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

// demoteHeir turns a queen who was passed over for the throne into an ordinary
// worker or nurse, keeping her ID and position. A colony holds exactly one
// queen, so an uncrowned heir has no role left to play and joins the workforce.
func demoteHeir(world *types.World, colony *types.Colony, heir *types.QueenAnt) {
	x, y := heir.Position.X, heir.Position.Y
	RemoveAnt(world, heir)

	var replacement types.AntInterface
	if world.Rng.Below(100) < nurseSpawnChance {
		replacement = SpawnNurseWithID(colony, heir.ID, x, y)
	} else {
		replacement = SpawnWorkerWithID(colony, heir.ID, x, y)
	}
	replacement.GetAnt().CurrentAction = "gave up the claim"
	PlaceAnt(world, replacement)
}
