package auth

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationScheme    = "basic"
)

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
		headerValue := c.Get(authorizationHeaderKey)
		headerParams := strings.Split(headerValue, " ")
		if len(headerParams) != 2 ||
			strings.ToLower(headerParams[0]) != authorizationScheme ||
			headerParams[0] == "" {

			return unauthorized(c)
		}
		role, err := a.store.GetRole(c.Context(), headerParams[1])
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
