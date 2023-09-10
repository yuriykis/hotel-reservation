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
	ctx        = context.Background()
)

func seedHotel(name, location string) {
	hotel := types.NewHotel("Hilton", "London", []primitive.ObjectID{})
	rooms := []*types.Room{
		types.NewRoom(types.SingleRoomType, 99.9, hotel.ID),
		types.NewRoom(types.DoubleRoomType, 199.9, hotel.ID),
		types.NewRoom(types.SeaSideRoomType, 199.9, hotel.ID),
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
	seedHotel("Hilton", "London")
	seedHotel("The cozy hotel", "The Netherlands")
	seedHotel("Dont die while sleeping", "Paris")
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
}
