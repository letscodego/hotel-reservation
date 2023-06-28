package main

import (
	"flag"

	"github.com/gofiber/fiber/v2"
	"github.com/lets-goo/hotel-reservation/api/v1"
)

func main() {
	listenAdd := flag.String("listenAdd", ":5000", "The listen address of the API server")
	app := fiber.New()
	apiv1 := app.Group("api/v1")
	apiv1.Get("/user", api.HandleGetUsers)
	apiv1.Get("/user/:id", api.HandleGetUser)
	app.Listen(*listenAdd)
}
