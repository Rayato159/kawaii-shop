package main

import (
	"os"

	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/servers"
	"github.com/Rayato159/kawaii-shop/pkg/databases"
)

func envPath() string {
	if len(os.Args) == 1 {
		return ".env"
	} else {
		return os.Args[1]
	}
}

func main() {
	// Setup config
	cfg := config.LoadConfig(envPath())

	// Db setup
	db := databases.DbConnect(cfg.Db())
	defer db.Close()

	// Server start
	servers.NewServer(cfg.App(), db).Start()
}
