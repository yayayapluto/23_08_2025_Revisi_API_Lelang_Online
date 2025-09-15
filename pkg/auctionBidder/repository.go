package auctionBidder

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.AuctionBidder, int64, error)
		Create(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error)
		Get(ctx context.Context, id uint) (*entities.AuctionBidder, error)
		Update(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, sortDir, sortBy *string) (*[]entities.AuctionBidder, int64, error) {
	validSortBy := []string{"id", "created_at", "user_id", "auction_id"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).Preload("User").Preload("Auction")

	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.AuctionBidder
	if err := query.Offset(offset).Limit(limit).Order(*orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil
}

func (r *repository) Create(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).
		Where("user_id = ? and auction_id = ?", e.UserID, e.AuctionID).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("auction bidder already exists")
	}

	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return nil, err
	}

	var result entities.AuctionBidder
	if err := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).
		Preload("User").Preload("Auction").
		First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.AuctionBidder, error) {
	var result entities.AuctionBidder
	if err := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).
		Preload("User").Preload("Auction").
		First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, e *entities.AuctionBidder) (*entities.AuctionBidder, error) {
	updateData := map[string]interface{}{}
	if e.UserID != 0 {
		updateData["user_id"] = e.UserID
	}
	if e.AuctionID != 0 {
		updateData["auction_id"] = e.AuctionID
	}
	if e.PhoneNumber != "" {
		updateData["phone_number"] = e.PhoneNumber
	}
	if e.BankName != "" {
		updateData["bank_name"] = e.BankName
	}
	if e.AccountNumber != "" {
		updateData["account_number"] = e.AccountNumber
	}
	if e.AccountName != "" {
		updateData["account_name"] = e.AccountName
	}

	if err := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).
		Where("id = ?", e.ID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.AuctionBidder
	if err := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).
		Preload("User").Preload("Auction").
		First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Model(&entities.AuctionBidder{}).Where("id = ?", id).Delete(&entities.AuctionBidder{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("auction bidder doesn't exist")
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
