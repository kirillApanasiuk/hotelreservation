package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/db"
	"hotelreservation/types"
	"log"
)

var client *mongo.Client
var roomStore db.RoomStore
var hotelStore db.HotelStore
var ctx = context.Background()

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Type:      types.SinglePersonRoomType,
			BasePrice: 99.9,
			Size:      "small",
		}, {
			Type:      types.DeluxeRoomType,
			BasePrice: 199.9,
			Size:      "small",
		}, {
			Type:      types.SeaSideRoomType,
			BasePrice: 150.9,
			Size:      "small",
		},
	}

	insertedHotel, err := hotelStore.Insert(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelId = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}

	fmt.Println(insertedHotel)
}

func main() {
	seedHotel("Santa-monica", "Africa", 4)
	seedHotel("The cozy hotel", "Netherlands", 5)
	seedHotel("Crazy", "Belarus", 6)
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoDBuri))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
