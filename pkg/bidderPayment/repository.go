package bidderPayment

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
)

type (
	Repository interface {
		FindById(ctx context.Context, id string) (*entities.BidderPayment, error)
		Insert(ctx context.Context, e *entities.BidderPayment) error
		Update(ctx context.Context, e *entities.BidderPayment) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) FindById(ctx context.Context, id string) (*entities.BidderPayment, error) {
	var result entities.BidderPayment
	if err := r.db.WithContext(ctx).First(&result, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Insert(ctx context.Context, e *entities.BidderPayment) error {
	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) Update(ctx context.Context, e *entities.BidderPayment) error {
	if err := r.db.WithContext(ctx).Model(&entities.BidderPayment{}).Where("id = ?", e.ID).Updates(e).Error; err != nil {
		return err
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
