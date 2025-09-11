package itemGrade

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.ItemGrade, int64, error)
		Create(ctx context.Context, e *entities.ItemGrade) (*entities.ItemGrade, error)
		Get(ctx context.Context, itemID uint) (*entities.ItemGrade, error)
		Update(ctx context.Context, e *entities.ItemGrade) (*entities.ItemGrade, error)
		Delete(ctx context.Context, itemID uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.ItemGrade, int64, error) {
	validSortBy := []string{"id", "interior", "exterior", "frame", "machine", "created_at"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.ItemGrade{})
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.ItemGrade
	if err := query.Offset(offset).Limit(limit).Order(orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil
}

func (r *repository) Create(ctx context.Context, e *entities.ItemGrade) (*entities.ItemGrade, error) {
	if e.ItemID == 0 {
		return nil, errors.New("item_id is required")
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", e.ItemID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, gorm.ErrDuplicatedKey
	}

	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return nil, err
	}

	var result entities.ItemGrade
	if err := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", e.ItemID).Find(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Get(ctx context.Context, itemID uint) (*entities.ItemGrade, error) {
	var result entities.ItemGrade
	if err := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", itemID).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, e *entities.ItemGrade) (*entities.ItemGrade, error) {
	updateData := map[string]interface{}{}

	if e.Interior != "" {
		updateData["interior"] = e.Interior
	}

	if e.Exterior != "" {
		updateData["exterior"] = e.Exterior
	}

	if e.Frame != "" {
		updateData["frame"] = e.Frame
	}

	if e.Machine != "" {
		updateData["machine"] = e.Machine
	}

	if err := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", e.ItemID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.ItemGrade
	if err := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", e.ItemID).First(&result).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, itemID uint) error {
	res := r.db.WithContext(ctx).Model(&entities.ItemGrade{}).Where("item_id = ?", itemID).Delete(&entities.ItemGrade{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected <= 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
