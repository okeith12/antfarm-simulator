package types

// cell.go - Defines the basic unit of the world grid
// Each cell represents a single position in the 2D ant farm world and contains
// information about terrain type, whether it's been dug into a tunnel, and what occupies it

// Soil represents different types of terrain in the world
type Soil int

const (
	Sand  Soil = iota // Easiest soil to dig through
	Dirt              // Standard soil to dig through
	Clay              // Hard soil to dig through
	Rock              // Soil that cannot be dug through
	Empty             // Tunnel/empty space
)

// Cell represents a single position in the world grid
// It tracks terrain type, whether it's been tunneled, and what occupies the space
type Cell struct {
	Soil     Soil
	IsTunnel bool
	Occupant AntInterface // nil if empty
	Food     int          // Food stored in this cell
	// Moisture  int  // 0-100
	// Stability int  // How stable the cell is (affects collapse)
}

// NewCell creates a new cell with the given soil type
// By default cells are not tunnels and have no occupants
func NewCell(soil Soil) *Cell {
	return &Cell{
		Soil:     soil,
		IsTunnel: false,
		Occupant: nil,
		Food:     0,
		// Moisture:  50,
		// Stability: 100,
	}
}

// GetIcon returns the visual character to display for this cell
func (c *Cell) GetIcon() rune {
	if c.Soil == Empty {
		if c.Food > 0 {
			return '🌾' // Grass with food
		}
		if c.Food == -1 {
			return ' ' // Harvested - empty
		}
		return '🌱' // Just grass
	}

	// soil textures
	switch c.Soil {
	case Sand:
		return '░'
	case Dirt:
		return '▒'
	case Clay:
		return '▓'
	case Rock:
		return '█'
	default:
		return ' '
	}
}
