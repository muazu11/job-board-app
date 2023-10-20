package advertisement

import (
	"context"
	"errors"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	"jobboard/backend/services/company"
	"jobboard/backend/services/user"
	jsonutil "jobboard/backend/util/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/sanggonlee/gosq"
)

const (
	apiPathRoot = "/advertisements"
)

type Advertisement struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CompanyID   int           `json:"companyID"`
	Wage        float64       `json:"wage"`
	Address     string        `json:"address"`
	ZipCode     string        `json:"zipCode"`
	City        string        `json:"city"`
	WorkTime    time.Duration `json:"workTimeNs"`
}

func DecodeAdvertisement(data jsonutil.Value) (advertisement Advertisement, err error) {
	advertisement.Title, err = data.Get("title").String()
	if err != nil {
		return
	}
	advertisement.Description, err = data.Get("description").String()
	if err != nil {
		return
	}
	advertisement.CompanyID, err = data.Get("companyID").Int()
	if err != nil {
		return
	}
	advertisement.Wage, err = data.Get("wage").Float()
	if err != nil {
		return
	}
	advertisement.Address, err = data.Get("address").String()
	if err != nil {
		return
	}
	advertisement.ZipCode, err = data.Get("zipCode").String()
	if err != nil {
		return
	}
	advertisement.City, err = data.Get("city").String()
	if err != nil {
		return
	}
	workTimeNs, err := data.Get("workTimeNs").Int()
	advertisement.WorkTime = time.Duration(workTimeNs)
	if err != nil {
		return
	}
	return
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

type AdvertisementPage []Advertisement

func (a *AdvertisementPage) Len() int {
	return len(*a)
}

func (a *AdvertisementPage) GetCursor(idx int) any {
	return (*a)[idx].ID
}

func (a *AdvertisementPage) Slice(start, end int) {
	*a = (*a)[start:end]
}

type CompanyAdvertisement struct {
	Advertisement
	Applied         bool
	company.Company `json:"Company"`
}

type CompanyAdvertisementPage []CompanyAdvertisement

func (a *CompanyAdvertisementPage) Len() int {
	return len(*a)
}

func (a *CompanyAdvertisementPage) GetCursor(idx int) any {
	return (*a)[idx].Advertisement.ID
}

func (a *CompanyAdvertisementPage) Slice(start, end int) {
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

	server.Get(apiPathRoot+"/with_detail", service.getAllWithDetailHandler)
}

func (s service) addHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	advertisement, err := DecodeAdvertisement(jsonVal)
	if err != nil {
		return err
	}
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

func (s service) getAllHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	page, err := db.DecodePage(jsonVal)
	if err != nil {
		return err
	}
	advertisements, cursors, err := s.getAll(c.Context(), page)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, advertisements))
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
	advertisement, err := DecodeAdvertisement(jsonVal)
	if err != nil {
		return err
	}
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

func (s service) getAllWithDetailHandler(c *fiber.Ctx) error {
	var user user.User
	token, err := auth.TokenFromContext(c)
	if err == nil {
		user, err = s.user.GetByToken(c.Context(), token)
		if err != nil {
			return err
		}
	} else if !errors.Is(err, auth.ErrInvalidToken) {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	page, err := db.DecodePage(jsonVal)
	if err != nil {
		return err
	}

	advertisements, cursors, err := s.getAllWithDetail(c.Context(), page, user.ID)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, advertisements))
}

func (s service) add(ctx context.Context, advertisement *Advertisement) error {
	return s.db.QueryRow(
		ctx, &advertisement.ID, `
		INSERT INTO advertisements
		VALUES (
			DEFAULT, @title, @description, @company_id, @wage, @address, @zip_code,
			@city, @work_time
		)
		RETURNING id`,
		advertisement.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (Advertisement, error) {
	var ret Advertisement
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM advertisements WHERE id = $1", id)
	return ret, err
}

func (s service) getAll(
	ctx context.Context, page db.Page,
) ([]Advertisement, db.Cursors, error) {

	var ret AdvertisementPage
	cursors, err := s.db.QueryPage(ctx, &ret, "SELECT * FROM advertisements", "id", page)
	return ret, cursors, err
}

func (s service) getAllWithDetail(
	ctx context.Context, page db.Page, userID int,
) ([]CompanyAdvertisement, db.Cursors, error) {

	withApplied := userID != 0
	query, err := gosq.Compile(`
		SELECT advertisements.*, companies.name, companies.logo_url, companies.siren
		{{ [if] .WithApplied [then] , EXISTS (
		  SELECT 1 FROM applications
		  WHERE applicant_id = $3 AND advertisement_id = advertisements.id
		) AS applied }}
		FROM advertisements 
		JOIN companies on companies.id = advertisements.company_id`,
		map[string]any{"WithApplied": withApplied},
	)

	var args []any
	if withApplied {
		args = append(args, userID)
	}
	var ret CompanyAdvertisementPage
	cursors, err := s.db.QueryPage(ctx, &ret, query, "advertisements.id", page, args...)
	return ret, cursors, err
}

func (s service) update(ctx context.Context, advertisement Advertisement) error {
	return s.db.Exec(
		ctx, `
		UPDATE advertisements
		SET title = @title, description = @description, company_id = @company_id, wage = @wage,
			address = @address, zip_code = @zip_code, city = @city, work_time = @work_time
		WHERE id = @id`,
		advertisement.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM advertisements WHERE id = $1", id)
}
