// pkg/routes/auth_routes.go

package routes

import (
	"github.com/ICOMP-UNC/newworld-francoriba/app/controllers"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	db := database.GetDB()

	authController := controllers.NewAuthController(db)
	app.Post("/auth/register", authController.Register)
	app.Post("/auth/login", authController.Login)

	// Protect these routes with JWTMiddleware
	app.Get("/auth/offers", middleware.JWTMiddleware, authController.GetOffers)
	app.Post("/auth/checkout", middleware.JWTMiddleware, authController.Checkout)
	app.Get("/auth/orders/:id", middleware.JWTMiddleware, authController.GetOrderStatus)
}
