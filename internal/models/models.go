package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"uniqueIndex"`
	PasswordHash string
}

type App struct {
	gorm.Model
	Name   string
	Secret string
}
