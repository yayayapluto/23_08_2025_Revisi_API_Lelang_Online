package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/bidderPayment"
)

type (
	BidderPaymentHandler interface {
		InitializePayment(ctx *fiber.Ctx) error
		CheckBidderPayment(ctx *fiber.Ctx) error
	}

	bidderPaymentHandler struct {
		bidderPaymentService bidderPayment.Service
	}
)

func (b *bidderPaymentHandler) CheckBidderPayment(ctx *fiber.Ctx) error {
	userID := ctx.QueryInt("user_id", 0)
	auctionID := ctx.QueryInt("auction_id", 0)

	if userID == 0 || auctionID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "user id and/or auction id are required", nil)
	}

	res, err := b.bidderPaymentService.FindByBidderData(ctx.UserContext(), uint(auctionID), uint(userID))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "error getting bidder payment", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "success getting bidder payment", res)
}

func (b *bidderPaymentHandler) InitializePayment(ctx *fiber.Ctx) error {
	var req domain.BidderPaymentRequest
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	res, err := b.bidderPaymentService.InitializePayment(ctx.UserContext(), req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to initialize payment", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully initialized payment", res)
}

func NewBidderPaymentHandler(service bidderPayment.Service) BidderPaymentHandler {
	return &bidderPaymentHandler{bidderPaymentService: service}
}
