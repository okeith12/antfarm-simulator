package pathfinder

import (
	"antfarm/types"
	"antfarm/util"
)

// NursePathfinder handles movement logic for nurse ants
type NursePathfinder struct{}

// NewNursePathfinder creates a new nurse pathfinder
func NewNursePathfinder() *NursePathfinder {
	return &NursePathfinder{}
}

// GuardNursery keeps the nurse near the queen when no larvae exist
func (np *NursePathfinder) GuardNursery(world *types.World, colony *types.Colony, nurse *types.NurseAnt) bool {
	distToQueen := ManhattanDistance(nurse.Position, colony.QueenPosition)

	// If within 2 cells of queen, stay put
	if distToQueen <= 2 {
		return true
	}

	// Move closer to queen
	np.MoveTowardQueen(world, colony, nurse)
	return false
}

// MoveTowardQueen moves the nurse closer to the queen
func (np *NursePathfinder) MoveTowardQueen(world *types.World, colony *types.Colony, nurse *types.NurseAnt) bool {
	return np.MoveTowardTarget(world, colony, nurse, colony.QueenPosition, false)
}

// MoveTowardLarvae moves the nurse toward a larvae
func (np *NursePathfinder) MoveTowardLarvae(world *types.World, colony *types.Colony, nurse *types.NurseAnt, larvaePos types.Position) bool {
	return np.MoveTowardTarget(world, colony, nurse, larvaePos, true)
}

// MoveTowardTarget moves nurse toward a target, going around queen if needed
func (np *NursePathfinder) MoveTowardTarget(world *types.World, colony *types.Colony, nurse *types.NurseAnt, target types.Position, goAroundQueen bool) bool {
	// All 8 directions
	directions := [][2]int{
		{0, -1}, {0, 1}, {-1, 0}, {1, 0},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	// Find the best direction that gets us closer to target
	bestDir := -1
	bestDist := ManhattanDistance(nurse.Position, target)

	for i, dir := range directions {
		newX := nurse.Position.X + dir[0]
		newY := nurse.Position.Y + dir[1]

		// Skip queen's position
		if newX == colony.QueenPosition.X && newY == colony.QueenPosition.Y {
			continue
		}

		// Check if we can move or dig there
		if !CanMoveTo(world, newX, newY) && !CanDigTo(world, newX, newY) {
			continue
		}

		// Calculate new distance to target
		newDist := util.Abs(newX-target.X) + util.Abs(newY-target.Y)

		// Only accept if it gets us closer
		if newDist < bestDist {
			bestDist = newDist
			bestDir = i
		}
	}

	// If we found a direction that gets us closer, take it
	if bestDir >= 0 {
		newX := nurse.Position.X + directions[bestDir][0]
		newY := nurse.Position.Y + directions[bestDir][1]

		if CanMoveTo(world, newX, newY) {
			MoveAnt(world, nurse, newX, newY)
			return true
		}
		if CanDigTo(world, newX, newY) {
			DigAndMove(world, nurse, newX, newY)
			return true
		}
	}

	// If no direction gets us closer, we need to go around
	// Try any direction that doesn't go backwards
	currentDist := ManhattanDistance(nurse.Position, target)

	for _, dir := range directions {
		newX := nurse.Position.X + dir[0]
		newY := nurse.Position.Y + dir[1]

		// Skip queen's position
		if newX == colony.QueenPosition.X && newY == colony.QueenPosition.Y {
			continue
		}

		// Accept same distance (going around)
		newDist := util.Abs(newX-target.X) + util.Abs(newY-target.Y)
		if newDist > currentDist+1 {
			continue // Don't go too far backwards
		}

		if CanMoveTo(world, newX, newY) {
			MoveAnt(world, nurse, newX, newY)
			return true
		}
		if CanDigTo(world, newX, newY) {
			DigAndMove(world, nurse, newX, newY)
			return true
		}
	}

	return false
}

// IsAdjacentToLarvae checks if nurse is adjacent to a larvae position
func (np *NursePathfinder) IsAdjacentToLarvae(nurse *types.NurseAnt, larvaePos types.Position) bool {
	xDist := util.Abs(nurse.Position.X - larvaePos.X)
	yDist := util.Abs(nurse.Position.Y - larvaePos.Y)
	return xDist <= 1 && yDist <= 1
}
