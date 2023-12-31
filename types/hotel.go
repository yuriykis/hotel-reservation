package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name"          json:"name"`
	Location string               `bson:"location"      json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms"         json:"rooms"`
	Rating   int                  `bson:"rating"        json:"rating"`
}

func NewHotel(name, location string, rooms []primitive.ObjectID, rating int) *Hotel {
	return &Hotel{
		Name:     name,
		Location: location,
		Rooms:    rooms,
		Rating:   rating,
	}
}

type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string             `bson:"size"          json:"size"`
	Seaside bool               `bson:"seaside"       json:"seaside"`
	Price   float64            `bson:"price"         json:"price"`
	HotelID primitive.ObjectID `bson:"hotelID"       json:"hotelID"`
}

func NewRoom(price float64, seaside bool, hotelID primitive.ObjectID, size string) *Room {
	return &Room{
		Price:   price,
		Seaside: seaside,
		Size:    size,
		HotelID: hotelID,
	}
}
