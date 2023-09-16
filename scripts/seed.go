package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/yuriykis/hotel-reservation/api"
	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/db/fixtures"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// err := godotenv.Load("../.env")
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	var (
		ctx         = context.Background()
		mongodburi  = os.Getenv("MONGODB_URI")
		mongodbname = os.Getenv("MONGO_DBNAME")
	)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodburi))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(mongodbname)
	if err := client.Database(mongodbname).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
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

	for i := 0; i < 100; i++ {
		name := fmt.Sprintf("Hotel %d", i)
		location := fmt.Sprintf("location %d", i)
		fixtures.AddHotel(store, name, location, rand.Intn(5)+1, nil)
	}
}
