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

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	coll := client.Database(DBNAME).Collection(userCollection)
	return &MongoUserStore{
		client:     client,
		collection: coll,
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var dbuser types.User
	err = s.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&dbuser)
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title %s\n", id)
		return nil, err
	}
	if err != nil {
		log.Fatal(err)
	}
	return &dbuser, nil
}

type PostgresUserStore struct {
}
