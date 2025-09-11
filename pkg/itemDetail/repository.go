package itemDetail

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ItemDetail, int64, error)
		Create(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error)
		Get(ctx context.Context, itemID uint) (*entities.ItemDetail, error)
		Update(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error)
		Delete(ctx context.Context, itemID uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, search, sortDir, sortBy *string) (*[]entities.ItemDetail, int64, error) {
	validSortBy := []string{"id", "brand", "year", "color", "created_at"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Preload("Item.ObjectType").Preload("Item.File")
	if search != nil && *search != "" {
		query = query.Where("brand LIKE ? or year LIKE ? or color LIKE ?", "%"+*search+"%", "%"+*search+"%", "%"+*search+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var ItemDetails []entities.ItemDetail
	if err := query.Offset(offset).Limit(limit).Order(orderStr).Find(&ItemDetails).Error; err != nil {
		return nil, 0, err
	}
	return &ItemDetails, total, nil
}

func (r *repository) Create(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error) {
	if e.ItemID == 0 {
		return nil, errors.New("item_id is required")
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Where("item_id = ?", e.ItemID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count != 0 {
		return nil, gorm.ErrDuplicatedKey
	}

	if e.Brand == "" {
		return nil, errors.New("brand is required")
	}
	if e.Year == 0 {
		return nil, errors.New("year is required")
	}
	if e.Color == "" {
		return nil, errors.New("color is required")
	}

	if err := r.db.WithContext(ctx).Create(&e).Error; err != nil {
		return nil, err
	}

	var itemDetailRes entities.ItemDetail
	if err := r.db.WithContext(ctx).First(&itemDetailRes, e.ID).Error; err != nil {
		return nil, err
	}

	return &itemDetailRes, nil
}

func (r *repository) Get(ctx context.Context, itemID uint) (*entities.ItemDetail, error) {
	var itemDetail entities.ItemDetail
	if err := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Where("item_id = ?", itemID).First(&itemDetail).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
	}
	return &itemDetail, nil
}

func (r *repository) Update(ctx context.Context, e *entities.ItemDetail) (*entities.ItemDetail, error) {
	updateData := map[string]interface{}{}

	if e.PlateNumber != nil {
		updateData["plate_number"] = *e.PlateNumber
	}
	if e.Brand != "" {
		updateData["brand"] = e.Brand
	}
	if e.Series != nil {
		updateData["series"] = *e.Series
	}
	if e.CC != nil {
		updateData["cc"] = *e.CC
	}
	if e.Type != nil {
		updateData["type"] = *e.Type
	}
	if e.Transmission != nil {
		updateData["transmission"] = *e.Transmission
	}
	if e.Model != nil {
		updateData["model"] = *e.Model
	}
	if e.Year != 0 {
		updateData["year"] = e.Year
	}
	if e.FrameNumber != nil {
		updateData["frame_number"] = *e.FrameNumber
	}
	if e.MachineNumber != nil {
		updateData["machine_number"] = *e.MachineNumber
	}
	if e.Kilometer != nil {
		updateData["kilometer"] = *e.Kilometer
	}
	if e.Fuel != nil {
		updateData["fuel"] = *e.Fuel
	}
	if e.Color != "" {
		updateData["color"] = e.Color
	}
	if e.DriveType != nil {
		updateData["drive_type"] = *e.DriveType
	}
	if e.StnkDate != nil {
		updateData["stnk_date"] = *e.StnkDate
	}

	if err := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Where("item_id = ?", e.ItemID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.ItemDetail
	if err := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Preload("Item.ObjectType").Preload("Item.File").Where("item_id = ?", e.ItemID).First(&result).Error; err != nil {
		return nil, err
	}

	return e, nil
}

func (r *repository) Delete(ctx context.Context, itemID uint) error {
	res := r.db.WithContext(ctx).Model(&entities.ItemDetail{}).Where("item_id = ?", itemID).Delete(&entities.ItemDetail{})
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
