package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}

func (h *HotelHandler) HandleGetRooms(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := db.Map{"hotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return ErrResourceNotFound("rooms")
	}
	return c.JSON(rooms)
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.QueryParser(&params); err != nil {
		return ErrBadRequest()
	}
	filter := db.Map{
		"rating": params.Rating,
	}
	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return ErrResourceNotFound("hotels")
	}
	return c.JSON(NewResourceResp(hotels, int(params.Page), len(hotels)))
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

func NewResourceResp(data any, page int, results int) ResourceResp {
	return ResourceResp{
		Results: results,
		Data:    data,
		Page:    page,
	}
}

func (h *HotelHandler) HandleGetHotel(c *fiber.Ctx) error {
	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}
