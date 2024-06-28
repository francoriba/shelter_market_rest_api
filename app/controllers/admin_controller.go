// app/controllers/admin_controller.go

package controllers

import (
	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// AdminController handles admin related requests
type AdminController struct {
	DB *gorm.DB
}

func NewAdminController(db *gorm.DB) *AdminController {
	return &AdminController{DB: db}
}

// GetDashboard returns the current status of all orders
// @Summary Get dashboard data
// @Description Get the current status of all orders
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} models.DashboardResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/dashboard [get]
func (adc *AdminController) GetDashboard(c *fiber.Ctx) error {
	db := database.GetDB()
	var orders []models.Order
	if err := db.Preload("OrderItems").Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch orders",
		})
	}

	var dashboardOrders []models.OrderDashboard
	for _, order := range orders {
		var orderItems []models.OrderItemDetails
		for _, item := range order.OrderItems {
			orderItems = append(orderItems, models.OrderItemDetails{
				OfferID:  item.OfferID,
				Quantity: item.Quantity,
				SubTotal: item.SubTotal,
			})
		}
		dashboardOrders = append(dashboardOrders, models.OrderDashboard{
			ID:          order.ID,
			Status:      order.Status,
			TotalAmount: order.TotalAmount,
			Items:       orderItems,
		})
	}

	return c.JSON(models.DashboardResponse{
		Status:  "success",
		Message: "Dashboard data fetched successfully",
		Orders:  dashboardOrders,
	})
}

// UpdateOrderStatus updates the status of a specific order
// @Summary Update order status
// @Description Update the status of a specific order
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param data body models.UpdateOrderStatusRequest true "Order status data"
// @Success 200 {object} models.UpdateOrderStatusResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/orders/{id} [patch]
func (adc *AdminController) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var request models.UpdateOrderStatusRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Bad request",
		})
	}

	// Validate the status value
	validStatus := map[string]bool{
		"preparing":  true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
	}
	if !validStatus[request.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Invalid status value",
		})
	}

	db := database.GetDB()
	var order models.Order
	if err := db.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Order not found",
		})
	}

	order.Status = request.Status
	if err := db.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to update order status",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.UpdateOrderStatusResponse{
		Code:    200,
		Message: "Order status updated successfully",
		Status:  request.Status,
	})
}

// GetAllUsers retrieves all buyers
// @Summary Retrieve all buyer users
// @Description Retrieves a list of all buyers, only accessible to administrators
// @Tags Admin
// @Accept json
// @Produce json
// @Success 200 {object} models.GetAllUsersResponse
// @Failure 500 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/users [get]
func (adc *AdminController) GetAllUsers(c *fiber.Ctx) error {
	db := database.GetDB()
	var users []models.User
	if err := db.Where("role = ?", "user").Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Server error",
		})
	}

	var userResponses []models.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, models.UserResponse{
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.GetAllUsersResponse{
		Code:    200,
		Message: userResponses,
	})
}

// DeleteUser handles the deletion of a user
// @Summary Remove a user
// @Description Delete a user by ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.DeleteUserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /admin/users/{id} [delete]
func (adc *AdminController) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")

	db := database.GetDB()
	var user models.User

	// Find the user by ID
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Code:    400,
				Message: "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to find user",
		})
	}

	// Delete the user
	if err := db.Delete(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to delete user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.DeleteUserResponse{
		Code:    200,
		Message: "User deleted successfully",
	})
}
