package types

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	minNameLen     = 5
	minLocationLen = 5
)

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

type UpdateHotelParams struct {
	Name     string               `json:"name"`
	Location string               `json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
}

func (params UpdateHotelParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.Name) >= minNameLen {
		m["name"] = params.Name
	}

	if len(params.Location) >= minLocationLen {
		m["location"] = params.Location
	}
	m["rooms"] = append(params.Rooms, params.Rooms...)

	return m
}

type RoomType int

const (
	SingleRoomType RoomType = iota + 1
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType           `bson:"type" json:"tpye"`
	BasePrice float64            `bson:"basePrice" json:"basePrice"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}
