package types

// colony.go - Defines ant colonies (collections of ants with shared resources)
// Each colony has a queen, worker ants, food stores, and a unique color/identity

// ColonyColor identifies a colony's palette slot. It is a plain enum so the
// simulation stays render-agnostic; the display layer maps it to an actual color.
type ColonyColor int

const (
	ColonyRed ColonyColor = iota
	ColonyBlue
	ColonyGreen
	ColonyPurple
)

// Colony represents a group of ants that work together
// Contains the queen, all ants, shared resources, and colony identity
type Colony struct {
	Name          string        // Colony identifier (e.g. "Red", "Black")
	Color         ColonyColor   // Palette slot for this colony's ants
	Queen         *QueenAnt     // The reigning queen (center of the colony)
	Queens        []*QueenAnt   // Spare queens raised from larvae, heirs to the throne
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
// Every colony starts with the same three founders: one queen, one nurse to
// tend her brood, and one worker so food starts arriving from the first tick.
// Only the terrain and the dice vary between runs, never the opening lineup.
func NewColony(name string, queenX, queenY int, color ColonyColor) *Colony {
	queen := NewQueen(0, queenX, queenY, name)
	headNurse := NewNurse(1, queenX+1, queenY, name)    // Head nurse starts next to queen
	firstWorker := NewWorker(2, queenX-1, queenY, name) // First worker starts on her other side

	return &Colony{
		Name:          name,
		Color:         color,
		Queen:         queen,
		HeadNurse:     headNurse,
		Queens:        []*QueenAnt{},
		Nurses:        []*NurseAnt{},
		Workers:       []*WorkerAnt{firstWorker},
		Soldiers:      []*SoldierAnt{},
		Larvae:        []*LarvaeAnt{},
		Food:          50 * FoodScale, // Starting food, 50 food
		Eggs:          0,
		NextAntID:     3, // Start at 3: queen=0, head nurse=1, first worker=2
		QueenPosition: Position{queenX, queenY},
	}
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
