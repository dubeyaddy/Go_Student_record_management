package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	driver := os.Getenv("DB_DRIVER")
	dbName := os.Getenv("DB_NAME")

	fmt.Println("Driver:", driver)
	fmt.Println("DB Name:", dbName)

	var err error

	switch driver {
	case "sqlite":
		DB, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	default:
		log.Fatalf("Unsupported DB_DRIVER: %s", driver)
	}

	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	fmt.Println("Successfully connected to database!")
}
