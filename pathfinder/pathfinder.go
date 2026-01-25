package pathfinder

import (
	"antfarm/types"
	"antfarm/util"
)

// Direction represents a movement direction
type Direction int

const (
	DirIdle Direction = iota
	DirUp
	DirDown
	DirLeft
	DirRight
	DirUpLeft
	DirUpRight
	DirDownLeft
	DirDownRight
)

// DirectionToOffset converts a Direction to x,y offset
func DirectionToOffset(dir Direction) (int, int) {
	switch dir {
	case DirUp:
		return 0, -1
	case DirDown:
		return 0, 1
	case DirLeft:
		return -1, 0
	case DirRight:
		return 1, 0
	case DirUpLeft:
		return -1, -1
	case DirUpRight:
		return 1, -1
	case DirDownLeft:
		return -1, 1
	case DirDownRight:
		return 1, 1
	default:
		return 0, 0
	}
}

// GetCardinalDirections returns the 4 main directions
func GetCardinalDirections() []Direction {
	return []Direction{DirUp, DirDown, DirLeft, DirRight}
}

// GetAllDirections returns all 8 directions
func GetAllDirections() []Direction {
	return []Direction{DirUp, DirDown, DirLeft, DirRight, DirUpLeft, DirUpRight, DirDownLeft, DirDownRight}
}

// IsAdjacent checks if two positions are adjacent (including diagonals)
func IsAdjacent(pos1, pos2 types.Position) bool {
	xDist := util.Abs(pos1.X - pos2.X)
	yDist := util.Abs(pos1.Y - pos2.Y)
	return xDist <= 1 && yDist <= 1 && !(xDist == 0 && yDist == 0)
}

// IsAdjacentOrSame checks if two positions are adjacent or the same
func IsAdjacentOrSame(pos1, pos2 types.Position) bool {
	xDist := util.Abs(pos1.X - pos2.X)
	yDist := util.Abs(pos1.Y - pos2.Y)
	return xDist <= 1 && yDist <= 1
}

// ManhattanDistance calculates the Manhattan distance between two positions
func ManhattanDistance(pos1, pos2 types.Position) int {
	return util.Abs(pos1.X-pos2.X) + util.Abs(pos1.Y-pos2.Y)
}

// CanMoveTo checks if a cell is valid to move into
func CanMoveTo(world *types.World, x, y int) bool {
	cell := world.GetCell(x, y)
	if cell == nil {
		return false
	}
	return cell.IsTunnel && cell.Occupant == nil
}

// CanDigTo checks if a cell can be dug into
func CanDigTo(world *types.World, x, y int) bool {
	cell := world.GetCell(x, y)
	if cell == nil {
		return false
	}
	return !cell.IsTunnel && cell.Soil != types.Rock
}

// MoveAnt relocates an ant from its current position to a new position
func MoveAnt(world *types.World, ant types.AntInterface, newX, newY int) {
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

// DigAndMove digs a tunnel and moves the ant there
func DigAndMove(world *types.World, ant types.AntInterface, newX, newY int) bool {
	cell := world.GetCell(newX, newY)
	if cell == nil {
		return false
	}
	if !cell.IsTunnel && cell.Soil != types.Rock {
		cell.IsTunnel = true
		MoveAnt(world, ant, newX, newY)
		return true
	}
	return false
}
