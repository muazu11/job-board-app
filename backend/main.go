package main

import (
	"jobboard/backend/config"
	"jobboard/backend/db"
	"jobboard/backend/server"
	"jobboard/backend/services/advertisement"
	"jobboard/backend/services/application"
	"jobboard/backend/services/company"
	"jobboard/backend/services/user"
)

func main() {
	config := config.New()
	server := server.New(config.Server)
	db := db.New(config.DB)

	user.Init(server, db)
	advertisement.Init(server, db)
	application.Init(server, db)
	company.Init(server, db)

	<-(chan struct{})(nil)
}
