package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"hotelreservation/types"
	"os"
)

const userColl = "users"

type Dropper interface {
	Drop(ctx context.Context) error
}

type UserStore interface {
	Dropper
	GetUserById(context.Context, string) (*types.User, error)
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(ctx context.Context, user *types.User) (*types.User, error)
	DeleteUser(ctx context.Context, id string) error
	UpdateUser(ctx context.Context, filter bson.M, valuesToBeUpdated types.UpdateUserParams) error
	GetUserByEmail(ctx context.Context, email string) (*types.User, error)
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {
	param := bson.M{"email": email}
	var user types.User
	err := s.coll.FindOne(ctx, param).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

type MongoUserStore struct {
	client *mongo.Client
	dbname string
	coll   *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	return &MongoUserStore{
		client: client,
		coll:   client.Database(dbname).Collection(userColl),
	}
}

func (s *MongoUserStore) GetUserById(ctx context.Context, userId string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	var user types.User
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	res, err := s.coll.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
func (s *MongoUserStore) Drop(ctx context.Context) error {
	fmt.Println("---dropping user collection")
	return s.coll.Drop(ctx)
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, ID string) error {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}

	//	TODO: maybe it's a good idea to handle if we did not delete any user. Maybe log it or something
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter bson.M, params types.UpdateUserParams) error {
	update := bson.D{
		{
			"$set", params.ToBSON(),
		},
	}
	_, err := s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}
