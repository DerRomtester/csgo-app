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
	DateTime   string   `bson:"datetime"`
	SteamID64  uint64   `bson:"steamid_64"`
	SteamID32  uint32   `bson:"steamid_32"`
	Playername string   `bson:"playername"`
	Demo       DemoInfo `bson:"demoinfo"`
}

type DemoInfo struct {
	FileStamp       string `bson:"filestamp"`
	Protocol        int    `bson:"protocol"`
	NetworkProtocol int    `bson:"networkprotocol"`
	Servername      string `bson:"servername"`
	ClientName      string `bson:"clientname"`
	MapName         string `bson:"mapname"`
	GameDirectory   string `bson:"gamedirectory"`
	PlaybackTime    int    `bson:"playbacktime"`
	PlaybackTicks   int    `bson:"playbackticks"`
	PlaybackFrames  int    `bson:"playbackframes"`
	SignonLength    int    `bson:"signonlength"`
	Crosshaircode   string `bson:"crosshaircode"`
}

func GetDemos(demopath string) []string {
	var demos []string
	err := filepath.Walk(demopath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() && filepath.Ext(path) == ".dem" {
			demos = append(demos, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return demos
}

func GetCrosshairs(demos []string, out chan []PlayerInfo) {
	var allPlayers []PlayerInfo
	globalcounter := 0
	for i := range demos {
		go func(counter int) {
			startanalyze := time.Now()
			demofile, err := os.Open(demos[counter])
			if err != nil {
				log.Fatal(err)
			}
			defer demofile.Close()
			demopath := demofile.Name()
			demoname := demopath[strings.LastIndex(demopath, "/")+1:]
			fmt.Printf("%s Analyzing Demo: %s \n", time.Now().Format("2006-01-02 15:04:05"), demoname)
			parse := dem.NewParser(demofile)
			header, err := parse.ParseHeader()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(header)
			parse.RegisterEventHandler(func(start events.AnnouncementLastRoundHalf) {
				for _, player := range parse.GameState().Participants().All() {
					// Continue if crosshaircode is empty
					if player.CrosshairCode() == "" {
						fmt.Printf("Could not read Crosshaircode from Player %s SteamID %d Code: %s \n", player.Name, player.SteamID64, player.CrosshairCode())
						continue
					}
					player_info := PlayerInfo{
						DateTime:   time.Now().UTC().String(),
						SteamID64:  player.SteamID64,
						SteamID32:  player.SteamID32(),
						Playername: player.Name,
						Demo: DemoInfo{
							FileStamp:       header.Filestamp,
							Protocol:        header.Protocol,
							NetworkProtocol: header.NetworkProtocol,
							Servername:      header.ServerName,
							ClientName:      header.ClientName,
							MapName:         header.MapName,
							GameDirectory:   header.GameDirectory,
							PlaybackTime:    int(header.PlaybackTime),
							PlaybackTicks:   header.PlaybackTicks,
							PlaybackFrames:  header.PlaybackFrames,
							SignonLength:    header.SignonLength,
							Crosshaircode:   player.CrosshairCode(),
						},
					}
					allPlayers = append(allPlayers, player_info)
				}
			})

			if err = parse.ParseToEnd(); err != nil {
				log.Fatal(err)
			}

			elapsed := time.Since(startanalyze)
			globalcounter++
			fmt.Printf("%s Analyzing Finished: - Duration: %s  \n", time.Now().Format("2006-01-02 15:04:05"), elapsed)
			if globalcounter == len(demos) {
				out <- allPlayers
				close(out)
			}
		}(i)
	}
}
