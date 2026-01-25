package types

import "fmt"

// lgoo.go - Tracks what ants are currently doing
// Used for logging and debugging ant behavior
// AntActivity tracks what an ant is currently doing
type AntActivity struct {
	AntID    int
	Action   string // "digging", "moving", "resting", "nursing"
	Location Position
}

func (a *Ant) GetActivityString(worldNum int) string {
	roleName := map[Role]string{
		Worker:  "Worker",
		Soldier: "Soldier",
		Nurse:   "Nurse",
		Queen:   "Queen",
		Larvae:  "Larvae",
	}

	return fmt.Sprintf("%s_%s_Ant_%d is currently at (%d,%d) in World_%d",
		a.ColonyID, roleName[a.Role], a.ID, a.Position.X, a.Position.Y, worldNum)
}
