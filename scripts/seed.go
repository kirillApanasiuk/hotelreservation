package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/api"
	"hotelreservation/db"
	"hotelreservation/db/fixtures"
	"log"
	"math/rand"
	"os"
	"time"
)

var database *db.Store
var client *mongo.Client
var roomStore db.RoomStore
var hotelStore db.HotelStore
var bookingStore db.BookingStore
var ctx = context.Background()
var userStore db.UserStore
var mongoDbName string

func main() {
	if err := client.Database(mongoDbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)

	database = &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}

	admin := fixtures.AddUser(database, "Kirill", "Apanasiuk", true)
	fmt.Printf("admin --> %v", api.CreateTokenFromUser(admin))

	user := fixtures.AddUser(database, "Ighor", "Schepau", false)
	fmt.Printf("user --> %v", api.CreateTokenFromUser(user))

	hotel := fixtures.AddHotel(database, "some hotel", "bermuda", 5, nil)

	room := fixtures.AddRoom(database, "small", true, 89.99, hotel.ID)
	fmt.Println(room)
	booking := fixtures.AddBooking(database, user.ID, room.ID, time.Now(), time.Now().AddDate(0, 1, 2), 4)
	fmt.Println(booking)

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("random hotel name %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(database, name, location, rand.Intn(5)+1, nil)
	}
}

func init() {
	var (
		err        error
		mongoDbUri = os.Getenv("MONGO_DB_URL")
	)

	mongoDbName = os.Getenv("MONGO_DB_NAME")
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbUri))
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Database(mongoDbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	bookingStore = db.NewMongoBookingStore(client)
	userStore = db.NewMongoUserStore(client)
}
