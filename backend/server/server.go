package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Config struct {
	Port int
	Logs bool
}

func New(config Config) *fiber.App {
	server := fiber.New()
	if config.Logs {
		server.Use(logger.New(), cors.New(cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept,Authorization",
		}))
	}
	return server
}
