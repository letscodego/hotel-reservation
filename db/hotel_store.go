package db

import (
	"context"

	"github.com/lets-goo/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotelRooms(context.Context, string, string) error
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	coll := client.Database(DBNAME).Collection(hotelCollection)
	return &MongoHotelStore{
		client:     client,
		collection: coll,
	}
}

func (s *MongoHotelStore) CreateHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	cur, err := s.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}

	hotel.ID = cur.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (s *MongoHotelStore) UpdateHotelRooms(ctx context.Context, filter string, roomId string) error {
	oid, err := primitive.ObjectIDFromHex(filter)
	if err != nil {
		return err
	}
	rid, err := primitive.ObjectIDFromHex(roomId)
	if err != nil {
		return err
	}

	id := bson.M{"_id": oid}
	values := bson.D{{Key: "$push", Value: bson.M{"rooms": rid}}}
	_, err = s.collection.UpdateOne(ctx, id, values)
	if err != nil {
		return err
	}
	//TODO: to handle if the user is not updated, maybe log it, or???
	return nil
}
