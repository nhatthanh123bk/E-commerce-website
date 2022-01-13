package controller

import (
	"github.com/blogs/app/model"
	"github.com/blogs/app/response"

	"github.com/blogs/app/repository"
	"github.com/blogs/helper"

	"github.com/labstack/echo/v4"
)

func CreateUser(c echo.Context) error {
	newUser := model.User{}
	if err := c.Bind(&newUser); err != nil {
		return response.BadRequestResponse(c)
	}
	newUser.Password, _ = helper.HashPassword(newUser.Password)
	if err := repository.CreateUser(&newUser); err != nil {
		return response.InternalServerErrorResponse(c)
	}

	return response.SuccessResponseData(c, &newUser)
}

// extract token from header request
// decode token and then authenticate info that brought by token
// check token is valid or not, compare given info and redis info
func FindUser(c echo.Context) error {
	return response.SuccessResponseData(c, "Authenticated!")
}

func LogoutUser(c echo.Context) error {
	err := repository.DeleteToken(c.Request())
	if err != nil {
		return response.InternalServerErrorResponse(c)
	}

	return response.SuccessResponseNonData(c)
}
