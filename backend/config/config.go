package config

import (
	"jobboard/backend/db"
	"jobboard/backend/server"

	"github.com/cristalhq/aconfig"
)

type Config struct {
	Server server.Config
	DB     db.Config
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
