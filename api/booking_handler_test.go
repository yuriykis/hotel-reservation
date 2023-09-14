package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/yuriykis/hotel-reservation/db/fixtures"
)

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	user := fixtures.AddUser(tdb.Store, "james", "foo", false)
	hotel := fixtures.AddHotel(tdb.Store, "Bar hotel", "a", 5, nil)
	room := fixtures.AddRoom(tdb.Store, 4.4, true, hotel.ID, "small")

	from := time.Now()
	till := from.AddDate(0, 0, 1)
	booking := fixtures.AddBooking(tdb.Store, room.ID, user.ID, 2, from, till)
	fmt.Println("booking -> ", booking)

	// app := fiber.New()
	// store := &db.Store{
	// 	User:    tdb.User,
	// 	Hotel:   tdb.Hotel,
	// 	Room:    tdb.Room,
	// 	Booking: tdb.Booking,
	// }
	// bookingHandler := NewBookingHandler(store)
	// app.Get("/", bookingHandler.HandleGetBookings)

	// token := CreateTokenFromUser(user)
	// req := httptest.NewRequest(http.MethodGet, "/", nil)
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// res, err := app.Test(req)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// if res.StatusCode != http.StatusUnauthorized {
	// 	t.Errorf("expected %d, got %d", http.StatusUnauthorized, res.StatusCode)
	// }
}
