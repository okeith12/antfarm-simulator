package types

import (
	"testing"
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
