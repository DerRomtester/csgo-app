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

type PlayersData struct {
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
	checkError("Cannot obtain Demos from Folder", err)

	return demos
}

func GetAllCrosshairs(demos []string) []interface{} {
	var allPlayers []PlayersData
	for democounter := 0; democounter < len(demos); democounter++ {
		demofile, err := os.Open(demos[democounter])
		fmt.Println("Parse Start: ", demos[democounter], "counter Count: ", democounter)
		checkError("Cannot parse Demo", err)
		defer demofile.Close()

		p := dem.NewParser(demofile)
		p.RegisterEventHandler(func(start events.AnnouncementWinPanelMatch) {
			for _, player := range p.GameState().Participants().All() {
				if player.CrosshairCode() == "" {
				} else {
					demopath := demofile.Name()
					demoname := demopath[strings.LastIndex(demopath, "/")+1:]
					player_info := PlayersData{
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

		err = p.ParseToEnd()
		checkError("Error while Parsing end", err)
		fmt.Println("Parse End: ", demos[democounter])
	}
	playerdata := make([]interface{}, len(allPlayers))
	for i := range allPlayers {
		playerdata[i] = allPlayers[i]
	}
	return playerdata
}

func (p PlayersData) GetCrosshair() string {
	return p.Crosshaircode
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
