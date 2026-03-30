package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type BranchWifiHandler struct {
	service *services.BranchWifiService
}

func NewBranchWifiHandler(service *services.BranchWifiService) *BranchWifiHandler {
	return &BranchWifiHandler{service: service}
}

// @Summary Create Branch Wifi
// @Description Create a new wifi configuration for a branch
// @Tags BranchWifi
// @Accept json
// @Produce json
// @Param request body models.CreateBranchWifiRequest true "Branch Wifi data"
// @Success 201 {object} response.Response{data=models.BranchWifi} "Branch Wifi created successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/branch-wifi [post]
func (h *BranchWifiHandler) Create(c echo.Context) error {
	var req models.CreateBranchWifiRequest
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

// @Summary Get Branch Wifi by ID
// @Description Get a single branch wifi by its ID
// @Tags BranchWifi
// @Produce json
// @Param id path int true "Branch Wifi ID"
// @Success 200 {object} response.Response{data=models.BranchWifi} "Branch Wifi details"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/branch-wifi/{id} [get]
func (h *BranchWifiHandler) GetByID(c echo.Context) error {
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

// @Summary List Branch Wifi by Branch ID
// @Description Get all wifi configurations for a specific branch
// @Tags BranchWifi
// @Produce json
// @Param branch_id path int true "Branch ID"
// @Success 200 {object} response.Response{data=models.BranchWifi} "List of branch wifi"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/branch-wifi/branch/{branch_id} [get]
func (h *BranchWifiHandler) GetByBranchID(c echo.Context) error {
	branchID, err := strconv.ParseUint(c.Param("branch_id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.GetByBranchID(c.Request().Context(), uint(branchID))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Update Branch Wifi
// @Description Update an existing branch wifi configuration
// @Tags BranchWifi
// @Accept json
// @Produce json
// @Param id path int true "Branch Wifi ID"
// @Param request body models.UpdateBranchWifiRequest true "Branch Wifi data"
// @Success 200 {object} response.Response{data=models.BranchWifi} "Branch Wifi updated successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/branch-wifi/{id} [put]
func (h *BranchWifiHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var req models.UpdateBranchWifiRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Delete Branch Wifi
// @Description Delete a branch wifi configuration
// @Tags BranchWifi
// @Produce json
// @Param id path int true "Branch Wifi ID"
// @Success 200 {object} response.Response{data=string} "Branch Wifi deleted successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 404 {object} response.Response "Not found"
// @Router /v1/branch-wifi/{id} [delete]
func (h *BranchWifiHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if err := h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "Branch wifi deleted successfully"})
}
