package pathfinder

import (
	"antfarm/random"
	"antfarm/types"
	"testing"
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
	world := types.NewWorld(20, 20, random.New(1))
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, types.ColonyRed)

	// Create tunnels around queen
	for x := 8; x <= 12; x++ {
		for y := 8; y <= 12; y++ {
			world.GetCell(x, y).IsTunnel = true
		}
	}

	// Nurse within 2 cells should guard
	nurse := types.NewNurse(2, 11, 10, "Red")
	world.GetCell(11, 10).Occupant = nurse

	isGuarding := np.GuardNursery(world, colony, nurse)
	if !isGuarding {
		t.Error("Nurse within 2 cells should be guarding")
	}
}

func TestNurseGuardNurseryMovesCloser(t *testing.T) {
	world := types.NewWorld(20, 20, random.New(1))
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, types.ColonyRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		for y := 8; y <= 12; y++ {
			world.GetCell(x, y).IsTunnel = true
		}
	}
	world.GetCell(10, 10).Occupant = colony.Queen

	// Nurse far from queen should move closer
	nurse := types.NewNurse(2, 5, 10, "Red")
	world.GetCell(5, 10).Occupant = nurse

	isGuarding := np.GuardNursery(world, colony, nurse)
	if isGuarding {
		t.Error("Nurse far from queen should not be guarding yet")
	}
}

func TestNurseMoveTowardQueen(t *testing.T) {
	world := types.NewWorld(20, 20, random.New(1))
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, types.ColonyRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		world.GetCell(x, 10).IsTunnel = true
	}
	world.GetCell(10, 10).Occupant = colony.Queen

	nurse := types.NewNurse(2, 5, 10, "Red")
	world.GetCell(5, 10).Occupant = nurse

	success := np.MoveTowardQueen(world, colony, nurse)

	if !success {
		t.Error("Nurse should be able to move toward queen")
	}
}

func TestNurseMoveTowardLarvae(t *testing.T) {
	world := types.NewWorld(20, 20, random.New(1))
	np := NewNursePathfinder()

	colony := types.NewColony("Red", 10, 10, types.ColonyRed)

	// Create tunnels
	for x := 5; x <= 15; x++ {
		world.GetCell(x, 10).IsTunnel = true
	}
	world.GetCell(10, 10).Occupant = colony.Queen

	nurse := types.NewNurse(2, 5, 10, "Red")
	world.GetCell(5, 10).Occupant = nurse

	larvaePos := types.Position{X: 15, Y: 10}
	success := np.MoveTowardLarvae(world, colony, nurse, larvaePos)

	if !success {
		t.Error("Nurse should be able to move toward larvae")
	}
}
