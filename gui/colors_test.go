package gui

import (
	"testing"

	"antfarm/types"

	"github.com/gdamore/tcell/v2"
)

func TestSoilColor(t *testing.T) {
	tests := []struct {
		soil     types.Soil
		isTunnel bool
		expected tcell.Color
	}{
		{types.Sand, false, tcell.ColorYellow},
		{types.Dirt, false, tcell.ColorMaroon},
		{types.Clay, false, tcell.ColorOlive},
		{types.Rock, false, tcell.ColorGray},
		{types.Sand, true, tcell.ColorBlack},
	}

	for _, tt := range tests {
		if got := SoilColor(tt.soil, tt.isTunnel); got != tt.expected {
			t.Errorf("SoilColor(%d, %v): expected %v, got %v", tt.soil, tt.isTunnel, tt.expected, got)
		}
	}
}

func TestColonyColor(t *testing.T) {
	tests := []struct {
		colony   types.ColonyColor
		expected tcell.Color
	}{
		{types.ColonyRed, tcell.ColorRed},
		{types.ColonyBlue, tcell.ColorBlue},
		{types.ColonyGreen, tcell.ColorGreen},
		{types.ColonyPurple, tcell.ColorPurple},
	}

	for _, tt := range tests {
		if got := ColonyColor(tt.colony); got != tt.expected {
			t.Errorf("ColonyColor(%d): expected %v, got %v", tt.colony, tt.expected, got)
		}
	}
}
