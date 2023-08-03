package api

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelreservation/db"
	"hotelreservation/types"
	"net/http"
	"time"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
	store *db.Store
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}

	return nil
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
func (h *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}

	return c.JSON(rooms)
}

func (h *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
	var params BookRoomParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	if err := params.validate(); err != nil {
		return err
	}
	roomOId, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}

	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(
			genericResponse{
				Type: "error",
				Msg:  "internal server error",
			},
		)
	}

	ok, err = h.isRoomAvailableForBooking(h.store.Booking, c.Context(), roomOId, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(
			genericResponse{
				Type: "error",
				Msg:  fmt.Sprintf("room %s already booked", c.Params("id")),
			},
		)
	}

	booking := &types.Booking{
		UserId:     user.ID,
		RoomId:     roomOId,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	res, err := h.store.Booking.InsertBooking(c.Context(), booking)
	if err != nil {
		return err
	}
	return c.JSON(res)
}

func (h *RoomHandler) isRoomAvailableForBooking(
	bs db.BookingStore, ctx context.Context, roomId primitive.ObjectID, params BookRoomParams,
) (bool, error) {
	where := bson.M{
		"roomId": roomId,
		"fromDate": bson.M{
			"$gte": params.FromDate,
		},
		"tillDate": bson.M{
			"$lte": params.TillDate,
		},
	}
	bookings, err := bs.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}

	ok := len(bookings) == 0
	return ok, nil
}
