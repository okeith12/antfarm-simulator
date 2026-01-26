package types

import "math/rand"

// world.go - Defines the game world (the entire ant farm environment)
// The world is a 2D grid of cells containing terrain, tunnels, and ants

// World represents the entire  environment
// Contains the grid of cells, all colonies, and tracks simulation time
type World struct {
	Width    int       // Width of the world
	Height   int       // Height of the world
	Grid     [][]*Cell // 2D grid of cells
	Colonies []*Colony // All ant colonies in this world
	Ticks    int       // Number of updates that have occurred
}

// NewWorld creates a new world with procedurally generated terrain
// The top rows are open air (surface), deeper layers have different soil types
func NewWorld(width, height int) *World {
	grid := make([][]*Cell, height)

	for y := 0; y < height; y++ {
		grid[y] = make([]*Cell, width)
		for x := 0; x < width; x++ {
			// Generate terrain
			soilType := generateSoilType(y, height)
			grid[y][x] = NewCell(soilType)
		}
	}

	// Create surface (top 2 rows are empty)
	for y := 0; y < 2; y++ {
		for x := 0; x < width; x++ {
			grid[y][x].Soil = Empty
			grid[y][x].IsTunnel = true
		}
	}

	// Scatter food on surface (top row)
	for x := 0; x < width; x++ {
		if rand.Float64() < 0.1 { // 10% chance of food
			grid[1][x].Food = 5 // Food pellet
		}
	}

	return &World{
		Width:    width,
		Height:   height,
		Grid:     grid,
		Colonies: []*Colony{},
		Ticks:    0,
	}
}

// generateSoilType determines what soil type should appear at a given depth
// Surface is sand, middle layers are dirt/clay, deep layers are rock
func generateSoilType(y, maxHeight int) Soil {
	// depth := float64(y) / float64(maxHeight)

	// if depth < 0.1 {
	// 	return Sand
	// } else if depth < 0.6 {
	// 	if rand.Float64() < 0.1 {
	// 		return Clay
	// 	}
	// 	return Dirt
	// } else if depth < 0.9 {
	// 	if rand.Float64() < 0.3 {
	// 		return Rock
	// 	}
	// 	return Clay
	// } else {
	// 	return Rock
	// }
	if y < 2 {
		return Empty // Top 2 rows are surface/grass
	}
	return Sand // Everything else is sand
}

// AddColony places a new colony in the world
// Digs an initial chamber and places ALL ants (queen, nurses, etc.) in the world
func (w *World) AddColony(colony *Colony) {
	w.Colonies = append(w.Colonies, colony)

	// Place queen in world
	if colony.Queen != nil {
		pos := colony.Queen.Position
		if w.IsValidPosition(pos.X, pos.Y) {
			cell := w.Grid[pos.Y][pos.X]
			cell.IsTunnel = true
			cell.Occupant = colony.Queen
		}
	}

	// Place head nurse in world
	if colony.HeadNurse != nil {
		pos := colony.HeadNurse.Position
		if w.IsValidPosition(pos.X, pos.Y) {
			cell := w.Grid[pos.Y][pos.X]
			cell.IsTunnel = true
			cell.Occupant = colony.HeadNurse
		}
	}

	// Place any other nurses
	for _, nurse := range colony.Nurses {
		pos := nurse.Position
		if w.IsValidPosition(pos.X, pos.Y) {
			cell := w.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = nurse
			}
		}
	}

	// Place any existing workers
	for _, worker := range colony.Workers {
		pos := worker.Position
		if w.IsValidPosition(pos.X, pos.Y) {
			cell := w.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = worker
			}
		}
	}

	// Place any existing soldiers
	for _, soldier := range colony.Soldiers {
		pos := soldier.Position
		if w.IsValidPosition(pos.X, pos.Y) {
			cell := w.Grid[pos.Y][pos.X]
			if cell.IsTunnel && cell.Occupant == nil {
				cell.Occupant = soldier
			}
		}
	}
}

// IsValidPosition checks if the given coordinates are within world bounds
func (w *World) IsValidPosition(x, y int) bool {
	return x >= 0 && x < w.Width && y >= 0 && y < w.Height
}

// GetCell safely retrieves a cell at the given position
// Returns nil if the position is out of bounds
func (w *World) GetCell(x, y int) *Cell {
	if !w.IsValidPosition(x, y) {
		return nil
	}
	return w.Grid[y][x]
}

// PlaceAnt places any ant type into the world at its current position
func (w *World) PlaceAnt(ant AntInterface) bool {
	pos := ant.GetAnt().Position
	if !w.IsValidPosition(pos.X, pos.Y) {
		return false
	}

	cell := w.Grid[pos.Y][pos.X]
	if cell.IsTunnel && cell.Occupant == nil {
		cell.Occupant = ant
		return true
	}
	return false
}

// RemoveAnt removes an ant from its current position in the world
func (w *World) RemoveAnt(ant AntInterface) {
	pos := ant.GetAnt().Position
	if w.IsValidPosition(pos.X, pos.Y) {
		cell := w.Grid[pos.Y][pos.X]
		if cell.Occupant != nil && cell.Occupant.GetAnt().ID == ant.GetAnt().ID {
			cell.Occupant = nil
		}
	}
}

// MoveAnt moves an ant from its current position to a new position
func (w *World) MoveAnt(ant AntInterface, newX, newY int) bool {
	if !w.IsValidPosition(newX, newY) {
		return false
	}

	newCell := w.GetCell(newX, newY)
	if newCell == nil || !newCell.IsTunnel || newCell.Occupant != nil {
		return false
	}

	// Remove from old position
	w.RemoveAnt(ant)

	// Update ant's position
	ant.GetAnt().Position.X = newX
	ant.GetAnt().Position.Y = newY

	// Place in new position
	newCell.Occupant = ant
	return true
}
