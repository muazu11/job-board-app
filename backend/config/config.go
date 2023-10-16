package config

import (
	"jobboard/backend/db"
	"jobboard/backend/server"
	"jobboard/backend/services"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	Server   server.Config
	DB       db.Config
	Services services.Config
}

func New() Config {
	var config Config
	loader := aconfig.LoaderFor(&config, aconfig.Config{})
	err := loader.Load()
	if err != nil {
		panic(err)
	}
	return config
}
