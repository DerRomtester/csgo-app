package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	dem "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v2/pkg/demoinfocs/events"
)

// Collects All Demos from the Demopath and checks if the demo
func GetDemo(demopath string) []string {
	var demos []string
	err := filepath.Walk(demopath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".dem") {
			demos = append(demos, path)
		}
		return nil
	})
	checkError("Cannot obtain Demos from Folder", err)

	return demos
}

// Returns the Crosshair from the Demos that get into the function
func GetCrosshair(demos []string) map[uint32][]string {
	players_crosshair := make(map[uint32][]string)
	for _, demo := range demos {
		f, err := os.Open(demo)
		checkError("Cannot parse Demo", err)
		defer f.Close()

		p := dem.NewParser(f)

		p.RegisterEventHandler(func(start events.MatchStart) {
			for _, pl := range p.GameState().Participants().All() {
				if pl.CrosshairCode() == "" {
				} else {
					players_crosshair[pl.SteamID32()] = []string{pl.Name, pl.CrosshairCode()}
				}
			}
		})

		// Parse to end
		err = p.ParseToEnd()
		checkError("Error while Parsing end", err)

	}
	return players_crosshair
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func main() {

	/*
		For Live Going!

			if len(os.Args) != 2 {
			os.Exit(1)
			}
			demos := GetDemo(os.Args[1])
	*/
	demos := GetDemo("C:/Demo/99dmg/Season 18")

	crosshairs := GetCrosshair(demos)
	for steamid, data := range crosshairs {
		fmt.Printf("Steamid \"%v\" Player \"%v\" Crosshair \"%v\"\n", steamid, data[0], data[1])
	}

}
