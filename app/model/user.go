package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" gorm:"size:200"`
}

type Token map[string]string

type TokenDetail struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   string
	RefreshUuid  string
	AtExpires    int64
	RfExpires    int64
}

type AccessTokenMetadata struct {
	AccessUuid string
	UserId     int
}

type RefreshTokenMetadata struct {
	RefreshUuid string
	UserId      int
}
