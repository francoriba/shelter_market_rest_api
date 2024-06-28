// app/controllers/auth_tests.go
package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ICOMP-UNC/newworld-francoriba/app/controllers"
	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db     *sql.DB
	mock   sqlmock.Sqlmock
	gormDB *gorm.DB
)

func setupMockDB(t *testing.T) {
	var err error
	db, mock, err = sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to set up mock database: %s", err)
	}

	dsn := "user=gorm dbname=gorm sslmode=disable"
	gormDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to open gorm db: %s", err)
	}

	database.SetDB(gormDB)

}

func TestMain(m *testing.M) {
	setupMockDB(nil)
	defer db.Close()

	code := m.Run()

	if err := mock.ExpectationsWereMet(); err != nil {
		panic("There were unfulfilled expectations")
	}

	os.Exit(code)
}

func TestGetOffers(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	os.Setenv("JWT_SECRET_KEY", "a6a6d01782e0cb082ad4b016a508d4a913c2556f38aa26303d0d00562e111aaa")
	defer os.Unsetenv("JWT_SECRET_KEY")

	app := fiber.New()

	rows := sqlmock.NewRows([]string{"id", "name", "quantity", "price", "category"}).
		AddRow(1, "Offer 1", 10, 20.5, "Category A").
		AddRow(2, "Offer 2", 5, 15.75, "Category B")
	mock.ExpectQuery(`SELECT \* FROM "offers"`).WillReturnRows(rows)

	ctrl := controllers.NewAuthController(database.DB)

	app.Get("/auth/offers", ctrl.GetOffers)

	user := models.User{
		Email: "user@example.com",
		Role:  "user",
	}

	token, err := utils.GenerateJWTTokenFunc(user.Email, user.Role)
	if err != nil {
		t.Fatalf("Failed to generate JWT token: %v", err)
	}

	req := httptest.NewRequest("GET", "/auth/offers", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %s", err)
	}

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response models.OfferResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	expectedResponse := models.OfferResponse{
		Code:    200,
		Message: []models.Offer{{ID: 1, Name: "Offer 1", Quantity: 10, Price: 20.5, Category: "Category A"}, {ID: 2, Name: "Offer 2", Quantity: 5, Price: 15.75, Category: "Category B"}},
	}
	assert.Equal(t, expectedResponse, response)
}

func TestLogin(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	tests := []struct {
		name           string
		loginRequest   models.LoginRequest
		mockAuthUser   func(models.LoginRequest) (models.User, error)
		mockGenToken   func(string, string) (string, error)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name: "Successful login",
			loginRequest: models.LoginRequest{
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			mockAuthUser: func(loginRequest models.LoginRequest) (models.User, error) {
				return models.User{
					Email: "valid@example.com",
					Role:  "user",
				}, nil
			},
			mockGenToken: func(email, role string) (string, error) {
				return "mockJWTToken", nil
			},
			expectedStatus: http.StatusOK,
			expectedBody: fiber.Map{
				"token": "mockJWTToken",
			},
		},
		{
			name: "Invalid credentials",
			loginRequest: models.LoginRequest{
				Email:    "invalid@example.com",
				Password: "invalidpassword",
			},
			mockAuthUser: func(loginRequest models.LoginRequest) (models.User, error) {
				return models.User{}, errors.New("invalid credentials")
			},
			mockGenToken: func(email, role string) (string, error) {
				return "", nil
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody: models.ErrorResponse{
				Code:    401,
				Message: "Invalid credentials",
			},
		},
		{
			name: "Failed to generate JWT token",
			loginRequest: models.LoginRequest{
				Email:    "valid@example.com",
				Password: "validpassword",
			},
			mockAuthUser: func(loginRequest models.LoginRequest) (models.User, error) {
				return models.User{
					Email: "valid@example.com",
					Role:  "user",
				}, nil
			},
			mockGenToken: func(email, role string) (string, error) {
				return "", errors.New("failed to generate JWT token")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: models.ErrorResponse{
				Code:    500,
				Message: "Failed to generate JWT token",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			utils.AuthenticateUserFunc = tt.mockAuthUser
			utils.GenerateJWTTokenFunc = tt.mockGenToken

			app := fiber.New()

			ctrl := &controllers.AuthController{}

			app.Post("/auth/login", ctrl.Login)

			requestBody, _ := json.Marshal(tt.loginRequest)
			req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to perform request: %s", err)
			}

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var responseBody interface{}
			if tt.expectedStatus == http.StatusOK {
				var tokenResp fiber.Map
				if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
					t.Fatalf("Failed to decode response: %s", err)
				}
				responseBody = tokenResp
			} else {
				var errorResp models.ErrorResponse
				if err := json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode response: %s", err)
				}
				responseBody = errorResp
			}

			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func TestRegister(t *testing.T) {
	setupMockDB(t)
	defer db.Close()

	app := fiber.New()

	ctrl := controllers.NewAuthController(database.DB)

	app.Post("/auth/register", ctrl.Register)

	requestData := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpassword",
	}

	mockPassword := "$2a$10$mockpasswordmockpasswordmockpass"

	utils.BcryptGenerateFromPassword = func(password []byte, cost int) ([]byte, error) {
		return []byte(mockPassword), nil
	}

	defer func() {
		utils.BcryptGenerateFromPassword = bcrypt.GenerateFromPassword
	}()

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE username = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(requestData.Username, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectQuery(`SELECT \* FROM "users" WHERE email = \$1 AND "users"\."deleted_at" IS NULL ORDER BY "users"\."id" LIMIT \$2`).
		WithArgs(requestData.Email, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "users" \("created_at","updated_at","deleted_at","username","email","password","role"\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\) RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			requestData.Username,
			requestData.Email,
			mockPassword,
			"user",
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()

	requestBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != fiber.StatusCreated {
		t.Fatalf("Expected status %d but got %d", fiber.StatusCreated, resp.StatusCode)
	}

	var response models.SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %s", err)
	}

	expectedResponse := models.SuccessResponse{
		Code:    201,
		Message: "User registered successfully",
	}
	if response != expectedResponse {
		t.Fatalf("Expected response %+v but got %+v", expectedResponse, response)
	}
}
