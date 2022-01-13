package middleware

import (
	"strconv"

	"github.com/go-redis/redis/v7"
	"github.com/labstack/echo/v4"
	"github.com/nhatthanh123bk/E-commerce-website/app/repository"
	"github.com/nhatthanh123bk/E-commerce-website/app/response"
	"github.com/nhatthanh123bk/E-commerce-website/db"
)

func UserAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		givenTokenMetadata, err := repository.ExtractTokenMetadata(req)
		if err != nil {
			return response.BadRequestResponse(c)
		}
		userId, err := db.RedisClient.Get(givenTokenMetadata.AccessUuid).Result()

		// if the given key isn't existing
		if err == redis.Nil {
			return response.UnAuthorization(c)
		}
		if err != nil {
			return response.InternalServerErrorResponse(c)
		}

		// if the given userId doesn't match userId that is stored in redis
		if userId != strconv.Itoa(givenTokenMetadata.UserId) {
			return response.UnAuthorization(c)
		}
		return next(c)
	}
}
