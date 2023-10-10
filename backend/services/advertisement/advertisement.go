package advertisement

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"jobboard/backend/db"
	"strconv"
	"time"
)

type Advertisement struct {
	ID          int
	Title       string
	Description string
	CompanyID   int
	Wage        float32
	Address     string
	ZipCode     string
	City        string
	WorkTime    time.Duration
}

const (
	apiPathRoot = "/advertisement"
	tableName   = "advertisement"
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
	advertisement := Advertisement{
		ID:          id,
		Title:       c.Context().Value("title").(string),
		Description: c.Context().Value("description").(string),
		CompanyID:   c.Context().Value("company_id").(int),
		Wage:        c.Context().Value("wage").(float32),
		Address:     c.Context().Value("address").(string),
		ZipCode:     c.Context().Value("zip_code").(string),
		City:        c.Context().Value("city").(string),
		WorkTime:    c.Context().Value("work_time").(time.Duration),
	}
	err = s.update(c.Context(), advertisement)
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
	advertisement, err := db.GetById[Advertisement](c.Context(), s.db, tableName, id)
	if err != nil {
		return err
	}

	advertisementJson, _ := json.Marshal(advertisement)
	c.Write(advertisementJson)
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
	advertisement := Advertisement{
		Title:       c.Context().Value("title").(string),
		Description: c.Context().Value("description").(string),
		CompanyID:   c.Context().Value("company_id").(int),
		Wage:        c.Context().Value("wage").(float32),
		Address:     c.Context().Value("address").(string),
		ZipCode:     c.Context().Value("zip_code").(string),
		City:        c.Context().Value("city").(string),
		WorkTime:    c.Context().Value("work_time").(time.Duration),
	}
	err := s.add(c.Context(), advertisement)
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

func (s service) getAll(ctx context.Context) ([]Advertisement, error) {
	return db.GetAll[Advertisement](ctx, s.db, tableName)
}

func (s service) add(ctx context.Context, advertisement Advertisement) error {
	return s.db.Query(ctx, &advertisement.ID, "INSERT INTO advertisements VALUES (.Title, .Description, .CompanyId, .Wage, .Address, .ZipCode, .City, .WorkTime) RETURNING id", advertisement)
}

func (s service) delete(ctx context.Context, id int) error {
	args := map[string]any{"id": id}
	return s.db.Exec(ctx, "DELETE FROM advertisements WHERE id = .id", args)
}

func (s service) get(ctx context.Context, id int) (Advertisement, error) {
	var dest Advertisement
	args := map[string]any{"id": id}
	err := s.db.Query(ctx, &dest, "SELECT * FROM advertisements WHERE id = .id", args)
	return dest, err
}

func (s service) update(ctx context.Context, advertisement Advertisement) error {
	return s.db.Exec(ctx, "UPDATE advertisements SET  title = .Title, description = .Description , company_id = .CompanyID, wage = .Wage, address = .Address, zip_code = .ZipCode, city = .city, work_time = WorkTime WHERE id = .ID", advertisement)
}
