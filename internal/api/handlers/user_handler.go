package handlers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/api/presenters"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/user"
	"gorm.io/gorm"
	"strings"
)

type (
	UserHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
		Login(ctx *fiber.Ctx) error
		Me(ctx *fiber.Ctx) error
	}

	userHandler struct {
		s user.Service
	}
)

func (u *userHandler) List(ctx *fiber.Ctx) error {
	meta := utils.GetRequestMeta(ctx)
	list, total, err := u.s.List(ctx.UserContext(), meta.Offset, meta.Size, &meta.Search, &meta.SortDir, &meta.SortBy)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to list user", err)
	}
	pagination := utils.BuildPagination(ctx, meta, *list, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully list users", &pagination)
}

func (u *userHandler) Create(ctx *fiber.Ctx) error {
	var req entities.User
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body", err)
	}

	err := u.s.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to create user", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully create user", nil)
}

func (u *userHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid id", err)
	}

	userRes, err := u.s.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to get user", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get user", userRes)
}

func (u *userHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid id", err)
	}

	var updateData domain.UpdateRequestUser
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body", err)
	}

	data := new(entities.User)
	data.ID = uint(id)
	if updateData.Username != nil {
		data.Username = *updateData.Username
	}

	err = u.s.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to update user", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully update user", nil)
}

func (u *userHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Invalid id", err)
	}

	if err := u.s.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "User not found", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to delete user", err)
	}

	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "successfully delete user", nil)
}

func (u *userHandler) Login(ctx *fiber.Ctx) error {
	var loginInput domain.LoginInput
	if err := ctx.BodyParser(&loginInput); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "failed to parse body", err)
	}

	token, err := u.s.Login(ctx.UserContext(), loginInput.Identity, loginInput.Password)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to login", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully login", token)
}

func (u *userHandler) Me(ctx *fiber.Ctx) error {
	bearerToken := ctx.Get("Authorization")
	if bearerToken == "" {
		return presenters.ErrorResponse(ctx, fiber.StatusUnauthorized, "missing authorization header", nil)
	}

	// hapus prefix "Bearer "
	tokenStr := strings.TrimPrefix(bearerToken, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)

	userRes, err := u.s.GetUserFromToken(ctx.UserContext(), tokenStr)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusUnauthorized, "failed to get user", err)
	}

	return presenters.SuccessResponse(ctx, fiber.StatusOK, "successfully get user", userRes)
}

func NewUserHandler(s user.Service) UserHandler {
	return &userHandler{s: s}
}
