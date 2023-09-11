package db

import (
	"context"

	"github.com/yuriykis/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomColl = "rooms"

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(DBNAME).Collection(roomColl),
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	cur, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var rooms []*types.Room
	// for cur.Next(ctx) {
	// 	var room types.Room
	// 	if err := cur.Decode(&room); err != nil {
	// 		return nil, err
	// 	}
	// 	rooms = append(rooms, &room)
	// }
	if err := cur.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	res, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	roomID, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, err
	}
	room.ID = roomID

	filter := bson.M{"_id": room.HotelID}
	update := bson.M{
		"$push": bson.M{
			"rooms": roomID,
		},
	}
	if err := s.HotelStore.Update(ctx, filter, update); err != nil {
		return nil, err
	}

	return room, nil
}
