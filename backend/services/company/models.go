package company

import (
	jsonutil "jobboard/backend/utils/json"

	"github.com/jackc/pgx/v5"
)

type Company struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Siren   string `json:"siren"`
	LogoURL string `json:"logoURL"`
}

func DecodeCompany(data jsonutil.Value) (company Company, err error) {
	company.Name, err = data.Get("name").String()
	if err != nil {
		return
	}
	company.Siren, err = data.Get("siren").String()
	if err != nil {
		return
	}
	company.LogoURL, err = data.Get("logoURL").String()
	if err != nil {
		return
	}
	return
}

func (c Company) toArgs() pgx.NamedArgs {
	return pgx.NamedArgs{
		"id":       c.ID,
		"name":     c.Name,
		"siren":    c.Siren,
		"logo_url": c.LogoURL,
	}
}

type CompanyPage []Company

func (c *CompanyPage) Len() int {
	return len(*c)
}

func (c *CompanyPage) GetCursor(idx int) any {
	return (*c)[idx].ID
}

func (c *CompanyPage) Slice(start, end int) {
	*c = (*c)[start:end]
}
