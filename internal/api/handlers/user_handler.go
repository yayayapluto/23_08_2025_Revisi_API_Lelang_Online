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
	"log"
)

type (
	UserHandler interface {
		List(ctx *fiber.Ctx) error
		Create(ctx *fiber.Ctx) error
		Get(ctx *fiber.Ctx) error
		Update(ctx *fiber.Ctx) error
		Delete(ctx *fiber.Ctx) error
		CreateOrganizerUser(ctx *fiber.Ctx) error
		Register(ctx *fiber.Ctx) error
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
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal mendapatkan daftar user", err)
	}
	pagination := utils.BuildPagination(ctx, meta, *list, total)
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Berhasil mendapatkan daftar user", &pagination)
}

func (u *userHandler) Create(ctx *fiber.Ctx) error {
	var req entities.User
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal parsing body", err)
	}
	err := u.s.Create(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal membuat user", err)
	}
	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Berhasil membuat user", nil)
}

func (u *userHandler) Get(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "ID tidak valid", err)
	}
	userRes, err := u.s.Get(ctx.UserContext(), uint(id))
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal mendapatkan user", err)
	}
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Berhasil mendapatkan user", userRes)
}

func (u *userHandler) Update(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "ID tidak valid", err)
	}
	var updateData domain.UpdateRequestUser
	if err := ctx.BodyParser(&updateData); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal parsing body", err)
	}
	data := new(entities.User)
	data.ID = uint(id)
	if updateData.Username != nil {
		data.Username = *updateData.Username
	}
	err = u.s.Update(ctx.UserContext(), data)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal memperbarui user", err)
	}
	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Berhasil memperbarui user", nil)
}

func (u *userHandler) Delete(ctx *fiber.Ctx) error {
	id, err := ctx.ParamsInt("id")
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "ID tidak valid", err)
	}
	if err := u.s.Delete(ctx.UserContext(), uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return presenters.ErrorResponse(ctx, fiber.StatusNotFound, "User tidak ditemukan", err)
		}
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal menghapus user", err)
	}
	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Berhasil menghapus user", nil)
}

func (u *userHandler) CreateOrganizerUser(ctx *fiber.Ctx) error {
	var req domain.CreateOrganizerUserInput
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal parsing body", err)
	}
	if err := u.s.CreateOrganizerUser(ctx.UserContext(), &req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal membuat user organizer", err)
	}
	return presenters.SuccessResponse[any](ctx, fiber.StatusCreated, "Berhasil membuat user organizer", nil)
}

func (u *userHandler) Register(ctx *fiber.Ctx) error {
	var req domain.RegisterInput
	if err := ctx.BodyParser(&req); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal parsing body", err)
	}
	if req.Password != req.ConfirmPassword {
		log.Println(req.Password, req.ConfirmPassword)
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Password dan konfirmasi password tidak sama", nil)
	}
	err := u.s.Register(ctx.UserContext(), &req)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return presenters.SuccessResponse[any](ctx, fiber.StatusOK, "Berhasil registrasi user", nil)
}

func (u *userHandler) Login(ctx *fiber.Ctx) error {
	var loginInput domain.LoginInput
	if err := ctx.BodyParser(&loginInput); err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Gagal parsing body", err)
	}
	token, err := u.s.Login(ctx.UserContext(), loginInput.Identity, loginInput.Password)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return presenters.ErrorResponse(ctx, fiber.StatusBadRequest, "Email / Username tidak ditemukan", err)
	}
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, err.Error(), nil)
	}
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Berhasil login", token)
}

func (u *userHandler) Me(ctx *fiber.Ctx) error {
	userRes, err := utils.GetUserFromToken(ctx)
	if err != nil {
		return presenters.ErrorResponse(ctx, fiber.StatusInternalServerError, "Gagal mendapatkan user dari token", err)
	}
	return presenters.SuccessResponse(ctx, fiber.StatusOK, "Berhasil mendapatkan user", userRes)
}

func NewUserHandler(s user.Service) UserHandler {
	return &userHandler{s: s}
}
