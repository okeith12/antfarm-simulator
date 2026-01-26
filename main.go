package main

import (
	"antfarm/gui"
	"antfarm/logic"
	"antfarm/types"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

const (
	simulationUpdatesPerSecond = 1  // How fast ants move/act (1 = once per second)
	renderFPS                  = 30 // Frames per seconnnddd (30 FPS)
)

func main() {
	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	// Create world
	width, height := screen.Size()
	world := types.NewWorld(width, height-5) // Reserve space for stats and controls

	// Create initial colony
	queenX, queenY := width/4, height/3
	colony := types.NewColony("Red", queenX, queenY, tcell.ColorRed)

	world.AddColony(colony)

	// Create renderer
	renderer := gui.NewRenderer(screen)

	// Game loop with separate simulation and render rates
	simulationTicker := time.NewTicker(time.Second / simulationUpdatesPerSecond)
	renderTicker := time.NewTicker(time.Second / renderFPS)
	defer simulationTicker.Stop()
	defer renderTicker.Stop()

	running := true
	needsRender := true // Flag to track if we need to redraw

	for running {
		// Check for input events immediately
		for screen.HasPendingEvent() {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
					running = false
				}
				if ev.Rune() == 'l' || ev.Rune() == 'L' {
					renderer.ToggleLog()
					needsRender = true
				}
			case *tcell.EventResize:
				screen.Sync()
				needsRender = true
			}
		}

		select {
		case <-simulationTicker.C:
			// Update world state (controlled by simulationUpdatesPerSecond)
			logic.UpdateWorld(world)
			needsRender = true

		case <-renderTicker.C:
			// Render frequently for better visuals and responsive UI
			if needsRender {
				renderer.Render(world)
				needsRender = false
			}

		default:
			// Small sleep to prevent busy-waiting and CPU running
			time.Sleep(10 * time.Millisecond)
		}
	}
}
