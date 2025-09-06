package pic

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.PIC, int64, error)
		Create(ctx context.Context, e *entities.PIC) (*entities.PIC, error)
		Get(ctx context.Context, id uint) (*entities.PIC, error)
		Update(ctx context.Context, e *entities.PIC) (*entities.PIC, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.PIC, int64, error) {
	validSortBy := []string{"id", "name", "created_at"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.PIC{})
	if search != nil {
		query = query.Where("name LIKE ? or phone_number LIKE ?", "%"+*search+"%", "%"+*search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.PIC
	if err := query.Offset(offset).Limit(limit).Order(*orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil
}

func (r *repository) Create(ctx context.Context, e *entities.PIC) (*entities.PIC, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.PIC{}).Where("name = ? or phone_number = ?", e.Name, e.PhoneNumber).Count(&count).Error; err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, gorm.ErrDuplicatedKey
	}

	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return nil, err
	}

	var result entities.PIC
	if err := r.db.WithContext(ctx).First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.PIC, error) {
	var result entities.PIC
	if err := r.db.WithContext(ctx).Model(&entities.PIC{}).First(&result, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, e *entities.PIC) (*entities.PIC, error) {
	updateData := map[string]interface{}{}
	if e.Name != "" {
		updateData["name"] = e.Name
	}
	if e.PhoneNumber != "" {
		updateData["phone_number"] = e.PhoneNumber
	}

	if err := r.db.WithContext(ctx).Model(&entities.PIC{}).Where("id = ?", e.ID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.PIC
	if err := r.db.WithContext(ctx).First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&entities.PIC{}).Where("id = ?", id).Delete(&entities.PIC{})
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
