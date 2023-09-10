package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db"
)

type HotelHandler struct {
	hotelStore db.HotelStore
	roomStore  db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

type HotelQueryParams struct {
	Rooms  bool `query:"rooms"`
	Rating int  `query:"rating"`
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qp HotelQueryParams
	if err := c.QueryParser(&qp); err != nil {
		return err
	}
	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
