package user

import (
	"context"
	"jobboard/backend/auth"
	"jobboard/backend/db"
	jsonutil "jobboard/backend/util/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

const (
	apiPathRoot = "/users"
)

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Phone       string    `json:"phone"`
	DateOfBirth time.Time `json:"dateOfBirthUTC"`
}

func DecodeUser(data jsonutil.Value) (user User, err error) {
	user.Email, err = data.Get("email").String()
	if err != nil {
		return
	}
	user.Name, err = data.Get("name").String()
	if err != nil {
		return
	}
	user.Surname, err = data.Get("surname").String()
	if err != nil {
		return
	}
	user.Phone, err = data.Get("phone").String()
	if err != nil {
		return
	}
	dateOfBirthUTC, err := data.Get("dateOfBirthUTC").String()
	if err != nil {
		return
	}
	user.DateOfBirth, err = time.Parse(time.DateOnly, dateOfBirthUTC)
	if err != nil {
		return
	}
	return
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

type Account struct {
	UserID       int    `json:"-"`
	PasswordHash string `json:"-"`
	AuthToken    string `json:"-"`
	Role         Role   `json:"role"`
}

func DecodeAccount(data jsonutil.Value) (account Account, err error) {
	role, err := data.Get("role").String()
	account.Role = Role(role)
	if err != nil {
		return
	}
	password, err := data.Get("password").String()
	if err != nil {
		return
	}
	account.PasswordHash, err = auth.HashPassword(password)
	if err != nil {
		return
	}
	return
}

func (a Account) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"user_id":       a.UserID,
		"password_hash": a.PasswordHash,
		"role":          a.Role,
	}
}

type UserAccount struct {
	User
	Account
}

type UserAccountPage []UserAccount

func (u *UserAccountPage) Len() int {
	return len(*u)
}

func (u *UserAccountPage) GetCursor(idx int) any {
	return (*u)[idx].User.ID
}

func (u *UserAccountPage) Slice(start, end int) {
	*u = (*u)[start:end]
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

type tokenWrap struct {
	Token string
}

type authStore struct {
	db db.DB
}

func NewAuthStore(db db.DB) auth.Store {
	return authStore{db: db}
}

type Service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB, adminAuthorizer fiber.Handler) Service {
	Service := Service{db: db}

	server.Post(apiPathRoot, Service.addHandler)
	server.Get(apiPathRoot+"/:id<int>", adminAuthorizer, Service.getHandler)
	server.Get(apiPathRoot, adminAuthorizer, Service.getAllHandler)
	server.Put(apiPathRoot+"/:id<int>", adminAuthorizer, Service.updateHandler)
	server.Put(apiPathRoot+"/password/:id<int>", adminAuthorizer, Service.updatePasswordHandler)
	server.Delete(apiPathRoot+"/:id<int>", adminAuthorizer, Service.deleteHandler)

	server.Get(apiPathRoot+"/me", Service.getMeHandler)
	server.Put(apiPathRoot+"/me", Service.updateMeHandler)
	server.Put(apiPathRoot+"/password/me", Service.updateMyPasswordHandler)
	server.Delete(apiPathRoot+"/me", Service.deleteMeHandler)
	server.Post(apiPathRoot+"/login", Service.loginHandler)

	return Service
}

// :TODO think about transaction
func (s Service) addHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	user, err := DecodeUser(jsonVal)
	if err != nil {
		return err
	}
	account, err := DecodeAccount(jsonVal)
	if err != nil {
		return err
	}

	err = s.add(c.Context(), &user)
	if err != nil {
		return err
	}
	account.UserID = user.ID
	err = s.addAccount(c.Context(), account)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusCreated)
}

func (s Service) getHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}

	userAccount, err := s.get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(userAccount)
}

func (s Service) getAllHandler(c *fiber.Ctx) error {
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	page, err := db.DecodePage(jsonVal)
	if err != nil {
		return err
	}
	userAccounts, cursors, err := s.getAll(c.Context(), page)
	if err != nil {
		return err
	}
	return c.JSON(db.NewCursorWrap(cursors, userAccounts))
}

func (s Service) updateHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	user, err := DecodeUser(jsonVal)
	if err != nil {
		return err
	}
	role, err := jsonVal.Get("role").String()
	if err != nil {
		return err
	}

	user.ID = id
	err = s.update(c.Context(), user)
	if err != nil {
		return err
	}
	err = s.updateRole(c.Context(), id, Role(role))
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Service) updatePasswordHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	password, err := jsonVal.Get("password").String()
	if err != nil {
		return err
	}

	token, err := s.updatePassword(c.Context(), id, password)
	if err != nil {
		return err
	}
	return c.JSON(tokenWrap{Token: token})
}

func (s Service) deleteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return err
	}
	err = s.delete(c.Context(), id)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Service) getMeHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	ret, err := s.GetByToken(c.Context(), token)
	if err != nil {
		return err
	}
	return c.JSON(ret)
}

func (s Service) updateMeHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	user, err := DecodeUser(jsonVal)
	if err != nil {
		return err
	}

	err = s.updateByToken(c.Context(), token, user)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Service) updateMyPasswordHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	jsonVal, err := jsonutil.Parse(c.Body())
	if err != nil {
		return err
	}
	password, err := jsonVal.Get("password").String()
	if err != nil {
		return err
	}

	token, err = s.updatePasswordByToken(c.Context(), token, password)
	if err != nil {
		return err
	}
	err = c.JSON(tokenWrap{Token: token})
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Service) deleteMeHandler(c *fiber.Ctx) error {
	token, err := auth.TokenFromContext(c)
	if err != nil {
		return err
	}
	err = s.deleteByToken(c.Context(), token)
	if err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (s Service) loginHandler(c *fiber.Ctx) error {
	email := c.Query("email")
	password := c.Query("password")
	account, err := s.getAccountByEmail(c.Context(), email)
	if err != nil {
		return err
	}

	err = auth.ValidatePassword(password, account.PasswordHash)
	if err != nil {
		return err
	}
	return c.JSON(tokenWrap{Token: account.AuthToken})
}

func (s Service) add(ctx context.Context, user *User) error {
	return s.db.QueryRow(
		ctx, &user.ID, `
		INSERT INTO users
		VALUES (DEFAULT, @email, @name, @surname, @phone, @date_of_birth)
		RETURNING id`,
		user.toArgs(),
	)
}

func (s Service) addAccount(ctx context.Context, account Account) error {
	return s.db.Exec(
		ctx,
		"INSERT INTO accounts VALUES (@user_id, @password_hash, DEFAULT, @role)",
		account.toArgs(),
	)
}

func (s Service) get(ctx context.Context, id int) (UserAccount, error) {
	var ret UserAccount
	err := s.db.QueryRow(
		ctx, &ret, `
		SELECT * FROM users
		JOIN accounts ON users.id = accounts.user_id
		WHERE id = $1`,
		id,
	)
	return ret, err
}

func (s Service) getAll(ctx context.Context, page db.Page) ([]UserAccount, db.Cursors, error) {
	var ret UserAccountPage
	cursors, err := s.db.QueryPage(
		ctx, &ret,
		"SELECT * FROM users JOIN accounts on users.id = accounts.user_id",
		"id", page,
	)
	return ret, cursors, err
}

func (s Service) update(ctx context.Context, user User) error {
	return s.db.Exec(
		ctx, `
		UPDATE users
		SET email = @email, name = @name, surname = @surname, phone = @phone,
			date_of_birth = @date_of_birth
		WHERE id = @id`,
		user.toArgs(),
	)
}

func (s Service) updateRole(ctx context.Context, id int, role Role) error {
	return s.db.Exec(
		ctx,
		"UPDATE accounts SET role = $2 WHERE user_id = $1",
		id, role,
	)
}

func (s Service) updatePassword(ctx context.Context, id int, password string) (token string, err error) {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	return token, s.db.QueryRow(
		ctx, &token, `
		UPDATE accounts
		SET password_hash = $2, auth_token = DEFAULT
		WHERE user_id = $1
		RETURNING auth_token`,
		id, passwordHash,
	)
}

func (s Service) delete(ctx context.Context, id int) error {
	return s.db.Exec(ctx, "DELETE FROM users WHERE id = $1", id)
}

func (s Service) GetByToken(ctx context.Context, token string) (User, error) {
	var user User
	err := s.db.QueryRow(
		ctx, &user, `
		SELECT users.* FROM users
		JOIN accounts on users.id = accounts.user_id
		WHERE auth_token = $1`,
		token,
	)
	return user, err
}

func (s Service) updateByToken(ctx context.Context, token string, user User) error {
	args := user.toArgs()
	args["auth_token"] = token

	return s.db.Exec(
		ctx, `
		UPDATE users
		SET email = @email, name = @name, surname = @surname, phone = @phone,
			date_of_birth = @date_of_birth
		FROM accounts
		WHERE users.id = accounts.user_id AND accounts.auth_token = @auth_token`,
		args,
	)
}

func (s Service) updatePasswordByToken(
	ctx context.Context, token, password string,
) (string, error) {

	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return "", err
	}

	return token, s.db.QueryRow(
		ctx, &token, `
		UPDATE accounts
		SET password_hash = $2, auth_token = DEFAULT
		WHERE auth_token = $1
		RETURNING auth_token`,
		token, passwordHash,
	)
}

func (s Service) deleteByToken(ctx context.Context, token string) error {
	return s.db.Exec(ctx, `
		DELETE FROM users
		USING accounts
		WHERE users.id = accounts.user_id AND auth_token = $1`,
		token,
	)
}

func (s Service) getAccountByEmail(ctx context.Context, email string) (Account, error) {
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

func (s authStore) GetRole(ctx context.Context, token string) (string, error) {
	var ret string
	err := s.db.QueryRow(ctx, &ret, "SELECT role FROM accounts WHERE auth_token = $1", token)
	return ret, err
}
