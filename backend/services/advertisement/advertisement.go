package advertisement

import (
	"context"
	"errors"
	"fmt"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	"jobboard/backend/server"
	"jobboard/backend/services/user"
	jsonutil "jobboard/backend/utils/json"

	"github.com/gofiber/fiber/v2"
	"github.com/sanggonlee/gosq"
)

const (
	apiPathRoot = "/advertisements"
)

type service struct {
	db   db.DB
	user user.Service
}

func Init(app *fiber.App, db db.DB, user user.Service, adminAuthorizer fiber.Handler) {
	service := service{db: db, user: user}

	app.Post(apiPathRoot, adminAuthorizer, server.Create, service.addHandler)
	app.Get(apiPathRoot+"/:id<int>", adminAuthorizer, service.getHandler)
	app.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	app.Put(apiPathRoot+"/:id<int>", adminAuthorizer, server.NoContent, service.updateHandler)
	app.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, server.NoContent, service.deleteHandler)

	app.Get(apiPathRoot+"/with_detail", service.getAllWithDetailHandler)
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

func (s service) get(ctx context.Context, id int) (Advertisement, error) {
	var ret Advertisement
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM advertisements WHERE id = $1", id)
	return ret, err
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	pageRef, err := db.PageRefFromContext(c)
	if err != nil {
		return err
	}
	advertisements, cursors, err := s.getAll(c.Context(), pageRef)
	if err != nil {
		return err
	}
	return c.JSON(db.NewPage(cursors, advertisements))
}

func (s service) getAll(
	ctx context.Context, pageRef db.PageRef,
) ([]Advertisement, db.Cursors, error) {

	var ret AdvertisementPage
	cursors, err := s.db.QueryPage(ctx, &ret, "SELECT * FROM advertisements", "id", pageRef)
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
	advertisement, err := DecodeAdvertisement(jsonVal)
	if err != nil {
		return err
	}
	advertisement.ID = id
	return s.update(c.Context(), advertisement)
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

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM advertisements WHERE id = $1", id)
}

func (s service) getAllWithDetailHandler(c *fiber.Ctx) error {
	var user user.User
	token, err := auth.TokenFromContext(c)
	fmt.Println(token, err)
	if err == nil {
		user, err = s.user.GetByToken(c.Context(), token)
		if err != nil {
			return err
		}
	} else if !errors.Is(err, auth.ErrInvalidToken) {
		return err
	}
	pageRef, err := db.PageRefFromContext(c)
	if err != nil {
		return err
	}

	fmt.Println(user.ID)
	advertisements, cursors, err := s.getAllWithDetail(c.Context(), pageRef, user.ID)
	if err != nil {
		return err
	}
	fmt.Println(advertisements[0].Applied)
	return c.JSON(db.NewPage(cursors, advertisements))
}

func (s service) getAllWithDetail(
	ctx context.Context, pageRef db.PageRef, userID int,
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
	cursors, err := s.db.QueryPage(ctx, &ret, query, "advertisements.id", pageRef, args...)
	return ret, cursors, err
}
