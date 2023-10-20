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
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	if config.Logs {
		server.Use(logger.New())
	}
	return server
}
