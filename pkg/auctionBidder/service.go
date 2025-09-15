package auctionBidder

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.AuctionBidder, int64, error)
		Create(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error)
		Get(ctx context.Context, id uint) (*entities.AuctionBidder, error)
		Update(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repo Repository
	}
)

func (s service) List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.AuctionBidder, int64, error) {
	return s.repo.List(ctx, offset, limit, sortDir, sortBy)
}

func (s service) Create(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error) {
	return s.repo.Create(ctx, e)
}

func (s service) Get(ctx context.Context, id uint) (*entities.AuctionBidder, error) {
	return s.repo.Get(ctx, id)
}

func (s service) Update(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error) {
	return s.repo.Update(ctx, e)
}

func (s service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
