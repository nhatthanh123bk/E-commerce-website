package model

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" gorm:"size:200"`
}
