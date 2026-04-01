package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type BranchHandler struct {
	service *services.BranchService
	rbac    *services.RBACService
}

func NewBranchHandler(service *services.BranchService, rbac *services.RBACService) *BranchHandler {
	return &BranchHandler{service: service, rbac: rbac}
}

// @Summary Create Branch
// @Description Create a new branch
// @Tags Branch
// @Accept json
// @Produce json
// @Param request body models.CreateBranchRequest true "Branch data"
// @Success 201 {object} response.Response{data=models.Branch} "Create Branch successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branches [post]
func (h *BranchHandler) Create(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var req models.CreateBranchRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if role == string(models.RoleManager) && req.ParentBranchCode != "" {
		branchIDs, err := h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
		if err != nil {
			return response.HandleError(c, err)
		}
		allowed := false
		for _, branchID := range branchIDs {
			branch, getErr := h.service.GetByID(c.Request().Context(), branchID)
			if getErr == nil && branch.BranchCode == req.ParentBranchCode {
				allowed = true
				break
			}
		}
		if !allowed {
			return response.Error(c, appErrors.ErrForbidden)
		}
	}

	result, err := h.service.Create(c.Request().Context(), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

// @Summary Get Branch by ID
// @Description Get a branch by ID
// @Tags Branch
// @Accept json
// @Produce json
// @Param        id  	  path      int        true  "Branch id"
// @Success 200 {object} response.Response{data=models.Branch} "Get Branch successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branches/{id} [get]
func (h *BranchHandler) GetByID(c echo.Context) error {
	role, _ := c.Get("role").(string)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
		return response.Error(c, appErrors.ErrForbidden)
	}

	result, err := h.service.GetByID(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Get Branch by ID
// @Description Get a branch by ID
// @Tags Branch
// @Accept json
// @Produce json
// @Param        id  	  path      int        true  "Branch id"
// @Param request body models.UpdateBranchRequest true "Branch data"
// @Success 200 {object} response.Response{data=models.Branch} "Update Branch successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branches/{id} [put]
func (h *BranchHandler) Update(c echo.Context) error {
	role, _ := c.Get("role").(string)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
		return response.Error(c, appErrors.ErrForbidden)
	}

	var req models.UpdateBranchRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Get Branch by ID
// @Description Get a branch by ID
// @Tags Branch
// @Accept json
// @Produce json
// @Param        id  	  path      int        true  "Branch id"
// @Success 200 {object} response.Response{data=models.Branch} "Delete Branch successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branches/{id} [delete]
func (h *BranchHandler) Delete(c echo.Context) error {
	role, _ := c.Get("role").(string)
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
		return response.Error(c, appErrors.ErrForbidden)
	}

	if err := h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "Branch deleted successfully"})
}

// @Summary Get Branch by ID
// @Description Get a branch by ID
// @Tags Branch
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "List Branch successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branches [get]
func (h *BranchHandler) List(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}
	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	var branchIDs []uint
	if role == string(models.RoleManager) {
		branchIDs, err = h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
		if err != nil {
			return response.HandleError(c, err)
		}
	}

	result, err := h.service.List(c.Request().Context(), pq, branchIDs)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
