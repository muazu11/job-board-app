package application

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"jobboard/backend/db"
	"strconv"
	"time"
)

type Application struct {
	ID              int
	AdvertisementID int
	ApplicantID     int
	Message         string
	CreatedAt       time.Time
}

const (
	apiPathRoot = "/application"
	tableName   = "application"
)

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB) {
	service := service{db: db}
	server.Post(apiPathRoot, service.addHandler)
	server.Get(apiPathRoot, service.getAllHandler)
	server.Get(apiPathRoot+"/:id", service.getHandler)
	server.Delete(apiPathRoot+"/:id", service.deleteHandler)
	server.Put(apiPathRoot+"/:id", service.updateHandler) //TODO: Think about rename the route to be more clear or not
}
func (s service) updateHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	application := Application{
		ID:              id,
		AdvertisementID: c.Context().Value("advertisement_id").(int),
		ApplicantID:     c.Context().Value("applicant_id").(int),
		Message:         c.Context().Value("message").(string),
		CreatedAt:       c.Context().Value("created_at").(time.Time),
	}
	err = s.update(c.Context(), application)
	if err != nil {
		return err
	}
	c.Write(nil)
	return nil
}
func (s service) getHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	application, err := db.GetById[Application](c.Context(), s.db, tableName, id)
	if err != nil {
		return err
	}

	applicationJson, _ := json.Marshal(application)
	c.Write(applicationJson)
	return nil
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = db.DeleteById(c.Context(), s.db, tableName, id)
	if err != nil {
		return err
	}
	c.Write(nil)
	return nil
}

// :TODO think about transaction
func (s service) addHandler(c *fiber.Ctx) error {
	application := Application{
		AdvertisementID: c.Context().Value("advertisement_id").(int),
		ApplicantID:     c.Context().Value("applicant_id").(int),
		Message:         c.Context().Value("message").(string),
		CreatedAt:       c.Context().Value("created_at").(time.Time),
	}
	err := s.add(c.Context(), application)
	if err != nil {
		return err
	}
	return nil
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	advertisements, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	advertisementsJson, _ := json.Marshal(advertisements)
	c.Write(advertisementsJson)
	return nil
}

func (s service) getAll(ctx context.Context) ([]Application, error) {
	return db.GetAll[Application](ctx, s.db, tableName)
}

func (s service) add(ctx context.Context, application Application) error {
	return s.db.Query(ctx, &application.ID, "INSERT INTO applications VALUES (.AdvertisementID, .ApplicantID, .Message, .CreatedAt) RETURNING id", application)
}

func (s service) delete(ctx context.Context, id int) error {
	args := map[string]any{"id": id}
	return s.db.Exec(ctx, "DELETE FROM applications WHERE id = .id", args)
}

func (s service) get(ctx context.Context, id int) (Application, error) {
	var dest Application
	args := map[string]any{"id": id}
	err := s.db.Query(ctx, &dest, "SELECT * FROM applications WHERE id = .id", args)
	return dest, err
}

func (s service) update(ctx context.Context, application Application) error {
	return s.db.Exec(ctx, "UPDATE applications SET  advertisement_id = .AdvertisementID, applicant_id = .ApplicantID , message = .Message, created_at = .CreatedAt WHERE id = .ID", application)
}
