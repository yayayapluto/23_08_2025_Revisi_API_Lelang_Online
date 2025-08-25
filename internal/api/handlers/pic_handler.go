package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/pic"
	"gorm.io/gorm"
)

type (
	PICHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	picHandler struct {
		service pic.Service
	}
)

func (p *picHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := p.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.Search, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve PICs", err)
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieve PICs", &pagination)
}

func (p *picHandler) Create(ctx *fiber.Ctx) error {
	var req entities.PIC
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	if req.Name == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Name is required", nil)
	}
	if req.PhoneNumber == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Phone number is required", nil)
	}

	res, err := p.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create PIC", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully created PICs", res)
}

func (p *picHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid PIC ID", err)
	}

	result, err := p.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve PIC", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieved PIC", result)
}

func (p *picHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid PIC ID", err)
	}

	var updateData domain.UpdateRequestPIC
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	data := &entities.PIC{ID: uint(id)}
	if updateData.Name != nil {
		data.Name = *updateData.Name
	}
	if updateData.PhoneNumber != nil {
		data.PhoneNumber = *updateData.PhoneNumber
	}

	_, err = p.service.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update PIC", err)
	}

	result, err := p.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve PIC", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully updated PIC", result)
}

func (p *picHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid PIC ID", err)
	}

	if err := p.service.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "PIC not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete PIC", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Successfully delete PIC", nil)
}

func NewPICHandler(service pic.Service) PICHandler {
	return &picHandler{service: service}
}
