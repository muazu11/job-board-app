package application

import (
	jsonutil "jobboard/backend/utils/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type Application struct {
	ID              int       `json:"id"`
	AdvertisementID int       `json:"advertismentID"`
	ApplicantID     int       `json:"applicantID"`
	Message         string    `json:"message"`
	CreatedAt       time.Time `json:"createdAt"`
}

func DecodeApplication(data jsonutil.Value) (application Application, err error) {
	application.AdvertisementID, err = data.Get("advertisementID").Int()
	if err != nil {
		return
	}
	application.ApplicantID, err = data.Get("applicantID").Int()
	if err != nil {
		return
	}
	application.Message, err = data.Get("message").String()
	if err != nil {
		return
	}
	return
}

func (a Application) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":               a.ID,
		"advertisement_id": a.AdvertisementID,
		"applicant_id":     a.ApplicantID,
		"message":          a.Message,
		"created_at":       a.CreatedAt,
	}
}

type ApplicationPage []Application

func (a *ApplicationPage) Len() int {
	return len(*a)
}

func (a *ApplicationPage) GetCursor(idx int) any {
	return (*a)[idx].ID
}

func (a *ApplicationPage) Slice(start, end int) {
	*a = (*a)[start:end]
}
