package handlers

import (
	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type DashboardHandler struct {
	service *services.DashboardService
	rbac    *services.RBACService
}

func NewDashboardHandler(service *services.DashboardService, rbac *services.RBACService) *DashboardHandler {
	return &DashboardHandler{service: service, rbac: rbac}
}

// @Summary Dashboard Overview
// @Description Get dashboard summary and trends
// @Tags Admin/Dashboard
// @Accept json
// @Produce json
// @Param branch_id query int false "Branch ID"
// @Param date query string false "Date YYYY-MM-DD"
// @Success 200 {object} response.Response{data=models.DashboardOverviewResponse} "Dashboard Overview"
// @Router /v1/admin/dashboard/overview [get]
func (h *DashboardHandler) GetOverview(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, _ := getContextUint(c, "branch_id")

	var branchID *int64
	if bID := c.QueryParam("branch_id"); bID != "" && bID != "0" {
		id, err := strconv.ParseInt(bID, 10, 64)
		if err == nil {
			branchID = &id
		}
	}

	dateStr := c.QueryParam("date")
	reqDate := time.Now()
	if dateStr != "" {
		if parsed, err := time.Parse("2006-01-02", dateStr); err == nil {
			reqDate = parsed
		}
	}

	var branchIDs []uint
	if role == string(models.RoleManager) {
		if branchID != nil && *branchID > 0 {
			targetBranchID := uint(*branchID)
			if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, targetBranchID); err != nil {
				return response.Error(c, appErrors.ErrForbidden)
			}
		} else {
			var err error
			branchIDs, err = h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
			if err != nil {
				return response.HandleError(c, err)
			}
		}
	}

	result, err := h.service.GetOverview(branchIDs, branchID, reqDate)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Dashboard Recent Activities
// @Description Get recent check-ins, check-outs, and late arrivals
// @Tags Admin/Dashboard
// @Accept json
// @Produce json
// @Param branch_id query int false "Branch ID"
// @Param limit query int false "Limit"
// @Success 200 {object} response.Response{data=models.DashboardRecentActivityResponse} "Recent Activities"
// @Router /v1/admin/dashboard/recent-activities [get]
func (h *DashboardHandler) GetRecentActivities(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, _ := getContextUint(c, "branch_id")
	var branchID *int64
	if bID := c.QueryParam("branch_id"); bID != "" && bID != "0" {
		id, err := strconv.ParseInt(bID, 10, 64)
		if err == nil {
			branchID = &id
		}
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	reqDate := time.Now() // Fetch for today context

	var branchIDs []uint
	if role == string(models.RoleManager) {
		if branchID != nil && *branchID > 0 {
			targetBranchID := uint(*branchID)
			if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, targetBranchID); err != nil {
				return response.Error(c, appErrors.ErrForbidden)
			}
		} else {
			var err error
			branchIDs, err = h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
			if err != nil {
				return response.HandleError(c, err)
			}
		}
	}

	result, err := h.service.GetRecentActivities(branchIDs, branchID, reqDate, limit)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
