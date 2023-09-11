package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yuriykis/hotel-reservation/db"
	"github.com/yuriykis/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(firstName, lastName, email, password string) {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "password",
	})
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.InsertUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}
func seedHotel(name, location string, rating int) {
	hotel := types.NewHotel("Hilton", location, []primitive.ObjectID{}, rating)
	rooms := []*types.Room{
		types.NewRoom(types.SingleRoomType, 99.9, hotel.ID, "small"),
		types.NewRoom(types.DoubleRoomType, 199.9, hotel.ID, "normal"),
		types.NewRoom(types.SeaSideRoomType, 199.9, hotel.ID, "kingsize"),
	}
	insertedHotel, err := hotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertedHotel)

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.InsertRoom(ctx, room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(insertedRoom)
	}
}

func main() {
	seedHotel("Hilton", "London", 5)
	seedHotel("The cozy hotel", "The Netherlands", 4)
	seedHotel("Dont die while sleeping", "Paris", 4)
	seedUser("John", "Doe", "james@foo.com", "password")
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
}
