package db

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hotelreservation/types"
	"os"
)

type HotelStore interface {
	Insert(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)
	Update(ctx context.Context, m bson.M, update bson.M) error
	GetHotels(ctx context.Context, filter bson.M, hotelFilter HotelFilter) ([]*types.Hotel, error)
	GetHotel(ctx context.Context, filter bson.M) (*types.Hotel, error)
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter bson.M, hotelFilter HotelFilter) (
	[]*types.Hotel, error,
) {

	opts := options.FindOptions{}
	opts.SetSkip((hotelFilter.Page - 1) * hotelFilter.Limit)
	opts.SetLimit(hotelFilter.Limit)

	resp, err := s.coll.Find(
		ctx, &bson.M{
			"rating": hotelFilter.Rating,
		}, &opts,
	)
	if err != nil {
		return nil, err
	}
	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
func (s *MongoHotelStore) GetHotel(ctx context.Context, filter bson.M) (*types.Hotel, error) {
	var hotel types.Hotel
	err := s.coll.FindOne(ctx, filter).Decode(&hotel)
	if err != nil {
		return nil, err
	}

	return &hotel, nil
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}

func (s *MongoHotelStore) Insert(ctx context.Context, Hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, Hotel)
	if err != nil {
		return nil, err
	}

	Hotel.ID = resp.InsertedID.(primitive.ObjectID)
	return Hotel, nil
}

func (s *MongoHotelStore) Update(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
