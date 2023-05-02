package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/DerRomtester/csgo-app/m/v2/internal/demo"
	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
	Player    string  `json:"playername" bson:"playername"`
	SteamID64 uint64  `json:"steamid64" bson:"steamid64"`
	SteamID32 uint32  `json:"steamid32" bson:"steamid32"`
	Matches   []Match `json:"matches" bson:"matches,omitempty"`
}

type Match struct {
	MapName       string `bson:"mapname"`
	PlaybackTime  int    `bson:"playbacktime"`
	Crosshaircode string `bson:"crosshaircode"`
}

type SQL struct {
	Player        string `json:"playername"`
	SteamID64     uint64 `json:"steamid32"`
	SteamID32     uint32 `json:"steamid64"`
	CrosshairCode string `json:"crosshaircode"`
}

const (
	MONGO_URI   = "mongodb://root:example@localhost:27017/"
	PG_USER     = "postgres"
	PG_PASSWORD = "example"
	PG_DATABASE = "csgo-dev"
)

// PostgreSQL Database Interaction
func Pg_ConnectDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", PG_USER, PG_PASSWORD, PG_DATABASE)
	db, err := sql.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to PgSQL")
	return db
}

func Pg_WriteCrosshairDB(players []demo.PlayerInfo, db *sql.DB) {
	var sqlStatement []string
	for i := range players {
		sqlStatement = append(sqlStatement, `INSERT INTO crosshairs(code) VALUES ($1);`)
		sqlStatement = append(sqlStatement, `INSERT INTO players(steamid64, steamid32, name) VALUES ($1, $2, $3);`)
		sqlStatement = append(sqlStatement, `INSERT INTO public.matches(
			"FileStamp", "Protocol", "NetworkProtocol", "Servername", "ClientName", "MapName", "GameDirectory", "PlaybackTime", "PlaybackTicks", "PlaybackFrames", "SignonLength")
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`)

		db.Exec(sqlStatement[0], players[i].Demo.Crosshaircode)
		db.Exec(sqlStatement[1], players[i].SteamID64, players[i].SteamID32, players[i].Playername)
		db.Exec(sqlStatement[2], players[i].Demo.FileStamp, players[i].Demo.Protocol, players[i].Demo.NetworkProtocol, players[i].Demo.Servername, players[i].Demo.ClientName, players[i].Demo.MapName, players[i].Demo.GameDirectory, players[i].Demo.PlaybackTime, players[i].Demo.PlaybackTicks, players[i].Demo.PlaybackFrames, players[i].Demo.SignonLength)
	}
	db.Close()
}

func Pg_ReadCrosshairCollection(db *sql.DB) []SQL {
	fmt.Println("Start Reading Crosshair Collection")
	rows, err := db.Query(`SELECT DISTINCT code, steamid64, steamid32, name from crosshairs
	inner join matches ON matches.crosshairid = crosshairs.id
	inner join players ON matches.playerid = players.id`)
	if err != nil {
		log.Fatal(err)
	}
	var PlayersData []SQL

	for rows.Next() {
		var PlayerData SQL
		rows.Scan(&PlayerData.CrosshairCode, &PlayerData.SteamID64, &PlayerData.SteamID32, &PlayerData.Player)
		PlayersData = append(PlayersData, PlayerData)
	}
	fmt.Println("Return Crosshair Collection")
	defer rows.Close()
	return PlayersData
}

// Mongo Database Interaction
func Mongo_WriteCrosshairDB(players []demo.PlayerInfo, c *mongo.Client) {
	collection := c.Database("crosshair-db").Collection("Crosshairdata")
	for i := range players {
		filter := Player{Player: players[i].Playername, SteamID64: players[i].SteamID64, SteamID32: players[i].SteamID32}
		match := demo.DemoInfo{
			FileStamp:       players[i].Demo.FileStamp,
			Protocol:        players[i].Demo.Protocol,
			NetworkProtocol: players[i].Demo.NetworkProtocol,
			Servername:      players[i].Demo.Servername,
			ClientName:      players[i].Demo.ClientName,
			MapName:         players[i].Demo.MapName,
			GameDirectory:   players[i].Demo.GameDirectory,
			PlaybackTime:    players[i].Demo.PlaybackTime,
			PlaybackTicks:   players[i].Demo.PlaybackTicks,
			PlaybackFrames:  players[i].Demo.PlaybackFrames,
			SignonLength:    players[i].Demo.SignonLength,
			Crosshaircode:   players[i].Demo.Crosshaircode,
		}
		update := bson.D{{"$addToSet", bson.D{{"matches", match}}}}
		opts := options.Update().SetUpsert(true)
		insertResult, err := collection.UpdateOne(context.TODO(), filter, update, opts)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Inserted Document ", insertResult.ModifiedCount)
		fmt.Println("Upserted Document ", insertResult.UpsertedCount)
	}
}

func Mongo_ReadCrosshairCollection(client *mongo.Client) ([]Player, error) {
	fmt.Println("Start Reading Crosshair Collection")
	collection := client.Database("crosshair-db").Collection("Crosshairdata")
	cur, err := collection.Find(context.Background(), bson.D{})
	var CrosshairCollection []Player
	if err != nil {
		return CrosshairCollection, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		var result Player
		if err := cur.Decode(&result); err != nil {
			log.Fatal(err)
		}
		CrosshairCollection = append(CrosshairCollection, result)
	}
	fmt.Println("Return Crosshair Collection")
	return CrosshairCollection, err
}

func Mongo_ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	if err = client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	defer cancel()
	return client
}
