package gui

import (
	"antfarm/types"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// controls.go - Handles rendering of control hints and UI controls
// Displays keyboard shortcuts and interactive elements

// renderControls displays the control hints at the bottom of the screen
func (r *Renderer) renderControls(world *types.World, paused bool, speed float64) {
	y := world.Height + 3 // Below stats

	// Build status string
	status := "RUNNING"
	statusColor := tcell.ColorGreen
	if paused {
		status = "PAUSED"
		statusColor = tcell.ColorRed
	}

	// Format speed display
	speedStr := fmt.Sprintf("%.2fx", speed)

	controls := fmt.Sprintf("[%s] Speed: %s | Q=Quit | L=Log | P=Pause | +/- =Speed", status, speedStr)

	// Draw status with color
	statusStyle := tcell.StyleDefault.Foreground(statusColor).Background(tcell.ColorDefault)
	normalStyle := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorDefault)

	x := 0
	// Draw opening bracket
	r.screen.SetContent(x, y, '[', nil, normalStyle)
	x++

	// Draw status with color
	for _, ch := range status {
		r.screen.SetContent(x, y, ch, nil, statusStyle)
		x++
	}

	// Draw rest of controls
	rest := fmt.Sprintf("] Speed: %s | Q=Quit | L=Log | P=Pause | +/- =Speed", speedStr)
	for _, ch := range rest {
		r.screen.SetContent(x, y, ch, nil, normalStyle)
		x++
	}

	// Clear rest of line
	width, _ := r.screen.Size()
	for ; x < width; x++ {
		r.screen.SetContent(x, y, ' ', nil, normalStyle)
	}

	_ = controls // Avoid unused variable warning
}
