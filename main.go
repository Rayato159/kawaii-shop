package main

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/servers"
)

func main() {
	// Setup config
	cfg := config.LoadConfig()

	// Server start
	servers.NewServer(cfg.App()).Start()
}
