package api

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/yuriykis/hotel-reservation/db/fixtures"
	"github.com/yuriykis/hotel-reservation/types"
)

func TestUserGetBooking(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		nonAuthUser = fixtures.AddUser(tdb.Store, "jimmy", "foobar", false)
		user        = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel       = fixtures.AddHotel(tdb.Store, "Bar hotel", "a", 5, nil)
		room        = fixtures.AddRoom(tdb.Store, 4.4, true, hotel.ID, "small")

		from    = time.Now()
		till    = from.AddDate(0, 0, 5)
		booking = fixtures.AddBooking(tdb.Store, room.ID, user.ID, 2, from, till)
		app     = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		route = app.Group(
			"/",
			JWTAutentication(tdb.User),
		)
		bookingHandler = NewBookingHandler(tdb.Store)
	)

	route.Get("/:id", bookingHandler.HandleGetBooking)

	req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/%s", booking.ID.Hex()), nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err := app.Test(req)

	var bookingResp *types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookingResp); err != nil {
		t.Fatal(err)
	}
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	if bookingResp.ID != booking.ID {
		t.Fatalf("expected the bookings to be equal, got %v", bookingResp)
	}
	if booking.UserID != user.ID {
		t.Fatalf("expected the bookings to be equal, got %v", bookingResp)
	}

	req = httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == fiber.StatusOK {
		t.Fatalf("expected status code not not to be 200, got %d\n", resp.StatusCode)
	}
}

func TestAdminGetBookings(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	var (
		adminUser = fixtures.AddUser(tdb.Store, "admin", "foo", true)
		user      = fixtures.AddUser(tdb.Store, "james", "foo", false)
		hotel     = fixtures.AddHotel(tdb.Store, "Bar hotel", "a", 5, nil)
		room      = fixtures.AddRoom(tdb.Store, 4.4, true, hotel.ID, "small")

		from    = time.Now()
		till    = from.AddDate(0, 0, 5)
		booking = fixtures.AddBooking(tdb.Store, room.ID, user.ID, 2, from, till)
		app     = fiber.New(fiber.Config{
			ErrorHandler: ErrorHandler,
		})
		route = app.Group(
			"/",
			JWTAutentication(tdb.User),
			AdminAuth,
		)
		bookingHandler = NewBookingHandler(tdb.Store)
	)
	_ = booking

	route.Get("", bookingHandler.HandleGetBookings)

	req := httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(adminUser))

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.StatusCode)
	}
	var bookings []*types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil {
		t.Fatal(err)
	}
	if len(bookings) != 1 {
		t.Fatalf("expected 1 booking, got %d", len(bookings))
	}
	if reflect.DeepEqual(booking, bookings[0]) {
		t.Fatalf("expected the bookings to be equal, got %v", bookings[0])
	}

	// test non admin can't get bookings
	req = httptest.NewRequest(fiber.MethodGet, "/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Api-Token", CreateTokenFromUser(user))

	resp, err = app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != fiber.StatusUnauthorized {
		t.Fatalf("expected status code 401, got %d", resp.StatusCode)
	}
}
