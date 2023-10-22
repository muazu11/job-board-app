package advertisement

import (
	"jobboard/backend/services/company"
	jsonutil "jobboard/backend/utils/json"
	"time"

	"github.com/jackc/pgx/v5"
)

type Advertisement struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	CompanyID   int           `json:"companyID"`
	Wage        float64       `json:"wage"`
	Address     string        `json:"address"`
	ZipCode     string        `json:"zipCode"`
	City        string        `json:"city"`
	WorkTime    time.Duration `json:"workTimeNs"`
}

func DecodeAdvertisement(data jsonutil.Value) (advertisement Advertisement, err error) {
	advertisement.Title, err = data.Get("title").String()
	if err != nil {
		return
	}
	advertisement.Description, err = data.Get("description").String()
	if err != nil {
		return
	}
	advertisement.CompanyID, err = data.Get("companyID").Int()
	if err != nil {
		return
	}
	advertisement.Wage, err = data.Get("wage").Float()
	if err != nil {
		return
	}
	advertisement.Address, err = data.Get("address").String()
	if err != nil {
		return
	}
	advertisement.ZipCode, err = data.Get("zipCode").String()
	if err != nil {
		return
	}
	advertisement.City, err = data.Get("city").String()
	if err != nil {
		return
	}
	workTimeNs, err := data.Get("workTimeNs").Int()
	advertisement.WorkTime = time.Duration(workTimeNs)
	if err != nil {
		return
	}
	return
}

func (a Advertisement) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":          a.ID,
		"title":       a.Title,
		"description": a.Description,
		"company_id":  a.CompanyID,
		"wage":        a.Wage,
		"address":     a.Address,
		"zip_code":    a.ZipCode,
		"city":        a.City,
		"work_time":   a.WorkTime,
	}
}

type AdvertisementPage []Advertisement

func (a *AdvertisementPage) Len() int {
	return len(*a)
}

func (a *AdvertisementPage) GetCursor(idx int) any {
	return (*a)[idx].ID
}

func (a *AdvertisementPage) Slice(start, end int) {
	*a = (*a)[start:end]
}

type CompanyAdvertisement struct {
	Advertisement
	Applied         bool
	company.Company `json:"Company"`
}

type CompanyAdvertisementPage []CompanyAdvertisement

func (a *CompanyAdvertisementPage) Len() int {
	return len(*a)
}

func (a *CompanyAdvertisementPage) GetCursor(idx int) any {
	return (*a)[idx].Advertisement.ID
}

func (a *CompanyAdvertisementPage) Slice(start, end int) {
	*a = (*a)[start:end]
}
