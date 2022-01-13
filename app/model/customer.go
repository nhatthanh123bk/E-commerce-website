package model

import "gorm.io/gorm"

type Customer struct {
	gorm.Model
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" gorm:"size:200"`
}
