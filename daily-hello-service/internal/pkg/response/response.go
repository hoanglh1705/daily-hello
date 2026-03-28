package response

import (
	"net/http"

	appErrors "daily-hello-service/internal/pkg/errors"

	"github.com/labstack/echo/v4"
)

// Response represents standard API response
type Response struct {
	Success      bool   `json:"success"`
	Data         interface{} `json:"data,omitempty"`
	ErrorCode    string `json:"error_code,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
}

func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    data,
	})
}

func Error(c echo.Context, err appErrors.AppError) error {
	return c.JSON(mapStatus(err.Code), Response{
		Success:      false,
		ErrorCode:    err.Code,
		ErrorMessage: err.Message,
	})
}

func HandleError(c echo.Context, err error) error {
	if appErr, ok := err.(appErrors.AppError); ok {
		return Error(c, appErr)
	}

	return c.JSON(http.StatusInternalServerError, Response{
		Success:      false,
		ErrorCode:    "INTERNAL_ERROR",
		ErrorMessage: "Something went wrong",
	})
}

func mapStatus(code string) int {
	switch code {
	case "INVALID_INPUT", "ALREADY_CHECKED_IN", "NOT_CHECKED_IN", "EMAIL_EXISTS", "FAKE_GPS":
		return http.StatusBadRequest
	case "UNAUTHORIZED", "INVALID_CREDENTIALS", "INVALID_TOKEN", "TOKEN_EXPIRED":
		return http.StatusUnauthorized
	case "FORBIDDEN", "ACCOUNT_INACTIVE":
		return http.StatusForbidden
	case "NOT_FOUND", "BRANCH_NOT_FOUND":
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
