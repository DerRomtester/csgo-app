package demo

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	dem "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs"
	events "github.com/markus-wa/demoinfocs-golang/v3/pkg/demoinfocs/events"
)

type PlayerInfo struct {
	DateTime      string `bson:"datetime"`
	Steamid       uint64 `bson:"steamid"`
	Playername    string `bson:"playername"`
	Crosshaircode string `bson:"crosshaircode"`
	Demoname      string `bson:"demoname"`
}

func GetDemos(demopath string) []string {
	var demos []string
	err := filepath.Walk(demopath, func(path string, info os.FileInfo, err error) error {
		if strings.Contains(path, ".dem") {
			demos = append(demos, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return demos
}

func GetCrosshairs(demos []string) []PlayerInfo {
	var allPlayers []PlayerInfo
	for democounter := 0; democounter < len(demos); democounter++ {
		start := time.Now()
		demofile, err := os.Open(demos[democounter])
		if err != nil {
			log.Fatal(err)
		}
		defer demofile.Close()
		demopath := demofile.Name()
		demoname := demopath[strings.LastIndex(demopath, "/")+1:]
		fmt.Printf("%s Analyzing Demo: %s \n", time.Now().Format(time.RFC850), demoname)

		parse := dem.NewParser(demofile)
		parse.RegisterEventHandler(func(start events.AnnouncementWinPanelMatch) {
			for _, player := range parse.GameState().Participants().All() {
				if player.CrosshairCode() != "" {
					player_info := PlayerInfo{
						DateTime:      time.Now().UTC().String(),
						Steamid:       player.SteamID64,
						Playername:    player.Name,
						Crosshaircode: player.CrosshairCode(),
						Demoname:      demoname,
					}
					allPlayers = append(allPlayers, player_info)
				}
			}
		})

		err = parse.ParseToEnd()

		if err != nil {
			log.Fatal(err)
		}
		elapsed := time.Since(start)
		fmt.Printf("%s Analyzing Finished: - Duration: %s  \n", time.Now().Format(time.RFC850), elapsed)
	}

	return allPlayers
}
