package logic

import (
	"antfarm/pathfinder"
	"antfarm/types"
	"antfarm/util"
	"fmt"
)

// ant.go - Defines how ants move and make decisions
// Each ant role has different behavior patterns

// Pathfinders - reusable instances
var (
	workerPathfinder = pathfinder.NewWorkerPathfinder()
	nursePathfinder  = pathfinder.NewNursePathfinder()
)

// updateWorker performs one tick of behavior for a worker ant
func updateWorker(world *types.World, colony *types.Colony, worker *types.WorkerAnt) {
	worker.Age++
	workerBehavior(world, colony, worker)
}

// updateSoldier performs one tick of behavior for a soldier ant
func updateSoldier(_ *types.World, soldier *types.SoldierAnt) {
	soldier.Age++
	// TODO: Implement soldier patrol/combat behavior
	soldier.CurrentAction = "patrolling"
}

// updateNurse performs one tick of behavior for the nurse ant
func updateNurse(world *types.World, colony *types.Colony, nurse *types.NurseAnt) {
	nurse.Age++
	nurseBehavior(world, colony, nurse)
}

// workerBehavior defines how worker ants act
// Workers dig tunnels, gather food, and bring it back to the queen
func workerBehavior(world *types.World, colony *types.Colony, worker *types.WorkerAnt) {
	// If carrying food, bring it back to queen
	if worker.CarryingFood {
		// Check if adjacent to queen
		xDist := util.Abs(colony.QueenPosition.X - worker.Position.X)
		yDist := util.Abs(colony.QueenPosition.Y - worker.Position.Y)

		if xDist <= 1 && yDist <= 1 {
			// Deposit food
			colony.Food += worker.FoodAmount
			worker.CarryingFood = false
			worker.FoodAmount = 0
			worker.CurrentAction = "deposited food"
			return
		}

		// Move toward queen using dedicated function
		worker.CurrentAction = fmt.Sprintf("bringing %d food to queen", worker.FoodAmount)
		if !workerPathfinder.BringFoodToQueen(world, colony, worker) {
			worker.CurrentAction = "stuck with food"
		}
		return
	}

	// Check current cell for food
	currentCell := world.GetCell(worker.Position.X, worker.Position.Y)
	if currentCell != nil && currentCell.Food > 0 {
		worker.CarryingFood = true
		worker.FoodAmount = 10
		currentCell.Food = 0
		worker.CurrentAction = "picked up food"
		return
	}

	// Check if on surface (grass) - grass gives 5 food and disappears
	if worker.Position.Y == 1 && currentCell != nil && currentCell.Soil == types.Empty {
		// Check if this cell has grass (not already harvested)
		// We'll use Food = -1 to mark harvested grass
		if currentCell.Food >= 0 {
			worker.CarryingFood = true
			worker.FoodAmount = 5
			currentCell.Food = -1 // Mark as harvested (no more grass)
			worker.CurrentAction = "foraged grass"
			return
		}
	}

	// Move randomly like a real ant (continues in same direction for several moves)
	if workerPathfinder.MoveRandomly(world, worker) {
		worker.CurrentAction = "exploring"
	} else {
		worker.CurrentAction = "resting"
	}
}

// nurseBehavior defines how nurse ants act
// Nurses guard the nursery and tend to larvae
func nurseBehavior(world *types.World, colony *types.Colony, nurse *types.NurseAnt) {
	// FIRST: If already nursing a larvae, stick with it until it becomes a worker
	if nurse.CurrentlyNursing != nil {
		// Check if this larvae still exists
		for _, larvae := range colony.Larvae {
			if larvae.ID == nurse.CurrentlyNursing.ID {
				// Check if adjacent
				if nursePathfinder.IsAdjacentToLarvae(nurse, larvae.Position) {
					// Adjacent - keep nursing
					nurse.CurrentAction = fmt.Sprintf("taking care of larvae #%d", larvae.ID)
					return
				} else {
					// Not adjacent anymore - larvae might have moved or we got pushed
					// Clear and re-acquire
					nurse.CurrentlyNursing = nil
					break
				}
			}
		}
		// Larvae no longer exists (became a worker), clear it
		nurse.CurrentlyNursing = nil
	}

	// SECOND: Find nearest larvae that needs care
	var targetLarvae *types.LarvaeAnt
	minDist := 9999

	for _, larvae := range colony.Larvae {
		if !larvae.HasNurseCare {
			// Make sure no other nurse is already nursing this larvae
			alreadyTargeted := false
			if colony.HeadNurse != nil && colony.HeadNurse != nurse &&
				colony.HeadNurse.CurrentlyNursing != nil &&
				colony.HeadNurse.CurrentlyNursing.ID == larvae.ID {
				alreadyTargeted = true
			}
			if !alreadyTargeted {
				for _, otherNurse := range colony.Nurses {
					if otherNurse != nurse && otherNurse.CurrentlyNursing != nil &&
						otherNurse.CurrentlyNursing.ID == larvae.ID {
						alreadyTargeted = true
						break
					}
				}
			}

			if !alreadyTargeted {
				dist := util.Abs(larvae.Position.X-nurse.Position.X) + util.Abs(larvae.Position.Y-nurse.Position.Y)
				if dist < minDist {
					minDist = dist
					targetLarvae = larvae
				}
			}
		}
	}

	// No larvae need care - guard the nursery
	if targetLarvae == nil {
		if nursePathfinder.GuardNursery(world, colony, nurse) {
			nurse.CurrentAction = "guarding nursery"
		} else {
			nurse.CurrentAction = "moving to nursery"
		}
		nurse.CurrentlyNursing = nil
		return
	}

	// Check if adjacent to target larvae
	if nursePathfinder.IsAdjacentToLarvae(nurse, targetLarvae.Position) {
		// Adjacent - start nursing
		targetLarvae.HasNurseCare = true
		targetLarvae.GrowthProgress = 100
		nurse.CurrentAction = fmt.Sprintf("taking care of larvae #%d", targetLarvae.ID)
		nurse.CurrentlyNursing = targetLarvae
		return
	}

	// Not adjacent - move toward larvae (going around queen if needed)
	nurse.CurrentAction = fmt.Sprintf("going to larvae #%d", targetLarvae.ID)
	nursePathfinder.MoveTowardLarvae(world, colony, nurse, targetLarvae.Position)
}
