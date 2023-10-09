package models

import "time"

type Advertisement struct {
	ID          int
	Title       string
	Description string
	CompanyId   int
	Wage        float32
	Address     string
	ZipCode     string
	City        string
	WorkTime    time.Duration
}

type Company struct {
	ID      int
	Name    string
	LogoURL string
}

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

type Application struct {
	ID              int
	AdvertisementID int
	ApplicantID     int
	Message         string
	CreatedAt       time.Time
}
