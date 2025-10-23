package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/auctionBidder"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/bid"
	"gorm.io/gorm"
	"log"
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
		service       bid.Service
		bidderService auctionBidder.Service
	}
)

func (b *bidHandler) List(ctx *fiber.Ctx) error {
	var auctionID *uint

	if ctx.QueryInt("auction_id") != 0 {
		queryAuctionID := uint(ctx.QueryInt("auction_id"))
		auctionID = &queryAuctionID
		log.Println("[handler] auction_id", *auctionID)
	} else {
		auctionID = nil
	}

	rm := utils.GetRequestMeta(ctx)
	collection, total, err := b.service.List(ctx.UserContext(), rm.Offset, rm.Size, &rm.SortDir, &rm.SortBy, auctionID)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve bid list", err)

	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully retrieve bid list", &pagination)
}

func (b *bidHandler) Create(ctx *fiber.Ctx) error {
	user := ctx.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userIDClaim, ok := claims["user_id"].(float64)
	// map[email:farras@email.com exp:1.76145091e+09 role:user user_id:1 username:farras]
	if !ok {
		return presenters.ErrorResponse(ctx, fiber.StatusForbidden, "user_id not found in token", nil)
	}

	auctionID := ctx.QueryInt("auction_id")

	// {{api_url}}/bids?auction_id=8
	if userIDClaim == 0 || auctionID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "user_id and auction_id are required", nil)
	}

	bidder, err := b.bidderService.GetByUserIDAndAuctionID(ctx.UserContext(), uint(userIDClaim), uint(auctionID))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve bid list", err)
	}

	if bidder == nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "cannot find bidder data", nil)
	}

	var req entities.Bid
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse request body", err)
	}
	if req.Value <= 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "value required", nil)
	}

	req.BidderID = bidder.ID

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

func NewBidHandler(service bid.Service, bidderService auctionBidder.Service) BidHandler {
	return &bidHandler{service: service, bidderService: bidderService}
}
