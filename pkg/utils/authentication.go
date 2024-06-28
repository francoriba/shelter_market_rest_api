// pkg/utils/authentication.go
package utils

import (
	"errors"
	"os"
	"time"

	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/ICOMP-UNC/newworld-francoriba/pkg/database"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Function variables to allow injection of mock implementations in tests
var (
	AuthenticateUserFunc       = authenticateUser
	GenerateJWTTokenFunc       = generateJWTToken
	BcryptGenerateFromPassword = bcrypt.GenerateFromPassword
)

// GenerateJWTToken generates a JWT token with the given email
func generateJWTToken(email, role string) (string, error) {
	// Retrieve the JWT secret key from the environment variable

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")

	if jwtSecretKey == "" {
		return "", errors.New("JWT secret key not found")
	}

	// Create a new token object, specifying signing method and claims
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expiration time (24 hours)

	// Sign the token with the JWT secret key
	tokenString, err := token.SignedString([]byte(jwtSecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// AuthenticateUser authenticates a user with the given login credentials
func authenticateUser(loginRequest models.LoginRequest) (models.User, error) {
	// Query the database to find the user by email
	db := database.GetDB()
	var user models.User
	if err := db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		// If user not found or an error occurs, return an error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, err
	}

	// Compare the hashed password from the database with the provided password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		// If passwords don't match, return an error
		return models.User{}, errors.New("invalid password")
	}

	// If authentication is successful, return the authenticated user
	return user, nil
}
