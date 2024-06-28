// app/controllers/auth_controller.go

package controllers

import (
	"fmt"
	"os"
	"strings"

	"github.com/ICOMP-UNC/newworld-francoriba/pkg/utils"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

const DefaultRole = "user"
const AdminRole = "admin"

// AuthController handles authentication related requests
type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{DB: db}
}

// Register handles the registration of a new user
// @Summary Register a new user
// @Description Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body models.RegisterRequest true "User data to register"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/register [post]
func (ac *AuthController) Register(c *fiber.Ctx) error {
	var requestData models.RegisterRequest
	if err := c.BodyParser(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "bad request",
		})
	}

	db := database.GetDB()
	// Validate username, email, and password
	if err := utils.ValidateRegistrationRequest(requestData, db); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

	hashedPassword, err := utils.BcryptGenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to hash password",
		})
	}

	// Create the user with default role
	user := models.User{
		Username: requestData.Username,
		Email:    requestData.Email,
		Password: string(hashedPassword),
		Role:     DefaultRole, // Set the default role
	}

	// Save the user to the database

	if err := db.Create(&user).Error; err != nil {
		fmt.Println("errorrrrrrrrrrrr", err)

		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to register user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.SuccessResponse{
		Code:    201,
		Message: "User registered successfully",
	})
}

// Login handles user authentication
// @Summary Authenticate a user
// @Description Authenticate a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body models.LoginRequest true "User credentials for login"
// @Success 200 {string} JWT "Authentication token"
// @Failure 400 {object} models.ErrorResponse
// @Router /auth/login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
	var loginRequest models.LoginRequest
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Bad request",
		})
	}

	// Authenticate user (check username/password)
	user, err := utils.AuthenticateUserFunc(loginRequest)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Invalid credentials",
		})
	}

	// If authentication is successful, generate a JWT token
	// Pass the user's role (e.g., "admin" or "regular") to the GenerateJWTToken function
	token, err := utils.GenerateJWTTokenFunc(user.Email, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to generate JWT token",
		})
	}

	// Set the token in the response header
	c.Set("Authorization", "Bearer "+token)

	// Return the token to the client
	return c.JSON(fiber.Map{"token": token})
}

// GetOffers handles the retrieval of available offers
// @Summary Retrieve a list of available offers
// @Description Retrieve a list of available offers
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.OfferResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /auth/offers [get]
func (ac *AuthController) GetOffers(c *fiber.Ctx) error {

	// Placeholder authentication check
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "JWT secret key not found",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Invalid Authorization header format",
		})
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Invalid or expired JWT",
		})
	}

	// If the token is valid, fetch the offers from the database
	db := database.GetDB()
	offers := []models.Offer{}

	if err := db.Find(&offers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch offers from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.OfferResponse{
		Code:    200,
		Message: offers,
	})
}

// Checkout handles the checkout process
// @Summary Process checkout
// @Description Process checkout and create an order
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body models.CheckoutRequest true "Checkout request data"
// @Success 200 {object} models.CheckoutResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /auth/checkout [post]
func (ac *AuthController) Checkout(c *fiber.Ctx) error {
	var checkoutRequest models.CheckoutRequest
	if err := c.BodyParser(&checkoutRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Bad request",
		})
	}

	// Validate the request
	if err := utils.ValidateCheckoutRequest(checkoutRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: err.Error(),
		})
	}

	db := database.GetDB()
	tx := db.Begin() // Start a transaction

	// Validate availability and calculate total amount
	var totalAmount float64
	var orderItems []models.OrderItem
	for _, item := range checkoutRequest.Items {
		var offer models.Offer
		if err := tx.First(&offer, item.OfferID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Code:    400,
				Message: fmt.Sprintf("Offer with ID %d not found", item.OfferID),
			})
		}
		if offer.Quantity < item.Quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
				Code:    400,
				Message: fmt.Sprintf("Not enough quantity for offer ID %d", item.OfferID),
			})
		}
		subTotal := float64(item.Quantity) * offer.Price
		totalAmount += subTotal
		orderItems = append(orderItems, models.OrderItem{
			OfferID:  item.OfferID,
			Quantity: item.Quantity,
			SubTotal: subTotal,
		})
	}

	// Create the order
	order := models.Order{
		Status:      "processing",
		OrderItems:  orderItems,
		TotalAmount: totalAmount,
	}
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to create order",
		})
	}

	// Update the stock of the offers
	for _, item := range checkoutRequest.Items {
		if err := tx.Model(&models.Offer{}).Where("id = ?", item.OfferID).Update("quantity", gorm.Expr("quantity - ?", item.Quantity)).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
				Code:    500,
				Message: "Failed to update offer stock",
			})
		}
	}

	tx.Commit() // Commit the transaction

	return c.Status(fiber.StatusOK).JSON(models.CheckoutResponse{
		Code:    200,
		Message: "Order created successfully",
		OrderID: order.ID,
	})
}

// GetOrderStatus handles the retrieval of the status of a specific order
// @Summary Get the status of a specific order
// @Description Retrieve the status of a specific order
// @Tags Auth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Success 200 {object} models.OrderResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security BearerAuth
// @Router /auth/orders/{id} [get]
func (ac *AuthController) GetOrderStatus(c *fiber.Ctx) error {
	// Authentication check
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Unauthorized",
		})
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "JWT secret key not found",
		})
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Invalid Authorization header format",
		})
	}

	tokenString := parts[1]

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(jwtSecretKey), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(models.ErrorResponse{
			Code:    401,
			Message: "Invalid or expired JWT",
		})
	}

	// Retrieve the order ID from the URL
	orderID := c.Params("id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{
			Code:    400,
			Message: "Order ID is required",
		})
	}

	// Search for the order in the database
	var order models.Order
	db := database.GetDB()
	if err := db.First(&order, orderID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.ErrorResponse{
			Code:    500,
			Message: "Failed to fetch order from the database",
		})
	}

	// Return the status of the order
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"code": 200,
		"message": fiber.Map{
			"status": order.Status,
		},
	})
}
