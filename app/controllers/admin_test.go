package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ICOMP-UNC/newworld-francoriba/app/controllers"
	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetAllUsers(t *testing.T) {
	app := fiber.New()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to open sqlmock database: %s", err)
	}
	defer db.Close()

	dsn := "user=gorm dbname=gorm sslmode=disable"
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	database.SetDB(gormDB)

	rows := sqlmock.NewRows([]string{"username", "email"}).
		AddRow("user1", "user1@example.com").
		AddRow("user2", "user2@example.com")
	mock.ExpectQuery(`SELECT \* FROM "users" WHERE role = \$1`).WithArgs("user").WillReturnRows(rows)

	ctrl := controllers.NewAdminController(database.DB)

	app.Get("/admin/users", ctrl.GetAllUsers)

	req := httptest.NewRequest("GET", "/admin/users", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %s", err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response models.GetAllUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	expectedResponse := models.GetAllUsersResponse{
		Code: 200,
		Message: []models.UserResponse{
			{Username: "user1", Email: "user1@example.com"},
			{Username: "user2", Email: "user2@example.com"},
		},
	}
	assert.Equal(t, expectedResponse, response)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetDashboard(t *testing.T) {
	// Setup Fiber app and mock database
	app := fiber.New()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to set up mock database: %s", err)
	}
	defer db.Close()

	// Configure mock database connection for GORM
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to initialize GORM: %s", err)
	}
	database.SetDB(gormDB)

	// Setup mock expectations
	mock.ExpectQuery(`SELECT \* FROM "orders"`).WillReturnError(fmt.Errorf("database error"))

	// Initialize admin controller with mock DB
	ctrl := controllers.NewAdminController(database.DB)

	// Register endpoint handler
	app.Get("/admin/dashboard", ctrl.GetDashboard)

	// Prepare request
	req := httptest.NewRequest("GET", "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer valid_token")

	// Perform request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %s", err)
	}

	// Assert response status code
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	// Decode error response
	var response models.ErrorResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode error response: %s", err)
	}

	// Verify error response content
	expectedError := models.ErrorResponse{
		Code:    500,
		Message: "Failed to fetch orders", // Adjust based on actual error handling
	}
	assert.Equal(t, expectedError, response)

	// Verify mock expectations
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Unfulfilled expectations: %s", err)
	}
}
