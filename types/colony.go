package types

import "github.com/gdamore/tcell/v2"

// colony.go - Defines ant colonies (collections of ants with shared resources)
// Each colony has a queen, worker ants, food stores, and a unique color/identity

// Colony represents a group of ants that work together
// Contains the queen, all ants, shared resources, and colony identity
type Colony struct {
	Name          string        // Colony identifier (e.g. "Red", "Black")
	Color         tcell.Color   // Display color for this colony's ants
	Queen         *QueenAnt     // The queen ant (center of the colony)
	HeadNurse     *NurseAnt     // The primary nurse ant (second in command)
	Nurses        []*NurseAnt   // All other nurse ants
	Workers       []*WorkerAnt  // All worker ants
	Soldiers      []*SoldierAnt // All soldier ants
	Larvae        []*LarvaeAnt  // All larvae waiting to grow
	Food          int           // Shared food stockpile
	Eggs          int           // Number of eggs waiting to hatch
	NextAntID     int           // Counter for generating unique ant IDs
	QueenPosition Position      // Position of the queen (center of colony)
}

// NewColony creates a new ant colony with a queen and head nurse at the specified position
// Starts with initial food and spawns the queen and head nurse as the first ants
func NewColony(name string, queenX, queenY int, color tcell.Color) *Colony {
	queen := NewQueen(0, queenX, queenY, name)
	headNurse := NewNurse(1, queenX+1, queenY, name) // Head nurse starts next to queen

	return &Colony{
		Name:          name,
		Color:         color,
		Queen:         queen,
		HeadNurse:     headNurse,
		Nurses:        []*NurseAnt{},
		Workers:       []*WorkerAnt{},
		Soldiers:      []*SoldierAnt{},
		Larvae:        []*LarvaeAnt{},
		Food:          50, // Starting food
		Eggs:          0,
		NextAntID:     2, // Start at 2 since queen=0, head nurse=1
		QueenPosition: Position{queenX, queenY},
	}
}

// SpawnWorker creates a new worker ant at the given position
func (c *Colony) SpawnWorker(x, y int) *WorkerAnt {
	worker := NewWorker(c.NextAntID, x, y, c.Name)
	c.NextAntID++
	c.Workers = append(c.Workers, worker)
	return worker
}

// SpawnWorkerWithID creates a new worker ant with a specific ID (used when larvae becomes worker)
func (c *Colony) SpawnWorkerWithID(id int, x, y int) *WorkerAnt {
	worker := NewWorker(id, x, y, c.Name)
	c.Workers = append(c.Workers, worker)
	return worker
}

// SpawnSoldier creates a new soldier ant at the given position
func (c *Colony) SpawnSoldier(x, y int) *SoldierAnt {
	soldier := NewSoldier(c.NextAntID, x, y, c.Name)
	c.NextAntID++
	c.Soldiers = append(c.Soldiers, soldier)
	return soldier
}

// SpawnNurse creates a new nurse ant at the given position
func (c *Colony) SpawnNurse(x, y int) *NurseAnt {
	nurse := NewNurse(c.NextAntID, x, y, c.Name)
	c.NextAntID++
	c.Nurses = append(c.Nurses, nurse)
	return nurse
}

// SpawnLarvae creates a new larvae at the given position
func (c *Colony) SpawnLarvae(x, y int) *LarvaeAnt {
	larvae := NewLarvae(c.NextAntID, x, y, c.Name)
	c.NextAntID++
	c.Larvae = append(c.Larvae, larvae)
	return larvae
}

// GetAllAnts returns all ants in the colony as AntInterface slice
// Useful for iteration when you need to process all ants regardless of role
func (c *Colony) GetAllAnts() []AntInterface {
	var all []AntInterface

	// Add queen
	if c.Queen != nil {
		all = append(all, c.Queen)
	}

	// Add head nurse
	if c.HeadNurse != nil {
		all = append(all, c.HeadNurse)
	}

	// Add other nurses
	for _, n := range c.Nurses {
		all = append(all, n)
	}

	// Add workers
	for _, w := range c.Workers {
		all = append(all, w)
	}

	// Add soldiers
	for _, s := range c.Soldiers {
		all = append(all, s)
	}

	// Add larvae
	for _, l := range c.Larvae {
		all = append(all, l)
	}

	return all
}

// GetAntCount returns the total number of ants in the colony
func (c *Colony) GetAntCount() int {
	count := 0
	if c.Queen != nil {
		count++
	}
	if c.HeadNurse != nil {
		count++
	}
	count += len(c.Nurses)
	count += len(c.Workers)
	count += len(c.Soldiers)
	count += len(c.Larvae)
	return count
}

// RemoveLarvae removes a larvae from the colony's larvae list
func (c *Colony) RemoveLarvae(larvae *LarvaeAnt) {
	for i, l := range c.Larvae {
		if l.ID == larvae.ID {
			c.Larvae = append(c.Larvae[:i], c.Larvae[i+1:]...)
			return
		}
	}
}
