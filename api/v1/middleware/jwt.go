package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(c *fiber.Ctx) error {
	token, ok := c.GetReqHeaders()["X-Api-Token"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	claims, err := validateJWTToken(token)
	if err != nil {
		return err
	}
	expirationTime, err := claims.GetExpirationTime()
	if err != nil {
		return err
	}
	if time.Now().After(expirationTime.Time) {
		fmt.Println("failed to parse JWT token", err)
		return fmt.Errorf("token is expired")
	}
	return c.Next()
}

func validateJWTToken(tokenString string) (jwt.MapClaims, error) {
	token, err := parsToken(tokenString)
	if err != nil {
		fmt.Println("failed to parse JWT token ", err)
		return nil, fmt.Errorf("unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return claims, nil
}

func parsToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		err := os.Setenv("JWT_SECRET", "my_secret_key")
		if err != nil {
			fmt.Println("failed to set env", err)
			return nil, fmt.Errorf("unauthorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("failed to parse JWT token, token is invalid")
		return nil, fmt.Errorf("unauthorized")
	}
	return token, nil
}
