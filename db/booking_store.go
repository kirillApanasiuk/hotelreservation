package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelreservation/types"
	"os"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error)
	GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error)
	GetBookingById(ctx context.Context, id string) (*types.Booking, error)
	UpdateBooking(ctx context.Context, id string, m bson.M) error
}

func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = s.coll.UpdateByID(ctx, oid, update)
	return err
}

func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	var booked []*types.Booking
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cur.All(ctx, &booked); err != nil {
		return nil, err
	}
	return booked, nil
}

func (s *MongoBookingStore) GetBookingById(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.coll.FindOne(
		ctx, bson.M{
			"_id": oid,
		},
	).Decode(&booking); err != nil {
		return nil, err
	}
	return &booking, nil
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {

	res, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.Id = res.InsertedID.(primitive.ObjectID)
	return booking, nil
}

type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
	BookingStore
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection("bookings"),
	}
}
