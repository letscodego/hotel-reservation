package api

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type AuthResponse struct {
	User  types.User `json:"user"`
	Token string     `json:"token"`
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
	resp := AuthResponse{
		User:  *user,
		Token: createClaimsFromUser(user),
	}
	return c.JSON(resp)

}

func createClaimsFromUser(user *types.User) string {
	now := time.Now()
	expirationTime := now.Add(time.Minute * 25)
	claims := &Claims{
		Username: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	err := os.Setenv("JWT_SECRET", "my_secret_key")
	if err != nil {
		fmt.Println("failed to set env", err)
	}
	secret := os.Getenv("JWT_SECRET")
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)

	}
	return tokenStr
}
