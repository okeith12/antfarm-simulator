package gui

import (
	"antfarm/types"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// stats.go - Handles rendering of colony statistics and activity logs
// Displays tick count, ant populations, food, eggs, and activity logs

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
	var logHeader string
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
		r.renderActivityLog(world, y)
	}
}

// renderActivityLog displays the activity log for ants
func (r *Renderer) renderActivityLog(world *types.World, startY int) {
	y := startY
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
