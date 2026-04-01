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
	rbac    *services.RBACService
}

func NewUserHandler(service *services.UserService, rbac *services.RBACService) *UserHandler {
	return &UserHandler{service: service, rbac: rbac}
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
// @Router /v1/users [post]
func (h *UserHandler) Register(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var req models.CreateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if req.BranchID != nil {
		if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, *req.BranchID); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
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
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	role, _ := c.Get("role").(string)
	currentUserID := uint(c.Get("user_id").(float64))
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if models.Role(role) == models.RoleEmployee && uint(id) != currentUserID {
		return response.Error(c, appErrors.ErrForbidden)
	}
	if models.Role(role) == models.RoleManager {
		if err := h.rbac.EnsureUserAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
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
// @Router /v1/users/me [get]
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
// @Router /v1/users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	role, _ := c.Get("role").(string)
	currentUserID := uint(c.Get("user_id").(float64))
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if models.Role(role) == models.RoleEmployee && uint(id) != currentUserID {
		return response.Error(c, appErrors.ErrForbidden)
	}
	if models.Role(role) == models.RoleManager {
		if err := h.rbac.EnsureUserAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
	}

	var req models.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if req.BranchID != nil {
		if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, *req.BranchID); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
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
// @Param request body models.UserListQuery true "Pagination query"
// @Success 200 {object} response.Response{data=models.PaginatedResponse{items=[]models.User}} "List of users"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/users [get]
func (h *UserHandler) List(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

	var q models.UserListQuery
	if err := c.Bind(&q); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if role == string(models.RoleManager) {
		currentBranchID, err := getContextUint(c, "branch_id")
		if err != nil {
			return response.Error(c, appErrors.ErrInvalidInput)
		}
		if currentBranchID == nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
		if q.BranchID != nil {
			if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, *q.BranchID); err != nil {
				return response.Error(c, appErrors.ErrForbidden)
			}
		} else {
			q.BranchIDs, err = h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
			if err != nil {
				return response.HandleError(c, err)
			}
		}
	}

	result, err := h.service.List(c.Request().Context(), q)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
