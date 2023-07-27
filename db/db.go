package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DBNAME     = "hotel-reservation"
	TestDBNAME = "hotel-reservation-test"
	MongoDBuri = "mongodb://localhost:27017"
)

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}

func ToObjectId(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}
