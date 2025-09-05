package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils/pagination"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/item"
	"gorm.io/gorm"
)

type (
	ItemHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	itemHandler struct {
		s item.Service
	}
)

func (i *itemHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	items, total, err := i.s.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.Search, &rm.SortDir, &rm.SortBy)

	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve item list", err)
	}

	//for i := range *items {
	//	newUrlPath, err := utils.BuildFileURL(ctx, &(*items)[i].File)
	//	if err != nil {
	//		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	//	}
	//	(*items)[i].File.Path = *newUrlPath
	//}

	paginationRes := utils.BuildPagination[entities.Item](ctx, rm, *items, total)
	return presenters.SuccessResponse[pagination.ResponseMetaData[entities.Item]](ctx, fiber.StatusOK, "successfully retrieve item list", &paginationRes)
}

func (i *itemHandler) Create(ctx *fiber.Ctx) error {
	var itemReq entities.Item
	if err := ctx.BodyParser(&itemReq); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	header, err := ctx.FormFile("thumbnail")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	itemRes, err := i.s.Create(ctx.UserContext(), &itemReq, header)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusConflict, "failed to create new item", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to create new item", err)
	}

	//newUrlPath, err := utils.BuildFileURL(ctx, &itemRes.File)
	//if err != nil {
	//	return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	//}
	//itemRes.File.Path = *newUrlPath

	return presenters.SuccessResponse[entities.Item](ctx, fiber.StatusCreated, "successfully created new item", itemRes)
}

func (i *itemHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	itemFound, err := i.s.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get item detail", err)
	}

	//newUrlPath, err := utils.BuildFileURL(ctx, &itemFound.File)
	//if err != nil {
	//	return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	//}
	//itemFound.File.Path = *newUrlPath

	return presenters.SuccessResponse[entities.Item](ctx, fiber.StatusOK, "successfully retrieved item detail", itemFound)
}

func (i *itemHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	var updateData domain.UpdateRequestItem
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	itemFound := &entities.Item{ID: uint(id)}
	if updateData.ObjectTypeID != nil {
		itemFound.ObjectTypeID = *updateData.ObjectTypeID
	}
	if updateData.Name != nil {
		itemFound.Name = *updateData.Name
	}
	if updateData.Price != nil {
		itemFound.Price = *updateData.Price
	}
	if updateData.DepositPrice != nil {
		itemFound.DepositPrice = *updateData.DepositPrice
	}
	if updateData.Description != nil {
		itemFound.Description = updateData.Description
	}
	if updateData.FileID != nil {
		itemFound.FileID = *updateData.FileID
	}

	itemRes, err := i.s.Update(ctx.UserContext(), itemFound)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to update item detail", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update item detail", err)
	}
	itemRes.ID = uint(id)

	//newUrlPath, err := utils.BuildFileURL(ctx, &itemRes.File)
	//if err != nil {
	//	return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to build file url", err)
	//}
	//itemRes.File.Path = *newUrlPath

	return presenters.SuccessResponse[entities.Item](ctx, fiber.StatusOK, "successfully updated item detail", itemRes)
}

func (i *itemHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	if err := i.s.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "failed to delete object type", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to delete item detail", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully deleted item detail", nil)
}

func NewItemHandler(service item.Service) ItemHandler {
	return &itemHandler{s: service}
}
