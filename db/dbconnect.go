package db

import (
	"log"

	"github.com/thuranaing10/RestfulAPI-Wth-Golang_Gin_Gorm/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Database *gorm.DB // Global database instance

// InitDB initializes the database connection
func InitDB() {
	dsn := "root:@tcp(127.0.0.1:3306)/golang_gin_gorm?charset=utf8mb4&parseTime=True&loc=Local"
	var err error

	// Open the database connection
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the models
	Database.AutoMigrate(&models.User{}, &models.Product{}, &models.Post{})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	if Database == nil {
		log.Fatal("Database not initialized") // This prevents nil pointer dereference
	}
	return Database
}
