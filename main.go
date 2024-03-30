package main

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	_ "github.com/KseniiaSalmina/tikkichest-notifications-service/docs"
	app "github.com/KseniiaSalmina/tikkichest-notifications-service/internal"
	"github.com/KseniiaSalmina/tikkichest-notifications-service/internal/config"
)

var (
	cfg config.Application
)

func init() {
	_ = godotenv.Load(".env")
	if err := env.Parse(&cfg); err != nil {
		panic(err)
	}
}

// @title Tikkichest notifications service
// @version 1.0.0
// @description part of tikkichest
// @host localhost:8080
// @BasePath /
func main() {
	application, err := app.NewApplication(cfg)
	if err != nil {
		log.Fatal(err)
	}

	application.Run()
}
