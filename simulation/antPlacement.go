package logic

import (
	"antfarm/types"
)

// antPlacement.go - World's Colony mutation logic
// Handles colony placement and ant positioning within the world

// AddColony places a new colony in the world
// Digs an initial chamber and places ALL ants (queen, nurses, etc.) in the world
func AddColony(world *types.World, colony *types.Colony) {
	world.Colonies = append(world.Colonies, colony)

	// Place queen in world
	if colony.Queen != nil {
		pos := colony.Queen.Position
		if world.IsValidPosition(pos.X, pos.Y) {
			cell := world.Grid[pos.Y][pos.X]
			cell.IsTunnel = true
			cell.Occupant = colony.Queen
		}
	}

	// Place head nurse in world
	if colony.HeadNurse != nil {
		pos := colony.HeadNurse.Position
		if world.IsValidPosition(pos.X, pos.Y) {
			cell := world.Grid[pos.Y][pos.X]
			cell.IsTunnel = true
			cell.Occupant = colony.HeadNurse
		}
	}

	// Place any other nurses
	for _, nurse := range colony.Nurses {
		pos := nurse.Position
		if world.IsValidPosition(pos.X, pos.Y) {
			cell := world.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = nurse
			}
		}
	}

	// Place any existing workers
	for _, worker := range colony.Workers {
		pos := worker.Position
		if world.IsValidPosition(pos.X, pos.Y) {
			cell := world.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = worker
			}
		}
	}

	// Place any existing soldiers
	for _, soldier := range colony.Soldiers {
		pos := soldier.Position
		if world.IsValidPosition(pos.X, pos.Y) {
			cell := world.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = soldier
			}
		}
	}
}

// PlaceAnt places any ant type into the world at its current position
func PlaceAnt(world *types.World, ant types.AntInterface) bool {
	pos := ant.GetAnt().Position
	if !world.IsValidPosition(pos.X, pos.Y) {
		return false
	}

	cell := world.Grid[pos.Y][pos.X]
	if cell.IsTunnel && cell.Occupant == nil {
		cell.Occupant = ant
		return true
	}
	return false
}

// RemoveAnt removes an ant from its current position in the world
func RemoveAnt(world *types.World, ant types.AntInterface) {
	pos := ant.GetAnt().Position
	if world.IsValidPosition(pos.X, pos.Y) {
		cell := world.Grid[pos.Y][pos.X]
		if cell.Occupant != nil && cell.Occupant.GetAnt().ID == ant.GetAnt().ID {
			cell.Occupant = nil
		}
	}
}

// MoveAnt moves an ant from its current position to a new position
func MoveAnt(world *types.World, ant types.AntInterface, newX, newY int) bool {
	if !world.IsValidPosition(newX, newY) {
		return false
	}

	newCell := world.GetCell(newX, newY)
	if newCell == nil || !newCell.IsTunnel || newCell.Occupant != nil {
		return false
	}

	// Remove from old position
	RemoveAnt(world, ant)

	// Update ant's position
	ant.GetAnt().Position.X = newX
	ant.GetAnt().Position.Y = newY

	// Place in new position
	newCell.Occupant = ant
	return true
}
