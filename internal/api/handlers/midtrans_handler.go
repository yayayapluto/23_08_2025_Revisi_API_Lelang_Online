package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/bidderPayment"
	midtrans2 "github.com/yayayapluto/revisi_api_lelang_online/pkg/midtrans"
)

type (
	MidtransHandler interface {
		PaymentHandler(ctx *fiber.Ctx) error
	}

	midtrans struct {
		midtransService midtrans2.MidtransService
		paymentService  bidderPayment.Service
	}
)

func (m *midtrans) PaymentHandler(ctx *fiber.Ctx) error {
	orderID := ctx.Params("order_id")

	success, err := m.midtransService.VerifyPayment(ctx.UserContext(), orderID)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to verify payment", err)
	}
	if success {
		err = m.paymentService.ConfirmedPayment(ctx.UserContext(), orderID)
		return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully confirmed payment", nil)
	}

	return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to verify payment", err)
}

func NewMidtransHandler(midtransService midtrans2.MidtransService, service bidderPayment.Service) MidtransHandler {
	return &midtrans{
		midtransService: midtransService,
		paymentService:  service,
	}
}
