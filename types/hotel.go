package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name"          json:"name"`
	Location string               `bson:"location"      json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms"         json:"rooms"`
}

func NewHotel(name, location string, rooms []primitive.ObjectID) *Hotel {
	return &Hotel{
		Name:     name,
		Location: location,
		Rooms:    rooms,
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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type"          json:"type"`
	BasePrice float64            `bson:"basePrice"     json:"basePrice"`
	Price     float64            `bson:"price"         json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID"       json:"hotelID"`
}

func NewRoom(roomType RoomType, basePrice float64, hotelID primitive.ObjectID) *Room {
	return &Room{
		Type:      roomType,
		BasePrice: basePrice,
		Price:     basePrice,
		HotelID:   hotelID,
	}
}
