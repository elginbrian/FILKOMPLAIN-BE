package config

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
