package application

import (
	"context"
	"jobboard/backend/db"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	apiPathRoot = "/applications"
)

type Application struct {
	ID              int
	AdvertisementID int
	ApplicantID     int
	Message         string
	CreatedAt       time.Time
}

func (a Application) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":               a.ID,
		"advertisement_id": a.AdvertisementID,
		"applicant_id":     a.ApplicantID,
		"message":          a.Message,
		"created_at":       a.CreatedAt,
	}
}

func applicationFromContext(c *fiber.Ctx) Application {
	return Application{
		AdvertisementID: c.QueryInt("advertisement_id"),
		ApplicantID:     c.QueryInt("applicant_id"),
		Message:         c.Query("message"),
	}
}

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB, adminAuthorizer fiber.Handler) {
	service := service{db: db}
	server.Post(apiPathRoot, adminAuthorizer, service.addHandler)
	server.Get(apiPathRoot+"/:id<int>", adminAuthorizer, service.getHandler)
	server.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	server.Put(apiPathRoot+"/:id<int>", adminAuthorizer, service.updateHandler)
	server.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, service.deleteHandler)
}

func (s service) addHandler(c *fiber.Ctx) error {
	application := applicationFromContext(c)
	return s.add(c.Context(), application)
}

func (s service) getHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	application, err := s.get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(application)
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	applications, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(applications)
}

func (s service) updateHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	application := applicationFromContext(c)
	application.ID = id
	return s.update(c.Context(), application)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) add(ctx context.Context, application Application) error {
	return s.db.QueryOne(
		ctx, &application.ID, `
		INSERT INTO applications
		VALUES (DEFAULT, @advertisement_id, @applicant_id, @message, DEFAULT)
		RETURNING id`,
		nil, application.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Application, error) {
	var ret Application
	err := s.db.QueryOne(ctx, &ret, "SELECT * FROM applications WHERE id = $1", nil, id)
	return ret, err
}

func (s service) getAll(ctx context.Context) ([]Application, error) {
	var ret []Application
	err := s.db.Query(ctx, &ret, "SELECT * FROM applications", nil)
	return ret, err
}

func (s service) update(ctx context.Context, application Application) error {
	return s.db.Exec(ctx, `
		UPDATE applications
		SET advertisement_id = @advertisement_id, applicant_id = @applicant_id,
			message = @message
		WHERE id = @id`,
		nil, application.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM applications WHERE id = $1", nil, id)
}
