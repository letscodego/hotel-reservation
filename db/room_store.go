package db

import (
	"context"

	"github.com/lets-goo/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	CreateRoom(context.Context, *types.Room) (*types.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection

	HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	coll := client.Database(DBNAME).Collection(roomCollection)
	return &MongoRoomStore{
		client:     client,
		collection: coll,
		HotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	cur, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}

	room.ID = cur.InsertedID.(primitive.ObjectID)
	err = s.HotelStore.UpdateHotelRooms(ctx, room.HotelID.Hex(), room.ID.Hex())
	if err != nil {
		return nil, err
	}
	return room, nil
}
