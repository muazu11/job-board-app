package company

import (
	"context"
	"jobboard/backend/db"
	"jobboard/backend/server"
	jsonutil "jobboard/backend/utils/json"

	"github.com/gofiber/fiber/v2"
)

const (
	apiPathRoot = "/companies"
)

type service struct {
	db db.DB
}

func Init(app *fiber.App, db db.DB, adminAuthorizer fiber.Handler) {
	service := service{db: db}
	app.Post(apiPathRoot, adminAuthorizer, server.Create, service.addHandler)
	app.Get(apiPathRoot+"/:id<int>", adminAuthorizer, service.getHandler)
	app.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	app.Put(apiPathRoot+"/:id<int>", adminAuthorizer, server.NoContent, service.updateHandler)
	app.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, server.NoContent, service.deleteHandler)
}

func (s service) addHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	company, err := DecodeCompany(jsonVal)
	if err != nil {
		return err
	}
	return s.add(c.Context(), &company)
}

func (s service) add(ctx context.Context, company *Company) error {
	return s.db.QueryRow(
		ctx, &company.ID,
		"INSERT INTO companies VALUES(DEFAULT, @name, @siren, @logo_url) RETURNING id",
		company.toArgs(),
	)
}

func (s service) getHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	company, err := s.get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(company)
}

func (s service) get(ctx context.Context, id int) (Company, error) {
	var ret Company
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM companies WHERE id = $1", id)
	return ret, err
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	pageRef, err := db.PageRefFromContext(c)
	if err != nil {
		return err
	}
	companies, cursors, err := s.getAll(c.Context(), pageRef)
	if err != nil {
		return err
	}
	return c.JSON(db.NewPage(cursors, companies))
}

func (s service) getAll(ctx context.Context, pageRef db.PageRef) ([]Company, db.Cursors, error) {
	var ret CompanyPage
	cursors, err := s.db.QueryPage(ctx, &ret, "SELECT * FROM companies", "id", pageRef)
	return ret, cursors, err
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
	company, err := DecodeCompany(jsonVal)
	if err != nil {
		return err
	}
	company.ID = id
	return s.update(c.Context(), company)
}

func (s service) update(ctx context.Context, company Company) error {
	return s.db.Exec(
		ctx,
		"UPDATE companies SET name = @name, logo_url = @logo_url WHERE id = @id",
		company.toArgs(),
	)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM companies WHERE id = $1", id)
}
