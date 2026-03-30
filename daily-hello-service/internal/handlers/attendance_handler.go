package handlers

import (
	"strconv"

	"daily-hello-service/internal/models"
	appErrors "daily-hello-service/internal/pkg/errors"
	"daily-hello-service/internal/pkg/response"
	"daily-hello-service/internal/services"

	"github.com/labstack/echo/v4"
)

type AttendanceHandler struct {
	service *services.AttendanceService
}

func NewAttendanceHandler(service *services.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service: service}
}

// @Summary Check In
// @Description Check in to a branch
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body models.AttendanceRequest true "Attendance request data"
// @Success 201 {object} response.Response{data=models.Attendance} "Check in successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/attendance/check-in [post]
func (h *AttendanceHandler) CheckIn(c echo.Context) error {
	var req models.AttendanceRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID := c.Get("user_id").(float64)
	branchID := c.Get("branch_id").(float64)
	req.BranchID = uint(branchID)

	result, err := h.service.CheckIn(c.Request().Context(), uint(userID), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Check Out
// @Description Check out from a branch
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body models.AttendanceRequest true "Attendance request data"
// @Success 201 {object} response.Response{data=models.Attendance} "Check out successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/attendance/check-out [post]
func (h *AttendanceHandler) CheckOut(c echo.Context) error {
	var req models.AttendanceRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID := c.Get("user_id").(float64)
	branchID := c.Get("branch_id").(float64)
	req.BranchID = uint(branchID)

	result, err := h.service.CheckOut(c.Request().Context(), uint(userID), req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

// @Summary Get Attendance History
// @Description Get attendance history
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body models.AttendanceFilter true "Attendance filter"
// @Success 200 {object} response.Response{data=models.PaginatedResponse} "Get attendance history successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/attendance/history [get]
func (h *AttendanceHandler) GetHistory(c echo.Context) error {
	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var filter models.AttendanceFilter
	if err := c.Bind(&filter); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	var dateRange struct {
		From string `query:"from"`
		To   string `query:"to"`
	}
	if err := c.Bind(&dateRange); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if filter.DateFrom == "" {
		filter.DateFrom = dateRange.From
	}
	if filter.DateTo == "" {
		filter.DateTo = dateRange.To
	}

	result, err := h.service.GetHistory(c.Request().Context(), filter, pq)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Get Attendance by ID
// @Description Get attendance by ID
// @Tags Attendance
// @Accept json
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} response.Response{data=models.Attendance} "Get attendance successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/attendance/{id} [get]
func (h *AttendanceHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	_ = id
	// TODO: implement get by ID via service
	return response.Success(c, nil)
}

// @Summary Get Today Attendance
// @Description Get attendance record for today
// @Tags Attendance
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=models.Attendance} "Get today attendance successfully"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/attendance/today [get]
func (h *AttendanceHandler) GetToday(c echo.Context) error {
	userID := c.Get("user_id").(float64)

	result, err := h.service.GetToday(c.Request().Context(), uint(userID))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Approve Check In (Admin)
// @Description Approve pending check-in status. Admin/Manager only.
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} response.Response{data=models.Attendance} "Check-in approved"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Attendance not found"
// @Router /v1/admin/attendance/{id}/check-in/approve [put]
func (h *AttendanceHandler) ApproveCheckIn(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.ApproveCheckIn(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Reject Check In (Admin)
// @Description Reject pending check-in status. Admin/Manager only.
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} response.Response{data=models.Attendance} "Check-in rejected"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Attendance not found"
// @Router /v1/admin/attendance/{id}/check-in/reject [put]
func (h *AttendanceHandler) RejectCheckIn(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.RejectCheckIn(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Approve Check Out (Admin)
// @Description Approve pending check-out status. Admin/Manager only.
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} response.Response{data=models.Attendance} "Check-out approved"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Attendance not found"
// @Router /v1/admin/attendance/{id}/check-out/approve [put]
func (h *AttendanceHandler) ApproveCheckOut(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.ApproveCheckOut(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}

// @Summary Reject Check Out (Admin)
// @Description Reject pending check-out status. Admin/Manager only.
// @Tags Attendance
// @Produce json
// @Param id path int true "Attendance ID"
// @Success 200 {object} response.Response{data=models.Attendance} "Check-out rejected"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 403 {object} response.Response "Forbidden"
// @Failure 404 {object} response.Response "Attendance not found"
// @Router /v1/admin/attendance/{id}/check-out/reject [put]
func (h *AttendanceHandler) RejectCheckOut(c echo.Context) error {
	role, _ := c.Get("role").(string)
	if role != "admin" && role != "manager" {
		return response.Error(c, appErrors.ErrForbidden)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	result, err := h.service.RejectCheckOut(c.Request().Context(), uint(id))
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, result)
}
