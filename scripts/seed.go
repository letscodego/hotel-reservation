package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        = context.Background()
	client     *mongo.Client
	hotelStore db.HotelStore
	roomStore  db.RoomStore
	userStore  db.UserStore
)

func seedUser(fname, lname, email string) {
	user, err := types.NewUserFromParams(
		types.CreateUserParams{
			FirstName: fname,
			LastName:  lname,
			Email:     email,
			Password:  "1234567",
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	userStore.CreateUser(context.TODO(), user)
}

func seedHotel(name, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []types.Room{
		{
			Type:  types.DeluxeRoomType,
			Size:  "large",
			Price: 89.99,
		},
		{
			Type:  types.SeaSideRoomType,
			Size:  "small",
			Price: 79.99,
		},
		{
			Type:  types.DoubleRoomType,
			Size:  "normal",
			Price: 69.99,
		},
		{
			Type:  types.SingleRoomType,
			Size:  "small",
			Price: 59.99,
		},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID
		insertedRoom, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
		// insertedHotel.Rooms = append(insertedHotel.Rooms, insertedRoom.ID)
		// hotelStore.UpdateHotel(ctx, insertedHotel.ID.Hex(), types.UpdateHotelParams{
		// 	Rooms: insertedHotel.Rooms,
		// })
		fmt.Println(insertedRoom)
	}
}

func main() {
	seedHotel("California", "USA", 5)
	seedHotel("Blue", "UK", 5)
	seedUser("alfredo", "diangelo", "alfredo@diangelo.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}
