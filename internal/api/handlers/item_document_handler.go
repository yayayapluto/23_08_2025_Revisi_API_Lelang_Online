package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDocument"
	"gorm.io/gorm"
)

type (
	ItemDocumentHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	itemDocumentHandler struct {
		service itemDocument.Service
	}
)

func (i *itemDocumentHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := i.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve item documents", err)
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve item documents", &pagination)
}

func (i *itemDocumentHandler) Create(ctx *fiber.Ctx) error {
	var req entities.ItemDocument
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
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create item document", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully create item document", result)
}

func (i *itemDocumentHandler) Get(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item document", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get item document", result)
}

func (i *itemDocumentHandler) Update(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	var updateData domain.UpdateRequestItemDocument
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	data := &entities.ItemDocument{ItemID: uint(itemId)}

	if updateData.Bpkb != nil {
		data.Bpkb = updateData.Bpkb
	}
	if updateData.Stnk != nil {
		data.Stnk = updateData.Stnk
	}
	if updateData.Facture != nil {
		data.Facture = updateData.Facture
	}
	if updateData.Receipt != nil {
		data.Receipt = updateData.Receipt
	}
	if updateData.OwnershipRelease != nil {
		data.OwnershipRelease = updateData.OwnershipRelease
	}
	if updateData.Warranty != nil {
		data.Warranty = updateData.Warranty
	}
	if updateData.Box != nil {
		data.Box = updateData.Box
	}

	_, err = i.service.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update item document", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item document", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully update item document", result)
}

func (i *itemDocumentHandler) Delete(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	if err := i.service.Delete(ctx.UserContext(), uint(itemId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "item not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete item document", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully delete item document", nil)
}

func NewItemDocumentHandler(service itemDocument.Service) ItemDocumentHandler {
	return &itemDocumentHandler{service: service}
}
