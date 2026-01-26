package gui

import (
	"antfarm/types"

	"github.com/gdamore/tcell/v2"
)

// controls.go - Handles rendering of control hints and UI controls
// Displays keyboard shortcuts and interactive elements

// renderControls displays the control hints at the bottom of the screen
func (r *Renderer) renderControls(world *types.World) {
	y := world.Height + 3 // Below stats

	controls := "Controls: Ctrl+Q=Quit |+=Speed Up |-=Slow Down |P=Pause"
	style := tcell.StyleDefault.Foreground(tcell.ColorYellow).Background(tcell.ColorDefault)

	for i, ch := range controls {
		r.screen.SetContent(i, y, ch, nil, style)
	}
}
