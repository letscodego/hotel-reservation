package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
)

func insertTestUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.NewUserFromParams(types.CreateUserParams{
		FirstName: "Alfredo",
		LastName:  "Del",
		Email:     "alfredo@diangelo.com",
		Password:  "1234567",
	})
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticate(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)
	insertedUser := insertTestUser(t, testdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "alfredo@diangelo.com",
		Password: "1234567",
	}
	body, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status of 200 but got %d", resp.StatusCode)
	}
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}
	if authResp.Token == "" {
		t.Fatal("expected jwt token to be presented but it's empty")
	}

	if insertedUser.Email != authResp.User.Email || insertedUser.FirstName != authResp.User.FirstName ||
		insertedUser.LastName != authResp.User.LastName {
		t.Fatal("expected user to be the same as inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	testdb := setup(t)
	defer testdb.teardown(t)
	insertTestUser(t, testdb.UserStore)

	app := fiber.New()
	authHandler := NewAuthHandler(testdb.UserStore)
	app.Post("/", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "alfredo@diangelo.com",
		Password: "123456788888",
	}
	body, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode == http.StatusOK {
		t.Fatalf("expected status <> 200 but got %d", resp.StatusCode)
	}
}
