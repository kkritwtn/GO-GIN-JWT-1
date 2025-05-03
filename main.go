package main

import (
	"go-jwt-crud/config"
	"go-jwt-crud/models"
	"go-jwt-crud/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" 
	}

	// Database and route setup
	config.ConnectDB()
	config.DB.AutoMigrate(&models.User{})
	config.DB.AutoMigrate(&models.Admin{})
	r := gin.Default()
	routes.SetupRoutes(r)

	// Run the server on the specified port
	r.SetTrustedProxies([]string{"127.0.0.1"})
	r.Run(":" + port) 
}
