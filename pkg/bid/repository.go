package bid

import (
	"context"
	"errors"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
	"log"
)

type (
	Repository interface {
		List(ctx context.Context, offset, limit int, sortDir, sortBy *string, auctionID *uint) (*[]entities.Bid, int64, error)
		Create(ctx context.Context, e *entities.Bid) (*entities.Bid, error)
		Get(ctx context.Context, id uint) (*entities.Bid, error)
		Update(ctx context.Context, e *entities.Bid) (*entities.Bid, error)
		Delete(ctx context.Context, id uint) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) List(ctx context.Context, offset, limit int, sortDir, sortBy *string, auctionID *uint) (*[]entities.Bid, int64, error) {
	validSortBy := []string{"id", "created_at", "value"}
	validSortDir := []string{"asc", "desc"}
	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).Model(&entities.Bid{}).Preload("Bidder.User")

	log.Println("[repo] auction_id", *auctionID)
	log.Printf("[repo] query where %v", &entities.Bid{Bidder: entities.AuctionBidder{AuctionID: *auctionID}})
	if auctionID != nil {
		query = query.
			Joins("JOIN auction_bidders ON auction_bidders.id = bids.bidder_id").
			Where("auction_bidders.auction_id = ?", *auctionID)
	}

	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.Bid
	if err := query.Offset(offset).Limit(limit).Order(*orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil

	/*
		2025/10/23 16:10:23 [handler] auction_id 9
		2025/10/23 16:10:23 [service] auction_id 9
		2025/10/23 16:10:23 [repo] auction_id 9
		2025/10/23 16:10:23 [repo] query where &{0 0 0 {0 0 9     {0     <nil> <nil> <nil> {0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}} <nil> <nil> <nil> {0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}} {0001-01-01 00:00:00 +0000 UTC 0001-01-01 00:00:00 +0000 UTC}}
	*/
}

func (r *repository) Create(ctx context.Context, e *entities.Bid) (*entities.Bid, error) {
	sortDir := "desc"
	sortBy := "value"
	highestBid, _, err := r.List(ctx, 0, 1, &sortDir, &sortBy, nil)
	if err != nil || highestBid == nil {
		return nil, err
	}

	log.Printf("highest bid is: %v", highestBid)

	if e.Value <= (*highestBid)[0].Value {
		return nil, errors.New("bid value cannot be the same or lower than highest bid")
	}

	if err := r.db.WithContext(ctx).Create(e).Error; err != nil {
		return nil, err
	}

	var result entities.Bid
	if err := r.db.WithContext(ctx).
		Model(&entities.Bid{}).
		Preload("Bidder.User").
		First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Get(ctx context.Context, id uint) (*entities.Bid, error) {
	var result entities.Bid
	if err := r.db.WithContext(ctx).
		Model(&entities.Bid{}).
		Preload("Bidder.User").
		First(&result, id).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(ctx context.Context, e *entities.Bid) (*entities.Bid, error) {
	updateData := map[string]interface{}{}
	if e.BidderID != 0 {
		updateData["bidder_id"] = e.BidderID
	}
	if e.Value != 0 {
		updateData["value"] = e.Value
	}

	if err := r.db.WithContext(ctx).
		Model(&entities.Bid{}).
		Where("id = ?", e.ID).
		Updates(updateData).Error; err != nil {
		return nil, err
	}

	var result entities.Bid
	if err := r.db.WithContext(ctx).
		Model(&entities.Bid{}).
		Preload("Bidder.User").
		First(&result, e.ID).Error; err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Model(&entities.Bid{}).
		Where("id = ?", id).
		Delete(&entities.Bid{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("bid doesn't exist")
	}
	return nil
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}
