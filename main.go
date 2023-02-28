package main

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/servers"
	"github.com/Rayato159/kawaii-shop/pkg/databases"
)

func main() {
	// Setup config
	cfg := config.LoadConfig()

	// Db setup
	db := databases.DbConnect(cfg.Db())

	// Server start
	servers.NewServer(cfg.App(), db).Start()
}
