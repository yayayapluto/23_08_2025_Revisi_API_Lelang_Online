package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/auction"
	"gorm.io/gorm"
	"time"
)

type (
	AuctionHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	auctionHandler struct {
		service auction.Service
	}
)

func (p *auctionHandler) List(ctx *fiber.Ctx) error {
	startDateQuery := ctx.Query("startDate")
	endDateQuery := ctx.Query("startDate")

	var startDate *time.Time
	if startDateQuery != "" {
		res, err := time.Parse(time.DateOnly, startDateQuery)
		if err != nil {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse date", err)
		}
		startDate = &res
	}

	var endDate *time.Time
	if endDateQuery != "" {
		res, err := time.Parse(time.DateOnly, endDateQuery)
		if err != nil {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse date", err)
		}
		endDate = &res
	}

	rm := utils.GetRequestMeta(ctx)
	collection, total, err := p.service.List(ctx.UserContext(), rm.Offset, rm.Size, startDate, endDate, &rm.SortDir, &rm.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auctions", err)
	}

	for i := range *collection {
		auc := &(*collection)[i]

		if auc.Item.File.ID != 0 {
			if url, err := utils.BuildFileURL(ctx, &auc.Item.File); err == nil {
				auc.Item.File.Path = *url
			}
		}

		if auc.Item.ItemThumbnails != nil {
			for j := range *auc.Item.ItemThumbnails {
				thumb := &(*auc.Item.ItemThumbnails)[j].File
				if thumb.ID != 0 {
					if url, err := utils.BuildFileURL(ctx, thumb); err == nil {
						thumb.Path = *url
					}
				}
			}
		}
	}

	pagination := utils.BuildPagination(ctx, rm, *collection, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieve Auctions", &pagination)
}

func (p *auctionHandler) Create(ctx *fiber.Ctx) error {
	var req entities.Auction
	user, err := utils.GetUserFromToken(ctx)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get user from token", err)
	}
	if user.Role == "organizer" {
		req.OrganizerID = user.ID
	}

	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	if req.ItemID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "item_id is required", nil)
	}
	if req.OrganizerID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "organizer_id is required", nil)
	}
	if req.PicID == 0 {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "pic_id is required", nil)
	}

	res, err := p.service.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create Auction", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully created Auctions", res)
}

func (p *auctionHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction ID", err)
	}

	auc, err := p.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auction", err)
	}

	if auc.Item.File.ID != 0 {
		if url, err := utils.BuildFileURL(ctx, &auc.Item.File); err == nil {
			auc.Item.File.Path = *url
		}
	}

	if auc.Item.ItemThumbnails != nil {
		for j := range *auc.Item.ItemThumbnails {
			thumb := &(*auc.Item.ItemThumbnails)[j].File
			if thumb.ID != 0 {
				if url, err := utils.BuildFileURL(ctx, thumb); err == nil {
					thumb.Path = *url
				}
			}
		}
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully retrieved Auction", auc)
}

func (p *auctionHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction ID", err)
	}

	var updateData domain.UpdateRequestAuction

	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Failed to parse request body", err)
	}

	data := &entities.Auction{ID: uint(id)}
	if updateData.ItemID != nil {
		data.ItemID = *updateData.ItemID
	}
	if updateData.OrganizerID != nil {
		data.OrganizerID = *updateData.OrganizerID
	}
	if updateData.PicID != nil {
		data.PicID = *updateData.PicID
	}
	if updateData.StartDate != nil {
		data.StartDate = *updateData.StartDate
	}
	if updateData.EndDate != nil {
		data.EndDate = *updateData.EndDate
	}

	_, err = p.service.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update Auction", err)
	}

	result, err := p.service.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to retrieve Auction", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Successfully updated Auction", result)
}

func (p *auctionHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid Auction ID", err)
	}

	if err := p.service.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "Auction not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete Auction", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Successfully delete Auction", nil)
}

func NewAuctionHandler(service auction.Service) AuctionHandler {
	return &auctionHandler{service: service}
}
