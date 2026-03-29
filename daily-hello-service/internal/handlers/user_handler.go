package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	service *services.UserService
}

func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// @Summary Register
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.CreateUserRequest true "User data"
// @Success 201 {object} response.Response{data=models.User} "User created successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/user [post]
func (h *UserHandler) Register(c echo.Context) error {
	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

// @Summary Get User by ID
// @Description Get user by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=models.User} "User details"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/auth/user/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Get Current User
// @Description Get current user
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.User} "User details"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/auth/user/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(float64)

	result, err := h.service.GetByID(c.Request().Context(), uint(userID))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Update User
// @Description Update user
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body models.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response{data=models.User} "User updated successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/auth/user/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary List Users
// @Description List users
// @Tags User
// @Accept json
// @Produce json
// @Param request body models.PaginationQuery true "Pagination query"
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "List of users"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/user [get]
func (h *UserHandler) List(c echo.Context) error {
	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.List(c.Request().Context(), pq)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
