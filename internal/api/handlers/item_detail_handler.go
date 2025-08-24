package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/itemDetail"
	"gorm.io/gorm"
)

type (
	ItemDetailHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	itemDetailHandler struct {
		service itemDetail.Service
	}
)

func (i *itemDetailHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	itemDetails, total, err := i.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.Search, &rm.SortDir, &rm.SortBy)

	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve item details", err)
	}

	paginationRes := utils.BuildPagination[entities.ItemDetail](ctx, rm, *itemDetails, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve item details", &paginationRes)
}

func (i *itemDetailHandler) Create(ctx *fiber.Ctx) error {
	var itemDetailReq entities.ItemDetail
	if err := ctx.BodyParser(&itemDetailReq); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	itemDetailRes, err := i.service.Create(ctx.UserContext(), &itemDetailReq)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to create item detail", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully created item details", &itemDetailRes)
}

func (i *itemDetailHandler) Get(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	itemDetailFound, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get item detail", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieved item details", &itemDetailFound)
}

func (i *itemDetailHandler) Update(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	var updateData domain.UpdateRequestItemDetail
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	itemDetailFound := &entities.ItemDetail{ItemID: uint(itemId)}

	if updateData.PlateNumber != nil {
		itemDetailFound.PlateNumber = updateData.PlateNumber
	}
	if updateData.Brand != nil {
		itemDetailFound.Brand = *updateData.Brand
	}
	if updateData.Series != nil {
		itemDetailFound.Series = updateData.Series
	}
	if updateData.CC != nil {
		itemDetailFound.CC = updateData.CC
	}
	if updateData.Type != nil {
		itemDetailFound.Type = updateData.Type
	}
	if updateData.Transmission != nil {
		itemDetailFound.Transmission = updateData.Transmission
	}
	if updateData.Model != nil {
		itemDetailFound.Model = updateData.Model
	}
	if updateData.Year != nil {
		itemDetailFound.Year = *updateData.Year
	}
	if updateData.FrameNumber != nil {
		itemDetailFound.FrameNumber = updateData.FrameNumber
	}
	if updateData.MachineNumber != nil {
		itemDetailFound.MachineNumber = updateData.MachineNumber
	}
	if updateData.Kilometer != nil {
		itemDetailFound.Kilometer = updateData.Kilometer
	}
	if updateData.Fuel != nil {
		itemDetailFound.Fuel = updateData.Fuel
	}
	if updateData.Color != nil {
		itemDetailFound.Color = *updateData.Color
	}
	if updateData.DriveType != nil {
		itemDetailFound.DriveType = updateData.DriveType
	}
	if updateData.StnkDate != nil {
		itemDetailFound.StnkDate = updateData.StnkDate
	}

	itemRes, err := i.service.Update(ctx.UserContext(), itemDetailFound)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update item detail", err)
	}
	itemRes.ItemID = uint(itemId)

	itemDetailResFound, err := i.service.Get(ctx.UserContext(), uint(itemId))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get item detail", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully updated item details", &itemDetailResFound)
}

func (i *itemDetailHandler) Delete(ctx *fiber.Ctx) error {
	itemId, err := ctx.ParamsInt("itemID")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid item_id param", err)
	}

	if err := i.service.Delete(ctx.UserContext(), uint(itemId)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "item not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to delete item detail", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully deleted item details", nil)
}

func NewItemDetailHandler(service itemDetail.Service) ItemDetailHandler {
	return &itemDetailHandler{service: service}
}
