package gui

import (
	"antfarm/logic"
	"antfarm/types"
	"time"

	"github.com/gdamore/tcell/v2"
)

// antfarm.go the main Antfarm struct which orchestrates the entire
// simulation loop, world updates, and rendering.

const (
	simulationUpdatesPerSecond = 1  // How fast icons move/act (1 = once per second)
	renderFPS                  = 30 // Frames per seconnnddd (30 FPS)
)

// Antfarm manages the game loop and is the main struct that ties the farm together,
type Antfarm struct {
	screen   tcell.Screen // Terminal screen to render everything
	world    *types.World // The simulated world containing the colonies, ants, and terrain
	renderer *Renderer    // Handles all drawing operations
	running  bool         // Controls the main loop - set false to exit
}

// NewAntfarm creates and initializes a new Antfarm instance.
// It sets up the terminal screen, creates the world sized to fit the terminal,
// spawns an initial colony, and prepares the renderer.
//
// Returns an error if screen initialization fails
func NewAntfarm() (*Antfarm, error) {
	// Initialize screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}
	if err := screen.Init(); err != nil {
		return nil, err
	}

	// Create world
	width, height := screen.Size()
	world := types.NewWorld(width, height-5) // Reserve space for stats and controls

	// Create initial colony
	// Queen position determines where the colony's nest begins
	queenX, queenY := width/4, height/3
	colony := types.NewColony("Red", queenX, queenY, tcell.ColorRed)
	logic.AddColony(world, colony)

	// Create renderer
	renderer := NewRenderer(screen)

	return &Antfarm{
		screen:   screen,
		world:    world,
		renderer: renderer,
		running:  false,
	}, nil
}

// Run starts the main simulation loop. This is a blocking call that runs until
// the user quits (Q/Escape) or an error occurs.
//
// The loop uses two independent timers:
//   - simulationTicker: Controls world updates (ant movement, food gathering, etc.)
//   - renderTicker: Controls screen redraws for smooth visuals
func (a *Antfarm) Run() {
	// Ensure we clean up the terminal when done
	defer a.screen.Fini()

	// Game loop with separate simulation and render rates
	simulationTicker := time.NewTicker(time.Second / simulationUpdatesPerSecond)
	renderTicker := time.NewTicker(time.Second / renderFPS)
	defer simulationTicker.Stop()
	defer renderTicker.Stop()

	a.running = true
	needsRender := true // Flag to track if we need to redraw

	for a.running {
		// Process any pending keyboard/mouse/resize events
		a.handleEvents(&needsRender)

		select {
		case <-simulationTicker.C:
			// Time to update the world state
			// This moves ants, processes food, hatches eggs, etc.
			logic.UpdateWorld(a.world)
			needsRender = true // World changed, redraw

			// Render controls for the visuals and UI
		case <-renderTicker.C:
			// Time to potentially redraw the screen
			// Only render if something changed
			if needsRender {
				a.renderer.Render(a.world)
				needsRender = false
			}
			// sleep briefly to avoid busy-waiting; keep CPU cycle down :)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// handleEvents processes all pending input events from the terminal.
// It handles:
//   - Q or Escape: Quit the application
//   - L: Toggle the activity log display
//   - Window resize: Sync the screen buffer
func (a *Antfarm) handleEvents(needsRender *bool) {
	for a.screen.HasPendingEvent() {
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Handle quit
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
				a.running = false
			}
			// Handle logs
			if ev.Rune() == 'l' || ev.Rune() == 'L' {
				a.renderer.ToggleLog()
				*needsRender = true
			}
		case *tcell.EventResize:
			// Terminal was resized - sync internal buffer to new size
			// Todo: update the entire world as well
			a.screen.Sync()
			*needsRender = true
		}
	}
}
