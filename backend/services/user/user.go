package user

import (
	"context"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	apiPathRoot = "/users"
)

type User struct {
	ID          int
	Email       string
	Name        string
	Surname     string
	Phone       string
	DateOfBirth time.Time
}

func (u User) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":            u.ID,
		"email":         u.Email,
		"name":          u.Name,
		"surname":       u.Surname,
		"phone":         u.Phone,
		"date_of_birth": u.DateOfBirth,
	}
}

func userFromContext(c *fiber.Ctx) (User, error) {
	dateOfBirth, err := time.Parse(time.DateOnly, c.Query("date_of_birth_utc"))
	if err != nil {
		return User{}, err
	}
	user := User{
		Email:       c.Query("email"),
		Name:        c.Query("name"),
		Surname:     c.Query("surname"),
		Phone:       c.Query("phone"),
		DateOfBirth: dateOfBirth,
	}
	return user, nil
}

type Account struct {
	UserID       int
	PasswordHash string
	AuthToken    string
	Role         Role
}

func (a Account) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"user_id":       a.UserID,
		"password_hash": a.PasswordHash,
		"role":          a.Role,
	}
}

func accountFromContext(c *fiber.Ctx) (Account, error) {
	password := c.Query("password")
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return Account{}, err
	}
	account := Account{
		PasswordHash: string(passwordHash),
		Role:         RoleUser,
	}
	return account, nil
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

type authStore struct {
	db db.DB
}

func NewAuthStore(db db.DB) auth.Store {
	return authStore{db: db}
}

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB, adminAuthorizer fiber.Handler) {
	service := service{db: db}

	server.Post(apiPathRoot, service.addHandler)
	server.Get(apiPathRoot+"/:id<int>", adminAuthorizer, service.getHandler)
	server.Get(apiPathRoot, adminAuthorizer, service.getAllHandler)
	server.Put(apiPathRoot+"/:id<int>", adminAuthorizer, service.updateHandler)
	server.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, service.deleteHandler)

	server.Post(apiPathRoot+"/login", service.loginHandler)
}

// :TODO think about transaction
func (s service) addHandler(c *fiber.Ctx) error {
	user, err := userFromContext(c)
	if err != nil {
		return err
	}
	err = s.add(c.Context(), &user)
	if err != nil {
		return err
	}

	account, err := accountFromContext(c)
	if err != nil {
		return err
	}
	account.UserID = user.ID
	return s.addAccount(c.Context(), account)
}

func (s service) getHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	user, err := s.get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	users, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}

func (s service) updateHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	user, err := userFromContext(c)
	if err != nil {
		return err
	}
	user.ID = id
	return s.update(c.Context(), user)
}

func (s service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	return s.delete(c.Context(), id)
}

func (s service) loginHandler(c *fiber.Ctx) error {
	email := c.Query("email")
	password := c.Query("password")
	account, err := s.getAccountByEmail(c.Context(), email)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return err
	}

	token := struct{ Token string }{
		Token: account.AuthToken,
	}
	return c.JSON(token)
}

func (s service) add(ctx context.Context, user *User) error {
	return s.db.QueryOne(
		ctx, &user.ID, `
		INSERT INTO users
		VALUES (DEFAULT, @email, @name, @surname, @phone, @date_of_birth)
		RETURNING id`,
		nil, user.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (User, error) {
	var ret User
	err := s.db.QueryOne(ctx, &ret, "SELECT * FROM users WHERE id = $1", nil, id)
	return ret, err
}

func (s service) getAll(ctx context.Context) ([]User, error) {
	var ret []User
	err := s.db.Query(ctx, &ret, "SELECT * FROM users", nil)
	return ret, err
}

func (s service) update(ctx context.Context, user User) error {
	return s.db.Exec(
		ctx, `
		UPDATE users
		SET email = @email, name = @name, surname = @surname, phone = @phone,
			date_of_birth = @date_of_birth
		WHERE id = @id`,
		nil, user.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", nil, id)
}

func (s service) addAccount(ctx context.Context, account Account) error {
	return s.db.Exec(
		ctx,
		"INSERT INTO accounts VALUES (@user_id, @password_hash, DEFAULT, @role)",
		nil, account.toArgs(),
	)
}

func (s service) getAccountByEmail(ctx context.Context, email string) (Account, error) {
	var ret Account
	err := s.db.QueryOne(
		ctx, &ret, `
		SELECT accounts.* FROM accounts
		JOIN users ON accounts.user_id = users.id
		WHERE users.email = $1`,
		nil, email,
	)
	return ret, err
}

func (s authStore) GetRole(ctx context.Context, token string) (string, error) {
	var ret string
	err := s.db.QueryOne(ctx, &ret, "SELECT role FROM accounts WHERE auth_token = $1", nil, token)
	return ret, err
}
