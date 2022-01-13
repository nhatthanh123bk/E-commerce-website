package route

import (
	"net/http"

	"github.com/blogs/app/controller"
	"github.com/blogs/app/middleware"
	"github.com/labstack/echo/v4"
)

func Init() *echo.Echo {
	e := echo.New()
	e.POST("/users", controller.CreateUser)
	e.POST("/users/login", controller.LoginUser)
	e.POST("users/refresh-token", controller.ReproduceTokenUser)

	// this group of routes need jwt to access
	jwtGr := e.Group("users/", middleware.UserAuth)
	jwtGr.Add(http.MethodGet, "profile", controller.FindUser)

	return e
}
