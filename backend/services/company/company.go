package company

import (
	"context"
	"jobboard/backend/db"
	jsonutil "jobboard/backend/util/json"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	apiPathRoot = "/companies"
)

type Company struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Siren   string `json:"siren"`
	LogoURL string `json:"logoURL"`
}

func DecodeCompany(data jsonutil.Value) (company Company, err error) {
	company.Name, err = data.Get("name").String()
	if err != nil {
		return
	}
	company.Siren, err = data.Get("siren").String()
	if err != nil {
		return
	}
	company.LogoURL, err = data.Get("logoURL").String()
	if err != nil {
		return
	}
	return
}

func (c Company) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":       c.ID,
		"name":     c.Name,
		"siren":    c.Siren,
		"logo_url": c.LogoURL,
	}
}

type CompanyPage []Company

func (c *CompanyPage) Len() int {
	return len(*c)
}

func (c *CompanyPage) GetCursor(idx int) any {
	return (*c)[idx].ID
}

func (c *CompanyPage) Slice(start, end int) {
	*c = (*c)[start:end]
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
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	page, err := db.DecodePage(jsonVal)
	if err != nil {
		return err
	}
	companies, cursors, err := s.getAll(c.Context(), page)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, companies))
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

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) add(ctx context.Context, company *Company) error {
	return s.db.QueryRow(
		ctx, &company.ID,
		"INSERT INTO companies VALUES(DEFAULT, @name, @siren, @logo_url) RETURNING id",
		company.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Company, error) {
	var ret Company
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM companies WHERE id = $1", id)
	return ret, err
}

func (s service) getAll(ctx context.Context, page db.Page) ([]Company, db.Cursors, error) {
	var ret CompanyPage
	cursors, err := s.db.QueryPage(ctx, &ret, "SELECT * FROM companies", "id", page)
	return ret, cursors, err
}

func (s service) update(ctx context.Context, company Company) error {
	return s.db.Exec(
		ctx,
		"UPDATE companies SET name = @name, logo_url = @logo_url WHERE id = @id",
		company.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM companies WHERE id = $1", id)
}
