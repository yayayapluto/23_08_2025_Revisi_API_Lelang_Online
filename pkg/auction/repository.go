package auction

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
	"time"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, startDate, endDate *time.Time, sortDir, sortBy *string) (*[]entities.Auction, int64, error)
		Create(ctx context.Context, e *entities.Auction) (*entities.Auction, error)
		Get(ctx context.Context, id uint) (*entities.Auction, error)
		Update(ctx context.Context, e *entities.Auction) (*entities.Auction, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, startDate, endDate *time.Time, sortDir, sortBy *string) (*[]entities.Auction, int64, error) {
	validSortBy := []string{"id", "start_date", "end_date", "created_at"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.Auction{}).Preload("Item.ObjectType").Preload("Item.File").Preload("Organizer").Preload("PIC")
	if startDate != nil {
		query = query.Where("start_date >= ?", startDate)
	}
	if endDate != nil {
		query = query.Where("end_date <= ?", endDate)
	}

	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.Auction
	if err := query.Offset(offset).Limit(limit).Order(*orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil
}

func (r *repository) Create(ctx context.Context, e *entities.Auction) (*entities.Auction, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.Auction{}).Where("item_id = ? and pic_id = ? and organizer_id = ? and start_date = ? and end_date = ?", e.ItemID, e.PicID, e.OrganizerID, e.StartDate, e.EndDate).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("auction already exists")
	}

	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return nil, err
	}

	var result entities.Auction
	if err := r.db.WithContext(ctx).Model(&entities.Auction{}).Preload("Item.ObjectType").Preload("Item.File").Preload("Item.ItemDetail").Preload("Item.ItemDocument").Preload("Item.ItemGrade").Preload("Item.ItemThumbnails").Preload("Organizer").Preload("PIC").First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.Auction, error) {
	var result entities.Auction
	if err := r.db.WithContext(ctx).Model(&entities.Auction{}).Preload("Item.ObjectType").Preload("Item.File").Preload("Item.ItemDetail").Preload("Item.ItemDocument").Preload("Item.ItemGrade").Preload("Item.ItemThumbnails").Preload("Organizer").Preload("PIC").First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, e *entities.Auction) (*entities.Auction, error) {
	updateData := map[string]interface{}{}
	if e.ItemID != 0 {
		updateData["item_id"] = e.ItemID
	}
	if e.PicID != 0 {
		updateData["pic_id"] = e.PicID
	}
	if e.OrganizerID != 0 {
		updateData["organizer_id"] = e.OrganizerID
	}
	updateData["start_date"] = e.StartDate
	updateData["end_date"] = e.EndDate

	if err := r.db.WithContext(ctx).Model(&entities.Auction{}).Preload("Item.ObjectType").Preload("Item.File").Preload("Item.ItemDetail").Preload("Item.ItemDocument").Preload("Item.ItemGrade").Preload("Item.ItemThumbnails").Preload("Organizer").Preload("PIC").Where("id = ?", e.ID).Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.Auction
	if err := r.db.WithContext(ctx).First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&entities.Auction{}).Where("id = ?", id).Delete(&entities.Auction{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("auction doesn't exists")
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
