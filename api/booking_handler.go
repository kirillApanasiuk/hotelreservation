package api

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"hotelreservation/api/utils"
	"hotelreservation/db"
	"hotelreservation/types"
	"net/http"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{store: store}
}

func (h *BookingHandler) HandleCancelBooking(c *fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}
	user, err := utils.GetAuthUser(c)

	if err != nil {
		return err
	}

	if booking.UserId != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(
			genericResponse{
				Type: "error",
				Msg:  "not authorized",
			},
		)
	}

	err = h.store.Booking.UpdateBooking(
		c.Context(), id, bson.M{
			"canceled": true,
		},
	)
	if err != nil {
		return err
	}
	return c.JSON(
		genericResponse{
			Type: "msg",
			Msg:  "ok",
		},
	)
}

// TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c *fiber.Ctx) error {
	bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrResourceNotFound("bookings")
	}
	return c.JSON(bookings)
}

// TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c *fiber.Ctx) error {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return fmt.Errorf("not authorized")
	}
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingById(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("booking")
	}

	if booking.UserId != user.ID {
		return ErrUnAuthorized()
	}
	return c.JSON(booking)
}
