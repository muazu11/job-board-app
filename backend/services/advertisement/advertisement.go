package advertisement

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"jobboard/backend/db"
	"jobboard/backend/services/company"
	"time"
)

const (
	apiPathRoot = "/advertisements"
)

type Advertisement struct {
	ID          int
	Title       string
	Description string
	CompanyID   int
	Wage        float64
	Address     string
	ZipCode     string
	City        string
	WorkTime    time.Duration `json:"WorkTimeNs"`
}

type CompanyAdvertisement struct {
	Advertisement
	company.Company `json:"Company"`
}

func (a Advertisement) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":          a.ID,
		"title":       a.Title,
		"description": a.Description,
		"company_id":  a.CompanyID,
		"wage":        a.Wage,
		"address":     a.Address,
		"zip_code":    a.ZipCode,
		"city":        a.City,
		"work_time":   a.WorkTime,
	}
}

func advertisementFromContext(c *fiber.Ctx) Advertisement {
	return Advertisement{
		Title:       c.Query("title"),
		Description: c.Query("description"),
		CompanyID:   c.QueryInt("company_id"),
		Wage:        c.QueryFloat("wage"),
		Address:     c.Query("address"),
		ZipCode:     c.Query("zip_code"),
		City:        c.Query("city"),
		WorkTime:    time.Duration(c.QueryInt("work_time_ns")),
	}
}

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB, adminAuthorizer fiber.Handler) {
	service := service{db: db}

	server.Post(apiPathRoot, adminAuthorizer, service.addHandler)
	server.Get(apiPathRoot+"/:id", adminAuthorizer, service.getHandler)
	server.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	server.Put(apiPathRoot+"/:id", adminAuthorizer, service.updateHandler)
	server.Delete(apiPathRoot+"/:id", adminAuthorizer, service.deleteHandler)

	server.Get(apiPathRoot+"WithCompany", service.showHandler)
}

func (s service) addHandler(c *fiber.Ctx) error {
	advertisement := advertisementFromContext(c)
	return s.add(c.Context(), &advertisement)
}

func (s service) getHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	advertisement, err := s.get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(advertisement)
}

func (s service) showHandler(c *fiber.Ctx) error {
	advertisements, err := s.show(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(advertisements)
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	advertisements, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(advertisements)
}

func (s service) updateHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	advertisement := advertisementFromContext(c)
	advertisement.ID = id
	return s.update(c.Context(), advertisement)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) add(ctx context.Context, advertisement *Advertisement) error {
	return s.db.QueryOne(
		ctx, &advertisement.ID, `
		INSERT INTO advertisements
		VALUES (
			DEFAULT, @title, @description, @company_id, @wage, @address, @zip_code, @city, @work_time
		)
		RETURNING id`,
		nil, advertisement.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Advertisement, error) {
	var ret Advertisement
	err := s.db.QueryOne(ctx, &ret, "SELECT * FROM advertisements WHERE id = $1", nil, id)
	return ret, err
}

func (s service) show(ctx context.Context) ([]CompanyAdvertisement, error) {
	var ret []CompanyAdvertisement
	err := s.db.Query(ctx, &ret, "SELECT advertisements.*, companies.logo_url, companies.name FROM advertisements JOIN companies on companies.id = advertisements.company_id", nil)
	return ret, err
}

func (s service) getAll(ctx context.Context) ([]Advertisement, error) {
	var ret []Advertisement
	err := s.db.Query(ctx, &ret, "SELECT * FROM advertisements", nil)
	return ret, err
}

func (s service) update(ctx context.Context, advertisement Advertisement) error {
	return s.db.Exec(
		ctx, `
		UPDATE advertisements
		SET title = @title, description = @description, company_id = @company_id, wage = @wage,
			address = @address, zip_code = @zip_code, city = @city, work_time = @work_time
		WHERE id = @id`,
		nil, advertisement.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM advertisements WHERE id = $1", nil, id)
}
