package bidderPayment

import (
	"context"
	"github.com/yayayapluto/revisi_api_lelang_online/entities"
	"github.com/yayayapluto/revisi_api_lelang_online/internal/utils"
	"gorm.io/gorm"
	"log"
)

type (
	Repository interface {
		FindByOrderId(ctx context.Context, orderID string) (*entities.BidderPayment, error)
		FindByUserId(ctx context.Context, userID uint, offset, limit int, sortDir, sortBy *string, status, paymentType, itemName *string) (*[]entities.BidderPayment, int64, error)
		FindByBidderData(ctx context.Context, auctionID, userID uint) (*entities.BidderPayment, error)
		Insert(ctx context.Context, e *entities.BidderPayment) error
		Update(ctx context.Context, e *entities.BidderPayment) error
	}

	repository struct {
		db *gorm.DB
	}
)

func (r *repository) FindByUserId(
	ctx context.Context,
	userID uint,
	offset, limit int,
	sortDir, sortBy *string,
	status, paymentType, itemName *string,
) (*[]entities.BidderPayment, int64, error) {
	validSortBy := []string{"id", "created_at", "amount"}
	validSortDir := []string{"asc", "desc"}

	orderStr, err := utils.BuildOrderQuery(validSortBy, validSortDir, sortBy, sortDir)
	if err != nil {
		return nil, 0, err
	}

	query := r.db.WithContext(ctx).
		Model(&entities.BidderPayment{}).
		Preload("Bidder.Auction.Item.ObjectType").
		Preload("Bidder.Auction.Item.File").
		Preload("Bidder.Auction.Organizer").
		Joins("JOIN auction_bidders ON auction_bidders.id = bidder_payments.bidder_id").
		Joins("JOIN auctions ON auctions.id = auction_bidders.auction_id").
		Joins("JOIN items ON items.id = auctions.item_id").
		Where("auction_bidders.user_id = ?", userID)

	if status != nil && len(*status) > 0 {
		query = query.Where("bidder_payments.status = ?", *status)
	}

	if paymentType != nil && len(*paymentType) > 0 {
		query = query.Where("bidder_payments.type = ?", *paymentType)
	}

	if itemName != nil && len(*itemName) > 0 {
		query = query.Where("items.name ILIKE ?", "%"+*itemName+"%")
	}

	log.Println(query.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&entities.BidderPayment{})
	}))

	var total int64
	if err = query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var collection []entities.BidderPayment
	if err := query.Offset(offset).Limit(limit).Order(*orderStr).Find(&collection).Error; err != nil {
		return nil, 0, err
	}

	return &collection, total, nil
}

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
	if err := r.db.WithContext(ctx).Preload("Bidder.Auction.Item.ObjectType").
		Preload("Bidder.Auction.Item.File").
		Preload("Bidder.Auction.Organizer").Preload("Bidder.Auction.PIC").First(&result, "order_id = ?", orderID).Error; err != nil {
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
