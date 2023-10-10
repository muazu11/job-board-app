package main

import (
	"jobboard/backend/config"
	"jobboard/backend/db"
	"jobboard/backend/server"
	"jobboard/backend/services/user"
)

func main() {
	config := config.New()
	server := server.New(config.Server)
	db := db.New(config.DB)

	user.Init(server, db)
}
