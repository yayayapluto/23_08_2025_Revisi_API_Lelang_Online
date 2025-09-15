package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/auctionBidder"
	"gorm.io/gorm"
)

type (
	AuctionBidderHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	auctionBidderHandler struct {
		service auctionBidder.Service
	}
)

func (a *auctionBidderHandler) List(ctx *fiber.Ctx) error {
	rm := utils.GetRequestMeta(ctx)
	collection, total, err := a.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auction Bidders", err)
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieved Auction Bidders", &pagination)
}

func (a *auctionBidderHandler) Create(ctx *fiber.Ctx) error {
	var req entities.AuctionBidder
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	if req.UserID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "user_id is required", nil)
	}
	if req.AuctionID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "auction_id is required", nil)
	}
	if req.PhoneNumber == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "phone_number is required", nil)
	}
	if req.BankName == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "bank_name is required", nil)
	}
	if req.AccountNumber == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "account_number is required", nil)
	}
	if req.AccountName == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "account_name is required", nil)
	}

	res, err := a.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create Auction Bidder", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully created Auction Bidder", res)
}

func (a *auctionBidderHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction Bidder ID", err)
	}

	bidder, err := a.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auction Bidder", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieved Auction Bidder", bidder)
}

func (a *auctionBidderHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction Bidder ID", err)
	}

	var updateData entities.AuctionBidder
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}
	updateData.ID = uint(id)

	_, err = a.service.Update(ctx.UserContext(), &updateData)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update Auction Bidder", err)
	}

	result, err := a.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auction Bidder", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully updated Auction Bidder", result)
}

func (a *auctionBidderHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction Bidder ID", err)
	}

	if err := a.service.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "Auction Bidder not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete Auction Bidder", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Successfully deleted Auction Bidder", nil)
}

func NewAuctionBidderHandler(service auctionBidder.Service) AuctionBidderHandler {
	return &auctionBidderHandler{service: service}
}
