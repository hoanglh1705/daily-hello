package handlers

import (
	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "Login request"
// @Success 200 {object} models.LoginResponse
// @Router /v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req models.LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// Logout godoc
// @Summary Logout (revoke refresh token)
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.LogoutRequest true "Logout request"
// @Security BearerTokenAuth
// @Success 200
// @Router /v1/auth/logout [post]
func (h *AuthHandler) Logout(c echo.Context) error {
	var req models.LogoutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if err := h.service.Logout(c.Request().Context(), req); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, nil)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body models.RefreshTokenRequest true "Refresh token request"
// @Success 200 {object} models.RefreshTokenResponse
// @Router /v1/auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req models.RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.RefreshToken(c.Request().Context(), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
