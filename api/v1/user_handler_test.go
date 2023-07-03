package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
	db.UserStore
}

func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setup(t *testing.T) *testdb {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
	if err != nil {
		log.Fatal(err)
	}
	return &testdb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := types.CreateUserParams{
		FirstName: "Michel",
		LastName:  "G Sccot",
		Email:     "m.sccot@dm.com",
		Password:  "123456789.",
	}
	body, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Status)

	var user types.User
	err = json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}
	if params.Email != user.Email {
		t.Errorf("expected email is %s but %s", params.Email, user.Email)
	}
	if params.LastName != user.LastName {
		t.Errorf("expected lastName is %s but %s", params.LastName, user.LastName)
	}
	if params.FirstName != user.FirstName {
		t.Errorf("expected firstName is %s but %s", params.FirstName, user.FirstName)
	}
	if len(user.EncryptedPassword) > 0 {
		t.Error("expecting encrypted password not to be included in json response")
	}
	if len(user.ID) == 0 {
		t.Errorf("expected user id to be set")
	}
}

func TestGetUser(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(testdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)
	app.Get("/hotel/:id", userHandler.HandleGetUser)

	params := types.CreateUserParams{
		FirstName: "Michel",
		LastName:  "G Sccot",
		Email:     "m.sccot@dm.com",
		Password:  "123456789.",
	}
	body, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp.Status)

	var user types.User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("/hotel/" + user.ID.Hex())
	fmt.Println(user.ID.Hex())

	get_user_req := httptest.NewRequest("GET", "/hotel/"+user.ID.Hex(), nil)
	resp, _ = app.Test(get_user_req)

	var user_from_db types.User
	err = json.NewDecoder(resp.Body).Decode(&user_from_db)
	if err != nil {
		t.Error(err)
	}

	if user.Email != user_from_db.Email {
		t.Errorf("expected email is %s but %s", user.Email, user_from_db.Email)
	}
	if user.LastName != user_from_db.LastName {
		t.Errorf("expected lastName is %s but %s", user.LastName, user_from_db.LastName)
	}
	if user.FirstName != user_from_db.FirstName {
		t.Errorf("expected firstName is %s but %s", user.FirstName, user_from_db.FirstName)
	}
	if user.ID != user_from_db.ID {
		t.Errorf("expected firstName is %s but %s", user.ID, user_from_db.ID)
	}
}
