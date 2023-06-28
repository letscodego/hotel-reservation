package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/types"
)

func HandleGetUsers(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "K1",
		LastName:  "letsgo",
	}

	return c.JSON(user)
}

func HandleGetUser(c *fiber.Ctx) error {
	return c.JSON("user")
}
