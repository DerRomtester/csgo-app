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

// Gets a Crosshair from a CSGO-DemoFile and returns a map of that demo
func GetCrosshair(demo string) map[uint32][]string {
	players_crosshair := make(map[uint32][]string)
	f, err := os.Open(demo)
	fmt.Println("Parse Start: ", demo)
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
	fmt.Println("Parse End: ", demo)

	return players_crosshair
}

// Returns the Crosshair from all the demos and returs a map
func ReturnCrosshair(demos []string) map[uint32][]string {
	players_crosshair := make(map[uint32][]string)

	for i := 0; i < len(demos); i++ {
		players_crosshair = GetCrosshair(demos[i])
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
	demos := GetDemo("/home/stefan/development/go/csgo/demos/")

	crosshairs := ReturnCrosshair(demos)
	for steamid, data := range crosshairs {
		fmt.Printf("Steamid \"%v\" Player \"%v\" Crosshair \"%v\"\n", steamid, data[0], data[1])
	}

}
