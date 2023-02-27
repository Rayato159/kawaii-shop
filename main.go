package main

import (
	"fmt"

	"github.com/Rayato159/kawaii-shop/config"
)

func main() {
	cfg := config.LoadConfig(".env.dev")
	fmt.Println(cfg)
}
