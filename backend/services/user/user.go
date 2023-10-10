package user

import (
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"jobboard/backend/db"
	"strconv"
	"time"
)

type User struct {
	ID          int
	Email       string
	Name        string
	Surname     string
	Phone       string
	DateOfBirth time.Time
}

type Account struct {
	UserID       int
	PasswordHash string
	Role         Role
}
type Role int

const (
	RoleUser Role = iota + 1
	RoleAdmin
)

const (
	apiPathRoot = "/users"
	tableName   = "users"
)

type service struct {
	db db.DB
}

func Init(server *fiber.App, db db.DB) {
	service := service{db: db}
	server.Post(apiPathRoot, service.add)
	server.Get(apiPathRoot, service.getAllHandler)
	server.Get(apiPathRoot+"/:id", service.get)
	server.Delete(apiPathRoot+"/:id", service.delete)
	server.Put(apiPathRoot+"/:id", service.update) //TODO: Think about rename the route to be more clear or not
}
func (s service) update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	userUpdated := User{
		ID:          id,
		Email:       c.Context().Value("email").(string),
		Name:        c.Context().Value("name").(string),
		Surname:     c.Context().Value("surname").(string),
		Phone:       c.Context().Value("phone").(string),
		DateOfBirth: c.Context().Value("dateOfBirth").(time.Time),
	}
	err = s.updateUserById(c.Context(), userUpdated)
	if err != nil {
		return err
	}
	userJson, _ := json.Marshal(map[string]User{"user": userUpdated})
	c.Write(userJson)
	return nil
}
func (s service) get(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	user, err := s.getUserById(c.Context(), id)
	if err != nil {
		return err
	}

	userJson, _ := json.Marshal(user)
	c.Write(userJson)
	return nil
}

func (s service) delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return err
	}
	err = s.deleteUser(c.Context(), id)
	if err != nil {
		return err
	}
	c.Write(nil)
	return nil
}

// :TODO think about transaction
func (s service) add(c *fiber.Ctx) error {
	user := User{
		Email:       c.Context().Value("email").(string),
		Name:        c.Context().Value("name").(string),
		Surname:     c.Context().Value("surname").(string),
		Phone:       c.Context().Value("phone").(string),
		DateOfBirth: c.Context().Value("dateOfBirth").(time.Time),
	}
	err := s.newUser(c.Context(), user)
	if err != nil {
		return err
	}
	password := c.Context().Value("password").(string)
	hashPwd, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	account := Account{
		UserID:       user.ID,
		PasswordHash: string(hashPwd),
		Role:         RoleUser,
	}
	err = s.newAccount(c.Context(), account)
	if err != nil {
		return err
	}
	c.Write(nil)
	return nil
}

func (s service) getAllHandler(c *fiber.Ctx) error {
	users, err := s.getAll(c.Context())
	if err != nil {
		return err
	}
	userJson, _ := json.Marshal(users)
	c.Write(userJson)
	return nil
}

func (s service) getAll(ctx context.Context) ([]User, error) {
	return db.GetAll[User](ctx, s.db, tableName)
}

func (s service) newUser(ctx context.Context, user User) error {
	/*args := map[string]any{"id": user.ID,
	"email":         user.Email,
	"name":          user.Name,
	"surname":       user.Surname,
	"phone":         user.Phone,
	"date_of_birth": user.DateOfBirth}
	*/
	return s.db.Query(ctx, &user.ID, "INSERT INTO users VALUES (.Email, .Name, .Surname, .Phone, .DateOfBirth) RETURNING id", user)
}

func (s service) newAccount(ctx context.Context, account Account) error {
	/*
		args := map[string]any{"user_id": account.UserID,
			"password_hash": account.PasswordHash,
			"role":          account.Role}
	*/
	return s.db.Exec(ctx, "INSERT INTO accounts VALUES (.UserID, .PasswordHash, .role)", account)
}

func (s service) deleteUser(ctx context.Context, id int) error {
	args := map[string]any{"id": id}
	return s.db.Exec(ctx, "DELETE FROM users WHERE id = .id", args)
}

func (s service) getUserById(ctx context.Context, id int) (User, error) {
	var dest User
	args := map[string]any{"id": id}
	err := s.db.Query(ctx, &dest, "SELECT * FROM users WHERE id = .id", args)
	return dest, err
}

func (s service) updateUserById(ctx context.Context, user User) error {
	/*
		args := map[string]any{"id": id,
			"email":         user.Email,
			"name":          user.Name,
			"surname":       user.Surname,
			"phone":         user.Phone,
			"date_of_birth": user.DateOfBirth}
	*/
	return s.db.Exec(ctx, "UPDATE users SET email = .Email, name = .Name, surname = .Surname, phone = .Phone, date_of_birth = .DateOfBirth WHERE id = .id", user)
}
