package db

import "go.mongodb.org/mongo-driver/bson/primitive"

const DBURI = "mongodb://localhost:27017"
const DBNAME = "hotel_reservertion"
const TestDBNAME = "hotel_reservertion_test"
const userCollection = "users"
const hotelCollection = "hotels"
const roomCollection = "rooms"

func ToObjectID(id string) primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}
	return oid
}

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
