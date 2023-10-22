package server

import (
	"errors"
	"fmt"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	jsonutil "jobboard/backend/util/json"
	"strings"
	"unicode"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type Config struct {
	Port int
	Logs bool
}

func New(config Config) *fiber.App {
	server := fiber.New(
		fiber.Config{
			ErrorHandler: errorHandler,
		},
	)
	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	if config.Logs {
		server.Use(logger.New())
	}
	return server
}

func errorHandler(ctx *fiber.Ctx, err error) error {
	code, message := parseError(err)
	err = ctx.Status(code).SendString(message)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return nil
}

func parseError(err error) (int, string) {
	for err != nil {
		switch err {
		case auth.ErrInvalidPassword, auth.ErrInvalidToken:
			return fiber.StatusUnauthorized, capitalize(err.Error())
		case auth.ErrPasswordTooShort, db.ErrInvalidCursor:
			return fiber.StatusUnprocessableEntity, capitalize(err.Error())
		}

		switch e := err.(type) {
		case *fiber.Error:
			return e.Code, e.Message
		case jsonutil.ErrInvalidField, jsonutil.ErrMissingField:
			return fiber.StatusUnprocessableEntity, capitalize(e.Error())
		case *pgconn.PgError:
			constraintParts := strings.Split(e.ConstraintName, "_")
			column := strings.Join(constraintParts[1:len(constraintParts)-1], "_")

			var message string
			switch e.Code {
			case pgerrcode.IntegrityConstraintViolation, pgerrcode.CheckViolation:
				message = fmt.Sprintf("Invalid %s", column)
			case pgerrcode.NotNullViolation:
				message = fmt.Sprintf("%s should not be null", column)
			case pgerrcode.ForeignKeyViolation:
				message = fmt.Sprintf("This %s doesn't exists", column)
			case pgerrcode.UniqueViolation:
				message = fmt.Sprintf("This %s is already taken", column)
			}

			return fiber.StatusUnprocessableEntity, message
		}

		err = errors.Unwrap(err)
	}

	return fiber.StatusInternalServerError, "Internal server error"
}

func capitalize(str string) string {
	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
