package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemGrade"
	"gorm.io/gorm"
)

type (
	ItemGradeHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	itemGradeHandler struct {
		service itemGrade.Service
	}
)

func (i *itemGradeHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := i.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve item grades", err)
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve item grades", &pagination)
}

func (i *itemGradeHandler) Create(ctx *fiber.Ctx) error {
	var req entities.ItemGrade
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}
	req.ItemID = uint(itemId)

	result, err := i.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create item grade", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully create item grade", result)
}

func (i *itemGradeHandler) Get(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item grade", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get item grade", result)
}

func (i *itemGradeHandler) Update(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	var updateData domain.UpdateRequestItemGrade
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	data := &entities.ItemGrade{ItemID: uint(itemId)}

	if updateData.Interior != nil {
		data.Interior = *updateData.Interior
	}
	if updateData.Exterior != nil {
		data.Exterior = *updateData.Exterior
	}
	if updateData.Frame != nil {
		data.Frame = *updateData.Frame
	}
	if updateData.Machine != nil {
		data.Machine = *updateData.Machine
	}

	_, err = i.service.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update item grade", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item grade", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully update item grade", result)
}

func (i *itemGradeHandler) Delete(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	if err := i.service.Delete(ctx.UserContext(), uint(itemId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "item not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete item grade", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully delete item grade", nil)
}

func NewItemGradeHandler(service itemGrade.Service) ItemGradeHandler {
	return &itemGradeHandler{service: service}
}
