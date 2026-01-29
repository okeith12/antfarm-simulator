package gui

import (
	"antfarm/logic"
	"antfarm/types"
	"time"

	"github.com/gdamore/tcell/v2"
)

// antfarm.go the main Antfarm struct which orchestrates the entire
// simulation loop, world updates, and rendering.

// Speed presets in ticks per second
var speedPresets = []float64{0.25, 0.5, 1, 2, 5, 10} // How fast icons move/act (1 = once per second)

const (
	defaultSpeedIndex = 2  // Index of 1.0 in speedPresets
	renderFPS         = 30 // Frames per seconnnddd (30 FPS)
)

// AntfarmState controls the current state of the Antfarm
type AntfarmState struct {
	running    bool // Controls the main loop - set false to exit
	paused     bool // When true, game stops but rendering continues
	speedIndex int  // Index into speedPreset array
}

// Antfarm manages the game loop and is the main struct that ties the farm together,
type Antfarm struct {
	screen   tcell.Screen // Terminal screen to render everything
	world    *types.World // The simulated world containing the colonies, ants, and terrain
	renderer *Renderer    // Handles all drawing operations
	state    AntfarmState
}

// GetSpeed returns the current simulation speed in ticks per second
func (a *Antfarm) GetSpeed() float64 {
	return speedPresets[a.state.speedIndex]
}

// IsPaused returns whether the simulation is paused
func (a *Antfarm) IsPaused() bool {
	return a.state.paused
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
		state: AntfarmState{
			running:    false,
			paused:     false,
			speedIndex: defaultSpeedIndex,
		},
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

	// Create ticker for rendering (fixed rate)
	renderTicker := time.NewTicker(time.Second / renderFPS)
	defer renderTicker.Stop()

	// Create initial simulation ticker
	simulationTicker := time.NewTicker(a.getTickDuration())
	defer simulationTicker.Stop()

	a.state.running = true
	needsRender := true // Flag to track if we need to redraw

	for a.state.running {
		speedChanged := false
		// Process any pending keyboard/mouse/resize events
		a.handleEvents(&needsRender, &speedChanged)

		// If speed changed, recreate the simulation ticker
		if speedChanged {
			simulationTicker.Stop()
			simulationTicker = time.NewTicker(a.getTickDuration())
		}

		select {
		case <-simulationTicker.C:
			// Time to update the world state
			// This moves ants, processes food, hatches eggs, etc.
			if !a.state.paused {
				logic.UpdateWorld(a.world)
				needsRender = true
			} // World changed, redraw

			// Render controls for the visuals and UI
		case <-renderTicker.C:
			// Time to potentially redraw the screen
			// Only render if something changed
			if needsRender {
				a.renderer.Render(a.world, a.state.paused, a.GetSpeed())
				needsRender = false
			}
			// sleep briefly to avoid busy-waiting; keep CPU cycle down :)
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// getTickDuration returns the duration between simulation ticks based on current speed
func (a *Antfarm) getTickDuration() time.Duration {
	return time.Duration(float64(time.Second) / a.GetSpeed())
}

// handleEvents processes all pending input events from the terminal.
// It handles:
//   - Q or Escape: Quit the application
//   - L: Toggle the activity log display
//   - Window resize: Sync the screen buffer
func (a *Antfarm) handleEvents(needsRender, speedChanged *bool) {
	for a.screen.HasPendingEvent() {
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			// Handle quit
			if ev.Key() == tcell.KeyEscape || ev.Rune() == 'q' || ev.Rune() == 'Q' {
				a.state.running = false
			}
			// Handle logs
			if ev.Rune() == 'l' || ev.Rune() == 'L' {
				a.renderer.ToggleLog()
				*needsRender = true
			}

			// Handle pause
			if ev.Rune() == 'p' || ev.Rune() == 'P' {
				a.state.paused = !a.state.paused
				*needsRender = true
			}

			// Handle speed up both are same key
			if ev.Rune() == '+' || ev.Rune() == '=' {
				if a.state.speedIndex < len(speedPresets)-1 {
					a.state.speedIndex++
					*speedChanged = true
					*needsRender = true
				}
			}
			// Handle slow down (-)
			if ev.Rune() == '-' {
				if a.state.speedIndex > 0 {
					a.state.speedIndex--
					*speedChanged = true
					*needsRender = true
				}
			}
		case *tcell.EventResize:
			// Terminal was resized - sync internal buffer to new size
			// Todo: update the entire world as well
			a.screen.Sync()
			*needsRender = true
		}
	}
}
