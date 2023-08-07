package db

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const MongoDBNameEnvName = "MONGO_DB_NAME"

type PaginationFilter struct {
	Limit int64
	Page  int64
}

type HotelFilter struct {
	PaginationFilter
	Rating int64
}

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}

func ToObjectId(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}
