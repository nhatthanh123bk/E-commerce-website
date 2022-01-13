package repository

import (
	"errors"

	"github.com/nhatthanh123bk/E-commerce-website/app/model"
	"github.com/nhatthanh123bk/E-commerce-website/db"
	"github.com/nhatthanh123bk/E-commerce-website/helper"
	"golang.org/x/crypto/bcrypt"
)

var ErrIncorrectEmail = errors.New("incorrect email")

func CreateUser(user *model.User) error {
	if err := db.DB.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

// LoginUser check whether a given password match the password that has been stored
// in database or not. If matching, proceeding to create a couple of tokens, stored
// token metadata in Redis and then sent tokens to end user
func LoginUser(password string, user *model.User) (model.Token, error) {
	err := db.DB.Where("email = ?", user.Email).First(&user).Error
	if err != nil {
		return nil, ErrIncorrectEmail
	}
	helper.Logger.Infow("user infomation!", "users: ", user, "password:", password)
	matchErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if matchErr != nil {
		helper.Logger.Errorw("Error by hashing!", "error:", matchErr)
		return nil, matchErr
	}
	td, err := CreateToken(user.ID)
	if err != nil {
		helper.Logger.Errorw(
			"Fail to create tokens",
			"error: ", err,
		)
		return nil, err
	}

	err = StoreTokenIntoRedis(td, user.ID)
	if err != nil {
		helper.Logger.Errorw(
			"Fail to store the tokens into redis",
			"error: ", err,
		)
		return nil, err
	}
	tokens := model.Token{
		"access_token":  td.AccessToken,
		"refresh_token": td.RefreshToken,
	}

	return tokens, nil
}
