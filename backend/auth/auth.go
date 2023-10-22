package auth

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slices"
)

const (
	authorizationHeaderKey = "Authorization"
	authorizationScheme    = "basic"

	minPasswordLen = 8
)

var (
	ErrInvalidToken     = fmt.Errorf("invalid or missing authorization header")
	ErrInvalidPassword  = fmt.Errorf("incorrect password")
	ErrPasswordTooShort = fmt.Errorf(
		"passwords must have at least %d characters", minPasswordLen,
	)
)

func HashPassword(password string) (string, error) {
	if utf8.RuneCountInString(password) < minPasswordLen {
		return "", ErrPasswordTooShort
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash), err
}

func ValidatePassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return ErrInvalidPassword
	}
	return nil
}

func TokenFromContext(c *fiber.Ctx) (string, error) {
	headerValue := c.Get(authorizationHeaderKey)
	headerParams := strings.Split(headerValue, " ")
	if len(headerParams) != 2 ||
		strings.ToLower(headerParams[0]) != authorizationScheme ||
		headerParams[0] == "" {

		return "", ErrInvalidToken
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
			return err
		}
		role, _ := a.store.GetRole(c.Context(), token)

		if !slices.Contains(allowedRoles, role) {
			return fiber.ErrUnauthorized
		}

		return c.Next()
	}
}
