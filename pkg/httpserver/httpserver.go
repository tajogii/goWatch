package httpserver

import (
	"github.com/gofiber/fiber/v3"
	logm "github.com/tajogii/goWatch/pkg/logger"
	"go.uber.org/zap"
)

type Error struct {
	Message     string `json:"message"`
	Description string `json:"description"`
}

type RegisterRoute interface {
	RegisterRoutes(r fiber.Router)
}

func NewHttpServer(l *zap.Logger, handlers ...RegisterRoute) *fiber.App {
	app := fiber.New()
	loggerMiddleware := createloggerMiddleware(l)
	app.Use(loggerMiddleware)
	api := app.Group("/api")

	for _, h := range handlers {
		h.RegisterRoutes(api)
	}

	return app
}

func createloggerMiddleware(l *zap.Logger) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := logm.SetLogger(c.Context(), l)
		c.SetContext(ctx)
		return c.Next()
	}

}
