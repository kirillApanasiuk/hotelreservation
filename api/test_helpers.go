package api

import (
	"context"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/db"
	"log"
	"os"
	"testing"
)

var mongoDbTestUri string

type Testdb struct {
	client *mongo.Client
	Store  *db.Store
}

func (tdb *Testdb) teardown(t *testing.T) {
	mongodbName := os.Getenv("MONGO_DB_NAME")
	if err := tdb.client.Database(mongodbName).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *Testdb {
	if err := godotenv.Load("../.env.dev"); err != nil {
		t.Error(err)
	}

	mongoDbTestUri := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDbTestUri))
	if err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &Testdb{
		client: client,
		Store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   db.NewMongoHotelStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
