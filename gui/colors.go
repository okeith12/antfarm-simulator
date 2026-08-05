package gui

import (
	"antfarm/types"

	"github.com/gdamore/tcell/v2"
)

// colors.go - The only place that maps simulation enums to display colors.
// types/ and simulation/ stay render-agnostic; nothing there imports tcell.

// SoilColor returns the background color for a cell's terrain.
// A dug tunnel reads as empty space regardless of the soil it was cut through.
func SoilColor(soil types.Soil, isTunnel bool) tcell.Color {
	if isTunnel {
		return tcell.ColorBlack
	}

	switch soil {
	case types.Sand:
		return tcell.ColorYellow
	case types.Dirt:
		return tcell.ColorMaroon
	case types.Clay:
		return tcell.ColorOlive
	case types.Rock:
		return tcell.ColorGray
	default:
		return tcell.ColorBlack
	}
}

// ColonyColor returns the foreground color used to draw a colony's ants.
func ColonyColor(c types.ColonyColor) tcell.Color {
	switch c {
	case types.ColonyRed:
		return tcell.ColorRed
	case types.ColonyBlue:
		return tcell.ColorBlue
	case types.ColonyGreen:
		return tcell.ColorGreen
	case types.ColonyPurple:
		return tcell.ColorPurple
	default:
		return tcell.ColorWhite
	}
}
