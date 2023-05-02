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
	chan_playerdata := make(chan []demo.PlayerInfo)
	demo.GetCrosshairs(demos, chan_playerdata)
	for i := range chan_playerdata {
		var Mg_DB = database.Mongo_ConnectDB()
		var Pg_DB = database.Pg_ConnectDB()
		database.Mongo_WriteCrosshairDB(i, Mg_DB)
		database.Pg_WriteCrosshairDB(i, Pg_DB)
	}

}
