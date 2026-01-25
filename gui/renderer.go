package gui

import (
	"antfarm/types"
	"fmt"

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
		maxAntsToLog: 5,
	}
}

// Render draws the entire world state to the screen
// Shows terrain, tunnels, ants, and colony statistics
func (r *Renderer) Render(world *types.World) {
	r.screen.Clear()

	// Draw queen chambers first (background layer)
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
				ch = cell.Occupant.GetAntIcon()

				// Find the ant's colony to get its color
				for _, colony := range world.Colonies {
					if colony.Name == cell.Occupant.GetAnt().ColonyID {
						fgColor = colony.Color
						break
					}
				}

				// Ants in tunnels have a terminal-style background
				if cell.IsTunnel {
					bgColor = Tunnel
				} else {
					bgColor = cell.GetColor()
				}
			} else {
				// No ant - draw the terrain
				ch = cell.GetCellIcon()
				fgColor = cell.GetColor()

				// Tunnels have black background with grey/white foreground
				if cell.IsTunnel {
					bgColor = tcell.ColorBlack
					fgColor = tcell.ColorWhite
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

// renderStats displays colony information and simulation statistics
// Shows tick count, ant populations, food, and eggs for each colony
func (r *Renderer) renderStats(world *types.World) {
	y := world.Height + 1
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorDefault)

	// Overall Simulation Stats
	statsLine := fmt.Sprintf("Ticks: %d | Press 'q' or ESC to quit", world.Ticks)
	for i, ch := range statsLine {
		r.screen.SetContent(i, y, ch, nil, style)
	}

	y++
	for _, colony := range world.Colonies {
		colonyStats := fmt.Sprintf("%s Colony: %d ants | Food: %d | Eggs: %d | Larvae: %d",
			colony.Name, colony.GetAntCount(), colony.Food, colony.Eggs, len(colony.Larvae))
		style = tcell.StyleDefault.Foreground(colony.Color).Background(tcell.ColorDefault)
		for i, ch := range colonyStats {
			r.screen.SetContent(i, y, ch, nil, style)
		}
		y++
	}

	// Activity Log Header (always visible)
	y++
	logStyle := tcell.StyleDefault.Foreground(tcell.ColorDarkCyan).Background(tcell.ColorDefault)
	logHeader := "Log"
	if r.logExpanded {
		logHeader = "Log (Press L to collapse)"
	} else {
		logHeader = "Log (Press L to expand)"
	}
	for i, ch := range logHeader {
		r.screen.SetContent(i, y, ch, nil, logStyle)
	}
	y++

	// Only show activity if expanded
	if r.logExpanded {
		antCount := 0
		for _, colony := range world.Colonies {
			allAnts := colony.GetAllAnts()
			for _, ant := range allAnts {
				if antCount >= r.maxAntsToLog {
					break
				}

				baseAnt := ant.GetAnt()
				roleStr := getRoleString(baseAnt.Role)

				action := baseAnt.CurrentAction
				if action == "" {
					action = "idle"
				}

				logLine := fmt.Sprintf("%s_%s_Ant_%d is %s at (%d,%d)",
					colony.Name, roleStr, baseAnt.ID, action, baseAnt.Position.X, baseAnt.Position.Y)

				style := tcell.StyleDefault.Foreground(colony.Color).Background(tcell.ColorDefault)
				for i, ch := range logLine {
					r.screen.SetContent(i, y, ch, nil, style)
				}
				y++
				antCount++
			}
			if antCount >= r.maxAntsToLog {
				break
			}
		}
	}
}

// getRoleString converts a Role enum to a display string
func getRoleString(role types.Role) string {
	switch role {
	case types.Worker:
		return "Worker"
	case types.Soldier:
		return "Soldier"
	case types.Queen:
		return "Queen"
	case types.Nurse:
		return "Nurse"
	case types.Larvae:
		return "Larvae"
	default:
		return "Unknown"
	}
}

func (r *Renderer) renderControls(world *types.World) {
	y := world.Height + 3 // Below stats

	controls := "Controls: Ctrl+Q=Quit | Ctrl+>=Speed Up | Ctrl+<=Slow Down | P=Pause"
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorDefault)

	for i, ch := range controls {
		r.screen.SetContent(i, y, ch, nil, style)
	}
}
