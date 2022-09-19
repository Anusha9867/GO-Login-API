package models

import "gorm.io/gorm"

type Password struct {
	gorm.Model
	UserId   uint   `json:"user_id"`
	Password []byte `json:"-"`
}
