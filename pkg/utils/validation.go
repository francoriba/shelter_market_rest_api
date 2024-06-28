// pkg/utils/validation.go

package utils

import (
	"errors"
	"regexp"
	"strings"

	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate

// Init initializes the validator
func init() {
	validate = validator.New()
}

// ValidationFunc defines a function type for validating checkout requests
type ValidationFunc func(models.CheckoutRequest) error

// ValidateCheckoutRequest validates the CheckoutRequest data
var ValidateCheckoutRequest ValidationFunc = defaultValidateCheckoutRequest

// ValidateRegistrationRequest validates the registration request data
func ValidateRegistrationRequest(requestData models.RegisterRequest, db *gorm.DB) error {
	if requestData.Username == "" || requestData.Email == "" || requestData.Password == "" {
		return errors.New("all fields are required")
	}

	// Check if the username already exists
	var existingUser models.User
	if db.Where("username = ?", requestData.Username).First(&existingUser).Error == nil {
		return errors.New("username already exists")
	}

	// Check if the email already exists
	if db.Where("email = ?", requestData.Email).First(&existingUser).Error == nil {
		return errors.New("email already exists")
	}

	// Validate email format
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if match, _ := regexp.MatchString(emailRegex, requestData.Email); !match {
		return errors.New("invalid email format")
	}

	// Check password length
	if len(requestData.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

// ValidateCheckoutRequest validates the CheckoutRequest data
func defaultValidateCheckoutRequest(request models.CheckoutRequest) error {
	// Validate the request struct
	if err := validate.Struct(request); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, validationErr := range validationErrors {
			errorMessages = append(errorMessages, validationErr.Error())
		}
		return errors.New("Invalid request: " + strings.Join(errorMessages, ", "))
	}

	// Additional custom validation can be added here
	if len(request.Items) == 0 {
		return errors.New("the order must contain at least one item")
	}

	return nil
}
