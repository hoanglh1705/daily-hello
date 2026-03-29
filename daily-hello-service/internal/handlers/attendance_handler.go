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
// @Param request body models.CheckInRequest true "Check in data"
// @Success 201 {object} response.Response{data=models.Attendance} "Check in successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/attendance/check-in [post]
func (h *AttendanceHandler) CheckIn(c echo.Context) error {
	var req models.CheckInRequest
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

	return response.Created(c, result)
}

// @Summary Check Out
// @Description Check out from a branch
// @Tags Attendance
// @Accept json
// @Produce json
// @Param request body models.CheckOutRequest true "Check out data"
// @Success 201 {object} response.Response{data=models.Attendance} "Check out successfully"
// @Failure 400 {object} response.Response "Invalid input"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/auth/attendance/check-out [post]
func (h *AttendanceHandler) CheckOut(c echo.Context) error {
	var req models.CheckOutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID := c.Get("user_id").(float64)

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
// @Router /v1/auth/attendance/history [get]
func (h *AttendanceHandler) GetHistory(c echo.Context) error {
	var pq models.PaginationQuery
	if err := c.Bind(&pq); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	var filter models.AttendanceFilter
	if err := c.Bind(&filter); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
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
