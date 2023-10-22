package application

import (
	"context"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	"jobboard/backend/services/user"
	jsonutil "jobboard/backend/util/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	apiPathRoot = "/applications"
)

type Application struct {
	ID              int       `json:"id"`
	AdvertisementID int       `json:"advertismentID"`
	ApplicantID     int       `json:"applicantID"`
	Message         string    `json:"message"`
	CreatedAt       time.Time `json:"createdAt"`
}

func DecodeApplication(data jsonutil.Value) (application Application, err error) {
	application.AdvertisementID, err = data.Get("advertisementID").Int()
	if err != nil {
		return
	}
	application.ApplicantID, err = data.Get("applicantID").Int()
	if err != nil {
		return
	}
	application.Message, err = data.Get("message").String()
	if err != nil {
		return
	}
	return
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
	db   db.DB
	user user.Service
}

func Init(server *fiber.App, db db.DB, user user.Service, adminAuthorizer fiber.Handler) {
	service := service{db: db, user: user}
	server.Post(apiPathRoot, adminAuthorizer, service.addHandler)
	server.Get(apiPathRoot+"/:id<int>", adminAuthorizer, service.getHandler)
	server.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	server.Put(apiPathRoot+"/:id<int>", adminAuthorizer, service.updateHandler)
	server.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, service.deleteHandler)
	server.Post(apiPathRoot+"/me", service.applyHandler)
}
func (s service) applyHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	user, err := s.user.GetByToken(c.Context(), token)
	if err != nil {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	application, err := DecodeApplication(jsonVal)
	if err != nil {
		return err
	}
	application.ApplicantID = user.ID

	return s.add(c.Context(), application)
}

func (s service) apply(ctx context.Context, application Application) error {
	return s.db.QueryRow(
		ctx, &application.ID, `
		INSERT INTO applications
		VALUES (DEFAULT, @advertisement_id, @applicant_id, @message, DEFAULT)
		RETURNING id`,
		application.toArgs(),
	)
}

func (s service) addHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	application, err := DecodeApplication(jsonVal)
	if err != nil {
		return err
	}
	err = s.add(c.Context(), application)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
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
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	page, err := db.DecodePage(jsonVal)
	if err != nil {
		return err
	}
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
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	application, err := DecodeApplication(jsonVal)
	if err != nil {
		return err
	}
	application.ID = id
	err = s.update(c.Context(), application)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	err = s.delete(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
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
