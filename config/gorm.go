package config

import (
	"log"

	"github.com/storyofhis/books-management/httpserver/repository/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectGorm() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database : %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&models.Author{}, &models.Book{}, &models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database : %v", err)
		return nil, err
	}
	return db, err
}
