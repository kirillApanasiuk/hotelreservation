package api

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"hotelreservation/db"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

type HotelQueryParams struct {
	Rooms []int
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"hotelId": oid}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}

type PaginationParams struct {
	Limit int `query:"limit"`
	Page  int `query:"page"`
}

type HotelFilter struct {
	PaginationParams
	Rating int `json:"rating"`
}

type ResourceResp[T any] struct {
	Results int  `json:"results"`
	Data    []*T `json:"data"`
	Page    int  `json:"page"`
}

func NewResourceResp[T any](data []*T, page int) *ResourceResp[T] {
	return &ResourceResp[T]{
		Results: len(data),
		Data:    data,
		Page:    page,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	queryFilter := new(HotelFilter)
	if err := c.QueryParser(queryFilter); err != nil {
		return ErrBadRequest()
	}

	hotels, err := h.store.Hotel.GetHotels(
		c.Context(), nil, db.HotelFilter{
			PaginationFilter: db.PaginationFilter{
				Limit: int64(queryFilter.Limit),
				Page:  int64(queryFilter.Page),
			},
			Rating: int64(queryFilter.Rating),
		},
	)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(NewResourceResp(hotels, queryFilter.Page))

}
func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {

	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return ErrInvalidId()
	}
	filter := bson.M{"_id": oid}
	rooms, err := h.store.Hotel.GetHotel(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}
