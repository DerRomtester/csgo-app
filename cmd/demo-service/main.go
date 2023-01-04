package main

import (
	"github.com/DerRomtester/csgo-app/m/v2/internal/database"
	"github.com/DerRomtester/csgo-app/m/v2/internal/demo"
)

func main() {

	/*
		For Live Going!

			if len(os.Args) != 2 {
			os.Exit(1)0
			}
			demos := GetDemo(os.Args[1])
	*/
	demos := demo.GetDemos("/home/stefan/csgo-demos")
	playerdata := demo.GetCrosshairs(demos)
	var DB = database.ConnectDB()
	database.WriteCrosshairDB(playerdata, DB)
}
