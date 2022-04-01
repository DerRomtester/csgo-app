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

// Collects all demos from the demopath
func GetDemo(demopath string) []string {
	var demos []string
	err := filepath.Walk(demopath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".dem") { // just checks for the .dem extension
			demos = append(demos, path)
		}
		return nil
	})
	checkError("Cannot obtain Demos from Folder", err)

	return demos
}

// Gets a Crosshair from a CSGO-DemoFile and returns a map of that demo to the channel
func GetCrosshairData(demos []string, out chan map[uint32][]string) {
	players_crosshair := make(map[uint32][]string)
	ParseFinished := 0
	for i := 0; i < len(demos); i++ {
		// Goroutine for Parsing Demos
		go func(i int) {
			f, err := os.Open(demos[i])
			fmt.Println("Parse Start: ", demos[i], "i: ", i)
			checkError("Cannot parse Demo", err)
			defer f.Close()

			p := dem.NewParser(f)

			p.RegisterEventHandler(func(start events.MatchStart) {
				for _, pl := range p.GameState().Participants().All() {
					if pl.CrosshairCode() == "" {
					} else {
						players_crosshair[pl.SteamID32()] = []string{pl.Name, pl.CrosshairCode()}
						out <- players_crosshair
					}
				}
			})

			// Parse to end
			err = p.ParseToEnd()
			checkError("Error while Parsing end", err)
			fmt.Println("Parse End: ", demos[i])
			ParseFinished = ParseFinished + 1

			if ParseFinished == len(demos) {
				close(out) // close the channel when all demos are parsed
			}
		}(i)
	}
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

	demo := GetDemo("/home/stefan/development/go/csgo/demos/")
	c := make(chan map[uint32][]string)

	GetCrosshairData(demo, c)
	players_crosshair := make(map[uint32][]string)

	// Make a map with all the crosshair data
	for msg := range c {
		for steamid, data := range msg {
			players_crosshair[steamid] = []string{data[0], data[1]}
		}
	}

	// Prints out the Crosshair Data of the players
	for steamid, data := range players_crosshair {
		fmt.Printf("Steamid \"%v\" Player \"%v\" Crosshair \"%v\"\n", steamid, data[0], data[1])
	}

}
