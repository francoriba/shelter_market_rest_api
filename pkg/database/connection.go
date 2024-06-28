// pkg/database/connection.go

package database

import (
	"log"

	"github.com/ICOMP-UNC/newworld-francoriba/app/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(connStr string) {
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	// Auto migrate models
	err = DB.AutoMigrate(&models.User{}, &models.Offer{}, &models.Order{}, &models.OrderItem{})
	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}

func GetDB() *gorm.DB {
	return DB
}

func CloseDB() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Failed to close database connection: ", err)
	}
	sqlDB.Close()
}

// SetDB sets the database connection (useful for testing)
func SetDB(database *gorm.DB) {
	DB = database
}
