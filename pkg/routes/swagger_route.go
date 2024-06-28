package routes

import (
	"github.com/gofiber/fiber/v2"
	swagger "github.com/gofiber/swagger"

	_ "github.com/ICOMP-UNC/newworld-francoriba/docs"
)

// SwaggerRoute func for describe group of API Docs routes.
func SwaggerRoute(a *fiber.App) {
	// Create routes group.
	route := a.Group("/swagger")

	// Routes for GET method:
	route.Get("*", swagger.HandlerDefault)
}
