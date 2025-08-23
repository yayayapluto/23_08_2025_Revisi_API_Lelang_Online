package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	presenters "github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils/pagination"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/objectType"
	"gorm.io/gorm"
	"math"
)

type (
	ObjectTypeHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
	}

	objectTypeHandler struct {
		s objectType.Service
	}
)

func (o *objectTypeHandler) List(ctx *fiber.Ctx) error {
	search := ctx.Query("search")
	page := ctx.QueryInt("page", 1)

	size := ctx.QueryInt("size")
	size = int(math.Min(math.Max(float64(size), 10), 100)) // min 10, max 100

	offset := (page - 1) * size // ex: page=2 size=10 -> (2 - 1) * 10 -> 10 | that means it start from offset 10

	sortBy := ctx.Query("sortBy", "id")
	sortDir := ctx.Query("sortDir", "asc")

	OTs, total, err := o.s.List(ctx.UserContext(), offset, size, &search, &sortDir, &sortBy) // OTs => object types | plural
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to retrieve object type list", err)
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

	paginationRes := pagination.NewResponseMetaData[entities.ObjectType](page, currentPageUrl, *OTs, firstPageUrl, nextPageUrl, size, prevPageUrl)
	return presenters.SuccessResponse[pagination.ResponseMetaData[entities.ObjectType]](ctx, fiber.StatusOK, "successfully retrieve object type list", &paginationRes)
}

func (o *objectTypeHandler) Create(ctx *fiber.Ctx) error {
	var ot entities.ObjectType
	if err := ctx.BodyParser(&ot); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	if err := o.s.Create(ctx.UserContext(), &ot); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to create new object type", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to create new object type", err)
	}

	return presenters.SuccessResponse[entities.ObjectType](ctx, fiber.StatusOK, "successfully create new object type", &ot)
}

func (o *objectTypeHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	ot, err := o.s.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to get object type detail", err)
	}

	return presenters.SuccessResponse[entities.ObjectType](ctx, fiber.StatusOK, "successfully get object type detail", ot)
}

func (o *objectTypeHandler) Update(ctx *fiber.Ctx) error {
	var ot entities.ObjectType

	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	ot.ID = uint(id)

	if err := ctx.BodyParser(&ot); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body request", err)
	}

	otRes, err := o.s.Update(ctx.UserContext(), &ot)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to update object type", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to update object type", err)
	}
	otRes.ID = uint(id)

	return presenters.SuccessResponse[entities.ObjectType](ctx, fiber.StatusOK, "successfully update object type", otRes)
}

func (o *objectTypeHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "invalid id param", err)
	}

	if err := o.s.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "failed to delete object type", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "failed to delete object type", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully remove object type", nil)
}

func NewObjectTypeHandler(service objectType.Service) ObjectTypeHandler {
	return &objectTypeHandler{s: service}
}
