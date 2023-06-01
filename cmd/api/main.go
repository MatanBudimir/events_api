package main

import (
	"b2match_api/pkg/api"
	"b2match_api/pkg/database"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	// connect to the database
	if err := database.Connect(); err != nil {
		log.Fatalln(err)
	}
	//start the API
	api.StartWebServer()
}
