package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddBooking(
	store *db.Store,
	roomID, userID primitive.ObjectID,
	numPersons int,
	fromDate, tillDate time.Time,
) *types.Booking {
	booking := types.NewBooking(roomID, userID, numPersons, fromDate, tillDate)
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(
	store *db.Store,
	price float64,
	seaside bool,
	hotelID primitive.ObjectID,
	size string,
) *types.Room {
	room := types.NewRoom(price, seaside, hotelID, size)
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(
	store *db.Store,
	name, location string,
	rating int,
	rooms []primitive.ObjectID,
) *types.Hotel {
	if rooms == nil {
		rooms = []primitive.ObjectID{}
	}
	hotel := types.NewHotel(name, location, rooms, rating)
	insertedHotel, err := store.Hotel.InsertHotel(context.Background(), hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: fn,
		LastName:  ln,
		Email:     fmt.Sprintf("%s@%s.com", fn, ln),
		Password:  fmt.Sprintf("%s_%s", fn, ln),
	})
	user.IsAdmin = admin
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := store.User.InsertUser(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
