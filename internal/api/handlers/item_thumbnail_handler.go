package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemThumbnail"
	"gorm.io/gorm"
)

type (
	ItemThumbnailHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	itemThumbnailHandler struct {
		service itemThumbnail.Service
	}
)

func (i *itemThumbnailHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := i.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to list item thumbnails", err)
	}

	for i := range *collection {
		newUrlPath, err := utils.BuildFileURL(ctx, &(*collection)[i].File)
		if err != nil {
			return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
		}
		(*collection)[i].File.Path = *newUrlPath
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve item thumbnails", &pagination)
}

func (i *itemThumbnailHandler) Create(ctx *fiber.Ctx) error {
	var req entities.ItemThumbnail
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	header, err := ctx.FormFile("thumbnail")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	result, err := i.service.Create(ctx.UserContext(), &req, header)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create item grade", err)
	}

	newUrlPath, err := utils.BuildFileURL(ctx, &result.File)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	}
	result.File.Path = *newUrlPath

	return presenters.SuccessResponse(ctx, fiber.StatusCreated, "successfully create item grade", result)
}

func (i *itemThumbnailHandler) Get(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item thumbnail", err)
	}

	newUrlPath, err := utils.BuildFileURL(ctx, &result.File)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	}
	result.File.Path = *newUrlPath

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get item thumbnail", result)
}

func (i *itemThumbnailHandler) Update(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	var updateData domain.UpdateRequestItemThumbnail
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	data := &entities.ItemThumbnail{ItemID: uint(itemId)}
	if updateData.Name != nil {
		data.Name = *updateData.Name
	}

	_, err = i.service.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update item thumbnail", err)
	}

	result, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get item thumbnail", err)
	}

	newUrlPath, err := utils.BuildFileURL(ctx, &result.File)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	}
	result.File.Path = *newUrlPath

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully update item thumbnail", result)
}

func (i *itemThumbnailHandler) Delete(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	if err := i.service.Delete(ctx.UserContext(), uint(itemId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "item not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete item thumbnail", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully delete item thumbnail", nil)
}

func NewItemThumbnailHandler(service itemThumbnail.Service) ItemThumbnailHandler {
	return &itemThumbnailHandler{service: service}
}
