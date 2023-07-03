package db

import (
	"context"
	"fmt"
	"log"

	"github.com/lets-goo/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	CreateHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotelRooms(context.Context, string, string) error
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	GetHotels(context.Context) ([]*types.Hotel, error)
	DeleteHotel(context.Context, string) error
	UpdateHotel(context.Context, string, types.UpdateHotelParams) (int64, error)
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

func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var dbhotel types.Hotel
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&dbhotel)
	if err == mongo.ErrNoDocuments {
		return nil, fmt.Errorf("no document was found with id: %s", id)
	}
	if err != nil {
		log.Fatal(err)
	}
	return &dbhotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	hotels := []*types.Hotel{}
	if err = cur.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}

func (s *MongoHotelStore) DeleteHotel(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	//TODO: to handle if the user is not deleted, maybe log it, or???
	_, err = s.collection.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}
	return nil
}

func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter string, params types.UpdateHotelParams) (int64, error) {
	oid, err := primitive.ObjectIDFromHex(filter)
	if err != nil {
		return 0, err
	}

	id := bson.M{"_id": oid}
	values := bson.D{{Key: "$set", Value: params.ToBSON()}}
	res, err := s.collection.UpdateOne(ctx, id, values)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
