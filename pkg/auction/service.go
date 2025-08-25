package auction

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"time"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, startDate, endDate *time.Time, sortDir, sortBy *string) (*[]entities.Auction, int64, error)
		Create(ctx context.Context, e *entities.Auction) (*entities.Auction, error)
		Get(ctx context.Context, id uint) (*entities.Auction, error)
		Update(ctx context.Context, e *entities.Auction) (*entities.Auction, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repo Repository
	}
)

func (s *service) List(ctx context.Context, offset, limit int, startDate, endDate *time.Time, sortDir, sortBy *string) (*[]entities.Auction, int64, error) {
	return s.repo.List(ctx, offset, limit, startDate, endDate, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, e *entities.Auction) (*entities.Auction, error) {
	return s.repo.Create(ctx, e)
}

func (s *service) Get(ctx context.Context, id uint) (*entities.Auction, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, e *entities.Auction) (*entities.Auction, error) {
	return s.repo.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
