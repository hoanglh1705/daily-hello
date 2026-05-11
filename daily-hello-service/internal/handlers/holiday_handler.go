package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type HolidayHandler struct {
	service *services.HolidayService
}

func NewHolidayHandler(service *services.HolidayService) *HolidayHandler {
	return &HolidayHandler{service: service}
}

// @Summary Create Holiday
// @Description Create a new holiday (admin/manager only)
// @Tags Holiday
// @Accept json
// @Produce json
// @Param request body models.CreateHolidayRequest true "Holiday data"
// @Success 201 {object} response.Response{data=models.Holiday} "Create Holiday successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/admin/holidays [post]
func (h *HolidayHandler) Create(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

	var req models.CreateHolidayRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID, err := getContextUint(c, "user_id")
	if err != nil || userID == nil {
		return response.Error(c, appErrors.ErrUnauthorized)
	}

	result, err := h.service.Create(c.Request().Context(), req, *userID)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

// @Summary List Holidays (Admin)
// @Description List all holidays with pagination (admin/manager only)
// @Tags Holiday
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "List Holidays successfully"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/admin/holidays [get]
func (h *HolidayHandler) AdminList(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

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

// @Summary Update Holiday
// @Description Update an existing holiday (admin/manager only)
// @Tags Holiday
// @Accept json
// @Produce json
// @Param id path int true "Holiday ID"
// @Param request body models.UpdateHolidayRequest true "Holiday data"
// @Success 200 {object} response.Response{data=models.Holiday} "Update Holiday successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Holiday not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/admin/holidays/{id} [put]
func (h *HolidayHandler) Update(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var req models.UpdateHolidayRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.Update(c.Request().Context(), uint(id), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Delete Holiday
// @Description Delete a holiday (admin/manager only)
// @Tags Holiday
// @Accept json
// @Produce json
// @Param id path int true "Holiday ID"
// @Success 200 {object} response.Response "Delete Holiday successfully"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Holiday not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/admin/holidays/{id} [delete]
func (h *HolidayHandler) Delete(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != string(models.RoleAdmin) && role != string(models.RoleManager) {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if err := h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "Holiday deleted successfully"})
}

// @Summary Get Holidays by Month
// @Description Get holidays for a specific month (accessible by all authenticated users)
// @Tags Holiday
// @Accept json
// @Produce json
// @Param year query int true "Year"
// @Param month query int true "Month (1-12)"
// @Success 200 {object} response.Response{data=[]models.Holiday} "Get Holidays successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/holidays [get]
func (h *HolidayHandler) GetByMonth(c echo.Context) error {
	var filter models.HolidayFilter
	if err := c.Bind(&filter); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if filter.Year <= 0 || filter.Month < 1 || filter.Month > 12 {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.GetByMonth(c.Request().Context(), filter.Year, filter.Month)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
