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
	}

	bidderPaymentHandler struct {
		bidderPaymentService bidderPayment.Service
	}
)

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
