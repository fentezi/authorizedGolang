package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go-auth/models"
	"go-auth/routes"
	"log"
	"os"
)

func main() {
	r := gin.Default()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	config := models.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DbName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	models.InitDB(config)
	routes.AuthRoutes(r)
	log.Fatal(r.Run(":8080"))
}
