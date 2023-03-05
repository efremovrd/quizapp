package main

import (
	"log"
	"quizapp/config"
	"quizapp/internal/server"
	"quizapp/pkg/postgres"
)

// @title Quiz REST API
// @version 1.0
// @description REST API of Quiz app with Golang, Gin and PostgreSQL
// @contact.name Efremov Roman
// @contact.url https://github.com/efremovrd
// @contact.email efremovrd@yandex.ru
// @BasePath  /api/v1
// @securityDefinitions.apikey JWTToken
// @in header
// @name Authorization
func main() {
	log.Println("Starting api server")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	psqlDB, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	s := server.New(cfg, psqlDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
