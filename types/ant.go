package types

// ant.go - Defines our ants and their properties
// Ants have many different roles and behaviors

// Role defines what job/type an ant has in the colony
type Role int

const (
	Worker  Role = iota // Digs tunnels, gathers food
	Soldier             // Defends colony, fights enemies
	Nurse               // Tends to eggs and larvae
	Queen               // Lays eggs, center of colony
	Larvae              // baby ants
)

// Position represents a coordinate in the world
type Position struct {
	X, Y int
}

type Ant struct {
	ID            int      // Unique identifier for this ant
	Role          Role     // What job this ant perfoms
	Position      Position // Current position in the world grid
	Health        int
	ColonyID      string // Which colony this ant belongs to
	Age           int    // How long the ant been alive
	CurrentAction string // What is the ant currently doing
}

// NewAnt creates a new ant with the given properties
func NewAnt(id int, role Role, x, y int, colonyID string) *Ant {
	return &Ant{
		ID:            id,
		Role:          role,
		Position:      Position{X: x, Y: y},
		Health:        100,
		ColonyID:      colonyID,
		Age:           0,
		CurrentAction: "idle",
	}
}

// GetAntIcon returns the symbol used to display this ant
func (a *Ant) GetAntIcon() rune {
	switch a.Role {
	case Queen:
		return '♛'
	case Soldier:
		return '⚔'
	case Worker:
		return '●'
	case Nurse:
		return '○'
	case Larvae:
		return '◦'
	default:
		return '●'
	}
}

// AntInterface defines common behavior all ant types must implement
type AntInterface interface {
	GetAnt() *Ant     // Returns the default Ant
	GetAntIcon() rune // Returns the display icon for this ant
	GetRole() Role    // Returns the ant's role
}
