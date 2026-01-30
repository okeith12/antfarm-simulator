package types

// ant.go - Defines our ants and their properties
// Ants have many different roles and behaviors

// Lifespan constants (in ticks) - different ant types live different lengths
const (
	WorkerMaxTick  = 500   // Workers live shorter lives due to hard labor
	SoldierMaxTick = 600   // Soldiers live a bit longer
	NurseMaxTick   = 700   // Nurses live longer, babysitters
	QueenMaxTick   = 20000 // Queens live basically forever until I make it the opp
	LarvaeMaxTick  = 200   // Larvae must become workers before this or they die
)

// Health constants - different ant types have different health pools
const (
	WorkerMaxHealth  = 100
	SoldierMaxHealth = 150
	NurseMaxHealth   = 80
	QueenMaxHealth   = 2000
	LarvaeMaxHealth  = 50
)

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
	MaxHealth     int    // Maximum health for this ant type
	ColonyID      string // Which colony this ant belongs to
	Age           int    // How long the ant been alive
	MaxAge        int    // Maximum age before dying of old age
	CurrentAction string // What is the ant currently doing
}

// NewAnt creates a new ant with the given properties
func NewAnt(id int, role Role, x, y int, colonyID string, maxHealth, maxAge int) *Ant {
	return &Ant{
		ID:            id,
		Role:          role,
		Position:      Position{X: x, Y: y},
		Health:        100,
		MaxHealth:     maxHealth,
		ColonyID:      colonyID,
		Age:           0,
		MaxAge:        maxAge,
		CurrentAction: "idle",
	}
}

// AntInterface defines common behavior all ant types must implement
type AntInterface interface {
	GetAnt() *Ant  // Returns the default Ant
	GetIcon() rune // Returns the display icon for this ant
	GetRole() Role // Returns the ant's role
}

// IsDead checks if an ant should die (health depleted or old age)
func (a *Ant) IsDead() bool {
	return a.Health <= 0 || a.Age >= a.MaxAge
}
