package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type Config struct {
	Port int
}

func New(config Config) *fiber.App {
	server := fiber.New()
	go func() {
		err := server.Listen(fmt.Sprintf(":%d", config.Port))
		if err != nil {
			panic(err)
		}
	}()
	return server
}
