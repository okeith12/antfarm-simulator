package pathfinder

import (
	"antfarm/types"
	"math/rand"
)

// WorkerPathfinder handles movement logic for worker ants
type WorkerPathfinder struct{}

// NewWorkerPathfinder creates a new worker pathfinder
func NewWorkerPathfinder() *WorkerPathfinder {
	return &WorkerPathfinder{}
}

// MoveRandomly makes the worker move like a real ant - continues in same direction
// for several steps before changing direction
func (wp *WorkerPathfinder) MoveRandomly(world *types.World, worker *types.WorkerAnt) bool {
	// If worker has a current direction and momentum, keep going that way
	if worker.MovesMade < worker.MovesInDirection && worker.CurrentDirection != int(DirIdle) {
		dx, dy := DirectionToOffset(Direction(worker.CurrentDirection))
		newX := worker.Position.X + dx
		newY := worker.Position.Y + dy

		// Try to continue in current direction
		if CanMoveTo(world, newX, newY) {
			MoveAnt(world, worker, newX, newY)
			worker.MovesMade++
			return true
		}

		// Try to dig in current direction
		if CanDigTo(world, newX, newY) {
			DigAndMove(world, worker, newX, newY)
			worker.MovesMade++
			return true
		}

		// Can't continue - pick new direction
		worker.MovesMade = worker.MovesInDirection // Force new direction
	}

	// Pick a new direction and momentum
	return wp.pickNewDirection(world, worker)
}

// pickNewDirection selects a new random direction for the worker
func (wp *WorkerPathfinder) pickNewDirection(world *types.World, worker *types.WorkerAnt) bool {
	// Get all cardinal directions
	directions := GetCardinalDirections()

	// Shuffle for randomness
	rand.Shuffle(len(directions), func(i, j int) {
		directions[i], directions[j] = directions[j], directions[i]
	})

	// Avoid going back the way we came if possible
	oppositeDir := wp.getOppositeDirection(Direction(worker.CurrentDirection))

	// Try each direction (preferring not to go backwards)
	for _, dir := range directions {
		if dir == oppositeDir {
			continue // Skip opposite direction first pass
		}

		dx, dy := DirectionToOffset(dir)
		newX := worker.Position.X + dx
		newY := worker.Position.Y + dy

		if CanMoveTo(world, newX, newY) || CanDigTo(world, newX, newY) {
			// Set new direction with random momentum (3-6 moves)
			worker.CurrentDirection = int(dir)
			worker.MovesInDirection = rand.Intn(4) + 3
			worker.MovesMade = 0

			if CanMoveTo(world, newX, newY) {
				MoveAnt(world, worker, newX, newY)
			} else {
				DigAndMove(world, worker, newX, newY)
			}
			worker.MovesMade++
			return true
		}
	}

	// If all else fails, try going backwards
	if oppositeDir != DirIdle {
		dx, dy := DirectionToOffset(oppositeDir)
		newX := worker.Position.X + dx
		newY := worker.Position.Y + dy

		if CanMoveTo(world, newX, newY) {
			worker.CurrentDirection = int(oppositeDir)
			worker.MovesInDirection = rand.Intn(4) + 3
			worker.MovesMade = 0
			MoveAnt(world, worker, newX, newY)
			worker.MovesMade++
			return true
		}
	}

	return false // Completely stuck
}

// getOppositeDirection returns the opposite of a direction
func (wp *WorkerPathfinder) getOppositeDirection(dir Direction) Direction {
	switch dir {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	default:
		return DirIdle
	}
}

// BringFoodToQueen moves the worker toward the queen to deposit food
func (wp *WorkerPathfinder) BringFoodToQueen(world *types.World, colony *types.Colony, worker *types.WorkerAnt) bool {
	target := colony.QueenPosition

	// Current position
	curX := worker.Position.X
	curY := worker.Position.Y

	// Build directions based on where we need to go
	var directions [][2]int

	// Add directions that move us closer first
	if curX < target.X {
		directions = append(directions, [2]int{1, 0}) // Right
	} else if curX > target.X {
		directions = append(directions, [2]int{-1, 0}) // Left
	}

	if curY < target.Y {
		directions = append(directions, [2]int{0, 1}) // Down
	} else if curY > target.Y {
		directions = append(directions, [2]int{0, -1}) // Up
	}

	// todo: comment out
	// // Add diagonal if both X and Y need to change
	// if curX != target.X && curY != target.Y {
	// 	dx := 1
	// 	if curX > target.X {
	// 		dx = -1
	// 	}
	// 	dy := 1
	// 	if curY > target.Y {
	// 		dy = -1
	// 	}
	// 	// Insert diagonal at the front
	// 	directions = append([][2]int{{dx, dy}}, directions...)
	// }

	// Add perpendicular directions for going around obstacles
	if curX == target.X {
		// On same X, add left/right for going around
		directions = append(directions, [2]int{1, 0}, [2]int{-1, 0})
	}
	if curY == target.Y {
		// On same Y, add up/down for going around
		directions = append(directions, [2]int{0, 1}, [2]int{0, -1})
	}

	// Try each direction
	for _, dir := range directions {
		newX := curX + dir[0]
		newY := curY + dir[1]

		// Don't step on queen
		if newX == target.X && newY == target.Y {
			continue
		}

		// Try to move
		if CanMoveTo(world, newX, newY) {
			MoveAnt(world, worker, newX, newY)
			return true
		}

		// Try to dig
		if CanDigTo(world, newX, newY) {
			DigAndMove(world, worker, newX, newY)
			return true
		}
	}

	return false
}

// MoveTowardTarget moves the worker toward a specific target
func (wp *WorkerPathfinder) MoveTowardTarget(world *types.World, worker *types.WorkerAnt, target types.Position) bool {
	// Calculate direction to target
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

	// Try directions prioritized toward target
	attempts := [][2]int{
		{dx, dy},  // Diagonal toward
		{dx, 0},   // Horizontal toward
		{0, dy},   // Vertical toward
		{dx, -dy}, // Alternate diagonal
		{-dx, dy}, // Alternate diagonal
		{0, -dy},  // Vertical away
		{-dx, 0},  // Horizontal away
	}

	// First pass: try empty tunnels
	for _, dir := range attempts {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]

		if CanMoveTo(world, newX, newY) {
			MoveAnt(world, worker, newX, newY)
			return true
		}
	}

	// Second pass: dig if needed
	for _, dir := range attempts {
		if dir[0] == 0 && dir[1] == 0 {
			continue
		}
		newX := worker.Position.X + dir[0]
		newY := worker.Position.Y + dir[1]

		if CanDigTo(world, newX, newY) {
			DigAndMove(world, worker, newX, newY)
			return true
		}
	}

	return false
}

// IsAdjacentToTarget checks if worker is next to target
func (wp *WorkerPathfinder) IsAdjacentToTarget(worker *types.WorkerAnt, target types.Position) bool {
	return IsAdjacentOrSame(worker.Position, target)
}
