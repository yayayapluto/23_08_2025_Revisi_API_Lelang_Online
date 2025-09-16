package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/bid"
	"gorm.io/gorm"
)

type (
	BidHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	bidHandler struct {
		service bid.Service
	}
)

func (b *bidHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := b.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve bid list", err)

	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve bid list", &pagination)
}

func (b *bidHandler) Create(ctx *fiber.Ctx) error {
	var req entities.Bid

	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}

	if req.BidderID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "bidder id required", nil)
	}

	if req.Value <= 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "value required", nil)
	}

	res, err := b.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to create bid", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully create bid", res)
}

func (b *bidHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "id required", nil)
	}

	bid, err := b.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get bid", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get bid", bid)
}

func (b *bidHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "id required", nil)
	}

	var updateData entities.Bid
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}
	updateData.ID = uint(id)

	_, err = b.service.Update(ctx.UserContext(), &updateData)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update bid", err)
	}

	result, err := b.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get bid", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully update bid", result)
}

func (b *bidHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "id required", nil)
	}

	if err = b.service.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "bid does not exist", nil)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to delete bid", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully delete bid", nil)
}

func NewBidHandler(service bid.Service) BidHandler {
	return &bidHandler{service: service}
}
