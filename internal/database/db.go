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
	Name      string `json:"name" bson:"name"`
	Crosshair string `json:"crosshair" bson:"crosshair"`
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

	err = db.Ping()
	if err != nil {
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
		sqlStatement = append(sqlStatement, `INSERT INTO matches(playerid, crosshairid, matchname)VALUES ((select id from players where steamid64 = $1 limit 1), (select id from crosshairs where code = $2 limit 1), $3);`)
		db.Exec(sqlStatement[0], players[i].Crosshaircode)
		db.Exec(sqlStatement[1], players[i].SteamID64, players[i].SteamID32, players[i].Playername)
		db.Exec(sqlStatement[2], players[i].SteamID64, players[i].Crosshaircode, players[i].Demoname)
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
		match := Match{Name: players[i].Demoname, Crosshair: players[i].Crosshaircode}
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
		err := cur.Decode(&result)
		if err != nil {
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
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	defer cancel()
	return client
}
