package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// FalseParamResponse is respond when param is valid
func FalseParamResponse(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "False Param",
	}

	return c.JSON(http.StatusBadRequest, result)
}

func BadRequestResponse(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "Bad Request",
	}

	return c.JSON(http.StatusBadRequest, result)
}

// AccessForbiddenResponse is responded when user has no bussiness to access the route
func AccessForbiddenResponse(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusForbidden,
		"message": "Access Forbidden",
	}

	return c.JSON(http.StatusForbidden, result)
}

// SuccessResponseData is responded with data when operating successfully
func SuccessResponseData(c echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Successful Operation",
		"data":    data,
	}

	return c.JSON(http.StatusOK, result)
}

// SuccessResponseNonData is responded without data when operating successfully
func SuccessResponseNonData(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Successful Operation",
	}

	return c.JSON(http.StatusOK, result)
}

// LoginFailedResponse is responded when having bad request
func LoginFailedResponse(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusBadRequest,
		"message": "Login Failed",
	}

	return c.JSON(http.StatusBadRequest, result)
}

// UnAuthorization is responded when user provided wrong credentail.
func UnAuthorization(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusUnauthorized,
		"message": "UnAuthorization",
	}

	return c.JSON(http.StatusUnauthorized, result)
}

// LoginSuccessResponse is responded when logging in successfully
func LoginSuccessResponse(c echo.Context, data interface{}) error {
	result := map[string]interface{}{
		"code":    http.StatusOK,
		"message": "Login Success",
		"data":    data,
	}

	return c.JSON(http.StatusOK, result)
}

func InternalServerErrorResponse(c echo.Context) error {
	result := map[string]interface{}{
		"code":    http.StatusInternalServerError,
		"message": "Internal Server Error",
	}

	return c.JSON(http.StatusInternalServerError, result)
}
