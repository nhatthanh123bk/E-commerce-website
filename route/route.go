package route

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/nhatthanh123bk/E-commerce-website/app/controller"
	"github.com/nhatthanh123bk/E-commerce-website/app/middleware"
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
