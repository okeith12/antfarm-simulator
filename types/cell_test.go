package types

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewCell(t *testing.T) {
	cell := NewCell(Sand)

	if cell.Soil != Sand {
		t.Errorf("Expected Soil Sand, got %d", cell.Soil)
	}
	if cell.IsTunnel {
		t.Error("Expected IsTunnel false")
	}
	if cell.Occupant != nil {
		t.Error("Expected Occupant nil")
	}
	if cell.Food != 0 {
		t.Errorf("Expected Food 0, got %d", cell.Food)
	}
}

func TestGetColor(t *testing.T) {
	tests := []struct {
		soil     Soil
		isTunnel bool
		expected tcell.Color
	}{
		{Sand, false, tcell.ColorYellow},
		{Dirt, false, tcell.ColorMaroon},
		{Clay, false, tcell.ColorOlive},
		{Rock, false, tcell.ColorGray},
		{Sand, true, tcell.ColorBlack},
	}

	for _, tt := range tests {
		cell := NewCell(tt.soil)
		cell.IsTunnel = tt.isTunnel
		if cell.GetColor() != tt.expected {
			t.Errorf("GetColor() for soil %d, tunnel %v: expected %v", tt.soil, tt.isTunnel, tt.expected)
		}
	}
}

func TestCellGetIcon(t *testing.T) {
	tests := []struct {
		soil     Soil
		food     int
		expected rune
	}{
		{Sand, 0, '░'},
		{Dirt, 0, '▒'},
		{Clay, 0, '▓'},
		{Rock, 0, '█'},
		{Empty, 5, '🌾'},
		{Empty, 0, '🌱'},
		{Empty, -1, ' '},
	}

	for _, tt := range tests {
		cell := NewCell(tt.soil)
		cell.Food = tt.food
		if cell.GetIcon() != tt.expected {
			t.Errorf("GetCellIcon() for soil %d, food %d: expected '%c'", tt.soil, tt.food, tt.expected)
		}
	}
}
