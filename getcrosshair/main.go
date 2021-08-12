package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

func main() {

	var demos []string
	players_crosshair := make(map[string]string)
	demopath := "C:/Demo/Pros"

	err := filepath.Walk(demopath, func(path string, info os.FileInfo, err error) error {
		demos = append(demos, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, demo := range demos {

		matched := strings.Contains(demo, ".dem")
		if matched {
			f, err := os.Open(demo)

			if err != nil {
				panic(err)
			}
			defer f.Close()

			p := dem.NewParser(f)

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
		}
	}

	for player, crosshair := range players_crosshair {
		if crosshair != "" {
			fmt.Println("Player:", player, "   Crosshair:", crosshair)
		}

	}
}
