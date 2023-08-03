package api

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/db"
	"log"
	"testing"
)

const (
	TestDbUri = "mongodb://localhost:27017"
	dbname    = "hotel-reservation-test"
)

type Testdb struct {
	client *mongo.Client
	Store  *db.Store
}

func (tdb *Testdb) teardown(t *testing.T) {
	if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *Testdb {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(TestDbUri))
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
