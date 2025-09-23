package bidderPayment

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"gorm.io/gorm"
)

type (
	Repository interface {
		FindByOrderId(ctx context.Context, orderID string) (*entities.BidderPayment, error)
		FindByBidderData(ctx context.Context, auctionID, userID uint) (*entities.BidderPayment, error)
		Insert(ctx context.Context, e *entities.BidderPayment) error
		Update(ctx context.Context, e *entities.BidderPayment) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) FindByBidderData(ctx context.Context, auctionID, userID uint) (*entities.BidderPayment, error) {
	var result entities.BidderPayment
	err := r.db.WithContext(ctx).
		Table("bidder_payments").
		Joins("JOIN auction_bidders ON auction_bidders.id = bidder_payments.bidder_id").
		Where("auction_bidders.auction_id = ? AND auction_bidders.user_id = ?", auctionID, userID).
		First(&result).Error

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) FindByOrderId(ctx context.Context, orderID string) (*entities.BidderPayment, error) {
	var result entities.BidderPayment
	if err := r.db.WithContext(ctx).First(&result, "order_id = ?", orderID).Error; err != nil {
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
