package itemThumbnail

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
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.ItemThumbnail, int64, error)
		Create(ctx context.Context, e *entities.ItemThumbnail, thumbnail *multipart.FileHeader) (*entities.ItemThumbnail, error)
		Get(ctx context.Context, itemID uint) (*entities.ItemThumbnail, error)
		Update(ctx context.Context, e *entities.ItemThumbnail) (*entities.ItemThumbnail, error)
		Delete(ctx context.Context, itemID uint) error
	}

	service struct {
		repo Repository
		fs   file.Service
	}
)

func (s *service) List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.ItemThumbnail, int64, error) {
	return s.repo.List(ctx, offset, limit, sortDir, sortBy)
}

func (s *service) Create(ctx context.Context, e *entities.ItemThumbnail, thumbnail *multipart.FileHeader) (*entities.ItemThumbnail, error) {
	fileUpload, err := thumbnail.Open()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to open thumbnail: %v\n", err))
	}
	defer func(fileUpload multipart.File) {
		err := fileUpload.Close()
		if err != nil {
			return
		}
	}(fileUpload)

	fileName := strings.ReplaceAll(thumbnail.Filename, " ", "_")
	res, err := s.fs.Save(ctx, &fileName, fileUpload)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to save thumbnail: %v\n", err))
	}

	e.FileID = res.ID
	return s.repo.Create(ctx, e)
}

func (s *service) Get(ctx context.Context, itemID uint) (*entities.ItemThumbnail, error) {
	return s.repo.Get(ctx, itemID)
}

func (s *service) Update(ctx context.Context, e *entities.ItemThumbnail) (*entities.ItemThumbnail, error) {
	return s.repo.Update(ctx, e)
}

func (s *service) Delete(ctx context.Context, itemID uint) error {
	return s.repo.Delete(ctx, itemID)
}

func NewService(repo Repository, fs file.Service) Service {
	return &service{repo: repo, fs: fs}
}
