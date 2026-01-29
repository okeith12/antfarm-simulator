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
func generateSoilType(y, _ int) Soil {
	// put back maxHeigt when using it
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
