package controller

import (
	"github.com/blogs/app/model"
	"github.com/blogs/app/repository"
	"github.com/blogs/app/response"
	"github.com/blogs/helper"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func LoginUser(c echo.Context) error {
	user := model.User{}
	if err := c.Bind(&user); err != nil {
		helper.Logger.Errorw("Invalid request provided!", "error", err)
		return response.BadRequestResponse(c)
	}

	token, err := repository.LoginUser(user.Password, &user)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		helper.Logger.Errorw("Incorrect password provided!", "provided password:", user.Password)
		return response.UnAuthorization(c)
	}

	if err == repository.ErrIncorrectEmail {
		helper.Logger.Errorw("Incorrect email provided!", "provided email:", user.Email)
		return response.UnAuthorization(c)
	}

	if err != nil {
		helper.Logger.Errorw("Bad request provided!", "error", err)
		return response.InternalServerErrorResponse(c)
	}

	return response.LoginSuccessResponse(c, &token)
}

func ReproduceTokenUser(c echo.Context) error {
	td, err := repository.GenerateTokenFromRefreshToken(c)
	if err == helper.ErrUnAuthorization {
		return response.UnAuthorization(c)
	}
	if err == helper.ErrInternal {
		return response.InternalServerErrorResponse(c)
	}

	return response.SuccessResponseData(c, td)
}
