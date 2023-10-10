package user

import (
	"context"
	"encoding/json"
	"jobboard/backend/db"
	"jobboard/backend/models"

	"github.com/gofiber/fiber/v2"
)

const (
	apiPathRoot = "/users"
	tableName   = "users"
)

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB) {
	service := service{db: db}

	server.Get(apiPathRoot, service.getAllHandler)
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	users, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	userJson, _ := json.Marshal(users)
	c.Write(userJson)
	return nil
}

func (s service) getAll(ctx context.Context) ([]models.User, error) {
	return db.GetAll[[]models.User](ctx, s.db, tableName)
}
