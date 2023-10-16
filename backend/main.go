package main

import (
	"fmt"
	"jobboard/backend/auth"
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
	db := db.New(config.DB, config.Services)
	auth := auth.New(user.NewAuthStore(db))

	adminAuthorizer := auth.NewMiddleware(user.RoleAdmin.String())

	user.Init(server, db, adminAuthorizer)
	advertisement.Init(server, db)
	application.Init(server, db)
	company.Init(server, db)

	server.Listen(fmt.Sprintf(":%d", config.Server.Port))
}
