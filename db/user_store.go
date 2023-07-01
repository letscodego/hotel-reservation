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
	GetUsers(context.Context) ([]*types.User, error)
	CreateUser(context.Context, *types.User) (*types.User, error)
	DeleteUser(context.Context, string) error
	UpdateUser(context.Context, string, types.UpdateUserParams) (int64, error)
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
		return nil, fmt.Errorf("no document was found with id: %s", id)
	}
	if err != nil {
		log.Fatal(err)
	}
	return &dbuser, nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	users := []*types.User{}
	if err = cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *MongoUserStore) CreateUser(ctx context.Context, user *types.User) (*types.User, error) {
	cur, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = cur.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (s *MongoUserStore) DeleteUser(ctx context.Context, id string) error {
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

func (s *MongoUserStore) UpdateUser(ctx context.Context, filter string, params types.UpdateUserParams) (int64, error) {
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
	fmt.Println(res)
	//TODO: to handle if the user is not updated, maybe log it, or???
	return res.ModifiedCount, nil
}

type PostgresUserStore struct {
}
