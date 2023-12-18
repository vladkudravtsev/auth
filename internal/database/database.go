package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(host, user, password, database string) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, database)

	config := &gorm.Config{
		TranslateError: true,
	}

	if os.Getenv("ENV") == "test" {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(connectionString), config)

	return db, err
}
