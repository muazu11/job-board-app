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

type UserPage []User

func (u *UserPage) Len() int {
	return len(*u)
}

func (u *UserPage) GetCursor(idx int) any {
	return (*u)[idx].ID
}

func (u *UserPage) Slice(start, end int) {
	*u = (*u)[start:end]
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

type UserAccount struct {
	User
	Account
}

func (a UserAccount) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":            a.ID,
		"email":         a.Email,
		"name":          a.Name,
		"surname":       a.Surname,
		"phone":         a.Phone,
		"date_of_birth": a.DateOfBirth,
		"password_hash": a.PasswordHash,
		"role":          a.Role,
		"auth_token":    a.AuthToken,
	}
}

func userAccountFromContext(c *fiber.Ctx) (userAccount UserAccount, err error) {
	userAccount.User.ID = c.QueryInt("id")
	userAccount.Account.UserID = userAccount.User.ID
	userAccount.User, err = userFromContext(c)
	if err != nil {
		return userAccount, err
	}
	userAccount.Account, err = accountFromContext(c)
	return userAccount, err
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
	server.Get(apiPathRoot+"GetMe", service.getMeHandler)
	server.Put(apiPathRoot+"UpdateMe", service.updateMeHandler)
	server.Delete(apiPathRoot+"DeleteMe", service.deleteMeHandler)
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

func (s service) getMeHandler(c *fiber.Ctx) error {
	var ret User
	ret, err := s.getMe(c.Context(), c)
	if err != nil {
		return err
	}
	return c.JSON(ret)
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
	page := db.PageFromContext(c, db.IntColumn)
	users, cursors, err := s.getAll(c.Context(), page)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, users))
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

func (s service) updateMeHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	userAccount, err := userAccountFromContext(c)
	if err != nil {
		return err
	}
	userAccount.Account.AuthToken = token
	return s.updateMe(c.Context(), userAccount)
}

func (s service) deleteMeHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	return s.deleteMe(c.Context(), c, token)
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
	return s.db.QueryRow(
		ctx, &user.ID, `
		INSERT INTO users
		VALUES (DEFAULT, @email, @name, @surname, @phone, @date_of_birth)
		RETURNING id`,
		user.toArgs(),
	)
}

func (s service) get(ctx context.Context, id int) (User, error) {
	var ret User
	err := s.db.QueryRow(ctx, &ret, "SELECT * FROM users WHERE id = $1", id)
	return ret, err
}

func (s service) getAll(ctx context.Context, page db.Page) ([]User, db.Cursors, error) {
	var ret UserPage
	cursors, err := s.db.QueryPage(
		ctx, &ret,
		// "SELECT * FROM users JOIN accounts on users.id = accounts.user_id",
		"SELECT * FROM users",
		"id", page,
	)
	return ret, cursors, err
}

func (s service) update(ctx context.Context, user User) error {
	return s.db.Exec(
		ctx, `
		UPDATE users
		SET email = @email, name = @name, surname = @surname, phone = @phone,
			date_of_birth = @date_of_birth
		WHERE id = @id`,
		user.toArgs(),
	)
}

func (s service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
}

func (s service) addAccount(ctx context.Context, account Account) error {
	return s.db.Exec(
		ctx,
		"INSERT INTO accounts VALUES (@user_id, @password_hash, DEFAULT, @role)",
		account.toArgs(),
	)
}

func (s service) getAccountByEmail(ctx context.Context, email string) (Account, error) {
	var ret Account
	err := s.db.QueryRow(
		ctx, &ret, `
		SELECT accounts.* FROM accounts
		JOIN users ON accounts.user_id = users.id
		WHERE users.email = $1`,
		email,
	)
	return ret, err
}

func (s service) getUserWithToken(ctx context.Context, c *fiber.Ctx) (User, error) {
	var user User
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return user, err
	}
	err = s.db.QueryOne(ctx, &user, "SELECT users.* FROM users JOIN accounts on users.id = accounts.user_id WHERE auth_token = $1", nil, token)
	return user, err
}

func (s service) getMe(ctx context.Context, c *fiber.Ctx) (User, error) {
	return s.getUserWithToken(ctx, c)
}

func (s service) updateMe(ctx context.Context, userAccount UserAccount) error {
	err := s.db.Exec(ctx, `UPDATE users
		SET email = @email, name = @name, surname = @surname, phone = @phone,
			date_of_birth = @date_of_birth
		FROM accounts
		WHERE users.id = accounts.user_id AND accounts.auth_token = @auth_token`, nil, userAccount.toArgs())

	if err != nil {
		return err
	}
	return s.db.Exec(
		ctx, `
		UPDATE accounts
		SET password_hash = @password_hash
		WHERE auth_token = @auth_token`,
		nil, userAccount.Account.toArgs(),
	)
}

func (s service) deleteMe(ctx context.Context, c *fiber.Ctx, token string) error {
	var user User
	user, err := s.getUserWithToken(ctx, c)
	err = s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", nil, user.ID)
	if err != nil {
		return err
	}
	return s.db.Exec(ctx, "DELETE FROM accounts WHERE auth_token = $1", nil, token)
}
func (s authStore) GetRole(ctx context.Context, token string) (string, error) {
	var ret string
	err := s.db.QueryRow(ctx, &ret, "SELECT role FROM accounts WHERE auth_token = $1", token)
	return ret, err
}
