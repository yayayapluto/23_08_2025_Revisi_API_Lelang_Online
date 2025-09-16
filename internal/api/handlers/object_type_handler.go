package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	presenters "github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/objectType"
	"gorm.io/gorm"
)

type (
	ObjectTypeHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	objectTypeHandler struct {
		service objectType.Service
	}
)

func (o *objectTypeHandler) List(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	rm := utils.GetRequestMeta(ctx)

	OTs, total, err := o.service.List(ctx.UserContext(), rm.Offset, rm.Size, &search, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve object type list", err)
	}

	pagination := utils.BuildPagination(ctx, rm, *OTs, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieve object type list", &pagination)
}

func (o *objectTypeHandler) Create(ctx *fiber.Ctx) error {
	var ot entities.ObjectType
	if err := ctx.BodyParser(&ot); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	if ot.Name == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "name is required", nil)
	}

	if err := o.service.Create(ctx.UserContext(), &ot); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Object type name already exists", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create object type", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully created object type", &ot)
}

func (o *objectTypeHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid object type ID", err)
	}

	ot, err := o.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve object type", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieved object type", ot)
}

func (o *objectTypeHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid object type ID", err)
	}

	var updateData domain.UpdateRequestObjectType
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	ot := &entities.ObjectType{ID: uint(id)}
	if updateData.Name != nil {
		ot.Name = *updateData.Name
	}

	updatedOT, err := o.service.Update(ctx.UserContext(), ot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Object type name already exists", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update object type", err)
	}

	// Kembalikan data yang sudah diperbarui
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully updated object type", updatedOT)
}

func (o *objectTypeHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid object type ID", err)
	}

	if err := o.service.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "Object type not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete object type", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Successfully deleted object type", nil)
}

func NewObjectTypeHandler(service objectType.Service) ObjectTypeHandler {
	return &objectTypeHandler{service: service}
}
