package file

import (
	"context"
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"io"
	"os"
	"time"
)

type (
	Service interface {
		Save(ctx context.Context, fileName *string, src io.Reader) (*entities.File, error)
		Delete(ctx context.Context, id uint) error
	}

	service struct {
		repo Repository
	}
)

func (s *service) Save(ctx context.Context, fileName *string, src io.Reader) (*entities.File, error) {
	dir := fmt.Sprintf("./public/uploads/%s", time.Now().Format(time.DateOnly))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%d_%s", dir, time.Now().Unix(), *fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			return
		}
	}(dst)

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	f := &entities.File{Path: filePath}
	return s.repo.Save(ctx, f)
}

func (s *service) Delete(ctx context.Context, id uint) error {
	file, err := s.repo.Get(ctx, id)
	if err != nil {
		return err
	}

	if err := os.Remove(file.Path); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
