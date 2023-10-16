package company

import (
	"context"
	"jobboard/backend/db"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	apiPathRoot = "/companies"
)

type Company struct {
	ID      int
	Name    string
	Siren   string
	LogoURL string
}

func (c Company) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":       c.ID,
		"name":     c.Name,
		"siren":    c.Siren,
		"logo_url": c.LogoURL,
	}
}

func companyFromContext(c *fiber.Ctx) Company {
	return Company{
		Name:    c.Query("name"),
		Siren:   c.Query("siren"),
		LogoURL: c.Query("logo_url"),
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
	company := companyFromContext(c)
	return s.add(c.Context(), &company)
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

func (s service) getAllHandler(c *fiber.Ctx) error {
	companies, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(companies)
}

func (s service) updateHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	company := companyFromContext(c)
	company.ID = id
	return s.update(c.Context(), company)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) add(ctx context.Context, company *Company) error {
	return s.db.QueryOne(
		ctx, &company.ID,
		"INSERT INTO companies VALUES(DEFAULT, @name, @siren, @logo_url) RETURNING id",
		nil, company.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Company, error) {
	var ret Company
	err := s.db.QueryOne(ctx, &ret, "SELECT * FROM companies WHERE id = $1", nil, id)
	return ret, err
}

func (s service) getAll(ctx context.Context) ([]Company, error) {
	var ret []Company
	err := s.db.Query(ctx, &ret, "SELECT * FROM companies", nil)
	return ret, err
}

func (s service) update(ctx context.Context, company Company) error {
	return s.db.Exec(
		ctx,
		"UPDATE companies SET name = @name, logo_url = @logo_url WHERE id = @id",
		nil, company.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM companies WHERE id = $1", nil, id)
}
