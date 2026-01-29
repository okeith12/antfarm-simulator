package gui

import (
	"antfarm/types"

	"github.com/gdamore/tcell/v2"
)

// renderer.go - Handles all visual display of the simulation
// Renders the world grid, ants, and statistics to the terminal

// Renderer manages drawing the simulation to the screen
type Renderer struct {
	screen       tcell.Screen
	logExpanded  bool
	maxAntsToLog int // How many ants to log
}

// NewRenderer creates a new renderer with the given screen
func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{
		screen:       screen,
		logExpanded:  false,
		maxAntsToLog: 10,
	}
}

// Render draws the entire world state to the screen
// Shows terrain, tunnels, ants, and colony statistics
func (r *Renderer) Render(world *types.World) {
	r.screen.Clear()

	// Draw queen chambers first (background layer) --- FIX IT
	for _, colony := range world.Colonies {
		if colony.Queen != nil {
			qx, qy := colony.Queen.Position.X, colony.Queen.Position.Y
			// Draw 3x3 chamber around queen
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					x, y := qx+dx, qy+dy
					if world.IsValidPosition(x, y) {
						// Draw chamber walls (skip center where queen sits)
						if dx != 0 || dy != 0 {
							style := tcell.StyleDefault.
								Foreground(tcell.ColorGold).
								Background(tcell.ColorBlack)
							r.screen.SetContent(x, y, '▢', nil, style)
						}
					}
				}
			}
		}
	}

	// Draw the world grid (terrain and ants)
	for y := 0; y < world.Height; y++ {
		for x := 0; x < world.Width; x++ {
			cell := world.Grid[y][x]

			var ch rune
			var fgColor tcell.Color
			var bgColor tcell.Color

			if cell.Occupant != nil {
				// There's an ant in this cell - draw the ant
				ch = cell.Occupant.GetIcon()

				// Find the ant's colony to get its color
				for _, colony := range world.Colonies {
					if colony.Name == cell.Occupant.GetAnt().ColonyID {
						fgColor = colony.Color
						break
					}
				}

				// Ants in tunnels have a default background
				if cell.IsTunnel {
					bgColor = tcell.ColorDefault
				} else {
					bgColor = cell.GetColor()
				}
			} else {
				// No ant - draw the terrain
				ch = cell.GetIcon()

				// Tunnels have black background with grey/white foreground
				if cell.IsTunnel {
					bgColor = tcell.ColorDefault
					fgColor = tcell.ColorBlack
				} else {
					// Soil cells: foreground and background same color for solid blocks
					bgColor = tcell.ColorBlack
					fgColor = tcell.ColorDefault
				}
			}

			style := tcell.StyleDefault.Foreground(fgColor).Background(bgColor)
			r.screen.SetContent(x, y, ch, nil, style)
		}
	}

	// Draw statistics at the bottom
	r.renderStats(world)
	r.renderControls(world)
	r.screen.Show()
}

// ToggleLog toggles the expanded log view
func (r *Renderer) ToggleLog() {
	r.logExpanded = !r.logExpanded
}
