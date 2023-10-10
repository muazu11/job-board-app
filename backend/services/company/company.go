package company

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"jobboard/backend/db"
	"strconv"
)

type Company struct {
	ID      int
	Name    string
	LogoURL string
}

const (
	apiPathRoot = "/companies"
	tableName   = "companies"
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
	company := Company{
		ID:      id,
		Name:    c.Context().Value("name").(string),
		LogoURL: c.Context().Value("logoUrl").(string),
	}
	err = s.update(c.Context(), company)
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
	company, err := db.GetById[Company](c.Context(), s.db, tableName, id)
	if err != nil {
		return err
	}

	companyJson, _ := json.Marshal(company)
	c.Write(companyJson)
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
	company := Company{
		Name:    c.Context().Value("name").(string),
		LogoURL: c.Context().Value("logoUrl").(string),
	}
	err := s.add(c.Context(), company)
	if err != nil {
		return err
	}
	return nil
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	companies, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	companyJson, _ := json.Marshal(companies)
	c.Write(companyJson)
	return nil
}

func (s service) getAll(ctx context.Context) ([]Company, error) {
	return db.GetAll[Company](ctx, s.db, tableName)
}

func (s service) add(ctx context.Context, company Company) error {
	return s.db.Query(ctx, &company.ID, "INSERT INTO companies VALUES (.Name, .LogoUrl) RETURNING id", company)
}

func (s service) delete(ctx context.Context, id int) error {
	args := map[string]any{"id": id}
	return s.db.Exec(ctx, "DELETE FROM companies WHERE id = .id", args)
}

func (s service) get(ctx context.Context, id int) (Company, error) {
	var dest Company
	args := map[string]any{"id": id}
	err := s.db.Query(ctx, &dest, "SELECT * FROM companies WHERE id = .id", args)
	return dest, err
}

func (s service) update(ctx context.Context, company Company) error {
	return s.db.Exec(ctx, "UPDATE companies SET  name = .Name, logo_url = .LogoUrl WHERE id = .id", company)
}
