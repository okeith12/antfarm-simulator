package pathfinder

import (
	"antfarm/types"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewNursePathfinder(t *testing.T) {
	np := NewNursePathfinder()
	if np == nil {
		t.Error("NewNursePathfinder should not return nil")
	}
}

func TestNurseIsAdjacentToLarvae(t *testing.T) {
	np := NewNursePathfinder()

	nurse := types.NewNurse(1, 10, 10, "Red")

	if !np.IsAdjacentToLarvae(nurse, types.Position{X: 11, Y: 10}) {
		t.Error("Nurse should be adjacent to (11,10)")
	}
	if !np.IsAdjacentToLarvae(nurse, types.Position{X: 10, Y: 10}) {
		t.Error("Nurse should be adjacent to same position")
	}
	if np.IsAdjacentToLarvae(nurse, types.Position{X: 15, Y: 10}) {
		t.Error("Nurse should not be adjacent to (15,10)")
	}
}

func TestNurseGuardNursery(t *testing.T) {
	world := types.NewWorld(20, 20)
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)

	// Create tunnels around queen
	for x := 8; x <= 12; x++ {
		for y := 8; y <= 12; y++ {
			world.Grid[y][x].IsTunnel = true
		}
	}

	// Nurse within 2 cells should guard
	nurse := types.NewNurse(2, 11, 10, "Red")
	world.Grid[10][11].Occupant = nurse

	isGuarding := np.GuardNursery(world, colony, nurse)
	if !isGuarding {
		t.Error("Nurse within 2 cells should be guarding")
	}
}

func TestNurseGuardNurseryMovesCloser(t *testing.T) {
	world := types.NewWorld(20, 20)
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		for y := 8; y <= 12; y++ {
			world.Grid[y][x].IsTunnel = true
		}
	}
	world.Grid[10][10].Occupant = colony.Queen

	// Nurse far from queen should move closer
	nurse := types.NewNurse(2, 5, 10, "Red")
	world.Grid[10][5].Occupant = nurse

	isGuarding := np.GuardNursery(world, colony, nurse)
	if isGuarding {
		t.Error("Nurse far from queen should not be guarding yet")
	}
}

func TestNurseMoveTowardQueen(t *testing.T) {
	world := types.NewWorld(20, 20)
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		world.Grid[10][x].IsTunnel = true
	}
	world.Grid[10][10].Occupant = colony.Queen

	nurse := types.NewNurse(2, 5, 10, "Red")
	world.Grid[10][5].Occupant = nurse

	success := np.MoveTowardQueen(world, colony, nurse)

	if !success {
		t.Error("Nurse should be able to move toward queen")
	}
}

func TestNurseMoveTowardLarvae(t *testing.T) {
	world := types.NewWorld(20, 20)
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, tcell.ColorRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		world.Grid[10][x].IsTunnel = true
	}
	world.Grid[10][10].Occupant = colony.Queen

	nurse := types.NewNurse(2, 5, 10, "Red")
	world.Grid[10][5].Occupant = nurse

	larvaePos := types.Position{X: 15, Y: 10}
	success := np.MoveTowardLarvae(world, colony, nurse, larvaePos)

	if !success {
		t.Error("Nurse should be able to move toward larvae")
	}
}
