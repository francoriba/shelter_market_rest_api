// pkg/routes/admin_routes.go

package routes

import (
	"github.com/ICOMP-UNC/newworld-francoriba/app/controllers"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {

	db := database.GetDB()
	adminController := controllers.NewAdminController(db)

	// Protect these routes with both JWTMiddleware and AdminMiddleware
	app.Get("/admin/dashboard", middleware.JWTMiddleware, middleware.AdminMiddleware, adminController.GetDashboard)
	app.Patch("/admin/orders/:id", middleware.JWTMiddleware, middleware.AdminMiddleware, adminController.UpdateOrderStatus)
	app.Get("/admin/users", middleware.JWTMiddleware, middleware.AdminMiddleware, adminController.GetAllUsers)
	app.Delete("/admin/users/:id", middleware.JWTMiddleware, middleware.AdminMiddleware, adminController.DeleteUser)
}
