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

func WriteCrosshairDB(player []interface{}, c *mongo.Client) {
	collection := c.Database("crosshair-db").Collection("Crosshairdata")
	insertResult, err := collection.InsertMany(context.TODO(), player)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted Multiple Documents ", insertResult.InsertedIDs)
}

func ReadCrosshairCollection(client *mongo.Client) []demo.PlayersData {
	collection := client.Database("crosshair-db").Collection("Crosshairdata")
	var crosshairs []demo.PlayersData
	cur, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.TODO())
	for cur.Next(context.TODO()) {
		var crosshair demo.PlayersData
		err := cur.Decode(&crosshair)
		if err != nil {
			log.Fatal(err)
		}
		crosshairs = append(crosshairs, crosshair)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	return crosshairs
}

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:example@localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	defer cancel()
	return client
}
