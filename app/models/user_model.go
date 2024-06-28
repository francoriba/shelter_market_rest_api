// app/models/user_model.go

package models

import "gorm.io/gorm"

// User model represents a user in the system
type User struct {
	gorm.Model        // Embeds fields `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`
	Username   string `json:"username" gorm:"uniqueIndex;not null"` // Unique username, cannot be null
	Email      string `json:"email" gorm:"uniqueIndex;not null"`    // Unique email, cannot be null
	Password   string `json:"password" gorm:"not null"`             // Password, cannot be null
	Role       string `json:"role" gorm:"not null"`                 // Role of the user (e.g., "admin" or "regular")
}
