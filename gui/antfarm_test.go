package gui

import (
	"antfarm/logic"
	"antfarm/types"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

// mockScreen creates a simulation screen for testing.
func mockScreen() tcell.SimulationScreen {
	screen := tcell.NewSimulationScreen("")
	screen.Init()
	screen.SetSize(80, 24)
	return screen
}

// mockAntfarm creates an Antfarm using a simulation screen for testing.
func mockAntfarm(screen tcell.SimulationScreen) *Antfarm {
	width, height := screen.Size()
	world := types.NewWorld(width, height-5)

	queenX, queenY := width/4, height/3
	colony := types.NewColony("Red", queenX, queenY, tcell.ColorRed)
	logic.AddColony(world, colony)

	renderer := NewRenderer(screen)

	return &Antfarm{
		screen:   screen,
		world:    world,
		renderer: renderer,
		running:  false,
	}
}

// TestNewAntfarm tests using the mock screen.
func TestNewAntfarm(t *testing.T) {
	screen := mockScreen()
	defer screen.Fini()

	// Verify screen is usable
	screen.SetSize(80, 24)
	width, height := screen.Size()

	if width != 80 || height != 24 {
		t.Errorf("Expected 80x24, got %dx%d", width, height)
	}
}

// TestAntfarmRunQuitsOnEscape tests that the Run loop exits when Escape is pressed.
func TestAntfarmRunQuitsOnEscape(t *testing.T) {
	screen := mockScreen()
	antfarm := mockAntfarm(screen)

	done := make(chan bool)
	go func() {
		antfarm.Run()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	// Inject Escape key event
	screen.InjectKey(tcell.KeyEscape, 0, tcell.ModNone)

	// Wait for Run to exit
	select {
	case <-done:
		// Success - Run exited
	case <-time.After(1 * time.Second):
		t.Error("Run did not exit after Escape key")
	}
}

// TestAntfarmRunQuitsOnQ tests that the Run loop exits when Q is pressed.
func TestAntfarmRunQuitsOnQ(t *testing.T) {
	screen := mockScreen()
	antfarm := mockAntfarm(screen)

	done := make(chan bool)
	go func() {
		antfarm.Run()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	screen.InjectKey(tcell.KeyRune, 'q', tcell.ModNone)

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Error("Run did not exit after Q key")
	}
}

// TestAntfarmHandleEventsToggleLog tests that L key toggles the log display.
func TestAntfarmHandleEventsToggleLog(t *testing.T) {
	screen := mockScreen()
	antfarm := mockAntfarm(screen)
	antfarm.running = true

	initialLogState := antfarm.renderer.logExpanded

	// Inject 'L' key
	screen.InjectKey(tcell.KeyRune, 'L', tcell.ModNone)

	needsRender := false
	antfarm.handleEvents(&needsRender)

	if antfarm.renderer.logExpanded == initialLogState {
		t.Error("Log should have toggled")
	}
	if !needsRender {
		t.Error("needsRender should be true after toggle")
	}
}

// TestAntfarmWorldCreated tests that the world is properly created.
func TestAntfarmWorldCreated(t *testing.T) {
	screen := mockScreen()
	antfarm := mockAntfarm(screen)

	if antfarm.world == nil {
		t.Fatal("World should not be nil")
	}
	if antfarm.world.Width != 80 {
		t.Errorf("Expected world width 80, got %d", antfarm.world.Width)
	}
	// Height is screen height - 5 for stats
	if antfarm.world.Height != 19 {
		t.Errorf("Expected world height 19, got %d", antfarm.world.Height)
	}
}

// TestAntfarmColonyCreated tests that a colony exists in the world.
func TestAntfarmColonyCreated(t *testing.T) {
	screen := mockScreen()
	antfarm := mockAntfarm(screen)

	if len(antfarm.world.Colonies) != 1 {
		t.Errorf("Expected 1 colony, got %d", len(antfarm.world.Colonies))
	}
	if antfarm.world.Colonies[0].Name != "Red" {
		t.Errorf("Expected colony name 'Red', got '%s'", antfarm.world.Colonies[0].Name)
	}
}
