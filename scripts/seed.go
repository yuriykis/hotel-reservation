package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yuriykis/hotel-reservation/api"
	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/db/fixtures"
	"github.com/yuriykis/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client       *mongo.Client
	roomStore    db.RoomStore
	hotelStore   db.HotelStore
	userStore    db.UserStore
	bookingStore db.BookingStore
	ctx          = context.Background()
)

func seedUser(isAdmin bool, firstName, lastName, email, password string) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "password",
	})
	user.IsAdmin = isAdmin
	if err != nil {
		log.Fatal(err)
	}
	insertedUser, err := userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}

func seedBooking(
	roomID, userID primitive.ObjectID,
	numPersons int,
	fromDate, tillDate time.Time,
) *types.Booking {
	booking := types.NewBooking(roomID, userID, numPersons, fromDate, tillDate)
	insertedBooking, err := bookingStore.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("booking: ", insertedBooking)
	return insertedBooking
}

func seedHotel(name, location string, rating int) *types.Hotel {
	hotel := types.NewHotel(name, location, []primitive.ObjectID{}, rating)
	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func seedRoom(
	price float64,
	ss bool,
	hotelID primitive.ObjectID,
	size string,
) *types.Room {
	room := types.NewRoom(price, ss, hotelID, size)
	insertedRoom, err := roomStore.InsertRoom(ctx, room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)

	store := &db.Store{
		User:    userStore,
		Hotel:   hotelStore,
		Room:    roomStore,
		Booking: bookingStore,
	}
	user := fixtures.AddUser(store, "james", "foo", false)
	fmt.Println("james -> ", api.CreateTokenFromUser(user))
	admin := fixtures.AddUser(store, "admin", "admin", true)
	fmt.Println("admin -> ", api.CreateTokenFromUser(admin))

	hotel := fixtures.AddHotel(store, "Hilton", "Bermuda", 5, nil)
	fmt.Println("hotel: ", hotel)

	room := fixtures.AddRoom(store, 88.4, true, hotel.ID, "large")
	fmt.Println("room: ", room)

	booking := fixtures.AddBooking(
		store,
		room.ID,
		user.ID,
		2,
		time.Now(),
		time.Now().AddDate(0, 0, 1),
	)
	fmt.Println("booking -> ", booking)
}

func init() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
	bookingStore = db.NewMongoBookingStore(client)
}
