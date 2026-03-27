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

func (h *AttendanceHandler) CheckIn(c echo.Context) error {
	var req models.CheckInRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID := c.Get("user_id").(uint)

	result, err := h.service.CheckIn(c.Request().Context(), userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

func (h *AttendanceHandler) CheckOut(c echo.Context) error {
	var req models.CheckOutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}
	if err := c.Validate(req); err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	userID := c.Get("user_id").(uint)

	result, err := h.service.CheckOut(c.Request().Context(), userID, req)
	if err != nil {
		return response.HandleError(c, err)
	}

	return response.Created(c, result)
}

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

func (h *AttendanceHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	_ = id
	// TODO: implement get by ID via service
	return response.Success(c, nil)
}
