package main

import (
	"fmt"
	"os"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func main() {
	f, err := os.Open("nuke-fpl.dem")

	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := dem.NewParser(f)
	players_crosshair := make(map[string]string)

	p.RegisterEventHandler(func(start events.MatchStart) {
		for _, pl := range p.GameState().Participants().All() {
			players_crosshair[pl.Name] = pl.CrosshairCode()
		}
	})

	// Parse to end
	err = p.ParseToEnd()
	if err != nil {
		panic(err)
	}

	for player, crosshair := range players_crosshair {
		if crosshair != "" {
			fmt.Println("Player:", player, "   Crosshair:", crosshair)
		}

	}
}
