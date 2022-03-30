package main

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Crosshairdata struct {
	Steamid   uint32
	Player    string
	Crosshair string
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("crosshair-db").Collection("Crosshairdata")
	player1 := Crosshairdata{153400465, "ZywOo", "CSGO-ESKFZ-vA79a-GnBoU-VONmR-XdTOG"}

	insertResult, err := collection.InsertOne(context.TODO(), player1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a Single Document: ", insertResult.InsertedID)

}
