package utils

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/burhangltekin/byfood/models"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("books.db"), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	err = DB.AutoMigrate(&models.Book{})
	if err != nil {
		return fmt.Errorf("failed to migrate DB schema: %w", err)
	}
	return nil
}
