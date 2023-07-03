package types

import (
	"fmt"

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
	Rating   int                  `bson:"rating" json:"rating"`
}

type UpdateHotelParams struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Rating   int    `json:"rating"`
}

func NewHotelFromParams(hotel Hotel) (*Hotel, error) {
	return &Hotel{
		Name:     hotel.Name,
		Location: hotel.Location,
		Rooms:    []primitive.ObjectID{},
		Rating:   hotel.Rating,
	}, nil
}

func (h Hotel) Validate() map[string]string {
	errors := map[string]string{}
	if len(h.Name) < minNameLen {
		errors["name"] = fmt.Sprintf("name length should be at least %d characters", minNameLen)
	}
	if len(h.Location) < minLocationLen {
		errors["location"] = fmt.Sprintf("location length should be at least %d characters", minLocationLen)
	}
	if h.Rating < 0 && h.Rating > 5 {
		errors["rating"] = "rating should be at between 1 and 5 "
	}
	return errors
}

func (params UpdateHotelParams) ToBSON() bson.M {
	m := bson.M{}
	if len(params.Name) >= minNameLen {
		m["name"] = params.Name
	}
	if len(params.Location) >= minLocationLen {
		m["location"] = params.Location
	}
	if params.Rating >= 0 && params.Rating <= 5 {
		m["rating"] = params.Rating
	}
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
