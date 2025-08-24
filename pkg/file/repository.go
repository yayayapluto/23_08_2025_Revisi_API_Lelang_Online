package file

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
)

type (
	Repository interface {
		Get(ctx context.Context, id uint) (*entities.File, error)
		Save(ctx context.Context, f *entities.File) (*entities.File, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) Get(ctx context.Context, id uint) (*entities.File, error) {
	var file entities.File
	if err := r.db.WithContext(ctx).First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *repository) Save(ctx context.Context, f *entities.File) (*entities.File, error) {
	if err := r.db.WithContext(ctx).Create(f).Error; err != nil {
		return nil, err
	}
	return f, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entities.File{}, id).Error; err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
