package user

import (
	"jobboard/backend/auth"
	"jobboard/backend/db"
	jsonutil "jobboard/backend/utils/json"
	"time"

	"github.com/jackc/pgx/v5"
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
