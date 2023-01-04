package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DerRomtester/csgo-app/m/v2/internal/demo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Player struct {
	Player  string  `json:"playername" bson:"playername"`
	SteamId uint64  `json:"steamid" bson:"steamid"`
	Matches []Match `json:"matches" bson:"matches,omitempty"`
}

type Match struct {
	Name      string `json:"name" bson:"name"`
	Crosshair string `json:"crosshair" bson:"crosshair"`
}

const uri = "mongodb://root:example@localhost:27017/"

func WriteCrosshairDB(players []demo.PlayerInfo, c *mongo.Client) {
	collection := c.Database("crosshair-db").Collection("Crosshairdata")
	for i := range players {
		filter := Player{Player: players[i].Playername, SteamId: players[i].Steamid}
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

func ReadCrosshairCollection(client *mongo.Client) ([]Player, error) {
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

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
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
