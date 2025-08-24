package item

import (
	"context"
	"errors"
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/pkg/file"
	"mime/multipart"
	"strings"
)

type (
	Service interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Item, int64, error)
		Create(ctx context.Context, ot *entities.Item, thumbnail *multipart.FileHeader) (*entities.Item, error)
		Get(ctx context.Context, id uint) (*entities.Item, error)
		Update(ctx context.Context, ot *entities.Item) (*entities.Item, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repo Repository
		fs   file.Service
	}
)

func (s *service) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Item, int64, error) {
	return s.repo.List(ctx, offset, limit, search, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, item *entities.Item, thumbnail *multipart.FileHeader) (*entities.Item, error) {
	f, err := thumbnail.Open()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to upload thumbnail: %v\n", err))
	}
	defer func(f multipart.File) {
		err := f.Close()
		if err != nil {
			return
		}
	}(f)

	fileName := strings.ReplaceAll(thumbnail.Filename, " ", "_")

	fe, err := s.fs.Save(ctx, &fileName, f) // fe => file entity
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to upload thumbnail: %v\n", err))
	}

	item.FileID = fe.ID
	return s.repo.Create(ctx, item)
}

func (s *service) Get(ctx context.Context, id uint) (*entities.Item, error) {
	return s.repo.Get(ctx, id)
}

func (s *service) Update(ctx context.Context, ot *entities.Item) (*entities.Item, error) {
	return s.repo.Update(ctx, ot)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

func NewService(repo Repository, fileService file.Service) Service {
	return &service{repo: repo, fs: fileService}
}
