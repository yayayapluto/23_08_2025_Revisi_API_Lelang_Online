package organizer

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Organizer, int64, error)
		Create(ctx context.Context, ot *entities.Organizer) error
		Get(ctx context.Context, id uint) (*entities.Organizer, error)
		Update(ctx context.Context, ot *entities.Organizer) (*entities.Organizer, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repo Repository
	}
)

func (s *service) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Organizer, int64, error) {
	return s.repo.List(ctx, offset, limit, search, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, ot *entities.Organizer) error {
	return s.repo.Create(ctx, ot)
}

func (s *service) Get(ctx context.Context, id uint) (*entities.Organizer, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, ot *entities.Organizer) (*entities.Organizer, error) {
	return s.repo.Update(ctx, ot)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
