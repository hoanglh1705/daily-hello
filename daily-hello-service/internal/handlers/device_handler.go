package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type DeviceHandler struct {
	service *services.DeviceService
	rbac    *services.RBACService
}

func NewDeviceHandler(service *services.DeviceService, rbac *services.RBACService) *DeviceHandler {
	return &DeviceHandler{service: service, rbac: rbac}
}

// @Summary Register Device
// @Description Register the current user's device. Returns existing record if already registered.
// @Tags Device
// @Accept json
// @Produce json
// @Param request body models.RegisterDeviceRequest true "Device data"
// @Success 201 {object} response.Response{data=models.Device} "Device registered"
// @Failure 400 {object} response.Response "Invalid input"
// @Router /v1/devices/register [post]
func (h *DeviceHandler) Register(c echo.Context) error {
	userID := uint(c.Get("user_id").(float64))

	var req models.RegisterDeviceRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Register(c.Request().Context(), userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

// @Summary Get My Device Status
// @Description Get the registration status of the current user's device
// @Tags Device
// @Produce json
// @Param device_id query string true "Device ID"
// @Success 200 {object} response.Response{data=models.Device} "Device status"
// @Failure 404 {object} response.Response "Device not found"
// @Router /v1/devices/status [get]
func (h *DeviceHandler) GetStatus(c echo.Context) error {
	userID := uint(c.Get("user_id").(float64))

	var q models.DeviceStatusQuery
	if err := c.Bind(&q); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(q); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.GetStatus(c.Request().Context(), userID, q.DeviceID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary List My Devices
// @Description List all devices registered by the current user
// @Tags Device
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "List of devices"
// @Router /v1/devices [get]
func (h *DeviceHandler) ListMyDevices(c echo.Context) error {
	userID := uint(c.Get("user_id").(float64))

	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.GetByUserID(c.Request().Context(), userID, pq)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary List Devices (Admin)
// @Description List all devices, optionally filtered by status. Admin/Manager only.
// @Tags Device
// @Produce json
// @Param status query string false "Filter by status: pending, approved, rejected"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "List of devices"
// @Failure 403 {object} response.Response "Forbidden"
// @Router /v1/admin/devices [get]
func (h *DeviceHandler) AdminList(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var q models.DeviceListQuery
	if err := c.Bind(&q); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var branchIDs []uint
	if role == string(models.RoleManager) {
		if q.BranchID != nil {
			if err := h.rbac.EnsureBranchAccess(c.Request().Context(), role, currentBranchID, *q.BranchID); err != nil {
				return response.Error(c, appErrors.ErrForbidden)
			}
		} else {
			branchIDs, err = h.rbac.GetAllowedBranchIDs(c.Request().Context(), role, currentBranchID)
			if err != nil {
				return response.HandleError(c, err)
			}
		}
	}

	result, err := h.service.ListByStatus(c.Request().Context(), q.Status, q.BranchID, branchIDs, pq)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Approve Device (Admin)
// @Description Approve a device registration request. Admin/Manager only.
// @Tags Device
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} response.Response{data=models.Device} "Device approved"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Device not found"
// @Router /v1/admin/devices/{id}/approve [put]
func (h *DeviceHandler) Approve(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	adminID := uint(c.Get("user_id").(float64))

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if role == string(models.RoleManager) {
		if err := h.rbac.EnsureDeviceAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
	}

	result, err := h.service.Approve(c.Request().Context(), adminID, uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Reject Device (Admin)
// @Description Reject a device registration request. Admin/Manager only.
// @Tags Device
// @Produce json
// @Param id path int true "Device ID"
// @Success 200 {object} response.Response{data=models.Device} "Device rejected"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Device not found"
// @Router /v1/admin/devices/{id}/reject [put]
func (h *DeviceHandler) Reject(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}
	currentBranchID, err := getContextUint(c, "branch_id")
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	adminID := uint(c.Get("user_id").(float64))

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if role == string(models.RoleManager) {
		if err := h.rbac.EnsureDeviceAccess(c.Request().Context(), role, currentBranchID, uint(id)); err != nil {
			return response.Error(c, appErrors.ErrForbidden)
		}
	}

	result, err := h.service.Reject(c.Request().Context(), adminID, uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
