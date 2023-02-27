package main

import (
	"github.com/Rayato159/kawaii-shop/config"
	"github.com/Rayato159/kawaii-shop/modules/servers"
)

func main() {
	cfg := config.LoadConfig()

	servers.NewServer(cfg.App()).Start()
}
