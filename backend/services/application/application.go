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

type ApplicationPage []Application

func (a *ApplicationPage) Len() int {
	return len(*a)
}

func (a *ApplicationPage) GetCursor(idx int) any {
	return (*a)[idx].ID
}

func (a *ApplicationPage) Slice(start, end int) {
	*a = (*a)[start:end]
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
	page := db.PageFromContext(c, db.IntColumn)
	applications, cursors, err := s.getAll(c.Context(), page)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, applications))
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
	return s.db.QueryRow(
		ctx, &application.ID, `
		INSERT INTO applications
		VALUES (DEFAULT, @advertisement_id, @applicant_id, @message, DEFAULT)
		RETURNING id`,
		application.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Application, error) {
	var ret Application
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM applications WHERE id = $1", id)
	return ret, err
}

func (s service) getAll(ctx context.Context, page db.Page) ([]Application, db.Cursors, error) {
	var ret ApplicationPage
	cursors, err := s.db.QueryPage(ctx, &ret, "SELECT * FROM applications", "id", page)
	return ret, cursors, err
}

func (s service) update(ctx context.Context, application Application) error {
	return s.db.Exec(ctx, `
		UPDATE applications
		SET advertisement_id = @advertisement_id, applicant_id = @applicant_id,
			message = @message
		WHERE id = @id`,
		application.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM applications WHERE id = $1", id)
}
