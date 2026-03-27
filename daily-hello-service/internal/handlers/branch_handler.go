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
}

func NewBranchHandler(service *services.BranchService) *BranchHandler {
	return &BranchHandler{service: service}
}

func (h *BranchHandler) Create(c echo.Context) error {
	var req models.CreateBranchRequest
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

func (h *BranchHandler) GetByID(c echo.Context) error {
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

func (h *BranchHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
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

func (h *BranchHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.Error(c, appErrors.ErrInvalidInput)
	}

	if err := h.service.Delete(c.Request().Context(), uint(id)); err != nil {
		return response.HandleError(c, err)
	}

	return response.Success(c, map[string]string{"message": "Branch deleted successfully"})
}

func (h *BranchHandler) List(c echo.Context) error {
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
