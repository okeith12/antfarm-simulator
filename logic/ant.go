package logic

import (
	"antfarm/types"
	"antfarm/util"
	"fmt"
	"math/rand"
)

// ant.go - Defines how ants move and make decisions
// Each ant role has different behavior patterns

// updateWorker performs one tick of behavior for a worker ant
func updateWorker(world *types.World, colony *types.Colony, worker *types.WorkerAnt) {
	worker.Age++
	workerBehavior(world, colony, worker)
}

// updateSoldier performs one tick of behavior for a soldier ant
func updateSoldier(world *types.World, soldier *types.SoldierAnt) {
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

		// Move toward queen
		worker.CurrentAction = fmt.Sprintf("bringing %d food to queen", worker.FoodAmount)
		moveTowardQueen(world, colony, worker)
		return
	}

	// Check current cell for food
	currentCell := world.GetCell(worker.Position.X, worker.Position.Y)
	if currentCell != nil && currentCell.Food > 0 {
		// Pick up food (10 for food item)
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

	// Look for food or explore randomly
	directions := [][2]int{
		{0, -1}, // Up (toward surface for food)
		{0, 1},  // Down
		{-1, 0}, // Left
		{1, 0},  // Right
	}

	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]
		cell := world.GetCell(newX, newY)

		if cell == nil {
			continue
		}

		// If cell has food, go there
		if cell.IsTunnel && cell.Occupant == nil && cell.Food > 0 {
			worker.CurrentAction = "going to food"
			moveAnt(world, worker, newX, newY)
			return
		}

		// If it's soil, dig it
		if !cell.IsTunnel && cell.Soil != types.Rock {
			cell.IsTunnel = true
			worker.CurrentAction = "digging tunnel"
			moveAnt(world, worker, newX, newY)
			return
		}

		// If it's already a tunnel and empty, move there
		if cell.IsTunnel && cell.Occupant == nil {
			worker.CurrentAction = "exploring tunnel"
			moveAnt(world, worker, newX, newY)
			return
		}
	}

	worker.CurrentAction = "resting"
}

// moveTowardWorker moves a worker one step closer to a target position
func moveTowardWorker(world *types.World, worker *types.WorkerAnt, target types.Position) {
	dx := 0
	dy := 0

	if target.X > worker.Position.X {
		dx = 1
	} else if target.X < worker.Position.X {
		dx = -1
	}

	if target.Y > worker.Position.Y {
		dy = 1
	} else if target.Y < worker.Position.Y {
		dy = -1
	}

	attempts := [][2]int{{dx, dy}, {dx, 0}, {0, dy}, {0, -dy}, {-dx, 0}}

	for _, dir := range attempts {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && cell.IsTunnel && cell.Occupant == nil {
			moveAnt(world, worker, newX, newY)
			return
		}
	}
}

// nurseBehavior defines how nurse ants act
// Nurses seek out larvae and tend to them
func nurseBehavior(world *types.World, colony *types.Colony, nurse *types.NurseAnt) {
	// FIRST: If already nursing a larvae, stick with it until it's ready to become a worker
	if nurse.CurrentlyNursing != nil {
		// Check if this larvae still exists and still needs care
		for _, larvae := range colony.Larvae {
			if larvae.ID == nurse.CurrentlyNursing.ID {
				// Check if adjacent
				xDist := util.Abs(larvae.Position.X - nurse.Position.X)
				yDist := util.Abs(larvae.Position.Y - nurse.Position.Y)

				if xDist <= 1 && yDist <= 1 {
					// Adjacent - keep nursing
					nurse.CurrentAction = fmt.Sprintf("taking care of larvae #%d", larvae.ID)
					return
				} else {
					// Not adjacent - need to move closer
					nurse.CurrentAction = fmt.Sprintf("going to larvae #%d", larvae.ID)
					moveTowardWithDig(world, nurse, larvae.Position)
					return
				}
			}
		}
		// Larvae no longer exists (became a worker), clear it
		nurse.CurrentlyNursing = nil
	}

	// SECOND: Find nearest larvae that needs care (hasn't been nursed yet)
	var targetLarvae *types.LarvaeAnt
	minDist := 9999

	for _, larvae := range colony.Larvae {
		if !larvae.HasNurseCare {
			// Make sure no other nurse is already targeting this larvae
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

	if targetLarvae == nil {
		// No larvae need tending, wander near queen
		nurse.CurrentAction = "patrolling nursery"
		nurse.CurrentlyNursing = nil
		wanderNearQueen(world, colony, nurse)
		return
	}

	// Adjacent to larvae (including diagonals)? Start tending to it
	xDist := util.Abs(targetLarvae.Position.X - nurse.Position.X)
	yDist := util.Abs(targetLarvae.Position.Y - nurse.Position.Y)
	if xDist <= 1 && yDist <= 1 {
		targetLarvae.HasNurseCare = true
		targetLarvae.GrowthProgress = 100
		nurse.CurrentAction = fmt.Sprintf("taking care of larvae #%d", targetLarvae.ID)
		nurse.CurrentlyNursing = targetLarvae
		return
	}

	// Move toward larvae
	nurse.CurrentAction = fmt.Sprintf("going to larvae #%d", targetLarvae.ID)
	nurse.CurrentlyNursing = targetLarvae
	moveTowardWithDig(world, nurse, targetLarvae.Position)
}

// wanderNearQueen makes the nurse move randomly but stay close to the queen
func wanderNearQueen(world *types.World, colony *types.Colony, nurse *types.NurseAnt) {
	// Check if already adjacent to queen - if so, just stay put
	distToQueen := util.Abs(nurse.Position.X-colony.QueenPosition.X) + util.Abs(nurse.Position.Y-colony.QueenPosition.Y)
	if distToQueen <= 2 {
		// Already close to queen, no need to move
		nurse.CurrentAction = "guarding nursery"
		return
	}

	// Not close to queen, move closer
	directions := [][2]int{{0, 1}, {-1, 0}, {1, 0}, {0, -1}, {1, 1}, {-1, 1}, {1, -1}, {-1, -1}}
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	for _, dir := range directions {
		newX := nurse.Position.X + dir[0]
		newY := nurse.Position.Y + dir[1]

		// Only move if it gets us closer to queen
		newDistToQueen := util.Abs(newX-colony.QueenPosition.X) + util.Abs(newY-colony.QueenPosition.Y)
		if newDistToQueen >= distToQueen {
			continue
		}

		cell := world.GetCell(newX, newY)
		if cell != nil && cell.IsTunnel && cell.Occupant == nil {
			moveAnt(world, nurse, newX, newY)
			return
		}

		// Dig if needed
		if cell != nil && !cell.IsTunnel && cell.Soil != types.Rock {
			cell.IsTunnel = true
			moveAnt(world, nurse, newX, newY)
			return
		}
	}
}

// moveToward moves an ant one step closer to a target position
func moveToward(world *types.World, ant types.AntInterface, target types.Position) {
	baseAnt := ant.GetAnt()

	dx := 0
	dy := 0

	if target.X > baseAnt.Position.X {
		dx = 1
	} else if target.X < baseAnt.Position.X {
		dx = -1
	}

	if target.Y > baseAnt.Position.Y {
		dy = 1
	} else if target.Y < baseAnt.Position.Y {
		dy = -1
	}

	// Try horizontal first, then vertical, then diagonal
	attempts := [][2]int{{dx, 0}, {0, dy}, {dx, dy}}
	for _, dir := range attempts {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := baseAnt.Position.X + dir[0]
		newY := baseAnt.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && cell.IsTunnel && cell.Occupant == nil {
			moveAnt(world, ant, newX, newY)
			return
		}
	}
}

// moveAnt relocates an ant from its current position to a new position
func moveAnt(world *types.World, ant types.AntInterface, newX, newY int) {
	baseAnt := ant.GetAnt()

	// Clear old position
	oldCell := world.GetCell(baseAnt.Position.X, baseAnt.Position.Y)
	if oldCell != nil {
		oldCell.Occupant = nil
	}

	// Move to new position
	baseAnt.Position.X = newX
	baseAnt.Position.Y = newY

	newCell := world.GetCell(newX, newY)
	if newCell != nil {
		newCell.Occupant = ant
	}
}

// moveTowardWithDig moves nurse toward target, trying multiple paths
func moveTowardWithDig(world *types.World, nurse *types.NurseAnt, target types.Position) {
	dx := sign(target.X - nurse.Position.X)
	dy := sign(target.Y - nurse.Position.Y)

	// All 8 directions, prioritized by direction to target
	allDirections := [][2]int{
		{dx, dy},   // Direct
		{dx, 0},    // Horizontal toward
		{0, dy},    // Vertical toward
		{dx, -dy},  // Diagonal alternate
		{-dx, dy},  // Diagonal alternate
		{0, -dy},   // Vertical away
		{-dx, 0},   // Horizontal away
		{-dx, -dy}, // Opposite diagonal
	}

	// First: try moving through empty tunnels
	for _, dir := range allDirections {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := nurse.Position.X + dir[0]
		newY := nurse.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && cell.IsTunnel && cell.Occupant == nil {
			moveAnt(world, nurse, newX, newY)
			return
		}
	}

	// Second: dig new tunnels
	for _, dir := range allDirections {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := nurse.Position.X + dir[0]
		newY := nurse.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && !cell.IsTunnel && cell.Soil != types.Rock {
			cell.IsTunnel = true
			moveAnt(world, nurse, newX, newY)
			return
		}
	}
	// If completely stuck, do nothing this tick
}

// moveTowardQueen moves worker toward queen, trying multiple paths
func moveTowardQueen(world *types.World, colony *types.Colony, worker *types.WorkerAnt) {
	target := colony.QueenPosition
	dx := sign(target.X - worker.Position.X)
	dy := sign(target.Y - worker.Position.Y)

	allDirections := [][2]int{
		{dx, dy},
		{dx, 0},
		{0, dy},
		{dx, -dy},
		{-dx, dy},
		{0, -dy},
		{-dx, 0},
		{-dx, -dy},
	}

	// First: try moving through empty tunnels
	for _, dir := range allDirections {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && cell.IsTunnel && cell.Occupant == nil {
			moveAnt(world, worker, newX, newY)
			return
		}
	}

	// Second: dig new tunnels
	for _, dir := range allDirections {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]

		cell := world.GetCell(newX, newY)
		if cell != nil && !cell.IsTunnel && cell.Soil != types.Rock {
			cell.IsTunnel = true
			moveAnt(world, worker, newX, newY)
			return
		}
	}
}

// sign returns -1, 0, or 1 based on the value
func sign(x int) int {
	if x > 0 {
		return 1
	} else if x < 0 {
		return -1
	}
	return 0
}
