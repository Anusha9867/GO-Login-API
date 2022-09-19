package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Id    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
}
