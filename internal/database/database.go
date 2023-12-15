package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(host, user, password, database string) (*gorm.DB, error) {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", host, user, password, database)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{TranslateError: true})

	return db, err
}
