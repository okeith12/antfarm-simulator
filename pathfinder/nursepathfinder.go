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

// MoveTowardTarget moves nurse toward a target, can pass through queen's cell
func (np *NursePathfinder) MoveTowardTarget(world *types.World, colony *types.Colony, nurse *types.NurseAnt, target types.Position, goAroundQueen bool) bool {
	curX := nurse.Position.X
	curY := nurse.Position.Y

	// All 8 directions
	directions := [][2]int{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {1, -1}, {-1, 1}, {1, 1},
	}

	currentDist := ManhattanDistance(nurse.Position, target)

	// Find direction that gets us closer
	for _, dir := range directions {
		newX := curX + dir[0]
		newY := curY + dir[1]

		newDist := util.Abs(newX-target.X) + util.Abs(newY-target.Y)

		// Only move if it gets us closer
		if newDist >= currentDist {
			continue
		}

		// If it's the queen's cell, swap positions with queen
		if newX == colony.QueenPosition.X && newY == colony.QueenPosition.Y {
			// Swap: nurse goes to queen's spot, queen goes to nurse's spot
			queenCell := world.GetCell(colony.QueenPosition.X, colony.QueenPosition.Y)
			nurseCell := world.GetCell(curX, curY)

			// Move queen to nurse's old position
			colony.Queen.Position.X = curX
			colony.Queen.Position.Y = curY
			nurseCell.Occupant = colony.Queen

			// Update colony's queen position so everyant knows
			colony.QueenPosition.X = curX
			colony.QueenPosition.Y = curY

			// Move nurse to queen's position
			nurse.Position.X = newX
			nurse.Position.Y = newY
			queenCell.Occupant = nurse

			return true
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
