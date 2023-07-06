package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/db"
	"github.com/lets-goo/hotel-reservation/types"
)

type AuthHandler struct {
	userStore db.UserStore
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), authParams.Email)
	if err != nil {
		return err
	}

	if !types.IsPasswordValid(user.EncryptedPassword, authParams.Password) {
		return fmt.Errorf("invalid credentials")
	}
	return c.JSON(user)
}
