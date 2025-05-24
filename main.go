package main

import (
	"log"
	"os"
	"project_restfulApi_go/database"
	"project_restfulApi_go/middleware"
	"project_restfulApi_go/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	app := fiber.New()

	// Initialize database
	database.ConnectDB()

	isDev := os.Getenv("APP_ENV") == "development"

	// Register middleware
	app.Use(middleware.Logger(isDev))
	// Register routes
	routes.UserRoutes(app, database.DB)

	// Start server
	port := os.Getenv("APP_SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	err := app.Listen(":" + port)
	if err != nil {

		log.Printf("Shutting down due to error: %v", err)
		os.Exit(1) // Means exit with error
	}
}
