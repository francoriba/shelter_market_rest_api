// /main.go

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/routes"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

// @title New World API - Operating Systems Lab 3
// @version 1.0
// @description This API allows users to register, login, and access resources based on their role.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load .env file only if we're not in a GitHub Actions environment
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file:", err)
			// Don't fatal here, as the environment variables might be set another way
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable port=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort)

	// Initialize database connection
	database.InitDB(connStr)
	defer database.CloseDB()
	fmt.Println("Successfully connected to the database!")

	// Create new Fiber server
	app := fiber.New()

	// Middleware for CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3001, http://localhost:3000", // 3001 for local dev and qa, 3002 for docker deployment
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// First supply fetch from /supplies endpoint of HPCPP lab
	utils.FetchAndStoreSupplies()

	// Start cron job
	utils.StartCronJob()

	// Register routes
	routes.SwaggerRoute(app)     // Register a route for API Docs (Swagger).
	routes.SetupAuthRoutes(app)  // Register routes for the Auth API.
	routes.SetupAdminRoutes(app) // Register routes for the Admin API.
	routes.NotFoundRoute(app)    // Register a route for 404 Not Found.

	// Get the port from the environment variable, default to 3000 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Listen on the specified port
	log.Fatal(app.Listen(":" + port))
}
