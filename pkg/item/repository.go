package item

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Item, int64, error)
		Create(ctx context.Context, ot *entities.Item) (*entities.Item, error)
		Get(ctx context.Context, id uint) (*entities.Item, error)
		GetByName(ctx context.Context, name string) (*entities.Item, error)
		Update(ctx context.Context, ot *entities.Item) (*entities.Item, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.Item, int64, error) {
	var items []entities.Item

	validSortBy := []string{"id", "name", "created_at"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.Item{}).Preload("ObjectType").Preload("File")
	if search != nil {
		query = query.Where("name LIKE ?", "%"+*search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(limit).Order(orderStr).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	return &items, total, nil
}

func (r *repository) Create(ctx context.Context, item *entities.Item) (*entities.Item, error) {
	if item.ObjectTypeID == 0 {
		return nil, errors.New("object_type_id is required")
	}

	if item.Name == "" {
		return nil, errors.New("name is required")
	}

	if item.Price == 0 {
		return nil, errors.New("price is required")
	}

	if item.DepositPrice == 0 {
		return nil, errors.New("deposit_price is required")
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Item{}).Where("name = ?", item.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, gorm.ErrDuplicatedKey
	}
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, err
	}

	var itemRes entities.Item
	if err := r.db.WithContext(ctx).Preload("ObjectType").Preload("File").First(&itemRes, item.ID).Error; err != nil {
		return nil, err
	}

	return &itemRes, nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.Item, error) {
	var item entities.Item
	if err := r.db.WithContext(ctx).Model(&entities.Item{}).Preload("ObjectType").Preload("File").First(&item, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &item, nil
}

func (r *repository) GetByName(ctx context.Context, name string) (*entities.Item, error) {
	var item entities.Item
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &item, nil
}

func (r *repository) Update(ctx context.Context, item *entities.Item) (*entities.Item, error) {
	updateData := map[string]interface{}{}
	if item.Name != "" {
		updateData["name"] = item.Name

		var count int64
		if err := r.db.WithContext(ctx).Model(&entities.Item{}).Where("name = ? and id != ?", item.Name, item.ID).Count(&count).Error; err != nil {
			return nil, err
		}

		if count != 0 {
			return nil, gorm.ErrDuplicatedKey
		}
	}

	if len(updateData) == 0 {
		return item, nil
	}

	if err := r.db.WithContext(ctx).Model(&entities.Item{}).Where("id = ?", item.ID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.Item
	if err := r.db.WithContext(ctx).Model(&entities.Item{}).Preload("ObjectType").Preload("File").First(&result, item.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&entities.Item{}).Where("id = ?", id).Delete(&entities.Item{})
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
