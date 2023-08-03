package fixtures

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelreservation/db"
	"hotelreservation/types"
	"log"
	"time"
)

func AddUser(store *db.Store, fname, lname string, admin bool) *types.User {
	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			FirstName: fname,
			LastName:  lname,
			Email:     fmt.Sprintf("%s@%s.com", fname, lname),
			Password:  fmt.Sprintf("%s_%s", fname, lname),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	user.IsAdmin = admin
	createdUser, err := store.User.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
	return createdUser
}

func AddBooking(
	store *db.Store, userId, roomId primitive.ObjectID, from, till time.Time, numPersons int,
) *types.Booking {
	booking := &types.Booking{
		UserId:   userId,
		RoomId:   roomId,
		FromDate: from,
		TillDate: till,
	}

	booked, _ := store.Booking.InsertBooking(context.TODO(), booking)
	return booked
}
func AddRoom(store *db.Store, size string, ss bool, price float64, hid primitive.ObjectID) *types.Room {
	room := &types.Room{
		Size:    size,
		Seaside: ss,
		Price:   price,
		HotelId: hid,
	}
	insertedHotel, err := store.Room.InsertRoom(context.TODO(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddHotel(store *db.Store, hotelName, location string, rating int, rooms []primitive.ObjectID) *types.Hotel {
	var roomsIds = rooms
	if rooms == nil {
		roomsIds = []primitive.ObjectID{}
	}

	hotel := &types.Hotel{
		Name:     hotelName,
		Location: location,
		Rooms:    roomsIds,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.Insert(context.TODO(), hotel)
	if err != nil {
		log.Fatal(err)
	}

	return insertedHotel
}
