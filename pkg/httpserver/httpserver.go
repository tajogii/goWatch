package httpserver

import (
	"github.com/gofiber/fiber/v3"
)

type Error struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

type RegisterRoute interface {
	RegisterRoutes(r fiber.Router)
}

func NewHttpServer(handlers ...RegisterRoute) *fiber.App {
	app := fiber.New()
	api := app.Group("/api")

	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	return app
}
