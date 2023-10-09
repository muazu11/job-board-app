package models

type Advertisement struct {
	ID          int
	Title       string
	Description string
	CompanyId   int
	Wage        float32
	Address     string
	ZipCode     string
	City        string
	WorkTime    int
}

type Company struct {
	ID      int
	Name    string
	LogoUrl string
}

type User struct {
	Id          int
	Email       string
	Name        string
	Surname     string
	Phone       string
	DateOfBirth string
}

type Account struct {
	UserId   int
	Password string
	Role     string
}

type Application struct {
	ID              int
	AdvertisementId int
	ApplicantId     int
	Message         string
	CreatedAt       string
}
