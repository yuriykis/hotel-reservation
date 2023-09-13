package api

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) validate() error {
	now := time.Now()
	if p.FromDate.Before(now) {
		return fmt.Errorf("fromDate must be in the future")
	}
	if p.TillDate.Before(p.FromDate) {
		return fmt.Errorf("tillDate must be after fromDate")
	}
	if p.NumPersons < 1 {
		return fmt.Errorf("numPersons must be greater than 0")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}

func (h *RoomHandler) HanldeGetRooms(c *fiber.Ctx) error {
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
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  err.Error(),
		})
	}

	roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(genericResp{
			Type: "error",
			Msg:  "user not found",
		})
	}
	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s is already booked for this period", roomID.Hex()),
		})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}
	insertedBooking, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", insertedBooking)
	return c.JSON(insertedBooking)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
	where := bson.M{
		"roomID":   roomID,
		"fromDate": bson.M{"$gte": params.FromDate},
		"tillDate": bson.M{"$lte": params.TillDate},
	}
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	return len(bookings) == 0, nil
}
