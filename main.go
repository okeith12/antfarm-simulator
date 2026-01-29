package main

import (
	"antfarm/gui"
	"log"
)

func main() {

	antfarm, err := gui.NewAntfarm()
	if err != nil {
		log.Fatalf("Failed to create antfarm: %v", err)
	}
	antfarm.Run()

}
