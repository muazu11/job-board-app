package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationScheme    = "basic"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func ValidatePassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func TokenFromContext(c *fiber.Ctx) (string, error) {
	headerValue := c.Get(authorizationHeaderKey)
	headerParams := strings.Split(headerValue, " ")
	if len(headerParams) != 2 ||
		strings.ToLower(headerParams[0]) != authorizationScheme ||
		headerParams[0] == "" {

		return "", errors.New("invalid or missing authorization header")
	}
	return headerParams[1], nil
}

type Store interface {
	GetRole(ctx context.Context, token string) (string, error)
}

type Auth struct {
	store Store
}

func New(store Store) Auth {
	return Auth{
		store: store,
	}
}

func (a *Auth) NewMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		token, err := TokenFromContext(c)
		if err != nil {
			return unauthorized(c)
		}
		role, err := a.store.GetRole(c.Context(), token)
		if err != nil {
			return unauthorized(c)
		}

		if !slices.Contains(allowedRoles, role) {
			return unauthorized(c)
		}

		return c.Next()
	}
}

func unauthorized(c *fiber.Ctx) error {
	c.Set(fiber.HeaderWWWAuthenticate, "Basic realm=Restricted")
	return c.SendStatus(fiber.StatusUnauthorized)
}
