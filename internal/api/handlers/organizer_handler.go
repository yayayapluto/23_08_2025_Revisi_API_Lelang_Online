package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils/pagination"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/organizer"
	"gorm.io/gorm"
	"math"
)

type (
	OrganizerHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	organizerHandler struct {
		s organizer.Service
	}
)

func (o *organizerHandler) List(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	page := ctx.QueryInt("page", 1)

	size := ctx.QueryInt("size")
	size = int(math.Min(math.Max(float64(size), 10), 100)) // min 10, max 100

	offset := (page - 1) * size // ex: page=2 size=10 -> (2 - 1) * 10 -> 10 | that means it start from offset 10

	sortBy := ctx.Query("sortBy", "id")
	sortDir := ctx.Query("sortDir", "asc")

	OTs, total, err := o.s.List(ctx.UserContext(), offset, size, &search, &sortDir, &sortBy) // OTs => organizers | plural
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve organizer list", err)
	}

	currentPageUrl := pagination.BuildPageURL(ctx, search, page, size, sortBy, sortDir)
	firstPageUrl := pagination.BuildPageURL(ctx, search, 1, size, sortBy, sortDir)

	var nextPageUrl, prevPageUrl *string
	if offset+len(*OTs) < int(total) {
		url := pagination.BuildPageURL(ctx, search, page+1, size, sortBy, sortDir)
		nextPageUrl = &url
	}

	if page > 1 {
		url := pagination.BuildPageURL(ctx, search, page-1, size, sortBy, sortDir)
		prevPageUrl = &url
	}

	totalPage := int(math.Ceil(float64(total) / float64(size)))

	paginationRes := pagination.NewResponseMetaData[entities.Organizer](page, currentPageUrl, *OTs, firstPageUrl, nextPageUrl, size, prevPageUrl, totalPage)
	return presenters.SuccessResponse[pagination.ResponseMetaData[entities.Organizer]](ctx, fiber.StatusOK, "successfully retrieve organizer list", &paginationRes)
}

func (o *organizerHandler) Create(ctx *fiber.Ctx) error {
	var or entities.Organizer
	if err := ctx.BodyParser(&or); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	if err := o.s.Create(ctx.UserContext(), &or); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to create new organizer", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to create new organizer", err)
	}

	return presenters.SuccessResponse[entities.Organizer](ctx, fiber.StatusOK, "successfully create new organizer", &or)
}

func (o *organizerHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	or, err := o.s.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get organizer detail", err)
	}

	return presenters.SuccessResponse[entities.Organizer](ctx, fiber.StatusOK, "successfully get organizer detail", or)
}

func (o *organizerHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	var od domain.UpdateRequestOrganizer
	if err := ctx.BodyParser(&od); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	or := &entities.Organizer{ID: uint(id)}

	if od.Name != nil {
		or.Name = *od.Name
	}

	if od.Address != nil {
		or.Address = *od.Address
	}

	if od.BankName != nil {
		or.BankName = *od.BankName
	}

	if od.AccountNumber != nil {
		or.AccountNumber = *od.AccountNumber
	}

	if od.AccountName != nil {
		or.AccountName = *od.AccountName
	}

	orRes, err := o.s.Update(ctx.UserContext(), or)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to update organizer", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update organizer", err)
	}
	orRes.ID = uint(id)

	return presenters.SuccessResponse[entities.Organizer](ctx, fiber.StatusOK, "successfully update organizer", orRes)
}

func (o *organizerHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	if err := o.s.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "failed to delete organizer", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to delete organizer", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully remove organizer", nil)
}

func NewOrganizerHandler(service organizer.Service) OrganizerHandler {
	return &organizerHandler{s: service}
}
