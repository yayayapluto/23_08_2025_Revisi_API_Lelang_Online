package itemDetail

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ItemDetail, int64, error)
		Create(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error)
		Get(ctx context.Context, itemID uint) (*entities.ItemDetail, error)
		Update(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error)
		Delete(ctx context.Context, itemID uint) error
	}

	service struct {
		repo Repository
	}
)

func (s *service) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ItemDetail, int64, error) {
	return s.repo.List(ctx, offset, limit, search, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error) {
	return s.repo.Create(ctx, e)
}

func (s *service) Get(ctx context.Context, itemID uint) (*entities.ItemDetail, error) {
	return s.repo.Get(ctx, itemID)
}

func (s *service) Update(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error) {
	return s.repo.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, itemID uint) error {
	return s.repo.Delete(ctx, itemID)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
