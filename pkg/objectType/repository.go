package objectType

import (
	"context"
	"errors"
	"fmt"
	"github.com/yayayapluto/revisi_api_lelang_online/domain"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
	"slices"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ObjectType, int64, error)
		Create(ctx context.Context, ot *entities.ObjectType) error
		Get(ctx context.Context, id uint) (*entities.ObjectType, error)
		GetByName(ctx context.Context, name string) (*entities.ObjectType, error)
		Update(ctx context.Context, ot *entities.ObjectType) (*entities.ObjectType, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ObjectType, int64, error) {
	var ots []entities.ObjectType

	validSortBy := []string{"id", "name", "created_at"}
	defSortBy := validSortBy[0]
	if sortBy != nil {
		if !slices.Contains(validSortBy, *sortBy) {
			return nil, 0, domain.ErrInvalidSortByColumn
		}
		defSortBy = *sortBy
	}

	validSortDir := []string{"asc", "descending", "created_at"}
	defSortDir := validSortDir[0]
	if sortDir != nil {
		if !slices.Contains(validSortDir, *sortDir) {
			return nil, 0, domain.ErrInvalidSortDir
		}
		defSortDir = *sortDir
	}

	orderStr := fmt.Sprintf("%s %s", defSortBy, defSortDir)
	query := r.db.WithContext(ctx).Model(&entities.ObjectType{})

	if search != nil {
		sq := "%" + *search + "%"
		query = query.Where("name LIKE ?", sq)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order(orderStr).Find(&ots).Error; err != nil {
		return nil, 0, err
	}

	return &ots, total, nil
}

func (r *repository) Create(ctx context.Context, ot *entities.ObjectType) error {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.ObjectType{}).Where("name = ?", ot.Name).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return gorm.ErrDuplicatedKey
	}
	if err := r.db.WithContext(ctx).Create(ot).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.ObjectType, error) {
	var ot entities.ObjectType
	if err := r.db.WithContext(ctx).First(&ot, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &ot, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*entities.ObjectType, error) {
	var ot entities.ObjectType
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&ot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &ot, nil
}

func (r *repository) Update(ctx context.Context, ot *entities.ObjectType) (*entities.ObjectType, error) {
	ud := map[string]interface{}{} // ud := update data
	if ot.Name != "" {
		ud["Name"] = ot.Name

		var count int64
		if err := r.db.WithContext(ctx).Model(&entities.ObjectType{}).Where("name = ? AND id != ?", ot.Name, ot.ID).Count(&count).Error; err != nil {
			return nil, err
		}

		if count > 0 {
			return nil, gorm.ErrDuplicatedKey
		}
	}

	if err := r.db.WithContext(ctx).Model(&entities.ObjectType{}).Where("id = ?", ot.ID).Updates(ud).Error; err != nil {
		return nil, err
	}

	var result entities.ObjectType
	if err := r.db.WithContext(ctx).First(&ot, ot.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ObjectType{})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
